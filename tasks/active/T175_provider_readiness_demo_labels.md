# T175 - Provider readiness demo labels

Status: IN_PROGRESS
Owner: Codex
Branch: codex/t175-provider-readiness-demo-labels
PR: -
Risk: frontend demo fallback labels
Created: 2026-04-27
Updated: 2026-04-27

## Summary

Prevent raw provider readiness demo plan codes such as `vps-linux-small`, `proxy-residential`, and `proxy-dc-shared` from being rendered in the admin UI.

## Scope

- Prefer human-readable provider readiness plan names or formatted technical labels in the UI.
- Add smoke coverage for provider readiness demo fallback labels.
- Do not change provider readiness API contracts or backend behavior.

## Acceptance Criteria

- Provider readiness demo fallback shows human-readable plan labels.
- Raw demo plan codes are covered by admin smoke forbidden-text checks.
- Frontend lint, sensitive-text check, production build, admin smoke, taskguard, and diff check pass.

## Notes

- Follow-up after T174 while continuing demo fallback label cleanup.

## Agent Log

- 2026-04-27: Task created and claimed by Codex.
- 2026-04-27: Humanized provider readiness plan labels and added smoke fallback coverage; local gates pass.
