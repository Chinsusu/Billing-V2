# T220 - Cloudmini dev mapping evidence

Status: BLOCKED
Owner: Codex
Branch: codex/t220-cloudmini-dev-mapping-evidence
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
- The first create/delete pilot remains blocked until this mapping evidence, owner approval fields, quota/cleanup owner, and read-only recheck are complete.

## Blocker

No approved Billing DB access path was found.

- Local `/opt/cred-cloudmini-dev.env` has dev SSH/provider keys but no `DB_DSN`.
- SSH to the dev host works, but discovery found the provider/manager repo under `/opt`, not a Billing repo or Billing runtime env.
- DB key-name search on the dev host found only provider/manager deployment scripts, not an approved Billing DB target.
- The local default Billing DSN from the runbook points to a database that does not exist on this runner.
- Applying Billing migrations or pilot mapping to the provider/manager DB would be unsafe, so no DB mapping was applied.

Needed to unblock:

- Provide an approved non-production Billing `DB_DSN`, or install/run Billing on the dev host and record the approved env path.
- Confirm `APP_ENV` is `dev`, `staging`, `sandbox`, or `local`.
- Then rerun T219 migration + `scripts/cloudmini_pilot_mapping.sh` and capture readiness evidence using display IDs only.

## Agent Log

- 2026-05-16: Task created and claimed by Codex from latest `origin/main`.
- 2026-05-16: Checked local secret key names only; no `DB_DSN` exists in `/opt/cred-cloudmini-dev.env`.
- 2026-05-16: Verified SSH to dev host with local secret material without printing host/user/password.
- 2026-05-16: Discovered no Billing repo or Billing DB env on the dev host; did not apply migrations or mapping to an unverified DB.
- 2026-05-16: Local default Billing DSN check failed because the `billing` database does not exist on this runner.
- 2026-05-16: Validation passed for blocker record: `go run ./cmd/taskguard`; `git diff --check`; changed-file secret pattern scan returned no matches.
