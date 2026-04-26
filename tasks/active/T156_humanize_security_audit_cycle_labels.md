# T156 - Humanize security audit cycle labels

Status: REVIEW
Owner: Codex
Branch: codex/t156-humanize-security-audit-cycle-labels
PR: https://github.com/Chinsusu/Billing-V2/pull/343
Risk: frontend
Created: 2026-04-26
Updated: 2026-04-26

## Summary

Show readable security, audit action, actor, and billing cycle labels instead of raw keys.

## Scope

- Add shared display helpers for audit actions, audit actors, security status, and billing cycles.
- Use the helpers in admin account, audit log, report, recovery audit, and reseller catalog displays.
- Keep filter values, API values, and mock source values unchanged.
- Update browser smoke coverage where useful for readable labels.

## Acceptance Criteria

- UI shows labels such as Two-factor enabled, Retry job, Provider webhook, and 30 days instead of raw keys.
- Audit action select values still send raw action keys to the backend.
- Frontend lint, sensitive-text check, smoke, build, taskguard, and diff check pass.

## Notes

- This is frontend-only and should not change any API payloads.

## Agent Log

- 2026-04-26: Codex created and claimed the task after T155 was marked done; starting label cleanup for security, audit, and billing cycle text.
- 2026-04-26: Added shared display helpers for audit actions, account actors, security status, and billing cycles; applied them to admin logs, reports, recovery audit, account rows, and reseller catalog.
- 2026-04-26: Validation passed: frontend lint, sensitive-text check, production build, admin browser smoke, taskguard, and diff check.
- 2026-04-26: Opened review PR https://github.com/Chinsusu/Billing-V2/pull/343.
