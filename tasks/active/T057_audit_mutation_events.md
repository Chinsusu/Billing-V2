# T057 - Audit mutation events

Status: REVIEW
Owner: Codex
Branch: feat/audit-mutation-events
PR: https://github.com/Chinsusu/Billing-V2/pull/129
Risk: audit/mutation
Created: 2026-04-23
Updated: 2026-04-23

## Summary

Write audit events for important state-changing billing operations so the audit read API has useful operational history.

## Scope

- Add audit writes for top-up approve/reject actions.
- Add audit writes for invoice wallet payment.
- Add audit writes for order status transitions.
- Include tenant, actor, target record, action, and before/after status where available.
- Keep audit failures from silently hiding money-flow errors; choose explicit error behavior.

## Acceptance Criteria

- Mutation tests verify audit rows are written for each covered operation.
- Audit records are tenant-scoped and visible through the existing audit read API.
- No secrets or sensitive provider credentials are written into audit metadata.
- Backend quality gates pass.

## Notes

- Prefer small service-level helpers over duplicating raw insert SQL in handlers.

## Agent Log

- 2026-04-23: Task created after audit read API landed.
- 2026-04-23: Codex claimed task from latest `origin/main` after T056 merged.
- 2026-04-23: PR #129 opened. Validation: `go test ./internal/modules/wallet ./internal/modules/order ./internal/modules/payment ./cmd/api`, `make fmt`, `make test`, `make build`, `make migrate-validate`.
