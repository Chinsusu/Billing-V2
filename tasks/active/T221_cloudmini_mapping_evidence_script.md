# T221 - Cloudmini mapping evidence script

Status: DONE
Owner: Codex
Branch: codex/t221-cloudmini-mapping-evidence-script
PR: https://github.com/Chinsusu/Billing-V2/pull/473
Risk: provider/provisioning/config/database
Created: 2026-05-16
Updated: 2026-05-16

## Summary

Add a read-only operator evidence script for the Cloudmini V3 pilot mapping so an approved non-production Billing DB can be verified without sharing DSNs, provider tokens, or raw provider identifiers in chat or git.

## Scope

- Add a script that verifies the Cloudmini pilot plan resolves to a `cloudmini_v3` source on an approved non-production DB.
- Fix the existing Cloudmini mapping script precheck SQL so its `psql` variables expand correctly.
- Require non-production `APP_ENV`, `DB_DSN`, and explicit operator approval before DB access.
- Output only display IDs, readiness state, source type/status, and redacted guardrail fields.
- Keep the script read-only and avoid checkout, worker, provider API, or mutating Cloudmini routes.
- Update the Cloudmini runbook/evidence docs with the operator command.

## Acceptance Criteria

- Script refuses production environment names.
- Script requires an approval flag and DB DSN before reading the DB.
- Script runs inside a read-only transaction and emits no DSN, token, auth header, raw group ID, raw provider payload, or proxy credential.
- Script fails when the mapped source is missing, not `cloudmini_v3`, not ready, or missing first-pilot guardrails.
- Task board validation and shell/script checks pass.

## Notes

- This does not unblock T220 by itself; an approved non-production Billing `DB_DSN` or operator-run output is still required.
- This task does not run migrations, mapping, checkout, worker provisioning, or provider create/delete.

## Agent Log

- 2026-05-16: Task created and claimed by Codex from latest `origin/main`.
- 2026-05-16: Local validation found the existing T219 mapping script prechecks failed because `psql -c` did not expand `:'plan_code'`; fixed the prechecks to use stdin SQL.
- 2026-05-16: Verified on a temporary local PostgreSQL DB: `go run ./cmd/smoke dev-db`; `scripts/cloudmini_pilot_mapping.sh`; `scripts/cloudmini_mapping_evidence.sh` PASS; expected source-display mismatch returns FAIL. Dropped the temporary DB after validation. No provider API or mutating Cloudmini route was called.
- 2026-05-16: Validation passed: `bash -n scripts/cloudmini_mapping_evidence.sh scripts/cloudmini_pilot_mapping.sh`; production/approval/display-id guard negative checks; `go run ./cmd/taskguard`; `go run ./cmd/migrate validate`; `go test ./...`; `git diff --check`; changed-file secret pattern scan returned no matches.
- 2026-05-16: Opened PR #473 for review.
- 2026-05-16: PR #473 merged; marking task DONE.
