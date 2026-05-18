# T260 - Cloudmini rate-limit runtime evidence

Status: DONE
Owner: Codex
Branch: codex/t260-cloudmini-rate-limit-runtime
PR: https://github.com/Chinsusu/Billing-V2/pull/552
Risk: provider provisioning, credentials, launch-readiness evidence
Created: 2026-05-18
Updated: 2026-05-18
Merged: 2026-05-18

## Summary

Deploy or verify the merged Cloudmini V3 rate-limit fixture on the approved dev manager, then collect one redacted Billing runtime evidence run for `RATE_LIMITED` if the fixture is safely reachable.

## Scope

- Verify the approved dev manager is running provider code with Cloudmini PR #9 and `VPM_BILLING_ERROR_FIXTURES_ENABLED=yes`.
- Run only the guarded Billing `cloudmini-error-evidence` smoke for the side-effect-free rate-limit fixture path.
- Record redacted evidence if the run passes.
- If deployment or fixture activation is unavailable, record the exact blocker without broad provider mutation.
- Do not create/delete proxies, do not trigger the shared limiter, and do not print raw provider payloads, provider credentials, tokens, cookies, DSNs, proxy credentials, or raw provider IDs.

## Acceptance Criteria

- Runtime evidence closes the rate-limit case with HTTP `429`, provider code `RATE_LIMITED`, normalized `PROVIDER_RATE_LIMITED`, retry safety `safe_retry`, and `mutating_routes_called=false`; or the task records a concrete deployment/activation blocker.
- Docs 66, 69, 70, and 77 reflect the new evidence or remaining blocker.
- Required local checks pass before PR.

## Notes

- Provider PR #9 is already merged. This task must not claim runtime evidence unless the approved dev manager is actually running that code with the fixture env enabled.
- Provider 5xx and cancel/delete rejected fixtures remain separate follow-up work unless safely available without expanding this task.

## Agent Log

- 2026-05-18: Task created and claimed by Codex from Billing `origin/main`.
- 2026-05-18: Verified the Billing target server does not run the Cloudmini manager; identified the approved Cloudmini manager host separately without printing secret contents.
- 2026-05-18: Initial guarded rate-limit smoke returned `FAIL` because the fixture path returned HTTP `200 text/html`, proving the provider manager was not yet serving the fixture as a V3 error envelope.
- 2026-05-18: Fast-forwarded the provider manager source from `f7a7ab2` to `7a63532` while preserving pre-existing dirty files; remote build was blocked by Go `1.18.1` not parsing provider `go 1.25.0`.
- 2026-05-18: Built the manager binary from a live-source copy with Go `1.26.2`, ran targeted provider API tests, backed up the existing manager binary, deployed the new binary, enabled `VPM_BILLING_ERROR_FIXTURES_ENABLED=yes` only for the evidence run, and restarted `vpm-manager`.
- 2026-05-18: Guarded Billing smoke passed with HTTP `429`, provider code `RATE_LIMITED`, normalized `PROVIDER_RATE_LIMITED`, retry safety `safe_retry`, one fixture request, and `mutating_routes_called=false`.
- 2026-05-18: Removed the fixture env drop-in after evidence collection, restarted `vpm-manager`, verified the service active, verified fixture env absent, and verified public/local capabilities still return `401 application/json`.
- 2026-05-18: Opened Billing PR #552 and moved task to REVIEW.
- 2026-05-18: Billing PR #552 merged and task marked DONE.
