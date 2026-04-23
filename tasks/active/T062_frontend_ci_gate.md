# T062 - Frontend CI gate

Status: TODO
Owner: -
Branch: ci/frontend-quality-gate
PR: -
Risk: CI/frontend
Created: 2026-04-23
Updated: 2026-04-23

## Summary

Add CI coverage for the runnable frontend so dependency and build regressions fail before merge.

## Scope

- Add GitHub Actions steps for frontend dependency install, production dependency audit, lint, and build.
- Reuse the repo's existing Node/npm versions and scripts.
- Keep backend CI behavior intact.
- Document any required local frontend validation command if missing.

## Acceptance Criteria

- CI fails on frontend build or production audit failures.
- Existing backend test/build job still runs.
- Frontend checks are scoped so they do not read unrelated generated files.

## Notes

- This should wait until active frontend branches are not conflicting.

## Agent Log

- 2026-04-23: Task created after frontend app shell and API integration landed.
