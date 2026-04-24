# T095 - Admin job detail timeline UI

Status: TODO
Owner: -
Branch: codex/t095-admin-job-detail-timeline-ui
PR: -
Risk: frontend/admin-ops
Created: 2026-04-24
Updated: 2026-04-24

## Summary

Give admins a focused view of a provisioning job, its attempts, latest error, and recovery context.

## Scope

- Work mainly in frontend admin fulfillment screens and shared job API helpers.
- Reuse the T088 attempts API and T090 job display patterns.
- Show display IDs first, with UUIDs only where needed for copy/debug.
- Keep the view compact and operational, not a marketing page.
- Keep each frontend file under 500 lines.

## Acceptance Criteria

- Admin can open a job detail/timeline view from the provisioning queue.
- The view shows job status, attempt count, next attempt time, manual review reason, and attempt history.
- Empty, loading, and API error states are handled cleanly.
- Reseller/client views are not given admin-only recovery detail controls.
- Frontend and backend validation commands pass.

## Notes

- Should follow T088, T090, and T093 if T093 is already merged.

## Agent Log

- 2026-04-24: Task created after T092 completed and the active board was empty.
