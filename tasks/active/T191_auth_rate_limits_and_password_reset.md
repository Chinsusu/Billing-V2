# T191 - Auth rate limits and password reset

Status: TODO
Owner: -
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
