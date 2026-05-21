# 83 - UAT Consolidated Evidence Packet

**Date:** 2026-05-21
**Scope:** Consolidated Client, Reseller, and Admin UAT evidence for the selected bounded non-production test environment.
**Decision:** PASS for automated selected-environment UAT continuation only. This does not approve production, production customer data, broader private beta, or broader real-provider scope.

## Boundary

This packet consolidates evidence from:

- `docs/03_execution_operations_launch/80_Client_UAT_Evidence.md`
- `docs/03_execution_operations_launch/81_Reseller_UAT_Evidence.md`
- `docs/03_execution_operations_launch/82_Admin_UAT_Evidence.md`

Do not broaden this result:

- No production database or production customer data was used.
- No raw DB DSN, password, cookie, session token, provider token, provider payload, Telegram token, TOTP value, or plaintext service credential is recorded here.
- No Cloudmini create/delete/action route was called by these UAT tasks.
- Telegram notification delivery evidence is inherited from the selected-host notification evidence tasks and was not rerun in this packet.

## Cross-Portal Entry Evidence

Selected target health passed for all public domains and local services:

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

Runtime and secret handling evidence passed:

```text
billing-api.service: active/running
billing-frontend.service: active/running
cloudflared.service: active/running
postgresql@14-main.service: active/running
protected secret files: 600 root:root
cloudflared token-file arg: present
process argv secret-pattern check: no matches
```

Domain mapping was verified for:

```text
billing.resvn.net
client.resvn.net
reseller.resvn.net
```

## Portal Results

| Portal | Result | Evidence doc | Covered |
| --- | --- | --- | --- |
| Client | PASS | `80_Client_UAT_Evidence.md` | Login/logout, cookie-only session, catalog, top-up, checkout, fake-provider service fulfillment, renewal, credential reveal no-store/audit, finance reconciliation, RBAC negative checks. |
| Reseller | PASS | `81_Reseller_UAT_Evidence.md` | Login/logout, reseller-only scope, catalog, customers, services, invoices, transactions, wallets, top-up visibility, jobs, service detail sensitive-field safety, finance reconciliation, RBAC negative checks. |
| Admin | PASS | `82_Admin_UAT_Evidence.md` | Login, live 2FA gate, admin read coverage, provider/source visibility, orders/services/invoices/transactions/top-ups/jobs/audit reads, top-up approve/reject fixture, finance reconciliation, RBAC negative checks. |

## Cross-Cutting Evidence

Finance and cleanup checks remained clean after the UAT sequence:

```text
finance reconciliation: PASS
wallet_mismatches=0
invoice_mismatches=0
duplicate_payment_references=0
claimable_or_running_jobs=0
cloudmini_nonterminated_services=0
provider_mutation_routes_called=no
production_notification_delivery_called=no
```

Notification evidence remains inherited from T279/T280/T281:

```text
selected-host Telegram preflight: PASS by prior evidence
queued Telegram delivery: PASS by prior evidence
failure/retry classification: PASS by prior evidence
manual fallback path: documented
```

## Target Auth Smoke Credential Handling

The generic `dev-target-auth-rbac` smoke supports protected environment overrides so the selected target can rotate or replace seed credentials without code changes.

Supported override names:

```text
BILLING_TARGET_AUTH_SMOKE_CLIENT_EMAIL
BILLING_TARGET_AUTH_SMOKE_CLIENT_PASSWORD
BILLING_TARGET_AUTH_SMOKE_ADMIN_EMAIL
BILLING_TARGET_AUTH_SMOKE_ADMIN_PASSWORD
```

Rules:

- Store override values only in the protected non-repo secret store for the selected dev/test environment.
- Do not commit the values to docs, task files, scripts, shell history, or PR text.
- The smoke output must continue to exclude passwords, cookies, session tokens, DSNs, provider payloads, and credentials.
- When overrides are absent, the smoke falls back to the original dev seed defaults for local developer environments.

## Not Verified

- Human Client, Reseller, or Admin sign-off was not collected in this consolidated task.
- Full 2FA-satisfied admin browser navigation was not executed; the live 2FA gate and admin API probes passed.
- Admin credential reveal mutation, admin service lifecycle mutations, job retry/cancel/manual-review mutations, support ticket mutations, and real Cloudmini provisioning were not executed.
- Production notification delivery was not executed.
- Cross-reseller isolation against a second real reseller tenant was not tested because the selected target seed currently has one reseller tenant.

## Decision

```text
Client UAT: PASS for automated selected-environment evidence
Reseller UAT: PASS for automated selected-environment evidence
Admin UAT: PASS for automated selected-environment evidence
Open P0 bugs: 0
Open P1 bugs: 0
Production GO: not approved by this packet
```

Next recommended step:

```text
Run dev-target-auth-rbac on the selected target with protected credential overrides configured, then record a small follow-up evidence note if the selected target credentials differ from seed defaults.
```
