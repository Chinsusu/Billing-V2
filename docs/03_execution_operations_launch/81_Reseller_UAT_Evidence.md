# 81 - Reseller UAT Evidence

**Date:** 2026-05-21
**Scope:** Reseller portal UAT evidence on the selected bounded non-production test environment.
**Decision:** PASS for automated Reseller UAT continuation on the selected test environment only. This does not approve production, production customer data, broader private beta, or broader real-provider scope.

## Boundary

This evidence uses the selected test runtime behind `reseller.resvn.net` and the protected local secret files on the selected host.

Do not broaden this result:

- No production database or production customer data was used.
- No raw DB DSN, password, cookie, session token, provider token, provider payload, Telegram token, TOTP value, or plaintext service credential is recorded here.
- No Cloudmini create/delete/action route was called by this task.
- This task used reseller session read checks and browser navigation only; it did not approve, reject, reveal, suspend, terminate, retry, or mutate provider resources.

## Entry Checks

Target health:

```text
billing.resvn.net / frontend: 200
client.resvn.net / frontend: 200
reseller.resvn.net / frontend: 200
billing.resvn.net /backend/healthz: 200
client.resvn.net /backend/healthz: 200
reseller.resvn.net /backend/healthz: 200
127.0.0.1:8080 /healthz: 200
127.0.0.1:3000 /: 200
```

Runtime services:

```text
billing-api.service: active/running
billing-frontend.service: active/running
cloudflared.service: active/running
postgresql@14-main.service: active/running
```

Secret handling:

```text
/etc/billing/secrets/billing-api.env: 600 root:root
/etc/billing/secrets/billing-frontend.env: 600 root:root
/etc/billing/secrets/cloudmini.env: 600 root:root
/etc/billing/secrets/telegram.env: 600 root:root
process argv secret-pattern check after UAT: no matches
```

Domain mapping:

```text
reseller.resvn.net was already mapped from T285 as verified/active to the selected reseller tenant.
No new domain mapping change was made in T286.
```

## Automated Reseller Evidence

Reseller browser login/scope check:

```text
domain: reseller.resvn.net
result: PASS
login=pass
portal_scope=reseller_only
logout=pass
auth_login_status=200
auth_logout_status=200
raw_cookies_printed=no
```

Reseller session API read checks:

```text
command: browser session API checks against https://reseller.resvn.net/backend
result: PASS
reseller/catalog count=2
reseller/catalog/master-plans count=4
reseller/customers count=1
reseller/orders count=2
reseller/services count=2
reseller/invoices count=3
reseller/transactions count=3
reseller/wallets count=1
reseller/topup-requests count=2
reseller/jobs count=2
request_id_present=yes
raw_ids_printed=no
raw_cookies_printed=no
```

Service detail safety:

```text
result: PASS
service_detail_safe=yes
forbidden_field_markers_visible=0
checked_markers=encrypted_payload,payload_json,raw_response,provider_account_id,authorization,cookie
```

Negative access checks:

```text
admin/invoices with reseller session: 403 auth.permission_denied
client/catalog with reseller session: 403 auth.permission_denied
```

Reseller UI navigation:

```text
result: PASS
screens=overview,accounts,services,invoices,transactions,products,reports,settings
reseller_read_failures=0
forbidden_text_visible=0
raw_cookies_printed=no
```

Finance reconciliation support evidence:

```text
command: dev-target-finance-reconciliation against https://reseller.resvn.net/backend
result: PASS
transaction_display_id=51001
invoice_display_id=44001
wallet_display_id=41001
ledger_display_id=50002
daily_date=2026-04-23
daily_status=balanced
wallets_checked=1
invoices_checked=1
payments_checked=1
wallet_mismatches=0
invoice_mismatches=0
duplicate_payment_references=0
money_mutation_routes_called=no
provider_mutation_routes_called=no
```

## Cleanup And Residual State

Post-run database and queue checks:

```text
claimable_or_running_jobs=0
cloudmini_nonterminated_services=0
wallet_mismatch_count=0
provider_jobs succeeded count=1
provider_jobs manual_review count=1
service_status active count=2
```

Interpretation:

- No claimable/running jobs remained after the run.
- No non-terminated Cloudmini-backed service remained after the run.
- Wallet projection matched posted ledger source-of-truth after the run.
- Existing provider job and service counts were not modified by this task.

## Reseller UAT Matrix

| Area | Result | Evidence |
| --- | --- | --- |
| Login/logout | PASS | Browser check on `reseller.resvn.net` returned login/logout `200`. |
| Session scope | PASS | Browser exposed only the reseller portal after reseller login. |
| Dashboard/read access | PASS | Browser reached reseller overview and read routes returned `200`. |
| Catalog | PASS | `reseller/catalog` and `reseller/catalog/master-plans` returned scoped records. |
| Customers | PASS | `reseller/customers` returned scoped records. |
| Services | PASS | `reseller/services` returned scoped records and service detail had no forbidden sensitive field markers. |
| Invoices/payments | PASS | `reseller/invoices` and `reseller/transactions` returned scoped records; finance reconciliation stayed balanced. |
| Wallet/top-up visibility | PASS | `reseller/wallets` and `reseller/topup-requests` returned scoped records. |
| Jobs/provisioning visibility | PASS | `reseller/jobs` returned scoped records; no provider mutation route was called. |
| RBAC negative | PASS | Reseller session was denied on admin and client routes. |
| UI navigation | PASS | Major reseller screens loaded with no reseller read failures or forbidden text markers. |
| Notification/fallback | PASS by prior selected-host evidence | T279/T280/T281 prove selected-host Telegram preflight, one queued delivery, and retry/terminal classification; this task did not read notification payloads. |

## Not Verified

- Human reseller tester sign-off was not collected in this task.
- Cross-reseller isolation against a second real reseller tenant was not tested because the selected target seed currently has one reseller tenant.
- Credential reveal mutation was not executed in this task to avoid exposing or handling plaintext credentials in UAT evidence.
- Reseller support ticket mutation and reseller lifecycle mutations were not executed.
- Admin UAT is not covered here.
- Real Cloudmini provisioning was not used in this task.

## Decision

```text
Reseller UAT result: PASS for automated selected-environment evidence
Open P0 bugs: 0
Open P1 bugs: 0
Residual risk: human UX sign-off, second-reseller isolation fixture, credential reveal mutation, support/lifecycle mutations, and admin UAT remain pending
Next recommended task: run admin UAT evidence
```
