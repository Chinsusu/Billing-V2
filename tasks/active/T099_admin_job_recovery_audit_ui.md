# T099 - Admin job recovery audit UI

Status: IN_PROGRESS
Owner: Codex
Branch: codex/t099-admin-job-recovery-audit-ui
PR: -
Risk: frontend/admin-ops
Created: 2026-04-24
Updated: 2026-04-24

## Summary

Add recovery action audit context to the admin job detail panel so operators can see who retried, cancelled, or moved a job to manual review.

## Scope

- Work mainly in frontend admin fulfillment screens and shared audit API helpers.
- Reuse existing admin audit log API and T095 job detail panel.
- Filter audit logs by job target when a live job is selected.
- Show actor type, action, display ID, created time, and correlation ID.
- Keep recovery controls admin-only.
- Keep each frontend file under 500 lines.

## Acceptance Criteria

- Admin job detail panel shows a compact audit trail for recovery actions.
- Loading, empty, and API error states are handled cleanly.
- Audit entries do not expose provider credentials or raw provider payloads.
- Reseller/client views do not receive admin-only audit/recovery controls.
- Frontend and backend validation commands pass.

## Notes

- Should follow T091 and T095.
- If backend audit filters are insufficient, keep the task frontend-only and document any follow-up.

## Agent Log

- 2026-04-24: Task created after T096 completed and the active board was fully DONE.
- 2026-04-24: Codex claimed the task on `codex/t099-admin-job-recovery-audit-ui`.
- 2026-04-24: Added the admin job recovery audit panel, wired job-target audit filtering, and kept existing audit screen filters typed separately from pageable API query params.
- 2026-04-24: Validation passed: `npm ci`, `npm audit --omit=dev`, `npm run lint`, `npm run build`, `go test ./...`, `go build ./cmd/api ./cmd/migrate ./cmd/seed ./cmd/smoke ./cmd/worker`, plus mocked Playwright desktop/mobile checks for the recovery audit panel.
