# 26 — MVP Scope Lock & Non-Goals

**Version:** v1.3  
**Date:** 2026-04-22  
**Scope:** MVP boundary control.

## 1. MVP statement

MVP is not a smaller version of the final dream. MVP is the first safe end-to-end tunnel through the mountain.

```text
The MVP proves that the platform can safely sell, charge, reserve, provision, renew, suspend, terminate, support, and audit one real VPS/proxy product flow.
```

## 2. MVP goal

```text
Build a white-label VPS/proxy platform where:
- platform admin manages tenants, provider sources, wallet approvals, and audit;
- reseller manages clients, catalog pricing, wallet, and storefront basics;
- client tops up, buys, views service, reveals credential, renews;
- wallet/ledger/settlement are correct;
- provisioning is idempotent and safe;
- tenant isolation and credential security are enforced.
```

## 3. MVP includes

### Core platform

```text
- Auth
- Tenant management
- Domain-to-tenant mapping basic
- RBAC
- Admin 2FA
- Audit log
- Basic notification
```

### Catalog

```text
- Master product
- Master plan
- Provider source
- Plan-source mapping
- Tenant catalog clone
- Reseller price override
- Price/cost/policy snapshot
```

### Wallet and billing

```text
- Client wallet
- Reseller wallet
- Immutable ledger
- Manual top-up request
- Manual approval/rejection
- Reseller settlement debit
- Refund/adjustment basic
- Daily reconciliation report
```

### Checkout and reservation

```text
- Client checkout
- Balance validation
- Reseller balance validation
- Atomic reservation
- Reservation expiry
- Out-of-stock handling
```

### Provisioning

```text
- Queue-based provisioning
- Provider adapter interface
- At least one VPS source implementation
- At least one proxy/manual source implementation
- Idempotency key
- Safe/unsafe retry classification
- Manual review for uncertain provider state
```

### Lifecycle

```text
- active
- expired
- grace
- suspended
- terminated
- renew
- manual suspend/unsuspend if supported
- lifecycle events
```

### Portals

```text
- Admin portal basic
- Reseller portal basic
- Client portal basic
- Wallet pages
- Catalog pages
- Order pages
- Service detail
- Credential reveal with audit
```

### Operations

```text
- Support SOP basic
- Abuse flag basic
- Provider health basic
- Launch checklist
- Backup/restore test
```

## 4. MVP excludes

Do not include unless the project owner explicitly changes scope.

```text
- Native mobile app
- Advanced affiliate system
- Complex coupon/promotion engine
- AI fraud scoring
- KYC automation
- Auto crypto confirmation before risk controls
- Multi-currency reseller payout
- Public marketplace
- Complex tax engine
- Advanced BI dashboard
- Full custom ticket system
- Complex Kubernetes/autoscaling if unnecessary
- Usage-based minute/hour billing before monthly billing is stable
- Postpaid reseller credit
```

## 5. Provider scope lock

MVP does not need ten providers. MVP needs one or two sources that work safely.

```text
MVP target:
- 1 VPS source
- 1 proxy source
- manual provider fallback
```

Do not add a new provider if:

```text
- Adapter contract is not stable.
- Idempotency is not proven.
- Manual review flow is missing.
- Provider onboarding score is below threshold.
- Support team cannot handle that provider's failure modes.
```

## 6. Security items that cannot be cut

```text
- Admin 2FA
- Tenant isolation
- Credential encryption
- Credential reveal audit
- Rate limits for login/top-up/checkout/reveal credential
- Provider secret protection
- No plaintext credential in logs
```

These are not nice-to-have. They are safety rails.

## 7. Billing items that cannot be cut

```text
- Immutable ledger
- Client wallet
- Reseller wallet
- Settlement debit
- Refund/adjustment entry
- Reconciliation report
```

Do not introduce postpaid/credit limit in MVP. Prepaid first keeps the system honest.

## 8. Abuse scope lock

MVP does not need advanced automation, but it needs manual control.

```text
- Abuse flag
- Abuse case
- Suspend reason
- Evidence log
- Provider takedown workflow
- Basic blacklist/risk flag
```

## 9. Change control questions

Every new scope request must answer:

```text
1. Which MVP goal does this serve?
2. Can we launch safely without it?
3. Does it touch money, tenant, provisioning, or security?
4. Does it delay the critical path?
5. Who will test and operate it?
```

If the answer is unclear, move it to a later phase.

## 10. MVP done statement

MVP is done when one client under one reseller can:

```text
1. Top up wallet manually.
2. Receive approval.
3. Buy a VPS/proxy.
4. Trigger correct client debit and reseller settlement debit.
5. Reserve stock safely.
6. Get a provisioned service.
7. View masked credentials.
8. Reveal credential with audit.
9. Renew service.
10. Expire/suspend/terminate according to policy.
11. Pass finance reconciliation.
12. Fail cross-tenant access attempts.
```

## 11. Anti-patterns

```text
- Pretty UI before money core passes.
- Adding providers because they look cheap.
- Postpaid reseller credit too early.
- Coupon engine before ledger safety.
- Public launch before private pilot.
- Advanced reports before reconciliation is correct.
```
