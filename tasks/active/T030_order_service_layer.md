# T030 - Order service layer

Status: IN_PROGRESS
Owner: Codex
Branch: feat/order-service-layer
PR: -
Risk: order/service
Created: 2026-04-23
Updated: 2026-04-23

## Summary

Add order service layer methods on top of the order store contracts.

## Scope

- Add `Service` wrapper for order, reservation, provisioning job, and service instance creates.
- Normalize and validate inputs before delegating to the store.
- Add focused unit tests for missing store, normalization, validation, and delegation.
- Out of scope: HTTP handlers, transaction orchestration, wallet/ledger debit, provider execution, or runtime wiring.

## Acceptance Criteria

- Service methods return a clear error when no store is configured.
- Invalid inputs fail before store delegation.
- Valid inputs are normalized before reaching the store.
- `go test ./internal/modules/order` passes.
- Full validation passes: `make fmt`, `make test`, `make build`, `make migrate-validate`.

## Notes

- This prepares the order module for API/runtime wiring without coupling handlers to persistence.

## Agent Log

- 2026-04-23: Task created and claimed from latest `origin/main` in `/tmp/Billing-T030`.
