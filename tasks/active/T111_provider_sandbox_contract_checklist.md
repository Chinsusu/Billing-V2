# T111 - Provider sandbox contract checklist

Status: REVIEW
Owner: Codex
Branch: codex/t111-provider-sandbox-contract-checklist
PR: pending
Risk: docs/provider-ops
Created: 2026-04-25
Updated: 2026-04-25

## Summary

Add a sandbox provider contract checklist so future real-provider work has clear readiness, redaction, retry, and rollback expectations before implementation.

## Scope

- Define what information is required before enabling a non-fake provider in sandbox.
- Cover credentials handling, capability mapping, idempotency, retry behavior, timeout/error taxonomy, inventory checks, audit logs, and rollback.
- Reference existing provider architecture docs and operational runbooks.
- Do not include real provider credentials or production account examples.
- Keep each edited file under 500 lines.

## Acceptance Criteria

- Checklist is practical enough for an agent to decide whether a provider task is ready.
- Checklist distinguishes local fake provider, sandbox provider, and production provider readiness.
- Existing validation commands pass.

## Notes

- This is documentation-only unless a small index link is needed.

## Agent Log

- 2026-04-25: Task created in the post-readiness hardening batch.
- 2026-04-25: Codex claimed the task; adding sandbox provider contract checklist docs.
- 2026-04-25: Added provider sandbox contract checklist, linked it from provisioning readiness docs, README, and manifest. Validation passed: `go run ./cmd/contractguard`, `go test ./...`, `go build ./cmd/api ./cmd/migrate ./cmd/seed ./cmd/smoke ./cmd/worker ./cmd/contractguard`.
