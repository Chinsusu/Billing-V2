# 20 - UI Wireflow And Screen Spec

## 1. Mục tiêu tài liệu

Tài liệu này định nghĩa wireflow và screen spec cho 3 portal:
```text
Admin Portal
Reseller Portal
Client Portal
```

Đây không phải tài liệu thiết kế mỹ thuật/Figma. Mục tiêu là khóa:
- Màn nào cần có.
- Ai được xem.
- Dữ liệu hiển thị.
- Action nào có trên màn.
- Empty/error/loading state.
- Audit/notification liên quan.
- Security rule đặc biệt.

---

## 2. Nguyên tắc UI

### 2.1 UI phải theo permission và capability

Frontend phải ẩn action nếu:
```text
user không có permission
service/source capability không hỗ trợ
service status không cho action
tenant policy không cho action
```

Backend vẫn phải chặn nếu user gọi API trực tiếp.

### 2.2 Credential luôn masked mặc định

Service detail không hiển thị plaintext credential.

Credential flow:
```text
masked -> click Reveal -> confirm/2FA nếu cần -> API reveal -> hiển thị tạm thời -> audit
```

### 2.3 Money state phải rõ

Các màn ví/order phải hiển thị:
```text
available balance
pending top-up
ledger history
order payment status
refund/reversal nếu có
```

### 2.4 Manual review phải dễ thấy

Admin/reseller nên có dashboard cảnh báo:
```text
failed provisioning
manual review
pending top-up
provider down
abuse case
low reseller balance
```

---

## 3. Client Portal wireflow

### 3.1 Public storefront

Path logic:
```text
/
 /products
 /products/{product}
 /login
 /register
```

Data:
- Brand logo/theme theo tenant.
- Product categories: VPS, Proxy.
- Plan cards.
- Price.
- Specs.
- Stock status summary.
- Billing cycle.
- Terms/AUP link.

Actions:
- Register.
- Login.
- Buy now.

Empty state:
```text
No active products available.
```

Error state:
```text
Tenant disabled/domain not verified.
```

Security:
- Không hiển thị reseller_cost/provider internals.

### 3.2 Client dashboard

Purpose:
```text
tổng quan ví, order, service, expiry
```

Data:
- Wallet balance.
- Active services count.
- Expiring soon services.
- Pending orders.
- Recent notifications.

Actions:
- Top up wallet.
- Buy service.
- Renew expiring service.
- Open support.

### 3.3 Client wallet screen

Data:
```text
available balance
pending top-up requests
ledger entries
top-up instructions
```

Actions:
- Create top-up request.
- Upload proof.
- Cancel draft/submitted top-up if policy allows.

States:
- Pending review.
- Approved.
- Rejected with reason.
- Expired.

Acceptance:
- Ledger entries must show reference/order/topup.
- Balance must not be editable by client.

### 3.4 Client checkout flow

Steps:
```text
1. Select plan
2. Confirm specs and price
3. Check wallet balance
4. Accept terms/AUP
5. Submit order
6. Show provisioning status
```

UI validations:
- Insufficient wallet -> show top-up CTA.
- Plan out of stock -> disable buy.
- Provider unavailable -> disable/notify.
- Risk/manual review -> show pending review.

Result states:
```text
order provisioning queued
service active
manual review
failed with refund/reversal
```

### 3.5 Client orders list/detail

List columns:
```text
order_number
product/plan
amount
payment_status
provisioning_status
created_at
```

Detail:
```text
order timeline
ledger references
reservation/provisioning status
service link if active
refund/reversal if any
```

Actions:
- Cancel if allowed.
- Contact support.

### 3.6 Client services list

Columns/cards:
```text
service name/type
status
identifier masked/summary
expiry date
renew button
actions available
```

Filters:
```text
active
expired
suspended
terminated
vps/proxy
```

### 3.7 Client service detail

Data:
```text
service_id
product/plan snapshot
status
term_start_at
term_end_at
billing_status
service identifier
credential masked_hint
capability actions
lifecycle timeline
support link
```

Actions:
```text
Reveal credential
Renew
Request cancel
Reset password if supported
Reinstall if supported
Change IP if supported
Open ticket
```

Rules:
- Reveal credential is separate action and audit.
- Reinstall/change IP require confirmation.
- Terminated service has no renew/action except view history.

### 3.8 Client support/ticket screen

Data:
```text
tickets
messages
related service/order
status
```

Actions:
- Create ticket.
- Reply.
- Attach proof/evidence if allowed.

Security:
- Attachment private.
- Do not paste credential automatically.

---

## 4. Reseller Portal wireflow

### 4.1 Reseller dashboard

Data:
```text
reseller wallet balance
low balance warning
client count
active services
pending top-ups
profit summary
manual review items
expiring services
abuse warnings
```

Actions:
- Top up reseller wallet.
- Approve client top-up.
- Manage catalog.
- View failed/manual review orders.
- Manage clients.

### 4.2 Branding and storefront settings

Data:
```text
logo
brand color
support email
telegram/support link
footer
terms/AUP overrides if allowed
```

Actions:
- Update branding.
- Preview storefront.

Audit:
```text
tenant.branding.updated
```

### 4.3 Domain settings

Data:
```text
system subdomain
custom domains
verification status
TLS status
primary domain
```

Actions:
- Add domain.
- Verify.
- Set primary.
- Disable/remove domain.

Error states:
- Domain already used.
- Verification failed.
- TLS pending/failed.

Audit:
```text
tenant.domain.created
tenant.domain.verified
tenant.domain.primary_changed
```

### 4.4 Reseller catalog management

Data:
```text
available master plans
tenant cloned plans
selling price
reseller cost
margin
visibility
status
stock/source summary
```

Actions:
- Clone plan.
- Update selling price.
- Hide/show plan.
- Sync with master version.
- Disable plan.

Warnings:
```text
selling_price < reseller_cost
master plan has new version
source out of stock
provider degraded
```

Acceptance:
- Reseller never sees provider secret.
- Cost shown must be reseller cost, not necessarily platform internal base cost.

### 4.5 Client management

Data:
```text
client list
status
wallet balance
orders count
services count
risk flags
last login
```

Actions:
- Create/invite client if allowed.
- Suspend/disable client.
- View client detail.
- View wallet/order/service within tenant.
- Add risk note.

Security:
- Only clients under reseller tenant.

### 4.6 Client top-up review

Data:
```text
request amount
payment method
payment reference
proof attachment
client info
history
```

Actions:
- Approve.
- Reject with reason.
- Mark under review.
- Request more info.

Controls:
- Finance permission required.
- Approved creates ledger credit.
- Cannot approve twice.

Audit:
```text
wallet.topup.approved
wallet.topup.rejected
```

### 4.7 Reseller wallet/settlement

Data:
```text
reseller settlement wallet balance
ledger entries
platform cost debits
top-up history
low balance threshold
```

Important message:
```text
If reseller balance is insufficient, new client orders may not be provisioned even if client wallet has funds.
```

### 4.8 Reseller report

Reports:
```text
gross revenue from clients
platform reseller cost
gross profit
refunds/adjustments
top products
top clients
expiring services
```

Must use:
```text
ledger/order snapshots
```

Not current plan price.

---

## 5. Admin Portal wireflow

### 5.1 Admin dashboard

Cards:
```text
total tenants
active services
today orders
provisioning failures
manual review queue
provider health
pending reseller top-ups
abuse cases
revenue summary
```

Critical alerts:
- Provider down.
- Queue backlog.
- Failed provisioning spike.
- Manual review aging.
- Low stock.
- Credential reveal spike.

### 5.2 Tenant management

Data:
```text
tenant name/type/status
owner
domain
wallet balance
clients/services
risk status
created_at
```

Actions:
- Create tenant.
- Suspend/disable tenant.
- View tenant detail.
- Emergency access.
- Manage owner.

Emergency access:
- Requires reason.
- 2FA.
- Audit.

### 5.3 Master catalog

Screens:
```text
products list/detail
plans list/detail
plan-source mapping
version history
```

Actions:
- Create/update product.
- Create/update plan.
- Disable/archive.
- Create new version.
- Map source.
- Set reseller cost/min price.

Warnings:
- Plan used by active services.
- Changing cost may affect reseller margin.
- Source disabled/out-of-stock.

### 5.4 Provider/source management

Data:
```text
provider account
source type
health status
capability profile
inventory mode
capacity
last sync
failure rate
```

Actions:
- Add provider account.
- Update encrypted credentials.
- Add source.
- Test health.
- Disable source.
- Sync inventory.
- View provider request logs redacted.

Security:
- Provider secrets masked.
- Credential update requires critical permission/2FA.

### 5.5 Top-up and finance operations

Screens:
```text
reseller top-up queue
client direct top-up queue
ledger search
adjustment request/create
finance reports
```

Actions:
- Approve/reject top-up.
- Create adjustment.
- Export report.
- Search by correlation_id.

Controls:
- Adjustment requires reason/reference.
- Cannot edit ledger.
- Audit critical.

### 5.6 Provisioning operations

Data:
```text
jobs queued/running/failed/manual_review
source
tenant
order/service
attempt_count
retry_safety
error_code
correlation_id
```

Actions:
- Retry safe.
- Resolve manual review.
- Link existing external resource.
- Mark failed and refund.
- Disable problematic source.

Warnings:
- Unsafe retry requires high permission and reason.
- Provider unknown status should not retry blindly.

### 5.7 Service operations

Data:
```text
service status
tenant/client
provider mapping
billing status
lifecycle events
credential masked
```

Actions:
- Suspend.
- Unsuspend.
- Terminate.
- Reset password if supported.
- Sync provider.
- Reveal credential if permission.

Controls:
- Terminate critical.
- Reveal audit.
- Reason required for suspend/terminate.

### 5.8 Audit logs

Filters:
```text
tenant
actor
action
target
correlation_id
date range
risk level
```

Display:
```text
redacted before/after
metadata
ip/user agent
job id if worker
```

No plaintext secrets.

### 5.9 Abuse/risk center

Data:
```text
open abuse cases
risk flags
high-risk orders
provider notices
blacklist
```

Actions:
- Create case.
- Attach evidence.
- Warn client.
- Suspend service.
- Close case.
- Blacklist marker.

---

## 6. Common UI components

### 6.1 Status badges

Standard badges:
```text
active
pending
provisioning
manual review
failed
suspended
expired
terminated
out of stock
provider degraded
```

### 6.2 Timeline component

Used for:
```text
order timeline
service lifecycle
top-up review
provisioning job
abuse case
```

Timeline event:
```text
time
action
actor/system
status
note
correlation_id optional
```

### 6.3 Confirmation modal

Required for:
```text
reveal credential
reinstall
change IP
suspend
terminate
ledger adjustment
provider disable
domain primary change
```

### 6.4 Empty states

Examples:
```text
No services yet. Buy your first VPS/proxy.
No pending top-ups.
No manual review items.
No active provider sources.
```

### 6.5 Error states

Error should show:
```text
human message
error code
support correlation_id if needed
next action
```

Do not expose provider raw secret/error.

---

## 7. UI acceptance criteria

UI spec đạt khi:
- 3 portal có screen map rõ.
- Mỗi màn có role/data/action/error state.
- Credential luôn masked, reveal qua action riêng.
- Checkout hiển thị rõ insufficient client/reseller balance.
- Reseller thấy margin/cost của tenant mình.
- Admin thấy manual review/provider health/finance alerts.
- Action destructive có confirmation + reason.
- UI dùng capability để ẩn action không hỗ trợ.
- Cross-tenant data không thể xuất hiện trong UI.
- Mỗi lỗi quan trọng có correlation_id để support tra.

Câu nền: **UI tốt không chỉ đẹp; nó làm trạng thái hệ thống trở nên thật, để người vận hành biết phải làm gì trước khi lỗi thành tiền mất.**
