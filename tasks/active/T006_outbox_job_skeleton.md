# T006 - Outbox Job Skeleton

Status: TODO
Owner: -
Branch: feat/outbox-job-skeleton
PR: -
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
