# 19 - Worker Queue And Cron Jobs Spec

## 1. Mục tiêu tài liệu

Tài liệu này định nghĩa các worker/job nền cần có cho platform VPS/Proxy. Đây là job behavior spec, không phải code.

Dự án này không nên xử lý tất cả bằng request realtime. Các việc như provisioning, provider sync, expiry, suspend/terminate, notification phải chạy qua queue/cron để:
- tránh request timeout.
- kiểm soát retry.
- giữ idempotency.
- ghi audit/correlation.
- xử lý provider chậm/lỗi.

---

## 2. Nguyên tắc worker/job

### 2.1 Idempotent

Job chạy lại không được:
```text
trừ tiền lần hai
reserve stock lần hai
allocate stock lần hai
tạo service trùng
gọi provider create trùng nếu không an toàn
gửi notification critical trùng vô hạn
```

### 2.2 Observable

Mỗi job phải có:
```text
job_id
job_type
status
attempt_count
max_attempts
correlation_id
started_at
finished_at
last_error_code
next_retry_at
manual_review_reason
```

### 2.3 Retry có phân loại

Không retry theo kiểu “cứ lỗi là chạy lại”.

Retry dựa vào:
```text
safe_retry
unsafe_retry
do_not_retry
manual_review_required
```

### 2.4 Tenant scoped

Job liên quan tenant phải lưu `tenant_id` và không xử lý resource ngoài tenant đó.

### 2.5 Audit important transitions

Job làm thay đổi:
```text
service_status
order_status
wallet/ledger
reservation
credential
provider/resource
```

thì phải audit.

---

## 3. Queue message contract

Message logic:

```text
job_id
job_type
tenant_id
reference_type
reference_id
source_id optional
idempotency_key
correlation_id
attempt_count
created_at
scheduled_at optional
```

Reference examples:
```text
order_id
reservation_id
service_id
topup_request_id
provider_source_id
abuse_case_id
```

---

## 4. Job status

```text
queued
running
success
failed
retry_scheduled
manual_review
cancelled
```

Transition:
```text
queued -> running
running -> success
running -> failed
running -> retry_scheduled
running -> manual_review
retry_scheduled -> queued
manual_review -> queued/success/failed/cancelled by operator
```

---

## 5. Core worker: `provisioning_worker`

### 5.1 Trigger

Created by checkout when:
```text
order paid
reservation created
wallet debit posted
risk check passed
```

### 5.2 Input

```text
provisioning_job_id
order_id
order_item_id
reservation_id
source_id
idempotency_key
correlation_id
```

### 5.3 Preconditions

- job status queued/retry_scheduled.
- order payment_status paid.
- reservation status reserved.
- source active.
- provider health not down unless manual override.
- no existing service for same order_item unless resolving.

### 5.4 Steps

```text
1. mark job running
2. load order/order_item/reservation/source snapshot
3. check idempotency/provider request history
4. call adapter.provision
5. if success:
   - create provider_resource_mapping
   - create encrypted service_credentials
   - create service status active
   - reservation reserved -> allocated
   - order provisioning_status success
   - order/service lifecycle event
   - send activation notification
6. if failed and retry_safety safe_retry:
   - schedule retry if attempts left
7. if failed do_not_retry:
   - release reservation
   - reverse wallet debit if policy requires
   - order failed or manual_review depending policy
8. if partial_success/unknown/unsafe_retry:
   - job manual_review
   - order provisioning_status manual_review
   - alert operator
```

### 5.5 Failure handling

| Situation | Handling |
|---|---|
| Provider timeout unknown | manual_review, no retry mù |
| Provider out of stock | release reservation, refund/reversal, notify |
| Provider auth fail | disable source, alert admin, do not retry |
| Credential missing | manual_review, attempt credential fetch if safe |
| DB update fails after provider success | retry internal reconciliation before new provider call |

### 5.6 Audit

```text
provisioning.job.started
provider.request.sent
provisioning.job.succeeded
service.activated
reservation.allocated
provisioning.job.failed
provisioning.job.manual_review
```

---

## 6. Worker: `provider_sync_worker`

### 6.1 Trigger

- scheduled by provider sync cron.
- manual admin action.
- after unknown/partial provisioning.

### 6.2 Purpose

Sync external provider status with internal service mapping.

### 6.3 Rules

- Không tự terminate service chỉ vì provider không thấy resource một lần.
- Nếu mismatch, tạo provider_state_drift risk/manual review.
- Nếu service internal active nhưng provider terminated, alert admin.
- Nếu provider suspended nhưng internal active, mark drift and review.

### 6.4 Audit

```text
provider.sync.started
provider.sync.completed
provider.state_drift.detected
```

---

## 7. Worker: `service_action_worker`

Handles:
```text
suspend
unsuspend
terminate
reset_password
reinstall
change_ip
renew provider-side
```

### 7.1 Preconditions

- action supported by capability_snapshot/current policy.
- actor/action permission already validated before job creation.
- service belongs to tenant.
- action not conflicting with another running action.

### 7.2 Conflict examples

```text
cannot reinstall while terminate running
cannot change_ip while service terminated
cannot unsuspend if abuse case still active unless override
cannot renew terminated service
```

### 7.3 Audit

```text
service.action.started
service.suspended
service.unsuspended
service.terminated
credential.rotated
service.ip_changed
service.action.failed
```

---

## 8. Worker: `notification_worker`

### 8.1 Trigger

Created by modules/jobs when notification should be sent.

### 8.2 Channels

```text
email
dashboard
telegram
webhook
```

### 8.3 Rules

- Notification payload redacted.
- Credential activation email should avoid plaintext password unless policy explicitly allows. Safer: send “login to reveal”.
- Critical admin alert should include correlation_id.
- Failed send retry with backoff.
- Do not spam duplicate event; use dedupe key.

### 8.4 Dedupe key examples

```text
service_expiry_reminder:{service_id}:{days_before}
reseller_low_balance:{tenant_id}:{date}
provider_down:{source_id}:{date_hour}
```

---

## 9. Cron: `reservation_expiry_job`

### 9.1 Frequency

```text
every 1 minute
```

### 9.2 Logic

Find:
```text
reservations status = reserved
expires_at < now
```

For each:
```text
mark reservation expired
decrement reserved_count
if order not paid/provisioning -> order expired/cancelled
if wallet already debited due to edge case -> reversal per policy
audit reservation.expired
```

### 9.3 Idempotency

- If reservation already allocated/released/expired, skip.
- Use lock on reservation row/source inventory.
- Running twice must not decrement stock twice.

---

## 10. Cron: `service_expiry_job`

### 10.1 Frequency

```text
every 5-15 minutes
```

### 10.2 Logic

Find active services:
```text
term_end_at < now
billing_status not overdue/grace
```

Then:
```text
service_status may remain active or become expired depending policy
billing_status -> overdue/grace
send expired/grace notification
create lifecycle event
```

If policy says auto suspend immediately:
```text
create service_action_job suspend
```

### 10.3 Notes

Do not terminate immediately on expiry unless product policy says no grace.

---

## 11. Cron: `suspension_job`

### 11.1 Frequency

```text
every 15 minutes
```

### 11.2 Logic

Find:
```text
services billing_status = grace/overdue
grace_until_at < now
service_status active/expired
```

Then:
```text
create suspend job
reason = billing_overdue
notify client/reseller
```

### 11.3 Guard

- Skip if service renewed.
- Skip if already suspended/terminated.
- If provider does not support suspend, mark internal suspended and alert if required.

---

## 12. Cron: `termination_job`

### 12.1 Frequency

```text
daily or every few hours depending product risk
```

### 12.2 Logic

Find:
```text
services suspended/expired
terminate_after_at < now
policy auto_terminate = true
```

Then:
```text
create terminate job
```

### 12.3 Guard

Terminate is destructive:
- require product policy.
- keep audit.
- allow admin hold flag to prevent auto terminate.
- if provider response unknown, manual review.

---

## 13. Cron: `renewal_reminder_job`

### 13.1 Frequency

```text
daily
```

### 13.2 Reminder windows

Suggested:
```text
7 days before
3 days before
1 day before
on expiry
after grace start
```

### 13.3 Dedupe

```text
service_id + reminder_type + date/window
```

---

## 14. Cron: `provider_health_check_job`

### 14.1 Frequency

```text
every 5-15 minutes for active sources
```

### 14.2 Logic

For each active source:
```text
adapter.checkHealth
update health_status
if down/degraded -> alert admin
if recovered -> alert/reopen source depending policy
```

### 14.3 Guard

- Do not disable source on one transient failure unless threshold reached.
- Disable auto-provision if repeated critical failure.

---

## 15. Cron: `provider_inventory_sync_job`

### 15.1 Frequency

```text
every 15-60 minutes
```

### 15.2 Logic

- Sync capacity for provider_live source.
- Update available_count_cache.
- Detect out_of_stock.
- Alert if stock below threshold.

### 15.3 Guard

Do not override reserved/allocated counts incorrectly. Provider live availability and internal reservations must reconcile.

---

## 16. Cron: `pending_topup_review_reminder_job`

### 16.1 Frequency

```text
every 30-60 minutes
```

### 16.2 Logic

Find:
```text
topup_requests submitted/under_review older than SLA threshold
```

Notify:
```text
finance/admin/reseller owner
```

---

## 17. Cron: `reseller_low_balance_job`

### 17.1 Frequency

```text
daily or every few hours
```

### 17.2 Trigger

```text
reseller settlement wallet < configured threshold
or reseller balance < estimated cost for next N orders/days
```

Notify reseller:
```text
new client orders may not provision if reseller wallet is insufficient
```

---

## 18. Cron: `audit_retention_job`

### 18.1 Frequency

```text
daily/weekly
```

### 18.2 Rule

- Financial/security audit retained long-term.
- Low-risk debug logs can be archived/deleted per retention.
- Never delete ledger.
- Never delete critical audit before retention policy allows.

---

## 19. Manual review queue

Manual review items should appear in Admin Portal:

Types:
```text
provisioning_unknown
provider_state_drift
credential_missing
payment_mismatch
abuse_case
high_risk_order
unsafe_retry
```

Each item must show:
```text
severity
age
tenant
order/service
correlation_id
recommended next action
```

Allowed resolutions:
```text
retry_safe
link_existing_resource
mark_failed_and_refund
mark_success
cancel_order
suspend_service
clear_flag
escalate
```

Every resolution requires note/audit.

---

## 20. Monitoring metrics

P0 metrics:
```text
jobs queued by type
jobs running by type
jobs failed by type
jobs manual_review by type
avg provisioning time
provider success/fail rate
reservation expiry count
service expiry/suspend/terminate count
notification failure count
topup pending age
```

Alert examples:
```text
provider down > 5 minutes
provisioning manual_review > threshold
queue backlog > threshold
reservation expiry spike
ledger adjustment spike
credential reveal spike
```

---

## 21. Worker acceptance criteria

Worker/cron spec đạt khi:
- Mọi job idempotent.
- Provisioning success tạo service/credential/reservation allocation đúng.
- Unknown provider response không retry mù.
- Reservation expired không double-release.
- Expiry/suspend/terminate job không đụng service đã renewed.
- Notification không spam trùng vô hạn.
- Provider health issue alert đúng.
- Manual review queue có resolution rõ.
- Có correlation_id xuyên job/provider/order/service.
- Job failure có trạng thái, error code, next step.

Câu nền: **worker là nơi hệ thống thể hiện bản lĩnh thật: không phải lúc mọi thứ chạy tốt, mà khi provider timeout, mạng rớt, queue chạy lại và dữ liệu vẫn không sai.**
