# 89 - Broader Private Beta Target Preflight Evidence

**Date:** 2026-05-21
**Scope:** Read-only target preflight evidence for the current broader private beta launch-candidate domains and selected test server.
**Decision:** PASS for current target health/runtime/process-secret/secret-file metadata checks only. Broader private beta remains `NO-GO`.

## Boundary

This evidence covers the current launch-candidate public domains and the selected test server only.

Do not broaden this result:

- No production database, production provider account, production notification channel, production customer data, or customer list was used.
- No raw DB DSN, password, cookie, session token, provider token, Telegram token, TOTP value, private key, provider payload, notification payload, proxy credential, or customer data is recorded here.
- No money mutation route, provider mutation route, Cloudmini create/delete/action route, credential reveal route, login smoke, or notification delivery route was called by this task.
- This evidence does not approve broader private beta, production launch, provider expansion, production customer data, or Telegram as a production primary path.

## Public Domain Health

Public checks were run from the repo host against the current launch-candidate domains.

```text
billing.resvn.net /backend/healthz: 200
client.resvn.net /backend/healthz: 200
reseller.resvn.net /backend/healthz: 200
billing.resvn.net /: 200
client.resvn.net /: 200
reseller.resvn.net /: 200
```

## Target Runtime State

Target checks were run over SSH on the selected test server using metadata-only commands.

```text
host_role=test-server
billing-api.service=active
billing-frontend.service=active
cloudflared.service=active
postgresql@14-main.service=active
127.0.0.1:8080 /healthz: 200
127.0.0.1:3000 /: 200
```

Worker-specific proof was not rerun in this task. Worker behavior remains covered by the relevant E2E, provider, and notification evidence packets for their specific scope.

## Secret Metadata Remediation

Initial metadata check found the target notification secret directory/file more permissive than the launch-candidate standard:

```text
/etc/billing/secrets mode=755 owner=root group=root
/etc/billing/secrets/telegram.env mode=644 owner=root group=root
```

Remediation applied on the selected test server:

```text
/etc/billing/secrets owner=root group=billing-svc mode=750
/etc/billing/secrets/telegram.env owner=root group=billing-svc mode=640
```

Post-remediation service and health checks passed:

```text
billing-api.service=active
billing-frontend.service=active
cloudflared.service=active
127.0.0.1:8080 /healthz: 200
127.0.0.1:3000 /: 200
billing.resvn.net /backend/healthz: 200
client.resvn.net /backend/healthz: 200
reseller.resvn.net /backend/healthz: 200
billing.resvn.net /: 200
client.resvn.net /: 200
reseller.resvn.net /: 200
```

Additional secret metadata:

```text
/opt/Billing/.env.dev owner=root group=billing-svc mode=640
/etc/cloudflared/tunnel.token owner=root group=root mode=600
```

No secret file contents were printed or copied.

## Process Argv Secret-Value Check

The target process scan used a strict check for raw secret-value flags or secret assignments in argv.

```text
strict_cmdline_secret_pattern_matches=0
cloudflared_unit_has_token_file=yes
cloudflared_unit_raw_jwt_like_matches=0
cloudflared_raw_token_value_on_argv=no
```

Note: a broad substring scan can match the safe `--token-file` flag name. This evidence uses the strict scan result above for raw token value exposure.

## Not Verified

- Owner approval for broader private beta v1 was not collected.
- Customer list or data classification was not approved.
- Launch window and final target ownership were not approved.
- Full UAT, full E2E, renewal, credential reveal, auth/RBAC rerun, and finance reconciliation were not rerun.
- Provider quota/SKU/group mapping, timeout/idempotency, cleanup, and pilot run evidence were not expanded.
- Notification primary/fallback SLA, escalation, delivery, and failure drill were not rerun for broader private beta.
- Production launch, production customer data, and broader provider scope remain unapproved.

## Decision

```text
Public domain health: PASS
Target runtime health: PASS
Secret metadata remediation: PASS
Strict process argv secret-value check: PASS
Raw secret values printed: no
Money mutation routes called: no
Provider mutation routes called: no
Notification delivery routes called: no
Broader private beta GO: not approved
Production GO: not approved
```
