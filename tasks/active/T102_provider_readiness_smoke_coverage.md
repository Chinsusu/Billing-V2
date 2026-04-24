# T102 - Provider readiness smoke coverage

Status: REVIEW
Owner: Codex
Branch: codex/t102-provider-readiness-smoke-coverage
PR: https://github.com/Chinsusu/Billing-V2/pull/232
Risk: QA/smoke
Created: 2026-04-24
Updated: 2026-04-24

## Summary

Add smoke coverage that verifies the provider readiness endpoint is reachable and safe before local paid-order provisioning smoke runs.

## Scope

- Extend the local smoke path or add a focused smoke helper for provider readiness.
- Assert the response shape includes display IDs, state, and reason.
- Assert the response does not include provider credentials, raw provider payloads, or capability JSON.
- Keep local fake provider support intact.
- Keep each file under 500 lines.

## Acceptance Criteria

- Smoke validation calls `GET /admin/catalog/provider-readiness` with admin headers.
- Smoke output is concise and uses display IDs for human-readable diagnostics.
- Failure messages are redacted and actionable.
- Backend and frontend validation commands pass.

## Notes

- Follows T100 and should stay local/sandbox friendly.
- Do not require production provider credentials.

## Agent Log

- 2026-04-24: Task created in the provider readiness follow-up batch.
- 2026-04-24: Codex claimed the task on `codex/t102-provider-readiness-smoke-coverage`.
- 2026-04-24: Added admin provider readiness smoke coverage with display-ID summary output and redacted blocked-field checks.
- 2026-04-24: Opened PR #232 for review.
