# T235 - Target top-up review E2E evidence

Status: IN_PROGRESS
Owner: Codex
Branch: codex/t235-topup-review-e2e
PR: -
Risk: wallet/RBAC/API/finance/audit
Created: 2026-05-17
Updated: 2026-05-17

## Summary

Prove the fixed reseller top-up review path with a controlled target-environment E2E run that creates, approves, and rejects top-up requests without touching provider provisioning.

## Scope

- Add a repeatable dev/test smoke command for top-up review only.
- Run the smoke on the approved test server after deploying current code.
- Verify approve creates one wallet ledger credit and audit evidence.
- Verify reject creates no wallet ledger credit and audit evidence records rejection.
- Do not create orders, services, provider jobs, or provider resources.
- Record redacted target evidence in launch docs.

## Acceptance Criteria

- Client top-up create succeeds through HTTP API.
- Reseller approve through `/reseller/topup-requests/{id}/approve` succeeds and posts one ledger credit.
- Reseller reject through `/reseller/topup-requests/{id}/reject` succeeds without posting a ledger credit.
- Disallowed provider/provisioning paths are not invoked.
- Tests and required validation commands pass.

## Notes

- This follows T234, which added and deployed the reseller review route.
- The target run may mutate only the dev/test wallet and top-up tables.
- Evidence must not include raw DSNs, secrets, provider IDs, provider payloads, or customer credentials.

## Agent Log

- 2026-05-17: Task created and claimed by Codex on branch `codex/t235-topup-review-e2e`.
- 2026-05-17: Added `dev-topup-review` smoke command using a temporary dev/test wallet; it verifies approve ledger/audit, reject no-ledger/audit, and no order/provider/service side effects.
- 2026-05-17: Deployed current branch to the approved test server and ran `./bin/smoke -timeout 90s dev-topup-review` with `APP_ENV=dev` and local API base URL. Result PASS: approve top-up display `10003`, ledger display `10005`, audit display `10015`; reject top-up display `10004`, audit display `10016`; reject ledger count `0`; wallet delta `111`; provider side effects `none`.
- 2026-05-17: Local focused validation passed: `gofmt`, `go test ./cmd/smoke`, `go run ./cmd/taskguard`, `git diff --check`.
