# T153 - Move status badge labels out of mocks

Status: DONE
Owner: Codex
Branch: codex/t153-status-badge-label-source
PR: https://github.com/Chinsusu/Billing-V2/pull/337
Risk: frontend
Created: 2026-04-26
Updated: 2026-04-26

## Summary

Make status badge labels a production display helper instead of mock data.

## Scope

- Move status labels and badge variants into the shared display label helper module.
- Update `StatusBadge` to use the production helper.
- Remove unused status label exports from mock sample data.
- Keep visible status text and colors unchanged.

## Acceptance Criteria

- `StatusBadge` does not import from mock data.
- Existing status labels such as Manual Review, Under review, Posted, and Retryable keep rendering.
- Frontend lint, sensitive-text check, smoke, build, taskguard, and diff check pass.

## Notes

- This is frontend-only and should not change routes or API payloads.

## Agent Log

- 2026-04-26: Codex created and claimed the task after T152 was marked done; starting status badge label source cleanup.
- 2026-04-26: Moved status labels and variants into the shared display helper module and updated StatusBadge to use it.
- 2026-04-26: Validation passed: frontend lint, sensitive-text check, production build, admin browser smoke, taskguard, and diff check.
- 2026-04-26: Opened review PR https://github.com/Chinsusu/Billing-V2/pull/337.
- 2026-04-26: Merged PR https://github.com/Chinsusu/Billing-V2/pull/337 into main; marking task done.
