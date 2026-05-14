# T212 - Cloudmini V3 worker wiring

Status: DONE
Owner: Codex
Branch: codex/t212-cloudmini-worker-wiring
PR: https://github.com/Chinsusu/Billing-V2/pull/455
Risk: provider/provisioning/credential/config
Created: 2026-05-14
Updated: 2026-05-14

## Summary

Wire the Cloudmini V3 provider adapter into the provisioning worker registry behind explicit environment configuration.

## Scope

- Keep `PROVIDER_DEFAULT_MODE=fake` as the default worker behavior.
- Add `PROVIDER_DEFAULT_MODE=cloudmini_v3` wiring that registers the real Cloudmini V3 adapter and keeps other provider types on fake adapters.
- Require explicit non-empty Cloudmini base URL, API token, Billing source id, group id, kind, protocol, and encryption key before the worker starts in Cloudmini mode.
- Add `.env.example` placeholders for Cloudmini V3 sandbox configuration without real credentials.
- Update provider sandbox evidence to state that worker wiring remains disabled by default and is not real sandbox proof.
- Do not run real Cloudmini API calls in this task.

## Acceptance Criteria

- Worker defaults to the existing fake provider registry with no Cloudmini env set.
- Cloudmini mode fails before starting if required env is missing or invalid.
- Cloudmini mode registers a `cloudmini_v3` adapter only for the configured source id mapping.
- No secret values are committed, logged, or printed in task notes.
- Worker and provider tests pass locally.

## Notes

- This task enables safe local/staging wiring only. Real sandbox readiness remains blocked until approved non-production credentials, quota, owner, source mapping, and cleanup evidence exist outside git.

## Agent Log

- 2026-05-14: Task created and claimed on `codex/t212-cloudmini-worker-wiring`.
- 2026-05-14: Added disabled-by-default Cloudmini V3 worker registry wiring, env placeholders, tests, and evidence note.
- 2026-05-14: Validation passed: `make fmt`; `go test ./cmd/worker`; `go test ./internal/modules/provider`; `go test ./internal/modules/provider -run SandboxContract`; `make test`; `make build`; `go run ./cmd/taskguard`; `git diff --check`; local CI secret-scan grep.
- 2026-05-14: Opened PR https://github.com/Chinsusu/Billing-V2/pull/455 and moved task to `REVIEW`.
- 2026-05-14: PR https://github.com/Chinsusu/Billing-V2/pull/455 merged; marking task `DONE`.
