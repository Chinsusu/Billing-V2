# T237 - Target credential reveal audit evidence

Status: REVIEW
Owner: Codex
Branch: codex/t237-target-credential-reveal
PR: https://github.com/Chinsusu/Billing-V2/pull/506
Risk: credentials/auth/RBAC/tenant/audit
Created: 2026-05-17
Updated: 2026-05-17

## Summary

Add and run a repeatable target-environment smoke that proves credential reveal uses real session auth, returns plaintext only through the reveal endpoint, sets no-store headers, updates reveal metadata, and writes a redacted `credential.revealed` audit record.

## Scope

- Add a `dev-target-credential-reveal` smoke command for approved dev/test environments.
- Create or refresh only one dev/test encrypted service credential fixture for the seeded demo service.
- Reveal the credential through the client API using an HttpOnly session cookie, not `X-Actor-*` dev helper headers.
- Verify response envelope, no-store headers, redaction boundaries, `last_revealed_by`, reveal rate-limit tracking, and audit metadata.
- Run the smoke on the approved test server and record redacted launch evidence.
- Do not call provider mutating routes, money mutating routes, or lifecycle worker commands.

## Acceptance Criteria

- Smoke fails unless `APP_ENV` is non-production and required local API/database/encryption config is present.
- Smoke output excludes plaintext credentials, encrypted payloads, raw credential IDs, session tokens, cookies, DSNs, provider payloads, and provider credentials.
- Client session cookie-only reveal succeeds for the seeded service and fixture credential.
- Reveal response includes `request_id`, `masked_hint`, `revealed_at`, `reveal_expires_message`, and expected payload content without logging the payload.
- Reveal response sets `Cache-Control: no-store` and `Pragma: no-cache`.
- Database verification proves `last_revealed_by` is the seeded client actor, reveal rate-limit state exists, and `credential.revealed` audit metadata includes the service display ID and no plaintext/encrypted credential material.
- Launch evidence docs reflect the target credential reveal audit result without changing the overall NO-GO decision.

## Notes

- This task may mutate only dev/test credential fixture rows, reveal metadata, reveal rate-limit rows, auth session rows, and audit rows.
- The smoke must never be run against production or real customer data.

## Agent Log

- 2026-05-17: Task created and claimed by Codex on branch `codex/t237-target-credential-reveal`.
- 2026-05-17: Added `dev-target-credential-reveal` smoke command and focused unit coverage for cookie-only reveal, response redaction checks, and audit leak detection errors.
- 2026-05-17: Deployed current branch to the approved test server and ran `./bin/smoke -timeout 90s dev-target-credential-reveal` with `APP_ENV=dev` and local API base URL. Result PASS: service display `43001`, credential type `recovery_code`, client session cookie-only reveal, no-store response headers, audit display `10017`, client actor reveal metadata, one reveal rate-limit attempt, provider mutation routes called `no`, and money mutation routes called `no`.
- 2026-05-17: Local validation passed: `make fmt`, `go test ./cmd/smoke`, `make test`, `make build`, `go run ./cmd/contractguard`, `go run ./cmd/errorcodeguard`, `go run ./cmd/taskguard`, CI-equivalent basic secret scan, and `git diff --check`.
- 2026-05-17: Opened PR #506 and moved task to `REVIEW`.
