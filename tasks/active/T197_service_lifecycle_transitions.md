# T197 - Service lifecycle transitions

Status: IN_PROGRESS
Owner: Codex
Branch: codex/t197-service-lifecycle-transitions
PR: -
Risk: service lifecycle, billing, provisioning, tenant isolation, and audit
Created: 2026-05-13
Updated: 2026-05-13

## Summary

Implement or complete service lifecycle transitions for renew, expire, grace, suspend, unsuspend, and terminate.

## Scope

- Define lifecycle transition rules from existing docs and code.
- Add service methods and API behavior for supported manual transitions.
- Record lifecycle events and audit entries where required.
- Do not add scheduler automation in this task; T198 owns recurring jobs.

## Acceptance Criteria

- Supported transitions are deterministic and reject invalid state changes.
- Renew uses correct term calculation and billing state.
- Tests cover allowed and denied transitions, tenant/RBAC checks, and audit behavior.
- Relevant backend validation and CI pass.

## Notes

- Stop and ask before changing money-impacting renewal or refund behavior if policy is unclear.

## Agent Log

- 2026-05-13: Implemented service lifecycle transition primitives, admin/reseller suspend-unsuspend-terminate APIs, service lifecycle events/audit, RBAC permission seed/migration, and API/error docs. Local checks so far: `make fmt`, targeted `go test ./internal/modules/order ./cmd/contractguard ./cmd/errorcodeguard ./internal/seed`, `make test`, `make build`, `make migrate-validate`, `make contract-guard`, `make error-code-guard`, `make task-guard`, `git diff --check`. `make smoke-dev-billing` was blocked because `DB_DSN`/`-dsn` is not configured.
- 2026-05-13: Claimed by Codex on `codex/t197-service-lifecycle-transitions`; starting lifecycle/API/doc review before implementation.
- 2026-05-13: Task created by Codex backlog planning.
