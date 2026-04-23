# T022 - HTTP tenant context middleware

Status: REVIEW
Owner: Codex
Branch: feat/http-tenant-context-middleware
PR: https://github.com/Chinsusu/Billing-V2/pull/53
Risk: API/tenant
Created: 2026-04-23
Updated: 2026-04-23

## Summary

Add a small HTTP tenant context middleware so tenant-scoped handlers read tenant scope from request context instead of parsing tenant headers directly.

## Scope

- Add tenant HTTP header constants and middleware that attaches `tenant.Context` to the request context.
- Apply the middleware to catalog reseller/client tenant routes.
- Update catalog tenant handlers to require tenant context before calling the service.
- Keep admin catalog routes independent from tenant context.
- Add focused tenant middleware and catalog handler tests.
- Out of scope: full auth, RBAC enforcement, signed tokens, production tenant resolution, session management, or frontend changes.

## Acceptance Criteria

- Catalog tenant routes read tenant id from `tenant.RequireContext`.
- Missing tenant context is rejected before service calls.
- Header-based tenant resolution is isolated in tenant middleware, not repeated inside catalog handlers.
- Existing header-driven tests still pass through the middleware.
- Full validation passes: `make fmt`, `make test`, `make build`, `make migrate-validate`.

## Notes

- This is still a skeleton. A later auth task should replace raw tenant headers with authenticated tenant resolution.

## Agent Log

- 2026-04-23: Task created and claimed from latest `origin/main` in `/tmp/Billing-T022`.
- 2026-04-23: Opened PR #53. Validation passed: `go test ./internal/modules/tenant ./internal/modules/catalog`, `make fmt`, `make test`, `make build`, `make migrate-validate`.
