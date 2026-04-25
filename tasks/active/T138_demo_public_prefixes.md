# T138 - Demo public ID prefixes

Status: REVIEW
Owner: Codex
Branch: codex/t138-demo-public-prefixes
PR: https://github.com/Chinsusu/Billing-V2/pull/307
Risk: frontend
Created: 2026-04-26
Updated: 2026-04-26

## Summary

Align visible demo and fallback IDs with the public ID prefixes used by the live frontend data.

## Scope

- Replace demo provider IDs that still use internal-looking `prv-*` values.
- Replace demo provisioning job IDs that still use lowercase `job-*` values.
- Replace demo transaction IDs that still use random `txn_*` values.
- Keep names, statuses, correlation IDs, and backend-only references unchanged.

## Acceptance Criteria

- Demo provider rows show `SRC-*` IDs.
- Demo provisioning job rows show `JOB-*` IDs.
- Demo transaction rows show `TX-*` IDs.
- Frontend lint, sensitive-text check, build, taskguard, and diff check pass.

## Notes

- This task only changes UI-facing demo identifiers. It does not change persisted backend UUIDs or integration IDs.

## Agent Log

- 2026-04-26: Codex created and claimed the task after T137 was marked done; starting demo public ID prefix cleanup.
- 2026-04-26: Replaced demo provider, provisioning job, and transaction IDs with `SRC-*`, `JOB-*`, and `TX-*` public prefixes.
- 2026-04-26: Local validation passed: `npm --prefix frontend run lint`, `npm --prefix frontend run check:sensitive-text`, `npm --prefix frontend run build`, `go run ./cmd/taskguard`, and `git diff --check`.
- 2026-04-26: Opened PR https://github.com/Chinsusu/Billing-V2/pull/307 for review.
