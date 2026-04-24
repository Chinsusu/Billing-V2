# T073 - Reseller catalog clone actions

Status: TODO
Owner: -
Branch: codex/t073-reseller-catalog-clone-actions
PR: -
Risk: frontend/reseller
Created: 2026-04-24
Updated: 2026-04-24

## Summary

Add reseller catalog actions that clone master products/plans into the reseller catalog through existing backend endpoints.

## Scope

- Work mainly in `frontend/src/lib/api/**/*` and `frontend/src/modules/reseller/screens/ResellerCatalog.tsx`.
- Add API wrappers for `POST /reseller/catalog/products/clone` and `POST /reseller/catalog/plans/clone`.
- Add a simple action flow to clone a master plan/product, set selling price, visibility, and status where the endpoint requires it.
- Refresh the live catalog view after a successful clone.
- Keep inline price editing or broader catalog management out of scope unless it is already present and only needs small wiring.

## Acceptance Criteria

- Reseller catalog screen exposes a clear clone/add action for master catalog rows not already present in the reseller catalog.
- Success and error states are visible and do not require a page reload.
- API request bodies match backend field names exactly.
- Live rows use numeric display IDs where available.
- `npm audit --omit=dev`, `npm run lint`, and `npm run build` pass in `frontend`.

## Notes

- Existing backend clone endpoints are already present in the catalog module; this task should not add new backend APIs.
- If the UI needs a modal/form, reuse existing small components or keep the markup local and under 500 lines.

## Agent Log

- 2026-04-24: Task created after closing stale PR #80 and refreshing the board for the next live workflow batch.
