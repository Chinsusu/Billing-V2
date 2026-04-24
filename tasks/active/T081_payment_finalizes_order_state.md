# T081 - Payment finalizes order state

Status: DONE
Owner: Codex
Branch: codex/t081-payment-finalizes-order-state
PR: https://github.com/Chinsusu/Billing-V2/pull/186
Risk: backend/payment/order
Created: 2026-04-24
Updated: 2026-04-24

## Summary

When a client pays an issued invoice from wallet, finalize the related order state so the order is no longer left as pending/unpaid.

## Scope

- Work mainly in `internal/modules/payment/**/*`, `internal/modules/order/**/*`, and related API docs/tests.
- Keep the invoice wallet payment API response compatible.
- Preserve idempotency when the same payment request is submitted twice.
- Keep tenant and buyer scoping enforced from existing invoice and wallet checks.

## Acceptance Criteria

- Successful wallet invoice payment marks the invoice paid and the related order paid when the invoice has `order_id`.
- Replaying the same payment idempotency key returns the same paid result without corrupting order state.
- Payment for invoices without an order remains compatible.
- Tests cover normal payment, duplicate-submit behavior, and conflict cases.
- `go test ./...` passes and backend binaries build.

## Notes

- Do not start provisioning in this task; only make billing/order state consistent.
- If a paid invoice already has an order in paid state, treat that as idempotent.

## Agent Log

- 2026-04-24: Task created after T080 showed the live checkout smoke no longer needs manual admin order status changes, but paid invoices still need order state sync.
- 2026-04-24: Codex claimed the task after the T081-T086 batch was merged and started inspecting payment/order stores for a scoped state-finalization path.
- 2026-04-24: Implemented wallet-payment order finalization in the payment transaction path, added idempotent paid-order handling, conflict mapping, API reference notes, and local validation.
- 2026-04-24: Opened PR https://github.com/Chinsusu/Billing-V2/pull/186 and moved the task to REVIEW.
- 2026-04-24: PR https://github.com/Chinsusu/Billing-V2/pull/186 passed CI and merged to main at 1a1d7de.
