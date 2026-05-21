# 85 - Domain-Aware Target Auth Smoke Evidence

**Date:** 2026-05-21
**Scope:** Post-T290 selected test-server deploy and domain-aware target auth/RBAC smoke evidence.
**Decision:** PASS for selected non-production test-server domain-aware auth/RBAC smoke only. This does not approve production, production customer data, broader private beta, real-provider provisioning, or production notification delivery.

## Boundary

This evidence covers the selected test server reached from `/opt/Billing-Dev-Server` and the public non-production domains already mapped to that host.

Do not broaden this result:

- No production database or production customer data was used.
- No raw DB DSN, password, cookie, session token, provider token, provider payload, Telegram token, TOTP value, or plaintext service credential is recorded here.
- No money mutation route, provider mutation route, or Cloudmini create/delete/action route was called.
- The evidence records status codes, boolean checks, service states, and redacted outcomes only.

## Deploy Actions

Merged source deployed to the selected test server:

```text
source_commit=407ed13
feature_merge=PR #612
marker_merge=PR #613
target_path=/opt/Billing
deploy_method=rsync deploy-copy without .git
remote_build=PASS
remote_smoke_binary=/opt/Billing/bin/smoke
```

Runtime status after deploy:

```text
billing-api.service=active
billing-frontend.service=active
cloudflared.service=active
postgresql@14-main.service=active
127.0.0.1:8080 /healthz: 200
127.0.0.1:3000 /: 200
billing.resvn.net /backend/healthz: 200
client.resvn.net /backend/healthz: 200
reseller.resvn.net /backend/healthz: 200
```

Secret handling metadata:

```text
/opt/Billing/.env.dev mode=640 owner=root group=billing-svc
target_auth_credential_values_loaded_by_allowlist=yes
secret_values_printed=no
cmdline_secret_pattern_matches=0
```

## Domain-Aware Auth/RBAC Evidence

The remote smoke binary was executed on the selected test server with separate public API base URLs for client and platform admin login:

```text
command=dev-target-auth-rbac
base_url=http://127.0.0.1:8080
client_base_url=https://client.resvn.net/backend
admin_base_url=https://billing.resvn.net/backend
result=PASS
client_session_cookie_only=pass
admin_2fa_gate=pass
invalid_session_denied=pass
actor_required_denied=pass
tenant_mismatch_denied=pass
rbac_denials=3
domain_aware_base_urls=pass
provider_mutation_routes_called=no
money_mutation_routes_called=no
```

Smoke output explicitly reported that raw session tokens, cookies, passwords, DSNs, provider payloads, and credentials were excluded.

## Not Verified

- Human Client, Reseller, Admin, or Security sign-off was not collected in this task.
- Full browser UAT was not rerun.
- Full 2FA-satisfied admin browser navigation was not executed.
- Real provider provisioning, money mutations, and production notification delivery were not executed.
- Production launch is not approved by this evidence.

## Decision

```text
Selected test-server deploy-copy: PASS
Remote smoke binary build: PASS
Domain-aware target auth/RBAC smoke: PASS
Open P0 bugs: 0
Open P1 bugs: 0
Production GO: not approved by this evidence
```
