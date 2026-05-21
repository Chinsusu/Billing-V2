# 88 - Broader Private Beta V1 Intake Packet

**Date:** 2026-05-21
**Scope:** Intake packet for a possible broader private beta v1 decision.
**Decision recommendation:** `NO-GO` until the missing broader-scope approvals and evidence below are completed.

## Safety Boundary

This packet is an intake record only. It does not approve broader private beta, production launch, production customer data, provider expansion, or Telegram as a production primary path.

Do not broaden this packet:

- No production database, production provider account, production notification channel, production customer data, or customer list is recorded here.
- No raw DB DSN, password, cookie, session token, provider token, Telegram token, TOTP value, private key, provider payload, notification payload, proxy credential, or customer data is recorded here.
- No money mutation route, provider mutation route, Cloudmini create/delete/action route, or notification delivery route was called by this task.
- Existing selected-pilot evidence is referenced only as prior evidence. It is not reused as broader private beta approval.

## Requested Scope

```text
Packet ID: broader-private-beta-v1
Requested decision: NO-GO review
Requested launch type: broader private beta
Requested launch window: TBD
Customer list or data classification: TBD
Production data present: no evidence recorded in this packet
Target hosts and domains: billing.resvn.net, client.resvn.net, reseller.resvn.net are prior launch-candidate domains; exact broader-beta target approval is TBD
Backend/API base URL: TBD for the broader-beta launch window
Frontend base URLs: TBD for the broader-beta launch window
Provider account scope: selected non-production Cloudmini scope exists; broader-beta provider account scope is TBD
Provider quota/spend/concurrency limits: selected pilot was one bounded dev resource; broader-beta limits are TBD
Notification path: selected Telegram/manual-fallback evidence exists; broader-beta primary or fallback path is TBD
Rollback trigger: any pause criterion in doc 86, plus any broader-beta customer-impact trigger approved by owner
Rollback owner: TBD
Evidence storage reference: this packet plus scope-specific successor evidence docs
```

## Current Owner Fields

The project has used an `Admin` single-owner model for selected non-production pilot decisions. This packet does not automatically extend that approval to broader private beta v1 because broader scope changes customer exposure, support load, provider quota, notification responsibility, and finance/security risk.

| Role | Current value | Status for broader private beta v1 |
| --- | --- | --- |
| Product Owner | Admin candidate | Needs explicit scope approval. |
| Engineering Lead | Admin candidate | Needs explicit target/runtime approval. |
| QA Lead | Admin candidate | Needs explicit UAT and E2E acceptance. |
| Ops Lead | Admin candidate | Needs explicit service, tunnel, backup, restore, and rollback acceptance. |
| Finance Lead | Admin candidate | Needs explicit wallet/ledger/payment reconciliation sign-off for the broader window. |
| Security Owner | Admin candidate | Needs explicit data classification, secret-store, 2FA, RBAC, and credential-reveal acceptance. |
| Support Owner | Admin candidate | Needs explicit support window, SLA, escalation channel, and pause authority. |
| Provider Owner | Admin candidate | Needs explicit provider quota, cleanup, timeout/idempotency, and manual-review acceptance. |
| Single-owner risk accepted | TBD | Must be accepted again for this broader scope if `Admin` owns all roles. |
| Approval timestamp | TBD | Required before any `GO` or `CONDITIONAL GO`. |

## Evidence Reuse Boundary

| Evidence area | Existing evidence | Broader private beta v1 decision |
| --- | --- | --- |
| Selected pilot GO | Docs 69 and 70 record selected bounded non-production pilot GO. | Reuse as background only; not approval for broader customer scope. |
| Production/private-beta decision map | Doc 86 keeps broader private beta and production as `NO-GO`. | This packet keeps that decision unchanged. |
| Scope intake procedure | Doc 87 defines required intake fields and safe preflight procedure. | This packet applies doc 87 to broader private beta v1. |
| Domain-aware auth smoke | Doc 85 records target auth/RBAC smoke on known public domains. | Must be rerun or explicitly accepted for the broader-beta launch window and user set. |
| UAT evidence | Docs 79 to 83 record client, reseller, and admin UAT evidence. | Must be tied to the broader-beta customer list, data classification, and support window. |
| Provider evidence | Provider evidence exists for selected bounded Cloudmini scope. | Needs broader-beta account, quota, SKU/group mapping, timeout/idempotency, cleanup, and pilot evidence. |
| Notification evidence | Telegram/manual-fallback evidence exists for selected scope. | Needs broader-beta primary/fallback owner, SLA, escalation, and delivery or drill evidence. |
| Target preflight evidence | Doc 89 records read-only health/runtime/process-secret/secret-file metadata evidence for current launch-candidate domains. | Use as current target evidence only; rerun if target, launch window, domains, services, or secret paths change. |
| Auth/RBAC evidence | Doc 90 records domain-aware target auth/RBAC smoke evidence for the current launch-candidate domains. | Use as current auth/RBAC evidence only; rerun if target, launch window, user set, auth config, tenant mapping, or RBAC policy changes. |

## Required Preflight Before Review

Run the matching checks from doc 87 and store only redacted outcomes. A broader private beta review is not ready until each row is `PASS` or has an explicit owner-approved exception.

| Area | Required broader-beta proof | Current status |
| --- | --- | --- |
| Task board | No conflicting launch tasks and taskguard passes. | Pending for the final review window. |
| Repo state | Diff reviewed, no whitespace errors, no committed secrets/customer data. | Pending for the final review window. |
| Target health | API, frontend, and public domains return expected health/readiness. | Current launch-candidate target passed in doc 89; rerun if scope target changes. |
| Runtime services | API, frontend, worker, tunnel, and database services are active on the approved target. | Current launch-candidate target passed applicable API, frontend, tunnel, and database checks in doc 89; worker-specific proof remains tied to E2E/notification/provider runs. |
| Process secrecy | Process argv has no DSN/token/password/credential patterns. | Current launch-candidate target strict secret-value scan passed in doc 89. |
| Secret files | Launch-scope secret paths have restrictive owner/mode metadata and no committed secret values. | Current launch-candidate target metadata passed after remediation in doc 89. |
| Backup/restore | Clean shared staging or production-equivalent restore proof applies to this scope. | Pending explicit applicability or rerun. |
| Full E2E | Checkout, wallet debit, provisioning, renewal, credential reveal, audit, and finance reconciliation pass for the launch scope. | Pending broader-beta E2E run. |
| Auth/RBAC | Client/admin/reseller domain auth, 2FA gate, tenant mismatch, and RBAC denials pass for broader-beta users. | Current client/admin domain-aware auth/RBAC smoke passed in doc 90; reseller-domain login and broader-beta user-set proof remain tied to UAT. |
| Credential reveal | No-store reveal response, audit, redaction, and rate-limit proof pass on approved target. | Pending broader-beta target approval and rerun. |
| Finance | Wallet, ledger, payment, and reconciliation report is balanced and accepted by Finance owner. | Pending broader-beta reconciliation sign-off. |
| Provider | Mapping, quota, timeout/idempotency, cleanup, and error taxonomy proof match broader-beta provider scope. | Pending provider owner approval and evidence. |
| Notification | Primary path or manual fallback has owner, SLA, escalation, delivery/drill, and failure evidence. | Pending notification owner approval and evidence. |
| Support/Ops | Coverage window, escalation channel, rollback authority, and pause criteria are named. | Pending support and ops sign-off. |

## GO Criteria For This Scope

Do not change broader private beta v1 to `GO` unless all criteria below are complete:

- Exact customer list or data classification is approved without committing customer data.
- Exact launch window, target host/domain set, API base URL, and frontend base URLs are approved.
- Product, Engineering, QA, Ops, Finance, Security, Support, and Provider owners sign off for this scope.
- If one `Admin` owns all roles, concentration-of-duty risk is explicitly accepted for this scope and timestamped.
- Target-environment preflight passes for the approved broader-beta target.
- Full E2E, renewal, auth/RBAC, credential reveal, audit, and finance reconciliation pass for the approved broader-beta target.
- Provider quota/spend/concurrency limits and cleanup procedure are approved for the broader-beta customer exposure.
- Notification primary path or manual fallback is approved with owner, SLA, escalation, and drill or delivery evidence.
- Open P0 issues are zero, and any P1 issue has an owner-approved mitigation and pause rule.
- Pause criteria from doc 86 are reviewed and accepted before launch.

## Current Gaps

- Customer list or data classification is not approved for broader private beta v1.
- Launch window, target environment, API base URL, and frontend base URLs are not approved for broader private beta v1.
- Owner sign-off is not complete for Product, Engineering, QA, Ops, Finance, Security, Support, and Provider in this broader scope.
- Single-owner concentration risk is not accepted for this broader scope.
- Target-environment health, runtime, process secrecy, secret-file metadata, and ingress checks passed for the current launch-candidate target in doc 89, but target/window owner approval remains incomplete.
- Backup/restore applicability has not been signed for this scope.
- Full E2E, renewal, credential reveal, audit, and finance reconciliation have not been rerun for this scope; current client/admin auth/RBAC smoke passed in doc 90 but does not replace UAT.
- Provider account, quota/spend/concurrency limits, SKU/group mapping, cleanup, timeout/idempotency, and manual-review rules are not approved for this scope.
- Notification primary/fallback path, SLA, escalation, and delivery/drill evidence are not approved for this scope.

## Decision Recommendation

```text
Packet ID: broader-private-beta-v1
Requested decision: NO-GO review
Owner approvals complete: no
Target-environment proof complete: no
Customer/data classification complete: no
Provider broader-scope proof complete: no
Notification broader-scope proof complete: no
Finance/security/support sign-off complete: no
Pause criteria reviewed for this scope: no
Decision recommendation: NO-GO
```

The next safe action is to fill the missing owner-approved values in this packet, then run only the applicable read-only preflight checks from doc 87. Mutating provider, money, credential, or notification actions require explicit owner approval for this exact broader private beta v1 scope before execution.
