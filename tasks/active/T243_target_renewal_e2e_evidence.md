# T243 - Target staging-equivalent renewal E2E evidence

Status: REVIEW
Owner: Codex
Branch: codex/t243-target-renewal-e2e-evidence
PR: https://github.com/Chinsusu/Billing-V2/pull/518
Risk: wallet/ledger/service-lifecycle/audit/full-e2e/launch-readiness
Created: 2026-05-17
Updated: 2026-05-17

## Summary

Add and capture staging-equivalent full E2E evidence for the T206 client service renewal path.

## Scope

- Extend the existing billing/full-E2E smoke path to cover direct client service renewal after provisioning.
- Verify renewal wallet debit, paid renewal invoice, posted payment transaction, ledger entry, active service state, increased `term_end`, and audit evidence.
- Run the updated gate on the approved target server using temporary non-production DBs only.
- Record redacted launch evidence without DSNs, passwords, tokens, provider payloads, service credentials, or customer data.

## Acceptance Criteria

- `dev-billing` covers checkout, wallet payment, provisioning, service activation, renewal, and renewal audits.
- Target staging-equivalent run passes against a temporary target DB and cleans up temporary DBs/artifacts.
- Launch docs show the T206 renewal path evidence and remaining blockers without marking GO prematurely.
- Required local validation passes.

## Notes

- User-provided owner assignment from T241: `Admin` owns Product, Engineering, QA, Ops, Finance, Security, Support, and Provider launch roles.
- Keep the target run non-production and use fake/manual provider paths only unless a separate provider-mutating approval exists.

## Agent Log

- 2026-05-17: Task created and claimed by Codex on branch `codex/t243-target-renewal-e2e-evidence`.
- 2026-05-17: Added renewal coverage to `dev-billing`: after fake-provider service activation, the smoke renews the service, verifies a paid renewal invoice, posted payment transaction, purchase ledger, active paid service state, increased `term_end`, and `service.renewed`/`invoice.wallet_paid` audit rows.
- 2026-05-17: Target staging-equivalent full E2E passed on temporary DB `billing_t243_e2e_20260517140625`: backend taskguard/test/contract/error/build passed, target deploy-copy skipped only `git diff --check`, `dev-db` passed 25 migrations and 20 checks, `dev-api` passed 35 checks, `dev-billing` passed with service `10000` and renewal invoice `10002`, transaction `10001`, ledger `10002`, frontend `npm ci`, audit, sensitive-text, lint, build, and admin browser smoke passed. Temporary DB cleanup verified count `0`.
- 2026-05-17: Local validation passed: `make fmt`, `go test ./cmd/smoke`, `bash -n scripts/full_e2e_quality_gate.sh`, `make test`, `make build`, `go run ./cmd/contractguard`, `go run ./cmd/errorcodeguard`, `go run ./cmd/taskguard`, `git diff --check`, and staged secret-pattern scan.
- 2026-05-17: Opened PR https://github.com/Chinsusu/Billing-V2/pull/518 for review.
