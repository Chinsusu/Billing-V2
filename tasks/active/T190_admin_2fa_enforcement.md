# T190 - Admin 2FA enforcement

Status: TODO
Owner: -
Branch: codex/t190-admin-2fa-enforcement
PR: -
Risk: authentication, admin security, RBAC, and audit behavior
Created: 2026-05-13
Updated: 2026-05-13

## Summary

Add admin 2FA enforcement for privileged admin access before pilot launch.

## Scope

- Add backend state and service behavior needed to require 2FA for admin users.
- Enforce 2FA for privileged admin routes without weakening tenant/RBAC checks.
- Add audit events for 2FA setup, success, and failure where appropriate.
- Add minimal frontend state/copy only if required to complete the flow.

## Acceptance Criteria

- Admin users without satisfied 2FA cannot access privileged admin actions.
- 2FA enforcement is tested for allowed and denied paths.
- Audit records do not expose secrets or raw tokens.
- Relevant backend/frontend validation and CI pass.

## Notes

- Stop and ask before selecting a 2FA method if the repo docs do not make the decision clear.

## Agent Log

- 2026-05-13: Task created by Codex backlog planning.
