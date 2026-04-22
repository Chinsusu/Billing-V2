# Tài liệu 10 - Tenant Security & Access Control Specification

## 1. Mục tiêu tài liệu
Tài liệu này khóa các rule bảo mật tenant, role, domain, credential và audit access cho phase 1.

Mục tiêu lớn nhất: không để client/reseller/staff đọc hoặc thao tác nhầm dữ liệu của tenant khác.

## 2. Nguyên tắc lõi
- Không tin `tenant_id` từ request body.
- Tenant context lấy từ domain/session/token đã xác thực.
- Backend phải enforce tenant scope; UI ẩn không đủ.
- Credential là secret, không phải dữ liệu text bình thường.
- Emergency access của Admin phải có reason và audit.

## 3. Request context chuẩn
Mỗi request sau khi auth nên có context:
- actor_id
- actor_type: admin, reseller_owner, reseller_staff, client, system
- actor_role
- tenant_id
- seller_scope
- domain_context
- permissions
- session_id
- request_id
- correlation_id nếu đang trong flow order/payment/provisioning

Rule:
- API không lấy tenant từ body nếu request đã có tenant context.
- Nếu domain và token conflict tenant, request bị reject và audit security event.

## 4. Domain-to-tenant mapping
### 4.1. Admin storefront
Admin có domain mặc định và có thể có brand riêng.

### 4.2. Reseller storefront
Mỗi custom domain map tới đúng một reseller tenant.

Trạng thái domain:
- requested
- pending_verification
- verified
- tls_pending
- active
- suspended
- removed

### 4.3. Xác minh domain
Flow:
1. Reseller khai báo domain.
2. Hệ thống sinh verification token.
3. Reseller thêm DNS TXT/CNAME.
4. Hệ thống/Admin xác minh.
5. TLS/certificate sẵn sàng.
6. Domain active.

Rule:
- Không cho một domain active ở nhiều tenant.
- Domain suspended không được route vào tenant.
- Domain mapping change phải audit.

## 5. Backend tenant enforcement
Mọi resource tenant-level cần check:
- resource.tenant_id == context.tenant_id
- hoặc actor là platform admin với emergency access hợp lệ.

Tài nguyên bắt buộc scope:
- users/clients
- wallets
- ledger entries
- payment/top-up requests
- catalog clone/tenant plans
- orders
- reservations
- provisioning jobs
- services
- credentials/access info
- invoices
- audit logs
- abuse/risk flags

## 6. Composite ownership rule
Data model nên thiết kế để dev khó query sai:
- wallets: tenant_id + wallet_id
- orders: tenant_id + order_id
- services: tenant_id + service_id
- ledger: tenant_id + ledger_entry_id
- users: tenant_id + user_id nếu user thuộc tenant

Khi lấy service, không chỉ `WHERE service_id = ?`; phải có tenant scope.

## 7. Role và permission phase 1
### 7.1. Platform Admin
Có toàn quyền platform, nhưng khi vào tenant reseller để hỗ trợ nên dùng emergency context.

P0 action:
- provider management
- master catalog
- reseller management
- payment verification platform-level
- refund/adjustment
- force suspend/terminate
- audit/report global

### 7.2. Reseller Owner
Có quyền tenant-level:
- white-label settings
- tenant pricing
- client/service management
- staff management
- tenant wallet/report/audit

### 7.3. Reseller Staff
Phase 1 dùng preset hạn chế:
- view dashboard
- view clients/services
- hỗ trợ thao tác vận hành cơ bản nếu Owner bật

Mặc định không có:
- refund
- manual wallet adjustment
- domain change
- white-label sensitive change
- pricing change
- staff permission change
- credential reveal cho service nhạy cảm nếu chưa bật quyền

### 7.4. Client
Chỉ được xem và thao tác tài nguyên của chính mình trong tenant hiện tại.

## 8. Emergency access
Admin emergency access phải có:
- reason bắt buộc
- time window
- actor_id
- target_tenant_id
- request_id
- audit action

Mặc định read-only. Write access cần action riêng và reason riêng.

Audit action:
- tenant.emergency.read
- tenant.emergency.write

## 9. Credential security
Credential/access info gồm:
- VPS username/password
- proxy endpoint/username/password/token
- VNC/console token
- private key
- provider token/API key

Rule:
- Encrypt at rest.
- Không log plaintext.
- Không lưu plaintext trong audit snapshot.
- UI masked by default.
- Reveal phải audit `credential.revealed`.
- Có thể yêu cầu re-auth/2FA trước khi reveal.
- Download/export credential nếu có phải audit riêng.

## 10. Provider secret management
Provider API key/token không nên lưu thô trong DB.

Rule đề xuất:
- Lưu trong secret vault hoặc encrypted config.
- Chỉ adapter runtime được quyền đọc.
- Rotate secret phải có audit.
- Không expose secret qua admin UI sau khi đã lưu; chỉ cho replace.

## 11. 2FA và session policy
P0:
- Admin bắt buộc 2FA.
- Reseller Owner bật mặc định/khuyến nghị bắt buộc.
- Staff có quyền tiền/credential/pricing nên bắt buộc 2FA.

Session:
- Admin/reseller owner session timeout ngắn hơn client.
- Action nhạy cảm yêu cầu re-auth/step-up auth.
- Login failed nhiều lần phải rate limit.

## 12. Rate limit P0
Áp dụng cho:
- login
- register
- password reset
- top-up request
- checkout
- renew
- change IP
- reinstall
- credential reveal
- domain verification

Rate limit nên có theo IP + actor + tenant.

## 13. Audit và security events
Cần audit:
- login success/failed
- tenant context mismatch
- permission denied
- credential reveal
- domain mapping change
- 2FA enable/disable
- staff permission change
- wallet/refund/adjustment
- emergency access

## 14. Acceptance criteria
- Client tenant A không thể gọi API lấy service tenant B.
- Staff không có quyền pricing không thể update tenant plan price dù biết API endpoint.
- Request body gửi tenant_id khác tenant context bị reject.
- Credential reveal luôn có audit.
- Admin emergency read có reason và log.
- Domain đã active ở tenant A không thể active ở tenant B.
- Provider API key không xuất hiện trong response/log/audit plaintext.
