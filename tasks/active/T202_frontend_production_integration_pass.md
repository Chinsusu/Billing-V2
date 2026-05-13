# T202 - Frontend production integration pass

Status: IN_PROGRESS
Owner: Codex
Branch: codex/t202-frontend-production-integration-pass
PR: -
Risk: frontend API integration, tenant/RBAC visibility, credentials, and billing actions
Created: 2026-05-13
Updated: 2026-05-13

## Summary

Replace remaining critical demo-only frontend behavior with production-safe API integration paths.

## Scope

- Audit admin, reseller, and client screens for fallback-only critical actions.
- Wire critical launch actions to backend APIs where the backend is ready.
- Keep safe demo fallback only for local/mock states and label it clearly.
- Ensure unauthorized actions are hidden in UI but still blocked by API.

## Acceptance Criteria

- Critical portal flows use backend data/actions instead of static mocks where launch requires it.
- Empty, loading, and error states are acceptable for pilot.
- Sensitive text guard, lint, build, admin smoke, and CI pass.
- Any backend gaps found are linked to existing or new tasks rather than hidden in UI.

## Notes

- Do not wire credential reveal until T193 is complete.

## Agent Log

- 2026-05-13: Task created by Codex backlog planning.
- 2026-05-13: Codex claimed task on `codex/t202-frontend-production-integration-pass`.
- 2026-05-13: Exposed API-backed client shop and wallet screens in navigation, routed dashboard CTAs to those production paths, and replaced the unsupported direct renew CTA with a safe checkout/support path.
- 2026-05-13: Created follow-up T206 for the missing direct client service renewal API/UI action instead of wiring a fake renewal action.
- 2026-05-13: Local validation passed: `npm --prefix frontend run lint`, `npm --prefix frontend run check:sensitive-text`, `npm --prefix frontend run build`, `npm --prefix frontend run smoke:admin:ci`, `make task-guard`, and `git diff --check`.
