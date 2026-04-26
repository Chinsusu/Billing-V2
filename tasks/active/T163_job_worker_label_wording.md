# T163 - Job worker label wording

Status: IN_PROGRESS
Owner: Codex
Branch: codex/t163-job-worker-label-wording
PR: -
Risk: frontend display labels
Created: 2026-04-26
Updated: 2026-04-26

## Summary

Clean up provisioning job timeline worker labels so readable worker names are not prefixed twice.

## Scope

- Keep raw worker keys hidden from the admin UI.
- Show a short readable worker label in the attempt timeline.
- Add smoke coverage for the final visible wording.

## Acceptance Criteria

- Timeline attempts show `Worker A` instead of `Worker Worker A`.
- Raw `worker-a` is still not visible.
- Admin browser smoke test covers the worker label.

## Notes

- Follow-up to T162 job timeline label cleanup.

## Agent Log

- 2026-04-26: Task created and claimed by Codex.
