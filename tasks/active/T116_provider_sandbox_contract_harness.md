# T116 - Provider sandbox contract harness

Status: REVIEW
Owner: Codex
Branch: codex/t116-provider-sandbox-contract-harness
PR: https://github.com/Chinsusu/Billing-V2/pull/262
Risk: provider/testing
Created: 2026-04-25
Updated: 2026-04-25

## Summary

Create the first provider sandbox contract test harness shape so future real provider integrations have a consistent validation path.

## Scope

- Add a small harness for provider adapter sandbox contract cases without requiring real credentials.
- Cover at least one fake provider scenario for quote, order, status, cancel, and idempotency behavior where the current adapter interface supports it.
- Document how real provider sandbox credentials should be plugged in later without committing secrets.
- Keep each edited file under 500 lines.

## Acceptance Criteria

- The harness can run in local/CI mode without external network credentials.
- It defines clear pass/fail expectations for provider adapter behavior.
- Docs link the harness to the provider sandbox checklist.
- Existing tests pass.

## Notes

- Do not add real provider credentials.
- This task should prefer test helpers over production abstractions unless a production type is already needed.

## Agent Log

- 2026-04-25: Task created in the board and delivery hardening batch.
- 2026-04-25: Codex claimed the task; adding a provider sandbox contract harness shape without real credentials.
- 2026-04-25: Added a local provider sandbox contract harness for health, quote/stock, provision, status, cancel/terminate, and idempotency repeat cases, with fake-adapter tests and checklist docs.
- 2026-04-25: Opened PR #262. Validation passed: `go test ./internal/modules/provider -run SandboxContract`, `go run ./cmd/contractguard`, `go run ./cmd/errorcodeguard`, `go test ./...`, `go build ./cmd/api ./cmd/migrate ./cmd/seed ./cmd/smoke ./cmd/worker ./cmd/contractguard ./cmd/taskguard ./cmd/errorcodeguard`, `git diff --check`.
