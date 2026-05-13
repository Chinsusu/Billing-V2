# T201 - Support and abuse basic backend

Status: REVIEW
Owner: Codex
Branch: codex/t201-support-abuse-basic-backend
PR: https://github.com/Chinsusu/Billing-V2/pull/433
Risk: support operations, abuse workflow, tenant isolation, service suspension, and audit
Created: 2026-05-13
Updated: 2026-05-13

## Summary

Add basic backend support and abuse control records needed for MVP operations.

## Scope

- Add minimal support ticket or support case records if current backend does not provide them.
- Add basic abuse flag/case workflow with evidence notes and service/account references.
- Enforce tenant/RBAC access and audit all sensitive support/abuse actions.
- Do not build a full custom ticket system beyond MVP needs.

## Acceptance Criteria

- Admin/reseller/client access follows tenant and permission boundaries.
- Abuse actions can record reason/evidence and trigger supported manual suspension path when applicable.
- Tests cover allowed/denied access and audit behavior.
- Relevant backend validation and CI pass.

## Notes

- Stop and ask if support data retention or abuse takedown policy is unclear.

## Agent Log

- 2026-05-13: Task created by Codex backlog planning.
- 2026-05-13: Codex claimed task on `codex/t201-support-abuse-basic-backend`.
- 2026-05-13: Added support ticket/note, risk flag, and abuse case backend records with tenant/RBAC service checks, redacted audit metadata, and abuse-driven service suspension hook.
- 2026-05-13: Local validation passed: `make fmt`, `go test ./internal/modules/support`, `make migrate-validate`, `make test`, `make build`, `make contract-guard`, `make error-code-guard`, `make task-guard`, and `git diff --check`; `make smoke-dev-db` was blocked because `DB_DSN` was not configured.
- 2026-05-13: Opened PR https://github.com/Chinsusu/Billing-V2/pull/433 for review.
