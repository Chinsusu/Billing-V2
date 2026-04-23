# T032 - Order read API

Status: IN_PROGRESS
Owner: Codex
Branch: feat/order-read-api
PR: -
Risk: API/order/DB
Created: 2026-04-23
Updated: 2026-04-23

## Summary

Add tenant-scoped order list/detail read APIs for the client portal.

## Scope

- Add order store/service read contracts for list and detail lookup.
- Implement PostgreSQL list/detail queries scoped by tenant and buyer.
- Add `GET /client/orders` with status, billing status, and limit filters.
- Add `GET /client/orders/{order_id}` for client order detail.
- Add focused tests for query builders, service delegation, and HTTP parsing.
- Out of scope: admin order views, cursor tokens, order update/cancel APIs, invoices, ledger, provider actions, or frontend changes.

## Acceptance Criteria

- `GET /client/orders` returns only orders for the current tenant and actor.
- `GET /client/orders/{order_id}` requires tenant context and actor context.
- Invalid status, billing status, order id, or limit returns standard API validation errors.
- Missing order returns standard not-found response.
- Full validation passes: `make fmt`, `make test`, `make build`, `make migrate-validate`.

## Notes

- This builds on T031 runtime wiring; no new runtime option is needed.

## Agent Log

- 2026-04-23: Task created and claimed from latest `origin/main` in `/tmp/Billing-T032`.
