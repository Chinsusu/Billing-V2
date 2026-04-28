# T180 - Audit demo source labels

Status: IN_PROGRESS
Owner: Codex
Branch: codex/t180-audit-demo-source-labels
PR: -
Risk: frontend demo labels and smoke coverage
Created: 2026-04-28
Updated: 2026-04-28

## Summary

Humanize audit demo actor and migration detail labels at the mock data source while preserving admin smoke coverage against raw backend-style labels.

## Scope

- Replace raw audit demo actor names such as `billing-worker`, `prov-worker`, and `health-worker`.
- Replace raw audit migration detail `0003_rbac` with human-readable text.
- Add smoke coverage for the humanized audit migration detail.
- Do not change live API contracts or backend behavior.

## Acceptance Criteria

- Audit demo source data no longer stores raw worker actor names or `0003_rbac`.
- Admin smoke still rejects raw audit fallback labels and verifies humanized audit text.
- Frontend lint, sensitive-text check, production build, admin smoke, taskguard, and diff check pass locally.

## Notes

- T179 handled provisioning demo error and trace values; this task covers audit demo source values.

## Agent Log

- 2026-04-28: Task created and claimed by Codex.
- 2026-04-28: Humanized audit demo source actor and migration labels; local gates pass.
