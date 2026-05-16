# T224 - Auth forwarded host domain resolution

Status: IN_PROGRESS
Owner: Codex
Branch: codex/t224-auth-forwarded-host
PR: -
Risk: auth/session/tenant-domain/reverse-proxy
Created: 2026-05-16
Updated: 2026-05-16

## Summary

Fix domain auth behind the frontend reverse proxy so `/auth/login` and password reset tenant resolution can use the public forwarded host instead of only the backend request host.

## Scope

- Prefer `X-Forwarded-Host` or RFC `Forwarded` host for auth domain resolution when present.
- Keep direct backend `Host` behavior as fallback.
- Add unit tests for forwarded-host login and password reset request handling.
- Validate against the live staging tunnel after merge/deploy.

## Acceptance Criteria

- `/backend/auth/login` through `billing.resvn.net` resolves the mapped tenant without local tenant headers.
- Auth and password reset handlers pass forwarded host into the auth service.
- Existing direct-host behavior still works.
- Required backend tests and taskguard pass.

## Notes

- Found during staging E2E: direct API login with `Host: billing.resvn.net` passed, but login through the Next `/backend/*` reverse proxy returned `tenant.context_missing`.

## Agent Log

- 2026-05-16: Task created and claimed from latest `origin/main`.
