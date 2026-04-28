# T181 - Alert demo source labels

Status: IN_PROGRESS
Owner: Codex
Branch: codex/t181-alert-demo-source-labels
PR: -
Risk: frontend demo labels and smoke coverage
Created: 2026-04-28
Updated: 2026-04-28

## Summary

Humanize alert demo source labels that still expose provider source or migration identifiers in visible admin alerts.

## Scope

- Replace raw alert source references such as `SRC-23001`, `SRC-23005`, and `DB migration 0003`.
- Add admin smoke coverage to reject those raw alert labels.
- Do not change live API contracts or backend behavior.

## Acceptance Criteria

- Admin alert demo copy uses human-readable source and migration labels.
- Admin smoke rejects raw alert source identifiers.
- Frontend lint, sensitive-text check, production build, admin smoke, taskguard, and diff check pass locally.

## Notes

- T180 handled audit demo source labels; this task covers alert demo source labels.

## Agent Log

- 2026-04-28: Task created and claimed by Codex.
- 2026-04-28: Humanized alert demo source labels and added admin alert smoke guards; local gates pass.
