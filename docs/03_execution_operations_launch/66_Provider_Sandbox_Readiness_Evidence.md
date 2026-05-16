# Provider Sandbox Readiness Evidence

**Tasks:** T199, T208, T211, T212, T213, T214, T215
**Date:** 2026-05-16
**Decision:** real provider sandbox is not launch-ready yet. Cloudmini V3 intake is partially known, but authenticated read-only checks are blocked at the provider edge/gateway and approved credential storage, owners, quota, source mapping, cleanup, and a real pilot run are still missing.

## Scope

This record separates local fake-provider evidence from real sandbox-provider readiness. Local fake evidence is useful for CI and developer smoke tests, but it is not approval to provision real sandbox or production resources.

## Current Readiness

| Target | Current status | Evidence | Pilot decision |
|---|---|---|---|
| VPS local fake path | `ready` for local validation only | Fresh seed maps `vps-cx23-40gb-monthly` to `Local Fake Hetzner Ready`; provider contract tests cover fake Hetzner create, status, terminate, idempotency, and timeout mappings. | OK for local/CI smoke only. |
| Proxy/manual local path | documented non-ready state | Fresh seed maps `proxy-static-10gb-monthly` first to an unsupported VPS-style source and has a manual fallback path. This proves the readiness API surfaces the gap instead of silently treating proxy as ready. | Not ready for real proxy sandbox provisioning. |
| Real provider sandbox | `blocked` | Cloudmini V3 non-production base URL and API version are known. Unauthenticated `GET /api/v3/capabilities` returned HTTP `401`, confirming the endpoint is reachable and auth-gated. Authenticated read-only checks using bearer, `X-API-Key`, and `X-ACCESS-CODE` did not reach a successful V3 app envelope; the provider edge/gateway returned redacted HTTP `403` JSON. Approved credential storage path, account owner, quota/cost limit, source/group mapping, timeout policy, cleanup owner, and real pilot run are still missing. | Do not run pilot provisioning against real providers. |

## Proxy Cloudmini API V3 Candidate

T211 inspected the local `/opt/proxy-cloudmini` source code and added a Billing adapter for its API V3 contract using local `httptest` coverage only. T212 added disabled-by-default worker registry wiring behind explicit environment config. T213 recorded partial non-production intake for the Cloudmini V3 API. T214 attempted authenticated read-only provider checks and found an edge/gateway access blocker. This is still not real sandbox pilot evidence.

Known Cloudmini V3 intake as of 2026-05-15:

- Provider/API candidate: Cloudmini V3.
- Non-production API base URL: `https://cz.resvn.net/`.
- API version: V3.
- Auth boundary check: unauthenticated `GET /api/v3/capabilities` returned HTTP `401` in `2.475843s`; no response body was captured.
- Authenticated read-only check: transient credential input was used for `GET /api/v3/capabilities` and `GET /api/v3/inventory/groups?kind=<kind>` only. Bearer and `X-API-Key` returned HTTP `403` for capabilities, `ipv4_dc` inventory, and `residential` inventory. `X-ACCESS-CODE` returned HTTP `403` for both inventory checks and timed out once for capabilities after `20795ms`.
- Edge error shape: the redacted HTTP `403` body had Cloudflare/gateway-style keys including `cloudflare_error`, `error_code`, `ray_id`, `owner_action_required`, and `retryable`, not the expected V3 app envelope.
- Credential status: credential material must stay outside git, task notes, PR text, logs, and raw command output. T214 used it only as transient process input; an approved secret path or secret-manager reference is still required before shared authenticated testing.
- Pilot status: no authenticated provider call, create, delete, cleanup, or Billing end-to-end pilot has been run from this repository.

## Cloudmini Edge/Gateway Unblock Runbook

T215 documents the provider-owner handoff needed before another authenticated Cloudmini read-only check. The unblock is outside Billing runtime code because T214 reached the public base URL but received provider edge/gateway HTTP `403` responses before a successful V3 app envelope.

Provider owner must confirm these items before Billing reruns authenticated checks:

- Public hostname `https://cz.resvn.net/` routes `/api/v3/*` to the Cloudmini manager origin through the approved tunnel or gateway path.
- Edge/WAF/Access policy allows non-browser server-to-server clients for `/api/v3/capabilities` and `/api/v3/inventory/groups`.
- Edge policy allows the headers required by the code-read contract: `Authorization`, `X-API-Key`, `X-ACCESS-CODE`, `X-Request-ID`, and `Idempotency-Key`.
- If IP allowlisting is required, the provider owner records a redacted allowlist reference for the Billing runner egress IP or approved staging egress, not the raw credential.
- The provider API key is active, scoped to sandbox/non-production, and has read permission for capabilities and inventory plus later explicit `proxy_crud` permission only when create/delete pilot is approved.
- The credential shared through chat is rotated or explicitly accepted by Security/Provider owner as a temporary sandbox-only credential before reuse.
- Query-string credentials such as `?token=` or `?access_code=` are avoided for Billing evidence because they can leak through URL logs. Use headers unless a Security Owner signs a temporary exception.

Safe read-only rerun after unblock:

1. Store the rotated credential outside git in an approved secret path or local-only `.env` file.
2. Run only `GET /api/v3/capabilities` and `GET /api/v3/inventory/groups?kind=ipv4_dc|residential`.
3. Capture only status codes, envelope success, feature keys, inventory counts, sell-state counts, and redacted group references.
4. Do not capture raw provider response bodies, raw auth headers, provider-private IDs, proxy credentials, or URL query credentials.
5. Keep pilot readiness blocked unless both read-only checks return a successful V3 app envelope and owner/quota/mapping/cleanup evidence is recorded.

Do not run these until read-only evidence passes and the remaining provider readiness items are complete:

- `POST /api/v3/proxies`
- `DELETE /api/v3/proxies/:id`
- `POST /api/v3/proxies/:id/actions/:action`
- Billing checkout/provisioning worker pilot with `PROVIDER_DEFAULT_MODE=cloudmini_v3`

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
- provider edge/gateway allowlist or access policy that lets Billing reach `/api/v3` and evidence that the safe read-only rerun passed;
- source-to-group/SKU mapping for each Billing provider source;
- quota, rate, concurrency, timeout, and spend guardrails;
- redacted real error examples and one approved pilot create/delete run.

## Evidence Packet Status

No approved real provider sandbox pilot evidence is stored in this repository as of 2026-05-15. The packet below is the minimum evidence to collect before changing the real provider sandbox decision from `blocked`.

| Evidence area | Required proof | Current repo status |
|---|---|---|
| Provider intake | Provider name, sandbox account owner, support contact, docs/API version, and sandbox base URL. | Partial: Cloudmini V3, API V3, and `https://cz.resvn.net/` are recorded. Sandbox account owner and support contact are missing. |
| Credential safety | Approved secret store or local-only `.env` path, least-privilege scopes, rotation/revocation owner, and confirmation that no secret is committed. | Partial: no credential is committed in repo evidence. T214 used transient process input only. Approved secret path, scope, and rotation/revocation owner are missing. |
| Quota and cost guardrail | Sandbox quota, rate/concurrency limits, maximum spend or credit exposure, and stop condition. | Missing. |
| Capability mapping | Product type, Billing plan code, provider SKU, location, inventory mode, auto/manual provisioning support, cancellation support, and credential retrieval behavior. | Missing: authenticated read-only inventory did not succeed because provider edge/gateway returned HTTP `403`. |
| Retry/idempotency | Duplicate create behavior, timeout-after-send behavior, request/status lookup support, and mapping to retry safety or manual review. | Missing. |
| Error examples | Redacted auth, permission, rate limit, validation, out-of-capacity, timeout, duplicate, 5xx, not-found, and cancel-rejected examples. | Partial: redacted provider edge/gateway HTTP `403` shape captured. V3 app-level auth, permission, rate limit, validation, capacity, timeout, duplicate, 5xx, not-found, and cancel-rejected examples are still missing. |
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
