# T050 - Payment reconciliation read API

Status: DONE
Owner: Codex
Branch: feat/payment-reconciliation-read-api
PR: https://github.com/Chinsusu/Billing-V2/pull/114
Risk: payment/API/admin
Created: 2026-04-23
Updated: 2026-04-23

## Summary

Add read APIs that let admins inspect payment transactions, wallet ledger links, and invoice payment status together.

## Scope

- Add tenant-scoped admin reconciliation list/detail endpoints.
- Join payment transactions with invoice and wallet ledger references where available.
- Support filters for status, provider, invoice, wallet, and created time window.
- Include numeric display IDs for all visible entities.
- Keep mutation and settlement logic out of scope.

## Acceptance Criteria

- Admin reads are tenant-scoped and permission protected.
- Response clearly shows invoice, transaction, and ledger identifiers when linked.
- Filters are validated and covered by query-builder tests.
- Handler/runtime tests cover admin access and validation errors.
- Full validation passes: `make fmt`, `make test`, `make build`, `make migrate-validate`.

## Notes

- This task can follow T049 or be done independently for current payment records.

## Agent Log

- 2026-04-23: Task created for operations visibility.
- 2026-04-23: Added admin payment reconciliation read API with filters and runtime coverage.
- 2026-04-23: PR #114 merged after GitHub checks passed.
