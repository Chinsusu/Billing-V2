# 30 — Provider Onboarding & Scoring Checklist

**Version:** v1.3  
**Date:** 2026-04-22  
**Scope:** Evaluation and onboarding of VPS/proxy providers.

## 1. Purpose

Provider selection must be based on operational reliability, not just cheap price. A cheap provider with unstable API, false stock, vague refund policy, or harsh abuse handling can destroy real margin.

## 2. Onboarding stages

```text
1. Business evaluation
2. Technical evaluation
3. Risk/abuse evaluation
4. Sandbox or manual test
5. Limited pilot
6. Production enablement
7. Ongoing score review
```

## 3. Provider profile

```text
Provider name:
Product type: VPS / Proxy / Dedicated / Other
Website:
Account/contact:
Billing method:
Currency:
Minimum deposit:
Refund policy:
Abuse policy:
API documentation:
Sandbox available: yes/no
Production account ready: yes/no
```

## 4. Business checklist

```text
- margin acceptable
- price stable enough
- stock availability acceptable
- payment terms fit prepaid model
- refund policy clear
- reseller/business use allowed
- intended use cases not prohibited
- support response time acceptable
- reputation acceptable
```

## 5. Technical checklist

```text
- API auth documented
- rate limit documented
- provision endpoint
- status endpoint
- suspend/unsuspend/terminate endpoints if supported
- renew endpoint or clear external billing model
- credential retrieval
- stock check
- idempotency key support
- external request ID/resource ID
- webhook if applicable
- documented error codes
- timeout behavior understood
```

## 6. Capability matrix

| Capability | Required for MVP | Provider value |
|---|---:|---|
| checkHealth | Yes | |
| checkStock | Yes | |
| provision | Yes | |
| getStatus | Yes | |
| suspend | Preferred | |
| unsuspend | Preferred | |
| terminate | Yes | |
| renew | Depends | |
| resetPassword | Optional | |
| reinstall | Optional | |
| changeIp | Optional | |
| usage/bandwidth | Optional | |
| console | Optional | |
| reverseDNS | Optional | |

## 7. Risk checklist

```text
- resale allowed?
- abuse response deadline?
- immediate suspension risk?
- IP reputation quality?
- stock shortage frequency?
- price change frequency?
- API instability?
- account freeze risk?
- payment/chargeback sensitivity?
```

## 8. Scoring formula

Score each 1–5:

```text
API Reliability
Automation Coverage
Stock Accuracy
Margin Quality
Support Speed
Refund Clarity
Abuse Policy Clarity
Operational Stability
Documentation Quality
Sandbox/Testability
```

Decision:

| Score | Decision |
|---:|---|
| 42–50 | Auto-provision candidate |
| 34–41 | Enable with monitoring |
| 26–33 | Pilot/manual review only |
| < 26 | Not production-ready |

## 9. Required tests before production

```text
1. checkHealth success/fail
2. checkStock available/out_of_stock
3. provision success
4. provision invalid request
5. provision timeout before response
6. provision timeout after resource created if mockable
7. getStatus active
8. terminate success
9. credential retrieval
10. rate limit behavior
11. auth failed behavior
```

## 10. Pilot limits

New provider should start with limits:

```text
- max orders/day
- max active services
- max reseller exposure
- manual review for high-value orders
- auto-provision only for low-risk plans
```

## 11. Provider source statuses

```text
draft: configured but not sellable
testing: internal only
pilot: limited exposure
active: production available
degraded: checkout blocked/warned depending policy
disabled: not sellable
retired: no new orders, existing lifecycle only
```

## 12. Circuit breaker

Degrade/disable if:

```text
- health check fails repeatedly
- provisioning failure rate exceeds threshold
- timeout rate spikes
- duplicate resource incident occurs
- provider reports stock/payment/account issue
- abuse/takedown volume spikes
```

Suggested:

```text
3 consecutive health failures -> degraded
5 consecutive health failures -> disabled for new orders
failure rate > 10% in 1h -> manual review
uncertain timeout > 3 in 30m -> pause auto-provision
```

## 13. Approval before active

```text
Business owner approves margin.
Technical owner approves adapter test.
Ops owner approves support process.
Abuse owner approves AUP fit.
Finance approves payment/refund handling.
```

## 14. Retirement SOP

```text
1. Disable new orders.
2. Keep existing services visible.
3. Define renewal/migration policy.
4. Notify affected resellers.
5. Archive provider config.
6. Keep resource mapping for audit.
```

## 15. Closing principle

```text
Provider is not just inventory. Provider is a risk partner inside the customer experience.
```
