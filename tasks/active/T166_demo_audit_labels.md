# T166 - Demo audit labels

Status: DONE
Owner: Codex
Branch: codex/t166-demo-audit-labels
PR: https://github.com/Chinsusu/Billing-V2/pull/363
Risk: frontend fallback display
Created: 2026-04-26
Updated: 2026-04-26

## Summary

Clean up audit log demo fallback labels so operators do not see raw worker keys or underscore-heavy detail text when the live audit API is unavailable.

## Scope

- Humanize demo audit actor names.
- Replace raw technical fragments in demo audit details.
- Add browser smoke coverage for the audit API fallback path.

## Acceptance Criteria

- Demo fallback shows readable actor names such as `Provisioning Worker`.
- Demo fallback shows `manual review threshold exceeded` instead of `manual_review threshold exceeded`.
- Raw keys such as `prov-worker`, `billing-worker`, `health-worker`, `manual_review`, and `0003_rbac` are not visible.

## Notes

- Live audit rows already use API view model labels.

## Agent Log

- 2026-04-26: Task created and claimed by Codex.
- 2026-04-26: Opened PR #363 after lint, sensitive-text guard, build, admin smoke, taskguard, and diff check passed.
- 2026-04-26: PR #363 merged into main.
