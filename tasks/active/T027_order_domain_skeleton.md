# T027 - Order domain skeleton

Status: IN_PROGRESS
Owner: Codex
Branch: feat/order-domain-skeleton
PR: -
Risk: order/service/lifecycle
Created: 2026-04-23
Updated: 2026-04-23

## Summary

Add the first order domain skeleton for order, reservation, provisioning, and service lifecycle records.

## Scope

- Add order module IDs, statuses, lifecycle transition guards, input structs, and validation rules.
- Keep separate statuses for order, billing, reservation, provisioning, service, and suspension reason.
- Include `DisplayID` fields for FE-visible records while keeping UUID-style string IDs for storage.
- Add store interfaces that future PostgreSQL stores and API handlers can implement.
- Add focused unit tests for normalization, validation, and transition guards.
- Out of scope: database migrations, HTTP handlers, wallet/ledger debit, provider execution, invoice records, or frontend changes.

## Acceptance Criteria

- Order creation input validates tenant, buyer, tenant plan, amount, currency, quantity, status, billing status, and idempotency key.
- Reservation, provisioning job, and service instance inputs validate their required ids and lifecycle statuses.
- Order/provisioning/service transition helpers reject invalid transitions from the lifecycle document.
- `go test ./internal/modules/order` passes.
- Full validation passes: `make fmt`, `make test`, `make build`, `make migrate-validate`.

## Notes

- This task creates contracts only; persistence and APIs should be separate follow-up tasks.
- Future migrations must follow the display ID rule: UUID primary key plus `display_id BIGINT` for UI records.

## Agent Log

- 2026-04-23: Task created and claimed from latest `origin/main` in `/tmp/Billing-T027`.
