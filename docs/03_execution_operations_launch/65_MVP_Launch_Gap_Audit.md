# 65 - MVP Launch Gap Audit

**Date:** 2026-05-14
**Scope:** Current-code audit against MVP scope and launch Go/No-Go gates.
**Decision:** NO-GO for pilot launch. T189-T206 closed many repo/local implementation gaps; T208-T210 define the evidence packets for the remaining external/provider/staging/owner blockers, but the required proof is still missing.

## Source Documents

- `docs/03_execution_operations_launch/26_MVP_Scope_Lock_And_Non_Goals.md`
- `docs/03_execution_operations_launch/33_Launch_Checklist_And_Go_No_Go_Criteria.md`
- `docs/03_execution_operations_launch/69_Pilot_Go_No_Go_Record.md`
- `docs/03_execution_operations_launch/70_Launch_Evidence_Completion_Packet.md`
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
| Tenant isolation P0 tests | `partial` | T189 session baseline merged; T204 `dev-api` smoke passed RBAC negative checks; `cmd/api/runtime_protection_test.go` covers runtime protection. | Repeat the gate against approved staging/session inputs before changing launch decision to GO. |
| Ledger reconciliation | `partial` | T194 refunds/adjustments, T195 daily reconciliation, and T204 billing mutation smoke passed in local/dev. | Assign finance launch owner and capture day-one reconciliation ownership. |
| Checkout debit/reservation/provisioning flow | `partial` | T196 reservation concurrency and T204 fake-provider billing/provisioning flow passed. | Real provider pilot remains blocked by T199 and doc 66. |
| Provisioning idempotency test | `partial` | Local fake provider and worker paths are covered by T199 and T204. | Real provider idempotency, timeout, quota, SKU mapping, and cleanup evidence are missing. |
| Credential encryption/redaction | `partial` | T192 encrypted credential storage and T193 controlled reveal/audit merged. | Prove target-environment secret/key handling and reveal audit access before GO. |
| Admin 2FA | `partial` | T190 TOTP setup/verify, encrypted TOTP secret storage, 2FA-satisfied sessions, admin route enforcement, and redacted audit events merged. | Verify production admin enrollment and policy enforcement in the target environment. |
| Backup restore test | `partial` | T203 local restore drill passed and repeatable runbook/script exist. T209 adds the shared staging evidence packet required for launch proof. | Repeat restore drill against approved shared staging/non-production target with redacted operator evidence and Ops/QA sign-off. |
| Provider pilot test | `blocked` | Local fake provider and readiness APIs exist; T199 and doc 66 explicitly block real provider sandbox. T208 adds the redacted evidence packet that must be filled before reconsidering GO. | Approved sandbox account, credentials path, quota/cost limit, SKU mapping, timeout policy, redacted examples, cleanup owner, and real pilot run evidence are missing. |
| Support SOP readiness | `partial` | SOP docs exist; T201 support/abuse backend basics merged; T200 notification foundation exists. T210 defines the notification/fallback evidence fields, T222 defines the manual fallback runbook packet, and T244 records Admin-owned manual fallback evidence for the pilot scope. | Production SMTP/Telegram remains unproven; keep manual fallback only if Admin coverage/SLA is active for the pilot window. |
| Incident owner assignment | `partial` | Incident roles are documented in `docs/03_execution_operations_launch/31_Incident_Response_And_Disaster_Recovery_Playbook.md`; T241 assigns Admin to all launch-day owner roles and T244 records manual fallback ownership for Support/Ops/Security. | Final launch window, escalation acceptance, and remaining owner sign-offs are not complete. |

## MVP Done Flow Matrix

| MVP done item | Status | Evidence | Follow-up |
| --- | --- | --- | --- |
| Client tops up wallet manually | `done` | T204 billing mutation smoke created and approved top-up in local/dev. | Re-run against approved staging before launch. |
| Admin approves top-up | `done` | T204 billing mutation smoke and T189 auth/session baseline cover local/dev admin path. | Prove notification delivery or manual fallback before launch. |
| Client buys a VPS/proxy | `partial` | T204 order, checkout, invoice, wallet payment, and fake-provider provisioning passed. | Real provider path remains blocked by T199 and doc 66. |
| Client debit and reseller settlement debit are correct | `done` | T194/T195 and T204 local/dev billing mutation evidence cover append-only ledger and reconciliation. | Assign finance owner for day-one checks. |
| Stock is reserved safely | `done` | T196 provider inventory counters, atomic reserve SQL, expiry release SQL, and concurrency tests merged. | Re-run full gate on approved staging inputs. |
| Service is provisioned | `partial` | Local fake worker path passed in T204. | Real provider sandbox/pilot evidence is missing. |
| Credentials are masked by default | `done` | T192/T193 merged encrypted storage, metadata, reveal controls, and audit; frontend redaction guards remain. | Prove target secret/key handling before GO. |
| Credential reveal is audited | `done` | T193 reveal API, rate limit, no-store response, tenant/owner scoping, RBAC, and audit merged. | Verify audit access in the target environment. |
| Service renewal works | `done` | T197/T198 lifecycle primitives, APIs, scheduler, and jobs merged; T206 direct client renewal API/UI action merged with wallet debit, invoice/payment records, lifecycle renewal, audit, and client UI action. | Re-run renewal path against approved staging/full E2E inputs before GO. |
| Expire/suspend/terminate policy works | `done` | T197/T198 lifecycle transition and scheduler jobs merged. | Re-run smoke on approved staging inputs. |
| Finance reconciliation passes | `done` | T195 daily reconciliation report and T204 local/dev full gate passed. | Assign finance launch owner. |
| Cross-tenant access attempts fail | `partial` | T189 auth/session baseline and T204 RBAC negative checks passed in local/dev. | Repeat against target session/auth configuration. |

## Category Gate Summary

| Gate | Status | Current evidence | Required next task |
| --- | --- | --- | --- |
| Auth and sessions | `partial` | T189 merged Argon2id password verification, DB-backed sessions, cookie login/logout, global session middleware, and dev-only actor headers. | Repeat launch gate on approved target auth/session inputs. |
| Admin 2FA | `partial` | T190 merged TOTP setup/verify, encrypted secret storage, 2FA-satisfied session state, route enforcement, and redacted audit. | Verify target admin enrollment and policy enforcement. |
| Login protection and password reset | `done` | T191 merged DB-backed login/password-reset rate limits, hashed reset tokens, reset confirmation, session revocation, and routes. | Normal target-environment verification. |
| Credential storage and reveal | `partial` | T192/T193 merged encrypted credential storage, controlled reveal, rate limit, no-store response, RBAC, and audit. | Verify target secret/key handling and reveal audit access. |
| Refunds and adjustments | `done` | T194 merged append-only refund/adjustment ledger routes, idempotency conflict checks, reasons, actor metadata, and audit. | Normal finance owner sign-off. |
| Daily reconciliation | `done` | T195 merged read-only daily reconciliation report with wallet balance and mismatch checks. | Assign finance launch owner. |
| Reservation TTL and concurrency | `done` | T196 merged provider inventory counters, atomic reservation SQL, expiry release SQL, and concurrency tests. | Re-run launch gate on approved staging inputs. |
| Service lifecycle | `done` | T197/T198 merged lifecycle transitions, admin/reseller APIs, scheduler, and worker jobs; T206 added the client renewal API/UI path. | Re-run lifecycle and renewal paths on approved staging inputs. |
| Provider sandbox | `blocked` | Fake provider, local sandbox contract, T199/doc 66 no-go evidence, and T208 evidence packet requirements exist. | External provider sandbox intake and real pilot evidence are still required. |
| Notifications | `partial` | T200 adds `notifications` schema, notification module, redacted launch-critical event builders, and local delivery runner; T210 defines production delivery/fallback evidence; T222 defines the manual fallback owner/SLA/evidence packet; T244 records an Admin-owned manual fallback drill for the pilot scope. | Production SMTP/Telegram remains unproven; manual fallback depends on Admin availability and SLA during the approved pilot window. |
| Support and abuse backend | `partial` | T201 support/abuse backend basics merged; SOP docs exist; T241 assigns Admin as Support owner; T244 records manual fallback coverage and SLA for the pilot scope. | Production notification delivery remains unproven and launch coverage depends on Admin availability. |
| Frontend production integration | `partial` | T202 integrated production API paths and frontend smokes; T206 wired the client service renewal UI to the production API; T210 requires staging/full E2E evidence. | Repeat frontend smoke/full E2E against approved staging or sandbox-equivalent inputs. |
| Backup/restore | `partial` | T203 local restore drill passed with repeatable script/runbook; T209 defines the shared staging evidence packet. | Execute approved shared staging/non-production evidence before GO. |
| Full E2E quality gate | `partial` | T204 local/dev full gate passed across backend, DB/API/billing smokes, and frontend checks. | Repeat against approved staging/sandbox-equivalent inputs or record signed exception. |
| Final Go/No-Go | `done` | T205 record exists in `docs/03_execution_operations_launch/69_Pilot_Go_No_Go_Record.md`. | Decision is NO-GO until blockers are cleared. |

## Recommended Execution Order

1. Complete `docs/03_execution_operations_launch/70_Launch_Evidence_Completion_Packet.md`, including doc 66 provider evidence and doc 67 restore evidence.
2. Include the T206 client renewal path in approved staging/full E2E evidence or record a signed staging-equivalent exception.
3. Re-run T205 and change decision only when every P0 row has passing evidence or acceptable non-P0 mitigation.

## Verification Scope For This Audit

This audit combines repository inspection with linked task evidence through T210. Local/dev smokes and CI evidence are recorded in the task files and runbooks, but this audit does not claim that approved shared staging, production, notification delivery, launch owner, target-environment, or real provider sandbox tests were executed.
