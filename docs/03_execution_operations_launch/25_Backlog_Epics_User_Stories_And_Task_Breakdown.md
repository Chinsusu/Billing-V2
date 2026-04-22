# 25 — Backlog, Epics, User Stories & Task Breakdown

**Version:** v1.3  
**Date:** 2026-04-22  
**Scope:** Product execution backlog for dev planning.

## 1. Purpose

This document converts the blueprint into a backlog that can be imported into Jira, Linear, ClickUp, Trello, or GitHub Projects.

Backlog format:

```text
Epic
  User story
    Acceptance criteria
    Tasks
    Priority
    Dependencies
```

## 2. Priority standard

| Priority | Meaning | Example |
|---|---|---|
| P0 | Cannot launch without it | Tenant isolation, ledger, idempotency |
| P1 | Required for stable MVP | Notifications, basic reports |
| P2 | Useful after MVP | Advanced analytics, coupons |
| P3 | Explicitly not MVP | Native app, affiliate engine |

## 3. Epic A — Tenant, Auth & RBAC

### A1 — Platform admin creates reseller tenant

```text
As a platform admin,
I want to create a reseller tenant,
so that a reseller can manage their own storefront, clients, wallet, and pricing.
```

Acceptance:

```text
- Admin creates tenant with name, slug, status, owner email.
- System creates reseller owner user.
- System creates default reseller wallet.
- Tenant slug/domain is unique.
- Audit event tenant.created is recorded.
```

Tasks:

```text
- tenants migration
- tenant creation service
- reseller owner creation flow
- admin tenant API
- admin tenant UI
- tenant.created audit event
```

Priority: P0

### A2 — Tenant context middleware

Acceptance:

```text
- Request tenant is resolved from domain/session/token.
- tenant_id from request body is ignored for authorization.
- Tenant resources are always queried with tenant scope.
- Cross-tenant access returns 403 or 404.
```

Tasks:

```text
- tenant resolver
- tenant middleware
- repository tenant filter
- cross-tenant API tests
```

Priority: P0

### A3 — RBAC permission enforcement

Acceptance:

```text
- Role/permission matrix exists.
- APIs check permission before write actions.
- UI hides actions without permission.
- Permission denied returns stable error code.
```

Tasks:

```text
- roles table
- permissions table
- user_roles table
- permission middleware
- RBAC test suite
```

Priority: P0

## 4. Epic B — Wallet, Ledger & Settlement

### B1 — Manual top-up request

Acceptance:

```text
- Client/reseller creates top-up request.
- Request enters pending state.
- Admin/reseller approver can approve/reject according to policy.
- Approved request creates one ledger credit.
- Rejected request changes no wallet balance.
```

Tasks:

```text
- topup_requests migration
- create top-up API
- approve/reject API
- wallet ledger credit service
- top-up UI
- audit events
```

Priority: P0

### B2 — Immutable wallet ledger

Acceptance:

```text
- Ledger entry cannot be updated/deleted via application.
- Balance is derived from ledger or reconciled to ledger.
- Adjustment creates new ledger entry.
- Every wallet balance change has ledger reason/reference.
```

Tasks:

```text
- wallet_ledger_entries migration
- wallet balance service
- ledger append-only guard
- daily reconciliation job
- finance exception report
```

Priority: P0

### B3 — Reseller settlement on client checkout

Acceptance:

```text
- Client wallet must cover selling_price.
- Reseller wallet must cover reseller_cost.
- Checkout fails if reseller wallet is insufficient.
- Client debit and reseller settlement debit are atomic.
- Order snapshot stores selling_price, reseller_cost, margin, tenant.
```

Tasks:

```text
- settlement calculation service
- checkout transaction wrapper
- reseller wallet debit ledger type
- insufficient reseller balance error
- settlement tests
```

Priority: P0

## 5. Epic C — Catalog, Pricing & Source

### C1 — Master product and plan

Acceptance:

```text
- Admin creates product and plan.
- Plan has billing cycle, base price, reseller cost, status.
- Product/plan changes are audited.
- Disabled plan cannot be purchased.
```

Tasks:

```text
- products migration
- plans migration
- catalog CRUD API
- admin catalog UI
- pricing audit events
```

Priority: P0

### C2 — Tenant catalog clone and price override

Acceptance:

```text
- Reseller clones allowed master plans.
- Reseller sets selling price.
- Margin floor is enforced when enabled.
- Old orders preserve old price snapshot.
```

Tasks:

```text
- tenant catalog tables
- price override API
- margin validation
- reseller catalog UI
- snapshot tests
```

Priority: P0

### C3 — Source assignment and capability snapshot

Acceptance:

```text
- Plan can be linked to one or more sources.
- Source declares capability mask.
- Checkout stores selected source snapshot.
- Service action buttons use service capability snapshot, not current provider guess.
```

Tasks:

```text
- provider_sources migration
- plan_sources migration
- capability schema
- source selection logic
- capability snapshot storage
```

Priority: P0

## 6. Epic D — Order, Checkout & Reservation

### D1 — Client checkout

Acceptance:

```text
- Validates user, tenant, plan, source, stock, client balance, reseller balance.
- Creates order, order item snapshot, reservation, ledger entries, provisioning job.
- Returns stable error codes for known failures.
```

Tasks:

```text
- orders migration
- order_items migration
- checkout service
- checkout API
- checkout UI
- error mapping
```

Priority: P0

### D2 — Atomic inventory reservation

Acceptance:

```text
- Concurrent checkouts cannot oversell.
- Reservation expires after TTL.
- Expired reservation releases stock once.
- Allocated reservation cannot be released by expiry job.
```

Tasks:

```text
- reservations migration
- atomic lock/decrement logic
- reservation expiry job
- stock counters
- concurrency tests
```

Priority: P0

## 7. Epic E — Provisioning

### E1 — Provisioning worker

Acceptance:

```text
- Job is created only after valid checkout.
- idempotency_key is unique.
- Worker processes job by status.
- Only safe errors are retried.
- Unsafe/uncertain errors go to manual review.
```

Tasks:

```text
- provisioning_jobs migration
- worker framework
- retry policy
- manual review status
- provider request log
```

Priority: P0

### E2 — Provider adapter v1

Acceptance:

```text
- Adapter implements checkHealth, checkStock, provision, getStatus, suspend, terminate.
- Provider errors are normalized.
- Credentials are encrypted before storage.
- Provider capability response is stable.
```

Tasks:

```text
- adapter interface
- provider error enum
- provider account secret config
- adapter test harness
- first provider implementation
```

Priority: P0

### E3 — Partial success handling

Acceptance:

```text
- Timeout after provider create is not blindly retried.
- external_request_id/resource_id is stored if available.
- Admin can reconcile uncertain external resource.
- Duplicate resource prevention test exists.
```

Tasks:

```text
- unsafe retry classifier
- provider_resource_mappings table
- manual review UI
- reconcile action
- duplicate prevention tests
```

Priority: P0

## 8. Epic F — Service Lifecycle

### F1 — Service activation

Acceptance:

```text
- Provider success creates service.
- Service stores term_start, term_end, service_status.
- Credential is encrypted.
- Client can view service detail under tenant scope.
```

Priority: P0

### F2 — Renewal

Acceptance:

```text
- Active renewal adds term from old_term_end.
- Grace renewal follows configured policy.
- Wallet debit and ledger entry are created.
- Provider renew action is called only when supported/required.
```

Priority: P0

### F3 — Expire, suspend, terminate

Acceptance:

```text
- Expired service enters grace.
- Grace exceeded triggers suspend/terminate policy.
- Manual suspend requires reason.
- Termination is irreversible or requires elevated permission.
```

Priority: P0

## 9. Epic G — Portals

| Story | Acceptance summary | Priority |
|---|---|---|
| G1 Admin portal | Admin manages tenants, catalog, wallets, orders, services, provider health, failed jobs | P1 |
| G2 Reseller portal | Reseller manages clients, wallet, catalog pricing, orders, services within tenant | P1 |
| G3 Client portal | Client can top up, buy, view, renew, reveal credentials with audit | P1 |

## 10. Epic H — Notifications

Acceptance:

```text
- Top-up approved/rejected notifications.
- Service activated/expiring/suspended notifications.
- Reseller low balance alert.
- Provider degraded alert.
- Manual review required alert.
```

Priority: P1

## 11. Epic I — QA & Launch

Acceptance:

```text
- P0 tenant isolation tests pass.
- P0 ledger tests pass.
- Checkout/reservation/provisioning tests pass.
- Credential security tests pass.
- Launch checklist is signed.
```

Priority: P0

## 12. Backlog hygiene rules

```text
- No story enters sprint without acceptance criteria.
- No money/tenant/provisioning story merges without tests.
- No new feature breaks MVP scope lock.
- UI is not done until API, error, permission, and audit are done.
```
