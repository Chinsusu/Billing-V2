# Provider Sandbox Readiness Evidence

**Tasks:** T199, T208
**Date:** 2026-05-14
**Decision:** real provider sandbox is not launch-ready yet.

## Scope

This record separates local fake-provider evidence from real sandbox-provider readiness. Local fake evidence is useful for CI and developer smoke tests, but it is not approval to provision real sandbox or production resources.

## Current Readiness

| Target | Current status | Evidence | Pilot decision |
|---|---|---|---|
| VPS local fake path | `ready` for local validation only | Fresh seed maps `vps-cx23-40gb-monthly` to `Local Fake Hetzner Ready`; provider contract tests cover fake Hetzner create, status, terminate, idempotency, and timeout mappings. | OK for local/CI smoke only. |
| Proxy/manual local path | documented non-ready state | Fresh seed maps `proxy-static-10gb-monthly` first to an unsupported VPS-style source and has a manual fallback path. This proves the readiness API surfaces the gap instead of silently treating proxy as ready. | Not ready for real proxy sandbox provisioning. |
| Real provider sandbox | `blocked` | Missing approved sandbox provider account, API base URL, credential storage path, quota/cost limit, provider SKU mapping, timeout policy, and cleanup owner. | Do not run pilot provisioning against real providers. |

## Evidence Packet Status

No approved real provider sandbox evidence is stored in this repository as of 2026-05-14. The packet below is the minimum evidence to collect before changing the real provider sandbox decision from `blocked`.

| Evidence area | Required proof | Current repo status |
|---|---|---|
| Provider intake | Provider name, sandbox account owner, support contact, docs/API version, and sandbox base URL. | Missing. |
| Credential safety | Approved secret store or local-only `.env` path, least-privilege scopes, rotation/revocation owner, and confirmation that no secret is committed. | Missing. |
| Quota and cost guardrail | Sandbox quota, rate/concurrency limits, maximum spend or credit exposure, and stop condition. | Missing. |
| Capability mapping | Product type, Billing plan code, provider SKU, location, inventory mode, auto/manual provisioning support, cancellation support, and credential retrieval behavior. | Missing. |
| Retry/idempotency | Duplicate create behavior, timeout-after-send behavior, request/status lookup support, and mapping to retry safety or manual review. | Missing. |
| Error examples | Redacted auth, permission, rate limit, validation, out-of-capacity, timeout, duplicate, 5xx, not-found, and cancel-rejected examples. | Missing. |
| Cleanup and rollback | How to list test resources, cancel/delete them, disable the provider source, and assign manual cleanup owner. | Missing. |
| Pilot run | Redacted evidence for one approved sandbox order through checkout, reservation, provider request, service activation, credential storage/reveal audit, and cleanup. | Missing. |

## Evidence Packet Template

Use this template in a task, runbook appendix, or external launch packet. Store only redacted values in git.

```text
Provider:
Sandbox account owner:
Provider support contact:
Sandbox API base URL:
Provider docs/API version:
Approved credential storage path:
Credential scope:
Credential rotation/revocation owner:
Quota/rate/concurrency limits:
Maximum sandbox spend or quota exposure:
Stop condition:
Billing plan code:
Provider SKU:
Sandbox location:
Inventory mode:
Auto provisioning supported:
Manual fallback supported:
Credential retrieval behavior:
Cancellation/cleanup supported:
Provider idempotency level:
Duplicate create behavior:
Timeout-after-send behavior:
Safe status lookup path:
Retry safety map:
Redacted error examples captured:
Pilot resource cleanup owner:
Pilot run date:
Pilot run reviewer:
```

## Required Pilot Run Evidence

Before approving real sandbox provisioning, capture redacted evidence for all of these:

- Source readiness is `ready` for the exact plan/source being tested.
- Checkout debits wallet and creates a single reservation and provisioning job.
- Provider create uses the job idempotency key or a documented equivalent.
- Success creates one active service and one provider resource mapping.
- Duplicate retry does not create a second provider resource.
- Timeout-after-send moves to manual review or safe status lookup, not blind retry.
- Provider credentials are stored encrypted and not printed in logs, provider request records, audit, task notes, or PR text.
- Credential reveal remains a separate audited action.
- Cleanup/cancel removes or disables the sandbox resource and records the cleanup owner.
- Any failed case maps to an internal provider error code and retry safety.

## Redaction Rules

Do not commit or paste:

- API keys, bearer tokens, signatures, cookies, private keys, passwords, root credentials, proxy credentials, or raw auth headers.
- Production provider account IDs, production DSNs, production customer data, or real customer order data.
- Raw provider request or response bodies if they contain secrets, customer data, or provider-private identifiers.

Use stable display IDs or redacted placeholders when a human-readable reference is needed.

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
