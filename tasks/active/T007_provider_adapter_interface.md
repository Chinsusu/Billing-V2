# T007 - Provider Adapter Interface

Status: REVIEW
Owner: Codex
Branch: feat/provider-adapter-interface
PR: https://github.com/Chinsusu/Billing-V2/pull/28
Risk: provider/credential
Created: 2026-04-22
Updated: 2026-04-22

## Summary

Add provider adapter interface, normalized provider error types, and fake adapter tests.

## Scope

- Define provider adapter interface.
- Define normalized provider error categories.
- Add fake adapter for tests.
- Avoid real provider credentials or production endpoints.

## Acceptance Criteria

- Error categories distinguish retryable, non-retryable, timeout, partial success, and manual review cases where needed.
- No secrets or provider credentials are committed.
- Fake adapter tests cover success and failure paths.
- `make test` passes.
- `make build` passes.

## Notes

- Follow provider runtime and credential security docs before implementation.
- Provider behavior affects provisioning and must stay explicit.

## Agent Log

- 2026-04-22: Task file created from `TASKS.md`.
- 2026-04-22: PR #22 was closed and branch `feat/provider-adapter-interface` was deleted per owner request because CI failed and the task file was not updated. Task remains TODO for a clean restart.
- 2026-04-22: Claimed by Codex. Rebuilding provider adapter interface, error taxonomy, and fake adapter tests from latest origin/main.
- 2026-04-22: Opened PR https://github.com/Chinsusu/Billing-V2/pull/28. Validation passed: gofmt, make fmt, make test, make build, make migrate-validate.
