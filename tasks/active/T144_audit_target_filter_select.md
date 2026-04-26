# T144 - Audit target filter select

Status: REVIEW
Owner: Codex
Branch: codex/t144-audit-target-filter-select
PR: https://github.com/Chinsusu/Billing-V2/pull/319
Risk: frontend
Created: 2026-04-26
Updated: 2026-04-26

## Summary

Replace the audit target type text filter with a plain dropdown.

## Scope

- Add a reusable admin filter select control.
- Replace the raw audit target type input with options such as Job, Order, Service, Provider, and Top-up.
- Keep API query values unchanged.
- Update admin browser smoke to verify the selected target type is sent.

## Acceptance Criteria

- Audit target filter uses a select menu instead of free text.
- Users see plain target labels rather than raw backend values.
- API requests still send the expected `target_type` value.
- Frontend lint, sensitive-text check, smoke, build, taskguard, and diff check pass.

## Notes

- This is frontend-only and does not change audit API filters.

## Agent Log

- 2026-04-26: Codex created and claimed the task after T143 was marked done; starting audit target filter cleanup.
- 2026-04-26: Added a reusable admin filter select, changed audit target type to human labels, and updated smoke coverage for the submitted API value.
- 2026-04-26: Validation passed: frontend lint, sensitive-text check, build, admin browser smoke, taskguard, and diff check.
- 2026-04-26: Opened PR https://github.com/Chinsusu/Billing-V2/pull/319 for review.
