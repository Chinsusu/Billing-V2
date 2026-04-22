# T011 - Persistence Repositories

Status: DONE
Owner: Codex
Branch: feat/persistence-repositories
PR: https://github.com/Chinsusu/Billing-V2/pull/30
Risk: tenant/RBAC/audit/DB
Created: 2026-04-23
Updated: 2026-04-23

## Summary

Add the first PostgreSQL repository layer for tenant, identity, RBAC, and audit records created by the initial migrations.

## Scope

- Add DB executor/transaction helper that store implementations can use with `*sql.DB` or `*sql.Tx`.
- Add tenant store for `tenants` and `tenant_domains`.
- Add identity user store scoped by `tenant_id`.
- Add RBAC store for user role and permission loading.
- Add audit append store for redacted audit logs.
- Keep service logic, HTTP routes, auth sessions, and seed data out of scope.

## Acceptance Criteria

- Store methods validate required tenant/user/action inputs before SQL.
- Tenant-owned reads include tenant scope where applicable.
- Audit append stores only redacted snapshots/metadata supplied by callers.
- Files stay under 500 lines.
- `make fmt`, `make test`, `make build`, and `make migrate-validate` pass.

## Notes

- This task creates persistence foundations only; later tasks should wire services and routes on top.

## Agent Log

- 2026-04-23: Codex claimed task from `origin/main` using isolated worktree `/tmp/Billing-T011`.
- 2026-04-23: Opened PR #30 after `make fmt`, `make test`, `make build`, and `make migrate-validate` passed.
- 2026-04-23: PR #30 merged; T011 marked DONE.
