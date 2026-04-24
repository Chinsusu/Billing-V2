# T072 - Admin accounts live views

Status: REVIEW
Owner: Codex
Branch: codex/t072-admin-accounts-live-views
PR: https://github.com/Chinsusu/Billing-V2/pull/166
Risk: frontend/admin
Created: 2026-04-24
Updated: 2026-04-24

## Summary

Wire the admin tenants and customers screens to live account read APIs and remove mock-only behavior from the primary table views.

## Scope

- Work mainly in `frontend/src/lib/api/**/*` and `frontend/src/modules/admin/screens/AdminTenants.tsx` / `AdminCustomers.tsx`.
- Add frontend API types and wrappers for the endpoints delivered by T071.
- Show numeric display IDs in the tables, with UUIDs kept only for internal action keys.
- Keep loading, empty, error, and explicit demo fallback states.
- Do not add create/edit/delete account actions unless the matching backend endpoint already exists and the implementation stays small.

## Acceptance Criteria

- `AdminTenants` reads live tenant/account data from the shared API client.
- `AdminCustomers` reads live customer/account data from the shared API client.
- Live rows are not mixed with demo rows after data loads successfully.
- API failure displays a clear fallback state without hiding that data is not live.
- `npm audit --omit=dev`, `npm run lint`, and `npm run build` pass in `frontend`.

## Notes

- This task depends on T071 for backend endpoints.
- Keep table mapping helpers small and local; split helpers if either screen approaches 500 lines.

## Agent Log

- 2026-04-24: Task created after closing stale PR #80 and refreshing the board for the next live workflow batch.
- 2026-04-24: Codex claimed the task after T071 merged and started wiring admin account screens to live APIs.
- 2026-04-24: Added admin tenant/account frontend API wrappers and wired AdminTenants/AdminCustomers to live data with explicit demo fallback states. Local frontend audit, lint, and build passed.
- 2026-04-24: Opened PR #166 for review.
