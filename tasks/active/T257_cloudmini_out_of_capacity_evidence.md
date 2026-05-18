# T257 - Cloudmini out-of-capacity runtime evidence

Status: REVIEW
Owner: Codex
Branch: codex/t257-cloudmini-out-of-capacity
PR: https://github.com/Chinsusu/Billing-V2/pull/546
Risk: provider provisioning, credentials, launch-readiness evidence
Created: 2026-05-18
Updated: 2026-05-18

## Summary

Add and run guarded Cloudmini V3 out-of-capacity evidence using an exhausted-group reservation probe without proxy create/delete side effects.

## Scope

- Extend the Cloudmini error evidence smoke to optionally find one exhausted group from read-only inventory.
- Call `POST /api/v3/capacity/reservations` at most once with quantity `1` and TTL no more than `60s`.
- Expect `CAPACITY_EXHAUSTED`, map it to `PROVIDER_OUT_OF_STOCK`, and record only redacted metadata.
- Fail closed if the probe unexpectedly creates a reservation, and do not print raw group IDs, raw provider payloads, provider credentials, cookies, DSNs, or proxy credentials.
- Update provider launch evidence docs with the runtime result or an explicit blocker if the runtime cannot be executed safely.
- Keep rate-limit, provider-5xx, and cancel/delete-rejected cases out of scope.

## Acceptance Criteria

- Out-of-capacity case maps provider `CAPACITY_EXHAUSTED` to `PROVIDER_OUT_OF_STOCK` with `do_not_retry`.
- Smoke output records bounded reservation-probe metadata and no secret values.
- The probe is guarded by explicit non-production approval env and max attempt `1`.
- Relevant Go tests and task/docs checks pass.
- PR/CI/merge workflow is followed.

## Notes

- This task may close only the out-of-capacity provider-controlled error case. Broader provider approval and the remaining provider-controlled cases remain separate.

## Agent Log

- 2026-05-18: Task created and claimed by Codex from `origin/main`.
- 2026-05-18: Added guarded out-of-capacity mode to `cloudmini-error-evidence`, with one exhausted-group inventory selection, one reservation probe, TTL capped at 60 seconds, and fail-closed cleanup if a reservation is unexpectedly created.
- 2026-05-18: Ran focused smoke tests with `go test ./cmd/smoke -run 'CloudminiErrorEvidence'`; passed.
- 2026-05-18: Ran approved dev/test Cloudmini out-of-capacity evidence; result `PASS`, `409 CAPACITY_EXHAUSTED` mapped to `PROVIDER_OUT_OF_STOCK`, retry safety `do_not_retry`, exhausted group selected, reservation probe attempted once, and reservation_created=false.
- 2026-05-18: Updated provider launch evidence docs with redacted runtime output; no token, raw group ID, raw provider payload, DSN, cookie, or proxy credential was recorded.
- 2026-05-18: Validation passed: `make fmt`, `go test ./cmd/smoke -run 'CloudminiErrorEvidence'`, approved `go run ./cmd/smoke cloudmini-error-evidence`, `make test`, `make build`, `go run ./cmd/taskguard`, `git diff --check`, and diff secret-like scan excluding fake test fixtures.
- 2026-05-18: Opened PR #546 and moved task to `REVIEW`.
