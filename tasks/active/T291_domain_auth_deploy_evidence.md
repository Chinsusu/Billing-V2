# T291 - Domain-aware target auth deploy evidence

Status: REVIEW
Owner: Codex
Branch: codex/t291-domain-auth-deploy-evidence
PR: https://github.com/Chinsusu/Billing-V2/pull/614
Risk: auth/RBAC, tenant domain resolution, credential handling, launch evidence
Created: 2026-05-21
Updated: 2026-05-21

## Summary

Record post-T290 selected test-server deploy and remote binary evidence for the domain-aware target auth/RBAC smoke.

## Scope

- In scope: document selected non-production deploy-copy, remote build, health, and split-domain `dev-target-auth-rbac` evidence.
- In scope: record only redacted statuses and safe command outcomes.
- Out of scope: production approval, production customer data, real provider provisioning, money mutations, notification delivery, or credential values.

## Acceptance Criteria

- Evidence doc records the selected test-server deploy and remote binary smoke outcome.
- Evidence excludes passwords, cookies, session tokens, DSNs, provider payloads, Telegram tokens, and credentials.
- Task board passes `taskguard`.
- Docs-only validation passes before PR: `go run ./cmd/taskguard`, `git diff --check`, touched-file line count, added-line secret-pattern scan, and added-line docs/task UUID scan.

## Notes

- This records evidence for the already-approved selected non-production test server only.

## Agent Log

- 2026-05-21: Task created and claimed on `codex/t291-domain-auth-deploy-evidence`.
- 2026-05-21: Recorded post-T290 selected test-server deploy, service health, redacted secret-handling metadata, and remote binary split-domain target auth smoke evidence.
- 2026-05-21: Validation passed: `go run ./cmd/taskguard`, `git diff --check`, touched-file line count, added-line secret-pattern scan, and added-line docs/task UUID scan.
- 2026-05-21: Opened PR #614 and moved task to `REVIEW`.
