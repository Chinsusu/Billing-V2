# T118 - Smoke runbook command matrix

Status: IN_PROGRESS
Owner: Codex
Branch: codex/t118-smoke-runbook-command-matrix
PR: -
Risk: docs/workflow
Created: 2026-04-25
Updated: 2026-04-25

## Summary

Unify smoke test and validation command documentation so agents know which checks to run for each change type.

## Scope

- Review existing backend, frontend, contract, and smoke documentation.
- Add or update one concise command matrix for local and CI validation.
- Call out Windows local equivalents when `make` is not available.
- Avoid duplicating long command explanations across many docs.
- Keep each edited file under 500 lines.

## Acceptance Criteria

- Agents can find the right validation commands for docs-only, backend, frontend, DB, provider, and full-stack changes.
- The matrix mentions command ordering when commands share build output.
- Links to related detailed docs are present.
- Existing docs guard commands pass.

## Notes

- This task should reduce confusion, not add a second process.

## Agent Log

- 2026-04-25: Task created in the board and delivery hardening batch.
- 2026-04-25: Codex claimed the task; consolidating smoke and validation commands into one command matrix.
- 2026-04-25: Added the validation command matrix, linked related docs, and corrected stale build command examples; validation passed for task, contract, error-code guards, Go tests, and diff check.
