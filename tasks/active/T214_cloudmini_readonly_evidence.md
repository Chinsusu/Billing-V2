# T214 - Cloudmini V3 read-only evidence

Status: IN_PROGRESS
Owner: Codex
Branch: codex/t214-cloudmini-readonly-evidence
PR: -
Risk: provider/provisioning/credential/config
Created: 2026-05-15
Updated: 2026-05-15

## Summary

Capture redacted authenticated read-only Cloudmini V3 sandbox evidence without storing credentials or mutating provider resources.

## Scope

- Use the provided Cloudmini V3 credential only as transient process input; do not write it to git, docs, task notes, PR text, shell command text, or command output.
- Call only read-only Cloudmini V3 endpoints: capabilities and inventory groups.
- Record status codes, high-level capability keys, inventory counts, and stock-state summaries with provider identifiers redacted.
- Keep real sandbox pilot readiness blocked until owner, quota, source/group mapping, cleanup, idempotency, and create/delete pilot evidence exist.
- Do not call create, delete, action, status-by-resource, or any Billing order/provisioning flow in this task.

## Acceptance Criteria

- Authenticated `GET /api/v3/capabilities` succeeds or records a redacted failure.
- Authenticated inventory checks for `ipv4_dc` and `residential` succeed or record redacted failures.
- No raw provider response body, token, auth header, proxy credential, or provider-private identifier is committed.
- Provider evidence docs and launch packet reflect the authenticated read-only result while preserving `NO-GO` for pilot readiness.
- Task guard and whitespace checks pass.

## Notes

- The credential was provided in chat. Treat it as exposed for shared environments and rotate before staging/pilot use if provider policy requires it.
- This task cannot prove create/delete idempotency, timeout-after-send behavior, cleanup, or Billing provisioning activation.

## Agent Log

- 2026-05-15: Task created and claimed on `codex/t214-cloudmini-readonly-evidence`.
- 2026-05-15: Bearer authenticated read-only checks returned HTTP `403` for capabilities, `ipv4_dc` inventory, and `residential` inventory; no mutating calls were made and no raw response bodies were stored.
- 2026-05-15: Header variants `X-API-Key` and `X-ACCESS-CODE` were tested read-only. `X-API-Key` returned HTTP `403` for all three read endpoints. `X-ACCESS-CODE` returned HTTP `403` for both inventory endpoints and timed out once for capabilities after `20795ms`.
- 2026-05-15: A schema-only inspection of one HTTP `403` response showed Cloudflare/gateway-style keys, not the expected V3 app envelope. Provider edge/gateway access approval remains required before inventory/capability mapping can proceed.
- 2026-05-15: Validation passed: `go run ./cmd/taskguard`; `git diff --check`; generic token-pattern scan returned no matches in `/opt/Billing`.
