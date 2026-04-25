# T123 - Public ID smoke coverage

Status: DONE
Owner: Codex
Branch: codex/t123-public-id-smoke-coverage
PR: https://github.com/Chinsusu/Billing-V2/pull/277
Risk: smoke/API/frontend
Created: 2026-04-25
Updated: 2026-04-25

## Summary

Add smoke coverage that proves public numeric IDs work in admin API/UI flows and raw backend references are not used as labels.

## Scope

- Extend backend smoke or frontend browser smoke with focused public ID checks.
- Cover at least one admin list filter by numeric display ID.
- Cover at least one frontend admin screen that should show public IDs.
- Keep smoke output concise and safe for logs.
- Avoid making smoke depend on production data or real provider accounts.

## Acceptance Criteria

- Smoke fails when a public ID filter stops working.
- Smoke fails when a protected raw backend reference appears in an admin label covered by the test.
- Smoke remains deterministic with local seed or mocked data.
- Existing validation command matrix remains accurate.

## Notes

- Prefer focused checks over broad UI crawling.

## Agent Log

- 2026-04-25: Task created in the public ID and validation hardening batch.
- 2026-04-25: Codex claimed the task after T122 merged; reviewing smoke fixtures and browser/API coverage for public ID filters.
- 2026-04-25: Added API smoke positive/miss checks for admin public ID filters and browser smoke coverage for the invoice public ID filter UI.
- 2026-04-25: Local validation passed: frontend smoke, frontend build/lint/sensitive-text guard, cmd/smoke tests, Go package tests, Go command builds, taskguard, and diff check.
- 2026-04-25: Opened PR #277 for review.
- 2026-04-25: PR #277 merged into main.
