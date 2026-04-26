# T147 - Admin status filter selects

Status: DONE
Owner: Codex
Branch: codex/t147-admin-status-filter-selects
PR: https://github.com/Chinsusu/Billing-V2/pull/325
Risk: frontend
Created: 2026-04-26
Updated: 2026-04-26

## Summary

Replace admin status free-text filters with dropdowns for screens that use fixed status values.

## Scope

- Add shared admin status filter option sets.
- Replace Status text inputs with select menus in admin customers, accounts, providers, provisioning, invoices, transactions, and top-ups.
- Keep API query values unchanged.
- Update admin browser smoke to verify selected status values are sent.

## Acceptance Criteria

- Admin status filters use select menus where status values are fixed.
- Users see readable status labels instead of placeholder lists of backend values.
- API requests still send the expected `status` values.
- Frontend lint, sensitive-text check, smoke, build, taskguard, and diff check pass.

## Notes

- This is frontend-only and does not change API filters.

## Agent Log

- 2026-04-26: Codex created and claimed the task after T146 was marked done; starting admin status filter cleanup.
- 2026-04-26: Added shared admin status filter options, replaced fixed status text filters with selects, and expanded smoke coverage for selected status query values.
- 2026-04-26: Validation passed: frontend lint, sensitive-text check, build, admin browser smoke, taskguard, and diff check.
- 2026-04-26: Opened PR https://github.com/Chinsusu/Billing-V2/pull/325 for review.
- 2026-04-26: PR https://github.com/Chinsusu/Billing-V2/pull/325 merged into main; marking task done.
