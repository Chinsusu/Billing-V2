# T228 - Cloudmini controlled mutating pilot

Status: IN_PROGRESS
Owner: Codex
Branch: codex/t228-cloudmini-controlled-pilot
PR: -
Risk: provider/provisioning/credential/wallet/order/database/ops
Created: 2026-05-17
Updated: 2026-05-17

## Summary

Run the first approved Cloudmini V3 mutating pilot through the Billing checkout/provisioning path on the non-production dev environment, with one-create guardrails, same-session cleanup, and redacted evidence only.

## Scope

- Deploy or verify latest `origin/main` on the approved non-production test server before the pilot.
- Re-run local quality checks, Cloudmini read-only mapping evidence, and T226 mutating preflight.
- Run at most one Billing checkout/provisioning pilot for `proxy-static-10gb-monthly` using the mapped `cloudmini_v3` source.
- Capture redacted evidence for order/job/source/service/provider state, idempotency presence, encrypted credential storage, and cleanup.
- Cleanup the created provider resource in the same session through the approved Billing/service path when available, or the documented approved provider delete path if Billing cleanup is unavailable.
- Update launch/provider evidence docs with redacted results and remaining risks.
- Do not broaden to production, residential inventory, multi-provider rollout, or more than one active provider resource.

## Acceptance Criteria

- Dev environment is confirmed non-production before any provider mutating route is enabled or called.
- T226 preflight passes with exact one-create, one-active-resource, one-worker-concurrency guardrails.
- The pilot creates no more than one provider resource and records no raw token, DSN, provider group id, raw provider payload, or proxy credential in repo evidence.
- Billing records show the checkout/order/provisioning job/service path, or the task records the precise blocker before any mutating provider call if the Billing path cannot safely run.
- Cleanup is attempted in the same session and the final provider/Billing state is recorded redacted.
- Required local checks, task guard, and diff checks pass before PR.

## Notes

- User approval to proceed with the controlled dev pilot was given in chat on 2026-05-17. Runtime guardrails and preflight must still pass before any mutating call.
- Stop immediately on auth/permission failure, gateway/rate-limit block, timeout-after-send, duplicate resource suspicion, manual review, cleanup failure, or any secret exposure.

## Agent Log

- 2026-05-17: Task created and claimed by Codex from latest `origin/main`.
- 2026-05-17: Local validation before pilot passed: `go test ./...`; `go run ./cmd/taskguard`.
- 2026-05-17: Test server prechecks passed with `APP_ENV=dev`, T226 preflight `PASS`, provider `proxy_crud` read precheck HTTP `200`/`success=true`, Cloudmini active Billing services `0`, queued provider jobs `0`, selected provider group sellable with allocatable `200`, active proxy count `0`, pending create count `0`, and reserved count `0`.
- 2026-05-17: Ran controlled one-resource pilot `T228-dev-20260517T004039Z`: stopped default worker loop, created order `10001`, invoice `10002`, wallet payment transaction `10001`, ledger `10002`, provisioning job `10001`, source `10012`; one-off worker with `PROVIDER_DEFAULT_MODE=cloudmini_v3` returned `claimed=1 succeeded=1 retried=0 manual_review=0 terminal_failed=0 cancelled=0`.
- 2026-05-17: Pilot created service `10001` with status `active`, billing status `paid`, one encrypted active credential with masked hint, and redacted provider resource ref `redacted:dc3d9457bf5b`. No raw provider token, DSN, group id, provider response, or proxy credential was recorded in repo evidence.
- 2026-05-17: Cleanup completed in the same session: provider V3 delete operation reached `succeeded`, provider GET after cleanup returned HTTP `404`, Billing reseller service terminate returned service status `terminated`, Cloudmini active Billing services returned to `0`, selected provider group returned to allocatable `200`, active proxy count `0`, pending create count `0`, reserved count `0`, and the default `billing-worker` service was restarted.
- 2026-05-17: Residual risks recorded in docs: Billing service terminate does not call provider delete, provider GET immediately after create returned resource status `creating` despite operation/worker success, and duplicate-create plus timeout-after-send evidence remain missing.
