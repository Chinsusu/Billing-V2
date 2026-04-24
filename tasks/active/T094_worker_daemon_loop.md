# T094 - Provisioning worker daemon loop

Status: DONE
Owner: Codex
Branch: codex/t094-worker-daemon-loop
PR: https://github.com/Chinsusu/Billing-V2/pull/214
Risk: backend/worker
Created: 2026-04-24
Updated: 2026-04-24

## Summary

Add a long-running local/sandbox worker mode on top of the existing `provision-once` command.

## Scope

- Work mainly in `cmd/worker/**/*`, worker docs, and focused worker tests.
- Add a loop mode that repeatedly runs provisioning batches with interval, timeout, and graceful shutdown.
- Keep `provision-once` behavior unchanged.
- Keep the worker fake-provider friendly for local/dev.
- Keep each file under 500 lines.

## Acceptance Criteria

- `cmd/worker` supports a loop command for local/sandbox provisioning workers.
- Loop mode respects context cancellation and does not busy-spin when no jobs are claimed.
- Logs or summaries are clear enough to diagnose claimed/succeeded/retry/manual-review counts per pass.
- Unit tests cover command parsing and loop cancellation.
- Backend and frontend validation commands pass.

## Notes

- Should follow T089 and T092.
- Do not introduce production deployment automation in this task.

## Agent Log

- 2026-04-24: Task created after T092 completed and the active board was empty.
- 2026-04-24: Codex claimed the task on `codex/t094-worker-daemon-loop`.
- 2026-04-24: Added `provision-loop` with idle interval, timeout/signal cancellation, per-pass summaries, worker tests, and local operations docs.
- 2026-04-24: Validation passed: `go test ./cmd/worker`, `go test ./...`, `go build ./cmd/api ./cmd/migrate ./cmd/seed ./cmd/smoke ./cmd/worker`, `npm ci`, `npm audit --omit=dev`, `npm run lint`, `npm run build`.
- 2026-04-24: Opened PR https://github.com/Chinsusu/Billing-V2/pull/214 for review.
- 2026-04-24: PR #214 passed CI and merged into `main` at `52c79296e9035f78fd570707dd584492a49e989e`.
