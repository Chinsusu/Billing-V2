# 31 — Incident Response & Disaster Recovery Playbook

**Version:** v1.3  
**Date:** 2026-04-22  
**Scope:** Handling production incidents and recovery.

## 1. Purpose

This playbook reduces damage during incidents by keeping the team focused on containment, evidence, recovery, and postmortem — especially around money, tenant isolation, credentials, provider resources, and wallet ledger.

## 2. Severity

| Severity | Description | Examples | Response |
|---|---|---|---|
| SEV0 | Whole platform/security/money critical | cross-tenant leak, admin compromised, credential exposure | immediate war room |
| SEV1 | Core flow broken or broad impact | ledger mismatch, queue stopped, duplicate provisioning | urgent |
| SEV2 | Limited group/provider affected | provider degraded, top-up delays | same day |
| SEV3 | Minor issue | UI bug, delayed notification | normal |

## 3. Roles

```text
Incident Commander: coordinates and makes priority decisions.
Technical Lead: investigates technical root cause.
Ops Lead: provider/support/customer process.
Finance Lead: wallet/ledger/refund exposure.
Security Lead: credential/access/token risk.
Communications Owner: internal/external updates.
```

## 4. General flow

```text
1. Detect
2. Triage severity
3. Assign incident commander
4. Contain damage
5. Preserve evidence/logs
6. Communicate internally
7. Fix or mitigate
8. Verify recovery
9. Communicate externally if needed
10. Postmortem
```

## 5. What not to do

```text
- Do not delete logs.
- Do not edit ledger entries.
- Do not mass retry provisioning before understanding failure.
- Do not restore production DB casually.
- Do not send credentials through insecure channels.
- Do not re-enable provider auto-provision without verification.
```

## 6. Playbook — Provider API down

Symptoms:

```text
health checks fail
provisioning timeout increases
provider returns 5xx/429
pending orders rise
```

Immediate action:

```text
1. Mark provider degraded.
2. Stop new auto-provision for affected source if severe.
3. Move unsafe pending jobs to manual review.
4. Do not blindly retry create calls.
5. Notify reseller/admin if delay exceeds threshold.
```

Recovery:

```text
1. Confirm provider status.
2. Retry only safe jobs.
3. Reconcile uncertain jobs.
4. Reactivate source gradually.
5. Monitor failure rate.
```

## 7. Playbook — Duplicate provisioning suspected

Immediate action:

```text
1. Pause affected provider auto-provision.
2. Identify order_id, idempotency_key, external_resource_id.
3. Do not terminate duplicate until ownership is confirmed.
4. Create SEV1 or SEV0 if widespread.
```

Recovery:

```text
1. Map correct resource to service.
2. Terminate duplicate resource if safe.
3. Adjust provider cost if needed.
4. Fix idempotency root cause.
5. Add regression test.
```

## 8. Playbook — Ledger mismatch

Immediate action:

```text
1. Freeze affected wallet if severe.
2. Disable manual adjustment for non-admin roles if needed.
3. Export ledger and wallet state.
4. Identify first mismatch timestamp.
5. Open finance incident.
```

Recovery:

```text
1. Determine whether materialized balance or ledger flow is wrong.
2. Rebuild materialized balance from ledger if needed.
3. If correction is needed, create new adjustment/reversal entry.
4. Run full reconciliation.
5. Add regression test.
```

Never update/delete old ledger rows.

## 9. Playbook — Tenant data leak suspicion

Immediate action:

```text
1. Treat as SEV0.
2. Disable affected endpoint if needed.
3. Preserve logs.
4. Identify impacted tenants/resources.
5. Rotate exposed credentials if needed.
6. Notify security/leadership.
```

Recovery:

```text
1. Patch tenant guard.
2. Add regression tests.
3. Audit similar endpoints.
4. Notify affected parties according to policy/law.
5. Postmortem.
```

## 10. Playbook — Credential exposure

Immediate action:

```text
1. Stop leak path.
2. Rotate affected credentials if possible.
3. Revoke/rotate provider keys if exposed.
4. Preserve evidence.
5. Restrict access to leaked logs/artifacts.
```

Recovery:

```text
1. Add redaction tests.
2. Update credential handling SOP.
3. Notify impacted parties if required.
```

## 11. Playbook — Admin account compromised

```text
1. Disable account/session.
2. Revoke tokens.
3. Force password reset and 2FA review.
4. Review audit actions.
5. Freeze risky actions if needed.
6. Rotate secrets if accessed.
7. Revert unauthorized changes via new audited corrections.
```

## 12. Disaster recovery priorities

```text
1. Database and ledger integrity
2. Tenant isolation/security
3. Wallet/order/service consistency
4. Provider resource mapping
5. Credential access
6. Portal functionality
7. Reporting/dashboard
```

Do not restore portal while ledger/database consistency is unknown.

## 13. Backup/restore minimum

Backup must include:

```text
database, ledger, orders, services, provider mappings, encrypted credentials, audit logs, tenant/domain config
```

Restore drill must verify:

```text
app boots, login works, tenant isolation works, wallet balances reconcile, service mappings exist, credentials decrypt with correct key
```

## 14. Postmortem template

```text
Incident ID:
Severity:
Start/detect/resolve time:
Customer impact:
Financial impact:
Root cause:
Trigger:
Timeline:
What worked:
What failed:
Immediate fixes:
Long-term fixes:
Owners/due dates:
Regression tests added:
```

## 15. Closing principle

```text
An incident playbook looks boring on normal days. On bad days, it prevents the team from creating the second fire.
```
