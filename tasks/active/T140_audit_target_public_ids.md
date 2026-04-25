# T140 - Audit target public IDs

Status: DONE
Owner: Codex
Branch: codex/t140-audit-target-public-ids
PR: https://github.com/Chinsusu/Billing-V2/pull/311
Risk: backend/frontend
Created: 2026-04-26
Updated: 2026-04-26

## Summary

Resolve service and provider-source public IDs for live audit log target labels and filters.

## Scope

- Add audit target display ID resolution for service instances.
- Add audit target display ID resolution for provider sources.
- Extend audit target display ID filtering for the same target types.
- Extend frontend audit target labels with `SVC-*` and `SRC-*` prefixes.

## Acceptance Criteria

- Audit log API can return `target_display_id` for service targets.
- Audit log API can return `target_display_id` for provider-source targets.
- Admin audit target display ID filter can match service and provider-source targets.
- Frontend audit rows label service and provider-source targets with public prefixes.
- Go tests, frontend lint/build, taskguard, and diff check pass.

## Notes

- Existing invoice, order, job, and top-up audit target behavior must remain unchanged.

## Agent Log

- 2026-04-26: Codex created and claimed the task after T139 was marked done; starting live audit target public ID support.
- 2026-04-26: Added service/provider-source target display ID joins, matching target display ID filters, and frontend `SVC-*` / `SRC-*` audit target labels.
- 2026-04-26: Local validation passed: `go test ./internal/modules/audit`, `go test ./...`, `npm --prefix frontend run lint`, `npm --prefix frontend run check:sensitive-text`, `npm --prefix frontend run build`, `go run ./cmd/taskguard`, and `git diff --check`.
- 2026-04-26: Opened PR https://github.com/Chinsusu/Billing-V2/pull/311 for review.
- 2026-04-26: PR https://github.com/Chinsusu/Billing-V2/pull/311 merged into `main`; marking task done.
