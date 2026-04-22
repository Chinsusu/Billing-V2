# T014 - API Contracts

Status: IN_PROGRESS
Owner: Codex
Branch: feat/api-contracts
PR: -
Risk: API
Created: 2026-04-23
Updated: 2026-04-23

## Summary

Add shared API response helpers for success lists, error details, validation fields, and cursor pagination.

## Scope

- Extend the existing response envelope with optional `page`, `details`, and validation `fields`.
- Add cursor pagination request parsing and response page helper.
- Keep current health/readiness response behavior compatible.
- Keep auth, tenant routes, OpenAPI generation, and frontend wiring out of scope.

## Acceptance Criteria

- Existing success/error response helpers remain compatible.
- List responses include `page.limit` and nullable `page.next_cursor`.
- Validation errors return `validation.failed` and field-level details.
- Cursor pagination enforces default and max limits.
- `make fmt`, `make test`, `make build`, and `make migrate-validate` pass.

## Notes

- T014 follows `docs/05_development_standards/50_API_Response_Error_Logging_Standard.md`.

## Agent Log

- 2026-04-23: Codex claimed task from `origin/main` using isolated worktree `/tmp/Billing-T014`.
