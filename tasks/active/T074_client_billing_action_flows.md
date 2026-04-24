# T074 - Client billing action flows

Status: IN_PROGRESS
Owner: Codex
Branch: codex/t074-client-billing-action-flows
PR: -
Risk: frontend/client
Created: 2026-04-24
Updated: 2026-04-24

## Summary

Add basic client-side billing actions for ordering, requesting top-ups, and paying invoices from wallet balance.

## Scope

- Work mainly in `frontend/src/lib/api/**/*` and client portal screens under `frontend/src/modules/client/screens/`.
- Add API wrappers for existing endpoints:
  - `POST /client/orders`
  - `POST /client/topup-requests`
  - `POST /client/invoice-wallet-payments`
- Add idempotency-key support to the shared API client for mutation calls that require it.
- Wire small forms/actions into `ClientShop`, `ClientWallet`, and invoice/payment UI where the existing screen layout supports it.
- Refresh related live data after successful actions.

## Acceptance Criteria

- Client can create an order from an available catalog plan without editing UUIDs manually.
- Client can submit a top-up request with amount, currency, method, and reference.
- Client can pay a payable invoice from wallet balance through the existing wallet payment endpoint.
- Mutation errors are shown clearly and do not silently fall back to mock success.
- `npm audit --omit=dev`, `npm run lint`, and `npm run build` pass in `frontend`.

## Notes

- The backend requires idempotency keys for order/top-up/payment mutations; generate stable per-submit keys in the UI layer.
- Keep the language in forms simple: amount, method, reference, pay, order.
- Keep files under 500 lines; split small form components if a screen grows too much.

## Agent Log

- 2026-04-24: Task created after closing stale PR #80 and refreshing the board for the next live workflow batch.
- 2026-04-24: Codex claimed the task after T073 completed and started wiring client billing mutations.
- 2026-04-24: Added shared idempotency-key support, client order/top-up/wallet-payment API wrappers, live order/pay actions in ClientShop, and a real top-up form in ClientWallet.
- 2026-04-24: Validation passed: `npm audit --omit=dev`, `npm run lint`, `npm run build`.
