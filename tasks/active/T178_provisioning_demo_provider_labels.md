# T178 - Provisioning demo provider labels

Status: IN_PROGRESS
Owner: Codex
Branch: codex/t178-demo-provider-job-labels
PR: -
Risk: frontend demo labels and smoke coverage
Created: 2026-04-28
Updated: 2026-04-28

## Summary

Humanize provisioning demo job provider names at the mock data source and cover the visible fallback label in admin smoke.

## Scope

- Replace raw provisioning demo provider names that look like internal keys.
- Add smoke coverage for the humanized provider label in provisioning fallback.
- Do not change live API contracts or backend behavior.

## Acceptance Criteria

- Provisioning demo fallback does not depend on raw provider names such as `proxy-cheap`.
- Admin smoke verifies the humanized provisioning fallback provider label.
- Frontend lint, sensitive-text check, production build, admin smoke, taskguard, and diff check pass locally.

## Notes

- T177 handled provider source fallback names; this task covers the provisioning demo job source data that references the same provider.

## Agent Log

- 2026-04-28: Task created and claimed by Codex.
- 2026-04-28: Humanized provisioning demo provider names and added smoke coverage for the fallback label; local gates pass.
