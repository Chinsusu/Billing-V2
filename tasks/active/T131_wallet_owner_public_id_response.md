# T131 - Wallet owner public ID response

Status: REVIEW
Owner: Codex
Branch: codex/t131-wallet-owner-public-id
PR: https://github.com/Chinsusu/Billing-V2/pull/293
Risk: API/frontend
Created: 2026-04-26
Updated: 2026-04-26

## Summary

Expose wallet owner public IDs so admin/reseller/client wallet reads can show account labels without relying on backend owner references.

## Scope

- Add `owner_display_id` to wallet read responses when the owner is a user in the same tenant.
- Update frontend wallet types and reseller wallet lookup logic to use owner public IDs where available.
- Extend API smoke checks and wallet tests for owner public ID responses.
- Update the API reference wallet resource shape.

## Acceptance Criteria

- Client/admin/reseller wallet list/detail responses include `owner_display_id` when available.
- Reseller client/dashboard wallet joins prefer `owner_display_id` over `owner_id`.
- Existing wallet tests, smoke tests, frontend build, taskguard, and diff check pass.

## Notes

- `owner_id` remains in API responses for actions/internal callers, but UI lookups should prefer `owner_display_id`.

## Agent Log

- 2026-04-26: Codex created and claimed the task after T130 was marked done; starting wallet owner public ID response support.
- 2026-04-26: Codex implemented `owner_display_id` in wallet read responses, frontend wallet types, reseller wallet lookups, API reference, smoke checks, and wallet tests. Local validation passed: `go test $(go run ./cmd/gopackages)`, `go build ./cmd/api ./cmd/worker ./cmd/smoke`, `npm --prefix frontend run lint`, `npm --prefix frontend run check:sensitive-text`, `npm --prefix frontend run build`, `go run ./cmd/taskguard`, and `git diff --check`.
- 2026-04-26: Opened PR https://github.com/Chinsusu/Billing-V2/pull/293 for review and CI.
