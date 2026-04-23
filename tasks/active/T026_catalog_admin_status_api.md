# T026 - Catalog admin status API

Status: REVIEW
Owner: Codex
Branch: feat/catalog-admin-update-api
PR: https://github.com/Chinsusu/Billing-V2/pull/62
Risk: API/catalog
Created: 2026-04-23
Updated: 2026-04-23

## Summary

Add admin catalog status update endpoints so operators can activate, disable, or archive catalog products, plans, and provider sources.

## Scope

- Add catalog store/service methods for product, plan, and provider source status updates.
- Add admin HTTP `PATCH` routes for product, plan, and provider source records.
- Return the updated record using the existing catalog response DTOs.
- Add focused tests for domain validation, SQL query builders, service delegation, and HTTP route parsing.
- Out of scope: full field edits, tenant catalog updates, plan source updates, audit event persistence, frontend changes, or provider runtime behavior changes.

## Acceptance Criteria

- `PATCH /admin/catalog/products/{product_id}` updates product status and returns the product DTO.
- `PATCH /admin/catalog/plans/{plan_id}` updates plan status and returns the plan DTO.
- `PATCH /admin/catalog/provider-sources/{source_id}` updates provider source status and returns the source DTO.
- Invalid ids/statuses return existing validation error responses.
- Missing records return existing catalog not found responses.
- Full validation passes: `make fmt`, `make test`, `make build`, `make migrate-validate`.

## Notes

- Create this branch from latest `origin/main`; do not base it on another task branch.

## Agent Log

- 2026-04-23: Task created and claimed from latest `origin/main` in `/tmp/Billing-T026`.
- 2026-04-23: Opened PR #62. Validation passed: `go test ./internal/modules/catalog`, `make fmt`, `make test`, `make build`, `make migrate-validate`.
