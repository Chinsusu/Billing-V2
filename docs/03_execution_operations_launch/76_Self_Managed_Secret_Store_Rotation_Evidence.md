# 76 - Self-Managed Secret Store Rotation Evidence

**Task:** T254  
**Date:** 2026-05-18  
**Scope:** Selected-host self-managed secret-store proof for Cloudmini and cloudflared secrets.  
**Decision:** approved for the current selected host and pilot scope only; provider-controlled error examples and broader provider approval remain launch blockers.

## Owner Record

| Role | Owner | Sign-off |
| --- | --- | --- |
| Security Owner | Admin | Accepted local-only self-managed secret-store boundary. |
| Ops Lead | Admin | Owns file permissions, process argv checks, and service restarts. |
| Provider Owner | Admin | Stated all secret/API keys were rotated before this evidence was recorded. |

Owner statement source: user message on 2026-05-18, "tất cả secret apikey tao rotate hết rồi".

## Approved Secret Store

| Secret area | Approved path | Metadata evidence | Notes |
| --- | --- | --- | --- |
| Cloudmini provider config/credential | `/etc/billing/secrets/cloudmini.env` | directory mode `700` owner `root:root`; file mode `600` owner `root:root` | Required Cloudmini keys present; `DB_DSN` absent; secret values were not printed. |
| Legacy dev Cloudmini credential | `/opt/cred-cloudmini-dev.env` | file mode `600` owner `root:root` | Retained as legacy dev/test source; canonical selected-host path is `/etc/billing/secrets/cloudmini.env`. |
| Cloudflared tunnel token | `/etc/cloudflared/tunnel.token` | file mode `600` owner `root:root` | Running `cloudflared` uses `--token-file`; exact `--token` arg absent. |

## Verification Output

```text
billing_secret_dir=700 root:root /etc/billing/secrets
cloudmini_secret_file=600 root:root /etc/billing/secrets/cloudmini.env
legacy_cloudmini_dev_file=600 root:root /opt/cred-cloudmini-dev.env
cloudflared_token_file=600 root:root /etc/cloudflared/tunnel.token
cloudmini_required_keys_present=yes
cloudmini_secret_contains_db_dsn=no
cloudflared_token_arg_present=no
cloudflared_token_file_arg_present=yes
secret_values_printed=no
```

## Required Handling

- Do not place provider credentials, app secrets, database DSNs, or tunnel tokens under `/opt/Billing`.
- Do not pass tokens, API keys, passwords, DSNs, or cookies through process argv.
- Load secrets through protected files, systemd environment files, or an approved future secret manager.
- Record only path, owner, mode, required-key presence, and process-argv checks.
- Repeat this evidence if the deployment moves to another host or if ownership/path/permissions change.

## Boundaries

This evidence does not approve:

- production customer data;
- new provider create/delete runs;
- provider-controlled permission-denied, rate-limited, out-of-capacity, provider-5xx, or cancel-rejected evidence gaps;
- increasing create limits, active-resource limits, worker concurrency, or provider rate limits;
- broader provider launch approval.

The launch decision remains NO-GO until the remaining P0 blockers in docs 69 and 70 are closed or explicitly accepted where policy allows.
