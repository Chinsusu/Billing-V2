# T189 - Auth and session baseline

Status: TODO
Owner: -
Branch: codex/t189-auth-session-baseline
PR: -
Risk: authentication, tenant isolation, RBAC, and public API behavior
Created: 2026-05-13
Updated: 2026-05-13

## Summary

Add a real authentication/session baseline so runtime identity no longer depends only on demo actor headers for production paths.

## Scope

- Define and implement login/session primitives in the backend using existing identity, tenant, and RBAC boundaries.
- Add secure session expiration behavior and tests for missing/invalid sessions.
- Keep demo/local actor-header behavior explicitly dev-only if it remains necessary for smoke tests.
- Do not add 2FA in this task; T190 owns that.

## Acceptance Criteria

- Authenticated requests receive stable actor and tenant context from a session mechanism.
- Missing, expired, or invalid sessions are rejected with standard API errors.
- Existing tenant/RBAC route protections still pass.
- Go tests, build, contract/error guards as applicable, and CI pass.

## Notes

- Stop and ask before choosing production session storage or cookie policy if the codebase does not answer it safely.

## Agent Log

- 2026-05-13: Task created by Codex backlog planning.
