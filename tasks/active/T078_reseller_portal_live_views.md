# T078 - Reseller portal live views

Status: TODO
Owner: -
Branch: codex/t078-reseller-portal-live-views
PR: -
Risk: frontend/reseller
Created: 2026-04-24
Updated: 2026-04-24

## Summary

Wire reseller portal screens to live reseller APIs for customers, services, billing, wallet, and dashboard summaries.

## Scope

- Work mainly in `frontend/src/lib/api/**/*` and `frontend/src/modules/reseller/screens/**/*`.
- Add frontend API wrappers for the reseller read endpoints delivered by T075-T077.
- Keep demo fallback explicit when APIs are unavailable.
- Show numeric display IDs in reseller tables.

## Acceptance Criteria

- Reseller clients, services, billing, wallet, and dashboard screens use live data when available.
- Live rows are not mixed with demo rows after a successful API load.
- API errors are visible and do not look like successful mock data.
- Each touched file stays under 500 lines.
- `npm audit --omit=dev`, `npm run lint`, and `npm run build` pass in `frontend`.

## Notes

- This task depends on T075-T077.
- Keep forms and mutations out of scope unless the backend endpoint already exists and the change is small.

## Agent Log

- 2026-04-24: Task created after T074 completed and the board needed the next reseller/live workflow batch.
