# T114 - Frontend admin smoke CI hardening

Status: IN_PROGRESS
Owner: Codex
Branch: codex/t114-frontend-admin-smoke-ci-hardening
PR: -
Risk: frontend/CI
Created: 2026-04-25
Updated: 2026-04-25

## Summary

Make the admin browser smoke flow easier to run in CI and less likely to fail because of local environment assumptions.

## Scope

- Review the current frontend admin smoke command and GitHub Actions workflow.
- Ensure the smoke command has clear prerequisites for browser installation and local server startup.
- Add CI wiring or a documented opt-in path if the current environment cannot run browser tests reliably.
- Avoid changing UI behavior unless required to make the smoke stable.
- Keep each edited file under 500 lines.

## Acceptance Criteria

- Agents can run one documented command locally for the admin browser smoke.
- CI either runs the smoke or clearly documents why it is intentionally manual/opt-in.
- The smoke command does not race with frontend build output.
- Existing frontend and backend checks pass.

## Notes

- This task builds on T107.
- Do not broaden the smoke to every portal unless it stays small and stable.

## Agent Log

- 2026-04-25: Task created in the board and delivery hardening batch.
- 2026-04-25: Codex claimed the task; hardening the admin browser smoke path for CI and local runs.
- 2026-04-25: Added a CI smoke script that runs against the standalone production artifact, added Playwright Chromium installation to the frontend CI gate, and documented local/CI smoke ordering.
