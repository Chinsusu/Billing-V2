# 16 - API Contract And Permission Spec

## 1. Mục tiêu tài liệu

Tài liệu này định nghĩa API contract logic cho backend/frontend. Đây là hợp đồng hành vi, không phải code.

Mỗi API cần rõ:
- Ai được gọi.
- Tenant scope lấy từ đâu.
- Request/response cần gì.
- Validate gì.
- Ghi audit gì.
- Error code nào trả về.
- Rate limit nào áp dụng.
- Có cần idempotency hay không.

---

## 2. API conventions

### 2.1 Base principles

```text
Không tin tenant_id từ body.
Không trả plaintext secret nếu không qua reveal action.
Không trả provider API key qua API thường.
Không cho action nếu capability snapshot không hỗ trợ.
Không tạo tài nguyên nếu chưa debit/settlement hợp lệ.
```

### 2.2 Request headers logic

Frontend nên gửi:
```text
Authorization
Idempotency-Key với action tạo giao dịch
X-Request-Id nếu có
```

Backend tự gắn:
```text
correlation_id
effective_tenant_id
actor_context
domain_context
```

### 2.3 Standard response

Thành công:
```text
{
  "success": true,
  "data": {},
  "request_id": "...",
  "correlation_id": "..."
}
```

Lỗi:
```text
{
  "success": false,
  "error": {
    "code": "INSUFFICIENT_BALANCE",
    "message": "Human readable message",
    "details": {}
  },
  "request_id": "...",
  "correlation_id": "..."
}
```

### 2.4 Pagination

List API dùng:
```text
page
page_size
sort
filter
```

Response:
```text
items
pagination.total
pagination.page
pagination.page_size
```

### 2.5 Standard error codes

| Code | Meaning |
|---|---|
| `UNAUTHENTICATED` | Chưa đăng nhập |
| `FORBIDDEN` | Không có quyền |
| `TENANT_MISMATCH` | Resource không thuộc tenant |
| `VALIDATION_ERROR` | Sai input |
| `RESOURCE_NOT_FOUND` | Không thấy hoặc không được thấy resource |
| `RATE_LIMITED` | Vượt giới hạn |
| `IDEMPOTENCY_CONFLICT` | Idempotency key đã dùng cho payload khác |
| `INSUFFICIENT_CLIENT_BALANCE` | Ví client thiếu tiền |
| `INSUFFICIENT_RESELLER_BALANCE` | Ví reseller settlement thiếu tiền |
| `OUT_OF_STOCK` | Hết stock |
| `PLAN_DISABLED` | Plan bị tắt |
| `PROVIDER_UNAVAILABLE` | Provider/source không sẵn sàng |
| `CAPABILITY_NOT_SUPPORTED` | Action không được source hỗ trợ |
| `PROVISIONING_MANUAL_REVIEW` | Cần review thủ công |
| `CREDENTIAL_REVEAL_DENIED` | Không được xem credential |

---

## 3. Permission model per API

Mỗi endpoint phải khai báo:
```text
required_permission
allowed_roles
tenant_scope
risk_level
audit_action
```

Risk level:
```text
low: đọc list thông thường
medium: tạo order, cập nhật profile
high: duyệt top-up, suspend service, reveal credential
critical: chỉnh ledger, đổi provider, emergency access, terminate service
```

Critical action phải có:
```text
2FA nếu role yêu cầu
audit bắt buộc
reason nếu manual/admin action
correlation_id
```

---

## 4. Auth APIs

### 4.1 Register client

```text
POST /auth/register
```

Allowed:
```text
public on tenant storefront
```

Tenant scope:
```text
derived from domain
```

Request:
```text
email
password
full_name
accept_terms
```

Validation:
- domain phải map tenant active.
- email chưa tồn tại trong tenant.
- password đạt policy.
- accept_terms = true.

Response:
```text
user_id
email
status
email_verification_required
```

Audit:
```text
auth.registered
```

Errors:
```text
TENANT_NOT_ACTIVE
EMAIL_ALREADY_EXISTS
WEAK_PASSWORD
TERMS_NOT_ACCEPTED
```

### 4.2 Login

```text
POST /auth/login
```

Validation:
- tenant/domain hợp lệ.
- user active.
- password đúng.
- 2FA nếu required/enabled.

Audit:
```text
auth.login.success
auth.login.failed
auth.2fa.challenge
```

Rate limit:
```text
per IP + per email + per tenant
```

### 4.3 Logout

```text
POST /auth/logout
```

Audit:
```text
auth.logout
```

---

## 5. Tenant and domain APIs

### 5.1 Admin create reseller tenant

```text
POST /admin/tenants
```

Allowed:
```text
platform_super_admin
platform_staff with tenant.create
```

Request:
```text
name
slug
owner_email
default_currency
timezone
initial_status
```

Validation:
- slug unique.
- owner email valid.
- actor has permission.

Audit:
```text
tenant.created
```

### 5.2 Reseller update branding

```text
PATCH /reseller/tenant/branding
```

Allowed:
```text
reseller_owner
reseller_staff with tenant.branding.update
```

Request:
```text
logo_url
brand_color
support_email
support_telegram
footer_text
```

Validation:
- tenant active.
- assets allowed.
- no unsafe links.

Audit:
```text
tenant.branding.updated
```

### 5.3 Add custom domain

```text
POST /reseller/tenant/domains
```

Request:
```text
domain
```

Validation:
- domain chưa thuộc tenant khác.
- domain hợp lệ.
- tạo verification token.

Response:
```text
domain_id
verification_record_type
verification_record_value
verification_status
```

Audit:
```text
tenant.domain.created
```

---

## 6. Catalog APIs

### 6.1 Admin create master product

```text
POST /admin/catalog/products
```

Allowed:
```text
catalog.product.create
```

Request:
```text
product_type
name
description
status
display_order
```

Audit:
```text
catalog.product.created
```

### 6.2 Admin create master plan

```text
POST /admin/catalog/plans
```

Request:
```text
product_id
plan_code
name
specs
billing_cycle_type
billing_cycle_value
base_cost
suggested_price
reseller_min_price
status
```

Validation:
- product active/draft.
- billing cycle valid.
- prices >= 0.
- plan_code unique for version.

Audit:
```text
catalog.plan.created
```

### 6.3 Reseller list available master plans

```text
GET /reseller/catalog/master-plans
```

Allowed:
```text
reseller_owner
reseller_staff with catalog.view
```

Response:
```text
plans visible to reseller
suggested_price
reseller_cost
margin hints
capabilities
stock status
```

Do not expose:
```text
provider_api_key
sensitive provider account detail
```

### 6.4 Reseller clone/sync plan

```text
POST /reseller/catalog/plans/clone
```

Request:
```text
master_plan_id
selling_price
visibility
```

Validation:
- master plan active.
- selling_price >= reseller_cost or margin policy allows warning.
- tenant allowed to sell this product.
- source available.

Audit:
```text
catalog.tenant_plan.cloned
```

### 6.5 Client list catalog

```text
GET /client/catalog
```

Allowed:
```text
public or authenticated depending tenant setting
```

Tenant scope:
```text
domain/session tenant
```

Response:
```text
products
plans
selling_price
public specs
stock availability summary
```

Do not expose:
```text
reseller_cost
provider source internal ID unless safe alias
internal margin
```

---

## 7. Wallet APIs

### 7.1 Client get wallet

```text
GET /client/wallet
```

Allowed:
```text
client
```

Response:
```text
currency
available_balance
locked_balance
recent_ledger_entries
pending_topups
```

Audit:
```text
optional wallet.viewed
```

### 7.2 Submit top-up request

```text
POST /client/wallet/topups
```

Request:
```text
amount
currency
payment_method
payment_reference
proof_attachment_id
```

Validation:
- amount >= min topup.
- currency allowed.
- method enabled.
- proof required if method requires proof.

Status:
```text
submitted
```

Audit:
```text
wallet.topup.submitted
```

### 7.3 Approve client top-up

```text
POST /reseller/wallet/topups/{topup_request_id}/approve
```

Allowed:
```text
reseller_owner
reseller_staff with wallet.topup.approve
platform_admin for emergency/support
```

Validation:
- top-up belongs to actor tenant.
- status submitted/under_review.
- amount/currency match actual payment.
- reviewer has 2FA if required.

Result:
- create ledger credit to client wallet.
- set topup_request approved.
- notify client.

Audit:
```text
wallet.topup.approved
wallet.ledger.posted
```

Errors:
```text
TOPUP_ALREADY_REVIEWED
PERMISSION_DENIED
```

### 7.4 Admin approve reseller top-up

```text
POST /admin/resellers/{tenant_id}/wallet/topups/{topup_request_id}/approve
```

Allowed:
```text
platform_super_admin
finance_agent
```

Result:
- credit reseller settlement wallet.
- audit high risk.

---

## 8. Checkout and order APIs

### 8.1 Create checkout order

```text
POST /client/orders
```

Allowed:
```text
client
```

Idempotency:
```text
required
```

Request:
```text
tenant_plan_id
quantity
billing_cycle
coupon_code optional
```

Validation:
- user active.
- tenant active.
- tenant_plan belongs to tenant.
- plan visible and active.
- source active.
- stock available.
- client wallet >= selling_price.
- if tenant is reseller: reseller settlement wallet >= reseller_cost.
- abuse/risk check passes or routes manual_review.
- idempotency key not reused with different payload.

Success flow:
```text
create order
create order_item snapshots
reserve inventory
debit client wallet
debit reseller wallet if reseller tenant
create provisioning_job
order_status = provisioning
provisioning_status = queued
```

Response:
```text
order_id
order_number
order_status
provisioning_status
estimated_activation_message
```

Audit:
```text
order.created
reservation.created
wallet.client.debited
wallet.reseller_cost.debited
provisioning.job.created
```

Errors:
```text
INSUFFICIENT_CLIENT_BALANCE
INSUFFICIENT_RESELLER_BALANCE
OUT_OF_STOCK
PLAN_DISABLED
PROVIDER_UNAVAILABLE
RISK_MANUAL_REVIEW_REQUIRED
```

### 8.2 Get order detail

```text
GET /client/orders/{order_id}
```

Scope:
```text
client can view own order only
reseller can view orders inside tenant
admin can view all with permission
```

Tenant mismatch returns:
```text
404 or FORBIDDEN depending security policy
```

### 8.3 Cancel pending order

```text
POST /client/orders/{order_id}/cancel
```

Allowed only if:
```text
order_status in draft, pending_payment, manual_review_before_provision
provisioning not started
reservation still reserved
```

Result:
```text
release reservation
reverse debit if any
order_status = cancelled
```

Audit:
```text
order.cancelled
reservation.released
wallet.reversal.posted
```

---

## 9. Service APIs

### 9.1 Client list services

```text
GET /client/services
```

Filters:
```text
status
service_type
expiring_before
```

Response:
```text
service_id
name/spec summary
status
expiry
public endpoint masked
actions available by capability_snapshot
```

### 9.2 Client service detail

```text
GET /client/services/{service_id}
```

Response:
```text
service info
billing info
capability actions
credential masked_hint only
lifecycle events public subset
```

Never return plaintext credential here.

### 9.3 Reveal credential

```text
POST /client/services/{service_id}/credentials/{credential_id}/reveal
```

Allowed:
```text
client owns service
reseller owner/staff with credential.reveal
platform admin with credential.reveal
```

Validation:
- service belongs to tenant.
- credential belongs to service.
- actor has permission.
- 2FA required for staff/admin/reseller owner if policy requires.
- rate limit reveal.

Response:
```text
plaintext credential once
masked hint
reveal_expires_message
```

Audit:
```text
credential.revealed
```

Security:
- response not cached.
- no plaintext in audit/log.
- optional require re-auth for high-risk roles.

### 9.4 Renew service

```text
POST /client/services/{service_id}/renew
```

Idempotency:
```text
required
```

Validation:
- service active/expired/grace depending policy.
- plan still renewable.
- wallet sufficient.
- reseller wallet sufficient if reseller tenant.
- renew cycle valid.
- not terminated.

Result:
```text
debit wallet(s)
extend term_end_at based on policy
create lifecycle event service.renewed
if provider requires renew call -> create provisioning_job type renew
```

Audit:
```text
service.renewed
wallet.client.debited
wallet.reseller_cost.debited
```

### 9.5 Cancel service / request termination

```text
POST /client/services/{service_id}/cancel-request
```

Request:
```text
reason
cancel_at: immediate/end_of_term
```

Validation:
- service belongs to client.
- policy allows cancel.
- if immediate termination affects refund, route review if required.

Audit:
```text
service.cancel_requested
```

### 9.6 Admin/reseller suspend service

```text
POST /admin/services/{service_id}/suspend
POST /reseller/services/{service_id}/suspend
```

Request:
```text
reason
notify_client
```

Validation:
- reason required.
- permission required.
- source supports suspend or internal suspension only.
- if provider action required, create job.

Audit:
```text
service.suspended
provisioning.job.created optional
```

---

## 10. Provider and provisioning APIs

### 10.1 Admin list provider sources

```text
GET /admin/provider-sources
```

Allowed:
```text
provider.view
```

Response:
```text
source status
health status
inventory summary
capability profile
last sync
```

No provider secrets.

### 10.2 Admin create/update provider source

```text
POST /admin/provider-sources
PATCH /admin/provider-sources/{source_id}
```

Allowed:
```text
provider.manage
```

Validation:
- provider account credentials stored encrypted.
- capability profile valid.
- inventory policy valid.

Audit:
```text
provider.source.created
provider.source.updated
```

### 10.3 Manual retry provisioning job

```text
POST /admin/provisioning-jobs/{job_id}/retry
```

Allowed:
```text
provisioning.retry
```

Validation:
- job status failed/manual_review.
- retry_safety != do_not_retry.
- if unsafe_retry, require explicit override reason and high permission.
- max attempts not exceeded unless override.

Audit:
```text
provisioning.job.retry_requested
```

### 10.4 Mark manual review resolved

```text
POST /admin/provisioning-jobs/{job_id}/resolve
```

Request:
```text
resolution_type: success_failed_cancelled_link_existing_resource
external_resource_id optional
service_id optional
note required
```

Validation:
- note required.
- if linking resource, external_resource_id unique.
- ledger/reservation/service consistency checked.

Audit:
```text
provisioning.manual_review.resolved
```

---

## 11. Audit and report APIs

### 11.1 List audit logs

```text
GET /admin/audit-logs
GET /reseller/audit-logs
```

Scope:
```text
admin can view platform/global with permission
reseller can view own tenant
client generally cannot view raw audit logs
```

Filters:
```text
actor
action
target_type
target_id
date range
correlation_id
```

Redaction:
```text
always redacted
```

### 11.2 Reseller profit report

```text
GET /reseller/reports/profit
```

Calculation:
```text
gross_revenue = sum(client selling_price posted)
platform_cost = sum(reseller_cost posted)
gross_profit = gross_revenue - platform_cost - refunds/adjustments
```

Must use:
```text
ledger/order snapshots, not current plan price
```

### 11.3 Admin provider health report

```text
GET /admin/reports/provider-health
```

Includes:
```text
source status
success rate
failure rate
manual_review count
avg provisioning time
out_of_stock count
```

---

## 12. Abuse and risk APIs

### 12.1 Create manual risk flag

```text
POST /admin/risk-flags
POST /reseller/risk-flags
```

Request:
```text
user_id/service_id/order_id
flag_type
severity
note
```

Audit:
```text
risk.flag.created
```

### 12.2 Create abuse case

```text
POST /admin/abuse-cases
```

Request:
```text
service_id
case_type
severity
evidence_private
recommended_action
```

Actions:
```text
warning
suspend
terminate
provider_notice
```

Audit:
```text
abuse.case.created
```

---

## 13. Rate limit matrix

| Action | Limit logic |
|---|---|
| login | per IP + email + tenant |
| register | per IP + domain |
| top-up submit | per user + tenant |
| checkout | per user + tenant + plan |
| credential reveal | per user + service |
| reset password/change IP | per service + source |
| admin retry provisioning | per admin + job |
| domain verification | per tenant + domain |
| support message | per user + ticket |

---

## 14. API acceptance criteria

API spec được xem là đạt khi:
- Endpoint nào cũng có role/permission rõ.
- Endpoint nào cũng có tenant scope rõ.
- Checkout có idempotency bắt buộc.
- Reveal credential là action riêng, không nằm trong service detail.
- Reseller/client không thể thấy dữ liệu tenant khác.
- Admin action critical yêu cầu audit + reason.
- Error codes thống nhất để frontend hiển thị đúng.
- Provider/provisioning retry không cho unsafe retry vô thức.
- Report dùng snapshot/ledger, không dùng giá hiện tại.
- API response không chứa provider secret hoặc credential plaintext ngoài reveal endpoint.

Câu nền: **API không chỉ là đường truyền dữ liệu; nó là cổng kiểm soát quyền, tiền, tenant và rủi ro.**
