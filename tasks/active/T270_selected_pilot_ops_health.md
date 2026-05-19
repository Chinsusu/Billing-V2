# T270 - Selected pilot ops health

Status: DONE
Owner: Codex
Branch: codex/t270-launch-ops-health
PR: https://github.com/Chinsusu/Billing-V2/pull/571
Risk: launch decision, ops, support, secret handling, target environment
Created: 2026-05-19
Updated: 2026-05-19

## Summary

Record selected bounded non-production pilot launch-window ops health and manual fallback coverage evidence.

## Scope

- Verify selected runtime services are active during the T268 launch/support window.
- Verify selected public domains and backend health/readiness return HTTP `200`.
- Verify process command lines and protected env/token file metadata without printing secret contents.
- Record manual fallback support coverage evidence for the selected support window.
- Keep production, broader private beta, broader provider scope, and production customer data out of scope.

## Acceptance Criteria

- `billing-api`, `billing-frontend`, and `cloudflared` are active.
- `billing.resvn.net`, `billing.resvn.net/backend/healthz`, `billing.resvn.net/backend/readyz`, `client.resvn.net`, and `reseller.resvn.net` return HTTP `200`.
- Process command-line checks report no DSN/token/password/credential patterns without printing command lines.
- Protected secret/token file metadata remains restrictive without printing file contents.
- Manual fallback coverage is recorded for the selected support window plus the continued pause-on-SLA-breach rule.
- Task board remains consistent and required checks pass.

## Notes

- If health/readiness fails, or secret-like command-line patterns are detected, mark the task `BLOCKED` and keep or move the selected pilot to paused state.
- No provider, money, credential reveal, or customer-data mutation is in scope.

## Agent Log

- 2026-05-19: Task created and claimed by Codex from Billing `origin/main` during the selected pilot launch/support window.
- 2026-05-19: Runtime checks passed: `billing-api`, `billing-frontend`, and `cloudflared` active/enabled with main PIDs present.
- 2026-05-19: Process command-line secret-pattern checks reported `none` for `billing-api`, `billing-frontend`, and `cloudflared` without printing command lines.
- 2026-05-19: Protected metadata check passed: `/etc/billing/secrets` mode `700`, runtime env files mode `600`, and `/etc/cloudflared/tunnel.token` mode `600`, all owner `root:root`; file contents were not printed.
- 2026-05-19: Domain checks returned HTTP `200` for `billing.resvn.net`, `/backend/healthz`, `/backend/readyz`, `client.resvn.net`, and `reseller.resvn.net`.
- 2026-05-19: Read-only notification summary found launch-critical notification total `0`; no payloads, customer data, DSNs, tokens, provider payloads, or credentials were read or recorded.
- 2026-05-19: Updated docs 69 and 70 with selected launch-window ops health and manual fallback coverage evidence.
- 2026-05-19: Opened PR #571. Local checks passed: service status, command-line secret-pattern checks, protected metadata checks, domain health/readiness checks, read-only notification count query, `go run ./cmd/taskguard`, `git diff --check`, raw-secret pattern scan, and changed-file line-count check.
- 2026-05-19: PR #571 merged after GitHub checks passed. Marking task DONE in marker PR.
