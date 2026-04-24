# T067 - Admin service inventory live views

Status: REVIEW
Owner: Codex
Branch: codex/t067-admin-service-inventory
PR: -
Risk: frontend/admin
Created: 2026-04-24
Updated: 2026-04-24

## Summary

Replace the admin service inventory demo tables with live `/admin/services` read models while keeping clear fallback states when the backend is unavailable.

## Scope

- Work inside `frontend/src/modules/admin/**/*` unless a small shared frontend helper is needed.
- Update `AdminServicesProxies.tsx`, `AdminServicesVPS.tsx`, and `AdminServicesBandwidth.tsx` to consume the shared frontend API client.
- Keep each service screen specific to its product family by deriving labels/specs from live service data where practical.
- Preserve loading, empty, error, and demo fallback states.
- Do not add service write actions in this task.

## Acceptance Criteria

- Each admin service inventory screen reads from `billingApi.listAdminServices`.
- Live rows use numeric display IDs (`SVC-<display_id>`) instead of UUIDs.
- The UI clearly labels when it is showing demo fallback data.
- No unrelated mock service data is shown after live data loads successfully.
- `npm run lint` and `npm run build` pass in `frontend`.

## Notes

- Keep files under 500 lines. Create module-local helpers if the mapping logic starts to grow.
- Product type may need to be inferred from service snapshots until a stronger backend field exists.

## Agent Log

- 2026-04-24: Task created after T066 completed and the admin frontend batch was refreshed.
- 2026-04-24: Codex claimed the task and started wiring admin service inventory screens to live API data.
- 2026-04-24: Admin proxy, VPS, and bandwidth service screens now use live admin service reads with explicit demo fallback states. Local audit, lint, and build passed.
