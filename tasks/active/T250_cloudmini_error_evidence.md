# T250 - Cloudmini error evidence

Status: REVIEW
Owner: Codex
Branch: codex/t250-cloudmini-error-evidence
PR: https://github.com/Chinsusu/Billing-V2/pull/532
Risk: provider provisioning, credentials, provider error handling, launch-readiness evidence
Created: 2026-05-18
Updated: 2026-05-18

## Summary

Capture safe redacted Cloudmini V3 provider error examples for launch-readiness docs.

## Scope

- Add a guarded smoke path for Cloudmini error evidence that prints only status/code/retry-safety metadata.
- Run safe non-production examples on the approved target dev/test provider account where no provider resource should be created.
- Record which real error examples are proven and which require provider-side controls or owner action.
- Do not print or commit DB DSNs, provider tokens, raw provider payloads, cookies, proxy credentials, raw provider IDs, or file contents.

## Acceptance Criteria

- Command refuses production and missing approval fields.
- Output is redacted and states whether mutating routes were called.
- Docs preserve launch `NO-GO` if provider error coverage is still incomplete.
- Required validation for touched files passes.

## Notes

- Do not force rate-limit, 5xx, or cancel-rejected behavior by abusing the provider. If the provider cannot supply those safely, document them as remaining blockers.

## Agent Log

- 2026-05-18: Task created and claimed by Codex on branch `codex/t250-cloudmini-error-evidence`.
- 2026-05-18: Added `go run ./cmd/smoke cloudmini-error-evidence` with fail-closed non-production approval, owner fields, optional malformed-create approval, redacted output, and local httptest coverage.
- 2026-05-18: Target dev/test evidence passed for auth missing, auth invalid, proxy not found, and malformed-create validation. The run printed no raw response body, token, provider ID, provider payload, proxy credential, cookie, or file contents.
- 2026-05-18: Confirmed target `billing-api` and `billing-worker` remained active after evidence; target app env and Cloudmini credential files kept restrictive metadata without reading contents.
- 2026-05-18: Updated provider readiness, Go/No-Go, launch packet, and validation matrix docs with redacted T250 stdout while preserving remaining NO-GO blockers for provider-controlled error examples, shared secret-store proof, usable-status sign-off, and broader owner approval.
- 2026-05-18: Opened PR #532 for review after local and target validation passed.
