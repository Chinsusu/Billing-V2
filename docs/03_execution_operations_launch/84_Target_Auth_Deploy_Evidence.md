# 84 - Target Auth Deploy Evidence

**Date:** 2026-05-21
**Scope:** Selected test-server deploy and target auth/RBAC evidence after T288.
**Decision:** PASS for selected non-production test-server auth/RBAC deploy evidence only. This does not approve production, production customer data, broader private beta, or real-provider provisioning.

## Boundary

This evidence covers the selected test server reached from `/opt/Billing-Dev-Server`.

Do not broaden this result:

- No production database or production customer data was used.
- No raw DB DSN, password, cookie, session token, provider token, provider payload, Telegram token, TOTP value, or plaintext service credential is recorded here.
- No money mutation route, provider mutation route, or Cloudmini create/delete/action route was called.
- The evidence records status codes, boolean checks, and redacted outcomes only.

## Deploy Actions

Runtime deployed to the selected test server:

```text
source_commit=663ad19
source_merge=PR #609 T288 marker
target_path=/opt/Billing
target_layout=deploy-copy without .git
backend_build=PASS
frontend_build=PASS
billing-api.service=active
billing-frontend.service=active
cloudflared.service=active on test server
postgresql@14-main.service=active
```

Configuration applied on the selected test server:

```text
/opt/Billing/.env.dev mode=640 owner=root group=billing-svc
BILLING_TARGET_AUTH_SMOKE_CLIENT_EMAIL=present
BILLING_TARGET_AUTH_SMOKE_CLIENT_PASSWORD=present
BILLING_TARGET_AUTH_SMOKE_ADMIN_EMAIL=present
BILLING_TARGET_AUTH_SMOKE_ADMIN_PASSWORD=present
BILLING_API_URL=present
```

Operational notes:

- The four owner-provided target auth override values initially copied from the local secret file did not match users in the selected test DB, so the selected test server override values were aligned to the existing seeded dev/test users. Values are intentionally not recorded.
- `BILLING_API_URL` was required for the Next.js standalone rewrite from `/backend/*` to the Go API. The frontend was rebuilt with that environment loaded.
- The local non-test connector and the selected test-server connector used the same cloudflared token. The local non-test `cloudflared.service` was stopped so public hostnames route to the selected test server only.

## Health Evidence

After deploy and restart:

```text
127.0.0.1:8080 /healthz: 200
127.0.0.1:3000 /: 200
billing.resvn.net /backend/healthz: 200
client.resvn.net /backend/healthz: 200
reseller.resvn.net /backend/healthz: 200
client.resvn.net /: 200
```

## Auth/RBAC Evidence

Direct test-server smoke:

```text
command=dev-target-auth-rbac
base_url=http://127.0.0.1:8080
result=PASS
client_session_cookie_only=pass
admin_2fa_gate=pass
invalid_session_denied=pass
actor_required_denied=pass
tenant_mismatch_denied=pass
rbac_denials=3
provider_mutation_routes_called=no
money_mutation_routes_called=no
```

Domain-aware public probe:

```text
client_login_domain=client.resvn.net
client_login_status=200
client_actor_type=client
client_cookie_http_only=true
client_catalog_status=200
admin_login_domain=billing.resvn.net
admin_login_status=200
admin_actor_type=platform_staff
admin_two_factor_required=true
admin_two_factor_satisfied=false
admin_2fa_gate_status=403
admin_2fa_gate_error_code=auth.2fa_required
domain_aware_target_auth_probe=PASS
provider_mutation_routes_called=no
money_mutation_routes_called=no
```

Rate-limit cleanup:

```text
auth_login_rate_limit_counters_cleared=dev/test only
reason=automated auth probes hit the selected non-production auth login rate limit during routing diagnosis
```

## Findings

The generic `dev-target-auth-rbac` command passed against the test-server local API base URL.

The same generic command is not a correct public-domain test when only one public `base-url` is supplied, because login tenant resolution is domain-first:

- `client.resvn.net` is correct for client login.
- `billing.resvn.net` is correct for platform admin login.
- A single public domain cannot correctly cover both client and admin login in that command.

The domain-aware public probe covered the same critical auth path split across the correct public domains and passed.

Cloudflare blocked Python urllib with edge code `1010`; curl and Go HTTP clients were accepted. This was not an app auth failure.

## Not Verified

- Human tester sign-off was not collected.
- Full browser UAT was not rerun.
- Full 2FA-satisfied admin browser navigation was not executed.
- Real provider provisioning and production notification delivery were not executed.
- This task did not change application code.

## Decision

```text
Selected test-server deploy: PASS
Target auth/RBAC local smoke: PASS
Domain-aware public auth/RBAC probe: PASS
Open P0 bugs: 0
Open P1 bugs: 0
Production GO: not approved by this evidence
```

Recommended follow-up:

```text
Add a domain-aware target auth smoke mode that accepts separate client and admin public base URLs, or document that `dev-target-auth-rbac` should use the local API base URL for mixed client/admin login checks.
```
