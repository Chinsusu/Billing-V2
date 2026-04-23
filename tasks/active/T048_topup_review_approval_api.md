# T048 - Top-up review and approval API

Status: DONE
Owner: Codex
Branch: feat/topup-review-approval-api
PR: https://github.com/Chinsusu/Billing-V2/pull/110
Risk: wallet/topup/money/API
Created: 2026-04-23
Updated: 2026-04-23

## Summary

Add admin review actions for wallet top-up requests, including approval that credits the wallet.

## Scope

- Add admin approve/reject endpoints for top-up requests.
- Use wallet ledger posting service to credit approved requests.
- Store reviewer, review time, review note, and linked ledger entry.
- Make approval idempotent for repeated review calls.
- Keep client create/read behavior unchanged.

## Acceptance Criteria

- Only submitted/under_review requests can be approved or rejected.
- Approval credits exactly one wallet ledger entry and marks the request approved.
- Reject stores reviewer and reason without changing wallet balance.
- Duplicate approval calls return the already approved request without duplicate credit.
- Handler/store/service tests cover approve, reject, duplicate, invalid status, and auth scoping.
- Full validation passes: `make fmt`, `make test`, `make build`, `make migrate-validate`.

## Notes

- This task depends on T044 and T047.

## Agent Log

- 2026-04-23: Task created for manual wallet funding workflow.
- 2026-04-23: Implemented admin approve/reject top-up review flow, ledger credit handoff, review middleware, and focused tests.
- 2026-04-23: PR #110 passed checks and was merged.
