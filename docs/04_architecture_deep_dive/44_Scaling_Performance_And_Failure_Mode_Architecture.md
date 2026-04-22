# 44 - Scaling Performance And Failure Mode Architecture

Version: v1.4 Architecture Deep Dive  
Status: Draft for technical alignment  
Scope: Scaling strategy, performance targets, bottlenecks, degradation, backpressure, provider failure, database pressure, launch capacity  
Related docs: 14, 19, 22, 24, 30, 31, 33, 37, 38, 39, 40, 43

---

## 1. Mục tiêu tài liệu

Tài liệu này khóa hướng scale và failure mode cho MVP và giai đoạn sau.

Không nên tối ưu sớm thành hệ thống phức tạp, nhưng phải tránh thiết kế không thể vận hành khi có:

```text
checkout tăng đột biến
worker backlog
provider down/rate limit
database lock contention
tenant lớn bất thường
credential reveal spike
manual review quá tải
```

Kết luận:

```text
Scale theo bottleneck thật.
Giữ modular monolith trong MVP.
PostgreSQL là source of truth.
Worker scale ngang nhưng phải tôn trọng provider/source concurrency.
Backpressure tốt hơn là tạo lỗi tiền/tài nguyên.
Failure mode phải fail closed ở tiền, tenant, credential, provisioning.
```

---

## 2. Scaling principles

Ưu tiên:

```text
correctness first
backpressure second
performance third
feature expansion last
```

Không scale bằng cách phá invariant:

```text
không bỏ transaction cho nhanh
không cache wallet balance để checkout không lock
không retry provider create để tăng success rate ảo
không bỏ audit để giảm latency
không bypass tenant guard trong report
```

---

## 3. Initial production topology

MVP/pilot:

```text
api replicas: 2
worker replicas: 1-3
scheduler replicas: 1 active
postgres: managed primary with backup
redis/cache: optional for rate limit/session/cache
object storage: proof attachments/report exports
monitoring/logging stack
secret manager or protected env secret system
```

Scale first:

```text
increase worker replicas
tune DB indexes
limit provider concurrency
add read replicas for reporting if needed
move notification fanout to queue if needed
```

Do not split billing/ledger/order into microservices during MVP.

---

## 4. Performance targets

Early targets:

```text
catalog list p95 < 500ms
service list p95 < 700ms
wallet ledger page p95 < 1s with pagination
checkout API p95 < 2s excluding provisioning
top-up approval p95 < 2s
credential reveal p95 < 1s
worker job claim p95 < 1s
provider provisioning async, not HTTP-bound
```

Operational targets:

```text
queue job oldest age visible
manual review age visible
provider request p95 visible by source
DB lock wait visible
```

Correctness targets:

```text
ledger mismatch = 0
reservation oversell = 0
cross-tenant leak = 0
credential plaintext leak = 0
duplicate provider resource due to retry = 0
```

---

## 5. Database scaling

PostgreSQL pressure points:

```text
wallet row locks during checkout/topup/refund
inventory row locks during reservation
job claim scans
audit log growth
ledger/report queries
service lifecycle cron scans
```

Controls:

```text
proper indexes
short transactions
deterministic wallet lock order
pagination on ledger/audit/service lists
batch jobs with limits
SKIP LOCKED for workers
avoid long report queries on primary during peak
```

Indexes should cover:

```text
tenant_id + created_at
tenant_id + status
wallet_id + created_at
source_id + status
job status + next_attempt_at + priority
service term_end_at + status
correlation_id
```

Future:

```text
read replica for reports
partition audit_logs/ledger by time if volume requires
archive old notifications/log-like data
materialized reporting tables for finance summaries
```

---

## 6. API scaling

API is stateless except session/cache strategy.

Must support:

```text
horizontal replicas
request timeout
rate limits
graceful shutdown
health/readiness checks
connection pool limits
```

API should not:

```text
call provider provisioning synchronously during checkout
hold long transactions while waiting external APIs
load huge tenant reports without pagination/export job
return large credential/provider raw payloads
```

Backpressure:

```text
checkout rate limit by tenant/user/source
credential reveal rate limit
admin export async if large
provider source disabled/degraded blocks new checkout if needed
```

---

## 7. Worker scaling

Worker can scale horizontally if:

```text
job claim uses row lock/SKIP LOCKED
jobs are idempotent
provider concurrency is limited
retry/backoff is enforced
manual review is explicit
```

Worker pool separation may be needed:

```text
provisioning workers
notification workers
sync workers
report/reconciliation workers
```

Do not let notification/report jobs starve provisioning jobs. Use:

```text
priority
separate queues/tables
worker type filters
```

Provider-bound scale is not the same as worker scale. More workers can worsen rate limits if source concurrency is not controlled.

---

## 8. Provider failure handling

Provider failure modes:

```text
down
degraded
rate limited
auth failed
stock drift
timeout unknown
partial success
state drift
credential missing
```

System behavior:

```text
new checkout blocked if source unavailable/out_of_stock by policy
existing paid order goes manual_review if provider uncertain
rate limit triggers backoff and concurrency reduction
auth failure disables source/provider account
state drift creates reconciliation/manual review
```

Do not auto-failover in MVP unless product has explicitly designed inventory/pricing/snapshot consequences.

---

## 9. Backpressure design

Backpressure is intentional refusal or slowing to keep invariants safe.

Examples:

```text
source degraded -> hide/disable checkout for source
provider rate limit -> reduce worker concurrency
queue too deep -> show provisioning queued/pending honestly
manual review aging -> pause provider/source expansion
ledger reconciliation mismatch -> freeze financial mutations if severe
credential reveal spike -> rate limit and alert
```

Do not convert backpressure into:

```text
blind retries
double debit
oversell
skipping audit
serving stale wallet balance for checkout
```

---

## 10. Cache strategy

Cache may be used for:

```text
public catalog read
tenant branding/domain mapping with short TTL
permission list with invalidation
provider health summary
rate limit counters
session cache
```

Do not use cache as authority for:

```text
wallet balance mutation
ledger state
inventory reservation
service credential
tenant security decision without fallback validation
provider resource truth after drift
```

If cache and PostgreSQL disagree:

```text
PostgreSQL wins.
```

---

## 11. Failure mode table

| Failure | Safe behavior |
|---|---|
| PostgreSQL unavailable | API rejects mutations; worker stops side effects; no provider create from stale job memory |
| Redis/cache unavailable | degrade rate limit/session if designed; DB still source of truth |
| Worker down | jobs remain queued; alert by age |
| Scheduler duplicate run | dedupe/idempotent job handlers prevent double effect |
| Provider down | block/degrade source; queued jobs retry safe or manual review |
| Provider timeout after create | manual review/status lookup; no blind retry |
| Ledger mismatch | freeze affected wallet/flow; reconcile; incident |
| Tenant mapping bug suspected | disable affected domain/tenant routes if needed; incident |
| Credential plaintext detected | rotate/contain; incident |
| Notification down | business state still commits; retry notification later |
| Object storage down | block proof upload/export; core checkout may continue if not dependent |

---

## 12. Graceful degradation

Functions that can degrade:

```text
notifications delayed
reports delayed/export queued
provider sync delayed
public catalog read cache stale for short TTL
support attachment upload paused
```

Functions that must fail closed:

```text
checkout without wallet/ledger transaction
reseller settlement without balance
credential reveal without auth/audit
tenant-scoped data read without tenant context
provider create after unknown timeout
ledger adjustment without permission/reason
```

---

## 13. Data growth management

High-growth tables:

```text
audit_logs
wallet_ledger_entries
provider_requests
jobs/job_attempts
notifications
service_lifecycle_events
```

Controls:

```text
pagination
retention policy by table
archive old job attempts
partition audit/provider_requests if volume grows
summary tables for reports
object storage for large redacted raw references
```

Never delete ledger entries to save space. Archive strategy must preserve financial truth.

---

## 14. Multi-tenant noisy neighbor controls

One tenant/reseller can overload:

```text
checkout requests
credential reveals
provider actions
support uploads
report exports
```

Controls:

```text
per-tenant rate limits
per-tenant job quotas or fair scheduling
per-source/provider concurrency
report export limits
manual review threshold
```

Metrics should be able to show top tenants by:

```text
request volume
job volume
provider failures
manual review count
credential reveals
```

---

## 15. Future scaling boundaries

Only split microservice when:

```text
module boundary is stable
load profile demands it
team/deploy independence matters
data consistency story is safe
observability and incident ownership are ready
```

Potential split candidates:

```text
notification service
reporting service
provider worker service
public API gateway
```

Avoid splitting early:

```text
wallet
ledger
order
reservation
provisioning orchestration
```

These need close transactional consistency in MVP.

---

## 16. Load testing plan

MVP load tests:

```text
concurrent checkout for last stock
top-up approval double submit
ledger page pagination with large history
service list with many services
worker provisioning batch with provider mock latency
provider 429 storm
credential reveal spike/rate limit
reservation expiry batch
service expiry batch
```

Load test must verify:

```text
no oversell
no double debit
no duplicate service
no unbounded retries
no secret logs
acceptable queue age under expected volume
```

---

## 17. Acceptance criteria

Scaling/failure architecture đạt khi:

```text
API, worker, scheduler can scale as separate processes.
Provider concurrency is controlled independently from worker count.
Checkout remains transaction-safe under concurrent load.
Job backlog, DB lock wait, provider errors, and manual review age are observable.
Cache is not source of truth for money/inventory/credential.
Failure modes fail closed for money, tenant, credential, and provisioning.
Load tests include concurrency and provider failure scenarios.
Runbooks exist for DB outage, provider outage, duplicate provisioning, ledger mismatch, credential exposure.
```

---

## 18. Tóm tắt quyết định

```text
Scale only after preserving correctness.
Modular monolith is enough until real pressure proves otherwise.
PostgreSQL remains source of truth.
Backpressure is safer than unsafe success.
Provider limits shape worker throughput.
Fail closed for money, tenant, credential, and resource creation.
```
