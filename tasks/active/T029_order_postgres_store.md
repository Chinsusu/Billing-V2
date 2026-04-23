# T029 - Order PostgreSQL store

Status: IN_PROGRESS
Owner: Codex
Branch: feat/order-postgres-store
PR: -
Risk: order/DB
Created: 2026-04-23
Updated: 2026-04-23

## Summary

Add PostgreSQL store implementation for the order domain create contracts.

## Scope

- Add `PostgresStore` for order, reservation, provisioning job, and service instance creates.
- Add scan helpers for the new order schema tables.
- Normalize and validate inputs before building database arguments.
- Add focused unit tests for store argument builders and validation behavior.
- Out of scope: list/update APIs, transaction orchestration, wallet/ledger debit, provider execution, or HTTP handlers.

## Acceptance Criteria

- `PostgresStore` satisfies the order `Store` interface.
- Create methods return the inserted record with UUID id and numeric display id fields.
- Invalid inputs fail before database execution.
- `go test ./internal/modules/order` passes.
- Full validation passes: `make fmt`, `make test`, `make build`, `make migrate-validate`.

## Notes

- This task assumes T028 migration is already merged.

## Agent Log

- 2026-04-23: Task created and claimed from latest `origin/main` in `/tmp/Billing-T029`.
