# T065 - Reseller portal dedicated screens

Status: REVIEW
Owner: Codex
Branch: codex/t065-reseller-portal-screens
PR: https://github.com/Chinsusu/Billing-V2/pull/149
Risk: frontend/reseller
Created: 2026-04-24
Updated: 2026-04-24

## Summary

Replace reseller portal screen reuse with dedicated views and wire the catalog screens to real reseller catalog endpoints.

## Scope

- Work only inside `frontend/src/modules/reseller/**/*` unless a small wiring fix is required after rebasing on T063.
- Replace reused mappings in `ResellerPortal.tsx` so accounts, services, invoices, transactions, reports, and settings each have their own screen component.
- Use live reseller catalog data for product/pricing screens when the backend is available.
- For reseller screens that still have no backend route in `main`, use module-local demo adapters or clear read-only placeholders instead of reusing the wrong screen.
- Do not redesign shared app shell or global navigation in this task.

## Acceptance Criteria

- `reseller-accounts`, `reseller-tickets`, `reseller-services-proxies`, `reseller-services-vps`, `reseller-services-bandwidth`, `reseller-invoices`, `reseller-transactions`, `reseller-products`, `reseller-reports`, and `reseller-settings` each render through a dedicated screen component.
- `ResellerPortal.tsx` no longer maps unrelated nav IDs to `ResellerClients`, `ResellerWallet`, or `ResellerDashboard`.
- `ResellerCatalog` or its replacement reads from the shared API client for `/reseller/catalog` or `/reseller/catalog/master-plans` once T063 is merged.
- Missing backend routes are handled honestly with module-local placeholder states, not by wiring the wrong data source.
- `npm run lint` and `npm run build` pass in `frontend`.

## Notes

- Rebase on top of T063 if that task lands first.
- Avoid editing `frontend/src/lib/api/**/*` in this task; consume shared helpers from `main`.
- Keep new reseller screen files focused; split local cards/tables when needed to stay under 500 lines.

## Agent Log

- 2026-04-24: Task created for reseller portal completion after the current frontend shell and admin filter work.
- 2026-04-24: Codex claimed the task after T064 merged and started splitting reseller portal screens.
- 2026-04-24: PR #149 opened with dedicated reseller service, ticket, billing, report screens and live reseller catalog API wiring.
- 2026-04-24: Validation passed in `frontend`: `npm audit --omit=dev`, `npm run lint`, and `npm run build`.
