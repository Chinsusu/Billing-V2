# T219 - Cloudmini pilot source mapping

Status: IN_PROGRESS
Owner: Codex
Branch: codex/t219-cloudmini-pilot-mapping
PR: -
Risk: provider/provisioning/credential/config/database
Created: 2026-05-16
Updated: 2026-05-16

## Summary

Add guarded tooling and docs for mapping the approved non-production Cloudmini V3 pilot source to the proxy plan without committing credentials or calling mutating provider APIs.

## Scope

- Add the missing database provider type support needed for a `cloudmini_v3` provider source.
- Add a guarded non-production script that creates/updates the Cloudmini V3 provider source and plan-source mapping only after explicit local approval.
- Keep raw API tokens, auth headers, raw provider group IDs, provider payloads, and DSNs out of git, task notes, PR text, and command output.
- Do not run Cloudmini create/delete or Billing checkout/provisioning mutating pilot in this task.
- Do not implement T217 multiple endpoint/API-key runtime config in this task.

## Acceptance Criteria

- Migration validation accepts the new Cloudmini provider type migration.
- The mapping script refuses to run without a non-production environment, `DB_DSN`, and explicit pilot approval.
- The mapping script does not print secrets, raw group IDs, or provider payloads.
- Docs identify how to apply and verify the mapping before the first mutating pilot.
- Task board validation passes.

## Notes

- Existing code has `provider.TypeCloudminiV3`, but the catalog enum must also support `cloudmini_v3` before a real provider source can be inserted.
- The selected pilot group remains local-only in `/opt/cred-cloudmini-dev.env`; use only a redacted group reference in repo docs.
- The first mutating pilot remains blocked until approval fields, quota/cleanup owners, read-only recheck, and source readiness verification are complete.

## Agent Log

- 2026-05-16: Task created and claimed by Codex from latest `origin/main`.
- 2026-05-16: Added `cloudmini_v3` enum migration, guarded pilot mapping script, and runbook/evidence updates. No provider create/delete or Billing mutating pilot was run.
- 2026-05-16: Local `/opt/cred-cloudmini-dev.env` metadata was aligned to the pilot source ID and guardrails without persisting the approval flag.
- 2026-05-16: Validation passed: `bash -n scripts/cloudmini_pilot_mapping.sh`; `go run ./cmd/migrate validate`; `go test ./cmd/migrate ./internal/platform/db ./internal/seed`; `go test ./...`; `go run ./cmd/taskguard`; `git diff --check`; guard checks for production refusal and missing approval refusal.
