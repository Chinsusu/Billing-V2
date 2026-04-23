# T046 - Invoice generation service

Status: REVIEW
Owner: Codex
Branch: feat/invoice-generation-service
PR: pending
Risk: invoice/order/money
Created: 2026-04-23
Updated: 2026-04-23

## Summary

Add an invoice generation service that can create an invoice from a paid order using order price snapshots.

## Scope

- Add service method to generate one invoice for one paid order.
- Use idempotency to avoid duplicate invoices for the same order/action.
- Create invoice items from available order/service data and snapshots.
- Emit an outbox event when invoice generation succeeds.
- Out of scope: PDF rendering, email delivery, gateway settlement, or frontend views.

## Acceptance Criteria

- Only tenant-scoped paid orders can generate invoices.
- Duplicate generation attempts return the existing invoice or fail with a clear conflict.
- Invoice totals match order totals and item totals.
- Domain/store/service tests cover paid, unpaid, duplicate, and cross-tenant cases.
- Full validation passes: `make fmt`, `make test`, `make build`, `make migrate-validate`.

## Notes

- This task depends on T039 and should follow T045 if it reuses invoice read store methods.

## Agent Log

- 2026-04-23: Task created for the next backend wallet/invoice batch.
- 2026-04-23: Implemented invoice generation service, order idempotency index, create store path, outbox event, and focused tests.
