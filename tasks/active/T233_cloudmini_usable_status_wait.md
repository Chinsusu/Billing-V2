# T233 - Cloudmini usable status wait policy

Status: TODO
Owner: Unassigned
Branch: codex/t233-cloudmini-usable-status-wait
PR: -
Risk: provider/provisioning/credential/ops
Created: 2026-05-17
Updated: 2026-05-17

## Summary

Add an approved Cloudmini V3 provisioning wait/read policy so Billing can handle provider resources that report `creating` immediately after operation success without weakening T229 fail-closed behavior.

## Scope

- Keep `creating`, empty, unknown, and other non-usable statuses fail-closed if they remain non-usable after the bounded wait.
- Poll/read the created proxy status after a successful Cloudmini create operation until it becomes one of `running`, `active`, `ready`, or `available`, or until the configured timeout expires.
- Preserve encrypted credential handling and redacted errors.
- Add focused provider adapter tests for status becoming usable, status staying non-usable, timeout, and credential-missing behavior.
- After merge/deploy, rerun the one-resource T232 lifecycle activation flow in a new owner-approved window.

## Acceptance Criteria

- Cloudmini create can produce an active Billing service only after the provider resource reaches a usable status and credential fields are present.
- A resource that stays `creating` or otherwise non-usable still returns manual review and does not create an active service.
- Adapter tests cover the wait/read policy and fail-closed timeout.
- No raw provider IDs, payloads, tokens, DSNs, or proxy credentials are logged or committed.
- T232 remains blocked until the target test-server activation rerun proves lifecycle-worker cleanup.

## Notes

- Do not solve this by manually inserting service records or broadening acceptable statuses.
- If provider semantics say `creating` should be considered billable/usable, record explicit Provider Owner and Security/Engineering approval before changing the usable-status list.

## Agent Log

- 2026-05-17: Follow-up created from T232 after the approved dev activation attempt reached Cloudmini create but provider status stayed `creating`, blocking active service creation and lifecycle cleanup.
