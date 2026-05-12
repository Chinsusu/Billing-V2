# T198 - Lifecycle scheduler jobs

Status: TODO
Owner: -
Branch: codex/t198-lifecycle-scheduler-jobs
PR: -
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
