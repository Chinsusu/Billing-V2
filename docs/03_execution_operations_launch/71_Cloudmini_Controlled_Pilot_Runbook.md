# 71 - Cloudmini Controlled Pilot Runbook

**Date:** 2026-05-16
**Scope:** Controlled pre-approval packet for the first Cloudmini V3 mutating pilot.
**Decision:** Not approved for create/delete until every approval field below is complete.

## Current Safe State

Read-only evidence is complete for the Billing Go-client-style path:

- `GET /api/v3/capabilities` without auth returns app-level HTTP `401`.
- Authenticated `GET /api/v3/capabilities` returns HTTP `200` V3 success envelopes.
- Authenticated `GET /api/v3/inventory/groups?kind=ipv4_dc` returns HTTP `200` V3 success envelopes.
- Authenticated `GET /api/v3/inventory/groups?kind=residential` returns HTTP `200` V3 success envelopes.
- No mutating provider route has been called from Billing evidence.

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

T220 target-environment discovery did not find an approved Billing DB access path:

- `/opt/cred-cloudmini-dev.env` has provider/dev SSH keys but no `DB_DSN`.
- The reachable dev host did not contain a Billing repo or Billing runtime environment under `/opt`.
- DB key-name discovery found only provider/manager deployment scripts, not an approved Billing DB target.
- The local default Billing DSN from the runbook was not available on this runner.
- No migration, mapping script, checkout, worker, or provider mutating call was run against an unverified DB.

Keep the pilot blocked until an approved non-production Billing `DB_DSN` or equivalent operator-run evidence is provided.

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
- The approved dev/staging Billing database still needs the T219 mapping script to be applied and verified through provider readiness evidence.
- Runtime configuration must fail closed when the configured source id does not match the Billing provider source used by the provisioning job.
