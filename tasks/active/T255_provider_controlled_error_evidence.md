# T255 - Cloudmini provider-controlled error evidence

Status: DONE
Owner: Codex
Branch: codex/t255-provider-controlled-errors
PR: https://github.com/Chinsusu/Billing-V2/pull/542
Risk: provider provisioning, credentials, launch-readiness evidence
Created: 2026-05-18
Updated: 2026-05-18

## Summary

Record safe, redacted evidence or owner-approved deferrals for remaining Cloudmini provider-controlled error cases.

## Scope

- Inspect existing Cloudmini error smoke and adapter error taxonomy coverage.
- Prefer non-mutating or provider-owner-controlled simulation paths for permission denied, rate limited, out-of-capacity, provider 5xx, and cancel/delete rejected.
- If a case cannot be safely triggered in the current provider environment, record a precise deferral and required provider-side support instead of guessing.
- Keep all secret values, raw provider IDs, raw payloads, proxy credentials, cookies, and DSNs out of repo evidence.
- Preserve NO-GO unless all provider-controlled cases and broader provider owner approval are closed.

## Acceptance Criteria

- Provider evidence docs list the status of each remaining provider-controlled error case.
- Any executed command is redacted and records only stable code/category/retry-safety metadata.
- No new provider create/delete run is performed without an explicit owner-approved bounded command path.
- `go run ./cmd/taskguard` passes.
- `git diff --check` passes.

## Notes

- This task may close only the error-evidence blocker if safe evidence exists; broader provider owner approval remains separate.

## Agent Log

- 2026-05-18: Task created and claimed by Codex for Cloudmini provider-controlled error evidence.
- 2026-05-18: Inspected Billing Cloudmini smoke/adapter code and provider V3 route/error source under `/opt/proxy-cloudmini` without reading secrets.
- 2026-05-18: Added doc 77 and updated provider launch evidence docs with a case-by-case safe execution plan; no new provider create/delete/reservation/API-key mutation was run.
- 2026-05-18: Updated docs index/manifest so the new provider-controlled error evidence packet is discoverable.
- 2026-05-18: Opened PR #542 and moved task to REVIEW after `go run ./cmd/taskguard`, `git diff --check`, line-count check, and secret-like diff scan passed.
- 2026-05-18: PR #542 passed CI and merged; marking task DONE.
