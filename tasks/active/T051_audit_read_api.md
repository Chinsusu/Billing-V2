# T051 - Audit read API

Status: TODO
Owner: -
Branch: feat/audit-read-api
PR: -
Risk: audit/API/admin
Created: 2026-04-23
Updated: 2026-04-23

## Summary

Add tenant-scoped audit log read endpoints for admin troubleshooting and compliance checks.

## Scope

- Add audit store list/detail methods for existing audit records.
- Add admin endpoints for audit list and detail.
- Support filters for actor, action, entity type, entity id, and time window.
- Include numeric display IDs when audit rows have them.
- Keep export/download out of scope.

## Acceptance Criteria

- Audit reads are tenant-scoped and protected by audit view permission.
- Filters are validated and avoid unbounded large reads.
- Responses avoid exposing sensitive payload fields by default.
- Query-builder, handler, and runtime wiring tests pass.
- Full validation passes: `make fmt`, `make test`, `make build`, `make migrate-validate`.

## Notes

- This task should reuse existing HTTP pagination and RBAC middleware.

## Agent Log

- 2026-04-23: Task created for admin observability.
