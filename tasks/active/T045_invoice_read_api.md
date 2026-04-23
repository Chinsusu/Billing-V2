# T045 - Invoice read API

Status: TODO
Owner: -
Branch: feat/invoice-read-api
PR: -
Risk: invoice/API/auth
Created: 2026-04-23
Updated: 2026-04-23

## Summary

Add invoice and invoice item read APIs so customers and admins can inspect generated invoice records.

## Scope

- Add invoice store methods for list and detail reads including invoice items.
- Add client routes for current account invoice list and detail.
- Add admin routes for tenant-scoped invoice list and detail.
- Include numeric invoice display IDs in responses.
- Out of scope: invoice generation, PDF export, payment matching, or frontend views.

## Acceptance Criteria

- Client reads are scoped to tenant and actor account.
- Admin reads are scoped to tenant and require wallet/payment/account read permission.
- Invoice detail includes line items and optional order references.
- Handler, query-builder, and runtime wiring tests cover client/admin access.
- Full validation passes: `make fmt`, `make test`, `make build`, `make migrate-validate`.

## Notes

- This task depends on T039 invoice schema.

## Agent Log

- 2026-04-23: Task created for the next backend wallet/invoice batch.
