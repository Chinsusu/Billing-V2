# T208 - Provider sandbox evidence packet

Status: DONE
Owner: Codex
Branch: codex/t208-provider-sandbox-evidence
PR: https://github.com/Chinsusu/Billing-V2/pull/447
Risk: provider provisioning, credentials, idempotency, manual review, and launch readiness
Created: 2026-05-14
Updated: 2026-05-14

## Summary

Turn the real-provider sandbox blocker into a concrete evidence packet that an operator/provider owner can fill before reconsidering launch readiness.

## Scope

- Extend the provider sandbox readiness evidence record with required proof fields, redaction rules, and pilot-run evidence requirements.
- Keep real provider sandbox readiness blocked until approved external provider account and redacted run evidence exist.
- Update launch audit/go-no-go docs only enough to reference the new evidence packet.
- Do not add provider credentials, raw provider payloads, or claim that a real provider sandbox test passed.

## Acceptance Criteria

- Provider sandbox evidence requirements are explicit for account, credential path, quota, SKU mapping, timeout/idempotency, cleanup, and redacted examples.
- Launch docs still say provider sandbox is blocked and pilot remains NO-GO until real evidence is provided.
- Task board and whitespace validation pass.

## Notes

- No real provider account, credential, quota approval, or sandbox run evidence is present in this repository as of task creation.

## Agent Log

- 2026-05-14: Codex created and claimed task on `codex/t208-provider-sandbox-evidence`.
- 2026-05-14: Local validation passed: `go run ./cmd/taskguard`, `git diff --check`, `env GOMODCACHE=/tmp/go-mod-cache GOCACHE=/tmp/go-build-cache go test ./internal/modules/provider -run SandboxContract`.
- 2026-05-14: Opened PR #447 for review.
- 2026-05-14: PR #447 merged into `main` with squash commit `68a78dc0f303a5790e4f30c93a6e6fadd6f5e9a1`; marking task DONE.
