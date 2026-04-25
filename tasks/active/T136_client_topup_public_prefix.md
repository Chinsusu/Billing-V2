# T136 - Client top-up public prefix

Status: DONE
Owner: Codex
Branch: codex/t136-client-topup-prefix
PR: https://github.com/Chinsusu/Billing-V2/pull/303
Risk: frontend
Created: 2026-04-26
Updated: 2026-04-26

## Summary

Align client wallet top-up labels with the shared `TUP-` public ID prefix used by admin, reseller, audit, mock data, and wallet ledger references.

## Scope

- Replace client wallet `TOP-` top-up labels with `TUP-`.
- Keep backend IDs only in API action bodies.
- Leave unrelated top-up copy and mocks unchanged.

## Acceptance Criteria

- Client top-up success notices use `TUP-`.
- Client top-up request rows use `TUP-`.
- Frontend lint, sensitive-text check, build, taskguard, and diff check pass.

## Notes

- `TUP-` is the existing public prefix for top-up requests across the rest of the UI.

## Agent Log

- 2026-04-26: Codex created and claimed the task after T135 was marked done; starting client top-up prefix cleanup.
- 2026-04-26: Replaced client wallet top-up notice and row labels from `TOP-` to `TUP-`.
- 2026-04-26: Local validation passed: `npm --prefix frontend run lint`, `npm --prefix frontend run check:sensitive-text`, `npm --prefix frontend run build`, `go run ./cmd/taskguard`, and `git diff --check`.
- 2026-04-26: Opened PR https://github.com/Chinsusu/Billing-V2/pull/303 for review.
- 2026-04-26: PR https://github.com/Chinsusu/Billing-V2/pull/303 merged into `main`; marking task done.
