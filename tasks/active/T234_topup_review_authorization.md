# T234 - Top-up review authorization alignment

Status: DONE
Owner: Codex
Branch: codex/t234-topup-review-authorization
PR: https://github.com/Chinsusu/Billing-V2/pull/500
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
- Decision for this task: implement the normal reseller review path for client top-ups. Platform-admin emergency target-tenant review remains separate from this change because it needs explicit target context, reason, 2FA, and audit policy.

## Agent Log

- 2026-05-17: Task created from T233 target activation residual risk.
- 2026-05-17: Claimed by Codex on branch `codex/t234-topup-review-authorization`; inspecting wallet top-up review RBAC/tenant path before coding.
- 2026-05-17: Added reseller top-up approve/reject route wiring, route/middleware tests, runtime protection test, smoke update, and API contract docs.
- 2026-05-17: Opened PR #500 for review after local validation: focused tests, `make fmt`, `make test`, `make build`, contract guard, error-code guard, taskguard, and `git diff --check` passed.
- 2026-05-17: PR #500 merged into `main`; marking task DONE with marker branch `codex/t234-done-marker`.
