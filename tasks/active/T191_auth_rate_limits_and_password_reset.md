# T191 - Auth rate limits and password reset

Status: IN_PROGRESS
Owner: Codex
Branch: codex/t191-auth-rate-limits-password-reset
PR: -
Risk: authentication, account security, rate limiting, and notification boundaries
Created: 2026-05-13
Updated: 2026-05-13

## Summary

Add launch-required login protection and password reset primitives.

## Scope

- Add or wire rate limits for login and password reset attempts.
- Implement secure password reset token lifecycle and tests.
- Ensure reset tokens are never logged or returned after creation.
- Do not implement broad notification delivery beyond what is required to hand off reset delivery safely.

## Acceptance Criteria

- Login and reset endpoints reject excessive attempts deterministically.
- Password reset tokens expire, are single-use, and are redacted from logs/audit.
- Tests cover success, expired token, replay, and rate-limit denial.
- Relevant backend validation and CI pass.

## Notes

- Coordinate with T200 if reset delivery needs notification infrastructure.

## Agent Log

- 2026-05-13: Task created by Codex backlog planning.
- 2026-05-13: Claimed by Codex on branch `codex/t191-auth-rate-limits-password-reset`.
- 2026-05-13: Implemented DB-backed login/password-reset rate limits, hashed single-use password reset tokens, reset confirm password update with session revocation, and password reset API routes.
- 2026-05-13: Local validation passed: focused Go tests, `make test`, `make build`, migration validate, API/error guards, frontend lint/sensitive-text/audit/build/smoke, `taskguard`, `diff --check`, and secret grep. `make test` was run outside sandbox because `httptest` local sockets may be blocked inside sandbox.
