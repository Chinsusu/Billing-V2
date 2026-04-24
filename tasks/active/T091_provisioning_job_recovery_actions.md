# T091 - Provisioning job recovery actions

Status: TODO
Owner: -
Branch: codex/t091-provisioning-job-recovery-actions
PR: -
Risk: backend/ops
Created: 2026-04-24
Updated: 2026-04-24

## Summary

Add guarded admin operations for provisioning job recovery after the read APIs and worker flow are in place.

## Scope

- Work mainly in `internal/modules/jobs/**/*`, `cmd/api/**/*`, audit logging, and operational docs.
- Add minimal actions only where state transitions are clear: retry retryable/manual-review jobs, cancel safe jobs, or mark manual review.
- Require high-risk permission middleware and audit records.
- Keep each file under 500 lines.

## Acceptance Criteria

- Admin can retry a safe failed/manual-review provisioning job with a new next-attempt time.
- Admin can move a job to manual review with a required reason.
- Admin can cancel only jobs that have not succeeded.
- Actions are idempotent enough for duplicate submits and return clear conflict errors.
- Backend and frontend validation commands pass.

## Notes

- Should follow T087-T089.
- Do not expose these actions to client users.

## Agent Log

- 2026-04-24: Task created in the provisioning operations batch after T086.
