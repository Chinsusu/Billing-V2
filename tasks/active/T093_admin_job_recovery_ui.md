# T093 - Admin job recovery UI actions

Status: DONE
Owner: Codex
Branch: codex/t093-admin-job-recovery-ui
PR: https://github.com/Chinsusu/Billing-V2/pull/212
Risk: frontend/admin-ops
Created: 2026-04-24
Updated: 2026-04-24

## Summary

Expose the guarded provisioning job recovery actions from T091 in the admin fulfillment UI.

## Scope

- Work mainly in `frontend/src/lib/api/**/*` and admin provisioning/fulfillment screens.
- Add admin-only actions for retry, manual review, and cancel on eligible provisioning jobs.
- Require explicit operator text for manual review and clear confirmation for retry/cancel.
- Reuse existing API helpers, status badges, and table patterns.
- Keep each frontend file under 500 lines.

## Acceptance Criteria

- Admin can retry a `failed_retryable` or `manual_review` job from the UI and see the refreshed row.
- Admin can move a safe job to manual review with a required reason.
- Admin can cancel a safe non-active job with clear confirmation.
- Reseller and client screens do not show mutation controls.
- API errors such as `job.status_conflict` and validation failures are shown in plain language.
- Frontend and backend validation commands pass.

## Notes

- Should follow T090 and T091.
- Do not add new backend recovery behavior in this task.

## Agent Log

- 2026-04-24: Task created after T092 completed and the active board was empty.
- 2026-04-24: Codex claimed the task on `codex/t093-admin-job-recovery-ui`.
- 2026-04-24: Added admin-only retry, manual-review, and cancel controls for live provisioning jobs, plus plain API error messages and local validation coverage.
- 2026-04-24: Implementation pushed for review in PR #212.
- 2026-04-24: PR https://github.com/Chinsusu/Billing-V2/pull/212 passed CI and merged to main at `29ba0e13e268e62d163d9a44d5cfa9fe16728654`.
