# T223 - Production auth UI session path

Status: IN_PROGRESS
Owner: Codex
Branch: codex/t223-production-auth-ui
PR: -
Risk: auth/session/RBAC/frontend/config
Created: 2026-05-16
Updated: 2026-05-16

## Summary

Add a production-oriented frontend login/session path so the UI can use backend `/auth/*` cookie sessions instead of relying only on the demo portal switcher and local actor headers.

## Scope

- Add frontend auth API helpers for login, logout, and TOTP setup/verify.
- Add a session-gated app shell that can run in auth mode while keeping demo portal mode available for existing smoke coverage.
- Update the frontend API client so cookie credentials are sent and local actor headers are used only when explicitly enabled.
- Make local dev seed users usable for auth smoke by replacing placeholder password hashes with dev-only Argon2id hashes.
- Out of scope: production domain/TLS, secret-store rollout, payment/provider mutating approvals, and real launch owner sign-off.

## Acceptance Criteria

- Production auth mode renders login, handles session cookies, supports required TOTP setup/verify, and logs out.
- Session-auth mode is enabled through explicit frontend configuration, while existing demo smoke mode remains available.
- Frontend API requests include credentials and omit dev actor headers unless the dev-header flag is enabled.
- Dev seed users have valid Argon2id hashes for local-only auth testing.
- Required lint/build/smoke/backend checks pass or blockers are documented.

## Notes

- External Go/No-Go items that require real owners, credentials, staging, or provider approvals must remain evidence tasks and cannot be marked complete by code-only changes.

## Agent Log

- 2026-05-16: Task created and claimed from latest `origin/main`.
