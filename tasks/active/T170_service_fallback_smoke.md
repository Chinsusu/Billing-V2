# T170 - Service fallback smoke coverage

Status: IN_PROGRESS
Owner: Codex
Branch: codex/t170-service-fallback-smoke
PR: -
Risk: frontend smoke coverage
Created: 2026-04-26
Updated: 2026-04-26

## Summary

Add browser smoke coverage for admin service inventory fallback labels so demo service data stays readable when the live service API is unavailable.

## Scope

- Cover the admin service inventory fallback path for VPS demo rows.
- Keep the guard focused on labels introduced by T169.
- Do not refactor the smoke script or split mock data in this task.

## Acceptance Criteria

- Admin browser smoke verifies the service inventory fallback message and readable VPS labels.
- Raw internal VPS labels remain forbidden by the smoke guard.
- Frontend lint, sensitive-text check, production build, admin smoke, taskguard, and diff check pass.

## Notes

- T169 added reseller/client fallback assertions. This task covers the remaining admin service fallback gap.

## Agent Log

- 2026-04-26: Task created and claimed by Codex.
