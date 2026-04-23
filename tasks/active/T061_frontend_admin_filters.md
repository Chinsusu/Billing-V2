# T061 - Frontend admin filters

Status: IN_PROGRESS
Owner: Codex
Branch: feat/frontend-admin-filters
PR: -
Risk: frontend/API
Created: 2026-04-23
Updated: 2026-04-24

## Summary

Wire admin list screens to the backend search filters for practical support workflows.

## Scope

- Add filter controls for display ID, account/customer, status, and amount where supported.
- Connect filter state to existing API client query params.
- Preserve loading, error, empty, and mock fallback behavior.
- Keep UI consistent with the current Next.js/Tailwind architecture.

## Acceptance Criteria

- Admin invoices, transactions, logs, and other supported lists can query backend filters.
- Filter inputs are usable on desktop and mobile.
- Frontend build and audit gates pass.

## Notes

- Coordinate with T010/T011 owners if those branches are still active.

## Agent Log

- 2026-04-23: Task created after backend search filters merged.
- 2026-04-24: Codex claimed the task and started wiring shared admin filter controls to the live billing API queries.
