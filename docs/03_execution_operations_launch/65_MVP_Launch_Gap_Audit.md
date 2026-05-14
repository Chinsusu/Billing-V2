# 65 - MVP Launch Gap Audit

**Date:** 2026-05-13
**Scope:** Current-code audit against MVP scope and launch Go/No-Go gates.
**Decision:** NO-GO for pilot launch until the P0 partial, missing, and blocked items below are closed and re-verified.

## Source Documents

- `docs/03_execution_operations_launch/26_MVP_Scope_Lock_And_Non_Goals.md`
- `docs/03_execution_operations_launch/33_Launch_Checklist_And_Go_No_Go_Criteria.md`
- Current backend, frontend, migration, smoke, CI, and runbook files in this repository.

## Status Legend

| Status | Meaning |
| --- | --- |
| `done` | Current repo evidence is enough for the checklist item, subject to normal CI. |
| `partial` | Some implementation or documentation exists, but launch-required behavior or proof is incomplete. |
| `missing` | No production-ready implementation or launch evidence was found. |
| `blocked` | The item depends on an external decision, provider/account evidence, or operator sign-off. |

## P0 No-Go Matrix

| P0 launch item | Status | Current evidence | Gap / follow-up |
| --- | --- | --- | --- |
| Tenant isolation P0 tests | `partial` | Tenant context and RBAC checks exist in `internal/modules/tenant/context.go`, `internal/modules/rbac/authorizer.go`, `cmd/api/runtime_protection_test.go`, and `cmd/smoke/api_rbac.go`. | Production auth/session is still header-based through `internal/modules/identity/http_middleware.go`; close with T189 and prove full E2E in T204. |
| Ledger reconciliation | `partial` | Append-only wallet ledger schema exists in `migrations/0011_create_wallet_ledger_entries.sql`; payment reconciliation read path exists in `internal/modules/payment/reconciliation.go`; smoke covers `/admin/payment-reconciliation`. | Refund/adjustment behavior and daily reconciliation report are not complete; close with T194 and T195. |
| Checkout debit/reservation/provisioning flow | `partial` | Wallet invoice payment runs in a DB transaction in `internal/modules/payment/postgres_store.go`; paid orders queue provisioning in `internal/modules/order/payment_provisioning.go` and `internal/modules/order/provisioning_queue.go`; local billing smoke exists in `cmd/smoke/billing_mutation.go`. | Reservation TTL/concurrency proof and full launch E2E proof remain open; close with T196 and T204. |
| Provisioning idempotency test | `partial` | Job idempotency constraints exist in `migrations/0004_create_outbox_jobs.sql`; provider retry taxonomy and fake sandbox contract exist in `internal/modules/provider/*`; provisioning worker maps unsafe outcomes to manual review in `internal/modules/order/provider_provisioning_worker.go`; T199 evidence exists in `docs/03_execution_operations_launch/66_Provider_Sandbox_Readiness_Evidence.md`. | Local fake paths are proven, but real sandbox provider intake is still missing; do not pilot real provider provisioning. |
| Credential encryption/redaction | `missing` | Logger redaction keys exist in `internal/platform/logger/logger.go`; provider result has a credential envelope type in `internal/modules/provider/operation.go`. | No service credential storage, encrypted-at-rest migration, reveal API, or reveal audit flow exists; close with T192 and T193. |
| Admin 2FA | `missing` | `users.two_factor_status` exists in `migrations/0002_create_identity_rbac.sql` and identity read models expose the field. | No 2FA setup/challenge/enforcement path exists for privileged admin access; close with T190. |
| Backup restore test | `partial` | DR guidance exists in `docs/03_execution_operations_launch/31_Incident_Response_And_Disaster_Recovery_Playbook.md`; local DB smoke exists in `docs/05_development_standards/55_Local_Development_Runbook.md`; repeatable local/sandbox drill exists in `docs/03_execution_operations_launch/67_Backup_Restore_Drill_Runbook.md` and `scripts/backup_restore_drill.sh`. | Launch evidence still requires an approved non-production run with redacted operator evidence; close final proof with T204/T205. |
| Provider pilot test | `blocked` | Local fake provider and readiness APIs exist in `internal/modules/provider/*`, `internal/modules/catalog/readiness.go`, `docs/05_development_standards/58_Provisioning_Ops_Readiness_Checklist.md`, and T199 no-go evidence in `docs/03_execution_operations_launch/66_Provider_Sandbox_Readiness_Evidence.md`. | Approved sandbox provider credentials, redacted provider evidence, and real provider pilot run are still not present; keep real provider pilot blocked. |
| Support SOP readiness | `partial` | Support and abuse SOP docs exist in `docs/03_execution_operations_launch/29_Customer_Support_SOP_And_Macro_Templates.md` and `docs/03_execution_operations_launch/32_Abuse_Compliance_Takedown_SOP.md`; frontend has demo ticket screens; T200 adds notification foundation records and redacted event builders. | Backend support/abuse records, tenant/RBAC controls, support-specific notification wiring, and audit behavior are missing; close with T201. |
| Incident owner assignment | `blocked` | Incident roles are documented in `docs/03_execution_operations_launch/31_Incident_Response_And_Disaster_Recovery_Playbook.md`. | Named launch-day owners and final Go/No-Go sign-off are not recorded; close with T205 after T189-T204 evidence exists. |

## MVP Done Flow Matrix

| MVP done item | Status | Evidence | Follow-up |
| --- | --- | --- | --- |
| Client tops up wallet manually | `partial` | Top-up API and approval flow exist in `internal/modules/wallet/*` and are used by `cmd/smoke/billing_mutation.go`. | Full launch E2E evidence in T204. |
| Admin approves top-up | `partial` | Approval/rejection model and audit support exist in `internal/modules/wallet/topup_review.go`. | Notification and production auth gaps remain: T189, T200. |
| Client buys a VPS/proxy | `partial` | Order, checkout, invoice, and wallet payment paths exist in `internal/modules/order`, `internal/modules/checkout`, `internal/modules/invoice`, and `internal/modules/payment`. | Reservation/concurrency and frontend production integration remain: T196, T202. |
| Client debit and reseller settlement debit are correct | `partial` | Purchase and reseller-cost ledger entry types exist in `migrations/0011_create_wallet_ledger_entries.sql`; wallet payment creates purchase ledger entries. | End-to-end settlement and reconciliation proof remain: T194, T195, T204. |
| Stock is reserved safely | `partial` | Reservation schema/model exist in `migrations/0007_create_order_tables.sql` and `internal/modules/order/reservation.go`; T196 adds provider inventory counters, atomic reserve SQL, expiry release SQL, and concurrency tests. | Full launch E2E proof remains: T204. |
| Service is provisioned | `partial` | Paid order provisioning queue and local fake worker exist in `internal/modules/order/provisioning_queue.go`, `internal/modules/order/provider_provisioning_worker.go`, and `cmd/worker/main.go`. | Provider sandbox readiness and full E2E proof remain: T199, T204. |
| Credentials are masked by default | `missing` | Frontend and smoke redaction guards exist for sensitive text. | Backend credential storage/reveal model is missing: T192, T193. |
| Credential reveal is audited | `missing` | Permission constant `service.credential.reveal` exists in `internal/modules/rbac/model.go`. | Reveal API, rate limit, and audit are missing: T193. |
| Service renewal works | `partial` | Provider capability includes renew and service statuses exist in `internal/modules/provider/capability.go` and `internal/modules/order/status.go`. | Lifecycle service behavior and scheduler jobs remain: T197, T198. |
| Expire/suspend/terminate policy works | `partial` | Service status transitions exist in `internal/modules/order/status.go`. | API, scheduler, audit, and provider action execution remain: T197, T198. |
| Finance reconciliation passes | `partial` | Payment reconciliation read path exists. | Daily report, mismatch detection, refunds/adjustments, and full E2E proof remain: T194, T195, T204. |
| Cross-tenant access attempts fail | `partial` | Tenant/RBAC checks and smoke exist. | Production session source and final E2E launch gate remain: T189, T204. |

## Category Gate Summary

| Gate | Status | Current evidence | Required next task |
| --- | --- | --- | --- |
| Auth and sessions | `missing` | Identity users, RBAC permissions, and header actor middleware exist. | T189 |
| Admin 2FA | `missing` | Database status field exists only. | T190 |
| Login protection and password reset | `missing` | No login/reset endpoints or rate limit primitives found. | T191 |
| Credential storage and reveal | `missing` | Redaction helpers exist; storage/reveal does not. | T192, T193 |
| Refunds and adjustments | `partial` | Ledger enum types and reason validation exist. | T194 |
| Daily reconciliation | `partial` | Payment reconciliation read model exists. | T195 |
| Reservation TTL and concurrency | `partial` | Provider inventory counters, reservation quantity, atomic reserve SQL, expiry release SQL, and concurrency tests exist. | Full E2E proof remains: T204 |
| Service lifecycle | `partial` | Status constants and simple transition helpers exist. | T197, T198 |
| Provider sandbox | `blocked` | Fake provider, local sandbox contract, and T199 no-go evidence exist. | External provider sandbox intake before T205 |
| Notifications | `partial` | T200 adds `notifications` schema, `internal/modules/notification`, redacted launch-critical event builders, and a local delivery runner that marks queued notifications sent without external SMTP/Telegram. | Production delivery channels and direct wiring into every launch flow remain outside this foundation; prove end-to-end notification behavior in T204/T205 before launch. |
| Support and abuse backend | `missing` | SOP docs and frontend demo screens exist only. | T201 |
| Frontend production integration | `partial` | Many screens attempt live API and fall back to mocks; smoke covers admin live/fallback paths. | T202 |
| Backup/restore | `partial` | Repeatable local/sandbox drill script and runbook exist. | Execute approved non-production evidence before T205 |
| Full E2E quality gate | `partial` | DB/API/billing/frontend smoke commands exist, but no single launch gate record exists. | T204 |
| Final Go/No-Go | `blocked` | Checklist template exists. | T205 |

## Recommended Execution Order

1. T189-T191: establish production auth/session, admin 2FA, rate limits, and password reset.
2. T192-T193: add encrypted credential storage and controlled reveal with audit.
3. T194-T196: close money, reconciliation, and reservation safety gaps.
4. T197-T198: complete lifecycle transitions and scheduler jobs.
5. T199-T201: record provider sandbox readiness/no-go evidence, notifications, and support/abuse backend.
6. T202-T204: finish frontend production integration and repeatable full E2E gate.
7. T205: execute final launch Go/No-Go only after all P0 rows are `done`.

## Verification Scope For This Audit

This audit is based on repository inspection only. It does not claim that database smoke, API smoke, browser smoke, or provider sandbox tests were executed during T188. Those validations belong to the linked implementation tasks and T204.
