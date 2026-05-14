# T209 - Staging backup restore evidence packet

Status: DONE
Owner: Codex
Branch: codex/t209-staging-restore-evidence
PR: https://github.com/Chinsusu/Billing-V2/pull/449
Risk: operations, database safety, backup/restore, launch readiness, and secret handling
Created: 2026-05-14
Updated: 2026-05-14

## Summary

Turn the staging backup/restore blocker into a concrete redacted evidence packet that can be filled only after an approved shared staging or staging-equivalent target is available.

## Scope

- Extend the backup/restore drill runbook with staging evidence requirements, approval fields, redaction rules, and pass/fail criteria.
- Keep launch readiness blocked until an approved non-production staging target is restored and smoke-verified.
- Update launch audit/go-no-go docs only enough to reference the new evidence packet.
- Do not run destructive restore commands, add DSNs, or claim that shared staging restore evidence exists.

## Acceptance Criteria

- Staging backup/restore evidence requirements are explicit for source/target approval, destructive restore confirmation, checksum, smoke result, cleanup, and owner sign-off.
- Launch docs still say backup/restore is incomplete for launch until shared staging/non-production evidence is provided.
- Task board and whitespace validation pass.

## Notes

- No approved shared staging DSN, target database, operator sign-off, or staging restore evidence is present in this repository as of task creation.

## Agent Log

- 2026-05-14: Codex created and claimed task on `codex/t209-staging-restore-evidence`.
- 2026-05-14: Local validation passed: `go run ./cmd/taskguard`, `git diff --check`, `bash scripts/backup_restore_drill.sh --plan`. Destructive restore was not run because no approved shared staging/non-production DSNs or target overwrite approval are present in repo.
- 2026-05-14: Opened PR #449 for review.
- 2026-05-14: PR #449 merged into `main` with squash commit `c228f828ebffaf2f6b8bc0a634258519f8aef4b2`; marking task DONE.
