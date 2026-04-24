# T069 - Admin top-up review actions

Status: TODO
Owner: -
Branch: feat/admin-topup-review-actions
PR: -
Risk: frontend/API
Created: 2026-04-24
Updated: 2026-04-24

## Summary

Turn the admin top-up review queue from read-only display into real approve/reject actions backed by the existing backend review endpoints.

## Scope

- Add frontend API client wrappers for:
  - `POST /admin/topup-requests/{id}/approve`
  - `POST /admin/topup-requests/{id}/reject`
- Update `AdminTopups.tsx` to submit approve/reject actions only for reviewable live rows.
- Require or capture a rejection note before reject when the UI offers that action.
- Refresh the top-up list after a successful review action.
- Keep demo fallback read-only; do not show fake write buttons for demo rows.

## Acceptance Criteria

- Approve and reject buttons call real admin top-up review endpoints for live records.
- The screen shows pending, success, and error feedback for each action.
- Buttons are disabled while an action is running and hidden or disabled for non-reviewable statuses.
- Demo fallback remains clearly read-only.
- `npm run lint` and `npm run build` pass in `frontend`.

## Notes

- Backend reviewable statuses are `submitted` and `under_review`.
- Use the UUID `id` for API calls but keep display IDs visible to users.

## Agent Log

- 2026-04-24: Task created after T066 completed and the admin frontend batch was refreshed.
