# T025 - Catalog admin read API

Status: IN_PROGRESS
Owner: Codex
Branch: feat/catalog-admin-read-api
PR: -
Risk: API/catalog
Created: 2026-04-23
Updated: 2026-04-23

## Summary

Add admin catalog read/list API endpoints for products, plans, and provider sources.

## Scope

- Add catalog store/service list methods for master products and provider sources.
- Add admin HTTP list endpoints for products, plans, and provider sources with basic filters.
- Keep existing create routes on the same paths by dispatching on HTTP method.
- Add focused tests for query builders, service delegation, and HTTP filter parsing.
- Out of scope: update/disable/delete endpoints, detail endpoints, cursor tokens, RBAC permission changes, seed data, order/invoice/service APIs, or frontend changes.

## Acceptance Criteria

- `GET /admin/catalog/products` returns product DTOs and supports `product_type`, `status`, and `limit`.
- `GET /admin/catalog/plans` returns plan DTOs and supports existing master plan filters.
- `GET /admin/catalog/provider-sources` returns provider source DTOs and supports `source_type`, `status`, and `limit`.
- Existing POST create routes still work on the same paths.
- Full validation passes: `make fmt`, `make test`, `make build`, `make migrate-validate`.

## Agent Log

- 2026-04-23: Task created and claimed from latest `origin/main` in `/tmp/Billing-T025`.
