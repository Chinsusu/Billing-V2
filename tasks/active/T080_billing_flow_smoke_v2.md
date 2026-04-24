# T080 - Billing flow smoke v2

Status: REVIEW
Owner: Codex
Branch: codex/t080-billing-flow-smoke-v2
PR: https://github.com/Chinsusu/Billing-V2/pull/183
Risk: QA/smoke
Created: 2026-04-24
Updated: 2026-04-24

## Summary

Add a smoke test or runbook update for the live billing flow across reseller catalog, client order/top-up/payment, and read views.

## Scope

- Work mainly in existing smoke test locations and operational docs.
- Cover the real API sequence available after T075-T079.
- Keep test data deterministic and compatible with local seed data.
- Make failures clear enough for another agent to diagnose quickly.

## Acceptance Criteria

- Smoke coverage exercises the current live flow instead of only static fixtures.
- Required seed data and environment variables are documented.
- The smoke test is included in the appropriate local validation path or has a documented command.
- Backend and frontend validation commands still pass.

## Notes

- This task should wait until the endpoint surface from T075-T079 is stable.
- Avoid brittle UI browser automation unless the API smoke path is not sufficient.

## Agent Log

- 2026-04-24: Task created after T074 completed and the board needed a follow-up validation task for the live flow.
- 2026-04-24: Codex claimed the task after T079 completed and started updating the billing mutation smoke path to use the live checkout endpoint.
- 2026-04-24: Updated `cmd/smoke dev-billing` to exercise top-up approval, order creation, `POST /client/checkouts`, duplicate checkout submit, invoice read, wallet payment, and audit verification.
- 2026-04-24: Added `make smoke-dev-billing` and documented the required local/dev seed/API environment in the local runbook.
- 2026-04-24: Validation passed: `go test ./...`, backend builds, `npm audit --omit=dev`, `npm run lint`, `npm run build`.
- 2026-04-24: Opened PR https://github.com/Chinsusu/Billing-V2/pull/183 for review and CI.
