# T109 - API contract drift guard

Status: DONE
Owner: Codex
Branch: codex/t109-api-contract-drift-guard
PR: https://github.com/Chinsusu/Billing-V2/pull/247
Risk: docs/API/CI
Created: 2026-04-25
Updated: 2026-04-25

## Summary

Add a lightweight guard that helps detect when backend route changes drift from the operational API reference.

## Scope

- Check that key billing routes in `cmd/api` or route wiring are represented in `docs/05_development_standards/56_Billing_API_Operational_Reference.md`.
- Focus on stable route groups, permission names, response redaction notes, and query names.
- Keep the guard simple and maintainable; avoid generating a full OpenAPI spec in this task.
- Add a documented command for agents and CI to run.
- Keep each edited file under 500 lines.

## Acceptance Criteria

- Guard fails when a tracked backend route is missing from the operational reference.
- Guard output is readable enough for another agent to fix docs quickly.
- Existing build/test commands pass.
- Docs explain how to update the guard when intentional API changes are made.

## Notes

- This task should not change API behavior.

## Agent Log

- 2026-04-25: Task created in the post-readiness hardening batch.
- 2026-04-25: Codex claimed the task; adding a lightweight route/reference drift guard.
- 2026-04-25: Added `cmd/contractguard`, `make contract-guard`, CI guard step, and docs for updating tracked API contracts. Validation passed: `go run ./cmd/contractguard`, `go test ./cmd/contractguard`, `go test ./...`, `go build ./cmd/api ./cmd/migrate ./cmd/seed ./cmd/smoke ./cmd/worker ./cmd/contractguard`. Local `make contract-guard` could not run because `make` is not installed on this Windows host; the equivalent Go command passed.
- 2026-04-25: PR #247 passed CI and merged into main at `200b9aa`; moved the task to done.
