# T232 - Cloudmini lifecycle cleanup activation

Status: DONE
Owner: Codex
Branch: codex/t232-cloudmini-lifecycle-activation
PR: https://github.com/Chinsusu/Billing-V2/pull/497, https://github.com/Chinsusu/Billing-V2/pull/499
Risk: provider/provisioning/lifecycle/credential/ops
Created: 2026-05-17
Updated: 2026-05-17

## Summary

Run a one-resource, owner-approved Cloudmini lifecycle cleanup activation on the approved dev test server, proving the lifecycle worker can call the real Cloudmini provider cleanup path before marking a Billing service terminated.

## Scope

- Run all non-mutating preflight checks before any create/delete route.
- Use at most one Cloudmini dev resource and one lifecycle worker claim.
- Keep regular worker service in the safe fake-provider mode outside the bounded activation window.
- Record only redacted evidence in provider/launch docs.
- Stop without mutation if queue state, mapping, credential, or service/resource preconditions are not clean.

## Acceptance Criteria

- Preflight confirms `APP_ENV=dev`, private credential path, Cloudmini registry activation, ready mapping, and no unrelated queued provider/lifecycle jobs.
- Exactly one approved test resource is available for lifecycle cleanup.
- `cmd/worker lifecycle-once` runs with `PROVIDER_DEFAULT_MODE=cloudmini_v3`, batch size `1`, and claims no more than one lifecycle job.
- Evidence shows provider cleanup reached a terminal safe state and Billing did not print raw DSNs, tokens, raw provider IDs, raw provider payloads, or proxy credentials.
- Broader pilot remains blocked unless duplicate/timeout, shared secret-store, and named owner gates are also complete.

## Notes

- Initial T232 activation was blocked on 2026-05-17 because Cloudmini returned provider resource status `creating`; T233 added the bounded wait/read policy needed to keep fail-closed behavior while allowing resources that become usable to activate.
- After T233 merged and was deployed to the approved test server, the owner-approved rerun created one active Cloudmini-backed Billing service and cleaned it up through `cmd/worker lifecycle-once` with the real Cloudmini registry.
- Broader pilot remains blocked by shared secret-store, named owners, live duplicate/timeout evidence, redacted provider error examples, and launch sign-off gaps.

## Agent Log

- 2026-05-17: Task created and claimed by Codex from latest `origin/main` on branch `codex/t232-cloudmini-lifecycle-activation`.
- 2026-05-17: Ran mutating preflight on the approved dev test server; result `PASS` with one-resource guardrails, private credential path, and redacted mapping evidence.
- 2026-05-17: Stopped the always-on fake worker, created one Billing dev order/invoice/payment, and ran `cmd/worker provision-once` with `PROVIDER_DEFAULT_MODE=cloudmini_v3`, batch size `1`.
- 2026-05-17: Provisioning worker claimed exactly one job and returned `manual_review` with error code `PROVIDER_PARTIAL_SUCCESS` because provider status was `creating`; no active service was created.
- 2026-05-17: Found the resource by provider `external_ref`, cleaned it up through V3 `DELETE`, verified final provider `GET` returned HTTP `404`, and restarted `billing-worker`.
- 2026-05-17: Opened PR #497 with blocked evidence and T233 follow-up task.
- 2026-05-17: T233 target rerun passed: `provision-once` claimed one job and succeeded, service display `10002` became active with one encrypted active credential, `lifecycle-once` claimed one terminate job and succeeded, final provider `GET` returned HTTP `404`, and `billing-worker` was restored active.
