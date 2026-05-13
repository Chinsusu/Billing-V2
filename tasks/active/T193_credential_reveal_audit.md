# T193 - Credential reveal audit

Status: DONE
Owner: Codex
Branch: codex/t193-credential-reveal-audit
PR: https://github.com/Chinsusu/Billing-V2/pull/417
Risk: credential security, tenant isolation, RBAC, rate limiting, and audit
Created: 2026-05-13
Updated: 2026-05-13

## Summary

Add controlled credential reveal behavior with masking, authorization, rate limiting, and audit.

## Scope

- Add backend reveal API for service credentials using encrypted storage from T192.
- Enforce tenant/RBAC permission and reveal rate limits.
- Audit every reveal without logging plaintext credential values.
- Add frontend masking/reveal behavior for the relevant client/admin/reseller service detail paths.

## Acceptance Criteria

- Credentials are masked by default and reveal only after authorized action.
- Cross-tenant and unauthorized reveal attempts are denied and tested.
- Reveal audit includes actor, tenant, service, reason/context where required, and no secret plaintext.
- Relevant backend/frontend validation and CI pass.

## Notes

- This task depends on T192.

## Agent Log

- 2026-05-13: Task created by Codex backlog planning.
- 2026-05-13: Claimed by Codex on branch `codex/t193-credential-reveal-audit`.
- 2026-05-13: Implemented masked service credential metadata, controlled reveal API, tenant/owner scoping, RBAC wiring, reveal rate limiting, no-store responses, audit without plaintext, frontend reveal controls, docs, and tests. Local validation passed through backend/frontend build and guards.
- 2026-05-13: Validation run: `make migrate-validate`, `make test`, `make build`, `make contract-guard`, `make error-code-guard`, `make task-guard`, `npm --prefix frontend audit --omit=dev`, `check:sensitive-text`, `lint`, `build`, `smoke:admin:ci`, `git diff --check`, and hardcoded secret grep.
- 2026-05-13: Opened PR https://github.com/Chinsusu/Billing-V2/pull/417 and moved task to REVIEW.
- 2026-05-13: PR #417 merged into `main`; moved task to DONE.
