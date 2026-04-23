# T047 - Wallet ledger posting service

Status: REVIEW
Owner: Codex
Branch: feat/wallet-ledger-posting-service
PR: pending
Risk: wallet/money/DB
Created: 2026-04-23
Updated: 2026-04-23

## Summary

Add an atomic wallet ledger posting path that updates wallet balances and writes ledger rows in one operation.

## Scope

- Add service/store method for posting credit and debit ledger entries.
- Update wallet balances atomically and store `balance_after_minor`.
- Keep idempotency per wallet/action to avoid duplicate balance movement.
- Reject debit entries that would make available balance negative.
- Reuse existing wallet and ledger read APIs.

## Acceptance Criteria

- Credit and debit postings update wallet balances and ledger rows consistently.
- Duplicate idempotency keys return the existing ledger entry without moving balance again.
- Insufficient balance fails with a clear conflict/domain error.
- Store/service tests cover credit, debit, duplicate, and insufficient balance cases.
- Full validation passes: `make fmt`, `make test`, `make build`, `make migrate-validate`.

## Notes

- This task should be done before approval/payment workflows that mutate wallet balances.

## Agent Log

- 2026-04-23: Task created after T046 to make wallet money movement safe.
- 2026-04-23: Implemented atomic ledger posting service/store path with balance updates, idempotent SQL, and focused tests.
