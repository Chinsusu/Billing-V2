# T155 - Remove unused sample data helper

Status: DONE
Owner: Codex
Branch: codex/t155-remove-unused-sample-data
PR: https://github.com/Chinsusu/Billing-V2/pull/341
Risk: frontend
Created: 2026-04-26
Updated: 2026-04-26

## Summary

Remove the unused mock sample data helper file after production helpers were moved out.

## Scope

- Delete `frontend/src/mocks/sampleData.ts` if it has no imports.
- Keep all billing mock records in `frontend/src/mocks/billingData.ts`.
- Verify frontend build and smoke still pass.

## Acceptance Criteria

- No frontend code imports `sampleData`.
- The unused sample data helper file is removed.
- Frontend lint, sensitive-text check, smoke, build, taskguard, and diff check pass.

## Notes

- This follows T153 and T154; those tasks moved status and money helpers into production modules.

## Agent Log

- 2026-04-26: Codex created and claimed the task after T154 was marked done; starting unused sample data cleanup.
- 2026-04-26: Removed the unused sample data helper file after confirming no frontend imports remain.
- 2026-04-26: Validation passed: frontend lint, sensitive-text check, production build, admin browser smoke, taskguard, and diff check.
- 2026-04-26: Opened review PR https://github.com/Chinsusu/Billing-V2/pull/341.
- 2026-04-26: Merged PR https://github.com/Chinsusu/Billing-V2/pull/341 into main; marking task done.
