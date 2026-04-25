# T130 - Reseller public ID labels

Status: DONE
Owner: Codex
Branch: codex/t130-reseller-public-id-labels
PR: https://github.com/Chinsusu/Billing-V2/pull/291
Risk: frontend
Created: 2026-04-26
Updated: 2026-04-26

## Summary

Clean up reseller-facing billing and service screens so client/order/service labels use public numeric IDs even when backend-reference lookups are unavailable.

## Scope

- Prefer `buyer_display_id` and `account_display_id` for reseller client labels.
- Prefer related public IDs for reseller order/service relation labels.
- Stop using service external provider references as primary reseller service labels.
- Keep backend IDs only for joins and action bodies where required.

## Acceptance Criteria

- Reseller billing rows show account public IDs from invoice/transaction responses when customer lookup data is incomplete.
- Reseller service rows show `SVC-` and `ORD-` public labels, with plan text as detail rather than provider external IDs.
- Reseller client counters can use public buyer IDs where available.
- Frontend lint, sensitive-text check, build, taskguard, and diff check pass.

## Notes

- External provider resource IDs may remain internal/support data, but not the primary table label.

## Agent Log

- 2026-04-26: Codex created and claimed the task after T129 was marked done; starting reseller public ID label cleanup.
- 2026-04-26: Codex updated reseller billing, clients, dashboard, and services to prefer public account/order/service labels and shared account label mapping. Local validation passed: `npm --prefix frontend run lint`, `npm --prefix frontend run check:sensitive-text`, `npm --prefix frontend run build`, `go run ./cmd/taskguard`, and `git diff --check`.
- 2026-04-26: Opened PR https://github.com/Chinsusu/Billing-V2/pull/291 for review and CI.
- 2026-04-26: PR https://github.com/Chinsusu/Billing-V2/pull/291 passed CI and was merged to `main`; task marked DONE.
