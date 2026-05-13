# T198 - Lifecycle scheduler jobs

Status: DONE
Owner: Codex
Branch: codex/t198-lifecycle-scheduler-jobs
PR: https://github.com/Chinsusu/Billing-V2/pull/427
Risk: scheduler, service lifecycle, billing, provisioning, and audit
Created: 2026-05-13
Updated: 2026-05-13

## Summary

Add bounded scheduler jobs for service expiry, grace, suspension, and termination processing.

## Scope

- Add scheduler/worker logic for lifecycle transitions defined by T197.
- Ensure jobs are idempotent and safe to retry.
- Add observability/audit for lifecycle job results.
- Do not bypass billing or provider capability rules.

## Acceptance Criteria

- Scheduler jobs process due services exactly once per effective transition.
- Retry does not duplicate lifecycle events or money effects.
- Tests cover no-op, success, retry, and failure/manual-review paths.
- Relevant backend validation and CI pass.

## Notes

- This task depends on T197.

## Agent Log

- 2026-05-13: Task created by Codex backlog planning.
- 2026-05-13: Codex claimed task on `codex/t198-lifecycle-scheduler-jobs`.
- 2026-05-13: Added bounded `service.lifecycle` scheduler/worker flow with stale job guards for term, billing, and suspension state.
- 2026-05-13: Local checks passed: `make fmt`, `go test ./internal/modules/order`, `go test ./cmd/worker ./internal/modules/jobs`, `make test`, `make build`, `make migrate-validate`, `make contract-guard`, `make error-code-guard`, `make task-guard`, `git diff --check`. `make smoke-dev-billing` blocked because `DB_DSN`/`-dsn` is not configured.
- 2026-05-13: Opened PR https://github.com/Chinsusu/Billing-V2/pull/427 and moved task to `REVIEW`.
- 2026-05-13: PR https://github.com/Chinsusu/Billing-V2/pull/427 merged as `12bc93c`; marking task `DONE`.
