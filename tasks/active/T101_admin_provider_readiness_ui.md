# T101 - Admin provider readiness UI

Status: DONE
Owner: Codex
Branch: codex/t101-admin-provider-readiness-ui
PR: https://github.com/Chinsusu/Billing-V2/pull/230
Risk: frontend/admin-ops
Created: 2026-04-24
Updated: 2026-04-24

## Summary

Wire the provider readiness API into the admin portal so operators can inspect plan/source readiness without using curl or database rows.

## Scope

- Add frontend API types and helper for `GET /admin/catalog/provider-readiness`.
- Add a compact readiness panel to the admin providers or provisioning area.
- Show plan/source display IDs first, product type, source type, state, and reason.
- Handle loading, empty, error, and live API fallback states.
- Do not expose provider credentials, raw provider payloads, or capability JSON.
- Keep each frontend file under 500 lines.

## Acceptance Criteria

- Admin UI shows readiness rows from the live API when available.
- State badges cover `ready`, `inactive_source`, `missing_plan_source`, `unsupported_capability`, and `fake_provider_only`.
- Operator-facing rows use numeric display IDs before UUIDs.
- Frontend and backend validation commands pass.
- Browser verification covers desktop and mobile layout.

## Notes

- Follows T100.
- This task should not change backend readiness semantics.

## Agent Log

- 2026-04-24: Task created in the provider readiness follow-up batch.
- 2026-04-24: Codex claimed the task on `codex/t101-admin-provider-readiness-ui`.
- 2026-04-24: Added frontend API types/helper for admin provider readiness and a Providers screen readiness panel with display-ID-first rows and safe demo fallback.
- 2026-04-24: Browser verification passed with mocked live readiness data on desktop and mobile; confirmed Ready, Fake only, and Missing source states render without capability JSON or raw provider payloads.
- 2026-04-24: Validation passed: `npm ci`, `npm audit --omit=dev`, `npm run lint`, `npm run build`, `go test ./...`, and `go build ./cmd/api ./cmd/migrate ./cmd/seed ./cmd/smoke ./cmd/worker`.
- 2026-04-24: Opened PR https://github.com/Chinsusu/Billing-V2/pull/230 for review/CI.
- 2026-04-24: CI passed on PR #230 and merged to main at `116bd25`.
