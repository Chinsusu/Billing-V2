# T185 - Overview demo signup label

Status: IN_PROGRESS
Owner: Codex
Branch: codex/t185-overview-demo-signup-label
PR: -
Risk: frontend demo labels and smoke coverage
Created: 2026-04-28
Updated: 2026-04-28

## Summary

Humanize the Overview demo signup activity that still exposes a raw email-style label.

## Scope

- Replace the visible `startup-dev-42@proton.me` signup activity text on the Overview fallback feed.
- Add Overview fallback smoke coverage to reject that raw email-style label.
- Do not change customer records, live API contracts, or backend behavior.

## Acceptance Criteria

- Overview demo fallback uses a readable customer signup label instead of the raw email-style label.
- Admin smoke verifies the readable signup label and rejects `startup-dev-42@proton.me` on Overview fallback.
- Frontend lint, sensitive-text check, production build, admin smoke, taskguard, and diff check pass locally.

## Notes

- This is a narrow follow-up to the Overview demo label cleanup.

## Agent Log

- 2026-04-28: Task created and claimed by Codex.
- 2026-04-28: Humanized Overview fallback signup activity and added smoke guard; local gates pass.
