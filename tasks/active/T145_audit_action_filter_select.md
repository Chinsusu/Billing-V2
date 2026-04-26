# T145 - Audit action filter select

Status: REVIEW
Owner: Codex
Branch: codex/t145-audit-action-filter-select
PR: https://github.com/Chinsusu/Billing-V2/pull/321
Risk: frontend
Created: 2026-04-26
Updated: 2026-04-26

## Summary

Replace the audit action free-text filter with a plain dropdown and show readable action labels.

## Scope

- Add audit action options with human labels.
- Replace the raw action text filter with a select menu.
- Keep API query values unchanged.
- Display known audit actions with readable labels in the audit table.
- Update admin browser smoke to verify the selected action value is sent.

## Acceptance Criteria

- Audit action filter uses a select menu instead of free text.
- Users see readable action labels rather than raw backend action keys.
- API requests still send the expected `action` value.
- Frontend lint, sensitive-text check, smoke, build, taskguard, and diff check pass.

## Notes

- This is frontend-only and does not change audit API filters.

## Agent Log

- 2026-04-26: Codex created and claimed the task after T144 was marked done; starting audit action filter cleanup.
- 2026-04-26: Added audit action dropdown options, readable action labels, and smoke coverage for the submitted action value.
- 2026-04-26: Validation passed: frontend lint, sensitive-text check, build, admin browser smoke, taskguard, and diff check.
- 2026-04-26: Opened PR https://github.com/Chinsusu/Billing-V2/pull/321 for review.
