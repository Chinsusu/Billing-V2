# T167 - Demo alert labels

Status: IN_PROGRESS
Owner: Codex
Branch: codex/t167-demo-alert-labels
PR: -
Risk: frontend fallback display
Created: 2026-04-26
Updated: 2026-04-26

## Summary

Clean up demo alert copy so operators do not see raw status keys in alert banners and the alert center.

## Scope

- Replace raw alert text such as `manual_review` with readable words.
- Keep demo alert records aligned between dashboard banners and the Alerts screen.
- Add browser smoke coverage for the Alerts screen.

## Acceptance Criteria

- Alerts show `manual review` instead of `manual_review`.
- Admin browser smoke covers the Alerts screen.
- Raw alert keys are not visible in the Alerts screen.

## Notes

- This is limited to static demo alert data; live alert APIs are not present yet.

## Agent Log

- 2026-04-26: Task created and claimed by Codex.
