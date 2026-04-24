# T103 - Catalog admin read permission split

Status: DONE
Owner: Codex
Branch: codex/t103-catalog-admin-read-permission-split
PR: https://github.com/Chinsusu/Billing-V2/pull/234
Risk: backend/RBAC
Created: 2026-04-24
Updated: 2026-04-24

## Summary

Split admin catalog read routes from catalog mutation routes so provider readiness and catalog inspection can use read-level permission.

## Scope

- Add or reuse an admin catalog read middleware option.
- Keep catalog create/update routes on `catalog.manage`.
- Move admin catalog GET routes and provider readiness to `catalog.view`.
- Update focused tests for middleware selection.
- Update API docs and permission notes.
- Keep each file under 500 lines.

## Acceptance Criteria

- Admin catalog read routes require `catalog.view`.
- Admin catalog mutation routes still require `catalog.manage`.
- Existing reseller/client catalog access is unchanged.
- Backend and frontend validation commands pass.

## Notes

- Follows T100.
- Keep route behavior and response bodies unchanged.

## Agent Log

- 2026-04-24: Task created in the provider readiness follow-up batch.
- 2026-04-24: Codex claimed the task on `codex/t103-catalog-admin-read-permission-split`.
- 2026-04-24: Split admin catalog GET/provider readiness routes onto `catalog.view` while keeping create/update routes on `catalog.manage`; updated permission docs and middleware-selection tests.
- 2026-04-24: Validation passed: `go test ./internal/modules/catalog`, `go test ./cmd/api`, `go test ./...`, `go build ./cmd/api ./cmd/migrate ./cmd/seed ./cmd/smoke ./cmd/worker`, `npm ci`, `npm audit --omit=dev`, `npm run lint`, and `npm run build`.
- 2026-04-24: Opened PR #234 for review.
- 2026-04-24: CI passed on PR #234 and merged to main at `da68e72`.
