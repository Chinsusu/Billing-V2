# T172 - Split admin smoke script

Status: IN_PROGRESS
Owner: Codex
Branch: codex/t172-split-admin-smoke
PR: -
Risk: frontend smoke test organization
Created: 2026-04-27
Updated: 2026-04-27

## Summary

Reduce `frontend/scripts/smoke-admin-browser.cjs` below the repository 500-line file limit by moving stable smoke fixtures/helpers into focused modules without changing smoke behavior.

## Scope

- Split mock API fixture/server setup or flow helpers out of `smoke-admin-browser.cjs`.
- Keep `npm --prefix frontend run smoke:admin:ci` behavior and coverage intact.
- Do not change UI behavior or production frontend code unless needed to preserve smoke coverage.

## Acceptance Criteria

- `frontend/scripts/smoke-admin-browser.cjs` is below 500 lines.
- Any new script files are below 500 lines.
- Admin smoke still covers the existing forbidden-text checks.
- Frontend lint, sensitive-text check, production build, admin smoke, taskguard, and diff check pass.

## Notes

- Follow-up to T171 after the smoke script was identified at 612 lines.

## Agent Log

- 2026-04-27: Task created and claimed by Codex.
- 2026-04-27: Split API mock routing into `smoke-admin-api-mocks.cjs`; all local gates pass.
