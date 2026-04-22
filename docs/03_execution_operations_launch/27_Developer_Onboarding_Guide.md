# 27 — Developer Onboarding Guide

**Version:** v1.3  
**Date:** 2026-04-22  
**Audience:** Backend, frontend, QA, DevOps, technical PM.

## 1. What this project is

This is a multi-tenant, white-label VPS/proxy selling platform.

```text
Platform Admin:
- manages tenants/resellers/provider sources/catalog/wallet approvals/audit.

Reseller:
- owns a storefront and client base;
- sets selling prices;
- controls client wallet/top-up flow;
- pays platform through reseller wallet.

Client:
- belongs to a tenant/reseller;
- tops up wallet;
- buys/renews services;
- views credentials.
```

## 2. Reading order for new developers

```text
00_README.md
01_Product_Scope_Business_Model.md
02_Tenant_Model_Role_Architecture.md
04_Billing_Wallet_Ledger_Spec.md
09_Reseller_Settlement_Ledger_Model.md
10_Tenant_Security_Access_Control_Spec.md
11_Provisioning_Idempotency_And_Inventory_Locking.md
15_Database_Schema_And_ERD.md
16_API_Contract_And_Permission_Spec.md
17_RBAC_Permission_Matrix.md
21_QA_Test_Cases_And_Acceptance_Plan.md
26_MVP_Scope_Lock_And_Non_Goals.md
```

If time is limited, read files 09, 10, 11, 15, 16, and 21 first.

## 3. Mental model

### Tenant is the data border

```text
Tenant A must never read Tenant B resource.
Do not trust tenant_id from request body.
Resolve tenant from domain/session/token.
Query tenant resources with tenant scope.
```

### Ledger is the financial truth

```text
Never update old ledger entries.
Never delete ledger entries.
Use adjustments/reversals.
wallet_balance = sum(ledger entries)
```

### Provider is not fully reliable

Provider can timeout, return partial success, run out of stock, create a resource without returning response, or change behavior. This is why provisioning must use queue, idempotency, retry classification, and manual review.

## 4. Money flow to memorize

Client under reseller buys service:

```text
Client wallet: -selling_price
Reseller wallet: -reseller_cost
Platform revenue: +reseller_cost
Reseller gross profit: selling_price - reseller_cost
```

Do not provision if reseller wallet cannot cover reseller_cost.

## 5. Ten rules you must not break

```text
1. Do not update/delete ledger entries.
2. Do not query tenant resources without tenant scope.
3. Do not log plaintext credential/password/provider token.
4. Do not blindly retry provider create after timeout.
5. Do not use current price for old order disputes; use snapshot.
6. Do not allow critical admin actions without audit.
7. Do not provision if reseller wallet is insufficient.
8. Do not activate service if reservation is not allocated.
9. Do not hardcode provider secret.
10. Do not expose raw provider internal errors to clients.
```

## 6. Backend responsibilities

```text
- tenant enforcement
- RBAC enforcement
- wallet transaction atomicity
- ledger append-only behavior
- provisioning idempotency
- credential encryption
- audit writing
- stable error codes
```

Frontend permission hiding is not security. API must enforce permission.

## 7. Frontend responsibilities

```text
- show data by role and capability
- hide unauthorized actions
- never keep credentials visible longer than necessary
- show clear pending/error states
- map stable error codes to user-friendly copy
```

## 8. Worker responsibilities

```text
- provisioning jobs
- provider sync
- expiry/suspension/termination
- safe retries
- manual review for unsafe states
```

Workers must be idempotent. Running the same job twice must not create duplicate resource or debit money twice.

## 9. Transaction pattern

Recommended checkout pattern:

```text
DB transaction:
- validate tenant/user/plan/source
- reserve stock atomically
- create order/order item snapshot
- create ledger debit entries
- create provisioning job
Commit

Worker:
- calls provider
- stores provider response/resource mapping
- activates service or manual review/fail path
```

Avoid long DB transactions around external provider API calls.

## 10. Error code style

Use stable codes:

```text
TENANT_NOT_FOUND
TENANT_INACTIVE
PERMISSION_DENIED
PLAN_DISABLED
SOURCE_DISABLED
OUT_OF_STOCK
INSUFFICIENT_CLIENT_BALANCE
INSUFFICIENT_RESELLER_BALANCE
RESERVATION_EXPIRED
PROVIDER_UNAVAILABLE
PROVISIONING_UNCERTAIN
CREDENTIAL_ACCESS_DENIED
RATE_LIMITED
```

## 11. Audit rule

Write audit for:

```text
login failures, tenant changes, role changes, top-up approval/rejection, wallet adjustment, order/refund, service activation/suspension/termination, credential reveal, provider config changes, price changes, emergency access.
```

Audit must redact:

```text
password, token, provider_api_key, credential, secret
```

## 12. Local development expectations

Local stack should include:

```text
API
frontend
worker
database
queue
mail/notification mock
provider mock
```

Provider mock should simulate:

```text
success
out_of_stock
rate_limit
auth_failed
timeout_before_response
timeout_after_resource_created
success_missing_credential
```

## 13. PR checklist

```text
- Tenant scope enforced?
- Permission checked?
- Audit event added?
- P0/P1 tests added?
- No secret/credential logs?
- Idempotent if job/worker?
- Stable error code?
- Migration safe?
```

## 14. Common traps

| Trap | Damage |
|---|---|
| Query by id only | Cross-tenant leak |
| Retry provider create after timeout | Duplicate VPS/proxy |
| Update ledger | Broken finance truth |
| UI-only permission | Direct API bypass |
| Store credential in audit diff | Data leak |
| Use current price for refund | Billing dispute |
| No capability snapshot | Wrong service actions |

## 15. Closing principle

```text
This project is not hard because it has many screens.
It is hard because money, tenant, provider, and lifecycle must be correct at the same time.
```
