# T105 - Provider readiness job context

Status: IN_PROGRESS
Owner: Codex
Branch: codex/t105-provider-readiness-job-context
PR: -
Risk: frontend/admin-ops
Created: 2026-04-24
Updated: 2026-04-24

## Summary

Connect provider readiness context to admin provisioning job inspection so failed or manual-review jobs show source readiness hints.

## Scope

- Reuse T100 readiness data in the admin provisioning/job detail area.
- Match job `source_id` to readiness rows when possible.
- Show a compact source readiness hint near job summary/timeline.
- Keep job attempts and audit payloads redacted.
- Do not change backend job recovery behavior.
- Keep each file under 500 lines.

## Acceptance Criteria

- Admin job detail can show readiness state/reason for the job source when data is available.
- Missing readiness data has a quiet fallback state.
- No provider credentials, raw provider payloads, or capability JSON reach the UI.
- Frontend and backend validation commands pass.
- Browser verification covers the job detail flow.

## Notes

- Follows T095, T099, and T101.

## Agent Log

- 2026-04-24: Task created in the provider readiness follow-up batch.
- 2026-04-24: Codex claimed the task on `codex/t105-provider-readiness-job-context`.
- 2026-04-24: Added admin job detail source readiness context matched by job provider source and order plan snapshot.
- 2026-04-24: Reused a shared provider readiness badge in provider readiness and job detail views.
- 2026-04-24: Browser verification passed on desktop and mobile with mocked live provisioning/readiness APIs; job detail showed ready source context without provider credentials, raw payload, or capability JSON text.
- 2026-04-24: Validation passed: `npm ci`, `npm audit --omit=dev`, `npm run lint`, `npm run build`, `go test ./...`, and `go build ./cmd/api ./cmd/migrate ./cmd/seed ./cmd/smoke ./cmd/worker`.
