# T033 - Admin order read API

Status: IN_PROGRESS
Owner: Codex
Branch: feat/admin-order-read-api
PR: -
Risk: API/order/admin
Created: 2026-04-23
Updated: 2026-04-23

## Summary

Add tenant-scoped admin order list/detail read APIs for platform and reseller operations.

## Scope

- Add `GET /admin/orders` with status, billing status, buyer, and limit filters.
- Add `GET /admin/orders/{order_id}` for tenant-scoped admin order detail.
- Reuse existing order store/service read contracts.
- Wire admin order read routes to `order.view` permission.
- Add focused handler and runtime tests.
- Out of scope: order mutation, refund/cancel, provisioning actions, invoices, ledger, or frontend changes.

## Acceptance Criteria

- Admin order reads require tenant context and actor context.
- Admin list/detail are scoped by tenant context, with optional buyer filter on list.
- Missing/invalid order id, status, billing status, or limit returns standard API errors.
- Runtime protects admin order routes with auth/RBAC middleware.
- Full validation passes: `make fmt`, `make test`, `make build`, `make migrate-validate`.

## Notes

- This complements T032 client read APIs and prepares admin portal order views.

## Agent Log

- 2026-04-23: Task created and claimed from latest `origin/main` in `/tmp/Billing-T033`.
