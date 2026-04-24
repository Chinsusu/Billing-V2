# T071 - Admin accounts read API

Status: TODO
Owner: -
Branch: codex/t071-admin-accounts-read-api
PR: -
Risk: backend/API
Created: 2026-04-24
Updated: 2026-04-24

## Summary

Add small admin read APIs for tenants and account users so admin account screens can stop using mock data.

## Scope

- Work mainly in `internal/modules/tenant`, `internal/modules/identity`, and app route wiring.
- Add list/read store methods needed for admin screens, with tenant-safe filtering and pagination where the local HTTP helpers already support it.
- Expose response fields that are useful for operations: UUID `id`, numeric `display_id`, tenant/account name, slug or email, type, status, parent/tenant relation, timestamps, and primary domain when available.
- Keep write/create account flows out of this task.
- Do not change frontend screens in this task.

## Acceptance Criteria

- Admin can list tenants through a documented backend endpoint such as `GET /admin/tenants`.
- Admin can list account users through a documented backend endpoint such as `GET /admin/accounts` or `GET /admin/customers`; choose one name and document it clearly.
- Responses include numeric `display_id` for FE display, while keeping UUIDs for internal actions.
- Unit tests cover route tenant context, filters, and response mapping.
- `make fmt`, `make test`, and `make build` pass.

## Notes

- Use plain field names and avoid leaking password hashes, token hashes, or security-only data.
- If route naming is unclear, prefer the name that maps cleanly to existing screens and record it in the task log.
- Keep each file under 500 lines; split HTTP types/handlers if needed.

## Agent Log

- 2026-04-24: Task created after closing stale PR #80 and refreshing the board for the next live workflow batch.
