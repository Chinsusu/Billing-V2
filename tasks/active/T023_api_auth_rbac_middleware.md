# T023 - API auth/RBAC middleware

Status: DONE
Owner: Codex
Branch: feat/api-auth-rbac-middleware
PR: https://github.com/Chinsusu/Billing-V2/pull/56
Risk: API/auth/RBAC
Created: 2026-04-23
Updated: 2026-04-23

## Summary

Add actor context and RBAC HTTP middleware, then protect catalog routes with catalog permissions.

## Scope

- Add header-based actor context middleware for the current API skeleton.
- Add RBAC HTTP middleware that reads actor and tenant context, calls an `rbac.Authorizer`, and maps auth errors to API responses.
- Add catalog permissions and a store-backed authorizer skeleton.
- Let catalog handlers accept route middleware options for admin, reseller, and client routes.
- Wire catalog runtime routes with identity actor middleware and RBAC permission checks.
- Add focused tests for actor context, RBAC middleware, catalog route protection, and runtime composition.
- Out of scope: JWT/session validation, password login, production tenant resolution, token signing, refresh tokens, audit writes, and seed data.

## Acceptance Criteria

- Missing actor context returns an auth error before protected catalog route handlers call the service.
- Catalog admin routes require catalog manage permission.
- Catalog reseller/client routes require catalog view/manage permissions after tenant context is present.
- Runtime wiring composes identity header context, RBAC middleware, catalog service, and catalog handler.
- Existing health/readiness behavior remains unchanged.
- Full validation passes: `make fmt`, `make test`, `make build`, `make migrate-validate`.

## Notes

- Header-based actor context is a development adapter only. A later auth task must replace it with authenticated token/session resolution.

## Agent Log

- 2026-04-23: Task created and claimed from latest `origin/main` in `/tmp/Billing-T023`.
- 2026-04-23: Opened PR #56. Validation passed: `go test ./internal/modules/identity ./internal/modules/rbac ./internal/modules/catalog ./cmd/api`, `make fmt`, `make test`, `make build`, `make migrate-validate`.
- 2026-04-23: PR #56 merged into `main`.
