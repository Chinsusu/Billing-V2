# 82 - Admin UAT Evidence

**Date:** 2026-05-21
**Scope:** Admin portal UAT evidence on the selected bounded non-production test environment.
**Decision:** PASS for automated Admin UAT continuation on the selected test environment only. This does not approve production, production customer data, broader private beta, or broader real-provider scope.

## Boundary

This evidence uses the selected test runtime behind `billing.resvn.net` and the protected local secret files on the selected host.

Do not broaden this result:

- No production database or production customer data was used.
- No raw DB DSN, password, cookie, session token, provider token, provider payload, Telegram token, TOTP value, or plaintext service credential is recorded here.
- No Cloudmini create/delete/action route was called by this task.
- This task used admin auth gate checks, admin read checks, one admin top-up approve/reject fixture, finance reconciliation, and cleanup/status checks.

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
/etc/cloudflared/tunnel.token: 600 root:root
process argv secret-pattern check after UAT: no matches
cloudflared token-file arg: present
```

## Automated Admin Evidence

Admin login and 2FA gate:

```text
domain: billing.resvn.net
result: PASS
admin_login_status=200
admin_actor_type=platform_staff
admin_two_factor_required=yes
admin_two_factor_satisfied=no
admin_two_factor_setup_required=no
browser_two_factor_panel_visible=yes
logout_from_2fa_panel=pass
raw_cookies_printed=no
```

Operational note:

```text
auth_login_rate_limit_counters_cleared=35
reason=the automated UAT probes themselves hit the selected non-production auth login rate limit before the browser gate check
scope=auth.login counters in the selected non-production DB only
```

Admin API read checks:

```text
command: admin session-equivalent API read probe against https://billing.resvn.net/backend
result: PASS
admin_catalog_provider_readiness_count=4
admin_provider_sources_count=2
admin_orders_count=2
admin_services_count=2
admin_invoices_count=3
admin_transactions_count=3
admin_topup_requests_count=2
admin_jobs_count=2
admin_audit_logs_count=6
admin_request_id_present=yes
admin_sensitive_field_markers_visible=0
raw_ids_printed=no
raw_cookies_printed=no
```

Negative access checks:

```text
client_to_admin_denial=403 auth.permission_denied
low_permission_admin_denial=403 auth.permission_denied
```

Admin top-up review:

```text
command: admin top-up approve/reject probe against https://billing.resvn.net/backend
result: PASS
wallet_display_id=10000
approve_topup_display_id=10001
approve_ledger_display_id=10003
approve_audit_display_id=10007
reject_topup_display_id=10002
reject_ledger_count=0
reject_audit_display_id=10008
wallet_balance_delta_minor=333
provider_side_effects=none
raw_ids_printed=no
```

Finance reconciliation support evidence after the admin top-up mutation:

```text
command: dev-target-finance-reconciliation against https://billing.resvn.net/backend
result: PASS
transaction_display_id=51001
invoice_display_id=44001
wallet_display_id=41001
ledger_display_id=50002
daily_date=2026-04-23
daily_status=balanced
wallets_checked=2
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

- The task intentionally left dev/test admin UAT top-up records in the selected target DB for traceability.
- No claimable/running jobs remained after the run.
- No non-terminated Cloudmini-backed service remained after the run.
- Wallet projection matched posted ledger source-of-truth after the run.
- Existing provider job and service counts were not modified by this task.

## Admin UAT Matrix

| Area | Result | Evidence |
| --- | --- | --- |
| Login/2FA gate | PASS | Browser login reached the 2FA-required panel and logged out without printing cookies. |
| Admin read access | PASS | Admin API read checks covered catalog readiness, provider sources, orders, services, invoices, transactions, top-ups, jobs, and audit logs. |
| Provider/source visibility | PASS | Admin read checks returned readiness/source rows without secret markers in evidence. |
| Service/provisioning visibility | PASS | Admin service and job reads returned scoped records; no provider mutation route was called. |
| Top-up approve/reject | PASS | Admin route approved one dev/test top-up with one posted ledger credit and audit, and rejected one top-up with no ledger credit and audit. |
| Invoice/payment/finance | PASS | Finance reconciliation stayed balanced after the admin top-up mutation. |
| Audit | PASS | Top-up approve/reject audit rows were present with public audit display IDs only. |
| RBAC negative | PASS | Client and low-permission actors were denied on admin routes. |
| Credential safety | PASS by no-reveal admin checks | This task did not reveal credentials; admin read probes checked for sensitive field markers and did not record plaintext or encrypted payloads. |
| Notification/fallback | PASS by prior selected-host evidence | T279/T280/T281 prove selected-host Telegram preflight, one queued delivery, and retry/terminal classification; this task did not read notification payloads. |

## Not Verified

- Human admin tester sign-off was not collected in this task.
- Full 2FA-satisfied admin browser navigation was not executed in this task; the task verified the live browser 2FA gate and used admin API probes for read coverage.
- `dev-target-auth-rbac` was not used as pass evidence for T287 because the current selected target client seed credential no longer matches that smoke command's hardcoded password expectation and returns `auth.invalid_credentials`; admin-specific login and RBAC probes passed.
- Admin credential reveal mutation was not executed to avoid handling plaintext credentials in this UAT evidence.
- Admin service lifecycle mutations, job retry/cancel/manual-review mutations, support ticket mutations, and real Cloudmini provisioning were not executed.
- Production notification delivery was not executed; selected-host Telegram/fallback evidence remains inherited from T279/T280/T281 and related launch evidence.

## Decision

```text
Admin UAT result: PASS for automated selected-environment evidence
Open P0 bugs: 0
Open P1 bugs: 0
Residual risk: human UX sign-off, full 2FA-satisfied browser navigation, stale generic target auth smoke credentials, admin credential reveal/lifecycle/support mutations, and real Cloudmini provisioning remain outside this task
Next recommended task: consolidate client/reseller/admin UAT packet and decide whether to refresh generic target auth smoke credentials or parameterize that smoke
```
