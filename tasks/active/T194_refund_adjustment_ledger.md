# T194 - Refund and adjustment ledger

Status: REVIEW
Owner: Codex
Branch: codex/t194-refund-adjustment-ledger
PR: https://github.com/Chinsusu/Billing-V2/pull/419
Risk: wallet, ledger, payment, settlement, audit, and finance reconciliation
Created: 2026-05-13
Updated: 2026-05-13

## Summary

Add refund and adjustment ledger behavior required by MVP billing safety.

## Scope

- Implement append-only refund and adjustment entries for wallet/ledger flows.
- Require reason and actor context for manual adjustments.
- Preserve historical ledger immutability.
- Add finance/audit evidence for refund and adjustment actions.

## Acceptance Criteria

- Refunds and adjustments create new ledger entries rather than mutating historical entries.
- Duplicate or replayed refund/adjustment requests are idempotent where required.
- Tests cover success, insufficient context, duplicate request, and audit behavior.
- Relevant backend validation and CI pass.

## Notes

- Stop and ask if refund policy, partial refund behavior, or approval policy is unclear.

## Agent Log

- 2026-05-13: Task created by Codex backlog planning.
- 2026-05-13: Claimed by Codex on branch `codex/t194-refund-adjustment-ledger`.
- 2026-05-13: Implemented append-only admin wallet refund/adjustment ledger routes with actor, reason, idempotency conflict checks, and audit evidence. Local validation passed: focused wallet/audit/API/guard tests, `make fmt`, `make test`, `make build`, `make migrate-validate`, `make contract-guard`, `make error-code-guard`, `make task-guard`, `git diff --check`, and secret grep review.
- 2026-05-13: Opened PR https://github.com/Chinsusu/Billing-V2/pull/419 and moved task to `REVIEW`.
