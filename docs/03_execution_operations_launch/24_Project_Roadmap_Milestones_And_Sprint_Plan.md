# 24 — Project Roadmap, Milestones & Sprint Plan

**Version:** v1.3  
**Date:** 2026-04-22  
**Scope:** Execution roadmap from blueprint to MVP/pilot launch.  
**Audience:** Founder, product owner, engineering lead, QA, ops, finance.

## 1. Purpose

This document turns the product and technical handoff into an execution path. The goal is to prevent the team from building the easiest screens first while leaving the dangerous core — tenant isolation, wallet ledger, settlement, reservation, provisioning, and credential security — until too late.

The project should be managed by **end-to-end safety flows**, not by isolated features.

```text
Wrong question: Is the dashboard done?
Right question: Can one client safely top up, buy, get provisioned, renew, expire, and be audited without money/tenant/provisioning bugs?
```

## 2. Roadmap principle

Build in this order:

```text
1. Identity and tenant boundary
2. Money and ledger
3. Catalog and pricing
4. Order and reservation
5. Provisioning and service lifecycle
6. Portals
7. QA, operations, and pilot
8. Controlled launch
```

Do not build advanced reseller growth features before the core money/provider system is safe.

## 3. Milestone map

| Milestone | Name | Output | Launch relevance |
|---:|---|---|---|
| M0 | Foundation | Repo/app skeleton, DB, auth, tenant, RBAC, audit base | Mandatory |
| M1 | Money Core | Wallet, ledger, manual top-up, reseller settlement | Mandatory |
| M2 | Catalog & Order | Product/plan/source, tenant price, checkout, reservation | Mandatory |
| M3 | Provisioning | Queue, adapter, idempotency, manual review, service activation | Mandatory |
| M4 | Lifecycle & Portals | Admin/reseller/client basic portal, renew, expire, suspend | Mandatory |
| M5 | QA & Operational Readiness | P0 tests, reconciliation, incident/support/abuse SOP | Mandatory |
| M6 | Private Beta | Limited reseller/client pilot with real orders | Mandatory before public |
| M7 | Public MVP Launch | Controlled public rollout | Only after gates pass |

## 4. Suggested sprint plan

Assume 2-week sprints. If the team is small, stretch each sprint; do not remove safety gates.

### Sprint 0 — Project setup and architecture lock

**Goal:** The team can run the project locally and understands the non-negotiable architecture.

```text
Deliverables:
- Environment setup
- Database migration baseline
- Backend app skeleton
- Frontend portal skeleton
- Worker/queue skeleton
- Provider mock skeleton
- Logging/audit base pattern
- Coding standards and PR checklist
```

**Definition of done:** local dev runs API + worker + DB + queue + frontend + provider mock.

### Sprint 1 — Tenant, auth, RBAC, audit

```text
Deliverables:
- Auth/login/session
- Tenant model
- Domain-to-tenant resolver basic
- RBAC middleware
- Admin/reseller/client role model
- Audit log write utility
- Tenant isolation test harness
```

**P0 tests:** user from Tenant A cannot read Tenant B resource by guessed ID.

### Sprint 2 — Wallet, ledger, manual top-up

```text
Deliverables:
- Wallets
- Immutable ledger
- Client top-up request
- Reseller top-up request
- Manual approve/reject
- Balance recalculation
- Finance exception report draft
```

**P0 tests:** approved top-up creates one credit; rejected top-up creates none; ledger entry cannot be edited via app.

### Sprint 3 — Reseller settlement and catalog pricing

```text
Deliverables:
- Master catalog
- Product/plan/source
- Tenant catalog clone
- Reseller selling price override
- Margin floor
- Price/cost/policy snapshots
- Reseller settlement calculation service
```

**P0 tests:** client wallet and reseller wallet are validated before checkout; insufficient reseller balance blocks provisioning.

### Sprint 4 — Order, checkout, reservation

```text
Deliverables:
- Order creation
- Order item snapshot
- Atomic stock reservation
- Reservation TTL expiry job
- Checkout API
- Checkout error codes
- Concurrency tests
```

**P0 tests:** 10 concurrent checkouts for 1 remaining stock allocate only 1 reservation.

### Sprint 5 — Provisioning core

```text
Deliverables:
- Provisioning jobs
- Worker processing
- Adapter interface
- Provider request log
- Idempotency key
- Safe/unsafe retry policy
- Manual review state
```

**P0 tests:** provider timeout after create request does not blindly retry.

### Sprint 6 — Service lifecycle

```text
Deliverables:
- Service activation
- Credential encrypted storage
- Credential reveal audit
- Renew active service
- Expiry/grace/suspend/terminate jobs
- Lifecycle event history
```

**P0 tests:** renew term calculation is deterministic; expired service transitions correctly; credential plaintext is not logged.

### Sprint 7 — Admin and reseller portals

```text
Deliverables:
- Admin dashboard basic
- Tenant management
- Catalog management
- Wallet/top-up approval
- Failed provisioning/manual review screen
- Reseller dashboard
- Reseller pricing
- Reseller client/service view
```

**P0 tests:** UI may hide buttons, but API must still enforce permission.

### Sprint 8 — Client portal and notifications

```text
Deliverables:
- Client registration/login
- Wallet page
- Catalog browse
- Checkout flow
- Service list/detail
- Renew
- Notification templates
- Reseller low-balance alert
```

**P0 tests:** client sees only own services and credentials.

### Sprint 9 — QA hardening and operational readiness

```text
Deliverables:
- P0/P1 QA test suite
- Finance reconciliation SOP executed on staging
- Backup/restore drill
- Incident playbook drill
- Support macros loaded
- Abuse workflow tested
```

**P0 tests:** all No-Go criteria in file 33 pass.

### Sprint 10 — Private beta

```text
Deliverables:
- 2–3 pilot resellers
- 10–30 pilot clients
- limited catalog
- daily reconciliation
- weekly feedback review
- provider score review
```

**Exit gate:** no unresolved ledger mismatch, cross-tenant issue, duplicate provisioning incident, or credential leak.

## 5. Critical path dependencies

```text
Tenant middleware -> every tenant API
Ledger -> checkout/refund/settlement
Catalog snapshot -> order snapshot
Reservation -> provisioning job
Provisioning job -> service activation
Credential encryption -> service detail
QA P0 -> launch gate
```

If a dependency is not ready, downstream work should use mock/stub only and must not be treated as production-ready.

## 6. Role ownership

| Area | Owner | Notes |
|---|---|---|
| Product scope | Product Owner | Controls scope and acceptance |
| Tenant/RBAC | Backend Lead | Security-critical |
| Wallet/ledger | Backend + Finance | Money-critical |
| Provisioning | Backend + Ops | Provider-critical |
| Portal UI | Frontend Lead | Must obey RBAC/capability |
| QA | QA Lead | Owns P0/P1 suite |
| DevOps | DevOps Lead | Deploy, backup, monitoring |
| Support SOP | Ops Lead | Macros, escalation |
| Abuse SOP | Abuse/Compliance Owner | Takedown policy |

## 7. Sprint entry criteria

A story should not enter sprint unless it has:

```text
- clear acceptance criteria
- role and permission requirement
- tenant scope rule
- audit event if critical
- error codes if API-facing
- test expectation
- dependency identified
```

## 8. Sprint exit criteria

A sprint is not done until:

```text
- code merged
- migrations reviewed
- P0 tests pass
- no secret/credential in logs
- error codes stable
- audit events written for critical actions
- QA verifies acceptance criteria
- known risks documented
```

## 9. Project risk board

| Risk | Early signal | Blocking action |
|---|---|---|
| Scope creep | Marketing/affiliate/coupon requests before MVP | Use file 26 Non-Goals |
| Money bug | Wallet balance differs from ledger | Stop launch, reconcile |
| Tenant leak | Endpoint queries by ID only | Block merge, add tenant tests |
| Duplicate provisioning | Timeout followed by blind retry | Pause provider, manual review |
| Provider instability | Job pending/failure rises | Circuit breaker |
| Support overload | Same ticket repeats | Improve macros/notifications |
| Launch too early | P0 still open | No-Go via file 33 |

## 10. Reporting cadence

```text
Daily:
- blockers
- failed provisioning
- ledger mismatch
- P0 bugs
- sprint burn status

Weekly:
- milestone progress
- scope change requests
- provider readiness
- QA pass rate
- pilot readiness

Pre-launch:
- Go/No-Go review
- incident drill
- backup restore confirmation
```

## 11. Golden rule

```text
Do not ask only: “Is this feature done?”
Ask: “Is this flow safe end-to-end under real money, real tenant, real provider failure?”
```
