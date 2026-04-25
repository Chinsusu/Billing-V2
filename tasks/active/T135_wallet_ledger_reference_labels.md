# T135 - Wallet ledger reference labels

Status: REVIEW
Owner: Codex
Branch: codex/t135-wallet-ledger-reference-labels
PR: https://github.com/Chinsusu/Billing-V2/pull/301
Risk: frontend
Created: 2026-04-26
Updated: 2026-04-26

## Summary

Use shared public reference labels for wallet ledger rows so reseller and client wallet screens show the linked record ID, not a backend or unrelated ledger reference.

## Scope

- Extract wallet ledger reference label formatting into a shared frontend helper.
- Keep client wallet ledger references on the existing public-label behavior.
- Update reseller wallet ledger references to use `reference_display_id` with type-specific prefixes.

## Acceptance Criteria

- Reseller wallet ledger references show `TUP-`, `INV-`, `ORD-`, or `TX-` labels when `reference_display_id` is present.
- Client and reseller wallet ledger screens share one formatter.
- Frontend lint, sensitive-text check, build, taskguard, and diff check pass.

## Notes

- Falling back to `LED-` is allowed only when the API has no linked `reference_display_id`.

## Agent Log

- 2026-04-26: Codex created and claimed the task after T134 was marked done; starting shared wallet ledger reference label cleanup.
- 2026-04-26: Added shared wallet ledger reference formatter and updated client/reseller wallet ledger rows to use it.
- 2026-04-26: Local validation passed: `npm --prefix frontend run lint`, `npm --prefix frontend run check:sensitive-text`, `npm --prefix frontend run build`, `go run ./cmd/taskguard`, and `git diff --check`.
- 2026-04-26: Opened PR https://github.com/Chinsusu/Billing-V2/pull/301 for review.
