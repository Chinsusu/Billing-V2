# 71 - Cloudmini Controlled Pilot Runbook

**Date:** 2026-05-17
**Scope:** Controlled pre-approval packet for the first Cloudmini V3 mutating pilot.
**Decision:** One controlled dev create/delete pilot has passed. T229 fixes the repo-side non-usable-status and lifecycle-worker cleanup gaps, T230 proves that hardening deploys/builds on the approved test server, and T231 proves non-mutating Cloudmini registry activation. Broader pilot remains blocked until live duplicate/timeout evidence, shared secret storage, owner-approved mutating/lifecycle activation, and owner sign-offs are complete.

## Current Safe State

Read-only evidence is complete for the Billing Go-client-style path:

- `GET /api/v3/capabilities` without auth returns app-level HTTP `401`.
- Authenticated `GET /api/v3/capabilities` returns HTTP `200` V3 success envelopes.
- Authenticated `GET /api/v3/inventory/groups?kind=ipv4_dc` returns HTTP `200` V3 success envelopes.
- Authenticated `GET /api/v3/inventory/groups?kind=residential` returns HTTP `200` V3 success envelopes.
- T228 ran one approved dev mutating pilot through Billing checkout/payment/provisioning worker and cleaned up the provider resource in the same session.
- T229 makes Cloudmini provisioning fail closed when the provider resource status is not usable and makes lifecycle-worker termination call provider delete before marking a service terminated.
- T230 deployed the T229 hardening to the approved Billing test server and verified build/service health without calling Cloudmini mutating routes.
- T231 verified the worker can build a real Cloudmini provider registry from the protected dev credential path without opening DB state, claiming jobs, or calling provider APIs.

## Pilot Mapping Candidate

Use this mapping only after owner approval:

```text
Billing plan candidate: proxy-static-10gb-monthly
Provider kind: ipv4_dc
Provider group reference: redacted:c6a7189f0a
Provider sell state: sellable
Observed allocatable units: 200
Protocol: socks5
Credential/config source: /opt/cred-cloudmini-dev.env
```

Do not use the existing seeded `Local Fake Hetzner Ready` source as the Cloudmini pilot source. The pilot needs an explicit Cloudmini V3 provider source or equivalent dev/staging source record whose `source_type` is `cloudmini_v3`.

Do not pilot `residential` yet. The read-only evidence shows `residential` inventory is exhausted.

## Catalog Mapping Procedure

T219 adds the missing `cloudmini_v3` catalog provider type migration and a guarded mapping script:

```text
migrations/0025_add_cloudmini_provider_type.sql
scripts/cloudmini_pilot_mapping.sh
```

Apply the mapping only on an approved non-production database, after the approval fields below are filled:

```bash
go run ./cmd/migrate -dsn "$DB_DSN" up
APP_ENV=dev \
BILLING_CLOUDMINI_PILOT_APPROVED=yes \
bash scripts/cloudmini_pilot_mapping.sh
```

The script is intentionally narrow:

- It refuses `prod` and `production`.
- It requires `DB_DSN`, `APP_ENV`, and `BILLING_CLOUDMINI_PILOT_APPROVED=yes`.
- It creates or updates a `cloudmini_v3` provider source for the pilot plan.
- It links `proxy-static-10gb-monthly` to that source at priority `1`.
- It moves only the seeded fake Hetzner link for that proxy plan to priority `5` so it cannot win the source selection tie.
- It records only redacted group reference and guardrail metadata in `capacity_policy`.
- It must not print or store raw API tokens, raw auth headers, raw provider group IDs, raw provider payloads, proxy credentials, or DSNs.

After applying the mapping, verify the admin provider-readiness API shows the pilot proxy plan as `ready`, with `source_type=cloudmini_v3`, using display IDs only in evidence.

For operator-run DB evidence without exposing the DSN, token, raw group id, or provider payloads, run the read-only evidence collector on the approved non-production Billing DB:

```bash
APP_ENV=dev \
BILLING_CLOUDMINI_EVIDENCE_APPROVED=yes \
DB_DSN="$DB_DSN" \
CLOUDMINI_V3_PLAN_CODE=proxy-static-10gb-monthly \
bash scripts/cloudmini_mapping_evidence.sh
```

Optional public-ID guards can be supplied after the mapping script prints them:

```bash
CLOUDMINI_V3_SOURCE_DISPLAY_ID=<source_display_id>
CLOUDMINI_V3_PLAN_SOURCE_DISPLAY_ID=<plan_source_display_id>
```

The evidence collector runs in a read-only transaction and passes only when the selected plan source is `cloudmini_v3`, readiness is `ready`, priority is `1`, and first-pilot guardrails remain `1` create, `1` active resource, and `1` worker concurrency.

T220 target-environment mapping evidence was applied on the approved Billing dev runtime env at `/opt/Billing/.env.dev`:

- `APP_ENV=dev` was confirmed before DB access.
- `DB_DSN` presence was confirmed without printing the DSN.
- Migration validation found `25` migrations, plan showed `0` pending, and `up` applied `0`.
- The guarded mapping script returned plan-source display `10024`, source display `10012`, source type `cloudmini_v3`, active statuses, provider-live inventory mode, and priority `1`.
- The read-only evidence collector returned `result=PASS`, plan display `10002`, product type `proxy`, readiness `ready`, redacted group ref `redacted:c6a7189f0a`, protocol `socks5`, one-create/one-active-resource/one-worker guardrails, and `failed_checks=none`.
- No checkout, worker provisioning, provider create/delete, provider action, raw provider group id, DSN, token, or proxy credential was printed or stored in repo evidence.

This unblocks the dev mapping evidence gate only. Keep broader mutating use blocked until owner sign-offs, timeout/idempotency evidence, target-environment lifecycle-worker cleanup evidence, and residual-risk follow-ups below are complete.

## Required Approval Fields

Fill these before any mutating call:

```text
Pilot ID:
Environment:
Billing source display ID:
Cloudmini source/account owner:
Engineering owner:
Ops owner:
Security owner:
Cleanup owner:
Finance/quota owner:
Approved credential path: /opt/cred-cloudmini-dev.env for dev only, or redacted shared secret reference
Maximum create calls:
Maximum active test resources:
Maximum spend/quota exposure:
Worker concurrency:
Provider rate limit:
Stop condition:
Cleanup deadline:
Reviewer sign-off:
```

Minimum guardrails for the first pilot:

- `Maximum create calls`: `1`
- `Maximum active test resources`: `1`
- `Worker concurrency`: `1`
- `Provider rate limit`: no parallel mutating calls
- `Cleanup deadline`: same session as create

## Mutating Pilot Preflight Guard

T226 adds a non-mutating preflight guard:

```text
scripts/cloudmini_mutating_pilot_preflight.sh
```

Run it only after the approval fields above are filled with real owner-approved values:

```bash
APP_ENV=dev \
BILLING_CLOUDMINI_MUTATING_PREFLIGHT_APPROVED=yes \
DB_DSN="$DB_DSN" \
CLOUDMINI_PILOT_ID=<approved-pilot-id> \
CLOUDMINI_PILOT_ENVIRONMENT=dev \
CLOUDMINI_SOURCE_ACCOUNT_OWNER=<owner> \
CLOUDMINI_ENGINEERING_OWNER=<owner> \
CLOUDMINI_OPS_OWNER=<owner> \
CLOUDMINI_SECURITY_OWNER=<owner> \
CLOUDMINI_CLEANUP_OWNER=<owner> \
CLOUDMINI_FINANCE_QUOTA_OWNER=<owner> \
CLOUDMINI_REVIEWER_SIGNOFF=<reviewer> \
CLOUDMINI_PILOT_CLEANUP_DEADLINE=<same-session-deadline> \
CLOUDMINI_PILOT_STOP_CONDITION=<approved-stop-condition> \
CLOUDMINI_PILOT_READONLY_EVIDENCE_REF=<redacted-readonly-evidence-ref> \
CLOUDMINI_PILOT_CLEANUP_PROCEDURE_REF=<redacted-cleanup-procedure-ref> \
CLOUDMINI_PILOT_CREDENTIAL_PATH=/opt/cred-cloudmini-dev.env \
CLOUDMINI_PILOT_MAX_CREATE_CALLS=1 \
CLOUDMINI_PILOT_MAX_ACTIVE_RESOURCES=1 \
CLOUDMINI_PILOT_WORKER_CONCURRENCY=1 \
CLOUDMINI_PILOT_PROVIDER_RATE_LIMIT=no-parallel-mutating-calls \
CLOUDMINI_PILOT_MAX_SPEND_EXPOSURE=single-dev-resource \
bash scripts/cloudmini_mutating_pilot_preflight.sh
```

The preflight guard:

- refuses `prod` and `production`;
- requires `DB_DSN`, explicit preflight approval, owner fields, cleanup fields, read-only evidence reference, cleanup procedure reference, and exact one-resource guardrails;
- requires the credential path to be outside the repo, readable, and not group/world accessible;
- reruns the read-only mapping evidence collector and requires `result=PASS`;
- prints only mapping display IDs, redacted guardrails, and presence flags;
- does not call provider create/delete/action routes, Billing checkout, Billing payment, or the provisioning worker.

Passing this preflight is not approval to run the mutating pilot by itself. It only proves the required fields are present and the mapped source is still ready.

## T248 Idempotency And Timeout Evidence Smoke

T248 adds a fail-closed smoke command for the remaining live Cloudmini retry/idempotency evidence:

```bash
go run ./cmd/smoke cloudmini-idempotency-evidence
```

Detailed run instructions live in `docs/03_execution_operations_launch/73_Cloudmini_Idempotency_Evidence_Runbook.md` to keep this pilot runbook focused. The command is not itself GO evidence until both scenarios are run against the approved target provider account, cleanup succeeds, raw cleanup references are retained outside git with restricted permissions, and the redacted outputs are recorded in doc 66/doc 70.

## T228 Controlled Dev Pilot Evidence

Pilot ID: `T228-dev-20260517T004039Z`

Environment and preflight:

- Target environment: `APP_ENV=dev` on the approved non-production Billing test server.
- Local checks before pilot: `go test ./...` passed; `go run ./cmd/taskguard` passed.
- T226 preflight: `preflight_result=PASS`.
- Provider `proxy_crud` read permission precheck: `GET /api/v3/proxies?external_ref=<redacted-empty-ref>` returned HTTP `200` with `success=true`.
- Before create: Cloudmini active Billing services `0`, queued provider jobs `0`, selected group sellable, allocatable `200`, active proxy count `0`, pending create count `0`, reserved count `0`.
- Worker loop was stopped before payment so the default fake worker could not claim the pilot job.

Billing path evidence:

- Plan display ID: `10002`.
- Source display ID: `10012`.
- Order display ID: `10001`.
- Invoice display ID: `10002`.
- Transaction display ID: `10001`.
- Ledger display ID: `10002`.
- Provisioning job display ID: `10001`.
- Job idempotency key: present, not printed.
- Provider external ref: `redacted:8794850e2b96`.
- One-off worker command used `PROVIDER_DEFAULT_MODE=cloudmini_v3` and returned `claimed=1 succeeded=1 retried=0 manual_review=0 terminal_failed=0 cancelled=0`.
- Provisioning result: display ID `10001`, status `provisioned`.
- Service display ID: `10001`, status `active`, billing status `paid`.
- External provider resource ref: `redacted:dc3d9457bf5b`.
- Credential storage: active credential count `1`, encrypted payload present, masked hint present.

Cleanup evidence:

- Provider `GET /api/v3/proxies/<resource>` after create returned HTTP `200`, kind `ipv4_dc`, provider resource status `creating`.
- Provider cleanup used the approved V3 provider delete path because Billing service terminate currently does not invoke provider deletion.
- Provider delete operation reached `succeeded`.
- Provider `GET /api/v3/proxies/<resource>` after cleanup returned HTTP `404`.
- Billing service cleanup used the reseller service terminate route with an access reason and returned service status `terminated`.
- After cleanup: Cloudmini active Billing services `0`; provider selected group allocatable `200`, active proxy count `0`, pending create count `0`, reserved count `0`.
- The regular `billing-worker` service was restarted after the pilot.

Residual risks before broader pilot:

- Direct HTTP service terminate remains a lifecycle transition API. For provider-backed cleanup, use the lifecycle worker path with provider registry configured. Direct provider delete remains a dev-pilot exception only when owner approved and redacted cleanup evidence is recorded.
- T229 changes Cloudmini provisioning to fail closed if the resource status is not usable. Operation success with provider status `creating` now moves to manual review instead of creating an active service.
- Duplicate-create and timeout-after-send behavior are still not proven against the live provider.
- Production/shared secret-store owner and named launch owners are still not recorded in repo evidence.
- The always-on target worker remains on `PROVIDER_DEFAULT_MODE=fake`; Cloudmini mutating or lifecycle cleanup activation requires a new owner-approved window.

## T229 Cleanup And Status Hardening

Repo behavior after T229:

- Cloudmini provisioning treats only provider resource statuses `running`, `active`, `ready`, and `available` as usable for service activation.
- Cloudmini resource statuses such as `creating`, `provisioning`, `pending`, empty, or unknown produce `PROVIDER_PARTIAL_SUCCESS`, retry safety `manual_review_required`, and no active service or credential record.
- Provider delete uses the Cloudmini V3 `DELETE /api/v3/proxies/:id` path with the job idempotency key and waits for the provider operation to reach `succeeded`.
- Lifecycle-worker termination calls provider `Terminate` before transitioning the Billing service to `terminated` when a provider registry is configured.
- If provider delete is timeout/unknown, missing cleanup metadata, or adapter lookup fails, the lifecycle job moves to manual review and the Billing service is not marked `terminated`.

Operator boundary:

- Use `cmd/worker lifecycle-once` or `cmd/worker lifecycle-loop` with the same approved Cloudmini provider registry env used for provisioning when validating provider-backed cleanup.
- Do not treat `POST /admin/services/:id/terminate` or `POST /reseller/services/:id/terminate` as provider cleanup evidence; those routes remain lifecycle APIs.
- A direct provider delete is allowed only as a documented dev-pilot exception with owner approval, one-resource scope, same-session cleanup, and redacted evidence.

## T230 Target Test-Server Deployment Evidence

T230 verified the T229 hardening on the approved Billing dev test server without mutating Cloudmini or Billing provisioning state.

Target state:

- Runtime path: `/opt/Billing`.
- Environment: `APP_ENV=dev`, `APP_HTTP_ADDR=:8080`.
- Always-on provider mode: `PROVIDER_DEFAULT_MODE=fake`.
- Runtime secrets: `DB_DSN` and `ENCRYPTION_KEY` present in the local env file but not printed.
- Cloudmini dev credential path: `/opt/cred-cloudmini-dev.env`, mode `0600`, owner `root:root`.
- Source markers confirmed on target: `cloudminiV3ProxyStatusUsable` and `NewProviderBackedServiceLifecycleRunner`.

Target validation:

- `go test ./internal/modules/provider ./internal/modules/order ./cmd/worker` passed.
- `go run ./cmd/taskguard` passed.
- `go build -o bin/api ./cmd/api` passed.
- `go build -o bin/worker ./cmd/worker` passed.
- `npm --prefix frontend run build` passed.
- `billing-api`, `billing-worker`, `billing-frontend`, and `cloudflared` were active after restart.
- `GET http://127.0.0.1:8080/healthz` returned HTTP `200`.
- `GET http://127.0.0.1:8080/readyz` returned HTTP `200`.
- `GET http://127.0.0.1:3000/` returned HTTP `200`.
- Ports `8080` and `3000` were listening.

Explicit non-actions:

- No `POST /api/v3/proxies`.
- No `DELETE /api/v3/proxies/:id`.
- No `POST /api/v3/proxies/:id/actions/:action`.
- No Billing checkout, payment, provisioning worker mutation, or lifecycle mutation was run.
- No raw DSN, token, provider group id, provider resource id, provider payload, or proxy credential was printed or stored in repo evidence.

## T231 Non-Mutating Runtime Activation Evidence

T231 added a worker command for activation preflight only:

```bash
cmd/worker provider-registry-check
```

The command:

- refuses `APP_ENV=prod` and `APP_ENV=production`;
- builds the worker provider registry from environment variables;
- does not require `DB_DSN`;
- does not open the database;
- does not claim provisioning or lifecycle jobs;
- does not call Cloudmini `GET`, `POST`, `DELETE`, or action routes;
- prints only redacted counts and boolean evidence.

Target test-server run:

- Runtime path: `/opt/Billing`.
- Environment: `APP_ENV=dev`.
- Credential path: `/opt/cred-cloudmini-dev.env`, mode `0600`, owner `root:root`.
- Command: `PROVIDER_DEFAULT_MODE=cloudmini_v3 go run ./cmd/worker provider-registry-check`.
- Result: `PASS`.
- Adapter: real Cloudmini V3 adapter.
- Source mappings: `1`.
- Account mappings: `0`.
- Provider API called: `no`.
- Mutating routes called: `no`.
- Jobs claimed: `0`.
- Secrets printed: `no`.

This proves config/runtime activation only. It is not proof that a lifecycle cleanup job can safely mutate Cloudmini; that still needs a separate one-resource owner-approved activation window.

## T232 Mutating/Lifecycle Activation Attempt

T232 ran an owner-approved dev activation window on the approved Billing test server. The goal was to create one active Cloudmini-backed Billing service, prepare exactly one termination lifecycle job, and run `cmd/worker lifecycle-once` with the real Cloudmini registry.

Preflight:

- `APP_ENV=dev`.
- Protected credential path `/opt/cred-cloudmini-dev.env`, mode `0600`, owner `root:root`.
- Mapping evidence `PASS` for plan display `10002`, plan-source display `10024`, source display `10012`, source type `cloudmini_v3`, readiness `ready`, priority `1`, and redacted group ref `redacted:c6a7189f0a`.
- One-resource guardrails present: max create `1`, max active resources `1`, worker concurrency `1`, provider rate limit `no-parallel-mutating-calls`, maximum exposure `single-dev-resource`.
- Before mutation: ready jobs `0`, provisioning nonterminal records `0`, Cloudmini active services `0`.

Activation result:

- The always-on `billing-worker` was stopped before payment and restarted after cleanup; it remains in fake-provider mode.
- Billing path created one dev order display `10002`, invoice display `10003`, payment transaction display `10002`, and provider provisioning job display `10002`.
- One-off worker command used `PROVIDER_DEFAULT_MODE=cloudmini_v3` and `cmd/worker provision-once` with batch size `1`.
- Worker result: `claimed=1`, `succeeded=0`, `manual_review=1`, `terminal_failed=0`, `cancelled=0`.
- Provisioning job result: `manual_review`, attempt count `1`, error code `PROVIDER_PARTIAL_SUCCESS`.
- Provider resource lookup by `external_ref` succeeded with redacted resource ref `redacted:6e3d4ecffc7f` and provider status `creating`.
- Same-session fallback cleanup used Cloudmini V3 delete and reached `succeeded`.
- Final provider `GET /api/v3/proxies/:id` returned HTTP `404`.

Lifecycle result:

- No active Billing service was created because T229 correctly treats Cloudmini status `creating` as non-usable.
- `cmd/worker lifecycle-once` was not run for provider cleanup because that would require manually inserting or faking a service record.
- Broader readiness remains blocked until Cloudmini returns a usable status by operation completion or Billing adds an approved wait/read policy that preserves fail-closed behavior.

## T233 Bounded Usable-Status Wait Policy

T233 adds the approved repo-side wait/read policy for Cloudmini create responses that finish the async operation before the proxy itself reports a usable status.

Behavior after T233:

- After `POST /api/v3/proxies` and operation state `succeeded`, Billing reads/polls `GET /api/v3/proxies/:id`.
- Billing only creates an active service after status is `running`, `active`, `ready`, or `available` and credential fields are present.
- If the proxy remains `creating`, empty, unknown, or otherwise non-usable until the configured timeout, Billing returns partial success/manual review and does not create an active service.
- If the proxy becomes usable but credential fields are missing, Billing fails closed with credential-missing handling.
- The target rerun after deployment proved this path for one dev resource without broadening allowed statuses.

## T233 Target Lifecycle Activation Rerun

After PR #498 merged, T233 was deployed to the approved Billing test server and rerun in a new one-resource activation window.

Preflight and deployment:

- Mapping evidence still passed for plan display `10002`, plan-source display `10024`, source display `10012`, and redacted group ref `redacted:c6a7189f0a`.
- Ready provider/lifecycle jobs were `0`, non-terminated Cloudmini services were `0`, and due lifecycle candidates were `0`.
- Target tests/build passed for `go test ./internal/modules/provider ./internal/modules/order ./cmd/worker`, `go build -o bin/api ./cmd/api`, and `go build -o bin/worker ./cmd/worker`.

Activation result:

- Existing dev wallet balance was used because top-up review authorization is tracked separately in T234.
- The always-on fake worker was stopped before payment and restored active after cleanup.
- Billing path created order display `10003`, invoice display `10004`, payment transaction display `10003`, and provider job display `10003`.
- `provision-once` ran with `PROVIDER_DEFAULT_MODE=cloudmini_v3`, batch size `1`, and returned `claimed=1`, `succeeded=1`, `manual_review=0`.
- Service display `10002` became `active`/`paid` with provider resource ref `redacted:52be893bcb0f`.
- Credential storage check passed: one active credential, encrypted payload present, masked hint present.
- One lifecycle terminate job display `10004` was prepared and `lifecycle-once` returned `claimed=1`, `succeeded=1`, `manual_review=0`.
- Final service display `10002` was `terminated`; final provider `GET /api/v3/proxies/:id` returned HTTP `404`.
- No raw DSN, token, provider id, provider payload, raw group id, or proxy credential was printed or recorded.

## Required Preflight

Run these before enabling a mutating pilot:

```bash
go test ./...
go run ./cmd/taskguard
```

Then rerun read-only provider evidence from the local dev credential file:

```bash
set -a
. /opt/cred-cloudmini-dev.env
set +a
VPM_BILLING_V3_BASE_URL="$CLOUDMINI_V3_BASE_URL" \
VPM_BILLING_V3_AUTH_HEADER="Authorization" \
VPM_BILLING_V3_USER_AGENT="$CLOUDMINI_V3_USER_AGENT" \
VPM_BILLING_API_TOKEN="$CLOUDMINI_V3_API_TOKEN" \
/tmp/proxy-cloudmini-billing-edge/scripts/check-billing-v3-edge.sh
```

The read-only result must show:

- capabilities HTTP `200` and `success=true`;
- `ipv4_dc` inventory HTTP `200` and `success=true`;
- selected group ref still sellable with positive allocatable units;
- no token, raw auth header, raw group id, proxy credential, or raw provider payload in captured evidence.

## Mutating Pilot Boundary

The first mutating pilot must run through the Billing checkout/provisioning path, not an ad hoc direct provider `POST`, unless Engineering and Security explicitly approve direct provider testing.

The first run must create at most one provider resource. It must capture only redacted evidence for:

- order display ID;
- provider source display ID;
- provisioning job display ID;
- idempotency key presence, not raw token;
- provider operation/result state;
- redacted external resource reference;
- service active state;
- encrypted credential storage;
- credential reveal audit, if reveal is tested;
- cleanup operation and final provider state.

## Stop Conditions

Stop immediately and do not retry automatically if any of these occurs:

- provider returns auth/permission failure;
- provider returns rate limit or gateway block;
- create request times out after being sent;
- operation id is returned but polling does not finish;
- provider returns a resource without credential data;
- Billing records a manual review status;
- a duplicate resource is suspected;
- cleanup/delete does not complete;
- wallet/ledger/reconciliation mismatch appears;
- any raw secret, proxy credential, or provider payload is exposed.

## Cleanup Procedure

Cleanup must happen in the same pilot session:

Cleanup owner for selected pilot runs is `Admin` per T253. If cleanup cannot be confirmed, `Admin` owns source disable, residual-risk decision, and incident/follow-up creation before any further create attempt.

1. Record the redacted external resource reference.
2. Prefer the lifecycle worker provider-backed cleanup path when the service is eligible for lifecycle termination, with the Cloudmini provider registry configured.
3. Poll provider operation status until terminal state.
4. Confirm the resource is deleted, disabled, or otherwise no longer billable.
5. Confirm Billing service state and provider mapping do not imply an active paid resource after cleanup.
6. Record cleanup owner, time, result, and residual risk.

If the lifecycle worker path is not applicable for a dev pilot, an approved direct provider delete may be used as the fallback cleanup exception. Record the owner approval and keep broader readiness blocked until target-environment lifecycle cleanup evidence exists.

If cleanup fails, keep the launch decision `NO-GO`, disable the source, and open an incident/follow-up before any further create attempt.

## Remaining Code/Config Work

Before broader pilot or multiple provider accounts:

- T217 supports multiple Cloudmini V3 endpoint/API-key mappings through `CLOUDMINI_V3_MAPPINGS_JSON`; keep secret values in approved env/secret storage only.
- T220 verifies the dev pilot mapping. Any broader staging or production-equivalent mapping still needs an approved target environment and owner sign-off before use.
- T227 makes runtime configuration fail closed when the configured source id does not match the Billing provider source used by the provisioning job.
- T228 proves one controlled dev create/delete pilot.
- T229 resolves the repo-side lifecycle cleanup and terminal resource status residual risks with fail-closed code and tests. Broader pilot still needs live duplicate/timeout evidence, shared secret ownership, target-environment lifecycle-worker cleanup evidence, and named owner sign-off.
