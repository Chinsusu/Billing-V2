# T238 - Target finance reconciliation evidence

Status: DONE
Owner: Codex
Branch: codex/t238-target-finance-reconciliation
PR: https://github.com/Chinsusu/Billing-V2/pull/508
Risk: finance/wallet/ledger/RBAC/audit
Created: 2026-05-17
Updated: 2026-05-17

## Summary

Add and run a repeatable target-environment smoke that proves finance reconciliation read paths can be executed against the approved dev/test target without mutating money, provider, order, or service state.

## Scope

- Add a `dev-target-finance-reconciliation` smoke command for approved dev/test environments.
- Select an existing posted wallet transaction from the target DB and verify payment reconciliation list/detail API responses.
- Verify daily reconciliation for the selected transaction date.
- Verify database counters and wallet balance projection do not change before/after the read-only API checks.
- Run the smoke on the approved test server and record redacted launch evidence.
- Do not create top-ups, payments, orders, provider jobs, services, credentials, or provider resources.

## Acceptance Criteria

- Smoke fails unless `APP_ENV` is non-production and required local API/database config is present.
- Smoke output excludes raw transaction, invoice, wallet, ledger, actor, session, cookie, DSN, provider payload, and credential values.
- Payment reconciliation list returns the selected posted wallet transaction by public display filters.
- Payment reconciliation detail returns matching public transaction, invoice, wallet, and ledger display IDs.
- Daily reconciliation for the selected date returns a coherent `balanced` or `mismatched` report; mismatched output must include mismatch counts and keeps GO blocked until Finance owner review.
- Before/after DB baseline proves no money, order, provider job, service, or provider-resource mutation happened during the smoke.
- Launch evidence docs reflect the target finance reconciliation evidence while still requiring named Finance owner sign-off before GO.

## Notes

- This task captures target finance reconciliation evidence only; it does not assign or replace a named Finance owner.
- The smoke is read-only and must never be run against production or real customer data.

## Agent Log

- 2026-05-17: Task created and claimed by Codex on branch `codex/t238-target-finance-reconciliation`.
- 2026-05-17: Added `dev-target-finance-reconciliation` smoke command and focused unit coverage for read-only finance headers, non-leaking status errors, and mismatched daily reconciliation evidence handling.
- 2026-05-17: Deployed current branch to the approved test server and ran `./bin/smoke -timeout 90s dev-target-finance-reconciliation` with `APP_ENV=dev` and local API base URL. Result PASS for evidence collection: transaction display `51001`, invoice display `44001`, wallet display `41001`, ledger display `50002`, daily date `2026-04-23`, daily status `mismatched`, wallets checked `2`, wallet mismatches `1`, invoices checked `1`, invoice mismatches `0`, payments checked `1`, duplicate payment references `0`, money mutation routes called `no`, and provider mutation routes called `no`. Finance owner review remains required because the report is not balanced.
- 2026-05-17: Local validation passed: `make fmt`, `go test ./cmd/smoke`, `make test`, `make build`, `go run ./cmd/contractguard`, `go run ./cmd/errorcodeguard`, `go run ./cmd/taskguard`, and `git diff --check`.
- 2026-05-17: Opened PR #508 and moved task to `REVIEW`.
- 2026-05-17: PR #508 merged into `main`; task marked `DONE`.
