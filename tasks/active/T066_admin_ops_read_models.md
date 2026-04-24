# T066 - Admin ops read models

Status: DONE
Owner: Codex
Branch: codex/t066-admin-ops-read-models
PR: https://github.com/Chinsusu/Billing-V2/pull/151
Risk: frontend/admin
Created: 2026-04-24
Updated: 2026-04-24

## Summary

Wire the admin operational billing surfaces to live read models where backend endpoints already exist, starting with top-ups and overview-level ops data.

## Scope

- Work only inside `frontend/src/modules/admin/**/*` unless a small wiring fix is required after rebasing on T063.
- Replace mock-only top-up queue data in `AdminTopups.tsx` with live admin top-up reads.
- Refresh `AdminOverview.tsx` to consume live admin wallet/order/service/top-up summary data where practical.
- Preserve loading, empty, error, and fallback states so the admin portal still renders when the backend is offline.
- Do not add write actions that the backend does not already support in `main`.

## Acceptance Criteria

- `AdminTopups` reads from the shared frontend API client for `/admin/topup-requests`.
- `AdminOverview` surfaces at least one live-backed summary or recent-activity section using existing admin read endpoints.
- Buttons that still depend on missing write APIs are clearly non-destructive and do not pretend to complete the action.
- No admin screen silently falls back to unrelated mock data when live data is available.
- `npm run lint` and `npm run build` pass in `frontend`.

## Notes

- Rebase on top of T063 if that task lands first.
- Avoid editing `frontend/src/lib/api/**/*` in this task; consume shared helpers from `main`.
- Provisioning job actions still need backend endpoints, so keep that part read-only unless a real API already exists on `main`.

## Agent Log

- 2026-04-24: Task created for the next admin frontend batch after invoices, transactions, and logs filters were completed.
- 2026-04-24: Codex claimed the task and started wiring admin read models to live API data.
- 2026-04-24: Admin top-ups and overview now consume live read APIs with explicit demo fallback states. Local audit, lint, and build passed.
- 2026-04-24: PR #151 merged into `main` with commit `35d4957372a15419df5aa8361e25bcb90b1b73a2`.
