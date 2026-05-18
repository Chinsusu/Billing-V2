# T256 - Cloudmini permission-denied runtime evidence

Status: IN_PROGRESS
Owner: Codex
Branch: codex/t256-cloudmini-permission-denied
PR: -
Risk: provider provisioning, credentials, launch-readiness evidence
Created: 2026-05-18
Updated: 2026-05-18

## Summary

Add and run guarded Cloudmini V3 permission-denied evidence without proxy create/delete side effects.

## Scope

- Extend the Cloudmini error evidence smoke to optionally create one temporary low-scope provider API key.
- Use that temporary key to call a V3 route requiring `proxy_crud` and record redacted `403` metadata.
- Revoke the temporary key in the same run and verify cleanup without printing plaintext keys, raw provider payloads, provider IDs, cookies, DSNs, or proxy credentials.
- Update provider launch evidence docs with the runtime result or an explicit blocker if the runtime cannot be executed safely.
- Keep rate-limit, out-of-capacity, provider-5xx, and cancel/delete-rejected cases out of scope.

## Acceptance Criteria

- Permission-denied case maps HTTP `403` to `PROVIDER_PERMISSION_DENIED` with `do_not_retry`.
- Smoke output records only redacted metadata and no secret values.
- Temporary provider API key creation is bounded to one key and revoked.
- Relevant Go tests and task/docs checks pass.
- PR/CI/merge workflow is followed.

## Notes

- This task may close only the permission-denied provider-controlled error case. Broader provider approval and the other provider-controlled cases remain separate.

## Agent Log

- 2026-05-18: Task created and claimed by Codex from `origin/main`.
- 2026-05-18: Added guarded permission-denied mode to `cloudmini-error-evidence`, with one temporary low-scope key, same-run revoke, and active key count restoration.
- 2026-05-18: Ran focused smoke tests with `go test ./cmd/smoke -run 'CloudminiErrorEvidence'`; passed.
- 2026-05-18: Ran approved dev/test Cloudmini permission-denied evidence; result `PASS`, `403` mapped to `PROVIDER_PERMISSION_DENIED`, temporary key revoked, active key count restored, and `mutating_routes_called=true` for API-key create/revoke support routes.
- 2026-05-18: Updated provider launch evidence docs with redacted runtime output; no token, key, raw provider ID, raw provider payload, DSN, cookie, or proxy credential was recorded.
- 2026-05-18: Validation passed: `make fmt`, `go test ./cmd/smoke -run 'CloudminiErrorEvidence'`, approved `go run ./cmd/smoke cloudmini-error-evidence`, `make test`, `make build`, `go run ./cmd/taskguard`, `git diff --check`, and diff secret-like scan excluding fake test fixtures.
