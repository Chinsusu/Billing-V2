# T177 - Provider demo labels

Status: REVIEW
Owner: Codex
Branch: codex/t177-provider-demo-labels
PR: https://github.com/Chinsusu/Billing-V2/pull/385
Risk: frontend demo labels and smoke coverage
Created: 2026-04-27
Updated: 2026-04-27

## Summary

Humanize provider source demo fallback names and add admin smoke coverage for the provider sources fallback path.

## Scope

- Replace raw provider demo names that look like internal keys.
- Add browser smoke coverage for `/backend/admin/catalog/provider-sources` fallback behavior.
- Do not change live API contracts or backend behavior.

## Acceptance Criteria

- Provider demo fallback does not show raw names such as `proxy-manager` or `proxy-cheap`.
- Admin smoke covers provider source fallback labels.
- Frontend lint, sensitive-text check, production build, admin smoke, taskguard, and diff check pass.

## Notes

- T176 split fallback smoke flows into a dedicated helper module; reuse that module for this coverage.

## Agent Log

- 2026-04-27: Task created and claimed by Codex.
- 2026-04-27: Humanized provider source demo names and added provider source fallback smoke coverage; local gates pass.
- 2026-04-27: Opened PR #385 for review.
