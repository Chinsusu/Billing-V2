# T097 - Admin provisioning summary API

Status: DONE
Owner: Codex
Branch: codex/t097-admin-provisioning-summary-api
PR: https://github.com/Chinsusu/Billing-V2/pull/221
Risk: backend/API
Created: 2026-04-24
Updated: 2026-04-24

## Summary

Expose a compact admin read model for provisioning queue health so operators do not have to count job rows by hand.

## Scope

- Work mainly in `internal/modules/jobs/**/*`, `cmd/api/**/*`, API docs, and focused tests.
- Add an admin-only summary route for `provider.provision` jobs in the effective tenant.
- Include counts by status, total jobs, attention counts, oldest queued age/timestamp, and latest failure context.
- Keep response fields simple and display-ID friendly.
- Keep each file under 500 lines.

## Acceptance Criteria

- Admin can call a provisioning summary endpoint scoped to the effective tenant.
- Summary includes queued/running/succeeded/retryable/manual-review/terminal/cancelled counts.
- Summary makes operator attention obvious without exposing secrets.
- Missing store/service dependencies return the standard API error envelope.
- Backend and frontend validation commands pass.

## Notes

- Should follow T087, T088, T091, T094, and T096.
- Do not add production monitoring integration in this task.

## Agent Log

- 2026-04-24: Task created after T096 completed and the active board was fully DONE.
- 2026-04-24: Codex claimed the task on `codex/t097-admin-provisioning-summary-api`.
- 2026-04-24: Added `GET /admin/jobs/summary` with tenant-scoped job counts, attention count, oldest queued age, latest redacted failure context, and `provisioning.job.view` middleware wiring.
- 2026-04-24: Added focused jobs API/query/unit tests and updated API/runbook docs.
- 2026-04-24: Validation passed: `go test ./internal/modules/jobs`, `go test ./...`, `go build ./cmd/api ./cmd/migrate ./cmd/seed ./cmd/smoke ./cmd/worker`, `npm ci`, `npm audit --omit=dev`, `npm run lint`, `npm run build`.
- 2026-04-24: Opened PR https://github.com/Chinsusu/Billing-V2/pull/221 for review/CI.
- 2026-04-24: PR #221 merged into `main` at `c68b3d437c30d0afaca2164a60cb4d8a6a6aefbf`.
