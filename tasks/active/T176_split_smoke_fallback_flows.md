# T176 - Split smoke fallback flows

Status: DONE
Owner: Codex
Branch: codex/t176-split-smoke-fallback-flows
PR: https://github.com/Chinsusu/Billing-V2/pull/383
Risk: frontend smoke test organization
Created: 2026-04-27
Updated: 2026-04-27

## Summary

Reduce `frontend/scripts/smoke-admin-browser.cjs` file-size risk by moving admin fallback smoke flows into a focused module without changing smoke coverage.

## Scope

- Move fallback flow helpers out of `smoke-admin-browser.cjs`.
- Preserve existing admin smoke assertions and command behavior.
- Do not change frontend UI or mock API payload behavior.

## Acceptance Criteria

- `smoke-admin-browser.cjs` is safely below the 500-line file limit.
- New smoke helper files stay below 500 lines.
- Admin smoke still passes.
- Frontend lint, sensitive-text check, production build, taskguard, and diff check pass.

## Notes

- Follow-up after T175 left `smoke-admin-browser.cjs` at 493 lines.

## Agent Log

- 2026-04-27: Task created and claimed by Codex.
- 2026-04-27: Split admin fallback smoke flows into `smoke-admin-fallback-flows.cjs`; local gates pass.
- 2026-04-27: Opened PR #383 for review.
- 2026-04-27: PR #383 merged into `main`; task marked DONE.
