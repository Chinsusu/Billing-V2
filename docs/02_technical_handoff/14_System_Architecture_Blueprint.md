# 14 - System Architecture Blueprint

## 1. Mục tiêu tài liệu

Tài liệu này mô tả kiến trúc tổng thể cho nền tảng VPS/Proxy hybrid multi-tenant, white-label reseller, wallet-first.

Mục tiêu là để team kỹ thuật thống nhất:
- Hệ thống gồm những khối nào.
- Dữ liệu đi qua các lớp nào.
- Tenant context được xác định và bảo vệ ở đâu.
- Billing, provisioning, inventory, credential và audit được đặt ở tầng nào.
- Module nào xử lý sync realtime, module nào xử lý bất đồng bộ qua worker/cron.
- Các nguyên tắc không được phá khi chọn framework hoặc ngôn ngữ triển khai.

Tài liệu này **chưa phải code**, không ràng buộc framework cụ thể.

---

## 2. Nguyên tắc kiến trúc nền

### 2.1 Wallet-first

Không tạo tài nguyên thật nếu chưa có trạng thái tiền hợp lệ.

Với client trực tiếp của platform:
```text
client_wallet đủ tiền -> debit wallet -> tạo order/provisioning job
```

Với client thuộc reseller:
```text
client_wallet đủ selling_price
và reseller_wallet đủ reseller_cost
-> debit cả 2 lớp theo rule settlement
-> tạo order/provisioning job
```

### 2.2 Tenant-first

Mọi request nghiệp vụ phải có tenant context hợp lệ trước khi đọc/ghi dữ liệu.

Không tin các field:
```text
tenant_id
seller_id
reseller_id
owner_id
```

nếu chúng đến từ request body/client. Các giá trị này phải lấy từ:
```text
domain mapping
authenticated session/token
server-side user membership
admin emergency access context
```

### 2.3 Ledger-first

Mọi thay đổi số dư phải đi qua ledger immutable.

Không có ledger entry thì xem như giao dịch không tồn tại.

### 2.4 Idempotency-first

Các thao tác có rủi ro lặp, đặc biệt là checkout/provisioning/payment webhook/retry, phải có idempotency key.

Cấm retry mù khi không biết provider đã tạo tài nguyên hay chưa.

### 2.5 Audit-by-default

Mọi thao tác liên quan tiền, credential, tenant, domain, quyền, provisioning, suspend/terminate phải ghi audit.

Audit không được chứa plaintext secret/credential.

### 2.6 Capability masking

UI và API không hiển thị/không cho gọi action nếu provider/source/service capability snapshot không hỗ trợ.

Ví dụ:
```text
Service từ source không hỗ trợ change_ip
-> UI không hiện nút change_ip
-> API vẫn phải chặn nếu user gọi trực tiếp endpoint change_ip
```

---

## 3. Kiến trúc tổng thể

```text
[ Admin Portal ]        [ Reseller Portal ]        [ Client Portal ]
       |                        |                         |
       +------------------------+-------------------------+
                                |
                         [ Web / API Layer ]
                                |
               [ Auth + Tenant Context + RBAC Guard ]
                                |
     +--------------------------+------------------------------+
     |                          |                              |
[ Core Domain Modules ]   [ Financial Core ]            [ Security Core ]
     |                          |                              |
     |                          |                              |
Catalog / Product         Wallet / Ledger                Audit Log
Order / Checkout          Top-up / Adjustment            Secret/Credential
Inventory / Reservation   Settlement                     Rate Limit
Service Lifecycle         Invoice Request                Abuse/Risk
                                |
                         [ Queue / Job Bus ]
                                |
                 +--------------+--------------+
                 |                             |
         [ Provisioning Worker ]       [ Cron/Scheduler ]
                 |                             |
          [ Provider Adapter Layer ]     Expiry/Renewal/Sync
                 |
       +---------+---------+----------+----------+
       |                   |                     |
   Proxmox              VPS API              Proxy Source
   OVH/Hetzner          Manual Source        Upstream Provider
```

---

## 4. Lớp portal

### 4.1 Admin Portal

Dành cho platform owner/team vận hành.

Phạm vi:
- Quản tenant/reseller.
- Quản master catalog.
- Quản provider/source.
- Quản top-up reseller và client trực tiếp.
- Xem toàn hệ thống, revenue, provider health, failed provisioning.
- Emergency access có reason và audit.

Không nên để Admin Portal dùng cùng routing logic với Client Portal mà thiếu guard. Admin có quyền cao nhưng phải bị audit nặng hơn.

### 4.2 Reseller Portal

Dành cho reseller owner/staff.

Phạm vi:
- Quản branding/storefront riêng.
- Quản client thuộc tenant của reseller.
- Clone sản phẩm từ master catalog.
- Set giá bán riêng.
- Duyệt top-up client nội bộ.
- Xem settlement cost/profit.
- Quản ticket/support/notification của tenant.

Reseller không được thấy:
```text
client của reseller khác
provider API key
master cost ngoài quyền được hiển thị
tenant data của platform
ledger platform-level không liên quan
```

### 4.3 Client Portal

Dành cho khách mua VPS/proxy.

Phạm vi:
- Đăng ký/đăng nhập.
- Nạp ví theo tenant.
- Mua/gia hạn dịch vụ.
- Xem dịch vụ, credential masked/reveal.
- Gửi support/request cancel.
- Xem invoice/transaction của chính mình.

Client không bao giờ được truyền tenant_id để tự chọn tenant. Tenant được xác định từ domain/session.

---

## 5. API layer

API layer chịu trách nhiệm:
- Xác thực user.
- Xác định tenant context.
- Kiểm tra RBAC/permission.
- Validate request.
- Gắn correlation_id/request_id.
- Ghi audit event cho action nhạy cảm.
- Không chứa secret provider trong response.

Mỗi request nghiệp vụ nên có context nội bộ:

```text
request_context:
- request_id
- actor_id
- actor_role
- actor_tenant_id
- actor_permissions
- domain_tenant_id
- effective_tenant_id
- is_emergency_access
- emergency_reason
- ip_address
- user_agent
```

Rule:
```text
effective_tenant_id phải do server xác định.
Nếu actor là platform admin, effective_tenant_id có thể là target tenant nhưng phải audit.
Nếu actor là reseller/client, effective_tenant_id phải bằng tenant của actor/domain.
```

---

## 6. Core domain modules

### 6.1 Catalog Module

Nhiệm vụ:
- Quản master product/plan/source.
- Clone catalog sang tenant.
- Lưu version/snapshot.
- Kiểm soát margin floor.
- Ẩn plan/source disabled/out-of-stock.

Không xử lý tiền trực tiếp. Catalog chỉ cung cấp giá/snapshot/capability để Checkout Module dùng.

### 6.2 Order & Checkout Module

Nhiệm vụ:
- Tạo order.
- Validate plan/source/capability.
- Reserve inventory.
- Debit wallet/settlement.
- Tạo provisioning job.
- Gắn order_status, payment_status, provisioning_status.

Checkout phải chạy theo transaction boundary rõ:
```text
validate -> reserve -> debit ledger -> create order/provisioning job
```

Nếu không thể hoàn thành một bước P0 thì phải rollback hoặc chuyển trạng thái manual_review có kiểm soát.

### 6.3 Inventory & Reservation Module

Nhiệm vụ:
- Quản capacity.
- Atomic reserve/release/allocate.
- Chống oversell.
- Quản reservation expiry.

Không cho provider worker tự trừ/tăng stock trực tiếp ngoài module này.

### 6.4 Service Lifecycle Module

Nhiệm vụ:
- Tạo service instance sau provisioning success.
- Quản active/suspended/expired/terminated.
- Quản renew/cancel/refund lifecycle.
- Gắn lifecycle event.

Mọi thay đổi service_status phải có transition hợp lệ và audit.

### 6.5 Provisioning Module

Nhiệm vụ:
- Tạo provisioning job.
- Gọi provider adapter qua worker.
- Xử lý retry/partial success/manual review.
- Gắn external_resource_id.
- Đồng bộ status từ provider.

API request của user không nên gọi provider trực tiếp trong cùng request lâu. Nên tạo job và worker xử lý.

---

## 7. Financial Core

### 7.1 Wallet Module

Nhiệm vụ:
- Quản wallet theo user/tenant/reseller.
- Tính available balance từ ledger.
- Không sửa balance bằng tay ngoài adjustment có audit.
- Phân biệt client wallet và reseller settlement wallet.

### 7.2 Ledger Module

Ledger entry là bất biến.

Thông tin tối thiểu:
```text
ledger_id
wallet_id
tenant_id
actor_id
entry_type
direction
amount
currency
balance_after
reference_type
reference_id
correlation_id
created_at
```

Không update/delete ledger entry sau khi posted.

### 7.3 Settlement Module

Nhiệm vụ:
- Khi client thuộc reseller mua hàng, ghi nhận selling_price, reseller_cost, margin.
- Debit đúng wallet.
- Tạo report profit theo snapshot.
- Xử lý refund/rollback theo policy.

---

## 8. Security Core

### 8.1 Credential Service

Nhiệm vụ:
- Encrypt credential at rest.
- Mask credential khi hiển thị.
- Ghi audit khi reveal.
- Không log plaintext.
- Quản rotation/reissue nếu provider hỗ trợ.

Các field nhạy cảm:
```text
provider_api_key
provider_secret
vps_root_password
proxy_username
proxy_password
ssh_private_key
console_url_token
```

### 8.2 Audit Service

Audit service nhận event từ API và worker.

Audit event phải có:
```text
actor
tenant
target
action
before/after redacted
request_id/correlation_id
ip/user_agent nếu từ request
worker_job_id nếu từ worker
```

### 8.3 Risk/Abuse Service

Phase 1 chưa cần scoring phức tạp, nhưng phải có:
- risk flag.
- abuse case.
- manual review.
- suspend reason.
- blacklist tối thiểu.

---

## 9. Queue và worker

Các thao tác nên đi qua queue:
- Provision VPS/proxy.
- Retry provisioning.
- Sync provider resource.
- Send notification.
- Generate report/export.
- Expiry reminder.
- Suspend/terminate due to expiry.
- Provider health check.

Queue message phải có:
```text
job_id
job_type
tenant_id
reference_type
reference_id
idempotency_key
attempt_count
correlation_id
created_at
```

Worker phải ghi:
```text
started_at
finished_at
status
provider_response_summary redacted
error_code
next_retry_at
manual_review_required
```

---

## 10. Cron/Scheduler

Các job định kỳ P0:
- reservation_expiry_job.
- service_expiry_job.
- suspension_job.
- termination_job.
- provider_health_check_job.
- provider_inventory_sync_job.
- renewal_reminder_job.
- pending_topup_review_reminder_job.
- audit_retention_job.

Cron job phải idempotent: chạy lại không được làm sai số dư, stock hoặc trạng thái.

---

## 11. Data boundary

### 11.1 Bảng thuộc tenant

Các bảng này bắt buộc có `tenant_id`:
```text
users
wallets
orders
order_items
reservations
services
service_credentials
tenant_products
tenant_plans
topup_requests
support_tickets
audit_logs
risk_flags
abuse_cases
notifications
```

### 11.2 Bảng platform-level

Các bảng platform-level:
```text
master_products
master_plans
provider_sources
provider_accounts
platform_settings
exchange_rates
system_jobs
```

Bảng platform-level vẫn cần audit khi sửa.

---

## 12. Observability

Mọi request/job quan trọng cần có trace:

```text
correlation_id:
topup -> ledger -> order -> reservation -> provisioning_job -> provider_request -> service -> notification
```

Log không được chứa:
```text
plaintext credential
provider secret
payment proof sensitive image raw URL nếu private
full token/session
```

Metric tối thiểu:
- checkout success rate.
- provisioning success/fail/manual_review rate.
- provider API latency/error rate.
- reservation expired count.
- wallet adjustment count.
- credential reveal count.
- abuse suspension count.
- top-up approval time.
- reseller low balance count.

---

## 13. Deployment logical units

Có thể triển khai mono-repo hoặc multi-service, nhưng logical unit nên giữ rõ:

```text
web_frontend
backend_api
worker
scheduler
database
cache/queue
object_storage
logging/monitoring
secret_manager
```

Không bắt buộc tách microservice phase 1. Một monolith có module boundary rõ + worker riêng thường thực dụng hơn.

Khuyến nghị phase 1:
```text
modular monolith + queue worker + scheduler
```

Lý do:
- Nhanh ra hàng.
- Ít overhead vận hành.
- Dễ giữ transaction billing/order/inventory.
- Sau này tách module khi tải thật sự lớn.

---

## 14. Nguyên tắc thất bại an toàn

Nếu hệ thống không chắc chắn, ưu tiên:
```text
không provision hơn provision sai
manual_review hơn retry mù
ẩn credential hơn lộ credential
không trừ tiền hơn trừ tiền không trace được
suspend có reason hơn terminate vội
```

Câu nền: **hệ thống hạ tầng không cần “ảo thuật nhanh”; nó cần không làm sai khi provider, payment, network và người dùng cùng gây nhiễu.**
