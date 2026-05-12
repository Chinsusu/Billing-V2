# T194 - Refund and adjustment ledger

Status: TODO
Owner: -
Branch: codex/t194-refund-adjustment-ledger
PR: -
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
