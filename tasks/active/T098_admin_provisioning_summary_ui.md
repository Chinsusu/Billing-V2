# T098 - Admin provisioning summary UI

Status: TODO
Owner: -
Branch: codex/t098-admin-provisioning-summary-ui
PR: -
Risk: frontend/admin-ops
Created: 2026-04-24
Updated: 2026-04-24

## Summary

Show provisioning queue health at the top of the admin provisioning screen using the summary API from T097.

## Scope

- Work mainly in `frontend/src/modules/admin/**/*` and shared API helpers.
- Add compact KPI/status tiles for queued, running, retryable, manual-review, terminal, and cancelled jobs.
- Keep the existing queue table and job detail timeline usable.
- Handle loading, empty, and API error states.
- Keep reseller/client views free of admin-only recovery detail.
- Keep each frontend file under 500 lines.

## Acceptance Criteria

- Admin provisioning screen shows live summary data when API is configured.
- If summary API is unavailable, the UI degrades cleanly without hiding the queue.
- Attention states use restrained warning/error styling and do not overlap on mobile.
- Frontend and backend validation commands pass.

## Notes

- Depends on T097.
- Follow existing `KpiCard`, `StatusBadge`, and table styling patterns.

## Agent Log

- 2026-04-24: Task created after T096 completed and the active board was fully DONE.
