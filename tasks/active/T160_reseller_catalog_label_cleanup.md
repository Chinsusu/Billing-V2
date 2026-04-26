# T160 - Reseller catalog label cleanup

Status: DONE
Owner: Codex
Branch: codex/t160-reseller-catalog-labels
PR: https://github.com/Chinsusu/Billing-V2/pull/351
Risk: frontend
Created: 2026-04-26
Updated: 2026-04-26

## Summary

Show readable source and status labels in the reseller catalog table instead of raw catalog keys.

## Scope

- Replace raw source keys such as catalog, master, ok, low, and out with readable labels.
- Use the shared status badge for reseller catalog status display.
- Keep catalog API values, clone payloads, and pricing behavior unchanged.

## Acceptance Criteria

- Reseller catalog Source column shows readable text such as Catalog, Master plan, Available, Low stock, or Out of stock.
- Reseller catalog Status column uses readable status labels through the shared badge.
- Frontend lint, sensitive-text check, production build, taskguard, and diff check pass.

## Notes

- This is frontend-only and should not change backend contracts.

## Agent Log

- 2026-04-26: Codex created and claimed the task after T159 was marked done; starting reseller catalog label cleanup.
- 2026-04-26: Replaced raw reseller catalog source/status rendering with readable source labels and the shared status badge.
- 2026-04-26: Validation passed: frontend lint, sensitive-text check, production build, admin browser smoke, and taskguard.
- 2026-04-26: Opened review PR https://github.com/Chinsusu/Billing-V2/pull/351.
- 2026-04-26: Merged PR https://github.com/Chinsusu/Billing-V2/pull/351 into main; marking task done.
