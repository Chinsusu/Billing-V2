# Provider Sandbox Readiness Evidence

**Tasks:** T199, T208, T211, T212, T213
**Date:** 2026-05-15
**Decision:** real provider sandbox is not launch-ready yet. Cloudmini V3 intake is partially known, but approved credential storage, owners, quota, source mapping, cleanup, and a real pilot run are still missing.

## Scope

This record separates local fake-provider evidence from real sandbox-provider readiness. Local fake evidence is useful for CI and developer smoke tests, but it is not approval to provision real sandbox or production resources.

## Current Readiness

| Target | Current status | Evidence | Pilot decision |
|---|---|---|---|
| VPS local fake path | `ready` for local validation only | Fresh seed maps `vps-cx23-40gb-monthly` to `Local Fake Hetzner Ready`; provider contract tests cover fake Hetzner create, status, terminate, idempotency, and timeout mappings. | OK for local/CI smoke only. |
| Proxy/manual local path | documented non-ready state | Fresh seed maps `proxy-static-10gb-monthly` first to an unsupported VPS-style source and has a manual fallback path. This proves the readiness API surfaces the gap instead of silently treating proxy as ready. | Not ready for real proxy sandbox provisioning. |
| Real provider sandbox | `blocked` | Cloudmini V3 non-production base URL and API version are known. Unauthenticated `GET /api/v3/capabilities` returned HTTP `401`, confirming the endpoint is reachable and auth-gated. Approved credential storage path, account owner, quota/cost limit, source/group mapping, timeout policy, cleanup owner, and real pilot run are still missing. | Do not run pilot provisioning against real providers. |

## Proxy Cloudmini API V3 Candidate

T211 inspected the local `/opt/proxy-cloudmini` source code and added a Billing adapter for its API V3 contract using local `httptest` coverage only. T212 added disabled-by-default worker registry wiring behind explicit environment config. T213 recorded partial non-production intake for the Cloudmini V3 API. This is still not real sandbox pilot evidence.

Known Cloudmini V3 intake as of 2026-05-15:

- Provider/API candidate: Cloudmini V3.
- Non-production API base URL: `https://cz.resvn.net/`.
- API version: V3.
- Auth boundary check: unauthenticated `GET /api/v3/capabilities` returned HTTP `401` in `2.475843s`; no response body was captured.
- Credential status: credential material must stay outside git, task notes, PR text, logs, and raw command output. An approved secret path or secret-manager reference is still required before authenticated testing.
- Pilot status: no authenticated provider call, create, delete, cleanup, or Billing end-to-end pilot has been run from this repository.

Code-read contract summary:

- Auth supports `Authorization: Bearer <token>` and API-key fallback headers in the provider service. Billing adapter uses bearer auth.
- Readiness/inventory endpoints: `GET /api/v3/capabilities`, `GET /api/v3/inventory/groups?kind=<kind>`.
- Supported proxy kinds from code: `ipv4_dc` and `residential`.
- Mutating V3 endpoints require `Idempotency-Key`.
- Create path: `POST /api/v3/proxies` returns `202 Accepted` with an async operation id and resource id.
- Status path: `GET /api/v3/operations/:id` polls `accepted/running/succeeded/failed/timed_out/cancelled`.
- Resource paths: `GET /api/v3/proxies/:id`, `DELETE /api/v3/proxies/:id`.
- Action path: `POST /api/v3/proxies/:id/actions/:action` supports `stop`, `start`, and residential-only `change-ip`.
- Credential-bearing proxy response fields are encrypted by Billing adapter tests before returning a provider `CredentialEnvelope`.
- Worker runtime stays on `PROVIDER_DEFAULT_MODE=fake` by default. `PROVIDER_DEFAULT_MODE=cloudmini_v3` requires Cloudmini base URL, API token, Billing source id, kind, group id, protocol, and `ENCRYPTION_KEY` before startup.

Still missing for real sandbox readiness:

- scoped credential storage path outside git;
- sandbox account owner and cleanup owner;
- source-to-group/SKU mapping for each Billing provider source;
- quota, rate, concurrency, timeout, and spend guardrails;
- redacted real error examples and one approved pilot create/delete run.

## Evidence Packet Status

No approved real provider sandbox pilot evidence is stored in this repository as of 2026-05-15. The packet below is the minimum evidence to collect before changing the real provider sandbox decision from `blocked`.

| Evidence area | Required proof | Current repo status |
|---|---|---|
| Provider intake | Provider name, sandbox account owner, support contact, docs/API version, and sandbox base URL. | Partial: Cloudmini V3, API V3, and `https://cz.resvn.net/` are recorded. Sandbox account owner and support contact are missing. |
| Credential safety | Approved secret store or local-only `.env` path, least-privilege scopes, rotation/revocation owner, and confirmation that no secret is committed. | Partial: no credential is committed in repo evidence. Approved secret path, scope, and rotation/revocation owner are missing. |
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
