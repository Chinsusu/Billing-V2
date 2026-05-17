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
