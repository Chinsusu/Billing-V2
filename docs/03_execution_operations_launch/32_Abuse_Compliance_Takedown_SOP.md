# 32 — Abuse, Compliance & Takedown SOP

**Version:** v1.3  
**Date:** 2026-04-22  
**Scope:** Manual abuse handling for VPS/proxy platform MVP.

## 1. Purpose

VPS/proxy products have high abuse risk. This SOP defines how to receive, classify, act on, and record abuse reports while protecting provider relationships, IP reputation, tenants, and platform continuity.

## 2. Abuse sources

```text
provider
datacenter/upstream network
email abuse desk
legal/law enforcement request
payment processor
internal monitoring
client/reseller report
third-party researcher
```

## 3. Abuse types

```text
spam
phishing
malware
botnet
brute_force
port_scanning
ddos
copyright
credential_theft
proxy_scraping_violation
payment_fraud
chargeback
illegal_content
aup_violation
provider_takedown
```

## 4. Severity

| Severity | Description | Action |
|---|---|---|
| Critical | legal/provider/security demands immediate action | suspend immediately |
| High | severe spam/phishing/malware/repeat abuse | suspend or short deadline |
| Medium | violation with remediation possible | warning + deadline |
| Low | unclear/minor first report | investigate |

## 5. Abuse case fields

```text
abuse_case_id
tenant_id
user_id/client_id
reseller_id
service_id
provider_source_id
external_resource_id
abuse_type
severity
report_source
report_received_at
deadline
evidence_summary
evidence_files_or_links
status
assigned_owner
actions_taken
final_resolution
```

## 6. Immediate suspension triggers

```text
- provider requires urgent takedown
- confirmed phishing/malware
- DDoS or active attack
- credential theft
- botnet/C2
- law enforcement request requiring action
- repeated abuse from same client/service
- payment fraud tied to service
```

## 7. Warning-first cases

Warning may be allowed when:

```text
- first-time low/medium violation
- evidence incomplete
- issue is likely misconfiguration
- provider deadline allows remediation
- no immediate network harm
```

Warning must include service, abuse type, required action, deadline, and consequence.

## 8. Workflow

```text
1. Receive report.
2. Create abuse case.
3. Validate service/provider mapping.
4. Classify severity.
5. Choose warn/suspend/terminate/escalate.
6. Notify reseller/client according to policy.
7. Record evidence and action.
8. Reply to provider if needed.
9. Monitor recurrence.
10. Close case with resolution.
```

## 9. Reseller responsibility

If client belongs to reseller:

```text
- notify reseller owner for medium/high cases
- platform can override reseller for severe cases
- repeated abuse can restrict reseller auto-provision
```

Possible reseller actions:

```text
warning
lower order limits
manual review for new orders
restrict risky products
suspend tenant
terminate relationship
```

## 10. Repeat offender policy

```text
1st low/medium: warning or temporary suspension
2nd: suspension + manual review for future orders
3rd: account restriction or termination
Critical: immediate suspension/termination
```

## 11. Evidence handling

Store timestamps, IP/domain/resource IDs, provider ticket IDs, report excerpts, screenshots/log snippets if needed, and action history.

Do not expose full evidence if it contains third-party private data, provider internals, security-sensitive indicators, or legal restrictions.

## 12. Abuse status lifecycle

```text
new
triaging
awaiting_client_action
suspended
resolved
rejected_false_positive
terminated
escalated
closed
```

## 13. False positive handling

```text
1. Mark rejected_false_positive.
2. Record evidence/reason.
3. Unsuspend if safe.
4. Notify affected user/reseller.
5. Reply to provider if needed.
```

## 14. Abuse and refund

Default:

```text
No automatic refund for confirmed severe abuse.
```

Exceptions may include false positive, provider error, platform fault, or mistaken suspension. Abuse-related refunds require approval and reference the abuse case.

## 15. Risk flags and blacklist

Track:

```text
email, domain, IP, payment reference, provider resource, client account, reseller tenant
```

Actions:

```text
manual review, purchase limit, product restriction, top-up review, account suspension
```

## 16. Compliance note

This SOP is operational, not legal advice. Legal requests should be escalated to the designated legal/compliance owner before disclosing user data unless emergency policy applies.

## 17. Metrics

```text
abuse cases per 100 services
abuse by provider
abuse by reseller
repeat offender rate
time to first action
provider takedown SLA hit rate
false positive rate
```

## 18. Closing principle

```text
Abuse is not an exception in VPS/proxy. It is an underground current that needs a channel, or it will break the source of supply.
```
