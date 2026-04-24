# T064 - Client portal dedicated screens

Status: DONE
Owner: Codex
Branch: codex/t064-client-portal-screens
PR: https://github.com/Chinsusu/Billing-V2/pull/147
Risk: frontend/client
Created: 2026-04-24
Updated: 2026-04-24

## Summary

Replace client portal screen reuse with dedicated screens so each navigation item maps to the correct view and live billing data.

## Scope

- Work only inside `frontend/src/modules/client/**/*` unless a small wiring fix is required after rebasing on T063.
- Replace reused mappings in `ClientPortal.tsx` so invoices, transactions, and each service category have dedicated screen components.
- Split invoice and wallet/transaction concerns instead of reusing `ClientShop` and `ClientWallet` for unrelated pages.
- Keep support tickets as a clearly labeled placeholder only if there is still no backend route for that feature.
- Keep fallback demo data only where the backend surface does not exist yet.

## Acceptance Criteria

- `client-services-proxies`, `client-services-vps`, `client-services-bandwidth`, `client-invoices`, `client-transactions`, and `client-settings` each render through a dedicated screen component.
- Client invoices and transactions use the live frontend API client when the backend is available.
- Service category screens clearly scope the displayed service rows instead of reusing the overview screen.
- `ClientPortal.tsx` no longer maps unrelated nav IDs to `ClientDashboard`, `ClientShop`, or `ClientWallet`.
- `npm run lint` and `npm run build` pass in `frontend`.

## Notes

- Rebase on top of T063 if that task lands first.
- Avoid editing `frontend/src/lib/api/**/*` in this task; consume shared helpers from `main`.
- Keep files under 500 lines by splitting new screens or local view helpers.

## Agent Log

- 2026-04-24: Task created for client portal screen cleanup after the shared API integration batch.
- 2026-04-24: Codex claimed the task after T063 merged and started splitting client portal screens.
- 2026-04-24: PR #147 opened with dedicated client service, invoice, transaction, and settings screens.
- 2026-04-24: Validation passed in `frontend`: `npm audit --omit=dev`, `npm run lint`, and `npm run build`.
- 2026-04-24: PR #147 merged into `main` with commit `24432636e877da1101a1b9ccd8ba9535279245c2`.
