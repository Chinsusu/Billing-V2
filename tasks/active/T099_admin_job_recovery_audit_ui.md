# T099 - Admin job recovery audit UI

Status: TODO
Owner: -
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
