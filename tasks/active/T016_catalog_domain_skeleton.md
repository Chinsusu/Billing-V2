# T016 - Catalog domain skeleton

Status: REVIEW
Owner: Codex
Branch: feat/catalog-domain-skeleton
PR: https://github.com/Chinsusu/Billing-V2/pull/41
Risk: catalog/pricing
Created: 2026-04-23
Updated: 2026-04-23

## Summary

Add the first backend catalog domain skeleton for master products, master plans, provider sources, plan-source mappings, and tenant catalog records.

## Scope

- Add catalog module IDs, statuses, input structs, and validation rules.
- Use minor-unit integer price fields; do not introduce floating-point money.
- Include `DisplayID` fields for FE-visible records.
- Add store interfaces that future PostgreSQL stores and API handlers can implement.
- Add focused unit tests for normalization and validation.
- Out of scope: database migration, HTTP handlers, pricing engine, checkout, order, invoice, or provider provisioning.

## Acceptance Criteria

- Catalog structs match the documented product/plan/source/tenant catalog shape.
- Invalid product type/status, billing cycle, currency, negative price, missing tenant, and invalid provider source inputs are rejected.
- Default JSON snapshots/policies normalize to `{}` where needed.
- `go test ./internal/modules/catalog` passes.
- Full validation passes: `make fmt`, `make test`, `make build`, `make migrate-validate`.

## Notes

- Catalog is money-adjacent; keep this PR limited to contracts and validation.
- Future migrations must follow the display ID rule: UUID PK plus `display_id BIGINT` for UI records.

## Agent Log

- 2026-04-23: Task created and claimed from latest `origin/main` in `/tmp/Billing-T016`.
- 2026-04-23: Opened PR #41. Validation passed: `go test ./internal/modules/catalog`, `make fmt`, `make test`, `make build`, `make migrate-validate`.
