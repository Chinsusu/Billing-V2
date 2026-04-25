# T127 - Audit log public ID filters

Status: DONE
Owner: Codex
Branch: codex/t127-audit-public-id-filters
PR: https://github.com/Chinsusu/Billing-V2/pull/285
Risk: API/frontend
Created: 2026-04-25
Updated: 2026-04-25

## Summary

Allow admin audit logs to be filtered by actor public ID and target public ID instead of backend references.

## Scope

- Add API query support for `actor_display_id` and `target_display_id` on admin audit logs.
- Update frontend audit log filters to use public ID fields.
- Extend smoke coverage for the new filter path.
- Keep existing raw ID filters backward compatible if they already exist.

## Acceptance Criteria

- `/admin/audit-logs` can filter by actor public ID.
- `/admin/audit-logs` can filter by target public ID for supported target types.
- Admin audit UI sends numeric public ID filters and still hides raw backend IDs.
- Backend tests, frontend smoke, and taskguard pass.

## Notes

- T124 exposed the related display IDs; this task makes them usable for operator search.

## Agent Log

- 2026-04-25: Codex created and claimed the task after T126 merged; starting audit public ID filter support.
- 2026-04-25: Added backend `actor_display_id` and `target_display_id` filters, frontend audit filter fields, and smoke coverage.
- 2026-04-25: Validation passed: audit/smoke tests, full Go tests, Go build, frontend smoke/lint/build/sensitive-text, taskguard, and diff check.
- 2026-04-25: Opened PR #285 for review.
- 2026-04-25: PR https://github.com/Chinsusu/Billing-V2/pull/285 merged into `main`.
