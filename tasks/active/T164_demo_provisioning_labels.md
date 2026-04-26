# T164 - Demo provisioning labels

Status: REVIEW
Owner: Codex
Branch: codex/t164-demo-provisioning-labels
PR: https://github.com/Chinsusu/Billing-V2/pull/359
Risk: frontend fallback display
Created: 2026-04-26
Updated: 2026-04-26

## Summary

Clean up provisioning demo fallback labels so operators do not see raw error keys when the live jobs API is unavailable.

## Scope

- Humanize demo provisioning error labels.
- Keep fallback provider names readable when mock data contains dashed keys.
- Add browser smoke coverage for the jobs API fallback path.

## Acceptance Criteria

- Demo fallback shows readable error text such as `Provider Timeout: Resource State Unknown`.
- Raw demo keys such as `provider_timeout`, `auth_failed`, `partial_success`, and `external_id` are not visible.
- Admin browser smoke covers the fallback path.

## Notes

- This is limited to the provisioning fallback table; live API rendering is already covered separately.

## Agent Log

- 2026-04-26: Task created and claimed by Codex.
- 2026-04-26: Opened PR #359 after lint, sensitive-text guard, build, admin smoke, taskguard, and diff check passed.
