# T226 - Cloudmini mutating pilot preflight

Status: IN_PROGRESS
Owner: Codex
Branch: codex/t226-cloudmini-mutating-pilot-preflight
PR: -
Risk: provider/provisioning/credential/config/database
Created: 2026-05-16
Updated: 2026-05-16

## Summary

Add a non-mutating Cloudmini V3 pilot preflight guard that verifies approval, owner, quota, cleanup, credential-path, and mapping evidence prerequisites before any create/delete pilot is attempted.

## Scope

- Add a local/operator script that refuses production environments.
- Require explicit preflight approval and all launch owner/cleanup/quota fields before passing.
- Require the credential path to be outside the repo, readable, and not group/world accessible.
- Reuse the T221 read-only mapping evidence collector to prove the selected pilot source is still `cloudmini_v3` and `ready`.
- Output only redacted readiness/guardrail evidence and presence flags.
- Do not call Cloudmini create/delete/action routes, Billing checkout, Billing payment, or provisioning worker pilot in this task.
- Do not fill real owner sign-offs with assumptions.

## Acceptance Criteria

- Preflight fails without non-production `APP_ENV`, `DB_DSN`, explicit approval, owner fields, cleanup fields, or exact one-resource guardrails.
- Preflight fails if the credential path is missing, inside the repo, or group/world accessible.
- Preflight fails if mapping evidence does not return `result=PASS`.
- Preflight output records no DSN, token, raw provider group ID, raw provider payload, proxy credential, auth header, or owner contact detail.
- Runbook documents the command and makes clear this is not approval to run mutating provider routes.
- Task board validation and relevant shell checks pass.

## Notes

- This task only guards the step before a mutating pilot. It does not run a mutating pilot.
- The first create/delete pilot remains blocked until a real operator supplies owner/sign-off values and accepts the documented cleanup procedure.

## Agent Log

- 2026-05-16: Task created and claimed by Codex from latest `origin/main`.
- 2026-05-16: Added non-mutating preflight script and runbook docs. The script requires non-production env, explicit approval, owner/cleanup/read-only evidence fields, private credential path, exact one-resource guardrails, and T221 mapping evidence `PASS`.
- 2026-05-16: Local negative checks passed for production refusal, missing approval, missing owner field, guardrail greater than one, and group/world-accessible credential path.
- 2026-05-16: Test server positive validation passed with a temporary mode `0600` placeholder credential file and existing dev DB mapping evidence. No provider mutating route, Billing checkout, payment, or worker pilot was run.
- 2026-05-16: Validation passed: `go run ./cmd/taskguard`; `bash -n scripts/cloudmini_mutating_pilot_preflight.sh scripts/cloudmini_mapping_evidence.sh scripts/cloudmini_pilot_mapping.sh`; `git diff --check`; changed-file secret pattern scan returned no matches.
