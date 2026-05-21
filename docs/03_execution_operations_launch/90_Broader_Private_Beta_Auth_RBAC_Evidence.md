# 90 - Broader Private Beta Auth RBAC Evidence

**Date:** 2026-05-21
**Scope:** Domain-aware target auth/RBAC smoke evidence for the current broader private beta launch-candidate domains and selected test server.
**Decision:** PASS for current client/admin target auth/RBAC smoke only. Broader private beta remains `NO-GO`.

## Boundary

This evidence covers the current launch-candidate public domains and selected test server only.

Do not broaden this result:

- No production database, production provider account, production notification channel, production customer data, or customer list was used.
- No raw DB DSN, password, cookie, session token, provider token, Telegram token, TOTP value, private key, provider payload, notification payload, proxy credential, or customer data is recorded here.
- No money mutation route, provider mutation route, Cloudmini create/delete/action route, credential reveal route, full browser UAT, or notification delivery route was called by this task.
- This evidence does not approve broader private beta, production launch, provider expansion, production customer data, or Telegram as a production primary path.

## Smoke Command

The remote smoke binary was executed on the selected test server with separate public API base URLs for client and platform admin login. Credential values were loaded from the target secret environment and were not printed.

```text
command=dev-target-auth-rbac
target_path=/opt/Billing
base_url=http://127.0.0.1:8080
client_base_url=https://client.resvn.net/backend
admin_base_url=https://billing.resvn.net/backend
credential_values_loaded_from_target_secret_env=yes
secret_values_printed=no
```

## Result

```text
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

## What This Proves

- Client login works through the client public backend base URL and uses an HTTP-only session cookie.
- Platform admin login works through the billing public backend base URL and correctly requires unsatisfied 2FA before protected admin access.
- Invalid session, missing actor context, tenant mismatch, and RBAC-negative checks deny access.
- The smoke used separate domain-aware base URLs instead of a single local-only URL.
- The smoke did not call money mutation routes or provider mutation routes.

## Not Verified

- Owner approval for broader private beta v1 was not collected.
- Customer list or data classification was not approved.
- Reseller-domain login and full reseller UAT were not rerun in this task.
- Full browser UAT, full E2E, renewal, credential reveal, audit trail review, and finance reconciliation were not rerun.
- Provider quota/SKU/group mapping, timeout/idempotency, cleanup, and pilot run evidence were not expanded.
- Notification primary/fallback SLA, escalation, delivery, and failure drill were not rerun for broader private beta.
- Production launch, production customer data, and broader provider scope remain unapproved.

## Decision

```text
Domain-aware client auth smoke: PASS
Domain-aware admin auth smoke: PASS
Admin 2FA gate: PASS
Tenant/RBAC negative checks: PASS
Raw secret values printed: no
Money mutation routes called: no
Provider mutation routes called: no
Broader private beta GO: not approved
Production GO: not approved
```
