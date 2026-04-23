# T021 - Catalog API runtime wiring

Status: REVIEW
Owner: Codex
Branch: feat/catalog-api-runtime-wiring
PR: https://github.com/Chinsusu/Billing-V2/pull/51
Risk: API/DB/catalog
Created: 2026-04-23
Updated: 2026-04-23

## Summary

Wire catalog HTTP routes into the API runtime by composing the PostgreSQL store, catalog service, and catalog handler when database config is present.

## Scope

- Open the PostgreSQL connection from `cmd/api` when `DB_DSN` is configured.
- Compose catalog store, service, and HTTP handler through the existing `app.APIOptions` hook.
- Keep health-only local startup available when `DB_DSN` is empty.
- Ensure opened database connections are closed on shutdown and on startup errors.
- Add focused tests for runtime option composition without requiring a real database.
- Out of scope: auth/RBAC enforcement, production tenant middleware, seed data, admin UI, checkout, order, invoice, worker startup, or provider provisioning.

## Acceptance Criteria

- `cmd/api` registers catalog routes only when a database connection is available.
- Missing `DB_DSN` keeps the existing health/readiness endpoints usable.
- Database open failures stop API startup with a clear error.
- Catalog route composition uses `catalog.NewPostgresStore`, `catalog.NewService`, and `catalog.NewHTTPHandler`.
- Full validation passes: `make fmt`, `make test`, `make build`, `make migrate-validate`.

## Notes

- This task wires runtime composition only. The temporary header-based tenant context from T020 remains until an auth/tenant middleware task replaces it.

## Agent Log

- 2026-04-23: Task created and claimed from latest `origin/main` in `/tmp/Billing-T021`.
- 2026-04-23: Opened PR #51. Validation passed: `go test ./cmd/api`, `make fmt`, `make test`, `make build`, `make migrate-validate`.
