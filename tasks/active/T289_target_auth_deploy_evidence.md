# T289 - Target auth deploy evidence

Status: REVIEW
Owner: Codex
Branch: codex/t289-target-auth-deploy-evidence
PR: -
Risk: deploy, auth/RBAC, credential handling, tunnel routing, launch evidence
Created: 2026-05-21
Updated: 2026-05-21

## Summary

Record selected test-server deploy and target auth/RBAC evidence after T288 added protected target auth smoke credential overrides.

## Scope

- In scope: redacted deploy evidence for the selected test server.
- In scope: target auth/RBAC smoke and domain-aware public auth probe outcomes.
- In scope: document tunnel routing and frontend backend rewrite fixes applied during test.
- Out of scope: production approval, production customer data, real provider provisioning, full UAT rerun, or committing credential values.

## Acceptance Criteria

- Evidence records deploy commit, service health, public backend health, and auth/RBAC results.
- Evidence records skipped/failed command limitations honestly.
- Evidence does not include raw DSNs, passwords, cookies, session tokens, provider credentials, provider payloads, TOTP values, or plaintext service credentials.
- Required validation passes before PR: `taskguard`, `git diff --check`, docs/task added-line secret scan, and touched-file line-count check.

## Notes

- This is an evidence/docs-only task.
- The generic `dev-target-auth-rbac` command has a one-base-url limitation for public domain testing when tenant resolution is domain-first.

## Agent Log

- 2026-05-21: Task created and claimed on `codex/t289-target-auth-deploy-evidence`.
- 2026-05-21: Recorded deploy, tunnel routing, frontend rewrite, direct target auth smoke, and domain-aware public auth probe evidence.
- 2026-05-21: Validation passed: `go run ./cmd/taskguard`, `git diff --check`, touched-file line-count check, docs/task added-line secret scan, and docs/task added-line UUID scan.
- 2026-05-21: Moved task to REVIEW pending PR.
