# T220 - Cloudmini dev mapping evidence

Status: IN_PROGRESS
Owner: Codex
Branch: codex/t220-cloudmini-dev-mapping-evidence-unblock
PR: -
Risk: provider/provisioning/credential/config/database
Created: 2026-05-16
Updated: 2026-05-16

## Summary

Apply and verify the guarded Cloudmini V3 pilot source mapping on an approved non-production dev database, recording only redacted readiness evidence.

## Scope

- Discover a non-production dev DB access path from approved local secret material without printing secret values.
- Apply migrations and the T219 mapping script only against non-production.
- Verify the proxy pilot plan resolves to a `cloudmini_v3` source with readiness `ready` using display IDs/redacted references only.
- Update the Cloudmini readiness/runbook docs with the applied mapping evidence.
- Do not call Cloudmini mutating routes, Billing checkout, Billing payment, or provisioning worker pilot in this task.
- Do not implement T217 multiple endpoint/API-key runtime config in this task.

## Acceptance Criteria

- Mapping apply is blocked or executed only with explicit non-production guardrails.
- Evidence records plan/source display IDs, readiness state, source type, and guardrail summary without secrets or raw provider IDs.
- If dev DB access is unavailable, record the blocker precisely and keep mutating pilot blocked.
- Task board validation and relevant local checks pass.

## Notes

- `/opt/cred-cloudmini-dev.env` contains dev/server/provider secret material; only key names or redacted evidence may be printed.
- T219 already added the migration and mapping script; this task is for target-environment application/evidence, not new runtime provider behavior.
- The first create/delete pilot remains blocked until owner approval fields, quota/cleanup owner, timeout/idempotency evidence, and same-session cleanup procedure are complete.

## Evidence

T220 was unblocked by the approved non-production Billing test server runtime env at `/opt/Billing/.env.dev`.

- Confirmed `APP_ENV=dev` and `DB_DSN` presence without printing the DSN.
- `go run ./cmd/migrate validate` reported `25` migrations.
- `go run ./cmd/migrate plan` reported `0` pending migrations.
- `go run ./cmd/migrate up` applied `0` migrations.
- `BILLING_CLOUDMINI_PILOT_APPROVED=yes bash scripts/cloudmini_pilot_mapping.sh` applied the pilot mapping and printed only display IDs/status fields.
- Mapping output: plan `proxy-static-10gb-monthly`, plan-source display `10024`, source display `10012`, source type `cloudmini_v3`, source status `active`, inventory mode `provider_live`, plan-source status `active`, priority `1`.
- `BILLING_CLOUDMINI_EVIDENCE_APPROVED=yes bash scripts/cloudmini_mapping_evidence.sh` returned `result=PASS`.
- Evidence output: plan display `10002`, product type `proxy`, readiness `ready`, pilot max create `1`, pilot max active resources `1`, worker concurrency `1`, provider rate limit `no-parallel-mutating-calls`, max exposure `single-dev-resource`, provider kind `ipv4_dc`, provider group ref `redacted:c6a7189f0a`, protocol `socks5`, failed checks `none`.
- No Cloudmini create/delete/action route, Billing checkout, Billing payment, or provisioning worker pilot was run.
- No DSN, token, raw provider group ID, raw provider payload, proxy credential, or auth header is recorded in this task.

## Remaining Mutating Pilot Blockers

- Approved shared secret-store path/owner is still missing; `/opt/cred-cloudmini-dev.env` remains dev-local only.
- Named provider/account owner, cleanup owner, security owner, finance/quota owner, and reviewer sign-off are still missing.
- Timeout-after-send, duplicate create/idempotency, redacted mutating error examples, and cleanup/delete evidence are still missing.
- The first create/delete pilot remains blocked until those fields are complete and the run uses the documented one-create, one-active-resource, one-worker-concurrency guardrails.

## Agent Log

- 2026-05-16: Task created and claimed by Codex from latest `origin/main`.
- 2026-05-16: Checked local secret key names only; no `DB_DSN` exists in `/opt/cred-cloudmini-dev.env`.
- 2026-05-16: Verified SSH to dev host with local secret material without printing host/user/password.
- 2026-05-16: Discovered no Billing repo or Billing DB env on the dev host; did not apply migrations or mapping to an unverified DB.
- 2026-05-16: Local default Billing DSN check failed because the `billing` database does not exist on this runner.
- 2026-05-16: Validation passed for blocker record: `go run ./cmd/taskguard`; `git diff --check`; changed-file secret pattern scan returned no matches.
- 2026-05-16: Task unblocked by approved non-production test server Billing runtime env at `/opt/Billing/.env.dev`; starting guarded mapping/evidence run without printing DSN or provider secrets.
- 2026-05-16: Test server run passed: `APP_ENV=dev`, migration validate/plan/up, guarded mapping apply, and read-only mapping evidence `PASS`. No provider mutating route, Billing checkout, payment, or worker pilot was run.
