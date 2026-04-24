# T106 - Provider readiness runbook examples

Status: REVIEW
Owner: Codex
Branch: codex/t106-provider-readiness-runbook-examples
PR: https://github.com/Chinsusu/Billing-V2/pull/240
Risk: docs/operations
Created: 2026-04-24
Updated: 2026-04-24

## Summary

Add practical provider readiness examples to the operations docs so agents know what to do for each readiness state.

## Scope

- Add short examples for `ready`, `inactive_source`, `missing_plan_source`, `unsupported_capability`, and `fake_provider_only`.
- Include recommended operator actions for each state.
- Reference display IDs and safe API checks.
- Avoid provider credentials, production DSNs, or raw provider payload examples.
- Keep docs concise and under 500 lines per edited file.

## Acceptance Criteria

- Runbook explains what each readiness state means and what action to take.
- Examples are local/sandbox friendly and do not require real provider credentials.
- Backend and frontend validation commands pass.

## Notes

- Follows T100.
- Keep this docs-only unless a small doc link update is required.

## Agent Log

- 2026-04-24: Task created in the provider readiness follow-up batch.
- 2026-04-24: Codex claimed the task on `codex/t106-provider-readiness-runbook-examples`.
- 2026-04-24: Added provider readiness examples and operator actions for `ready`, `inactive_source`, `missing_plan_source`, `unsupported_capability`, and `fake_provider_only`.
- 2026-04-24: Validation passed: `npm ci`, `npm audit --omit=dev`, `npm run lint`, `npm run build`, `go test ./...`, and `go build ./cmd/api ./cmd/migrate ./cmd/seed ./cmd/smoke ./cmd/worker`.
- 2026-04-24: Opened PR #240 and moved the task to review.
