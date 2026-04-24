# T063 - Frontend API surface expansion

Status: TODO
Owner: -
Branch: feat/frontend-api-surface-expansion
PR: -
Risk: frontend/API
Created: 2026-04-24
Updated: 2026-04-24

## Summary

Expand the shared frontend API layer for the next portal batch without mixing screen work into the same PR.

## Scope

- Extend `frontend/src/lib/api/config.ts` to support the reseller actor cleanly.
- Extend `frontend/src/lib/api/types.ts` with typed models and query objects for data already exposed by the backend.
- Extend `frontend/src/lib/api/billing.ts` with wrappers for the current backend endpoints needed by the next UI tasks.
- Keep edits inside `frontend/src/lib/api/**/*` unless a small compile fix is unavoidable.
- Do not redesign portal screens in this task.

## Acceptance Criteria

- Shared API helpers support `admin`, `reseller`, and `client` actors.
- Typed wrappers exist for already available backend routes that the next UI batch needs, including:
  - `/admin/orders`
  - `/admin/services`
  - `/admin/wallets`
  - `/admin/topup-requests`
  - `/reseller/catalog`
  - `/reseller/catalog/master-plans`
  - `/client/catalog`
- Query objects cover only real backend filters or pagination params already supported by the current API.
- `npm run lint` and `npm run build` pass in `frontend`.

## Notes

- This is the preferred base task for the rest of the portal batch.
- Do not invent new backend routes in the frontend client.
- Keep response models practical; add only fields that current UI work will consume.

## Agent Log

- 2026-04-24: Task created as the shared API foundation for the next frontend batch.
