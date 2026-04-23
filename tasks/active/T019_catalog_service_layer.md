# T019 - Catalog service layer

Status: REVIEW
Owner: Codex
Branch: feat/catalog-service-layer
PR: https://github.com/Chinsusu/Billing-V2/pull/47
Risk: catalog/pricing
Created: 2026-04-23
Updated: 2026-04-23

## Summary

Add a catalog service/use-case layer above the store so API handlers do not call persistence directly.

## Scope

- Add catalog service methods for admin master product/plan/source setup.
- Add tenant product and tenant plan clone flows.
- Add margin guard behavior for tenant plan clones where selling price is lower than reseller cost.
- Add catalog list flows with tenant-scope validation.
- Add fake-store unit tests for service validation, normalization, and guard behavior.
- Out of scope: HTTP handlers, auth/RBAC enforcement, audit/outbox writes, seed data, checkout, order, invoice, or provisioning.

## Acceptance Criteria

- Service rejects nil store.
- Service normalizes and validates inputs before store calls.
- Tenant catalog list requires `tenant_id` before hitting store.
- Tenant plan clone with selling price below reseller cost becomes `margin_risk` instead of active.
- Full validation passes: `make fmt`, `make test`, `make build`, `make migrate-validate`.

## Notes

- API handlers should use this service in future tasks instead of calling `Store` directly.

## Agent Log

- 2026-04-23: Task created and claimed from latest `origin/main` in `/tmp/Billing-T019`.
- 2026-04-23: Opened PR #47. Validation passed: `go test ./internal/modules/catalog`, `make fmt`, `make test`, `make build`, `make migrate-validate`.
