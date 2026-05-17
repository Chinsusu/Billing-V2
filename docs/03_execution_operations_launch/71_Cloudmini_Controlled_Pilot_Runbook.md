# 71 - Cloudmini Controlled Pilot Runbook

**Date:** 2026-05-17
**Scope:** Controlled pre-approval packet for the first Cloudmini V3 mutating pilot.
**Decision:** One controlled dev create/delete pilot has passed. Broader pilot remains blocked until the T228 residual risks are fixed or explicitly accepted.

## Current Safe State

Read-only evidence is complete for the Billing Go-client-style path:

- `GET /api/v3/capabilities` without auth returns app-level HTTP `401`.
- Authenticated `GET /api/v3/capabilities` returns HTTP `200` V3 success envelopes.
- Authenticated `GET /api/v3/inventory/groups?kind=ipv4_dc` returns HTTP `200` V3 success envelopes.
- Authenticated `GET /api/v3/inventory/groups?kind=residential` returns HTTP `200` V3 success envelopes.
- T228 ran one approved dev mutating pilot through Billing checkout/payment/provisioning worker and cleaned up the provider resource in the same session.

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

This unblocks the dev mapping evidence gate only. Keep broader mutating use blocked until owner sign-offs, timeout/idempotency evidence, cleanup hardening, and residual-risk follow-ups below are complete.

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

- Billing service terminate is lifecycle-only today; it does not call provider `DELETE`, so T228 cleanup needed the provider delete path directly.
- Provider `GET` immediately after Billing activation returned resource status `creating` even though the provider operation and worker completed successfully. A follow-up must decide whether Billing should wait for a terminal usable resource status before marking the service active, or document that `creating` after operation success is provider-compatible.
- Duplicate-create and timeout-after-send behavior are still not proven against the live provider.
- Production/shared secret-store owner and named launch owners are still not recorded in repo evidence.

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

1. Record the redacted external resource reference.
2. Call the approved Billing/service cleanup path or approved provider delete path.
3. Poll provider operation status until terminal state.
4. Confirm the resource is deleted, disabled, or otherwise no longer billable.
5. Confirm Billing service state and provider mapping do not imply an active paid resource after cleanup.
6. Record cleanup owner, time, result, and residual risk.

If cleanup fails, keep the launch decision `NO-GO`, disable the source, and open an incident/follow-up before any further create attempt.

## Remaining Code/Config Work

Before broader pilot or multiple provider accounts:

- T217 supports multiple Cloudmini V3 endpoint/API-key mappings through `CLOUDMINI_V3_MAPPINGS_JSON`; keep secret values in approved env/secret storage only.
- T220 verifies the dev pilot mapping. Any broader staging or production-equivalent mapping still needs an approved target environment and owner sign-off before use.
- T227 makes runtime configuration fail closed when the configured source id does not match the Billing provider source used by the provisioning job.
- T228 proves one controlled dev create/delete pilot, but broader pilot must first resolve or accept the lifecycle cleanup and terminal resource status residual risks.
