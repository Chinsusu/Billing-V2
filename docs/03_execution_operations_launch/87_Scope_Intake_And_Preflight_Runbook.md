# 87 - Scope Intake and Preflight Runbook

**Date:** 2026-05-21
**Scope:** Required intake and preflight procedure before changing any production, private-beta, provider, notification, customer-data, or target-environment `NO-GO` row to `GO`.
**Decision:** This runbook does not approve launch. It defines the minimum evidence package required before a broader scope decision can be reviewed.

## Safety Boundary

Use this runbook before any broadened launch decision.

Allowed:

- Recording redacted scope, owner, target, command, and result metadata.
- Running read-only or non-mutating preflight commands on an approved non-production or launch-candidate target.
- Linking to existing evidence docs by task/doc ID.

Forbidden:

- Committing raw DB DSNs, API keys, tokens, passwords, cookies, TOTP secrets/codes, private keys, provider payloads, proxy credentials, customer data, or notification payloads.
- Running production mutations, provider create/delete/action routes, money mutation routes, or external notification delivery without explicit owner approval for that exact scope.
- Reusing selected non-production pilot evidence as production approval without a documented scope exception.

## Intake Fields

Every broader GO request must define these fields before evidence collection:

```text
Scope ID:
Requested decision: GO / CONDITIONAL GO / NO-GO review
Requested launch type: broader private beta / production / provider expansion / notification primary path / customer-data use / target host change
Requested launch window:
Customer list or data classification:
Production data present: no / yes with masking or owner approval reference
Target hosts and domains:
Backend/API base URL:
Frontend base URLs:
Provider account scope:
Provider quota/spend/concurrency limits:
Notification path: manual fallback / Telegram / SMTP / other
Rollback trigger:
Rollback owner:
Evidence storage reference:
```

## Required Owner Approvals

Record named owners for the requested scope. `Admin` may own multiple roles only if the concentration-of-duty risk is explicitly accepted for this new scope.

```text
Product Owner:
Engineering Lead:
QA Lead:
Ops Lead:
Finance Lead:
Security Owner:
Support Owner:
Provider Owner:
Single-owner risk accepted: yes/no/not applicable
Approval timestamp:
```

## Preflight Checklist

Run only the checks that match the requested scope. Store redacted results in the evidence packet.

| Area | Required proof | Safe command or source |
| --- | --- | --- |
| Task board | No open in-flight/review launch tasks unless intentionally scoped. | `go run ./cmd/taskguard` |
| Repo state | Branch diff reviewed and no whitespace errors. | `git diff --check` |
| Target health | API, frontend, and public domains return expected health/readiness status. | `curl -fsS -o /dev/null -w '%{http_code}' <redacted-url>` |
| Runtime services | API/frontend/worker/tunnel/database services active for the target. | `systemctl is-active ...` |
| Process secrecy | Process argv does not contain DSN/token/password/credential patterns. | redacted `/proc/*/cmdline` pattern count only |
| Secret files | Secret directories/files have restrictive owner/mode metadata. | `stat -c '%a %U %G' <path>`; do not print file contents |
| Backup/restore | Clean shared staging or production-equivalent restore proof exists. | `docs/03_execution_operations_launch/67_Backup_Restore_Drill_Runbook.md` |
| Full E2E | Launch-scope checkout, provisioning, renewal, audit, and frontend flow pass. | `docs/03_execution_operations_launch/68_Full_E2E_Quality_Gate_Runbook.md` or scope-specific successor |
| Auth/RBAC | Client/admin domain auth, 2FA gate, tenant mismatch, and RBAC denials pass. | `dev-target-auth-rbac` with scope-approved base URLs |
| Credential reveal | No-store reveal response, audit, redaction, and rate-limit proof pass. | `dev-target-credential-reveal` on approved target only |
| Finance | Wallet/ledger/payment reconciliation is balanced for launch target. | `dev-target-finance-reconciliation` or Finance owner reviewed report |
| Provider | Mapping, quota, timeout/idempotency, cleanup, and error taxonomy proof match requested provider scope. | docs 66, 73, 77, or new provider packet |
| Notification | Primary path or manual fallback has owner, SLA, escalation, delivery/drill, and failure evidence. | docs 72, 78, or new notification packet |

## Evidence Packet Template

Store this as a new evidence doc or as an update to `docs/03_execution_operations_launch/86_Production_Private_Beta_Decision_Packet.md`.

```text
Packet ID:
Date/time UTC:
Operator:
Requested scope:
Requested decision:
Target environment:
Data classification:
Production data present:
Owner approvals complete: yes/no
Task board result:
Repo diff result:
Target health result:
Runtime service result:
Process argv secret-pattern result:
Secret-file metadata result:
Backup/restore evidence reference:
Full E2E evidence reference:
Auth/RBAC evidence reference:
Credential reveal evidence reference:
Finance reconciliation evidence reference:
Provider evidence reference:
Notification evidence reference:
Open P0 issues:
Open P1 issues:
Pause criteria reviewed: yes/no
Decision recommendation: GO / CONDITIONAL GO / NO-GO
Reviewer sign-off:
```

## Review Rule

Do not change a `NO-GO` row in doc 86 to `GO` until:

- every applicable intake field is filled;
- every required owner is named and has accepted the requested scope;
- every applicable preflight area has passing evidence or a documented owner-approved exception;
- no raw secret, customer data, cookie, provider payload, notification payload, or credential is committed;
- pause criteria are accepted before launch.

If any P0 evidence is missing or based on assumption, the requested scope remains `NO-GO`.
