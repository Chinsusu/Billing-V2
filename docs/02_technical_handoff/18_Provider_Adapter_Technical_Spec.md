# 18 - Provider Adapter Technical Spec

## 1. Mục tiêu tài liệu

Tài liệu này định nghĩa contract kỹ thuật cho provider adapter layer. Đây là tài liệu interface/behavior, **không phải code triển khai**.

Provider adapter có nhiệm vụ biến các provider khác nhau thành một lớp hành vi thống nhất cho hệ thống:
```text
Proxmox
OVH
Hetzner
manual VPS provider
proxy upstream
preloaded proxy pool
custom API provider
```

Điểm quan trọng: không ép mọi provider hỗ trợ cùng một action. Adapter phải khai báo capability rõ ràng, UI/API phải dựa vào capability để cho phép hoặc chặn action.

---

## 2. Vai trò của provider adapter

Adapter layer nằm giữa:
```text
provisioning_worker/service_lifecycle_module
và
provider/source thật
```

Adapter chịu trách nhiệm:
- Kiểm tra health provider.
- Kiểm tra stock/capacity.
- Provision tài nguyên.
- Lấy trạng thái tài nguyên.
- Suspend/unsuspend/terminate nếu provider hỗ trợ.
- Reset password/reinstall/change IP nếu provider hỗ trợ.
- Fetch/generate credential.
- Chuẩn hóa lỗi provider thành error code nội bộ.
- Đánh dấu retry safety.
- Không làm lộ provider secret/credential trong log.

Adapter không chịu trách nhiệm:
- Debit wallet.
- Quyết định giá.
- Quyết định tenant permission.
- Quyết định order có được mua hay không.
- Tự ý retry ngoài policy của job.

---

## 3. Adapter capability profile

Mỗi source phải khai báo `capability_profile`.

### 3.1 Capability chung

```text
supports_health_check
supports_live_stock_check
supports_auto_provision
supports_manual_provision
supports_status_sync
supports_suspend
supports_unsuspend
supports_terminate
supports_renew
supports_reset_password
supports_reinstall
supports_change_ip
supports_bandwidth_usage
supports_console
supports_reverse_dns
supports_snapshot
supports_backup
supports_credential_fetch
supports_credential_rotation
```

### 3.2 VPS-specific capability

```text
supports_os_template_selection
supports_custom_hostname
supports_ipv6
supports_private_network
supports_resize
supports_rescue_mode
supports_vnc_console
supports_ssh_key_injection
```

### 3.3 Proxy-specific capability

```text
supports_proxy_protocol_http
supports_proxy_protocol_socks5
supports_rotating_proxy
supports_static_proxy
supports_geo_selection
supports_ip_whitelist
supports_userpass_auth
supports_bandwidth_quota
supports_thread_limit
supports_change_exit_ip
```

### 3.4 Capability rule

Nếu capability = false:
```text
UI không hiện action
API trả CAPABILITY_NOT_SUPPORTED nếu user gọi trực tiếp
worker không tạo job action đó
```

Nếu capability thay đổi sau khi service đã bán:
```text
service dùng capability_snapshot tại thời điểm mua
nhưng admin có thể override nếu provider/source bị disable/maintenance
```

---

## 4. Adapter operation contract

Mỗi operation phải trả về kết quả chuẩn hóa:

```text
operation_result:
- status: success / failed / partial_success / unknown / manual_review_required
- external_request_id
- external_resource_id
- provider_status
- credentials_encrypted_payload optional
- public_service_identifier optional
- retry_safety: safe_retry / unsafe_retry / do_not_retry / manual_review_required
- error_code optional
- error_message_redacted optional
- raw_response_reference optional
```

Không trả:
```text
provider_api_secret
plaintext credential trong log
full raw response nếu chứa secret
```

---

## 5. Required operations

### 5.1 `checkHealth`

Mục tiêu:
```text
xác định provider/source còn gọi được không
```

Output:
```text
healthy
degraded
down
unknown
```

Nếu provider trả 401/403:
```text
source health = down
provider credential invalid
alert admin
không retry provisioning mới qua source này
```

### 5.2 `checkStock`

Mục tiêu:
```text
xác định source còn capacity cho plan/source không
```

Output:
```text
available
out_of_stock
unknown
capacity_count optional
```

Rule:
- Nếu source finite stock: dùng DB inventory là chính, có thể sync provider.
- Nếu provider live stock: phải check live trước reserve hoặc theo cache TTL.
- Nếu unknown: tùy policy, có thể cho manual_review hoặc chặn checkout.

### 5.3 `provision`

Input logic:
```text
tenant context
order item snapshot
plan specs
source config
idempotency_key
correlation_id
```

Output:
```text
external_resource_id
service_identifier
credentials
provider_status
```

Status:
- `success`: provider đã tạo tài nguyên, có thể tạo service active.
- `failed`: provider chắc chắn không tạo tài nguyên.
- `partial_success`: provider có thể đã tạo hoặc tạo thiếu dữ liệu.
- `unknown`: mất kết nối/timeout không biết kết quả.
- `manual_review_required`: cần người xử lý.

### 5.4 `getStatus`

Mục tiêu:
```text
sync trạng thái tài nguyên thật với service trong hệ thống
```

Output:
```text
external_status
is_running
is_suspended
is_terminated
usage metrics optional
```

Không tự đổi service_status nếu không qua lifecycle module.

### 5.5 `suspend` / `unsuspend`

Rule:
- Nếu provider không hỗ trợ: trả `CAPABILITY_NOT_SUPPORTED`.
- Nếu suspend do billing/abuse, reason phải được truyền vào job.
- Nếu suspend thành công provider nhưng hệ thống fail update, worker phải retry update nội bộ trước khi gọi provider lần nữa.

### 5.6 `terminate`

Terminate là destructive action.

Rule:
```text
không retry mù terminate nếu provider response unknown
manual review nếu không chắc tài nguyên đã bị xóa hay chưa
service chỉ chuyển terminated khi xác nhận hoặc admin resolve
```

### 5.7 `renew`

Một số provider cần gọi renew API, một số chỉ cần internal billing.

Rule:
```text
nếu provider renew required -> tạo job renew
nếu không -> chỉ extend term_end_at internal
```

### 5.8 `resetPassword`

Output:
```text
new encrypted credential
provider_status
```

Rule:
```text
credential cũ chuyển rotated/revoked nếu reset thành công
audit credential.rotated
```

### 5.9 `reinstall`

Reinstall có thể xóa dữ liệu.

Rule:
```text
client confirmation required
warning required
provider capability required
audit service.reinstall_requested
```

### 5.10 `changeIp`

Rule:
```text
check capability
check rate limit
possible extra cost
provider may return new endpoint/credential
update service_identifier and credential if needed
audit service.ip_changed
```

---

## 6. Error normalization

Provider error phải map về internal error.

| Provider situation | Internal code | Retry safety | Action |
|---|---|---|---|
| API key invalid / 401 | `PROVIDER_AUTH_FAILED` | do_not_retry | disable/alert |
| Rate limit / 429 | `PROVIDER_RATE_LIMITED` | safe_retry | backoff |
| Out of stock | `PROVIDER_OUT_OF_STOCK` | do_not_retry | release reservation |
| Timeout before response | `PROVIDER_TIMEOUT_UNKNOWN` | unsafe_retry | check resource/manual review |
| Timeout after provider request id known | `PROVIDER_TIMEOUT_REQUEST_KNOWN` | manual_review_required | query status first |
| 500 temporary | `PROVIDER_TEMPORARY_ERROR` | safe_retry | retry limited |
| Plan/template not found | `PROVIDER_PLAN_NOT_FOUND` | do_not_retry | disable source/plan mapping |
| Success but missing credential | `PROVIDER_CREDENTIAL_MISSING` | manual_review_required | fetch credential/manual |
| Resource already exists by idempotency | `PROVIDER_RESOURCE_ALREADY_EXISTS` | safe if mapped | link existing resource |
| Provider says terminated but system active | `PROVIDER_STATE_DRIFT` | manual_review_required | reconcile |

---

## 7. Retry policy

### 7.1 Safe retry

Có thể retry nếu:
```text
provider chắc chắn chưa tạo tài nguyên
hoặc operation idempotent
hoặc provider hỗ trợ idempotency key mạnh
```

Ví dụ:
```text
429 rate limit
temporary 500 trước khi gửi create
health check fail
status sync fail
```

### 7.2 Unsafe retry

Không retry tự động nếu:
```text
timeout sau khi đã gửi create request
provider không hỗ trợ idempotency
không biết external_resource_id
response mâu thuẫn
```

Action:
```text
manual_review_required
query provider by external_request_id nếu có
query by metadata/tag/idempotency nếu provider hỗ trợ
```

### 7.3 Do not retry

Không retry nếu:
```text
API key invalid
plan/template không tồn tại
source disabled
out_of_stock confirmed
permission denied
provider account suspended
```

---

## 8. Idempotency with provider

### 8.1 Internal idempotency

Mỗi job có:
```text
job_id
idempotency_key
correlation_id
```

Adapter phải nhận idempotency_key từ job, không tự tạo mới khi retry.

### 8.2 Provider idempotency

Nếu provider hỗ trợ:
```text
truyền idempotency_key vào provider request
lưu external_request_id
lưu external_resource_id
```

Nếu provider không hỗ trợ:
```text
phải gắn metadata/tag nếu được
hoặc dùng manual review khi unknown
```

### 8.3 Duplicate detection

Trước khi gọi create lần nữa, worker/adapter phải kiểm tra:
```text
provider_resource_mappings
provider_requests with same idempotency_key
provider external lookup if supported
```

---

## 9. Provider request logging

Mỗi provider call tạo `provider_requests`.

Lưu:
```text
request_type
external_request_id
external_resource_id
status
retry_safety
request_summary_redacted
response_summary_redacted
error_code
sent_at
received_at
correlation_id
```

Không lưu:
```text
API secret
root password plaintext
proxy password plaintext
full auth header
session token
```

Nếu raw response cần lưu để debug:
```text
lưu private object storage
redact trước
set retention
restrict permission
```

---

## 10. Credential handling

Adapter có thể nhận credential từ provider hoặc tự generate.

Rule:
```text
credential phải được encrypt trước khi lưu DB
plaintext chỉ tồn tại trong memory ngắn hạn
không đưa credential vào audit/provider_request log
service detail API chỉ trả masked_hint
reveal credential là API/action riêng
```

Credential payload logic:
```text
vps:
- hostname/ip
- username
- password or ssh key
- port
- console link if supported

proxy:
- host
- port
- protocol
- username
- password
- geo
- expiry/quota if any
```

---

## 11. Manual provider/source

Phase 1 có thể có provider manual.

Manual source flow:
```text
checkout -> reserve -> debit wallet -> provisioning_job manual_review
operator nhập external_resource_id/service_identifier/credential
system tạo service active
```

Manual source vẫn phải:
- có source_id.
- có reservation.
- có ledger.
- có audit.
- có credential encrypted.
- có lifecycle event.

Không được bypass hệ thống bằng cách tạo service tay không có order/ledger nếu là paid service.

---

## 12. Proxmox adapter notes

Proxmox thường có:
```text
node
storage
template
vmid
network bridge
ip allocation
cloud-init/user/password
```

Cần chú ý:
- VMID unique.
- clone template có thể mất thời gian.
- lock VM trong lúc clone/reinstall.
- password/cloud-init credential có thể không fetch lại dễ dàng.
- partial success có thể xảy ra khi clone thành công nhưng API timeout.

Capability thường:
```text
auto_provision yes
suspend yes via stop/disable policy
terminate yes
reset_password depends cloud-init/guest agent
reinstall yes but destructive
console maybe yes
```

---

## 13. Proxy upstream adapter notes

Proxy source có thể là:
```text
preloaded list
upstream API
rotating proxy account
dedicated static proxy
```

Cần chú ý:
- Một proxy credential có thể là tài nguyên thật.
- Oversell dễ xảy ra nếu preloaded list không lock atomic.
- Change IP có thể là rotate endpoint hoặc cấp proxy mới.
- Bandwidth/quota sync nếu provider hỗ trợ.

Preloaded proxy inventory:
```text
proxy_pool_items
status: available/reserved/allocated/disabled
reservation_id
service_id
```

---

## 14. Provider onboarding checklist

Trước khi bật source active:
- Provider account credential được encrypt.
- Health check pass.
- Capability profile được xác nhận.
- Test provision sandbox/manual nhỏ.
- Test timeout behavior.
- Test credential retrieval.
- Test terminate/suspend nếu hỗ trợ.
- Define retry safety map.
- Define stock mode.
- Define plan/source mapping.
- Define provider cost.
- Define abuse/takedown process.
- Add monitoring alert.

---

## 15. Adapter acceptance criteria

Adapter layer đạt khi:
- Mọi provider action trả result chuẩn hóa.
- Capability profile điều khiển UI/API/worker.
- Timeout create không retry mù.
- Credential không xuất hiện trong logs/audit.
- Idempotency key được truyền xuyên job/provider request.
- Out-of-stock release reservation đúng.
- Provider auth fail disable/alert đúng.
- Manual provider vẫn có order/ledger/reservation/service/audit.
- Có error mapping table cho từng provider.
- Có manual review flow cho partial_success/unknown.

Câu nền: **adapter tốt không làm provider giống nhau; adapter tốt làm hệ thống biết provider khác nhau ở đâu và phản ứng an toàn.**
