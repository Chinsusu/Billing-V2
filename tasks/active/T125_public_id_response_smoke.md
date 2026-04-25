# T125 - Public ID response smoke coverage

Status: REVIEW
Owner: Codex
Branch: codex/t125-public-id-response-smoke
PR: https://github.com/Chinsusu/Billing-V2/pull/281
Risk: API/frontend
Created: 2026-04-25
Updated: 2026-04-25

## Summary

Add smoke coverage for related public display IDs returned by admin APIs and rendered by admin UI screens.

## Scope

- Extend API smoke checks for related public ID response fields added in T124.
- Extend admin browser smoke mocks so live UI paths display related public IDs.
- Assert the admin UI shows public labels for linked records and does not fall back to backend references in high-value screens.
- Keep the task focused on coverage; do not add new API response fields.

## Acceptance Criteria

- Smoke checks cover at least invoices, transactions, services, top-up requests, provisioning jobs, and audit logs.
- Admin browser smoke verifies linked public ID labels are visible in at least invoice, transaction, service, top-up, provisioning, and audit screens.
- Existing frontend and backend gates pass.

## Notes

- This task follows T124 and protects the new response contract from accidental regressions.

## Agent Log

- 2026-04-25: Codex created and claimed the task after T124 merged; starting smoke coverage for related public IDs.
- 2026-04-25: Added API smoke assertions, dev seed job/audit records, and admin browser smoke checks for related public ID labels.
- 2026-04-25: Validation passed: `go test ./cmd/smoke ./internal/seed`, full Go package tests, Go build, frontend lint/build/sensitive-text, `npm run smoke:admin`, taskguard, and diff check.
- 2026-04-25: Opened PR #281 for review.
