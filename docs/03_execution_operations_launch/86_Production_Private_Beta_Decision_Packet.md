# 86 - Production and Private-Beta Decision Packet

**Date:** 2026-05-21
**Scope:** Current launch decision map for selected non-production pilot, private beta, production, provider, notification, and target-environment scopes.
**Decision:** GO remains limited to the selected bounded non-production pilot. Production, broader private beta, production customer data, broader provider scope, and production notification primary-path operation remain NO-GO until separately proven and signed off.

## Boundary

This packet consolidates current repo evidence into explicit scope decisions. It does not add new runtime evidence and does not approve any broader launch.

Do not broaden this packet:

- No production database, production provider account, production notification channel, or production customer data was used.
- No raw DB DSN, password, cookie, session token, provider token, provider payload, Telegram token, TOTP value, or plaintext service credential is recorded here.
- No money mutation route, provider mutation route, Cloudmini create/delete/action route, or notification delivery route was called by this task.
- Evidence links and outcomes are redacted summaries only.

## Source Evidence

Primary evidence used for this decision:

- `docs/03_execution_operations_launch/69_Pilot_Go_No_Go_Record.md`
- `docs/03_execution_operations_launch/70_Launch_Evidence_Completion_Packet.md`
- `docs/03_execution_operations_launch/83_UAT_Consolidated_Evidence_Packet.md`
- `docs/03_execution_operations_launch/85_Domain_Aware_Target_Auth_Smoke_Evidence.md`

Supporting evidence is referenced from those packets, including provider sandbox/error evidence, backup/restore evidence, full E2E evidence, notification evidence, target runtime evidence, finance reconciliation evidence, Admin 2FA evidence, and selected owner sign-off records.

## Decision Matrix

| Scope | Decision | Why |
| --- | --- | --- |
| Selected bounded non-production pilot | GO | Docs 69 and 70 record owner-accepted selected-pilot GO with target evidence, one-resource Cloudmini limit, manual fallback, support-window closeout, finance reconciliation, protected runtime, and domain-aware auth/RBAC evidence. |
| Continued selected test-server validation | GO with controls | T291 proves the selected test server can run the deployed domain-aware auth/RBAC smoke. Continued validation must stay non-production, redacted, and within existing provider/notification guardrails. |
| Broader private beta | NO-GO | Current evidence does not approve a broader customer list, broader support scope, production customer data, higher provider quota, or production notification primary path. |
| Production launch | NO-GO | Current evidence is non-production and staging-equivalent. Production host, secret-store, backup/restore, E2E, notification, owner sign-off, and finance/security proof must be repeated or explicitly approved for production scope. |
| Production customer data | NO-GO | No production data classification, backup/restore proof, access control proof, or customer-data handling sign-off exists for production data. |
| Broader provider or production-like real provisioning | NO-GO | Cloudmini approval is limited to selected bounded non-production scope and one active test resource. Higher quota, additional accounts, or production-like provisioning require a new owner-approved packet. |
| Telegram as sole production primary notification path | NO-GO | Selected-host Telegram preflight, one queued delivery, and failure classification passed, but broader production primary-path approval is still scope-specific and missing. |
| Manual notification fallback for selected pilot | GO | T244/T264/T270-T276 accept manual fallback for the historical selected-pilot window and record support coverage/closeout. |

## Required Evidence To Broaden Scope

Any broader private-beta or production GO requires a new packet or an update to this packet with actual evidence for the requested scope. Minimum required items:

- Scope approval: named launch scope, customer list or customer-data classification, communication plan, rollback rule, and owner acceptance.
- Environment proof: target host/domain, service process ownership, protected secret files, process argv secret-pattern check, health/readiness, and cloudflared or ingress proof for that exact target.
- Secret-store proof: production or launch-scope credential path, file mode/owner, rotation state, and no committed secret values.
- Backup/restore proof: clean shared staging or production-equivalent restore drill, or signed staging-equivalent exception that explicitly applies to the requested scope.
- Full E2E proof: launch-scope full E2E including checkout, wallet debit, provisioning path, renewal, credential reveal, audit, and finance reconciliation where applicable.
- Provider proof: approved account, credential path, SKU/source/group mapping, quota/spend/concurrency guardrails, timeout/idempotency behavior, cleanup owner/procedure, and pilot run evidence for the requested provider scope.
- Notification proof: owner-approved primary notification path or manual fallback for the requested launch window, with SLA, escalation, delivery or drill evidence, and redacted failure/retry behavior.
- Security proof: target Admin 2FA enrollment/enforcement, credential reveal no-store/audit/redaction, RBAC/tenant negative checks, and incident owner sign-off.
- Finance proof: wallet/ledger/payment reconciliation reviewed by the Finance owner for the requested launch window and target environment.
- Support/Ops proof: named coverage window, escalation channel, support owner, ops owner, and pause criteria for health/readiness, provider, ledger, credential, or notification failures.

## Non-Negotiable Pause Criteria

Pause or keep NO-GO if any of these occur:

- Ledger mismatch, duplicate debit, duplicate provider resource, or unbalanced finance reconciliation.
- Cross-tenant access, RBAC bypass, missing 2FA enforcement, or session/cookie leakage.
- Secret exposure in git, logs, process argv, docs, PR text, shell history, or evidence output.
- Provider timeout, partial success, unknown status, or cleanup uncertainty outside an approved manual-review path.
- Notification primary path failure without an approved fallback owner and SLA for the requested scope.
- Missing named owner for Product, Engineering, QA, Ops, Finance, Security, Support, or Provider decisions.

## Decision

```text
Selected bounded non-production pilot: GO
Continued selected non-production validation: GO with existing controls
Broader private beta: NO-GO
Production launch: NO-GO
Production customer data: NO-GO
Broader provider or production-like real provisioning: NO-GO
Telegram as sole production primary notification path: NO-GO
Manual notification fallback for selected pilot: GO
Open P0 bugs for selected pilot evidence: 0
Open P1 bugs for selected pilot evidence: 0
```

This packet is intentionally conservative. To change any NO-GO row to GO, create a new scoped task, collect the missing evidence, run the required validation, obtain owner sign-off for that exact scope, and update the decision record through PR.

Use `docs/03_execution_operations_launch/87_Scope_Intake_And_Preflight_Runbook.md` as the required intake and preflight procedure before changing any broader production, private-beta, provider, notification, customer-data, or target-environment row from `NO-GO` to `GO`.

The first broader private beta v1 intake is tracked in `docs/03_execution_operations_launch/88_Broader_Private_Beta_V1_Intake_Packet.md`; it remains `NO-GO` until that packet is completed with scope-specific approvals and evidence.
