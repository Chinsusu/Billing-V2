# T053 - Dev DB smoke

Status: IN_PROGRESS
Owner: Codex
Branch: test/dev-db-smoke
PR: -
Risk: DB/seed/dev
Created: 2026-04-23
Updated: 2026-04-23

## Summary

Add a repeatable dev database smoke command that runs migrations, applies dev seed data, and checks the seeded billing flow on a real PostgreSQL database.

## Scope

- Add a command or Make target that requires `DB_DSN` and uses the real PostgreSQL driver.
- Apply all migrations before seed data.
- Apply dev seed data idempotently.
- Verify core seeded records exist for tenant, catalog, order, service, wallet, ledger, top-up, invoice, payment, and audit-readable tables where applicable.
- Keep the smoke command safe for local/dev databases only and document the expected environment.

## Acceptance Criteria

- `make smoke-dev-db` or an equivalent documented command runs migrations and seed against `DB_DSN`.
- Running the command twice does not create duplicate logical seed records.
- The command fails with a clear error if required seeded billing records are missing.
- `make fmt`, `make test`, `make build`, and `make migrate-validate` pass.
- Real DB validation result is recorded in the task log or PR body.

## Notes

- Use a fresh local/dev database when possible.
- Do not embed secrets, production DSNs, or real customer data.
- If PostgreSQL is not available on the worker host, still add the command and record the environment blocker clearly.

## Agent Log

- 2026-04-23: Task created after billing seed flow landed.
- 2026-04-23: Claimed by Codex; adding repeatable dev DB smoke command and Make target.
