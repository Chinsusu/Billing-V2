# 29 — Customer Support SOP & Macro Templates

**Version:** v1.3  
**Date:** 2026-04-22  
**Scope:** Support operations for VPS/proxy platform.

## 1. Purpose

This document gives support a consistent way to handle tickets without inventing policy, exposing sensitive data, or making unsafe promises about provider provisioning.

Support must not:

```text
- edit wallet balances directly unless explicitly authorized
- reveal credentials through insecure channels
- blindly retry provisioning
- promise refund when provider state is uncertain
- disclose cross-tenant/internal provider data
```

## 2. Ticket categories

```text
Billing / Wallet
Top-up
Order / Checkout
Provisioning
Service access
Credential
Renewal / Expiry
Suspension / Termination
Provider issue
Abuse / Takedown
Account / Login
Reseller setup
Feature request
```

## 3. Priority levels

| Priority | Description | Suggested first response | Examples |
|---|---|---:|---|
| P0 | money/data/security risk | 15–30 min | ledger mismatch, tenant leak, credential exposure |
| P1 | core service unavailable | 1–2 h | active VPS inaccessible, queue stuck |
| P2 | delayed order/billing review | 4–8 h | provisioning pending, top-up mismatch |
| P3 | normal how-to | 24 h | renewal guide |
| P4 | feature request | 2–5 days | new provider request |

## 4. Escalation matrix

| Case | Escalate to |
|---|---|
| Ledger mismatch | Finance Lead + Engineering |
| Duplicate provisioning | Engineering + Ops |
| Provider outage | Ops + Provider Owner |
| Abuse takedown | Abuse/Compliance Owner |
| Credential exposure | Security Owner |
| Admin account compromise | Security Owner + Founder |
| Refund above threshold | Finance Lead/Super Admin |
| Reseller dispute | Account Manager/Founder |

## 5. Investigation checklist

Before replying:

```text
1. Identify tenant.
2. Identify actor role.
3. Identify order/service/top-up.
4. Check audit events.
5. Check ledger if billing.
6. Check provisioning job if pending.
7. Check provider health if service issue.
8. Check lifecycle status and expiry.
9. Check abuse/suspension reason.
10. Confirm what can be disclosed to this actor.
```

## 6. Macro — Top-up pending

```text
Hi {name},

Your top-up request is currently pending manual verification.

We are checking the payment reference and amount against the submitted request. Once verified, your wallet will be credited and you will receive a notification.

Reference: {topup_id}
Submitted amount: {amount}
Current status: Pending
```

## 7. Macro — Top-up approved

```text
Hi {name},

Your top-up has been approved and your wallet has been credited.

Amount credited: {amount}
Reference: {topup_id}

You can now use your wallet balance to place orders or renew services.
```

## 8. Macro — Top-up rejected

```text
Hi {name},

Your top-up request could not be approved.

Reason: {rejection_reason}

No wallet balance was changed for this rejected request. You may submit a new request with the correct payment reference.
```

## 9. Macro — Order pending provisioning

```text
Hi {name},

Your order has been received and is waiting for provisioning.

Order: {order_id}
Product: {product_name}
Status: {status}

Some providers require confirmation time. We do not blindly retry uncertain provider requests because that can create duplicate resources.
```

## 10. Macro — Provisioning failed/review

```text
Hi {name},

Your order could not be provisioned automatically.

Current action: {action_summary}

If the provider confirms no resource was created, the refund process can follow policy. If the provider state is uncertain, our team will verify first.
```

## 11. Macro — Service activated

```text
Hi {name},

Your service is now active.

Service: {service_name}
Expiry: {term_end}

Access details are available in your dashboard. Credentials are masked by default and reveal actions are logged for security.
```

## 12. Macro — Service expiring soon

```text
Hi {name},

Your service will expire soon.

Service: {service_name}
Expiry: {term_end}

Please renew before expiry to avoid interruption.
```

## 13. Macro — Service suspended for expiry

```text
Hi {name},

Your service has been suspended because it passed the expiry/grace period.

Service: {service_name}
Reason: Expired

You may renew if the product policy still allows reactivation.
```

## 14. Macro — Credential access

```text
Hi {name},

For security, credentials are visible only from the authorized service detail page and are masked by default.

Please log in, open the service detail page, and use the reveal action if your role has permission. We cannot send sensitive credentials through insecure support channels unless policy explicitly allows it.
```

## 15. Macro — Abuse warning

```text
Hi {name},

We received an abuse report related to your service.

Service: {service_name}
Type: {abuse_type}
Deadline: {deadline}

Please resolve this before the deadline. Severe or repeated abuse may result in suspension or termination.
```

## 16. Internal note template

```text
Tenant:
User:
Role:
Order:
Service:
Wallet:
Provider:
Current status:
Audit checked:
Ledger checked:
Provider job checked:
Risk:
Next action:
Owner:
Deadline:
```

## 17. Support metrics

```text
first response time
resolution time
tickets per 100 orders
billing ticket rate
provisioning ticket rate
refund rate
abuse ticket count
repeat issue count
```

## 18. Closing principle

```text
Good support protects the truth of the system. It does not create a bigger incident while trying to sound helpful.
```
