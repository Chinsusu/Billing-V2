# T229 - Cloudmini pilot hardening

Status: TODO
Owner: -
Branch: codex/t229-cloudmini-pilot-hardening
PR: -
Risk: provider/provisioning/lifecycle/credential
Created: 2026-05-17
Updated: 2026-05-17

## Summary

Harden the Cloudmini V3 pilot path before broader provider use by resolving the T228 residual risks around provider-backed cleanup and terminal resource status.

## Scope

- Decide and implement or document how Billing should handle Cloudmini resources whose provider operation succeeds but `GET /api/v3/proxies/:id` still reports `status=creating`.
- Add tests or runbook evidence so Billing does not mark a service active before the resource is considered usable, unless the provider owner explicitly signs off that `creating` after operation success is acceptable.
- Add a provider-backed cleanup path for Cloudmini service termination, or document and guard an explicit direct-provider cleanup exception for pilot-only use.
- Cover success, provider failure, timeout/unknown, and cleanup failure paths with tests or documented operator evidence.
- Do not run another mutating provider pilot in this task unless a new one-resource approval is recorded.

## Acceptance Criteria

- The T228 `creating`-after-success status behavior is resolved by code, provider sign-off, or an explicit fail-closed rule.
- Billing cleanup no longer depends on an undocumented manual provider delete for broader pilot, or the exception is documented with owner approval and guardrails.
- Relevant provider/provisioning/lifecycle tests pass.
- Docs are updated with the chosen behavior and remaining risks.

## Notes

- T228 proved one controlled dev create/delete pilot but exposed these hardening gaps.
- Keep production and residential inventory out of scope.

## Agent Log

- 2026-05-17: Task created as follow-up from T228 controlled dev pilot residual risks.
