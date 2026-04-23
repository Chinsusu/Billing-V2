# T041 - Wallet domain schema

Status: DONE
Owner: Codex
Branch: feat/wallet-domain-schema
PR: https://github.com/Chinsusu/Billing-V2/pull/95
Risk: wallet/migration
Created: 2026-04-23
Updated: 2026-04-23

## Summary

Add wallet domain models and PostgreSQL schema so account balances have a stable tenant-scoped storage base.

## Scope

- Add wallet owner/type/status domain models with validation.
- Add PostgreSQL schema for wallets with UUID primary keys and numeric display IDs.
- Add currency and balance cache constraints.
- Add rollback artifact and focused domain/migration validation.
- Out of scope: ledger posting, top-up workflow, payment gateway calls, or frontend views.

## Acceptance Criteria

- Wallets support tenant, user, reseller settlement, and platform owner types.
- Wallet rows include tenant, owner, currency, status, available balance cache, and locked balance cache.
- Numeric display IDs are generated for account-facing wallet records.
- Migration validation passes and rollback notes are clear.
- Full validation passes: `make fmt`, `make test`, `make build`, `make migrate-validate`.

## Notes

- This task should start from latest `origin/main` after T040.

## Agent Log

- 2026-04-23: Task created for the next backend wallet/invoice batch.
- 2026-04-23: Claimed by Codex from latest `origin/main` in `/tmp/Billing-T041`.
- 2026-04-23: Added wallet domain models, wallet schema, and rollback script. Validation passed: `make fmt`, `make test`, `make build`, `make migrate-validate`.
- 2026-04-23: Opened PR #95.
- 2026-04-23: PR #95 merged into `main`.
