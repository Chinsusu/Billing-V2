# 34 — Beta Pilot Program & Feedback Loop

**Version:** v1.3  
**Date:** 2026-04-22  
**Scope:** Controlled pilot before public launch.

## 1. Purpose

Pilot is controlled risk discovery with real usage, not marketing theater.

```text
Pilot goal: find operational bugs while volume is small enough that money, trust, and provider relationships are protected.
```

## 2. Pilot goals

```text
- validate order-to-service flow
- validate wallet/ledger settlement
- validate reseller/client behavior
- validate provider reliability
- validate support load
- validate renewal/expiry flow
- validate abuse handling
```

## 3. Pilot participants

Recommended first batch:

```text
- 2 trusted resellers
- 10–30 clients total
- 1 VPS product
- 1 proxy product
- 1–2 provider sources
```

Participant criteria:

```text
trusted relationship
willing to report issues
not high-risk abuse profile
accepts pilot limits
understands manual review may occur
```

## 4. Pilot limits

```text
- max active services per client
- max orders per day
- max reseller wallet exposure
- max refund threshold
- manual review for high-value orders
- no postpaid
- no custom provider requests
```

Example:

```text
Client limit: 3 active services
Reseller limit: 30 active services
Daily order cap: 20 orders
High-value threshold: manual review above 100 USD
```

## 5. Pilot scope

Include:

```text
one basic VPS plan
one basic proxy plan/source
monthly billing
manual top-up
renew
suspend/expire
```

Exclude:

```text
complex coupons
custom enterprise deals
advanced provider actions
public reseller signup
high-risk proxy categories
```

## 6. Onboarding flow

```text
1. Invite reseller.
2. Create reseller tenant.
3. Configure storefront basics.
4. Clone allowed catalog.
5. Set reseller price.
6. Top up reseller wallet.
7. Create/allow pilot clients.
8. Walk through first top-up/order.
9. Monitor first service activation.
10. Collect feedback after first use.
```

## 7. Feedback channels

```text
dedicated support chat/channel
feedback form
weekly reseller call
internal bug board
daily ops review
```

Feedback form:

```text
Role:
Tenant:
Action attempted:
Expected result:
Actual result:
Screenshot/log:
Severity:
Can reproduce? yes/no
Business impact:
```

## 8. Pilot metrics

Flow metrics:

```text
registration success rate
top-up approval time
checkout success rate
provisioning success rate
average provisioning time
failed provisioning count
manual review count
renewal success rate
```

Finance metrics:

```text
ledger mismatch count
refund count
adjustment count
reseller low-balance events
paid-not-provisioned count
```

Support/risk metrics:

```text
tickets per 100 orders
top categories
first response time
resolution time
abuse reports
provider warnings
fraud flags
```

## 9. Success criteria

```text
- provisioning success rate >= 95%
- 0 unresolved ledger mismatch
- 0 cross-tenant incidents
- 0 duplicate provider resource incidents
- support tickets per order below target
- top-up approval flow stable
- renewal flow works end-to-end
- provider failure modes understood
```

## 10. Daily review

```text
1. Review orders in last 24h.
2. Review failed/pending provisioning.
3. Review wallet/ledger reconciliation.
4. Review support tickets.
5. Review provider health.
6. Review abuse/risk flags.
7. Decide fixes before next batch.
```

## 11. Expansion rule

Move to next batch only if:

```text
all blocker bugs fixed
no unresolved finance mismatch
no tenant/security issue
provider stable
support capacity available
```

## 12. Exit decisions

```text
Proceed: ready for wider launch.
Extend: core stable but needs more data.
Pause: serious unresolved P0/P1 risk.
Rollback: stop pilot and fix architecture/process.
```

## 13. Closing principle

```text
A good pilot is not one with no bugs. It is one where bugs appear small, early, visible, and survivable.
```
