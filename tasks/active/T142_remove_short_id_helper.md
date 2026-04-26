# T142 - Remove short ID helper

Status: DONE
Owner: Codex
Branch: codex/t142-remove-short-id-helper
PR: https://github.com/Chinsusu/Billing-V2/pull/315
Risk: frontend
Created: 2026-04-26
Updated: 2026-04-26

## Summary

Remove the unused frontend helper that shortens internal IDs for display.

## Scope

- Delete the unused `shortID` helper from frontend formatting utilities.
- Confirm no frontend call sites still depend on shortened backend IDs.
- Keep public display ID helpers unchanged.

## Acceptance Criteria

- No frontend source call site references `shortID`.
- Frontend formatting still exposes `recordLabel` for public ID labels.
- Frontend lint, sensitive-text check, build, taskguard, and diff check pass.

## Notes

- This reduces the chance of future UI code displaying shortened UUIDs or correlation IDs instead of public IDs.

## Agent Log

- 2026-04-26: Codex created and claimed the task after T141 was marked done; removing unused internal ID display helper.
- 2026-04-26: Removed `shortID` from frontend format utilities and confirmed there are no frontend call sites left.
- 2026-04-26: Local validation passed: `npm --prefix frontend run lint`, `npm --prefix frontend run check:sensitive-text`, `npm --prefix frontend run build`, `go run ./cmd/taskguard`, and `git diff --check`.
- 2026-04-26: Opened PR https://github.com/Chinsusu/Billing-V2/pull/315 for review.
- 2026-04-26: PR https://github.com/Chinsusu/Billing-V2/pull/315 merged into `main`; marking task done.
