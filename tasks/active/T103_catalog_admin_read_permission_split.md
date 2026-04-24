# T103 - Catalog admin read permission split

Status: TODO
Owner: -
Branch: codex/t103-catalog-admin-read-permission-split
PR: -
Risk: backend/RBAC
Created: 2026-04-24
Updated: 2026-04-24

## Summary

Split admin catalog read routes from catalog mutation routes so provider readiness and catalog inspection can use read-level permission.

## Scope

- Add or reuse an admin catalog read middleware option.
- Keep catalog create/update routes on `catalog.manage`.
- Move admin catalog GET routes and provider readiness to `catalog.view`.
- Update focused tests for middleware selection.
- Update API docs and permission notes.
- Keep each file under 500 lines.

## Acceptance Criteria

- Admin catalog read routes require `catalog.view`.
- Admin catalog mutation routes still require `catalog.manage`.
- Existing reseller/client catalog access is unchanged.
- Backend and frontend validation commands pass.

## Notes

- Follows T100.
- Keep route behavior and response bodies unchanged.

## Agent Log

- 2026-04-24: Task created in the provider readiness follow-up batch.
