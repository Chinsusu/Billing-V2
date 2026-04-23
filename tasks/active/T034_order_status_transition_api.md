# T034 - Order status transition API

Status: TODO
Owner: -
Branch: feat/order-status-transition-api
PR: -
Risk: API/order/status
Created: 2026-04-23
Updated: 2026-04-23

## Summary

Add a tenant-scoped admin API to move an order from one status to another through the order service instead of direct database edits.

## Scope

- Add a service method that validates allowed order status changes.
- Add a PostgreSQL update method that changes one tenant-scoped order only when the expected current status matches.
- Add an admin route such as `PATCH /admin/orders/{order_id}/status`.
- Require actor, tenant, and RBAC checks for the admin mutation route.
- Keep billing status explicit in the request so payment-related changes are easy to audit.
- Out of scope: refunds, payment capture, provider provisioning, invoices, or frontend changes.

## Acceptance Criteria

- Invalid status changes return a standard validation or conflict error.
- Cross-tenant updates cannot change an order.
- Runtime wiring protects the route with admin auth/RBAC middleware.
- Focused service, store, handler, and runtime tests cover success and failure paths.
- Full validation passes: `make fmt`, `make test`, `make build`, `make migrate-validate`.

## Notes

- Start from latest `origin/main`, not from another task branch.
- This task is the first backend task in the 6-hour queue.

## Agent Log

- 2026-04-23: Task created for the next backend batch.
