# 17 - RBAC Permission Matrix

## 1. Mục tiêu tài liệu

Tài liệu này định nghĩa quyền truy cập theo vai trò cho platform VPS/Proxy multi-tenant.

RBAC phải bảo vệ 5 thứ:
```text
tiền
credential
tenant data
provider/source
hành động phá hủy như suspend/terminate
```

Không được dựa vào frontend để ẩn nút là đủ. Backend bắt buộc kiểm tra permission.

---

## 2. Role chuẩn phase 1

### 2.1 Platform Super Admin

Quyền cao nhất trên platform.

Có thể:
- Quản mọi tenant.
- Quản master catalog.
- Quản provider/source.
- Duyệt top-up reseller.
- Xem report toàn hệ thống.
- Emergency access tenant.

Yêu cầu:
```text
2FA bắt buộc
mọi critical action có audit
không dùng cho tác vụ hằng ngày nếu có role thấp hơn
```

### 2.2 Platform Admin / Staff

Dành cho nhân sự vận hành.

Có thể được cấp theo permission:
- support.
- finance.
- catalog.
- provisioning.
- abuse.
- read-only audit.

Không mặc định có toàn quyền.

### 2.3 Finance Agent

Chuyên xử lý top-up, refund, adjustment.

Có thể:
- xem top-up request.
- approve/reject top-up theo scope.
- tạo adjustment nếu được cấp.
- xem ledger/report.

Không nên có quyền:
- xem credential.
- sửa provider.
- terminate service.

### 2.4 Support Agent

Có thể:
- xem user/order/service metadata.
- hỗ trợ ticket.
- tạo note.
- request suspend/unsuspend theo workflow.

Không mặc định có quyền:
- reveal credential.
- approve top-up.
- adjustment ledger.
- đổi provider source.

### 2.5 Provisioning Operator

Có thể:
- xem failed jobs.
- retry safe jobs.
- resolve manual review.
- link existing provider resource nếu có quyền cao.

Không nên có quyền:
- duyệt tiền.
- chỉnh catalog price.
- xem toàn bộ ledger nếu không cần.

### 2.6 Reseller Owner

Owner của reseller tenant.

Có thể:
- quản branding/domain.
- quản staff.
- quản client trong tenant.
- clone catalog/set giá.
- duyệt top-up client.
- xem report profit tenant.
- suspend service trong tenant theo policy.

Yêu cầu:
```text
2FA bật mặc định, khuyến nghị bắt buộc.
critical action cần audit.
```

### 2.7 Reseller Staff

Nhân sự của reseller.

Quyền theo role con:
- sales.
- support.
- finance.
- catalog manager.
- read-only.

Không được có quyền owner mặc định.

### 2.8 Client

Khách mua dịch vụ.

Có thể:
- xem/mua/gia hạn dịch vụ của mình.
- xem ledger/top-up của mình.
- reveal credential của service thuộc mình.
- gửi support/cancel request.

Không thể:
- xem client khác.
- xem reseller cost.
- chỉnh giá.
- duyệt top-up.
- xem provider/source nội bộ.

### 2.9 Read-only Auditor

Có thể xem dữ liệu phục vụ kiểm toán nhưng không sửa.

Không được:
- reveal credential.
- approve tiền.
- retry provisioning.
- sửa tenant/domain/provider.

---

## 3. Permission naming convention

Format:
```text
module.resource.action
```

Ví dụ:
```text
tenant.view
tenant.create
tenant.update
tenant.domain.manage

catalog.master.view
catalog.master.create
catalog.master.update
catalog.tenant.clone
catalog.tenant.price_update

wallet.view
wallet.topup.submit
wallet.topup.approve
wallet.adjustment.create
wallet.ledger.view

order.view
order.create
order.cancel
order.refund

service.view
service.credential.reveal
service.renew
service.suspend
service.unsuspend
service.terminate

provider.view
provider.manage
provisioning.job.view
provisioning.job.retry
provisioning.manual_review.resolve

audit.view
report.view
risk.flag.manage
abuse.case.manage
```

---

## 4. Permission risk level

| Risk level | Examples | Control |
|---|---|---|
| Low | view catalog, view own order | auth + tenant scope |
| Medium | create order, submit top-up, renew service | auth + validation + audit optional |
| High | approve top-up, reveal credential, suspend service | permission + audit + rate limit |
| Critical | ledger adjustment, provider manage, terminate, emergency access | 2FA + reason + audit + restricted role |

---

## 5. Matrix tổng quan

| Action | Super Admin | Platform Staff | Finance | Support | Provisioning | Reseller Owner | Reseller Staff | Client | Auditor |
|---|---:|---:|---:|---:|---:|---:|---:|---:|---:|
| Xem tất cả tenant | Yes | Permission | No | Permission-limited | No | No | No | No | Read-only |
| Tạo reseller tenant | Yes | Permission | No | No | No | No | No | No | No |
| Sửa branding tenant | Yes | Permission | No | No | No | Own tenant | Permission | No | No |
| Quản custom domain | Yes | Permission | No | No | No | Own tenant | Permission | No | Read-only |
| Tạo master product/plan | Yes | Permission | No | No | No | No | No | No | Read-only |
| Clone plan về reseller | Yes | Permission | No | No | No | Yes | Permission | No | Read-only |
| Sửa giá bán reseller | Yes | Permission | No | No | No | Yes | Permission | No | Read-only |
| Xem reseller cost | Yes | Permission | Finance view | No | No | Own tenant | Permission | No | Read-only |
| Xem client wallet | Yes | Permission | Permission | Permission | No | Own tenant | Permission | Own wallet | Read-only |
| Submit top-up | No | No | No | No | No | For own reseller wallet | No | Own wallet | No |
| Approve client top-up | Yes | Permission | Permission | No | No | Own tenant | Permission | No | Read-only |
| Approve reseller top-up | Yes | Permission | Yes | No | No | No | No | No | Read-only |
| Ledger adjustment | Yes | Permission critical | Permission critical | No | No | Optional own tenant only | No | No | Read-only |
| Create order | No | No | No | No | No | For client via admin action optional | No | Own account | No |
| View order | Yes | Permission | Finance-related | Support-related | Provisioning-related | Own tenant | Permission | Own order | Read-only |
| Cancel pending order | Yes | Permission | No | Permission | No | Own tenant | Permission | Own order if policy | No |
| View service | Yes | Permission | No/limited | Permission | Permission | Own tenant | Permission | Own service | Read-only |
| Reveal credential | Yes + audit | Permission + audit | No | Usually No | Optional + audit | Own tenant + audit | Permission + audit | Own service + audit | No |
| Renew service | Yes + audit | Permission + audit | No | No | No | Own tenant + audit | Permission + audit | Own service + wallet | No |
| Suspend service | Yes | Permission | No | Request/Permission | Permission if technical | Own tenant | Permission | No | No |
| Terminate service | Yes | Critical permission | No | No | Critical permission | Own tenant if policy | Permission critical | No | No |
| Retry provisioning job | Yes | Permission | No | No | Yes | No | No | No | Read-only |
| Resolve manual provisioning | Yes | Permission | No | No | Yes | No | No | No | Read-only |
| Manage provider source | Yes | Permission critical | No | No | No | No | No | No | Read-only |
| View audit logs | Yes | Permission | Finance scope | Support scope | Provisioning scope | Own tenant | Permission | No | Read-only |
| Manage abuse case | Yes | Permission | No | Permission | No | Own tenant limited | Permission | No | Read-only |

---

## 6. Data visibility rules

### 6.1 Platform admin

Can see:
```text
all tenants
all users
all orders
all services
all provider health
all financial reports
```

But:
- Credential reveal still requires explicit action.
- Provider secrets hidden by default.
- Emergency access into tenant must have reason.

### 6.2 Reseller

Can see:
```text
tenant clients
tenant orders
tenant services
tenant wallet/client wallet
tenant catalog clones
tenant reports
```

Cannot see:
```text
other tenant data
platform provider secret
global profit except own tenant
master cost beyond reseller cost exposed to them
```

### 6.3 Client

Can see:
```text
own profile
own wallet
own order
own service
own notifications
own ticket
```

Cannot see:
```text
reseller margin
provider source details
other users
raw audit logs
```

---

## 7. Emergency access policy

Emergency access là khi platform admin/staff truy cập dữ liệu tenant để hỗ trợ/sửa sự cố.

Required:
```text
permission: tenant.emergency_access
2FA: required
reason: required
target_tenant_id: required
time-bound session: recommended
audit action: tenant.emergency_access.started / ended
```

Forbidden:
```text
không dùng emergency access để duyệt tiền không có chứng từ
không reveal credential nếu không có reason rõ
không chỉnh ledger nếu không có finance permission
```

Audit metadata:
```text
actor_id
target_tenant_id
reason
started_at
ended_at
actions performed
correlation_id
```

---

## 8. Credential reveal policy

Credential reveal là action high risk.

Allowed by default:
```text
client owns service
reseller owner for own tenant
platform super admin
staff with service.credential.reveal
```

Controls:
```text
2FA for admin/reseller/staff
rate limit
audit every reveal
mask by default
no plaintext in log
```

Optional policy:
```text
support staff không được reveal, chỉ được reset password nếu provider hỗ trợ.
```

---

## 9. Finance permission policy

Top-up approval và ledger adjustment phải tách quyền.

### 9.1 Approve top-up

Yêu cầu:
```text
wallet.topup.approve
scope: client tenant hoặc reseller wallet
audit
review note recommended
```

### 9.2 Ledger adjustment

Yêu cầu:
```text
wallet.adjustment.create
reason required
supporting reference required
2FA required
audit critical
```

Không cho cùng một staff vừa tạo adjustment vừa approve nếu cần four-eyes policy phase sau.

---

## 10. Provider permission policy

Provider/source là tài sản platform-level.

Actions critical:
```text
provider.manage
provider.credentials.update
provider.source.disable
provider.source.delete/archive
```

Required:
```text
2FA
audit
reason
change summary
```

Không cấp provider.manage cho reseller.

---

## 11. Service destructive action policy

### 11.1 Suspend

Required:
```text
reason
notify_client decision
audit
provider capability check
```

### 11.2 Unsuspend

Required:
```text
reason
verify billing/abuse resolved
audit
```

### 11.3 Terminate

Required:
```text
critical permission
reason
confirmation
audit
check refund/cancel policy
provider terminate capability
```

Terminated service không được renew.

---

## 12. Role templates phase 1

### 12.1 Platform Finance

```text
wallet.view
wallet.topup.approve
wallet.ledger.view
report.finance.view
audit.view.finance
```

### 12.2 Platform Support

```text
user.view
order.view
service.view
ticket.manage
audit.view.support
risk.flag.create
```

### 12.3 Platform Provisioning

```text
provider.view
provisioning.job.view
provisioning.job.retry
provisioning.manual_review.resolve
service.view
audit.view.provisioning
```

### 12.4 Reseller Finance

```text
wallet.view
wallet.topup.approve
wallet.ledger.view
report.tenant_profit.view
```

### 12.5 Reseller Support

```text
client.view
order.view
service.view
ticket.manage
risk.flag.create
```

### 12.6 Reseller Catalog Manager

```text
catalog.master.view
catalog.tenant.clone
catalog.tenant.price_update
catalog.tenant.visibility_update
```

---

## 13. RBAC acceptance criteria

RBAC được xem là đạt khi:
- Backend check permission ở mọi API nhạy cảm.
- Tenant scope luôn được kiểm tra cùng permission.
- Client không thể truy cập service/order/wallet của user khác.
- Reseller không thể truy cập tenant khác.
- Staff không có quyền mặc định ngoài role được cấp.
- Credential reveal luôn tạo audit.
- Ledger adjustment cần reason + permission critical.
- Admin emergency access có reason + audit.
- UI ẩn nút theo permission, nhưng backend vẫn chặn nếu gọi trực tiếp.
- Test case cross-tenant access phải pass.

Câu nền: **quyền không phải để làm khó người dùng; quyền là ranh giới giữa platform có thể mở rộng và platform tự đâm xuyên dữ liệu của chính mình.**
