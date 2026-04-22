# T006 - Outbox Job Skeleton

Status: DONE
Owner: Codex
Branch: feat/outbox-job-skeleton
PR: https://github.com/Chinsusu/Billing-V2/pull/27
Risk: worker/migration
Created: 2026-04-22
Updated: 2026-04-22

## Summary

Add outbox/jobs table model and worker claim interface after the DB skeleton exists.

## Scope

- Add outbox/job data shape.
- Add worker claim interface.
- Add retry/state fields needed by the documented worker model.
- Avoid real provider calls in this task.

## Acceptance Criteria

- Claim interface is safe for concurrent workers by design.
- Job states are clear and match docs.
- Retry behavior can be tested with fake implementations.
- `make test` passes.
- `make build` passes.

## Notes

- Follow async worker and database consistency docs before implementation.
- Keep provider behavior mocked or interface-only in this task.

## Agent Log

- 2026-04-22: Task file created from `TASKS.md`.
- 2026-04-22: Claimed by Codex. Adding outbox/jobs migration, job data shapes, and atomic claim interface skeleton.
- 2026-04-22: Opened PR https://github.com/Chinsusu/Billing-V2/pull/27. Validation passed: gofmt, make fmt, make test, make build, make migrate-validate.
- 2026-04-22: PR #27 merged into `main`; task completed.
