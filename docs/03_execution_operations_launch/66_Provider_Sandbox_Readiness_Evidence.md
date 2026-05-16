# Provider Sandbox Readiness Evidence

**Tasks:** T199, T208, T211, T212, T213, T214, T215, T216, T217, T218, T219, T220, T221, T226
**Date:** 2026-05-16
**Decision:** real provider sandbox is not launch-ready yet. Cloudmini V3 intake and authenticated read-only reachability are proven through the public hostname, T218 defines a controlled pilot approval packet, T219/T221 provide guarded mapping/evidence tooling, T220 applied the pilot mapping on an approved non-production Billing dev DB, and T226 adds a non-mutating pilot preflight guard. Approved shared credential storage, named owners, cleanup sign-off, idempotency evidence, and a real create/delete pilot run are still missing.

## Scope

This record separates local fake-provider evidence from real sandbox-provider readiness. Local fake evidence is useful for CI and developer smoke tests, but it is not approval to provision real sandbox or production resources.

## Current Readiness

| Target | Current status | Evidence | Pilot decision |
|---|---|---|---|
| VPS local fake path | `ready` for local validation only | Fresh seed maps `vps-cx23-40gb-monthly` to `Local Fake Hetzner Ready`; provider contract tests cover fake Hetzner create, status, terminate, idempotency, and timeout mappings. | OK for local/CI smoke only. |
| Proxy/manual local path | documented non-ready state | Fresh seed maps `proxy-static-10gb-monthly` first to an unsupported VPS-style source and has a manual fallback path. This proves the readiness API surfaces the gap instead of silently treating proxy as ready. | Not ready for real proxy sandbox provisioning. |
| Real provider sandbox | `blocked` | Cloudmini V3 non-production base URL and API version are known. A 2026-05-16 read-only rerun through `https://cz.resvn.net/` reached the V3 app: unauthenticated capabilities returned HTTP `401`, and authenticated capabilities plus inventory returned HTTP `200` V3 success envelopes using bearer, `X-API-Key`, and `X-ACCESS-CODE`. T218 selected a redacted `ipv4_dc` pilot mapping candidate from sellable read-only inventory and defines quota/cleanup approval requirements. T220 applied the guarded Cloudmini pilot mapping on the approved Billing dev runtime env and T221 evidence passed with plan display `10002`, plan-source display `10024`, source display `10012`, source type `cloudmini_v3`, readiness `ready`, priority `1`, and first-pilot guardrails of one create, one active resource, and one worker concurrency. Approved shared credential storage path, account owner, timeout policy, cleanup owner, and real create/delete pilot run are still missing. | Do not run pilot provisioning against real providers until the remaining approval fields, owner sign-offs, timeout/idempotency evidence, cleanup procedure, and same-session cleanup plan are complete. |

## Proxy Cloudmini API V3 Candidate

T211 inspected the local `/opt/proxy-cloudmini` source code and added a Billing adapter for its API V3 contract using local `httptest` coverage only. T212 added disabled-by-default worker registry wiring behind explicit environment config. T213 recorded partial non-production intake for the Cloudmini V3 API. T214 attempted authenticated read-only provider checks and found an edge/gateway access blocker. T216 reran read-only checks with the local dev credential source and reached successful V3 app envelopes through the public hostname. This is still not real sandbox pilot evidence.

Known Cloudmini V3 intake as of 2026-05-16:

- Provider/API candidate: Cloudmini V3.
- Non-production API base URL: `https://cz.resvn.net/`.
- API version: V3.
- Auth boundary check: unauthenticated `GET /api/v3/capabilities` returned HTTP `401` with app JSON in `711ms`, confirming the public hostname routes to the Cloudmini manager and keeps auth enforced.
- Authenticated read-only rerun: T216 used the local dev credential source at `/opt/cred` without printing the raw key. The run used the Billing Go-client-style user-agent plus `X-Request-ID`; the checker was limited to read-only endpoints.
- Header forwarding result: bearer `Authorization`, `X-API-Key`, and `X-ACCESS-CODE` each returned HTTP `200` V3 success envelopes for `GET /api/v3/capabilities`, `GET /api/v3/inventory/groups?kind=ipv4_dc`, and `GET /api/v3/inventory/groups?kind=residential`.
- Capability summary: feature keys returned were `inventory_webhooks`, `prefer_wait`, `reservations`, and `tombstones`.
- Inventory summary: `ipv4_dc` returned `2` groups, with `1` sellable and `1` exhausted, totaling `200` allocatable units. `residential` returned `4` groups, all exhausted, totaling `0` allocatable units.
- Controlled pilot mapping candidate: T218 selected the sellable `ipv4_dc` group as `redacted:c6a7189f0a` with `200` allocatable units and `socks5` protocol. The raw group id is stored only in `/opt/cred-cloudmini-dev.env`.
- Edge note: provider-side evidence reported Cloudflare still blocks the generic `Python-urllib/3.12` user-agent with HTTP `403` code `1010`. The Billing Go-client-style path passed, so launch evidence should use the Billing adapter or the provider checker user-agent override, not a generic scripting user-agent.
- Credential status: credential material must stay outside git, task notes, PR text, logs, and raw command output. The local dev provider credential has been split into `/opt/cred-cloudmini-dev.env` with mode `0600`; an approved shared secret path or secret-manager reference is still required before shared authenticated testing or pilot provisioning.
- Multi-endpoint status: T217 adds runtime support for multiple Cloudmini V3 endpoint/API-key mappings through `CLOUDMINI_V3_MAPPINGS_JSON`, keyed by provider source and optionally provider account. T220 covers only the single dev pilot source mapping; no multi-account target mapping or real provider pilot was run.
- Catalog mapping status: T219 adds `migrations/0025_add_cloudmini_provider_type.sql` and `scripts/cloudmini_pilot_mapping.sh` so an approved non-production DB can create the pilot `cloudmini_v3` provider source and plan-source mapping. T220 ran it on the approved Billing dev runtime env with `APP_ENV=dev`; migration plan showed `0` pending migrations and migration apply reported `0` applied migrations.
- Mapping evidence collector status: T221 adds `scripts/cloudmini_mapping_evidence.sh` so an operator can verify the applied mapping on an approved non-production Billing DB without sharing DSNs, tokens, raw group IDs, or raw provider payloads in repo evidence. T220 ran it read-only and recorded `result=PASS`, readiness `ready`, source type `cloudmini_v3`, and redacted guardrails only.
- Mutating pilot preflight status: T226 adds `scripts/cloudmini_mutating_pilot_preflight.sh` to fail closed unless non-production env, mapping evidence, owner fields, cleanup fields, private credential path, and exact one-resource guardrails are present. The script does not call provider mutating routes.
- Target DB access status: T220 used the approved test server Billing runtime env at `/opt/Billing/.env.dev`; the run confirmed `APP_ENV=dev` and `DB_DSN` presence without printing the DSN or provider secrets.
- Pilot status: no create, delete, cleanup, or Billing end-to-end pilot has been run from this repository.

## Cloudmini Edge/Gateway Unblock Runbook

T215 documents the provider-owner handoff needed for authenticated Cloudmini read-only checks. The unblock is outside Billing runtime code because T214 reached the public base URL but initially received provider edge/gateway HTTP `403` responses before a successful V3 app envelope. T216 shows the required read-only route/header path now works for the Billing Go-client-style path.

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

Although read-only evidence now passes, do not run these until the remaining provider readiness items are complete:

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
- provider edge/gateway approval record for the read-only route/header policy;
- owner-approved source-to-group/SKU mapping beyond the dev pilot evidence;
- active Cloudmini V3 provider source and plan source readiness evidence for any additional target provider source;
- multi-endpoint/account config if more than one Cloudmini V3 URL or API key is needed;
- a filled T226 preflight run with real owner-approved values and owner-approved timeout/spend guardrail sign-off beyond the dev defaults;
- redacted real error examples and one approved pilot create/delete run.

## Evidence Packet Status

No approved real provider sandbox mutating pilot evidence is stored in this repository as of 2026-05-16. The packet below is the minimum evidence to collect before changing the real provider sandbox decision from `blocked`.

| Evidence area | Required proof | Current repo status |
|---|---|---|
| Provider intake | Provider name, sandbox account owner, support contact, docs/API version, and sandbox base URL. | Partial: Cloudmini V3, API V3, and `https://cz.resvn.net/` are recorded. Sandbox account owner and support contact are missing. |
| Credential safety | Approved secret store or local-only `.env` path, least-privilege scopes, rotation/revocation owner, and confirmation that no secret is committed. | Partial: no credential is committed in repo evidence. T216 used `/opt/cred` as a local dev-only source without printing the raw key. Approved shared secret path, scope approval, and rotation/revocation owner are missing. |
| Quota and cost guardrail | Sandbox quota, rate/concurrency limits, maximum spend or credit exposure, and stop condition. | Partial: T220 recorded first-pilot dev guardrails of one create, one active resource, one worker concurrency, no parallel mutating calls, and single-dev-resource exposure. T226 adds a preflight guard that requires these fields, but owner sign-off, sandbox quota, and stop-condition approval are still missing. |
| Capability mapping | Product type, Billing plan code, provider SKU, location, inventory mode, auto/manual provisioning support, cancellation support, and credential retrieval behavior. | Partial: authenticated read-only inventory succeeds and shows `ipv4_dc` has sellable capacity while `residential` is exhausted. T220 applied dev pilot mapping for `proxy-static-10gb-monthly` to `cloudmini_v3` with readiness `ready`; owner-approved SKU/source mapping and cleanup/cancellation behavior are still missing. |
| Retry/idempotency | Duplicate create behavior, timeout-after-send behavior, request/status lookup support, and mapping to retry safety or manual review. | Missing. |
| Error examples | Redacted auth, permission, rate limit, validation, out-of-capacity, timeout, duplicate, 5xx, not-found, and cancel-rejected examples. | Partial: redacted provider edge/gateway HTTP `403` shape from the earlier blocked run and app-level unauthenticated HTTP `401` were captured. V3 app-level permission, rate limit, validation, capacity, timeout, duplicate, 5xx, not-found, and cancel-rejected examples are still missing. |
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
