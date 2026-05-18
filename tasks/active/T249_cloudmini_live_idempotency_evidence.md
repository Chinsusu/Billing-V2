# T249 - Cloudmini live idempotency evidence

Status: REVIEW
Owner: Codex
Branch: codex/t249-cloudmini-live-idempotency-evidence
PR: https://github.com/Chinsusu/Billing-V2/pull/530
Risk: provider provisioning, credentials, mutating Cloudmini routes, launch-readiness evidence
Created: 2026-05-18
Updated: 2026-05-18

## Summary

Run the guarded T248 Cloudmini idempotency evidence smoke against the approved target dev/test provider account and record only redacted results.

## Scope

- Verify the target dev/test runtime has the latest T248 smoke command available.
- Run the duplicate-create scenario with one active-resource guardrail and same-session cleanup.
- Run the timeout-after-send scenario with one active-resource guardrail and same-session cleanup.
- Keep raw cleanup references outside git with restricted permissions.
- Update launch-readiness docs with redacted stdout evidence or a precise blocker.
- Do not print or commit DB DSNs, provider tokens, raw provider IDs, provider payloads, proxy credentials, cookies, or raw cleanup reference contents.

## Acceptance Criteria

- Evidence runs refuse production and require explicit owner/approval fields.
- Duplicate-create evidence records two create attempts, one distinct redacted resource, successful cleanup, and no sensitive stdout.
- Timeout-after-send evidence records request-known timeout/manual-review behavior, successful cleanup, and no sensitive stdout.
- Docs keep launch `NO-GO` if either live scenario cannot be completed.
- Required validation for touched files passes.

## Notes

- Owner-approved values for this dev/test run are `Admin` for source/account, engineering, ops, security, cleanup, finance/quota, and reviewer signoff.
- Use the protected target credential/env paths only; do not copy secret contents into the repo.

## Agent Log

- 2026-05-18: Task created and claimed by Codex on branch `codex/t249-cloudmini-live-idempotency-evidence`.
- 2026-05-18: Synced the T248 smoke command to the approved target dev/test server at `/opt/Billing`, preserving local env and credential files outside git.
- 2026-05-18: Added a separate cleanup poll timeout for the idempotency evidence smoke so forced timeout-after-send evidence does not force cleanup to use the same short timeout.
- 2026-05-18: Target duplicate-create evidence passed with two create attempts, one distinct redacted resource, `duplicate_same_resource=true`, raw cleanup ref mode `0600`, and cleanup success.
- 2026-05-18: Target timeout-after-send evidence passed with `PROVIDER_TIMEOUT_REQUEST_KNOWN`, `manual_review_required`, raw cleanup ref mode `0600`, and cleanup success.
- 2026-05-18: Confirmed target `billing-api` and `billing-worker` remained active after evidence; target app env and Cloudmini credential files kept restrictive metadata without reading contents.
- 2026-05-18: Updated provider readiness, Go/No-Go, launch packet, and idempotency runbook docs with redacted T249 evidence while preserving remaining NO-GO blockers.
- 2026-05-18: Local validation passed: `go test ./cmd/smoke -count=1`, `go test ./internal/modules/provider -count=1`, `go build ./cmd/api ./cmd/migrate ./cmd/seed ./cmd/smoke ./cmd/worker`, `make test`, `go run ./cmd/taskguard`, `git diff --check`, changed-file line counts under 500, and staged added-line secret-pattern scan.
- 2026-05-18: Opened PR https://github.com/Chinsusu/Billing-V2/pull/530 for review.
