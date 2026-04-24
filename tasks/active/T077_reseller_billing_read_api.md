# T077 - Reseller billing read APIs

Status: REVIEW
Owner: Codex
Branch: codex/t077-reseller-billing-read-api
PR: https://github.com/Chinsusu/Billing-V2/pull/177
Risk: backend/API
Created: 2026-04-24
Updated: 2026-04-24

## Summary

Add reseller-scoped read endpoints for invoices, transactions, wallets, ledger, and top-up requests.

## Scope

- Work mainly in `internal/modules/invoice/**/*`, `internal/modules/payment/**/*`, `internal/modules/wallet/**/*`, `cmd/api/main.go`, and API docs.
- Add reseller routes that mirror the useful admin read endpoints but are tenant-scoped.
- Keep owner/tenant scoping enforced in backend queries.
- Return numeric display IDs for visible records.

## Acceptance Criteria

- Reseller can list tenant invoices, transactions, wallets, wallet ledger, and top-up requests without admin routes.
- Routes use reseller wallet/billing view permission middleware.
- Tests cover tenant scoping and bad filter validation.
- `go test ./...` passes and backend binaries build.

## Notes

- Do not add top-up approval, wallet adjustment, refund, or reconciliation mutations in this task.
- Keep the route set small if an endpoint needs missing store support; create follow-up work instead of forcing a large refactor.

## Agent Log

- 2026-04-24: Task created after T074 completed and the board needed the next reseller/live workflow batch.
- 2026-04-24: Codex claimed the task after T076 completed and started adding reseller billing/wallet read routes.
- 2026-04-24: Added reseller read routes for invoices, transactions, wallets, wallet ledger, and top-up requests with runtime middleware wiring and docs.
- 2026-04-24: Split runtime protection tests out of `cmd/api/main_test.go` to keep files under 500 lines.
- 2026-04-24: Validation passed: `go test ./...` and `go build ./cmd/api ./cmd/migrate ./cmd/seed ./cmd/smoke`.
- 2026-04-24: Opened PR #177 for review/CI.
