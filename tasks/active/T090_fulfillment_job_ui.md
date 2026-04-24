# T090 - Fulfillment job UI

Status: IN_PROGRESS
Owner: Codex
Branch: codex/t090-fulfillment-job-ui
PR: -
Risk: frontend/admin-reseller
Created: 2026-04-24
Updated: 2026-04-24

## Summary

Use the provisioning job read API to show real fulfillment job state in admin/reseller screens instead of inferred job labels.

## Scope

- Work mainly in `frontend/src/lib/api/**/*` and admin/reseller fulfillment-related screens.
- Prefer shared API types and helpers over per-screen hardcoding.
- Keep row displays numeric-display-ID first.
- Show explicit fallback text when job API is unavailable or partial data fails.
- Keep each file under 500 lines.

## Acceptance Criteria

- Admin and reseller can see job display ID and status for paid orders pending provisioning.
- Screens link order display ID, service display ID when present, and job display ID when present.
- Failed, retryable, terminal, manual review, and active service states are visually distinct.
- Frontend validation commands pass; backend gates still pass.

## Notes

- Should follow T087.
- Do not add mutation controls in this task.

## Agent Log

- 2026-04-24: Task created in the provisioning operations batch after T086.
- 2026-04-24: Codex claimed the task after T089 merged and started wiring real provisioning job state into admin/reseller fulfillment UI.
- 2026-04-24: Added shared frontend job types/API helpers, live admin provisioning queue rows, reseller fulfillment job labels, pending paid-order rows, and explicit unavailable-state text when job API data is missing. Validation passed: frontend lint/build/audit plus backend test/build gates.
