# T215 - Cloudmini edge unblock runbook

Status: REVIEW
Owner: Codex
Branch: codex/t215-cloudmini-edge-unblock
PR: https://github.com/Chinsusu/Billing-V2/pull/461
Risk: provider/provisioning/credential/config
Created: 2026-05-16
Updated: 2026-05-16

## Summary

Document the Cloudmini edge/gateway unblock requirements and safe read-only rerun procedure after T214 found provider edge HTTP `403` responses.

## Scope

- Record the Cloudmini edge/gateway access requirements needed before another authenticated `/api/v3` read-only check.
- Document safe rerun boundaries for capability and inventory checks after access is fixed.
- Document credential handling and token rotation expectations for the chat-shared provider credential.
- Preserve `NO-GO` for real sandbox pilot until read-only checks, mapping, quota, cleanup, idempotency, and pilot evidence are complete.
- Do not call Cloudmini, create/delete provider resources, or change runtime code in this task.

## Acceptance Criteria

- Provider evidence docs include a clear unblock checklist for Cloudflare/gateway/provider owner action.
- Rerun procedure avoids query-string credentials and raw response logging unless explicitly approved as an exception.
- Launch evidence packet references edge/gateway unblock proof as a required real-provider sandbox evidence item.
- Task guard and whitespace checks pass.

## Notes

- `/opt/proxy-cloudmini` code read shows `/api/v3` is behind Cloudflare Tunnel and auth supports bearer plus API-key headers. T214 showed requests did not reach a successful V3 app envelope.

## Agent Log

- 2026-05-16: Task created and claimed on `codex/t215-cloudmini-edge-unblock`.
- 2026-05-16: Validation passed: `go run ./cmd/taskguard`; `git diff --check`; generic token-pattern scan returned no matches in `/opt/Billing`.
- 2026-05-16: Opened PR https://github.com/Chinsusu/Billing-V2/pull/461 and moved task to `REVIEW`.
