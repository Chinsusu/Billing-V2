# 43 - Observability Logging Metrics Tracing Spec

Version: v1.4 Architecture Deep Dive  
Status: Draft for technical alignment  
Scope: Structured logging, metrics, tracing, dashboards, alerts, audit/observability boundary, retention, incident signals  
Related docs: 08, 19, 21, 22, 31, 37, 38, 39, 40, 41, 42, 44

---

## 1. Mục tiêu tài liệu

Observability của dự án không chỉ để xem server còn sống. Nó phải giúp phát hiện và điều tra các lỗi P0:

```text
ledger mismatch
tenant access leak
duplicate provisioning
provider timeout/partial success
credential exposure
queue backlog
worker crash loop
service lifecycle stuck
backup failure
```

Kết luận:

```text
Every request/job needs request_id/correlation_id.
Logs are structured and redacted.
Metrics cover money, tenant, provisioning, worker, provider, credential, lifecycle.
Tracing links API -> DB transaction -> outbox/job -> worker -> provider.
Alerts focus on safety and customer-impacting failure.
```

---

## 2. Observability vs audit

Observability và audit khác nhau.

Observability:

```text
debug runtime behavior
detect failure
measure latency/error/throughput
support incident response
```

Audit:

```text
prove who did what
track sensitive action
financial/security evidence
compliance/forensics
```

Do not rely on logs as audit trail for money/credential/security actions. Audit is stored in `audit_logs` with explicit schema and redaction.

Do not put secrets in either logs or audit.

---

## 3. Correlation model

Required IDs:

```text
request_id: per HTTP request
correlation_id: links business flow across API/job/provider/notification
job_id: async job
outbox_event_id: outbox event
provider_request_id: external provider call
audit_id: audit event
order_id/service_id/wallet_id where relevant
```

Checkout example:

```text
HTTP request_id A
correlation_id C
order_id O
ledger entries L1/L2
reservation_id R
provisioning_job_id J
provider_request_id P
service_id S
notification N
```

All logs/metrics/traces should allow searching by `correlation_id`.

---

## 4. Structured logging

Use structured JSON logs in production.

Required common fields:

```text
timestamp
level
environment
service/process
version
message
request_id
correlation_id
tenant_id
actor_id optional
job_id optional
source_id optional
provider_type optional
error_code optional
duration_ms optional
```

Log levels:

```text
debug: local/dev or temporary targeted investigation
info: normal lifecycle events
warn: recoverable abnormal state
error: failed operation requiring attention
critical: launch-blocking/security/money/credential incident
```

Do not log full request/response body by default.

---

## 5. Redaction in logs

Always redact:

```text
password
token
secret
api_key
api_secret
credential
private_key
authorization
cookie
set_cookie
otp
totp
recovery_code
encrypted_payload
```

Provider adapter logs must redact both request and response.

Safe references:

```text
credential_id
provider_account_id
source_id
external_resource_id
request_payload_hash
response_reference
```

Unsafe:

```text
root password
proxy auth
provider token
full raw provider response with credentials
```

Log redaction should be implemented centrally and tested with known sample secret values.

---

## 6. Metrics taxonomy

### 6.1 API metrics

```text
http_requests_total by route, method, status, error_code
http_request_duration_ms by route, method
http_request_inflight
rate_limited_requests_total by action
auth_login_failures_total by tenant
```

Do not label metrics with high-cardinality values like raw user email, order number, IP, or service id.

### 6.2 Database metrics

```text
db_query_duration_ms by operation/module
db_transaction_duration_ms by flow
db_connections_open
db_connections_in_use
db_lock_wait_ms
db_deadlocks_total
migration_duration_ms
```

Critical transaction metrics:

```text
checkout_transaction_duration_ms
topup_approval_transaction_duration_ms
renewal_transaction_duration_ms
provisioning_finalize_transaction_duration_ms
```

### 6.3 Wallet/ledger metrics

```text
ledger_entries_created_total by entry_type,direction
wallet_debit_total by entry_type,currency
wallet_credit_total by entry_type,currency
ledger_reconciliation_mismatch_total
wallet_negative_balance_detected_total
topup_approval_total by result
refund_total by result
```

Alert on any ledger mismatch, not just high rate.

### 6.4 Checkout/reservation metrics

```text
checkout_attempts_total by result,error_code
checkout_idempotency_conflict_total
reservation_created_total
reservation_expired_total
reservation_allocated_total
reservation_oversell_detected_total
out_of_stock_total by source_id
```

### 6.5 Worker/outbox metrics

```text
jobs_queued_count by job_type
jobs_running_count by job_type
jobs_succeeded_total by job_type
jobs_failed_total by job_type,error_code
job_duration_ms by job_type
job_oldest_queued_age_seconds by job_type
outbox_pending_count by event_type
outbox_oldest_pending_age_seconds
manual_review_count by reason
```

### 6.6 Provider metrics

```text
provider_requests_total by provider_type,source_id,operation,result,error_code
provider_request_duration_ms by provider_type,operation
provider_rate_limited_total
provider_timeout_total
provider_unknown_result_total
provider_auth_failed_total
provider_circuit_state by source_id
provider_health_status by source_id
```

### 6.7 Credential/security metrics

```text
credential_reveal_total by tenant,actor_type,result
credential_reveal_denied_total
credential_reveal_rate_limited_total
credential_redaction_test_fail_total
admin_2fa_failure_total
emergency_access_started_total
permission_denied_total by permission
cross_tenant_denied_total
```

---

## 7. Distributed tracing

Trace spans should cover:

```text
HTTP request
tenant resolution
auth/RBAC check
DB transaction
outbox insert
job claim
provider request
service activation finalization
notification dispatch
```

Span attributes:

```text
tenant_id
module
operation
error_code
job_type
provider_type
source_id
retry_safety
```

Avoid high-cardinality or sensitive attributes:

```text
email
credential
raw token
full request body
provider secret
```

Tracing should allow answering:

```text
Why is this order stuck?
Did worker call provider?
Was provider success persisted?
Which transaction failed?
Which notification was emitted?
```

---

## 8. Dashboards

### 8.1 Executive safety dashboard

```text
orders today
checkout success/failure
active services
provider health
manual review count
ledger mismatch count
queue backlog
credential reveal spike
```

### 8.2 Finance dashboard

```text
top-up approvals
ledger debit/credit totals
wallet reconciliation mismatch
refunds/adjustments
reseller settlement debit total
negative balance exceptions
```

### 8.3 Provisioning dashboard

```text
provisioning queued/running/succeeded/failed/manual_review
provider request latency/error
timeout unknown count
provider circuit state
oldest provisioning job age
```

### 8.4 Security dashboard

```text
login failures
permission denied
cross-tenant denied
credential reveal total/spike
emergency access sessions
admin 2FA failures
rate limited actions
```

### 8.5 Worker dashboard

```text
job depth by type
oldest job age
retry count
dead/manual review count
worker heartbeat
outbox pending age
```

---

## 9. Alerts

P0 alerts:

```text
ledger reconciliation mismatch > 0
reservation oversell detected > 0
provider timeout unknown spike
duplicate provisioning suspected
credential plaintext detected in logs/audit scan
cross-tenant access success suspected
admin 2FA disabled/failing unexpectedly
backup failed
restore test failed
DB unavailable
```

P1 alerts:

```text
provider down/degraded
manual review item older than threshold
outbox oldest pending age too high
worker crash loop
queue depth growing continuously
checkout failure rate high
credential reveal spike
top-up pending aging
```

Alert must include:

```text
severity
service/process
tenant/source if relevant
runbook link
dashboard link
correlation examples
```

---

## 10. SLOs and thresholds

Suggested early thresholds:

```text
API p95 latency < 500ms for normal reads
checkout API p95 < 2s excluding worker provisioning
worker provisioning job claim delay p95 < 60s
manual review P0 age < 30 minutes
ledger mismatch = 0
cross-tenant leak = 0
credential plaintext leak = 0
provider unknown timeout count monitored per source
```

Do not make provider provisioning success latency a user-facing synchronous SLO if provider can take minutes. Track queued/running/manual states clearly instead.

---

## 11. Retention

Suggested retention:

```text
application logs: 14-30 days hot
security logs: 90 days or more
audit_logs: 1-7 years depending policy
metrics: 30-90 days high resolution, longer downsampled
traces: 7-14 days, longer for sampled incidents
provider raw debug references: short retention, encrypted/restricted
```

Retention must respect:

```text
privacy
credential safety
storage cost
incident investigation needs
```

---

## 12. Health checks

API health:

```text
/health/live - process alive
/health/ready - can serve traffic: DB reachable, config loaded
```

Worker health:

```text
process alive
DB reachable
can claim/heartbeat
queue/outbox access ok
```

Scheduler health:

```text
process alive
lock held or standby known
last enqueue time by job type
```

Provider health is separate from app health. A provider down should not mark the API process unhealthy; it should mark source/provider degraded.

---

## 13. Incident investigation queries

For any order:

```text
order_id -> correlation_id
correlation_id -> ledger entries
correlation_id -> reservation
correlation_id -> provisioning_job
correlation_id -> provider_requests
correlation_id -> service/lifecycle
correlation_id -> notifications
correlation_id -> audit_logs
```

For credential reveal:

```text
credential_id/service_id -> audit credential.revealed
actor_id -> permissions/role at time
request_id/correlation_id -> API logs
rate limit logs -> suspicious behavior
```

For ledger mismatch:

```text
wallet_id -> ledger entries
wallet cache -> ledger sum
topup/order/refund references
adjustment records
audit actor/reason
```

---

## 14. Acceptance criteria

Observability đạt khi:

```text
Every API request has request_id and correlation_id.
Every job/provider request carries correlation_id.
Logs are structured and centrally redacted.
Metrics cover API, DB, wallet, checkout, worker, provider, credential, and security.
Dashboards exist for finance, provisioning, worker, security, and launch safety.
Alerts exist for P0 No-Go failures.
Tracing links checkout to worker/provider/service activation.
Logs/audit are tested for secret redaction.
Runbooks link from alerts.
```

P0 tests/checks:

```text
known secret sample does not appear in logs.
checkout correlation_id can trace order to service.
provider timeout unknown alert fires.
ledger mismatch alert fires in test scenario.
worker backlog alert fires when jobs age past threshold.
```

---

## 15. Tóm tắt quyết định

```text
Observability is part of safety, not a launch decoration.
Correlation_id is the spine of debugging.
Audit proves sensitive actions; logs explain runtime behavior.
Metrics must expose money, tenant, provider, worker, and credential risk.
No secret may be logged for the sake of convenience.
```
