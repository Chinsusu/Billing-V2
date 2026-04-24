# T105 - Provider readiness job context

Status: TODO
Owner: -
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
