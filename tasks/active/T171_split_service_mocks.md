# T171 - Split service mock data

Status: IN_PROGRESS
Owner: Codex
Branch: codex/t171-split-service-mocks
PR: -
Risk: frontend mock data organization
Created: 2026-04-26
Updated: 2026-04-26

## Summary

Reduce `frontend/src/mocks/billingData.ts` size by moving service-related mock types and data into a focused module while preserving the existing import surface.

## Scope

- Move service/customer/catalog mock records used by service-facing screens out of `billingData.ts`.
- Keep `@/mocks/billingData` exports compatible for existing screens.
- Do not change demo values or frontend UI behavior.

## Acceptance Criteria

- `billingData.ts` is safely below the 500-line limit with service mock data split out.
- Existing imports from `@/mocks/billingData` continue to work.
- Frontend lint, sensitive-text check, production build, admin smoke, taskguard, and diff check pass.

## Notes

- Follow-up to T169/T170 after service demo label work pushed `billingData.ts` near the file-size threshold.

## Agent Log

- 2026-04-26: Task created and claimed by Codex.
- 2026-04-26: Split service/customer/reseller mock data into `serviceData.ts`; local gates pass.
