# T199 - Provider sandbox readiness

Status: REVIEW
Owner: Codex
Branch: codex/t199-provider-sandbox-readiness
PR: https://github.com/Chinsusu/Billing-V2/pull/429
Risk: provider provisioning, credentials, idempotency, manual review, and operations
Created: 2026-05-13
Updated: 2026-05-13

## Summary

Prove provider sandbox readiness for one VPS source and one proxy/manual source before pilot.

## Scope

- Execute or automate provider sandbox contract checks for approved non-production sources.
- Verify idempotency, timeout/manual-review behavior, health/readiness, and redacted attempts.
- Document provider readiness evidence using display IDs and redacted errors only.
- Do not add a new provider unless the adapter contract and approval are clear.

## Acceptance Criteria

- One VPS source and one proxy/manual source have documented sandbox readiness status.
- Provider timeout after create does not blindly retry.
- Provider errors and attempts are redacted.
- Provider tests/smoke and CI pass.

## Notes

- Never use production provider credentials or raw provider responses in tasks, PRs, logs, or docs.

## Agent Log

- 2026-05-13: Task created by Codex backlog planning.
- 2026-05-13: Codex claimed task on `codex/t199-provider-sandbox-readiness`.
- 2026-05-13: Added provider sandbox readiness/no-go evidence. Real sandbox provider remains blocked because approved account, credentials, quota, SKU mapping, timeout policy, redacted examples, and cleanup owner are not present.
- 2026-05-13: Added local fake provider contract coverage for proxy and request-known timeout mapping; added provisioning worker test proving timeout-after-create moves to manual review instead of blind retry.
- 2026-05-13: Local checks passed: `make fmt`, `go test ./internal/modules/provider -run SandboxContract`, `go test ./internal/modules/order -run ProviderProvisioningHandler`, `go test ./internal/modules/provider ./internal/modules/order ./cmd/smoke`, `make test`, `make build`, `make migrate-validate`, `make contract-guard`, `make error-code-guard`, `make task-guard`, `git diff --check`. `make smoke-dev-billing` blocked because `DB_DSN`/`-dsn` is not configured.
- 2026-05-13: Opened PR #429 for review.
