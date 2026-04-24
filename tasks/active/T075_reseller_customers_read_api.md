# T075 - Reseller customers read API

Status: IN_PROGRESS
Owner: Codex
Branch: codex/t075-reseller-customers-read-api
PR: -
Risk: backend/API
Created: 2026-04-24
Updated: 2026-04-24

## Summary

Add a reseller-scoped customer account read endpoint so reseller screens can list their own clients with numeric display IDs.

## Scope

- Work mainly in `internal/modules/identity/**/*`, `cmd/api/main.go`, and API docs if endpoint tables need updates.
- Add `GET /reseller/customers` using the same response shape as admin account reads.
- Scope reads to the effective tenant from headers/context and default `user_type` to `client`.
- Support simple filters already available for account reads: `display_id`, `status`, `email`, and `limit`.

## Acceptance Criteria

- Reseller users can list only customers in their own tenant.
- The route uses reseller permission middleware and does not bypass actor/RBAC checks.
- Response includes `display_id` for frontend tables.
- Handler/runtime tests cover route registration, tenant scoping, default client type, and bad filter validation.
- `go test ./...` passes and backend binaries build.

## Notes

- Keep admin account behavior unchanged.
- Do not add create/edit/delete customer actions in this task.

## Agent Log

- 2026-04-24: Task created after T074 completed and the board needed the next reseller/live workflow batch.
- 2026-04-24: Codex claimed the task and started adding the reseller customer read route from latest `main`.
- 2026-04-24: Added `GET /reseller/customers`, reseller account middleware wiring, handler/runtime tests, and API contract documentation.
- 2026-04-24: Validation passed: `go test ./...` and `go build ./cmd/api ./cmd/migrate ./cmd/seed ./cmd/smoke`.
