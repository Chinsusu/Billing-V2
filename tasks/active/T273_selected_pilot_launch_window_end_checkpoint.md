# T273 - Selected pilot launch-window end checkpoint

Status: REVIEW
Owner: Codex
Branch: codex/t273-launch-window-end-checkpoint
PR: https://github.com/Chinsusu/Billing-V2/pull/577
Risk: launch decision, ops, support, finance, secret handling, target environment
Created: 2026-05-19
Updated: 2026-05-19

## Summary

Record the selected bounded non-production pilot checkpoint at the end of the launch window and before final support-window closeout.

## Scope

- Verify selected runtime services and public domains at the launch-window boundary.
- Verify process command-line secret-pattern status and protected secret/token file metadata without printing secret contents.
- Re-run the read-only target finance reconciliation smoke.
- Check launch-critical notification counts without reading payloads.
- Record that the launch window ended but final support-window closeout remains due after `2026-05-19 22:00 Asia/Ho_Chi_Minh`.
- Keep production, broader private beta, broader provider scope, and production customer data out of scope.

## Acceptance Criteria

- `billing-api`, `billing-frontend`, and `cloudflared` are active.
- Selected domains and backend health/readiness return HTTP `200`.
- Process command-line checks report no DSN/token/password/credential patterns without printing command lines.
- Protected secret/token file metadata remains restrictive without printing file contents.
- Finance reconciliation smoke passes with status `balanced` and zero mismatch counts.
- Launch-critical notification count is captured without reading payloads.
- Docs clearly state final support-window closeout remains due after `2026-05-19 22:00 Asia/Ho_Chi_Minh`.
- Task board remains consistent and required checks pass.

## Notes

- This task can record the end of the `18:00-20:00 Asia/Ho_Chi_Minh` launch window, but it cannot truthfully close the support window before `22:00`.
- If health/readiness, finance, secret-pattern, or support coverage fails, mark the task `BLOCKED` and keep or move the selected pilot to paused state.

## Agent Log

- 2026-05-19: Task created and claimed by Codex from Billing `origin/main` at `20:00 Asia/Ho_Chi_Minh`; final support-window closeout remains due after `22:00 Asia/Ho_Chi_Minh`.
- 2026-05-19: Runtime checks passed at `20:01 Asia/Ho_Chi_Minh`: `billing-api`, `billing-frontend`, and `cloudflared` active/enabled with main PIDs present.
- 2026-05-19: Process command-line secret-pattern checks reported `none` for `billing-api`, `billing-frontend`, and `cloudflared` without printing command lines.
- 2026-05-19: Protected metadata check passed: `/etc/billing/secrets` mode `700`, runtime env files mode `600`, and `/etc/cloudflared/tunnel.token` mode `600`, all owner `root:root`; file contents were not printed.
- 2026-05-19: Domain checks returned HTTP `200` for `billing.resvn.net`, `/backend/healthz`, `/backend/readyz`, `client.resvn.net`, and `reseller.resvn.net`.
- 2026-05-19: Read-only finance reconciliation smoke stayed `balanced` with wallets/invoices/payments checked `1/1/1`, zero mismatch counts, and no money or provider mutation routes called.
- 2026-05-19: Read-only notification summary found launch-critical notification total `0`; no payloads, customer data, DSNs, tokens, provider payloads, or credentials were read or recorded.
- 2026-05-19: Updated docs 69 and 70 with the selected launch-window end checkpoint evidence and explicit note that final support-window closeout remains due after `22:00 Asia/Ho_Chi_Minh`.
- 2026-05-19: Opened PR #577. Local checks passed: service/domain health, command-line secret-pattern checks, protected metadata checks, finance reconciliation smoke, read-only notification count query, `go run ./cmd/taskguard`, `git diff --cached --check`, staged raw-secret pattern scan, and changed-file line-count check.
