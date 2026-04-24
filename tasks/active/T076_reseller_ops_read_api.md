# T076 - Reseller ops read APIs

Status: TODO
Owner: -
Branch: codex/t076-reseller-ops-read-api
PR: -
Risk: backend/API
Created: 2026-04-24
Updated: 2026-04-24

## Summary

Add reseller-scoped read endpoints for orders and service inventory so reseller operators can monitor their tenant without using admin routes.

## Scope

- Work mainly in `internal/modules/order/**/*`, `cmd/api/main.go`, and API docs.
- Add reseller routes for listing orders and services using existing read models where possible.
- Scope all queries to the effective tenant and keep UUIDs internal while returning numeric display IDs.
- Reuse existing pagination and filters where they already exist.

## Acceptance Criteria

- `GET /reseller/orders` returns only orders in the reseller tenant.
- `GET /reseller/services` returns only services in the reseller tenant.
- Routes use reseller view permission middleware.
- Tests cover tenant scoping, filter parsing, and route registration.
- `go test ./...` passes and backend binaries build.

## Notes

- Do not add mutation endpoints for order status or service lifecycle in this task.
- Keep response contracts aligned with existing admin/client order and service response shapes.

## Agent Log

- 2026-04-24: Task created after T074 completed and the board needed the next reseller/live workflow batch.
