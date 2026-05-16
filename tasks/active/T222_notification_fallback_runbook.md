# T222 - Notification fallback runbook

Status: REVIEW
Owner: Codex
Branch: codex/t222-notification-fallback-runbook
PR: https://github.com/Chinsusu/Billing-V2/pull/475
Risk: launch readiness, notification delivery, support operations, security
Created: 2026-05-16
Updated: 2026-05-16

## Summary

Define the production notification delivery or manual fallback evidence required before launch, including launch-critical events, owner/SLA fields, safe message rules, and redacted evidence packet.

## Scope

- Add a notification delivery/manual fallback runbook under launch operations docs.
- Update launch evidence docs to reference the fallback packet.
- Update documentation indexes for the new runbook.
- Do not add SMTP/Telegram credentials or production delivery code.
- Do not claim production delivery proof or named owner sign-off.

## Acceptance Criteria

- Runbook identifies launch-critical notification events and manual fallback priorities.
- Runbook defines required owner/SLA fields and explicit pass/fail criteria.
- Runbook forbids secrets, credentials, DSNs, raw provider payloads, and private abuse evidence in fallback messages/evidence.
- Launch evidence docs still state notification delivery/fallback is not proven until owners/sign-off exist.
- Task board validation and docs checks pass.

## Notes

- This reduces the notification/fallback blocker to an executable evidence packet, but does not clear the blocker without named owners and a real production delivery or fallback drill.

## Agent Log

- 2026-05-16: Task created and claimed by Codex from latest `origin/main`.
- 2026-05-16: Added notification delivery/manual fallback runbook, launch evidence references, and documentation index entries. Validation passed: `go run ./cmd/taskguard`; `git diff --check`; changed-file secret pattern scan returned no matches.
- 2026-05-16: Opened PR #475 for review.
