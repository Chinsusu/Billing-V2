# T028 - Order schema migration

Status: REVIEW
Owner: Codex
Branch: feat/order-schema-migration
PR: https://github.com/Chinsusu/Billing-V2/pull/66
Risk: migration/order/service
Created: 2026-04-23
Updated: 2026-04-23

## Summary

Add PostgreSQL tables for order, reservation, provisioning, and service lifecycle records.

## Scope

- Add migration `0007_create_order_tables.sql` with enum types, display ID sequences, tables, constraints, and indexes.
- Add rollback artifact under `migrations/rollback/`.
- Keep UUID primary keys and add `display_id BIGINT` on FE-visible records.
- Align table statuses with the order domain skeleton.
- Out of scope: repository/store implementation, HTTP handlers, wallet/ledger tables, invoice tables, provider execution, or frontend changes.

## Acceptance Criteria

- Migration creates tables for orders, order reservations, order provisioning jobs, and service instances.
- All new FE-visible tables have numeric `display_id` sequences starting at 10000.
- Core tenant, user, tenant plan, provider source, status, amount, currency, idempotency, and term-window constraints are present.
- Rollback artifact documents the destructive down path for clean/dev environments.
- Full validation passes: `make fmt`, `make test`, `make build`, `make migrate-validate`.

## Notes

- Do not edit older migration files; add only a new forward migration and rollback file.

## Agent Log

- 2026-04-23: Task created and claimed from latest `origin/main` in `/tmp/Billing-T028`.
- 2026-04-23: Opened PR #66. Validation passed: `make fmt`, `make test`, `make build`, `make migrate-validate`.
