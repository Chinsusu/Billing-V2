# T168 - Overview activity labels

Status: REVIEW
Owner: Codex
Branch: codex/t168-overview-activity-labels
PR: https://github.com/Chinsusu/Billing-V2/pull/367
Risk: frontend fallback display
Created: 2026-04-26
Updated: 2026-04-26

## Summary

Clean up demo Overview activity feed copy so the first admin screen does not show internal VPS names.

## Scope

- Replace internal demo service names in the Overview activity feed with readable business text.
- Add smoke coverage that keeps the internal VPS name off the Overview screen.

## Acceptance Criteria

- Overview activity feed shows readable provisioning activity.
- Raw internal label `vps-scrape-02` is not visible on the Overview screen.
- Admin browser smoke covers this guard.

## Notes

- This is limited to Overview demo feed copy; service inventory labels can be handled in a separate task.

## Agent Log

- 2026-04-26: Task created and claimed by Codex.
- 2026-04-26: Opened PR #367 after lint, sensitive-text guard, build, admin smoke, taskguard, and diff check passed.
