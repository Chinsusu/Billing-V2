# T044 - Top-up request schema API

Status: DONE
Owner: Codex
Branch: feat/topup-request-schema-api
PR: https://github.com/Chinsusu/Billing-V2/pull/101
Risk: wallet/topup/API
Created: 2026-04-23
Updated: 2026-04-23

## Summary

Add top-up request records and basic read/create APIs so manual wallet funding can be tracked before approval.

## Scope

- Add top-up request domain models and PostgreSQL schema.
- Add client create/list/detail APIs for current account.
- Add admin list/detail APIs for tenant-scoped review.
- Add numeric display IDs for top-up requests.
- Out of scope: approval, wallet crediting, file upload, or gateway integration.

## Acceptance Criteria

- Top-up requests support draft/submitted/under_review/approved/rejected/expired/cancelled statuses.
- Client writes use tenant and actor from context, not request body.
- Admin reads require wallet top-up review/read permission.
- Migration validation and focused handler/store tests pass.
- Full validation passes: `make fmt`, `make test`, `make build`, `make migrate-validate`.

## Notes

- This task can start after wallet schema exists.

## Agent Log

- 2026-04-23: Task created for the next backend wallet/invoice batch.
- 2026-04-23: Implemented top-up request schema, store/read service methods, client/admin APIs, rollback runbook, and focused tests.
- 2026-04-23: PR #101 passed checks and was merged.
