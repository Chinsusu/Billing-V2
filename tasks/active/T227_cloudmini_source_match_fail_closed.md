# T227 - Cloudmini source match fail closed

Status: IN_PROGRESS
Owner: Codex
Branch: codex/t227-cloudmini-source-match-fail-closed
PR: -
Risk: provider/provisioning/credential/config
Created: 2026-05-16
Updated: 2026-05-16

## Summary

Make Cloudmini V3 runtime selection fail closed when an operation carries a Billing provider source ID that is not explicitly mapped in runtime config, even if an account-level endpoint exists.

## Scope

- Require explicit source match before using Cloudmini runtime config for operations with `SourceID`.
- Keep account-level endpoint fallback available only for account-level operations that do not carry a source ID.
- Cover source/account mismatch with tests that prove no provider endpoint is called.
- Update the Cloudmini runbook/readiness docs to remove the remaining runtime-config blocker.
- Do not call real provider APIs or run any mutating pilot in this task.

## Acceptance Criteria

- Cloudmini adapter returns `PROVIDER_CONFIG_INVALID` when `operation.SourceID` is set but not configured.
- A configured account endpoint is not used to bypass a mismatched source ID.
- Account endpoint checks without a source ID still work for account-level read checks.
- Relevant provider/worker tests, task guard, and diff checks pass.

## Notes

- This is a safety/code-readiness task before any real create/delete pilot.

## Agent Log

- 2026-05-16: Task created and claimed by Codex from latest `origin/main`.
- 2026-05-16: Updated Cloudmini runtime selection to require explicit source mapping whenever an operation carries `SourceID`; account endpoint fallback remains available only when no source ID is present.
- 2026-05-16: Added provider tests proving account endpoint checks still work without a source ID and source/account mismatch returns `PROVIDER_CONFIG_INVALID` without calling the provider endpoint.
- 2026-05-16: Validation passed: `go test ./internal/modules/provider ./cmd/worker`; `go test ./internal/modules/order`; `go test ./cmd/worker ./internal/modules/provider ./internal/modules/order`; `go test ./...`; `go run ./cmd/taskguard`; `git diff --check`; added-line secret pattern scan returned no matches. Changed-file scan matches an existing proxy URI test fixture that was already present and was not added by this diff.
