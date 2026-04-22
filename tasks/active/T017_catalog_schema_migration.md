# T017 - Catalog schema migration

Status: IN_PROGRESS
Owner: Codex
Branch: feat/catalog-schema-migration
PR: -
Risk: migration/catalog
Created: 2026-04-23
Updated: 2026-04-23

## Summary

Add PostgreSQL schema for catalog records after T016 introduced the catalog domain contracts.

## Scope

- Add migration for master products, master plans, provider sources, plan-source mappings, tenant products, and tenant plans.
- Include `display_id BIGINT` for all FE-visible catalog records.
- Add constraints and indexes for status, tenant scope, unique catalog keys, and common list queries.
- Add manual rollback artifact for clean/dev environments.
- Out of scope: repository implementation, seed data, API handlers, pricing engine, checkout, order, invoice, or provisioning.

## Acceptance Criteria

- Migration validates with `make migrate-validate`.
- Tables use UUID primary keys plus numeric `display_id` sequences starting at 10000.
- Tenant-scoped tables include `tenant_id` and tenant indexes.
- Money-like amounts are stored as minor-unit integers with non-negative checks.
- Manual rollback drops catalog tables, sequences, enum types, and schema migration row for version `0006`.
- Full validation passes: `make fmt`, `make test`, `make build`, `make migrate-validate`.

## Notes

- This migration has data impact once used; rollback is destructive and clean/dev only unless an owner approves backup/restore.

## Agent Log

- 2026-04-23: Task created and claimed from latest `origin/main` in `/tmp/Billing-T016`.
