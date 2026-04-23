# T038 - Service instance store API

Status: DONE
Owner: Codex
Branch: feat/service-instance-store-api
PR: https://github.com/Chinsusu/Billing-V2/pull/88
Risk: service/API/DB
Created: 2026-04-23
Updated: 2026-04-23

## Summary

Add store methods and read APIs for provisioned service instances so customer and admin portals can list active services.

## Scope

- Add PostgreSQL store methods for service instance create, list, and detail reads.
- Add client routes for the current tenant/account service list and detail.
- Add admin routes for tenant-scoped service instance list and detail.
- Include numeric display IDs in API responses.
- Out of scope: service renewals, suspension, provider sync, or frontend views.

## Acceptance Criteria

- Client reads are scoped to tenant and account context.
- Admin reads are scoped to tenant context and require service read permission.
- Missing or cross-tenant service IDs return standard not-found errors.
- Handler and runtime tests cover client and admin access.
- Full validation passes: `make fmt`, `make test`, `make build`, `make migrate-validate`.

## Notes

- This task can start after the service instance schema and order lifecycle paths are stable.

## Agent Log

- 2026-04-23: Task created for the next backend batch.
- 2026-04-23: Claimed by Codex from latest `origin/main` in `/tmp/Billing-T038`.
- 2026-04-23: Opened PR #88. Validation passed: `go test ./internal/modules/order ./cmd/api`, `make fmt`, `make test`, `make build`, `make migrate-validate`.
- 2026-04-23: PR #88 merged into `main`.
