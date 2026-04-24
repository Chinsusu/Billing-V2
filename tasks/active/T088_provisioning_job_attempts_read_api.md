# T088 - Provisioning job attempts read API

Status: TODO
Owner: -
Branch: codex/t088-provisioning-job-attempts-read-api
PR: -
Risk: backend/API
Created: 2026-04-24
Updated: 2026-04-24

## Summary

Expose job attempt history so operations can understand why a provisioning job retried, failed, or moved to manual review.

## Scope

- Work mainly in `internal/modules/jobs/**/*`, `cmd/api/**/*`, and API docs.
- Add read store methods and HTTP routes for job attempts.
- Scope attempts through the parent job tenant.
- Keep error messages redacted and avoid provider credential leakage.
- Keep each file under 500 lines.

## Acceptance Criteria

- Admin can list attempts for a job in the effective tenant.
- Reseller can list attempts for a job in the effective tenant.
- Attempt responses include display ID, worker ID, attempt number, result, redacted error, duration, correlation ID, and timestamps.
- Missing or cross-tenant job IDs return clear API errors.
- Backend and frontend validation commands pass.

## Notes

- Depends on or should follow T087 route/module structure.
- Keep this read-only.

## Agent Log

- 2026-04-24: Task created in the provisioning operations batch after T086.
