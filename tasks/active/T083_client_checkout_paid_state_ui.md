# T083 - Client checkout paid state UI

Status: TODO
Owner: -
Branch: codex/t083-client-checkout-paid-state-ui
PR: -
Risk: frontend/client
Created: 2026-04-24
Updated: 2026-04-24

## Summary

Polish the client checkout UI around order, checkout invoice, wallet payment, and paid/pending-provisioning states.

## Scope

- Work mainly in `frontend/src/modules/client/**/*` and frontend API wrappers if needed.
- Keep the flow on live APIs introduced through T074-T081.
- Avoid a new design system; follow the current shell and component style.
- Keep demo fallback explicit where API data is unavailable.

## Acceptance Criteria

- Client can see a clear order -> invoice -> payment state after ordering a plan.
- Wallet payment refreshes invoices, transactions, wallet balance, and order status without a manual reload.
- Numeric display IDs are shown for orders, invoices, transactions, and services.
- Errors for checkout/payment are visible and do not look like success states.
- `npm audit --omit=dev`, `npm run lint`, and `npm run build` pass.

## Notes

- This task should wait until T081 is merged if it needs paid order status from the API.
- Do not add browser automation unless a simple component/API flow check is not enough.

## Agent Log

- 2026-04-24: Task created after T079 added checkout API wiring and T081 was planned for paid order state consistency.

