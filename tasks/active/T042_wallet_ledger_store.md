# T042 - Wallet ledger store

Status: REVIEW
Owner: Codex
Branch: feat/wallet-ledger-store
PR: https://github.com/Chinsusu/Billing-V2/pull/97
Risk: wallet/ledger/DB
Created: 2026-04-23
Updated: 2026-04-23

## Summary

Add immutable wallet ledger models, schema, and store methods so balance movements can be recorded without mutating history.

## Scope

- Add ledger entry domain models for credit, debit, adjustment, refund, reversal, and purchase movements.
- Add PostgreSQL schema for wallet ledger entries with idempotency constraints.
- Add store methods for create and account-scoped reads.
- Keep wallet balance cache changes out of scope unless needed for transaction safety.
- Out of scope: top-up approval workflow, gateway integration, or service provisioning integration.

## Acceptance Criteria

- Ledger entries are append-only and link to wallet, tenant, reference type/id, and idempotency key.
- Posted entries cannot be changed by store methods.
- Store reads are tenant-scoped and wallet-scoped.
- Migration validation and focused store/domain tests pass.
- Full validation passes: `make fmt`, `make test`, `make build`, `make migrate-validate`.

## Notes

- This task depends on T041 wallet schema.

## Agent Log

- 2026-04-23: Task created for the next backend wallet/invoice batch.
- 2026-04-23: Claimed by Codex from latest `origin/main` in `/tmp/Billing-T042`.
- 2026-04-23: Added wallet ledger domain models, schema, store methods, and scoped read query builders. Validation passed: `make fmt`, `make test`, `make build`, `make migrate-validate`.
- 2026-04-23: Opened PR #97.
