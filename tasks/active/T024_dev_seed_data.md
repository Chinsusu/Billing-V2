# T024 - Development seed data

Status: IN_PROGRESS
Owner: Codex
Branch: feat/dev-seed-data
PR: -
Risk: seed/catalog/RBAC
Created: 2026-04-23
Updated: 2026-04-23

## Summary

Add an idempotent development seed runner for RBAC permissions, demo users, master catalog data, and a demo reseller catalog.

## Scope

- Add a `cmd/seed` entrypoint with `dev` and `plan` commands.
- Add reusable seed statements under `internal/seed`.
- Seed RBAC permissions, system roles, a platform tenant, a platform admin user, and a demo reseller user.
- Seed master catalog products/plans/provider sources and demo reseller tenant catalog clones.
- Add Makefile targets for building and running the seed command.
- Add tests that validate seed statement coverage and idempotency without requiring a real database.
- Out of scope: production seed data, secrets, real password setup, provider credentials, auth token generation, order/invoice/service data, or destructive reseeding.

## Acceptance Criteria

- `go run ./cmd/seed plan` works without `DB_DSN` and lists seed statement names.
- `go run ./cmd/seed dev` requires `DB_DSN` and applies idempotent seed statements in order.
- Seed SQL uses stable UUIDs and `ON CONFLICT` so reruns are safe.
- Seed data includes `catalog.view` and `catalog.manage` permissions required by T023.
- Full validation passes: `make fmt`, `make test`, `make build`, `make migrate-validate`.

## Notes

- Seed users use placeholder password hashes for local development only. They are not valid production credentials.

## Agent Log

- 2026-04-23: Task created and claimed from latest `origin/main` in `/tmp/Billing-T024`.
