# T248 - Cloudmini idempotency evidence

Status: IN_PROGRESS
Owner: Codex
Branch: codex/t248-cloudmini-idempotency-evidence
PR: -
Risk: provider provisioning, idempotency, credentials, launch-readiness evidence
Created: 2026-05-18
Updated: 2026-05-18

## Summary

Add a guarded evidence path for Cloudmini duplicate-create, timeout-after-send, and redacted provider error behavior.

## Scope

- Inspect existing Cloudmini adapter, preflight, pilot, and evidence tooling.
- Add the smallest fail-closed script/docs needed to collect duplicate/timeout/error evidence safely.
- Keep all raw DSNs, provider tokens, raw provider IDs, provider payloads, proxy credentials, cookies, and file contents out of repo output.
- Do not run a mutating provider scenario unless the script enforces explicit non-production approval and one-resource cleanup boundaries.

## Acceptance Criteria

- Evidence path refuses production and missing approval fields.
- Evidence output is redacted and documents whether mutating routes were called.
- Provider/Go-No-Go docs keep launch `NO-GO` unless live evidence is actually captured.
- Required validation for touched files passes.

## Notes

- This task may produce tooling/docs only if a safe live run cannot be completed without stronger provider-side controls.

## Agent Log

- 2026-05-18: Task created and claimed by Codex on branch `codex/t248-cloudmini-idempotency-evidence`.
- 2026-05-18: Added `go run ./cmd/smoke cloudmini-idempotency-evidence` with fail-closed approval, owner, non-production, guardrail, redacted-output, raw-cleanup-reference, duplicate-create, timeout-after-send, and cleanup handling.
- 2026-05-18: Added local httptest coverage for duplicate-create redaction/cleanup, timeout-after-send manual-review mapping/cleanup, and missing approval refusal.
- 2026-05-18: Updated provider readiness, Go/No-Go, launch evidence, controlled pilot, and validation command docs to record that T248 adds tooling only; live Cloudmini duplicate/timeout evidence remains missing until the command is run against an approved non-production provider account.
- 2026-05-18: Local validation passed: `go test ./cmd/smoke -count=1`, `go test ./internal/modules/provider -count=1`, `go build ./cmd/api ./cmd/migrate ./cmd/seed ./cmd/smoke ./cmd/worker`, `make test`, `go run ./cmd/taskguard`, `git diff --check`, changed-file line counts under 500, and staged added-line secret-pattern scan.
