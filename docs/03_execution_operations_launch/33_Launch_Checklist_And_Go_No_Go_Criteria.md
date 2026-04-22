# 33 — Launch Checklist & Go/No-Go Criteria

**Version:** v1.3  
**Date:** 2026-04-22  
**Scope:** Launch readiness gate for MVP/pilot.

## 1. Purpose

This document defines when the platform is allowed to launch. Do not launch because “it looks okay.” Launch only when core safety gates pass.

```text
A feature can be ugly and still launch.
A money, tenant, provisioning, or security bug cannot launch.
```

## 2. Launch stages

```text
Internal alpha: team only, mock + limited real provider.
Private beta: selected reseller/client, low volume, daily reconciliation.
Pilot launch: limited production with real payment/order/service.
Public MVP launch: broader rollout after gates pass.
```

## 3. P0 No-Go items

Launch is NO-GO if any item below fails:

```text
- tenant isolation P0 tests
- ledger reconciliation
- checkout debit/reservation/provisioning flow
- provisioning idempotency test
- credential encryption/redaction
- admin 2FA
- backup restore test
- provider pilot test
- support SOP readiness
- incident owner assignment
```

## 4. Security checklist

```text
[ ] Admin 2FA enforced
[ ] Reseller owner 2FA available/recommended
[ ] Password reset secure
[ ] Login rate limit enabled
[ ] Top-up/checkout/reveal credential rate limit enabled
[ ] Provider API keys not hardcoded
[ ] Secrets stored outside repo
[ ] Credential encrypted at rest
[ ] Credential reveal audited
[ ] Audit redacts sensitive fields
[ ] Cross-tenant tests pass
[ ] Emergency admin access requires reason
[ ] Session/token expiration configured
```

## 5. Billing checklist

```text
[ ] Wallet ledger append-only
[ ] Wallet balance equals ledger sum
[ ] Approved top-up creates ledger credit
[ ] Rejected top-up creates no credit
[ ] Client checkout debits client wallet
[ ] Reseller settlement debits reseller wallet
[ ] Insufficient reseller balance blocks provisioning
[ ] Refund creates reversal/credit entry
[ ] Adjustment requires reason
[ ] Daily reconciliation report works
[ ] Duplicate payment reference blocked/warned
```

## 6. Catalog/order checklist

```text
[ ] Product/plan/source active status works
[ ] Tenant clone pricing works
[ ] Margin floor works if enabled
[ ] Order stores price/cost/policy snapshot
[ ] Disabled plan cannot checkout
[ ] Disabled source cannot checkout
[ ] Out-of-stock returns correct error
[ ] Reservation TTL works
[ ] Concurrent reservation test passes
```

## 7. Provisioning checklist

```text
[ ] Queue worker running
[ ] Job idempotency key unique
[ ] Provider adapter health check works
[ ] Provider success creates service
[ ] Provider out-of-stock follows refund/release policy
[ ] Provider timeout uncertain moves manual review
[ ] No blind retry for unsafe create timeout
[ ] Provider resource mapping saved
[ ] Manual review process ready
[ ] Failed provisioning alert works
```

## 8. Service lifecycle checklist

```text
[ ] Service activation works
[ ] Renew uses correct term calculation
[ ] Expiry job works
[ ] Grace period works
[ ] Suspension works
[ ] Termination works
[ ] Lifecycle event history works
[ ] Capability-based action buttons work
```

## 9. Portal checklist

```text
[ ] Admin can manage tenant/catalog/wallet/orders/services
[ ] Reseller can manage clients/catalog/wallet/services within tenant
[ ] Client can top up/order/view/renew service
[ ] UI hides unauthorized actions
[ ] API blocks unauthorized actions
[ ] Credentials masked by default
[ ] Empty/error/loading states acceptable
```

## 10. Operations checklist

```text
[ ] Support macros ready
[ ] Finance reconciliation SOP ready
[ ] Abuse SOP ready
[ ] Provider onboarding checklist complete
[ ] Incident playbook ready
[ ] Admin alert channel ready
[ ] Monitoring dashboard ready enough
[ ] Error logs accessible
[ ] Backup schedule configured
[ ] Restore drill completed
```

## 11. Beta limits

Recommended initial pilot:

```text
2–3 resellers
10–30 clients
1 VPS product
1 proxy product
max 3 active services per client
daily finance reconciliation
manual review for high-risk orders
no postpaid
```

## 12. Go meeting agenda

```text
1. Review P0 QA status.
2. Review open bugs.
3. Review finance reconciliation.
4. Review provider readiness.
5. Review support/incident readiness.
6. Confirm pilot limits.
7. Assign launch-day owners.
8. Decide Go / Conditional Go / No-Go.
```

## 13. Decision types

```text
GO: all P0 pass; P1 issues acceptable.
CONDITIONAL GO: minor P1 issues with mitigation and owner.
NO-GO: any P0 fail or unclear money/tenant/provisioning/security behavior.
```

## 14. Launch day monitoring

```text
login errors
checkout errors
top-up approvals
wallet mismatch
provisioning queue latency
provider error rate
failed jobs
manual review queue
support tickets
abuse reports
```

## 15. Pause triggers

```text
ledger mismatch
cross-tenant access
provider duplicate
severe credential issue
provisioning failure rate above threshold
support volume exceeds capacity
provider account/payment issue
```

## 16. Sign-off

```text
Product Owner:
Engineering Lead:
QA Lead:
Ops Lead:
Finance Lead:
Security Owner:
Support Owner:
Launch Date:
Decision:
```

## 17. Closing principle

```text
Launch is not opening the door. Launch is knowing which doors open, which stay locked, and who holds each key.
```
