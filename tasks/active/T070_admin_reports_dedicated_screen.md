# T070 - Admin reports dedicated screen

Status: REVIEW
Owner: Codex
Branch: codex/t070-admin-reports-screen
PR: -
Risk: frontend/admin
Created: 2026-04-24
Updated: 2026-04-24

## Summary

Replace the admin reports route reuse of `AdminOverview` with a dedicated reports screen backed by existing admin billing and audit read models.

## Scope

- Add a dedicated `AdminReports.tsx` under `frontend/src/modules/admin/screens/`.
- Update `AdminPortal.tsx` so `admin-reports` renders the new screen.
- Use existing frontend API client methods for transactions, payment reconciliation, invoices, top-ups, and/or audit logs.
- Keep charts or summaries simple; prefer clear operational tables and KPI summaries over decorative visuals.
- Preserve loading, empty, error, and fallback states.

## Acceptance Criteria

- The admin reports route no longer renders `AdminOverview`.
- Reports include at least two live-backed sections from existing admin read endpoints.
- Live rows use numeric display IDs where available.
- Missing backend data is handled honestly with explicit fallback or empty states.
- `npm run lint` and `npm run build` pass in `frontend`.

## Notes

- This task should not add backend endpoints.
- Keep the screen under 500 lines; split helpers if needed.

## Agent Log

- 2026-04-24: Task created after T066 completed and the admin frontend batch was refreshed.
- 2026-04-24: Codex claimed the task and started building a dedicated live-backed admin reports screen.
- 2026-04-24: Added a dedicated AdminReports screen backed by reconciliation, transactions, invoices, and audit read APIs. Local audit, lint, and build passed.
