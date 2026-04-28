# T179 - Provisioning demo error labels

Status: DONE
Owner: Codex
Branch: codex/t179-demo-provisioning-error-labels
PR: https://github.com/Chinsusu/Billing-V2/pull/389
Risk: frontend demo labels and smoke coverage
Created: 2026-04-28
Updated: 2026-04-28

## Summary

Humanize provisioning demo job error and trace labels at the mock data source while preserving admin smoke coverage against raw backend-style codes.

## Scope

- Replace raw provisioning demo error values such as `provider_timeout`, `auth_failed`, and `partial_success`.
- Replace raw demo correlation identifiers with a non-trace placeholder.
- Preserve smoke coverage that rejects the raw backend-style values in the visible fallback UI.
- Do not change live API contracts or backend behavior.

## Acceptance Criteria

- Provisioning demo source data no longer stores raw error codes or `cor_` trace identifiers.
- Admin smoke still verifies the humanized provisioning fallback error labels.
- Frontend lint, sensitive-text check, production build, admin smoke, taskguard, and diff check pass locally.

## Notes

- T178 handled provisioning demo provider names; this task covers the remaining provisioning demo error and trace strings.

## Agent Log

- 2026-04-28: Task created and claimed by Codex.
- 2026-04-28: Humanized provisioning demo error and trace source values; local gates pass.
- 2026-04-28: Opened PR #389 for review.
- 2026-04-28: PR #389 merged into `main`; task marked DONE.
