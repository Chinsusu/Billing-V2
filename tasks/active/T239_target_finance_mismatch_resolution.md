# T239 - Target finance mismatch resolution

Status: IN_PROGRESS
Owner: Codex
Branch: codex/t239-finance-mismatch-investigation
PR: -
Risk: finance/wallet/ledger/database/audit
Created: 2026-05-17
Updated: 2026-05-17

## Summary

Investigate and resolve the target-environment finance reconciliation mismatch surfaced by T238 so launch evidence can distinguish a real code defect from dev/test fixture drift.

## Scope

- Investigate the T238 target daily reconciliation mismatch with read-only database/API checks first.
- Identify whether the mismatch is caused by code logic, wallet projection drift, missing ledger entries, stale dev/test fixture data, or unsupported launch evidence assumptions.
- If code is wrong, fix the owned module and add focused tests.
- If target dev/test data is wrong, document the exact safe correction path and only apply a correction when it preserves append-only ledger rules and remains non-production.
- Re-run `dev-target-finance-reconciliation` on the approved target server and update launch evidence.
- Do not update posted ledger amounts or delete ledger rows.
- Do not print or commit raw DSNs, secrets, backend UUIDs, raw provider payloads, or customer data.

## Acceptance Criteria

- Root cause for wallet display `41001` mismatch is recorded with redacted/public evidence.
- Any code change has focused tests and does not weaken ledger, tenant, RBAC, audit, or credential safety.
- Any data correction is non-production only, append-only where money state is involved, and documented with rollback/residual risk.
- `dev-target-finance-reconciliation` is re-run on the approved target server after the resolution path.
- GO status is updated honestly: `GO` only if reconciliation evidence is balanced and remaining P0 owner/sign-off gates are satisfied; otherwise the specific remaining blocker is documented.

## Notes

- T238 result: target daily reconciliation for `2026-04-23` returned `mismatched` with one wallet mismatch while read-only API checks and before/after DB baselines passed.
- The task may end with a documented blocker if the required correction needs Finance owner approval.

## Agent Log

- 2026-05-17: Task created and claimed by Codex on branch `codex/t239-finance-mismatch-investigation`.
- 2026-05-17: Read required workflow, testing, database, ledger, and reconciliation docs before code/data changes.
- 2026-05-17: Read-only target investigation found wallet display `41001` had available balance `950`, posted ledger sum `1750`, difference `-800`, seven posted ledger entries, and no invoice or duplicate-payment mismatch. Seed baseline was inconsistent: public ledger `50001` credit `5000` and `50002` debit `1400` imply wallet projection `3600`, but seed expected `3200`.
- 2026-05-17: Fixed dev billing seed baseline to `3600`, updated smoke expectations, added seed test coverage, and added `scripts/dev_wallet_projection_repair.sh` as a non-production approved repair tool that updates wallet projection from posted ledger source-of-truth and writes an audit row without inserting or updating ledger rows.
- 2026-05-17: Deployed current branch to the approved target test server, ran repair plan for wallet display `41001`, then applied projection repair with `APP_ENV=dev` and `BILLING_DEV_WALLET_PROJECTION_REPAIR_APPROVED=yes`. Result: before `950`, after `1750`, ledger `1750`, audit display `10018`, ledger rows inserted `0`, posted ledger rows updated `0`, mutating routes called `no`, secrets printed `no`.
- 2026-05-17: Target validation passed: `go test ./internal/seed ./cmd/smoke`; `APP_ENV=dev API_BASE_URL=http://127.0.0.1:8080 ./bin/smoke -timeout 90s dev-target-finance-reconciliation` returned `balanced` for daily date `2026-04-23` with wallet mismatches `0`, invoice mismatches `0`, duplicate payment references `0`, money mutation routes called `no`, and provider mutation routes called `no`.
- 2026-05-17: Local validation passed: `make fmt`, `bash -n scripts/dev_wallet_projection_repair.sh`, `go test ./internal/seed ./cmd/smoke`, `make test`, `make build`, `go run ./cmd/migrate validate`, `go run ./cmd/taskguard`, `go run ./cmd/contractguard`, `go run ./cmd/errorcodeguard`, and `git diff --check`. Local `dev-db` smoke was not run because local `DB_DSN` is not set; target finance smoke covered the approved test-server DB after repair.
