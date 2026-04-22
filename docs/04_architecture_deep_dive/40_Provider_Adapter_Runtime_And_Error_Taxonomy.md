# 40 - Provider Adapter Runtime And Error Taxonomy

Version: v1.4 Architecture Deep Dive  
Status: Draft for technical alignment  
Scope: Provider adapter runtime, timeout, rate limit, circuit breaker, normalized error taxonomy, retry safety, provider request logging  
Related docs: 05, 11, 18, 19, 21, 30, 37, 38, 39, 43, 44

---

## 1. Mục tiêu tài liệu

File 18 định nghĩa interface/behavior của provider adapter. File này khóa phần runtime: adapter chạy thế nào, timeout ra sao, lỗi provider được phân loại thế nào, khi nào retry, khi nào manual review.

Provider là vùng không đáng tin nhất của hệ thống:

```text
API có thể timeout.
API có thể tạo resource nhưng không trả response.
Stock có thể thay đổi giữa lúc reserve và provision.
Credential có thể thiếu.
Provider có thể đổi error format.
Provider có thể rate limit hoặc block account.
```

Kết luận:

```text
Adapter không được trả raw error lung tung.
Mọi lỗi phải map về taxonomy nội bộ.
Mọi operation phải gắn retry_safety.
Provider request phải được log redacted.
Unknown result không được coi là failed terminal.
```

---

## 2. Runtime components

Provider runtime gồm:

```text
adapter registry
provider account resolver
secret loader/decrypter
HTTP/API client factory
timeout policy
rate limiter
circuit breaker
provider request logger
error normalizer
credential sanitizer/encrypter
capability resolver
```

Trong Go codebase, các phần này nên nằm quanh:

```text
/internal/modules/provider
/internal/modules/provisioning
/internal/platform/crypto
/internal/platform/ratelimit
/internal/platform/httpclient
```

Rule:

```text
Provider module không debit tiền.
Provider module không tự tạo order/service.
Provider module không quyết định permission.
Provider module chỉ trả normalized result cho provisioning/service lifecycle.
```

---

## 3. Adapter registry

Mỗi provider type đăng ký một adapter implementation:

```text
proxmox
ovh
hetzner
manual
proxy_upstream
preloaded_proxy_pool
custom_api
```

Registry lookup:

```text
source_id -> provider_source -> provider_account -> provider_type -> adapter
```

Nếu không tìm thấy adapter:

```text
error_code = PROVIDER_ADAPTER_NOT_FOUND
retry_safety = do_not_retry
action = disable source or manual review
```

Adapter phải expose:

```text
ProviderType()
CapabilityProfile()
CheckHealth()
CheckStock()
Provision()
GetStatus()
Suspend()
Unsuspend()
Terminate()
Renew()
ResetPassword()
ChangeIp()
```

Không ép provider nào cũng support mọi action. Unsupported action là capability issue, không phải runtime exception.

---

## 4. Operation envelope

Mọi operation nhận envelope thống nhất:

```text
operation_context:
- operation_id
- tenant_id
- source_id
- provider_account_id
- actor_or_system_id
- idempotency_key
- correlation_id
- request_timeout
- deadline
- attempt_number
- capability_snapshot
- provider_source_snapshot
```

Mọi operation trả result thống nhất:

```text
operation_result:
- status
- external_request_id
- external_resource_id
- service_identifier
- credential_payload optional
- provider_status
- retry_safety
- error_code
- error_message_redacted
- raw_response_reference
- observed_at
```

Allowed status:

```text
success
failed
partial_success
unknown
manual_review_required
capability_not_supported
```

---

## 5. Timeout policy

Mỗi provider/source cần timeout cấu hình được:

```text
connect_timeout
request_timeout
status_timeout
health_timeout
max_total_deadline
```

Default gợi ý:

```text
health: 5s
stock: 8s
provision: 30s-120s tùy provider
status: 10s
suspend/terminate: 30s
```

Rule:

```text
Timeout trước khi request rời process có thể safe_retry.
Timeout sau khi request đã gửi đi là unknown/unsafe nếu operation tạo/sửa resource.
```

Adapter phải phân biệt tốt nhất có thể:

```text
request_not_sent
request_sent_no_response
response_partial
response_received
```

Nếu không phân biệt được, chọn hướng an toàn:

```text
unknown + manual_review_required
```

---

## 6. Error taxonomy

### 6.1 Auth and configuration

| Code | Meaning | Retry safety | Action |
|---|---|---|---|
| `PROVIDER_AUTH_FAILED` | API key/token invalid | do_not_retry | disable source/account, alert admin |
| `PROVIDER_PERMISSION_DENIED` | Account lacks permission | do_not_retry | operator fix permission |
| `PROVIDER_ACCOUNT_SUSPENDED` | Provider account suspended | do_not_retry | business/operator action |
| `PROVIDER_CONFIG_INVALID` | Missing config/template/region | do_not_retry | disable plan-source mapping |
| `PROVIDER_ADAPTER_NOT_FOUND` | No adapter implementation | do_not_retry | config/deploy fix |

### 6.2 Capacity and product mapping

| Code | Meaning | Retry safety | Action |
|---|---|---|---|
| `PROVIDER_OUT_OF_STOCK` | Provider confirmed no stock | do_not_retry | release reservation/refund if needed |
| `PROVIDER_PLAN_NOT_FOUND` | Provider plan/template missing | do_not_retry | disable mapping |
| `PROVIDER_REGION_UNAVAILABLE` | Location unavailable | do_not_retry or safe_retry by policy | source degraded |
| `PROVIDER_CAPABILITY_NOT_SUPPORTED` | Action unsupported | do_not_retry | API returns capability error |

### 6.3 Transient runtime

| Code | Meaning | Retry safety | Action |
|---|---|---|---|
| `PROVIDER_RATE_LIMITED` | 429/rate limit | safe_retry | backoff, reduce concurrency |
| `PROVIDER_TEMPORARY_ERROR` | 5xx or temporary fail before side effect certainty | safe_retry | retry limited |
| `PROVIDER_NETWORK_ERROR_BEFORE_SEND` | request not sent | safe_retry | retry |
| `PROVIDER_MAINTENANCE` | provider maintenance | safe_retry/manual_review | backoff/alert |

### 6.4 Uncertain and dangerous

| Code | Meaning | Retry safety | Action |
|---|---|---|---|
| `PROVIDER_TIMEOUT_UNKNOWN` | Create/action request may have reached provider | unsafe_retry | status lookup/manual review |
| `PROVIDER_TIMEOUT_REQUEST_KNOWN` | external_request_id known but no final status | manual_review_required | query by request id |
| `PROVIDER_PARTIAL_SUCCESS` | Resource may exist but data incomplete | manual_review_required | reconcile/fetch missing fields |
| `PROVIDER_STATE_DRIFT` | Provider state conflicts with DB | manual_review_required | reconciliation |
| `PROVIDER_RESOURCE_ALREADY_EXISTS` | Duplicate/idempotent resource found | safe if mapped | link existing resource |

### 6.5 Credential-specific

| Code | Meaning | Retry safety | Action |
|---|---|---|---|
| `PROVIDER_CREDENTIAL_MISSING` | Success but no credential | manual_review_required | fetch credential/manual |
| `PROVIDER_CREDENTIAL_INVALID` | Credential returned unusable | manual_review_required | reset/fetch/manual |
| `PROVIDER_CREDENTIAL_ROTATION_FAILED` | Reset/rotate failed | safe_retry or manual_review | depends if request was sent |

---

## 7. Retry safety rules

### 7.1 Safe retry

Safe retry only when:

```text
operation is read-only
request did not leave process
provider operation is idempotent with strong idempotency key
provider explicitly says no resource was created
error is rate limit/temporary before side effect
```

### 7.2 Unsafe retry

Unsafe retry when:

```text
create/provision request may have reached provider
terminate/suspend request may have reached provider
provider lacks idempotency support
external_resource_id unknown
response is contradictory or incomplete
```

Action:

```text
manual_review_required
or safe status lookup before retry
```

### 7.3 Do not retry

Do not retry:

```text
auth failed
permission denied
plan/template missing
confirmed out of stock
source disabled
provider account suspended
capability unsupported
```

---

## 8. Provider idempotency

Provider support levels:

```text
level 0: no idempotency support
level 1: can tag resource with internal idempotency/correlation metadata
level 2: can query by external_request_id
level 3: provider guarantees idempotency key on create
```

Policy:

```text
level 0 providers require stricter manual review on unknown result.
level 1 providers must tag resource if API allows metadata/name/comment.
level 2 providers should query request status before manual review.
level 3 providers can retry create only if provider guarantee is documented/tested.
```

Internal idempotency key should be included where safe:

```text
hostname prefix
metadata tag
order reference
client reference redacted
provider request field
```

Do not leak sensitive tenant/client info into provider-visible metadata.

---

## 9. Provider request logging

Every external call creates/updates `provider_requests`.

Required:

```text
provider_request_id
job_id
tenant_id
source_id
operation
idempotency_key
external_request_id nullable
external_resource_id nullable
request_payload_hash
request_summary_redacted
response_summary_redacted
status
retry_safety
error_code
sent_at
received_at
duration_ms
correlation_id
```

Do not store:

```text
provider_api_key
password
root credential
proxy credential
session token
full raw response containing secret
```

If raw response is needed for debugging:

```text
store in restricted object storage
encrypt it
redact first if possible
store reference only in DB
short retention
```

---

## 10. Rate limiting and concurrency

Runtime must support concurrency limits:

```text
global provider concurrency
per provider_account concurrency
per source concurrency
per operation concurrency
```

Examples:

```text
provision max 2 concurrent per source
status sync max 10 concurrent per provider account
terminate max 1-2 concurrent if dangerous
```

Provider 429:

```text
respect Retry-After
increase backoff
emit metric
consider circuit breaker degraded
```

---

## 11. Circuit breaker

Circuit states:

```text
closed
degraded
open
half_open
manual_disabled
```

Triggers:

```text
auth failure -> open/manual_disabled
high timeout rate -> degraded/open
high 5xx rate -> degraded/open
confirmed provider outage -> open
manual operator disable -> manual_disabled
```

Behavior:

```text
closed: normal
degraded: allow limited jobs, warn UI/admin
open: block new provisioning, allow status/recovery checks if safe
half_open: allow small test requests
manual_disabled: only operator can re-enable
```

Circuit breaker state must not silently hide failed jobs. It should produce:

```text
admin alert
provider health report
source status update
audit event for manual change
```

---

## 12. Capability runtime

Capability is not just UI decoration.

Runtime checks:

```text
API checks service capability_snapshot before enqueue action job.
Worker checks capability before calling provider.
Adapter checks provider runtime support before operation.
```

If service capability snapshot says action was supported but current source is disabled:

```text
block or manual_review depending operator policy
return PROVIDER_UNAVAILABLE or CAPABILITY_TEMPORARILY_UNAVAILABLE
```

If current provider newly supports a feature, old service does not automatically get it unless product policy enables capability migration.

---

## 13. Manual provider/source

Manual provider is first-class, not a hack.

Manual provision flow:

```text
checkout creates provisioning_job
job enters manual_review/manual_waiting_provider
operator creates resource outside system
operator enters external_resource_id/service_identifier/credential
system encrypts credential
system activates service through same service activation transaction
```

Manual provider must still follow:

```text
reservation allocation
wallet debit already posted
credential encryption
audit
lifecycle event
provider_resource_mapping if external id exists
```

---

## 14. Testing contract

Provider mock must simulate:

```text
success
out_of_stock
auth_failed
rate_limited
temporary_500
timeout_before_send
timeout_after_send_unknown
success_missing_credential
partial_success_with_resource_id
resource_already_exists
state_drift
slow_response
```

Adapter tests:

```text
maps raw provider error to internal code
sets retry_safety correctly
redacts request/response
does not return plaintext secret in logs
honors timeout/context cancellation
handles idempotency key duplicate
```

---

## 15. Acceptance criteria

Provider runtime đạt khi:

```text
Every adapter operation returns normalized operation_result.
Every provider error maps to internal error_code and retry_safety.
Provider request records are redacted and traceable by correlation_id.
Timeout after create request moves to manual_review or safe status lookup.
Provider auth failure disables/degrades source and alerts admin.
Rate limit triggers backoff, not tight retry loop.
Capability unsupported is blocked by API and worker.
Manual provider flow activates service through same audited service path.
Provider mock covers success/failure/timeout/partial cases.
```

P0 tests:

```text
timeout after send does not create duplicate resource.
success with missing credential does not expose blank credential as active secret.
provider auth failure is do_not_retry.
raw provider secret never appears in provider_requests/logs.
```

---

## 16. Tóm tắt quyết định

```text
Provider adapter is a normalization boundary.
Runtime errors must be classified, not thrown upward raw.
Retry safety is part of every provider result.
Unknown provider state is manual review by default.
Circuit breaker protects users, provider accounts, and support team.
Manual provider remains controlled by the same ledger/reservation/audit rules.
```
