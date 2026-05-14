# T203 - Backup and restore ops drill

Status: IN_PROGRESS
Owner: Codex
Branch: codex/t203-backup-restore-ops-drill
PR: -
Risk: operations, database safety, launch readiness, and disaster recovery
Created: 2026-05-13
Updated: 2026-05-13

## Summary

Add and verify a backup/restore drill for launch readiness.

## Scope

- Document or script a local/sandbox backup and restore procedure for the database.
- Validate restore on a non-production database.
- Record operator checklist and expected evidence for launch sign-off.
- Do not touch production data or production DSNs.

## Acceptance Criteria

- Backup and restore steps are documented and reproducible in local/sandbox.
- Restore validation confirms migrations and core smoke can run after restore.
- Docs clearly forbid production DSNs and unmasked production data.
- Relevant docs validation, DB smoke where available, and CI pass.

## Notes

- Stop before any destructive database command unless the target is explicitly safe and approved.

## Agent Log

- 2026-05-13: Task created by Codex backlog planning.
- 2026-05-14: Codex claimed task on `codex/t203-backup-restore-ops-drill`.
- 2026-05-14: Added guarded backup/restore drill script, Makefile targets, and DR evidence runbook.
- 2026-05-14: Ran local non-production drill: source seed `dev-db` smoke passed 19 checks, restore target `dev-db` smoke passed 19 checks, temp DBs and dump files cleaned up.
