# T216 - Cloudmini dev read-only evidence

Status: DONE
Owner: Codex
Branch: codex/t216-cloudmini-dev-evidence
PR: https://github.com/Chinsusu/Billing-V2/pull/463
Risk: provider/provisioning/credential/config
Created: 2026-05-16
Updated: 2026-05-16

## Summary

Record the Cloudmini V3 dev read-only evidence collected from the local dev credential source without storing secrets or running mutating provider calls.

## Scope

- Use the local dev credential source only for read-only Cloudmini V3 checks.
- Record status codes, envelope success, feature keys, inventory counts, and remaining blockers.
- Keep raw provider credentials, raw auth headers, provider-private IDs, and raw response bodies out of git.
- Do not SSH to provider hosts or run create/delete/action provider routes.
- Keep pilot readiness blocked until approved shared secret storage, source mapping, quota, cleanup, idempotency, and pilot evidence exist.

## Acceptance Criteria

- Provider evidence docs show the Billing Go-client-style read-only path succeeds through `https://cz.resvn.net/`.
- Launch completion packet distinguishes read-only reachability from real pilot readiness.
- Task guard, whitespace, and secret-pattern checks pass.

## Notes

- `/opt/cred` was treated as a local dev-only credential source. Production entries in that file were not used.
- The checker used a Billing Go-client-style user-agent because provider-side evidence reported Cloudflare blocks generic `Python-urllib/3.12` with HTTP `403` code `1010`.

## Agent Log

- 2026-05-16: Task created and claimed on `codex/t216-cloudmini-dev-evidence`.
- 2026-05-16: Ran unauthenticated `GET /api/v3/capabilities`; result was HTTP `401` JSON from the app in `711ms`.
- 2026-05-16: Ran read-only authenticated checks for capabilities plus `ipv4_dc` and `residential` inventory using `Authorization`, `X-API-Key`, and `X-ACCESS-CODE`; all nine authenticated checks returned HTTP `200` V3 success envelopes. No mutating provider routes were called.
- 2026-05-16: Validation passed: `go run ./cmd/taskguard`; `git diff --check`; file length check; secret-pattern scan against changed files. Opened PR https://github.com/Chinsusu/Billing-V2/pull/463 and moved task to `REVIEW`.
- 2026-05-16: PR https://github.com/Chinsusu/Billing-V2/pull/463 merged; marking task `DONE`.
