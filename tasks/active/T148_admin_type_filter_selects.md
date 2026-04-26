# T148 - Admin type filter selects

Status: DONE
Owner: Codex
Branch: codex/t148-admin-type-filter-selects
PR: https://github.com/Chinsusu/Billing-V2/pull/327
Risk: frontend
Created: 2026-04-26
Updated: 2026-04-26

## Summary

Replace fixed-value admin type/product/source text filters with dropdowns.

## Scope

- Add shared option sets for account type, tenant type, product type, and provider source type.
- Replace matching text inputs with select menus in admin customers, accounts, providers, and provider readiness.
- Keep API query values unchanged.
- Update admin browser smoke to verify selected type/source/product values are sent.

## Acceptance Criteria

- Fixed-value admin type/product/source filters use select menus.
- Users see readable option labels rather than placeholder lists of backend values.
- API requests still send the expected query values.
- Frontend lint, sensitive-text check, smoke, build, taskguard, and diff check pass.

## Notes

- This is frontend-only and does not change API filters.

## Agent Log

- 2026-04-26: Codex created and claimed the task after T147 was marked done; starting admin type/source/product filter cleanup.
- 2026-04-26: Added shared type/product/source option sets, replaced matching text filters with selects, and expanded provider smoke coverage for selected query values.
- 2026-04-26: Validation passed: frontend lint, sensitive-text check, build, admin browser smoke, taskguard, and diff check.
- 2026-04-26: Opened PR https://github.com/Chinsusu/Billing-V2/pull/327 for review.
- 2026-04-26: PR https://github.com/Chinsusu/Billing-V2/pull/327 merged into main; marking task done.
