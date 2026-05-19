# T277 - Notification local worker command

Status: DONE
Owner: Codex
Branch: codex/t277-notification-local-worker
PR: https://github.com/Chinsusu/Billing-V2/pull/585
Risk: notifications, workers, launch readiness, secret handling
Created: 2026-05-19
Updated: 2026-05-19

## Summary

Expose the existing notification local delivery runner through `cmd/worker` so local/dev notification delivery jobs can be exercised through the same worker entrypoint used by other job types.

## Scope

- Add `notification-local-once` and `notification-local-loop` worker commands.
- Wire the commands to `notification.NewLocalDeliveryRunner`.
- Fail closed outside `APP_ENV=local` or `APP_ENV=dev` because local delivery marks notifications sent without SMTP/Telegram.
- Add worker command tests for config, output, loop behavior, and environment guard.
- Update notification runbook notes to clarify that this is local/dev delivery plumbing only, not production SMTP/Telegram delivery proof.
- Keep production notification credentials, real customer messages, SMTP/Telegram integration, and launch decision broadening out of scope.

## Acceptance Criteria

- `cmd/worker notification-local-once` can run a local/dev notification delivery batch using `DB_DSN` or `-dsn`.
- `cmd/worker notification-local-loop` can run repeated local/dev notification delivery batches and uses existing loop output format.
- The new commands reject staging and production environments.
- Existing provisioning/lifecycle worker commands keep their behavior.
- Tests cover the new command dispatch and guard behavior.
- Docs do not claim automated production delivery proof.
- Task board remains consistent and required checks pass.

## Notes

- This task reduces the local/dev notification execution gap but does not close the broader production SMTP/Telegram evidence gap.
- Production SMTP/Telegram delivery still needs real channel config, secrets owner approval, safe delivery proof, and redacted evidence before broader GO.

## Agent Log

- 2026-05-19: Task created and claimed by Codex from Billing `origin/main`; selected next gap is notification worker execution plumbing, not production delivery proof.
- 2026-05-19: Added `notification-local-once` and `notification-local-loop` worker commands backed by `notification.NewLocalDeliveryRunner`.
- 2026-05-19: Added a local/dev-only guard so the commands reject staging and production because local delivery does not send SMTP/Telegram.
- 2026-05-19: Added worker command tests for once, loop, staging guard, and production guard; split notification local worker code into a separate file to keep `cmd/worker/main.go` under 500 lines.
- 2026-05-19: Updated the notification fallback runbook to document local/dev command usage and explicitly state this is not production SMTP/Telegram delivery proof.
- 2026-05-19: Validation passed: `go test ./cmd/worker`, `GOFLAGS=-buildvcs=false make test`, `GOFLAGS=-buildvcs=false go build ./cmd/api ./cmd/migrate ./cmd/seed ./cmd/smoke ./cmd/worker`, `GOFLAGS=-buildvcs=false go run ./cmd/taskguard`, `git diff --check`, diff raw-secret pattern scan, and changed-file line-count check. Raw `make test`/`go build` without `GOFLAGS` hit Go VCS stamping errors in the worktree, so the documented no-VCS equivalent was used.
- 2026-05-19: Opened PR #585.
- 2026-05-19: PR #585 merged after GitHub checks passed. Marking T277 DONE in marker PR.
