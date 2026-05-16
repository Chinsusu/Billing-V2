# T225 - RBAC route surface isolation

Status: DONE
Owner: Codex
Branch: codex/t225-rbac-route-surface-isolation
PR: https://github.com/Chinsusu/Billing-V2/pull/481
Risk: auth/RBAC/tenant-isolation/credential-safety
Created: 2026-05-16
Updated: 2026-05-16

## Summary

Fix portal route isolation so client sessions cannot access reseller or admin API surfaces, and reseller sessions cannot access admin API surfaces.

## Scope

- Enforce the expected actor type for admin, reseller, and client route middleware.
- Keep existing tenant scoping, permission checks, and admin 2FA behavior.
- Add backend tests for client-to-reseller/admin and reseller-to-admin denial.
- Validate the fix locally and against the live test domains after deployment.

## Acceptance Criteria

- Client session can access client-owned service routes but receives a forbidden response for reseller/admin service routes.
- Reseller session can access reseller service routes but receives a forbidden response for admin service routes.
- Platform admin routes still require satisfied 2FA when session auth is active.
- Focused RBAC/auth tests, full backend tests, build, taskguard, and diff checks pass.

## Notes

- Found during target-environment RBAC diagnostic on `2026-05-16`: client and reseller sessions could call higher-scope route surfaces with status `200`.

## Agent Log

- 2026-05-16: Task created and claimed from latest `origin/main`.
- 2026-05-16: Opened PR #481. Validation passed: `go test ./internal/modules/rbac ./cmd/api`; `make test`; `make build`; `go run ./cmd/contractguard`; `go run ./cmd/errorcodeguard`; `go run ./cmd/taskguard`; `git diff --check`.
- 2026-05-16: PR #481 merged into `main`; marking task DONE.
