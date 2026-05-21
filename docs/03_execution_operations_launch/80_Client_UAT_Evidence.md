# 80 - Client UAT Evidence

**Date:** 2026-05-21
**Scope:** Client portal UAT evidence on the selected bounded non-production test environment.
**Decision:** PASS for automated Client UAT continuation on the selected test environment only. This does not approve production, production customer data, broader private beta, or broader real-provider scope.

## Boundary

This evidence uses the selected test runtime behind `client.resvn.net` and the protected local secret files on the selected host.

Do not broaden this result:

- No production database or production customer data was used.
- No raw DB DSN, password, cookie, session token, provider token, provider payload, Telegram token, TOTP value, or plaintext service credential is recorded here.
- No Cloudmini create/delete/action route was called by this task.
- The checkout/provisioning path used the existing dev/test smoke path and fake-provider fulfillment.

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
cloudflared argv exact token arg: absent
cloudflared token-file arg: present
process argv secret-pattern check after UAT: no matches
```

## Target Domain Mapping

Initial browser login on `client.resvn.net` failed with `tenant.context_missing` because the selected target DB had no rows for the selected public domains in `tenant_domains`.

The selected non-production target DB was updated with these verified mappings before rerunning browser UAT:

```text
domain_mapping_upserted count=3
billing.resvn.net verification=verified tls=active primary=true domain_display_id=10000 tenant_display_id=10000
client.resvn.net verification=verified tls=active primary=true domain_display_id=10001 tenant_display_id=10001
reseller.resvn.net verification=verified tls=active primary=false domain_display_id=10002 tenant_display_id=10001
```

This was a target-environment data/config fix only. Repeat it for any new database, restored database, host, or launch scope.

## Automated Client Evidence

Auth/RBAC smoke:

```text
command: dev-target-auth-rbac against https://client.resvn.net/backend
result: PASS
client_session_cookie_only=pass
admin_2fa_gate=pass
invalid_session_denied=pass
actor_required_denied=pass
tenant_mismatch_denied=pass
rbac_denials=3
provider_mutation_routes_called=no
money_mutation_routes_called=no
```

Client billing path smoke:

```text
command: dev-billing against https://client.resvn.net/backend
result: PASS
topup=10000
order=10000
invoice=10000
transaction=10000
ledger=10001
service=10000
renewal_invoice=10002
renewal_transaction=10001
renewal_ledger=10002
provider path: fake-provider fulfillment
```

Credential reveal smoke:

```text
command: dev-target-credential-reveal against https://client.resvn.net/backend
result: PASS
service_display_id=43001
credential_type=recovery_code
client_session_cookie_only=pass
no_store=pass
audit_display_id=10006
last_revealed_by=client_actor
rate_limit_attempts=1
provider_mutation_routes_called=no
money_mutation_routes_called=no
```

Finance reconciliation support evidence:

```text
command: dev-target-finance-reconciliation against https://client.resvn.net/backend
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

Client browser check:

```text
domain: client.resvn.net
result: PASS
login=pass
portal_scope=client_only
logout=pass
auth_login_status=200
auth_logout_status=200
raw_cookies_printed=no
```

## Cleanup And Residual State

Post-run database and queue checks:

```text
recent_client_uat_smoke_records topups=1 invoices=2 services=1
provider_jobs succeeded count=1
provider_jobs manual_review count=1
claimable_or_running_jobs=0
cloudmini_nonterminated_services=0
wallet_mismatch_count=0
service_status active count=2
```

Interpretation:

- The task intentionally left dev/test UAT records in the selected target DB for traceability.
- No claimable/running jobs remained after the run.
- No non-terminated Cloudmini-backed service remained after the run.
- Wallet projection matched posted ledger source-of-truth after the run.
- The existing manual-review provider job count was not modified by this task.

## Client UAT Matrix

| Area | Result | Evidence |
| --- | --- | --- |
| Login/logout | PASS | Browser check on `client.resvn.net` returned login/logout `200`. |
| Session scope | PASS | Browser exposed only the client portal after client login. |
| Dashboard/read access | PASS | Browser reached client overview after login; backend read paths passed through smokes. |
| Catalog | PASS | `dev-target-auth-rbac` verified cookie-only `/client/catalog`. |
| Top-up | PASS | `dev-billing` created and approved one dev/test top-up. |
| Checkout/payment | PASS | `dev-billing` created order, checkout invoice, wallet payment, and paid order. |
| Provisioning/service | PASS | `dev-billing` fulfilled one fake-provider service. |
| Renewal | PASS | `dev-billing` renewed the created service and recorded invoice/payment/ledger IDs. |
| Credential reveal | PASS | Credential reveal no-store/audit smoke passed without printing plaintext. |
| Invoice/finance | PASS | Finance reconciliation smoke remained balanced with zero mismatch counts. |
| RBAC negative | PASS | Invalid session, actor-required, tenant mismatch, admin 2FA gate, and 3 RBAC denials passed. |
| Notification/fallback | PASS by prior selected-host evidence | T279/T280/T281 prove selected-host Telegram preflight, one queued delivery, and retry/terminal classification; this task did not read notification payloads. |

## Not Verified

- Human client tester sign-off was not collected in this task.
- Deep browser navigation and form submission for every client screen was not used as pass evidence; one extra navigation probe was flaky and was discarded instead of weakening assertions.
- Reseller and admin UAT are not covered here.
- Real Cloudmini provisioning was not used in this task.

## Decision

```text
Client UAT result: PASS for automated selected-environment evidence
Open P0 bugs: 0
Open P1 bugs: 0
Residual risk: human UX sign-off and reseller/admin UAT remain pending
Next recommended task: run reseller UAT evidence, then admin UAT evidence
```
