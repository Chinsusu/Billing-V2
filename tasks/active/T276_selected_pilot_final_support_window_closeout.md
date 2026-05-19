# T276 - Selected pilot final support-window closeout

Status: DONE
Owner: Codex
Branch: codex/t276-final-support-window-closeout
PR: https://github.com/Chinsusu/Billing-V2/pull/583
Risk: launch decision, ops, support, finance, secret handling, target environment
Created: 2026-05-19
Updated: 2026-05-19

## Summary

Record final selected bounded non-production pilot support-window closeout after `2026-05-19 22:00 Asia/Ho_Chi_Minh`.

## Scope

- Verify selected runtime services and public domains after the approved support window.
- Verify process command-line secret-pattern status and protected secret/token file metadata without printing secret contents.
- Re-run the read-only target finance reconciliation smoke.
- Check launch-critical notification counts without reading payloads.
- Record final support-window closeout evidence in the launch docs.
- Keep production, broader private beta, broader provider scope, and production customer data out of scope.

## Acceptance Criteria

- Local timestamp proves the check ran after `2026-05-19 22:00 Asia/Ho_Chi_Minh`.
- `billing-api`, `billing-frontend`, and `cloudflared` are active.
- Selected domains and backend health/readiness return HTTP `200`.
- Process command-line checks report no DSN/token/password/credential patterns without printing command lines.
- Protected secret/token file metadata remains restrictive without printing file contents.
- Finance reconciliation smoke passes with status `balanced` and zero mismatch counts.
- Launch-critical notification count is captured without reading payloads.
- Docs clearly state final support-window closeout is complete for the selected bounded non-production pilot.
- Task board remains consistent and required checks pass.

## Notes

- If health/readiness, finance, secret-pattern, or support coverage fails, mark the task `BLOCKED` and keep or move the selected pilot to paused state.
- This task does not approve production SMTP/Telegram delivery or broader provider/customer launch scope.

## Agent Log

- 2026-05-19: Task created and claimed by Codex from Billing `origin/main` after `22:00 Asia/Ho_Chi_Minh` for final support-window closeout.
- 2026-05-19: Runtime checks passed at `22:06 Asia/Ho_Chi_Minh`: `billing-api`, `billing-frontend`, and `cloudflared` active/enabled with main PIDs present.
- 2026-05-19: Process command-line secret-pattern checks reported `none` for `billing-api`, `billing-frontend`, and `cloudflared` without printing command lines.
- 2026-05-19: Protected metadata check passed: `/etc/billing/secrets` mode `700`, runtime env files mode `600`, and `/etc/cloudflared/tunnel.token` mode `600`, all owner `root:root`; file contents were not printed.
- 2026-05-19: Domain checks returned HTTP `200` for `billing.resvn.net`, `/backend/healthz`, `/backend/readyz`, `client.resvn.net`, and `reseller.resvn.net`.
- 2026-05-19: Read-only finance reconciliation smoke stayed `balanced` with wallets/invoices/payments checked `1/1/1`, zero mismatch counts, and no money or provider mutation routes called.
- 2026-05-19: Read-only notification summary found launch-critical notification total `0`; no payloads, customer data, DSNs, tokens, provider payloads, or credentials were read or recorded.
- 2026-05-19: Updated docs 69 and 70 to record final support-window closeout as complete for the selected bounded non-production pilot.
- 2026-05-19: Opened PR #583. Local checks passed: service/domain health, command-line secret-pattern checks, protected metadata checks, finance reconciliation smoke, read-only notification count query, `go run ./cmd/taskguard`, `git diff --check`, diff raw-secret pattern scan, and changed-file line-count check.
- 2026-05-19: PR #583 merged after GitHub checks passed. Marking T276 DONE in marker PR.
