# T040 - Payment transaction API

Status: TODO
Owner: -
Branch: feat/payment-transaction-api
PR: -
Risk: payment/API/DB
Created: 2026-04-23
Updated: 2026-04-23

## Summary

Add payment transaction records and read APIs so account history can show charges, refunds, and manual adjustments.

## Scope

- Add payment transaction domain models and PostgreSQL schema.
- Add store methods for create, list, and detail reads.
- Add tenant-scoped client and admin read APIs.
- Include numeric display IDs in API responses.
- Out of scope: real payment gateway calls, refunds, wallet balance math, or frontend views.

## Acceptance Criteria

- Transactions support charge, refund, and adjustment types.
- Transactions link to account, order, and invoice records when available.
- Client reads are scoped to tenant and account context.
- Admin reads are scoped to tenant context and require payment/account read permission.
- Full validation passes: `make fmt`, `make test`, `make build`, `make migrate-validate`.

## Notes

- This task should start after T039 if it links to invoices.

## Agent Log

- 2026-04-23: Task created for the next backend batch.
