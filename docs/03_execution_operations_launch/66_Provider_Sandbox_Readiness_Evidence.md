# Provider Sandbox Readiness Evidence

**Task:** T199
**Date:** 2026-05-13
**Decision:** real provider sandbox is not launch-ready yet.

## Scope

This record separates local fake-provider evidence from real sandbox-provider readiness. Local fake evidence is useful for CI and developer smoke tests, but it is not approval to provision real sandbox or production resources.

## Current Readiness

| Target | Current status | Evidence | Pilot decision |
|---|---|---|---|
| VPS local fake path | `ready` for local validation only | Fresh seed maps `vps-cx23-40gb-monthly` to `Local Fake Hetzner Ready`; provider contract tests cover fake Hetzner create, status, terminate, idempotency, and timeout mappings. | OK for local/CI smoke only. |
| Proxy/manual local path | documented non-ready state | Fresh seed maps `proxy-static-10gb-monthly` first to an unsupported VPS-style source and has a manual fallback path. This proves the readiness API surfaces the gap instead of silently treating proxy as ready. | Not ready for real proxy sandbox provisioning. |
| Real provider sandbox | `blocked` | Missing approved sandbox provider account, API base URL, credential storage path, quota/cost limit, provider SKU mapping, timeout policy, and cleanup owner. | Do not run pilot provisioning against real providers. |

## Verification Evidence

Provider contract expectations are covered locally by:

```bash
go test ./internal/modules/provider -run SandboxContract
```

Provisioning worker timeout safety is covered by:

```bash
go test ./internal/modules/order -run ProviderProvisioningHandler
```

The required local billing smoke before pilot remains:

```bash
go run ./cmd/smoke -dsn "$DB_DSN" -base-url "$API_BASE_URL" dev-billing
```

That smoke requires a running API and a non-production database. Do not treat unit tests as a substitute for that smoke before pilot.

## No-Go Until Fixed

Before changing this decision to ready, record the required provider intake from `docs/05_development_standards/60_Provider_Sandbox_Contract_Checklist.md`:

- sandbox provider name, account owner, and support contact;
- sandbox API base URL and docs version;
- sandbox-only credential storage outside git;
- supported VPS/proxy product types, locations, quota, and rate limits;
- idempotency behavior for duplicate create and timeout-after-create;
- redacted examples for auth, rate limit, validation, timeout, and provider 5xx errors;
- cleanup and rollback plan for resources created by sandbox tests.

If any item is missing, keep real provider sandbox readiness blocked.
