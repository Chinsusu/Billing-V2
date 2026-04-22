# T012 - Outbox Store And Worker Runner

Status: REVIEW
Owner: Codex
Branch: feat/outbox-store-runner
PR: https://github.com/Chinsusu/Billing-V2/pull/32
Risk: worker/outbox/retry
Created: 2026-04-23
Updated: 2026-04-23

## Summary

Implement PostgreSQL job/outbox claim stores and a small worker runner that can process claimed jobs with retry-safe completion behavior.

## Scope

- Add PostgreSQL implementations for `jobs.Store` and `jobs.OutboxStore`.
- Use atomic claim SQL with row locking and lock expiry.
- Add a worker runner that claims jobs, calls a handler, records attempts, and completes jobs with backoff/manual review behavior.
- Keep provider provisioning, notification publishing, scheduler, and external queue integrations out of scope.

## Acceptance Criteria

- Claim SQL only claims `queued`/`failed_retryable` jobs or `pending`/`failed_retryable` outbox events whose retry time has arrived.
- Claimed rows set `locked_by`, `locked_until`, and processing status atomically.
- Runner records attempts and completes success, retryable failure, and manual review paths.
- Tests cover runner success, retry, exhausted attempts, and validation guard behavior.
- `make fmt`, `make test`, `make build`, and `make migrate-validate` pass.

## Notes

- T012 builds on T006 model/migration and T011 DB executor helpers.

## Agent Log

- 2026-04-23: Codex claimed task from `origin/main` using isolated worktree `/tmp/Billing-T012`.
- 2026-04-23: Opened PR #32 after `make fmt`, `make test`, `make build`, and `make migrate-validate` passed.
