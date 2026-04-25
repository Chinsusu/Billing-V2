# T128 - Order buyer public ID response

Status: REVIEW
Owner: Codex
Branch: codex/t128-order-buyer-public-id-response
PR: https://github.com/Chinsusu/Billing-V2/pull/287
Risk: API/frontend
Created: 2026-04-25
Updated: 2026-04-25

## Summary

Expose buyer public ID on order responses so admin/client flows can show account labels without backend user references.

## Scope

- Add `buyer_display_id` to order read responses.
- Keep create/update write paths compatible with existing `RETURNING` column counts.
- Update frontend API types.
- Extend API smoke checks for client and admin order responses.

## Acceptance Criteria

- Admin order list/detail include `buyer_display_id` for seeded orders.
- Client order list/detail include `buyer_display_id` for the current buyer.
- Existing order tests, smoke tests, frontend type build, and taskguard pass.

## Notes

- Order filters already support `buyer_display_id`; this task completes the read response contract.

## Agent Log

- 2026-04-25: Codex created and claimed the task after T127 merged; starting order buyer public ID response support.
- 2026-04-25: Codex implemented `buyer_display_id` in order read responses, frontend API type, smoke checks, and API reference. Local validation passed: `go test $(go run ./cmd/gopackages)`, `go build ./cmd/api ./cmd/worker ./cmd/smoke`, `npm --prefix frontend run lint`, `npm --prefix frontend run check:sensitive-text`, `npm --prefix frontend run build`, `go run ./cmd/taskguard`, and `git diff --check`.
- 2026-04-25: Opened PR https://github.com/Chinsusu/Billing-V2/pull/287 for review and CI.
