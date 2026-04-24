# T092 - Billing smoke worker fulfillment

Status: TODO
Owner: -
Branch: codex/t092-billing-smoke-worker-fulfillment
PR: -
Risk: QA/smoke
Created: 2026-04-24
Updated: 2026-04-24

## Summary

Extend billing smoke so it can process a paid order provisioning job through the local worker and verify the resulting service state.

## Scope

- Work mainly in `cmd/smoke/**/*`, worker command docs, and local runbooks.
- Reuse the existing `dev-billing` flow and worker `run-once` command.
- Keep the smoke deterministic and fake-provider based.
- Keep each file under 500 lines.

## Acceptance Criteria

- Smoke can create/pay an order, run the provisioning worker once, and verify service visibility.
- Smoke reports clear failure messages for missing job, failed job, missing service, or wrong order linkage.
- Smoke documents required local seed/API/worker environment.
- Backend and frontend validation commands pass.

## Notes

- Should follow T089 and can use T087/T088 APIs if available.
- Avoid real provider credentials.

## Agent Log

- 2026-04-24: Task created in the provisioning operations batch after T086.
