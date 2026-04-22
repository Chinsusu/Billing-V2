# 21 - QA Test Cases And Acceptance Plan

## 1. Mục tiêu tài liệu

Tài liệu này định nghĩa test plan và acceptance criteria cho dự án VPS/Proxy.

Mục tiêu:
- Test đúng nghiệp vụ tiền, tenant, provisioning, credential.
- Phát hiện lỗi trước khi lên production.
- Giảm kiểu test “click thấy chạy là được”.
- Làm căn cứ nghiệm thu dev.

---

## 2. Test principles

### 2.1 Test theo rủi ro

Ưu tiên test:
```text
money
tenant isolation
credential
provisioning idempotency
inventory oversell
refund/reversal
permission
abuse/suspend/terminate
```

### 2.2 Test negative case bắt buộc

Không chỉ test happy path. Dự án này chết ở edge case:
```text
provider timeout
reseller wallet thiếu
cross-tenant access
double click checkout
reservation expired
top-up approve twice
unsafe retry
credential reveal unauthorized
```

### 2.3 Test bằng correlation_id

Mỗi test nghiệp vụ phức tạp phải kiểm tra trace:
```text
order -> ledger -> reservation -> provisioning_job -> provider_request -> service -> audit
```

---

## 3. Test environment requirements

Minimum:
```text
staging environment
test tenant platform
test reseller A
test reseller B
test clients under each reseller
mock/manual provider source
at least one finite inventory source
test wallet balances
email/notification sandbox
```

Provider test:
```text
successful provider
out-of-stock provider
timeout/unknown provider
auth-failed provider
manual provider
```

---

## 4. Tenant isolation test cases

### TI-001 Client cannot access another client's service

Setup:
```text
Client A and Client B in same tenant.
Client B owns service S.
```

Action:
```text
Client A requests service S detail.
```

Expected:
```text
403 or 404.
No credential returned.
Audit/security event optional.
```

### TI-002 Reseller A cannot access Reseller B client

Setup:
```text
Reseller A tenant.
Reseller B tenant.
Client B belongs to Reseller B.
```

Action:
```text
Reseller A staff requests Client B profile/order/service.
```

Expected:
```text
403/404.
No data leak.
```

### TI-003 Body tenant_id injection ignored

Action:
```text
Client under tenant A submits request with tenant_id = tenant B in body.
```

Expected:
```text
Backend uses tenant A from context.
Request fails if resource not in tenant A.
No write under tenant B.
```

### TI-004 Domain maps correct tenant

Action:
```text
Open storefront from reseller custom domain.
```

Expected:
```text
Catalog/user registration belongs to mapped tenant.
```

---

## 5. RBAC test cases

### RBAC-001 Support cannot approve top-up

Setup:
```text
User role = support_agent without wallet.topup.approve.
```

Action:
```text
Approve top-up request.
```

Expected:
```text
FORBIDDEN.
No ledger entry.
Audit permission denied optional.
```

### RBAC-002 Finance cannot reveal credential

Action:
```text
Finance agent calls credential reveal.
```

Expected:
```text
FORBIDDEN.
No plaintext returned.
No reveal audit except denied event optional.
```

### RBAC-003 Reseller staff only with permission can update plan price

Expected:
```text
Staff without catalog.tenant.price_update cannot update.
Staff with permission can update own tenant plan only.
```

### RBAC-004 Admin emergency access requires reason

Action:
```text
Admin attempts tenant emergency access without reason.
```

Expected:
```text
VALIDATION_ERROR.
No emergency session.
```

---

## 6. Wallet and ledger test cases

### WL-001 Top-up approval creates exactly one ledger entry

Setup:
```text
Client submits top-up 100.
```

Action:
```text
Reseller approves top-up.
```

Expected:
```text
topup status approved.
one credit ledger entry amount 100.
wallet balance +100.
audit wallet.topup.approved.
```

### WL-002 Top-up cannot be approved twice

Action:
```text
Approve same top-up twice.
```

Expected:
```text
second attempt returns TOPUP_ALREADY_REVIEWED.
no second ledger entry.
balance unchanged.
```

### WL-003 Ledger posted cannot be edited

Action:
```text
Attempt update/delete posted ledger entry through API/admin.
```

Expected:
```text
forbidden/not supported.
adjustment required.
```

### WL-004 Manual adjustment requires reason

Action:
```text
Finance creates adjustment without reason.
```

Expected:
```text
VALIDATION_ERROR.
No ledger entry.
```

### WL-005 Client wallet sufficient but reseller wallet insufficient

Setup:
```text
Client wallet = 100.
Plan selling_price = 20.
Reseller settlement wallet = 5.
Reseller cost = 12.
```

Action:
```text
Client checkout.
```

Expected:
```text
Error INSUFFICIENT_RESELLER_BALANCE
No provisioning job.
No service.
No stock allocated.
Client wallet not debited, unless policy explicitly creates pending order without debit.
```

---

## 7. Checkout and reservation test cases

### CO-001 Successful checkout creates order/reservation/ledger/job

Expected:
```text
order created.
order_item snapshots stored.
reservation status reserved then allocated after provisioning.
client wallet debited.
reseller wallet debited if reseller tenant.
provisioning_job queued.
correlation_id present.
```

### CO-002 Double click checkout with same idempotency key

Action:
```text
Submit same checkout twice same idempotency key.
```

Expected:
```text
Only one order.
Only one wallet debit.
Only one reservation.
Only one provisioning job.
Second response returns same result.
```

### CO-003 Same idempotency key with different payload

Expected:
```text
IDEMPOTENCY_CONFLICT.
No new order.
```

### CO-004 Out of stock prevents checkout

Setup:
```text
source available = 0.
```

Expected:
```text
OUT_OF_STOCK.
No wallet debit.
No reservation.
No provisioning job.
```

### CO-005 Reservation expires

Setup:
```text
reservation created but payment/provisioning not completed.
expires_at passed.
```

Action:
```text
reservation_expiry_job runs.
```

Expected:
```text
reservation expired.
reserved_count decremented once.
order expired/cancelled.
audit reservation.expired.
```

### CO-006 Concurrent checkout for last stock

Setup:
```text
source capacity available = 1.
Two clients checkout same plan concurrently.
```

Expected:
```text
Only one succeeds.
Other gets OUT_OF_STOCK.
No oversell.
reserved/allocated counts correct.
```

---

## 8. Provisioning test cases

### PR-001 Provider success activates service

Expected:
```text
provider_request success.
external_resource_id stored.
service active.
credential encrypted.
reservation allocated.
order provisioning_status success.
activation notification queued.
```

### PR-002 Provider out of stock after reservation

Expected:
```text
job failed do_not_retry.
reservation released.
wallet reversal/refund per policy.
order failed.
notification sent.
```

### PR-003 Provider timeout unknown

Action:
```text
Adapter returns timeout unknown after create request sent.
```

Expected:
```text
job manual_review.
No automatic retry.
Order provisioning_status manual_review.
Operator alert.
No duplicate provider create.
```

### PR-004 Provider auth failure

Expected:
```text
source marked degraded/down or disabled per policy.
job failed/manual_review.
alert admin.
no retry.
```

### PR-005 Success but credential missing

Expected:
```text
job manual_review or credential fetch attempt if safe.
service not shown as fully active until credential ready, unless policy allows active_pending_credential.
```

### PR-006 Manual provider activation

Setup:
```text
manual source.
```

Expected:
```text
checkout creates paid order + reservation + manual_review job.
operator enters resource/credential.
service active.
credential encrypted.
audit manual resolution.
```

---

## 9. Service lifecycle test cases

### SL-001 Renew active service

Setup:
```text
service active, term_end_at future.
wallet sufficient.
```

Expected:
```text
wallet debit.
term_end_at extended from old term_end_at.
lifecycle event service.renewed.
```

### SL-002 Renew expired service in grace

Expected according to policy:
```text
term_end_at extends from old term_end_at if policy says no free days.
service unsuspended/active if required.
```

### SL-003 Cannot renew terminated service

Expected:
```text
VALIDATION_ERROR or SERVICE_NOT_RENEWABLE.
No wallet debit.
```

### SL-004 Auto expiry

Action:
```text
service_expiry_job runs after term_end_at.
```

Expected:
```text
billing_status overdue/grace.
notification queued.
no immediate terminate unless policy.
```

### SL-005 Auto suspend after grace

Expected:
```text
suspend job created.
service suspended.
suspension_reason = billing_overdue.
audit service.suspended.
```

### SL-006 Terminate after hold period

Expected:
```text
terminate job created only if policy auto_terminate true.
reason required.
provider terminate result handled.
```

---

## 10. Credential security test cases

### CR-001 Service detail returns masked credential only

Expected:
```text
no plaintext password/token in response.
masked_hint present.
```

### CR-002 Reveal credential creates audit

Action:
```text
Client reveals own credential.
```

Expected:
```text
plaintext returned only in reveal response.
audit credential.revealed.
last_revealed_at updated.
```

### CR-003 Unauthorized reveal denied

Action:
```text
Client tries reveal credential of another service.
```

Expected:
```text
403/404.
No plaintext.
No last_revealed update.
```

### CR-004 Logs/audit do not contain secret

Action:
```text
Search provider_request/audit/notification payload after activation.
```

Expected:
```text
no root password/proxy password/API key plaintext.
```

---

## 11. Catalog/pricing test cases

### CP-001 Reseller cannot sell disabled master plan

Expected:
```text
PLAN_DISABLED.
```

### CP-002 Price snapshot unaffected by future price update

Setup:
```text
Client buys plan at 20.
Admin later changes suggested price to 25.
```

Expected:
```text
existing order price_snapshot remains 20.
reports for order use 20.
```

### CP-003 Margin risk warning

Setup:
```text
reseller selling_price < reseller_cost.
```

Expected:
```text
plan status margin_risk or update blocked depending policy.
checkout blocked if policy requires non-negative margin.
```

### CP-004 Capability masking

Setup:
```text
source does not support change_ip.
```

Expected:
```text
UI hides change_ip.
API returns CAPABILITY_NOT_SUPPORTED if called.
```

---

## 12. Abuse/risk test cases

### AB-001 New high-risk order routes manual review

Setup:
```text
risk rule triggered.
```

Expected:
```text
order manual_review before provisioning.
no provider job until approved.
```

### AB-002 Abuse suspend requires reason/evidence

Action:
```text
Admin suspends service for abuse without reason.
```

Expected:
```text
VALIDATION_ERROR.
```

### AB-003 Abuse case workflow

Expected:
```text
case open -> investigating -> warning/suspended/resolved.
audit all transitions.
client notification if policy.
```

---

## 13. Notification test cases

### NT-001 Activation notification

Expected:
```text
sent/queued after service active.
does not contain plaintext credential unless policy explicitly allows.
contains link to service detail.
```

### NT-002 Expiry reminder dedupe

Action:
```text
renewal_reminder_job runs twice same day/window.
```

Expected:
```text
only one notification for same service/window.
```

### NT-003 Reseller low balance warning

Expected:
```text
notification sent when settlement wallet below threshold.
not spammed repeatedly beyond dedupe policy.
```

---

## 14. Report test cases

### RP-001 Reseller profit calculation

Setup:
```text
client purchase selling_price 20.
reseller_cost 12.
```

Expected:
```text
gross revenue 20.
platform cost 12.
gross profit 8 before refund/adjustment.
```

### RP-002 Refund affects report

Expected:
```text
refund/reversal reflected.
uses ledger entries and snapshots.
```

### RP-003 Admin provider health report

Expected:
```text
failed/manual_review/success count matches provisioning_jobs/provider_requests.
```

---

## 15. Deployment smoke tests

Before production:
- Login/register works per tenant domain.
- Client top-up submit works.
- Top-up approval posts ledger.
- Checkout success with manual/mock provider.
- Service active detail shows masked credential.
- Credential reveal works and audits.
- Cross-tenant access blocked.
- Reservation expiry job works.
- Service expiry reminder works.
- Provider health job reports.
- Backup job completed and restore tested in staging.

---

## 16. Acceptance gates

### Gate 1: Foundation

Pass if:
```text
tenant isolation tests pass
RBAC high-risk tests pass
wallet ledger tests pass
```

### Gate 2: Commerce

Pass if:
```text
checkout/reservation/idempotency tests pass
reseller settlement tests pass
refund/reversal tests pass
```

### Gate 3: Provisioning

Pass if:
```text
provider success/fail/timeout/manual tests pass
credential security tests pass
service lifecycle tests pass
```

### Gate 4: Production readiness

Pass if:
```text
notification/report tests pass
monitoring/backup/restore smoke tests pass
abuse/manual review tests pass
```

---

## 17. QA acceptance criteria

Project không nên production nếu fail bất kỳ P0:
- Cross-tenant access leaks data.
- Ledger double posts.
- Checkout double debits.
- Reservation oversells.
- Provider timeout creates duplicate resources.
- Credential plaintext appears in logs/audit.
- Reseller client order provisions while reseller wallet lacks cost.
- Staff without permission can approve money/reveal credential.
- Terminated service can be renewed accidentally.
- Backup/restore not tested.

Câu nền: **QA của hệ thống hạ tầng không phải tìm lỗi giao diện; QA là đóng những cánh cửa nơi tiền, tenant và credential có thể rò ra ngoài.**
