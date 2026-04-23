# T031 - Order API handlers

Status: IN_PROGRESS
Owner: Codex
Branch: feat/order-api-handlers
PR: -
Risk: API/order
Created: 2026-04-23
Updated: 2026-04-23

## Summary

Add the first order HTTP API handler and runtime wiring for client order creation.

## Scope

- Add order HTTP handler for `POST /client/orders`.
- Read tenant scope from request context/header, buyer from actor context, and idempotency key from `Idempotency-Key`.
- Return order response DTOs with both UUID id and numeric display id.
- Wire order routes into API runtime when database config is present.
- Add focused handler and runtime tests.
- Out of scope: order list/detail APIs, reservation/provisioning/service APIs, wallet/ledger debit, provider execution, or frontend changes.

## Acceptance Criteria

- `POST /client/orders` creates an order through the order service.
- The handler does not trust tenant id or buyer id from the request body.
- Missing tenant, actor, idempotency key, or invalid fields return standard API errors.
- Runtime registers order routes only when `DB_DSN` is present.
- Full validation passes: `make fmt`, `make test`, `make build`, `make migrate-validate`.

## Notes

- Follow-up tasks should add order list/detail APIs, admin provisioning actions, and service management APIs.

## Agent Log

- 2026-04-23: Task created and claimed from latest `origin/main` in `/tmp/Billing-T031`.
