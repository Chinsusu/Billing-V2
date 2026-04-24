# T095 - Admin job detail timeline UI

Status: DONE
Owner: Codex
Branch: codex/t095-admin-job-detail-timeline-ui
PR: https://github.com/Chinsusu/Billing-V2/pull/216
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
- 2026-04-24: Codex claimed the task on `codex/t095-admin-job-detail-timeline-ui`.
- 2026-04-24: Added admin job timeline panel, admin attempts API helper, and live queue selection from the provisioning table.
- 2026-04-24: Validation passed: `go test ./...`, `go build ./cmd/api ./cmd/migrate ./cmd/seed ./cmd/smoke ./cmd/worker`, `npm ci`, `npm audit --omit=dev`, `npm run lint`, `npm run build`.
- 2026-04-24: Opened PR https://github.com/Chinsusu/Billing-V2/pull/216 for review.
- 2026-04-24: PR #216 passed CI and merged into `main` at `f22b55818d5c376e75d3f49e02cd2878aef32344`.
