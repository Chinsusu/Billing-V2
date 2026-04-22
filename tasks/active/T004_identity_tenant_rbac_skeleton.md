# T004 - Identity Tenant RBAC Skeleton

Status: TODO
Owner: -
Branch: feat/identity-tenant-rbac-skeleton
PR: -
Risk: tenant/RBAC
Created: 2026-04-22
Updated: 2026-04-22

## Summary

Add skeleton interfaces/types for identity, tenant context, and RBAC checks without persistence.

## Scope

- Define identity/user context types.
- Define tenant context helpers.
- Define RBAC check interfaces.
- Avoid database persistence in this task.

## Acceptance Criteria

- Types are small and placed in clear owner modules.
- Tenant context cannot be silently absent in protected flows.
- RBAC interfaces are testable without real persistence.
- `make test` passes.
- `make build` passes.

## Notes

- Follow tenant/RBAC docs before implementation.
- Do not add real auth or production session handling in this task.

## Agent Log

- 2026-04-22: Task file created from `TASKS.md`.
