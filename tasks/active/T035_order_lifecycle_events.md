# T035 - Order lifecycle events

Status: REVIEW
Owner: Codex
Branch: feat/order-lifecycle-events
PR: https://github.com/Chinsusu/Billing-V2/pull/82
Risk: order/outbox
Created: 2026-04-23
Updated: 2026-04-23

## Summary

Publish outbox events when orders are created and when order status changes, so later workers can react without polling every table.

## Scope

- Define clear event names and payloads for order creation and status changes.
- Write outbox records from the order service or store path in the same logical operation as the order change.
- Add tests that verify expected outbox records are created.
- Document which event starts provisioning and which event is only for audit/history.
- Out of scope: running provider calls, invoice generation, or email/webhook delivery.

## Acceptance Criteria

- Creating an order emits one order-created event.
- Changing order status emits one status-changed event with old and new status values.
- Event payloads use stable IDs and include the numeric display ID when available.
- Failed order changes do not emit events.
- Full validation passes: `make fmt`, `make test`, `make build`, `make migrate-validate`.

## Notes

- This task should start after T034 is merged so it can hook into the status transition path.

## Agent Log

- 2026-04-23: Task created for the next backend batch.
- 2026-04-23: Claimed by Codex from latest `origin/main` in `/tmp/Billing-T035`.
- 2026-04-23: Opened PR #82. Validation passed: `go test ./internal/modules/order`, `make fmt`, `make test`, `make build`, `make migrate-validate`.
