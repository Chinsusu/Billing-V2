# T132 - Ledger reference public ID response

Status: REVIEW
Owner: Codex
Branch: codex/t132-ledger-reference-public-id
PR: https://github.com/Chinsusu/Billing-V2/pull/295
Risk: API/frontend
Created: 2026-04-26
Updated: 2026-04-26

## Summary

Expose related public IDs on wallet ledger entries so wallet history can label invoice/top-up references without backend UUIDs.

## Scope

- Add `reference_display_id` to ledger read responses for supported reference types.
- Preserve ledger create/post scans that use the original `RETURNING` column list.
- Update frontend API types and client wallet reference labels.
- Extend smoke checks, wallet read tests, and API reference docs.

## Acceptance Criteria

- Client wallet ledger responses include `reference_display_id` for seeded top-up and invoice entries.
- Client wallet UI uses related public ID labels such as `TUP-` or `INV-` when available.
- Existing wallet tests, smoke tests, frontend build, taskguard, and diff check pass.

## Notes

- `reference_id` remains available for internal callers and backend joins.

## Agent Log

- 2026-04-26: Codex created and claimed the task after T131 was marked done; starting ledger reference public ID response support.
- 2026-04-26: Codex implemented `reference_display_id` in ledger read responses, frontend wallet labels, API reference, smoke checks, and wallet tests while preserving ledger create/post scans. Local validation passed: `go test $(go run ./cmd/gopackages)`, `go build ./cmd/api ./cmd/worker ./cmd/smoke`, `npm --prefix frontend run lint`, `npm --prefix frontend run check:sensitive-text`, `npm --prefix frontend run build`, `go run ./cmd/taskguard`, and `git diff --check`.
- 2026-04-26: Opened PR https://github.com/Chinsusu/Billing-V2/pull/295 for review and CI.
