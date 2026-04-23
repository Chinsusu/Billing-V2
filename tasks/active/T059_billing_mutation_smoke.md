# T059 - Billing mutation smoke

Status: TODO
Owner: -
Branch: test/billing-mutation-smoke
PR: -
Risk: test/API/audit
Created: 2026-04-23
Updated: 2026-04-23

## Summary

Add repeatable integration smoke coverage for billing mutations and audit visibility.

## Scope

- Extend or add a smoke command that performs safe seeded mutations.
- Cover top-up approve or reject, invoice wallet payment, and order status transition.
- Verify resulting audit records are visible through the audit read API.
- Keep the smoke idempotent or isolate it in a fresh test database.

## Acceptance Criteria

- Smoke can run repeatedly in local dev without corrupting shared seed data.
- Mutations verify expected status/display IDs after completion.
- Audit list/detail checks confirm tenant-scoped mutation logs.
- Backend quality gates pass.

## Notes

- Prefer fresh database setup in smoke scripts over mutating a long-lived dev DB.

## Agent Log

- 2026-04-23: Task created after audit mutation events landed.
