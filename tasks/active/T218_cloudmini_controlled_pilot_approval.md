# T218 - Cloudmini controlled pilot approval

Status: REVIEW
Owner: Codex
Branch: codex/t218-cloudmini-pilot-approval
PR: https://github.com/Chinsusu/Billing-V2/pull/466
Risk: provider/provisioning/credential/config
Created: 2026-05-16
Updated: 2026-05-16

## Summary

Prepare the controlled Cloudmini V3 mutating pilot approval packet with mapping, quota, cleanup, and stop conditions before any create/delete call is allowed.

## Scope

- Record the Cloudmini V3 pilot mapping candidate using redacted provider references only.
- Define quota, concurrency, spend, cleanup, and owner approval requirements.
- Keep raw API tokens, provider-private group IDs, proxy credentials, and raw provider payloads out of git.
- Do not run `POST /api/v3/proxies`, `DELETE /api/v3/proxies/:id`, provider action routes, or Billing worker pilot in this task.
- Keep T217 as the separate runtime follow-up for multiple Cloudmini V3 endpoint/API-key mappings.

## Acceptance Criteria

- A controlled pilot runbook exists with explicit preflight, approval, stop, and cleanup gates.
- The selected dev mapping uses a sellable `ipv4_dc` inventory group from read-only evidence and records only a redacted group reference.
- Docs keep pilot status blocked until approval and one create/delete run are completed.
- Task guard, whitespace, and secret-pattern checks pass.

## Notes

- `/opt/cred-cloudmini-dev.env` was updated as the local-only dev secret/config file with raw provider mapping values and mode `0600`.
- Existing seeded source `Local Fake Hetzner Ready` is not a Cloudmini source and must not be reused as the Cloudmini pilot source.

## Agent Log

- 2026-05-16: Task created and claimed on `codex/t218-cloudmini-pilot-approval`.
- 2026-05-16: Read-only inventory selected one sellable `ipv4_dc` group with redacted ref `c6a7189f0a` and `200` allocatable units. Raw group id stayed only in `/opt/cred-cloudmini-dev.env`.
- 2026-05-16: Validation passed: `go run ./cmd/taskguard`; `git diff --check`; file length check; secret-pattern scan against changed files. Opened PR https://github.com/Chinsusu/Billing-V2/pull/466 and moved task to `REVIEW`.
