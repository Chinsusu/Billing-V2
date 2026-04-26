# T154 - Move money format helpers out of mocks

Status: REVIEW
Owner: Codex
Branch: codex/t154-money-format-source
PR: https://github.com/Chinsusu/Billing-V2/pull/339
Risk: frontend
Created: 2026-04-26
Updated: 2026-04-26

## Summary

Make money formatting helpers production utilities instead of mock data exports.

## Scope

- Move `fmtMoney` and `fmtMoneyShort` into the shared frontend format helper module.
- Update frontend screens to import money format helpers from the production helper.
- Remove unused money format exports from mock sample data.
- Keep visible money formatting unchanged.

## Acceptance Criteria

- Production screens no longer import `fmtMoney` or `fmtMoneyShort` from mock data.
- Money values keep the same display format for demo and live fallback rows.
- Frontend lint, sensitive-text check, smoke, build, taskguard, and diff check pass.

## Notes

- This is frontend-only and does not change mock records or API payloads.

## Agent Log

- 2026-04-26: Codex created and claimed the task after T153 was marked done; starting money format helper cleanup.
- 2026-04-26: Moved money format helpers into the production format module and updated frontend screen imports away from mock sample data.
- 2026-04-26: Validation passed: frontend lint, sensitive-text check, production build, admin browser smoke, taskguard, and diff check.
- 2026-04-26: Opened review PR https://github.com/Chinsusu/Billing-V2/pull/339.
