# T236 - Target auth session and RBAC evidence

Status: REVIEW
Owner: Codex
Branch: codex/t236-target-auth-rbac-smoke
PR: https://github.com/Chinsusu/Billing-V2/pull/504
Risk: auth/RBAC/tenant/security/audit
Created: 2026-05-17
Updated: 2026-05-17

## Summary

Prove target-environment auth/session and RBAC negative checks with a repeatable dev/test smoke that does not touch money, provider provisioning, service lifecycle, or credentials.

## Scope

- Add a repeatable target auth/RBAC smoke command for approved dev/test environments.
- Verify a real `billing_session` cookie can access a client route without `X-Actor-*` dev helper headers.
- Verify an unsatisfied platform admin session is blocked from admin routes by 2FA enforcement.
- Verify invalid session, missing actor, missing permission, and cross-tenant mismatch errors are denied with stable envelopes.
- Run the smoke on the approved test server.
- Record redacted evidence in launch docs.

## Acceptance Criteria

- Client session login succeeds and cookie-only `/client/catalog` succeeds.
- Platform admin session without satisfied 2FA receives `auth.2fa_required` on an admin route.
- Invalid session receives `auth.session_invalid`.
- Missing actor receives `auth.actor_required`.
- Cross-tenant mismatch receives `tenant.context_mismatch`.
- Low-permission RBAC checks still receive `auth.permission_denied`.
- Evidence excludes raw session tokens, cookies, passwords, DSNs, provider payloads, and credentials.

## Notes

- This task may create dev/test auth session, login audit, and rate-limit rows only.
- Seed login uses the existing dev/test seed users and must never be run against production.

## Agent Log

- 2026-05-17: Task created and claimed by Codex on branch `codex/t236-target-auth-rbac-smoke`.
- 2026-05-17: Added `dev-target-auth-rbac` smoke command and focused unit coverage for cookie-only session access and non-leaking status errors.
- 2026-05-17: Deployed current branch to the approved test server and ran `./bin/smoke -timeout 90s dev-target-auth-rbac` with `APP_ENV=dev` and local API base URL. Result PASS: client session cookie-only access, admin 2FA gate, invalid session denial, missing actor denial, cross-tenant mismatch denial, three low-permission RBAC denials, provider mutation routes called `no`, and money mutation routes called `no`.
- 2026-05-17: Local focused validation passed: `gofmt`, `go test ./cmd/smoke`, `go run ./cmd/taskguard`, `git diff --check`.
- 2026-05-17: Opened PR #504 and moved task to `REVIEW`.
