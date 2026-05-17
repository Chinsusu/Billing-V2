# T234 - Top-up review authorization alignment

Status: TODO
Owner: -
Branch: codex/t234-topup-review-authorization
PR: -
Risk: wallet/RBAC/API/finance
Created: 2026-05-17
Updated: 2026-05-17

## Summary

Fix or deliberately document the top-up review authorization path so approved dev/staging operators can approve or reject top-up requests through the API.

## Scope

- Inspect wallet top-up review routes, route middleware, tenant context, RBAC permissions, and existing tests.
- Decide whether reseller owners should approve tenant top-ups, whether platform admins need an emergency target-tenant context path, or both.
- Implement the smallest safe route/middleware/test change for the intended product behavior.
- Do not bypass wallet ledger append-only rules or weaken RBAC broadly.

## Acceptance Criteria

- The intended operator role can approve and reject a top-up in dev/staging through HTTP API.
- Disallowed actor types and cross-tenant requests still fail.
- Tests cover allowed and denied review paths at the route/middleware level.
- Required wallet/RBAC validation commands pass.

## Notes

- During T233 target activation, `/admin/topup-requests/:id/approve` returned `auth.permission_denied` for the demo reseller owner because the route uses admin review middleware restricted to platform actor types.
- Retrying with platform admin headers returned `tenant.context_mismatch` because the generic tenant header middleware does not build a platform emergency target-tenant context for this route.
- T233 activation used existing dev wallet balance instead of creating a new approved top-up, so provider lifecycle evidence is valid but top-up review remains unproven for launch.

## Agent Log

- 2026-05-17: Task created from T233 target activation residual risk.
