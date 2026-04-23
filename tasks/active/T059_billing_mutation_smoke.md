# T059 - Billing mutation smoke

Status: DONE
Owner: Codex
Branch: feat/billing-mutation-smoke
PR: #134
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
- 2026-04-23: Codex started the repeatable billing mutation smoke flow.
- 2026-04-23: Added `dev-billing` smoke command for top-up approval, order payment, invoice generation, wallet payment, and audit verification.
- 2026-04-23: Fixed invoice audit payload SQL casts so smoke invoice creation and payment writes work on PostgreSQL.
- 2026-04-23: Verified smoke on a fresh `billing_smoke` database and passed `make fmt`, `make test`, `make build`, and `make migrate-validate`.
