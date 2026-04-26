# T146 - Humanize top-up status labels

Status: REVIEW
Owner: Codex
Branch: codex/t146-humanize-topup-status-labels
PR: https://github.com/Chinsusu/Billing-V2/pull/323
Risk: frontend
Created: 2026-04-26
Updated: 2026-04-26

## Summary

Show readable top-up status labels and replace the admin top-up status text filter with a dropdown.

## Scope

- Add shared status labels for top-up review states such as `under_review` and `submitted`.
- Replace the admin top-up status free-text filter with a select menu.
- Keep API query values unchanged.
- Update admin browser smoke to verify the selected status value is sent.

## Acceptance Criteria

- Top-up status badges show readable labels rather than raw backend values.
- Admin top-up status filter uses a select menu.
- API requests still send the expected `status` value.
- Frontend lint, sensitive-text check, smoke, build, taskguard, and diff check pass.

## Notes

- This is frontend-only and does not change top-up API filters.

## Agent Log

- 2026-04-26: Codex created and claimed the task after T145 was marked done; starting top-up status label cleanup.
- 2026-04-26: Added readable top-up status labels, changed admin top-up status filter to a dropdown, and updated smoke coverage for the submitted status value.
- 2026-04-26: Validation passed: frontend lint, sensitive-text check, build, admin browser smoke, taskguard, and diff check.
- 2026-04-26: Opened PR https://github.com/Chinsusu/Billing-V2/pull/323 for review.
