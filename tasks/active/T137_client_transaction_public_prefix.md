# T137 - Client transaction public prefix

Status: DONE
Owner: Codex
Branch: codex/t137-client-transaction-prefix
PR: https://github.com/Chinsusu/Billing-V2/pull/305
Risk: frontend
Created: 2026-04-26
Updated: 2026-04-26

## Summary

Align client transaction labels with the shared `TX-` public ID prefix used by admin, reseller, reports, and wallet ledger references.

## Scope

- Replace client `TXN-` transaction labels with `TX-`.
- Keep invoice, order, wallet, and top-up prefixes unchanged.
- Leave backend IDs only in API action bodies or internal joins.

## Acceptance Criteria

- Client transaction list uses `TX-`.
- Client checkout payment summary uses `TX-`.
- Frontend lint, sensitive-text check, build, taskguard, and diff check pass.

## Notes

- `TX-` is the existing public prefix for payment transactions across the rest of the UI.

## Agent Log

- 2026-04-26: Codex created and claimed the task after T136 was marked done; starting client transaction prefix cleanup.
- 2026-04-26: Replaced client checkout and transaction list labels from `TXN-` to `TX-`.
- 2026-04-26: Local validation passed: `npm --prefix frontend run lint`, `npm --prefix frontend run check:sensitive-text`, `npm --prefix frontend run build`, `go run ./cmd/taskguard`, and `git diff --check`.
- 2026-04-26: Opened PR https://github.com/Chinsusu/Billing-V2/pull/305 for review.
- 2026-04-26: PR https://github.com/Chinsusu/Billing-V2/pull/305 merged into `main`; marking task done.
