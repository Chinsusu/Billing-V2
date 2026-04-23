# T054 - Billing API integration smoke

Status: REVIEW
Owner: Codex
Branch: test/billing-api-integration-smoke
PR: https://github.com/Chinsusu/Billing-V2/pull/123
Risk: API/billing
Created: 2026-04-23
Updated: 2026-04-23

## Summary

Add a repeatable HTTP smoke path that exercises seeded billing read APIs through the running backend.

## Scope

- Start from a migrated and seeded dev database.
- Cover customer wallet, ledger, invoice, payment, and service reads.
- Cover admin order, payment reconciliation, and audit reads.
- Use the documented local actor and tenant headers.
- Keep the smoke path scriptable for local developers and agents.

## Acceptance Criteria

- A documented command verifies the key seeded billing endpoints return success responses.
- The command reports which endpoint failed and why.
- Smoke checks use local/dev data only and no real credentials.
- Backend quality gates still pass.

## Notes

- This task should follow T053 so it can rely on a known seeded DB state.

## Agent Log

- 2026-04-23: Task created for API-level validation after DB smoke.
- 2026-04-23: Claimed by Codex; adding `smoke dev-api` coverage for seeded client/admin billing reads.
- 2026-04-23: Opened PR #123. Validation passed: `make fmt`, `make test`, `make build`, `make migrate-validate`, plus real `make smoke-dev-db` and `make smoke-dev-api` against local PostgreSQL/API.
