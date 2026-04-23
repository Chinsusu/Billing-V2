# T055 - Frontend API integration

Status: REVIEW
Owner: Codex
Branch: feat/frontend-api-integration
PR: https://github.com/Chinsusu/Billing-V2/pull/125
Risk: frontend/API
Created: 2026-04-23
Updated: 2026-04-23

## Summary

Connect billing and admin frontend screens to real backend API clients using the agreed Next.js App Router, React, TypeScript, and Tailwind stack.

## Scope

- Add typed API client modules for wallet, invoice, payment, reconciliation, audit, order, and service reads.
- Wire existing frontend screens to backend data behind local/dev configuration.
- Add loading, empty, and error states that do not break layout.
- Keep mock/reference data separated from production API paths.
- Avoid adding route behavior that is not backed by existing backend APIs.

## Acceptance Criteria

- `npm run build` passes in `frontend`.
- Screens can read seeded local data when the backend is running.
- Shared API/client helpers are separated from page components.
- No frontend file exceeds 500 lines.

## Notes

- UI reference artifacts are not the architecture source of truth.
- This task should follow T053 and T054 if it needs verified local data.

## Agent Log

- 2026-04-23: Task created for frontend/backend integration.
- 2026-04-23: Claimed by Codex; adding typed API client, hooks, and live-data wiring with mock fallback.
- 2026-04-23: Opened PR #125. Validation passed: `npm ci`, `npm run build`, `npm audit --omit=dev`, backend gates, and Playwright live API probe against local seed data.
