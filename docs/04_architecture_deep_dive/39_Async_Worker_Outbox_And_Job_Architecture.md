# 39 - Async Worker Outbox And Job Architecture

Version: v1.4 Architecture Deep Dive  
Status: Draft for technical alignment  
Scope: Async execution model, outbox, job table, scheduler, worker claim, retry, idempotency, failure recovery  
Related docs: 11, 14, 18, 19, 21, 22, 37, 38, 40, 43, 44

---

## 1. Mục tiêu tài liệu

Tài liệu này khóa kiến trúc async cho nền tảng VPS/Proxy.

Dự án không được xử lý mọi thứ trong HTTP request vì các flow quan trọng có thể chậm, lỗi ngoài hệ thống, hoặc cần retry:

```text
provider provisioning
provider status sync
service action
notification
expiry/suspension/termination
finance reconciliation
credential hygiene check
report export
abuse/takedown workflow
```

Kết luận kiến trúc:

```text
MVP dùng PostgreSQL-backed outbox + jobs table trước.
Queue ngoài có thể thêm sau khi có tải thật.
Business state và outbox/job record phải được tạo trong cùng DB transaction.
Worker phải idempotent.
Provider create không được retry mù khi kết quả không chắc.
```

Câu chốt:

```text
Async không phải để "chạy nền cho nhanh".
Async là cách giữ transaction ngắn, retry có kiểm soát, và recover được sau lỗi.
```

---

## 2. Vì sao cần outbox và jobs

Nếu API checkout làm trực tiếp:

```text
debit wallet
reserve stock
call provider API
send email
return response
```

hệ thống sẽ gặp các lỗi nguy hiểm:

```text
DB commit thành công nhưng provider call timeout.
Provider tạo VPS thành công nhưng API process crash trước khi ghi service.
Queue publish thành công nhưng DB rollback.
Email gửi thành công nhưng order rollback.
HTTP request timeout trong khi trạng thái phía sau vẫn chạy.
```

Outbox giải quyết lỗi DB/queue lệch nhau:

```text
BEGIN
  create order
  debit ledger
  reserve inventory
  create provisioning_job
  insert outbox_events
COMMIT
```

Worker/dispatcher chỉ xử lý thứ đã commit.

---

## 3. Khối runtime async

Trong modular monolith Go:

```text
/cmd/api
  tạo business mutation + outbox/job

/cmd/worker
  claim job
  dispatch outbox
  gọi provider
  gửi notification
  chạy reconciliation async

/cmd/scheduler
  enqueue recurring job
  giữ scheduler lock

/cmd/cli
  internal recovery/manual tool nếu cần
```

Ban đầu không bắt buộc dùng Kafka/RabbitMQ/SQS. PostgreSQL đủ cho MVP nếu:

```text
job volume còn thấp/vừa
claim query có index đúng
worker batch nhỏ
retry/backoff rõ
monitoring queue depth đầy đủ
```

Queue ngoài chỉ thêm khi có lý do thật:

```text
job throughput vượt khả năng Postgres-backed queue
cần fanout lớn
cần consumer độc lập theo team/service
cần regional queue
```

---

## 4. Bảng dữ liệu đề xuất

### 4.1 `outbox_events`

```text
outbox_event_id UUID PK
tenant_id UUID NULL
aggregate_type TEXT NOT NULL
aggregate_id UUID NOT NULL
event_type TEXT NOT NULL
payload_json JSONB NOT NULL DEFAULT '{}'
status TEXT NOT NULL
dedupe_key TEXT NOT NULL
attempt_count INT NOT NULL DEFAULT 0
max_attempts INT NOT NULL DEFAULT 10
next_attempt_at TIMESTAMPTZ NOT NULL DEFAULT now()
locked_by TEXT NULL
locked_until TIMESTAMPTZ NULL
last_error_code TEXT NULL
last_error_message_redacted TEXT NULL
correlation_id TEXT NOT NULL
created_at TIMESTAMPTZ NOT NULL DEFAULT now()
published_at TIMESTAMPTZ NULL
```

Status:

```text
pending
processing
published
failed_retryable
failed_terminal
discarded
```

Constraint:

```text
unique(dedupe_key)
payload_json must be redacted
```

Index:

```text
status, next_attempt_at, created_at
tenant_id, event_type, created_at
correlation_id
```

### 4.2 `jobs`

Một bảng `jobs` tổng quát có thể dùng cho nhiều loại worker, hoặc giữ bảng riêng như `provisioning_jobs` cho domain quan trọng. MVP nên có:

```text
provisioning_jobs riêng vì có logic rủi ro cao
generic jobs cho notification, sync, report, hygiene
```

Generic fields:

```text
job_id UUID PK
tenant_id UUID NULL
job_type TEXT NOT NULL
reference_type TEXT NOT NULL
reference_id UUID NOT NULL
source_id UUID NULL
payload_json JSONB NOT NULL DEFAULT '{}'
status TEXT NOT NULL
priority INT NOT NULL DEFAULT 100
idempotency_key TEXT NOT NULL
attempt_count INT NOT NULL DEFAULT 0
max_attempts INT NOT NULL DEFAULT 5
next_attempt_at TIMESTAMPTZ NOT NULL DEFAULT now()
locked_by TEXT NULL
locked_until TIMESTAMPTZ NULL
last_error_code TEXT NULL
last_error_message_redacted TEXT NULL
manual_review_reason TEXT NULL
correlation_id TEXT NOT NULL
created_at TIMESTAMPTZ NOT NULL DEFAULT now()
updated_at TIMESTAMPTZ NULL
finished_at TIMESTAMPTZ NULL
```

Constraint:

```text
unique(tenant_id, job_type, idempotency_key)
```

### 4.3 `job_attempts`

Lưu lịch sử attempt để debug và audit vận hành:

```text
job_attempt_id UUID PK
job_id UUID NOT NULL
worker_id TEXT NOT NULL
attempt_number INT NOT NULL
started_at TIMESTAMPTZ NOT NULL
finished_at TIMESTAMPTZ NULL
result TEXT NOT NULL
error_code TEXT NULL
error_message_redacted TEXT NULL
duration_ms INT NULL
correlation_id TEXT NOT NULL
```

Không lưu secret, credential, raw provider payload chưa redact.

### 4.4 `processed_events`

Dùng khi consumer cần chống xử lý trùng:

```text
processed_event_id UUID PK
consumer_name TEXT NOT NULL
dedupe_key TEXT NOT NULL
processed_at TIMESTAMPTZ NOT NULL DEFAULT now()
result_summary_json JSONB NOT NULL DEFAULT '{}'
```

Constraint:

```text
unique(consumer_name, dedupe_key)
```

### 4.5 `scheduler_locks`

Dùng để chỉ một scheduler active enqueue recurring job:

```text
lock_name TEXT PK
locked_by TEXT NOT NULL
locked_until TIMESTAMPTZ NOT NULL
heartbeat_at TIMESTAMPTZ NOT NULL
```

Hoặc dùng PostgreSQL advisory lock nếu implementation đơn giản hơn.

---

## 5. Job status lifecycle

Status chuẩn:

```text
queued
claimed
running
succeeded
failed_retryable
failed_terminal
manual_review
cancelled
expired
```

Transition hợp lệ:

```text
queued -> claimed
claimed -> running
running -> succeeded
running -> failed_retryable
running -> failed_terminal
running -> manual_review
failed_retryable -> queued
manual_review -> queued
manual_review -> succeeded
manual_review -> failed_terminal
queued/claimed/running -> cancelled by operator/system
```

Rule:

```text
manual_review không tự quay lại queued nếu không có operator/system action rõ.
failed_terminal không retry.
cancelled job không được worker xử lý tiếp.
```

---

## 6. Worker claim pattern

Worker claim job bằng row lock:

```sql
WITH picked AS (
  SELECT job_id
  FROM jobs
  WHERE status IN ('queued', 'failed_retryable')
    AND next_attempt_at <= now()
  ORDER BY priority ASC, created_at ASC
  FOR UPDATE SKIP LOCKED
  LIMIT $1
)
UPDATE jobs j
SET status = 'claimed',
    locked_by = $2,
    locked_until = now() + interval '2 minutes',
    updated_at = now()
FROM picked
WHERE j.job_id = picked.job_id
RETURNING j.*;
```

Rule:

```text
claim không có nghĩa là job đã chạy xong.
locked_until cho phép reclaim khi worker chết.
worker phải kiểm tra trạng thái/domain precondition trước khi gọi external side effect.
```

Worker ID nên gồm:

```text
hostname
process id
deployment version
random suffix
```

---

## 7. Outbox dispatcher

Outbox dispatcher có 2 chế độ:

```text
internal dispatch: gọi handler trong cùng worker process
external publish: publish ra queue/broker sau này
```

MVP có thể dùng internal dispatch:

```text
outbox_event order.paid -> enqueue provisioning job nếu chưa có
outbox_event service.activated -> enqueue notification
outbox_event wallet.topup.approved -> enqueue notification
```

Nếu sau này thêm queue ngoài:

```text
outbox_event pending -> publish queue -> mark published
consumer queue -> idempotent handler
```

Không được mark `published` trước khi publish thành công.

---

## 8. Idempotency rule theo loại job

### 8.1 Provisioning

Idempotency key phải gắn với:

```text
tenant_id
order_item_id
source_id
operation provision
```

Worker phải kiểm tra:

```text
service cho order_item đã tồn tại chưa
provider_resource_mapping đã tồn tại chưa
provider_request cùng idempotency_key đã có kết quả chưa
reservation còn reserved không
```

Không được gọi provider create lần nữa nếu kết quả provider trước đó là unknown/unsafe.

### 8.2 Notification

Dedupe theo:

```text
tenant_id
recipient_id
event_type
business_reference
channel
```

Ví dụ:

```text
service_expiring:tenant:service:days_before:email
```

Duplicate notification critical có thể gây rối vận hành, nhưng không được chặn audit/security alert quan trọng nếu nó là tín hiệu thật.

### 8.3 Expiry/suspension/termination

Scheduled job chạy lặp phải no-op nếu state đã chuyển:

```text
service already expired -> no-op
service already suspended -> no-op
service terminated -> no-op
```

### 8.4 Reconciliation

Reconciliation không tự sửa tiền/tài nguyên bằng phỏng đoán. Nó tạo:

```text
exception report
risk flag
manual review item
audit event
```

---

## 9. Retry và backoff

Retry classification:

```text
safe_retry
unsafe_retry
do_not_retry
manual_review_required
```

Backoff đề xuất:

```text
attempt 1: immediate or 10 seconds
attempt 2: 1 minute
attempt 3: 5 minutes
attempt 4: 15 minutes
attempt 5: 1 hour
```

Provider 429/rate limit:

```text
respect Retry-After nếu có
limit concurrency per provider/source
do not hammer provider
```

Unknown create result:

```text
manual_review_required
query provider status if safe
do not schedule normal retry create
```

---

## 10. Scheduler architecture

Scheduler tạo recurring job, không xử lý business logic nặng trực tiếp.

Cron jobs phase 1:

```text
reservation_expiry_job
service_expiry_job
suspension_job
termination_job
renewal_reminder_job
provider_health_check_job
inventory_sync_job
finance_reconciliation_job
credential_hygiene_job
notification_retry_job
audit_retention_job
```

Scheduler phải:

```text
giữ lock để tránh enqueue trùng quá mức
generate dedupe_key ổn định
enqueue batch nhỏ
ghi metric enqueue count/error
```

Ngay cả khi scheduler enqueue trùng, job handler vẫn phải idempotent.

---

## 11. Transaction boundaries

Không làm external side effect trong transaction dài.

Đúng:

```text
API transaction:
- debit wallet
- reserve stock
- create order
- create job/outbox
- commit

Worker:
- claim job
- call provider outside DB transaction
- short transaction persist result
```

Sai:

```text
BEGIN
  debit wallet
  call provider API
  send notification
COMMIT
```

Provider success finalization phải là transaction ngắn:

```text
lock provisioning_job
lock reservation
upsert provider mapping
create service
store encrypted credential
allocate reservation
update order/service status
audit/outbox
commit
```

---

## 12. Dead letter và manual review

Không phải mọi lỗi đều là dead letter.

`manual_review` dùng khi:

```text
provider result unknown
partial success
credential missing
state drift
operator decision required
```

`failed_terminal` dùng khi:

```text
provider chắc chắn không tạo resource
source/plan invalid
permission/auth failure không thể tự phục hồi
retry exhausted với safe retry
```

Manual review record phải có:

```text
reason
last safe state
recommended action
links to order/job/provider_request/service
correlation_id
```

---

## 13. Observability

Worker logs phải có:

```text
job_id
job_type
tenant_id
source_id
attempt_count
correlation_id
worker_id
duration_ms
result
error_code
```

Metrics tối thiểu:

```text
jobs_queued_total by job_type
jobs_running_total by job_type
jobs_failed_total by job_type/error_code
job_duration_ms by job_type
outbox_pending_count
outbox_oldest_pending_age_seconds
manual_review_count
provider_retry_count
worker_claim_conflict_count
```

Alert:

```text
oldest pending provisioning job > threshold
manual_review aging > threshold
outbox pending grows continuously
provider auth failures
worker crash loop
```

Chi tiết observability xem file 43.

---

## 14. Migration path sang queue ngoài

Khi cần queue ngoài, không bỏ outbox.

Path:

```text
1. Giữ outbox_events trong DB transaction.
2. Dispatcher publish outbox_event sang broker.
3. Consumer xử lý idempotent theo dedupe_key.
4. processed_events chống xử lý trùng.
5. Postgres vẫn là source of truth.
```

Không chuyển sang:

```text
API publish queue trực tiếp trước DB commit
worker tin message thay vì load DB state
queue message chứa secret/credential plaintext
```

---

## 15. Failure modes

| Failure | Expected behavior |
|---|---|
| API commits but worker down | Job/outbox remains pending; alert by age |
| Worker dies after claim | `locked_until` expires; another worker reclaims |
| Worker dies after provider success before DB update | Reconcile by provider_request/idempotency/resource lookup; no blind create retry |
| Outbox dispatch fails | Retry with backoff; business state remains committed |
| Scheduler runs twice | Dedupe key/idempotent handler prevents duplicate effect |
| DB unavailable | Worker stops claiming, does not call provider based on stale memory |
| Provider 429 | Backoff and reduce source concurrency |
| Unknown provider result | Manual review or safe status lookup |

---

## 16. Acceptance criteria

Async architecture đạt khi:

```text
Checkout creates order/reservation/ledger/job/outbox in one transaction.
Worker claim uses row lock or equivalent safe mechanism.
Every job has idempotency key and correlation_id.
Provisioning unknown result never blindly retries create.
Outbox payloads are redacted.
Scheduler has duplicate-run protection.
Job retry policy distinguishes safe/unsafe/do_not_retry/manual_review.
Worker crash recovery is tested.
Manual review queue is visible to admin/operator.
Metrics and alerts cover queue depth, age, failure, and manual review.
```

P0 tests:

```text
double dispatch outbox does not duplicate service or notification.
worker crash after provider success does not create second resource.
same idempotency key with different payload is rejected.
reservation expiry job can run twice without double release.
```

---

## 17. Tóm tắt quyết định

```text
PostgreSQL-backed outbox/jobs first.
External queue later, only when pressure proves it is needed.
Outbox remains even with external queue.
Workers load and verify DB state before side effects.
Provider create is never retried blindly after unknown result.
Scheduler enqueues jobs; handlers stay idempotent.
Manual review is a first-class state, not an afterthought.
```
