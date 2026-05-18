# T258 - Cloudmini rate-limit runtime evidence

Status: IN_PROGRESS
Owner: Codex
Branch: codex/t258-cloudmini-rate-limit-evidence
PR: -
Risk: provider provisioning, credentials, launch-readiness evidence
Created: 2026-05-18
Updated: 2026-05-18

## Summary

Add or record safe Cloudmini V3 rate-limit runtime evidence without tripping the shared provider limiter.

## Scope

- Inspect Billing Cloudmini error evidence support and the non-production provider source for a safe rate-limit fixture or low-limit dev path.
- If a safe provider fixture exists, capture redacted `RATE_LIMITED` runtime metadata and map it to `PROVIDER_RATE_LIMITED` with `safe_retry`.
- If no safe fixture exists, record the exact provider-side blocker and required fixture contract in launch evidence docs.
- Do not generate high-volume traffic, do not call production, and do not print raw provider payloads, provider credentials, tokens, cookies, DSNs, proxy credentials, or raw provider IDs.
- Keep provider 5xx and cancel/delete-rejected evidence out of scope.

## Acceptance Criteria

- Rate-limit evidence is either safely collected with explicit non-production guardrails or explicitly blocked with a provider-side fixture requirement.
- Billing docs accurately reflect whether the rate-limit case is closed or still blocked.
- Relevant tests/checks pass for any code or docs changed.
- PR/CI/merge workflow is followed.

## Notes

- This task does not authorize intentionally tripping the shared Cloudmini limiter.
- If the provider fixture is absent, the task may only add Billing runner support and record the blocker.

## Agent Log

- 2026-05-18: Task created and claimed by Codex from `origin/main`.
- 2026-05-18: Inspected `/opt/proxy-cloudmini/internal/api/router.go`; the inspected provider source has v1/v2 global limiter responses with `RATE_LIMITED`, but no safe V3 low-limit fixture route. Did not call provider or induce rate limits.
- 2026-05-18: Added guarded Billing `cloudmini-error-evidence` rate-limit fixture mode requiring one approved GET request to a side-effect-free `/api/v3/...` fixture path.
- 2026-05-18: Added focused `httptest` coverage for `429 RATE_LIMITED` mapping to `PROVIDER_RATE_LIMITED` with retry safety `safe_retry` and redacted stdout.
- 2026-05-18: Updated launch evidence docs to keep the live provider rate-limit case blocked until a provider-owned fixture is deployed.
- 2026-05-18: Did not run live provider rate-limit evidence because the inspected provider source has no safe V3 fixture and inducing the shared limiter is prohibited.
- 2026-05-18: Validation passed: `make fmt`, `go test ./cmd/smoke -run 'CloudminiErrorEvidence'`, `make test`, `make build`, `go run ./cmd/taskguard`, `git diff --check`, and targeted diff secret-value scan excluding fake test fixtures.
