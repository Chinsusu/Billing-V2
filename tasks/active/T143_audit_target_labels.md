# T143 - Audit target labels

Status: REVIEW
Owner: Codex
Branch: codex/t143-audit-target-labels
PR: https://github.com/Chinsusu/Billing-V2/pull/317
Risk: frontend
Created: 2026-04-26
Updated: 2026-04-26

## Summary

Use plain audit target labels in the frontend instead of raw backend target type names.

## Scope

- Map audit target type names to user-facing labels.
- Keep public target IDs and prefixes unchanged.
- Update browser smoke expectations for the new labels.

## Acceptance Criteria

- Live audit rows show labels such as `Job`, `Order`, `Service`, and `Provider`.
- Live audit rows do not expose raw labels such as `provider_source` or `service_instance`.
- Frontend lint, sensitive-text check, smoke, build, taskguard, and diff check pass.

## Notes

- This is display-only. API response fields and filters keep using existing target type values.

## Agent Log

- 2026-04-26: Codex created and claimed the task after T142 was marked done; starting audit target label cleanup.
- 2026-04-26: Mapped live audit target types to plain labels and updated admin smoke expectation from `job JOB-*` to `Job JOB-*`.
- 2026-04-26: Local validation passed: `npm --prefix frontend run lint`, `npm --prefix frontend run check:sensitive-text`, `npm --prefix frontend run build`, `npm --prefix frontend run smoke:admin:ci`, `go run ./cmd/taskguard`, and `git diff --check`.
- 2026-04-26: Opened PR https://github.com/Chinsusu/Billing-V2/pull/317 for review.
