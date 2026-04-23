# T049 - Invoice wallet payment service

Status: TODO
Owner: -
Branch: feat/invoice-wallet-payment-service
PR: -
Risk: invoice/wallet/payment/money
Created: 2026-04-23
Updated: 2026-04-23

## Summary

Add a service that pays invoices from wallet balance and records the related ledger/payment records.

## Scope

- Add service method to pay one issued invoice from a wallet owned by the buyer.
- Debit wallet through the ledger posting service.
- Mark invoices paid when the payment covers the full total.
- Create or link payment transaction records for auditability.
- Emit an outbox event when payment succeeds.

## Acceptance Criteria

- Only tenant-scoped issued/overdue invoices can be paid.
- Wallet ownership, currency, and available balance are validated before debit.
- Duplicate payment attempts with the same idempotency key do not double debit.
- Invoice, wallet ledger, and payment transaction state stay consistent.
- Tests cover paid, duplicate, insufficient balance, currency mismatch, and cross-tenant cases.
- Full validation passes: `make fmt`, `make test`, `make build`, `make migrate-validate`.

## Notes

- This task depends on T045, T046, and T047.

## Agent Log

- 2026-04-23: Task created after invoice generation and wallet posting foundation.
