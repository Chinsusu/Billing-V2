# T184 - Audit demo provider target labels

Status: IN_PROGRESS
Owner: Codex
Branch: codex/t184-audit-demo-provider-target-labels
PR: -
Risk: frontend demo labels and smoke coverage
Created: 2026-04-28
Updated: 2026-04-28

## Summary

Humanize audit demo provider-source target labels that still expose raw source identifiers.

## Scope

- Display readable provider target labels for audit fallback rows backed by `SRC-23004` and `SRC-23001`.
- Keep mock filtering behavior backed by the original raw target values.
- Add audit fallback smoke coverage to reject those raw source IDs.
- Do not change live API contracts or backend behavior.

## Acceptance Criteria

- Audit demo fallback displays readable provider target labels instead of raw `SRC-*` source IDs.
- Audit fallback filtering still works from the original mock data.
- Admin smoke verifies the readable labels and rejects the raw source IDs.
- Frontend lint, sensitive-text check, production build, admin smoke, taskguard, and diff check pass locally.

## Notes

- This follows T180/T182 and covers remaining provider-source target labels in audit fallback UI.

## Agent Log

- 2026-04-28: Task created and claimed by Codex.
- 2026-04-28: Humanized audit demo provider-source targets and added smoke guards; local gates pass.
