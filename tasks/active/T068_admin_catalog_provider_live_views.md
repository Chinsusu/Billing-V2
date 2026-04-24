# T068 - Admin catalog and provider live views

Status: TODO
Owner: -
Branch: feat/admin-catalog-provider-live-views
PR: -
Risk: frontend/API
Created: 2026-04-24
Updated: 2026-04-24

## Summary

Expose admin catalog/provider-source read APIs through the frontend API client and wire the product/provider admin screens to live data.

## Scope

- Add frontend API types and client wrappers for existing admin catalog read endpoints:
  - `/admin/catalog/products`
  - `/admin/catalog/plans`
  - `/admin/catalog/provider-sources`
- Update `AdminProducts.tsx` and `AdminProviders.tsx` to prefer live read data.
- Keep explicit loading, empty, error, and demo fallback states.
- Do not add catalog/provider mutation UI unless the frontend client and screen clearly call an existing backend endpoint.

## Acceptance Criteria

- `AdminProducts` reads live products/plans from the shared API client.
- `AdminProviders` reads live provider sources from the shared API client.
- Numeric display IDs are shown where backend records expose them.
- Fallback data is visibly labeled and is not mixed with live rows after live data loads.
- `npm run lint` and `npm run build` pass in `frontend`.

## Notes

- Backend admin catalog status PATCH endpoints exist, but this task should stay read-focused unless the implementation remains small and testable.
- Split local row mappers if either screen approaches 500 lines.

## Agent Log

- 2026-04-24: Task created after T066 completed and the admin frontend batch was refreshed.
