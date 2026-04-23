# T043 - Wallet read API

Status: TODO
Owner: -
Branch: feat/wallet-read-api
PR: -
Risk: wallet/API/auth
Created: 2026-04-23
Updated: 2026-04-23

## Summary

Add tenant-scoped wallet and ledger read APIs so clients and admins can inspect balances and account history.

## Scope

- Add wallet service/store reads for list and detail.
- Add ledger history reads with pagination limit support.
- Add client routes for current account wallet summary and ledger history.
- Add admin routes for tenant-scoped account wallet reads.
- Out of scope: balance mutations, top-up approval, gateway calls, or frontend views.

## Acceptance Criteria

- Client reads are scoped to tenant and actor account.
- Admin reads are scoped to tenant and require wallet/account read permission.
- Responses include numeric display IDs for wallet and ledger records where applicable.
- Handler, query-builder, and runtime wiring tests cover client/admin access.
- Full validation passes: `make fmt`, `make test`, `make build`, `make migrate-validate`.

## Notes

- This task depends on T041 and T042.

## Agent Log

- 2026-04-23: Task created for the next backend wallet/invoice batch.
