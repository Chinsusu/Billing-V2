# T199 - Provider sandbox readiness

Status: TODO
Owner: -
Branch: codex/t199-provider-sandbox-readiness
PR: -
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
