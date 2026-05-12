# T200 - Notification foundation

Status: TODO
Owner: -
Branch: codex/t200-notification-foundation
PR: -
Risk: notification delivery, account security, billing operations, and secrets
Created: 2026-05-13
Updated: 2026-05-13

## Summary

Add basic notification infrastructure for launch-critical account, billing, provisioning, and support events.

## Scope

- Add a backend notification abstraction and safe local/dev implementation.
- Cover password reset, top-up status, provisioning failure/manual review, and service lifecycle notifications.
- Ensure payloads do not leak secrets, credentials, or sensitive provider details.
- Do not add broad marketing or campaign features.

## Acceptance Criteria

- Notification events can be emitted from launch-critical flows with redacted payloads.
- Tests cover event creation and secret redaction.
- Local/dev delivery is deterministic for smoke or operator inspection.
- Relevant backend validation and CI pass.

## Notes

- Coordinate with T191 and T198 for password reset and lifecycle notification use.

## Agent Log

- 2026-05-13: Task created by Codex backlog planning.
