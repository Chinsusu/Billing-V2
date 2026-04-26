# T161 - Ticket priority labels

Status: DONE
Owner: Codex
Branch: codex/t161-ticket-priority-labels
PR: https://github.com/Chinsusu/Billing-V2/pull/353
Risk: frontend
Created: 2026-04-26
Updated: 2026-04-26

## Summary

Show readable support ticket priority labels in admin and reseller ticket tables.

## Scope

- Add a shared frontend display helper for ticket priority values.
- Use the helper in admin and reseller support ticket screens.
- Keep source mock/API values unchanged.

## Acceptance Criteria

- Ticket tables show High, Medium, and Low instead of raw values such as high, medium, and low.
- Existing priority coloring still applies.
- Frontend lint, sensitive-text check, production build, taskguard, and diff check pass.

## Notes

- This is frontend-only and should not change backend contracts.

## Agent Log

- 2026-04-26: Codex created and claimed the task after T160 was marked done; starting ticket priority label cleanup.
- 2026-04-26: Added shared ticket priority label helper and applied it to admin and reseller support ticket tables.
- 2026-04-26: Validation passed: frontend lint, sensitive-text check, production build, admin browser smoke, and taskguard.
- 2026-04-26: Opened review PR https://github.com/Chinsusu/Billing-V2/pull/353.
- 2026-04-26: Merged PR https://github.com/Chinsusu/Billing-V2/pull/353 into main; marking task done.
