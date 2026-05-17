# T240 - Target secret/key handling evidence

Status: IN_PROGRESS
Owner: Codex
Branch: codex/t240-target-secret-key-evidence
PR: -
Risk: secrets/config/security/ops
Created: 2026-05-17
Updated: 2026-05-17

## Summary

Capture target-environment secret/key handling evidence and remove the visible cloudflared tunnel token from process argv.

## Scope

- Verify target dev/test secret-bearing files are outside git and have restricted file permissions.
- Move cloudflared from token-on-argv usage to token-file usage without printing the token.
- Restart cloudflared and verify target local/domain HTTP reachability.
- Update launch evidence docs with redacted proof.
- Do not print, commit, rotate, or disclose raw DSNs, tokens, API keys, passwords, cookies, or credential payloads.

## Acceptance Criteria

- Cloudflared process no longer exposes `--token` with a token value in argv.
- Cloudflared uses `/etc/cloudflared/tunnel.token` through `--token-file`.
- Token file and Cloudmini dev credential file permissions are recorded with owner/mode only.
- Billing local frontend and configured tunnel domains return HTTP `200`.
- Launch docs remain honest that Security owner sign-off and approved shared secret-store evidence are still required before GO.

## Notes

- This task records target dev/test evidence only. It does not approve production secret handling or replace Security owner sign-off.

## Agent Log

- 2026-05-17: Task created and claimed by Codex on branch `codex/t240-target-secret-key-evidence`.
- 2026-05-17: Target pre-check found `/opt/Billing/.env.dev` mode `640` owner `root:billing-svc`, `/opt/cred-cloudmini-dev.env` mode `600` owner `root:root`, `/etc/cloudflared/tunnel.token` mode `600` owner `root:root`, services active, but cloudflared still exposed token usage through process argv.
- 2026-05-17: Updated target cloudflared systemd service from token flag usage to `--token-file /etc/cloudflared/tunnel.token`, ran daemon reload, restarted cloudflared, and verified cloudflared active with `cloudflared_token_in_argv=no`.
- 2026-05-17: Target reachability check passed with HTTP `200` for `http://localhost:3000`, `https://billing.resvn.net`, `https://reseller.resvn.net`, and `https://client.resvn.net`.
