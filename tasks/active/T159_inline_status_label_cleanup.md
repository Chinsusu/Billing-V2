# T159 - Inline status label cleanup

Status: IN_PROGRESS
Owner: Codex
Branch: codex/t159-inline-status-labels
PR: -
Risk: frontend
Created: 2026-04-26
Updated: 2026-04-26

## Summary

Show readable labels for inline status text that appears outside the shared status badge component.

## Scope

- Use shared status labels for wallet summary status text in client and reseller views.
- Use shared status labels for live top-up activity text in the admin overview feed.
- Keep status values, filters, API payloads, and `StatusBadge` behavior unchanged.

## Acceptance Criteria

- Inline wallet summaries show labels such as Active instead of raw keys like active.
- Admin overview top-up activity shows labels such as Under review instead of raw keys like under_review.
- Frontend lint, sensitive-text check, production build, taskguard, and diff check pass.

## Notes

- This is frontend-only and should not change backend contracts.

## Agent Log

- 2026-04-26: Codex created and claimed the task after T158 was marked done; starting inline status label cleanup.
- 2026-04-26: Applied shared status labels to admin top-up activity and client/reseller wallet summaries; added overview smoke coverage for readable top-up status.
- 2026-04-26: Validation passed: frontend lint, sensitive-text check, production build, admin browser smoke, and taskguard.
