# T169 - Demo service VPS labels

Status: DONE
Owner: Codex
Branch: codex/t169-demo-service-vps-labels
PR: https://github.com/Chinsusu/Billing-V2/pull/369
Risk: frontend fallback display
Created: 2026-04-26
Updated: 2026-04-26

## Summary

Clean up demo service and VPS labels so client, reseller, and admin screens do not show internal VPS host-style names.

## Scope

- Replace raw demo service labels such as `vps-prod-01`, `vps-scrape-01`, `vps-scrape-02`, `vps-test`, `vps-api-gateway`, `vps-db-replica`, and `vps-worker-03`.
- Keep the change display-only in frontend mock data and smoke coverage.
- Do not split `frontend/src/mocks/billingData.ts`; that remains a follow-up task.

## Acceptance Criteria

- Service-facing demo screens show readable service labels instead of internal VPS names.
- Admin browser smoke guards against the raw VPS labels.
- Frontend lint, sensitive-text check, production build, admin smoke, taskguard, and diff check pass.

## Notes

- Follow-up to T168. This is limited to service/VPS demo labels and does not change backend API behavior.

## Agent Log

- 2026-04-26: Task created and claimed by Codex.
- 2026-04-26: Opened PR #369 after npm ci, frontend lint, sensitive-text guard, build, admin smoke, taskguard, and diff check passed.
- 2026-04-26: PR #369 merged into main.
