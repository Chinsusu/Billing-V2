# T242 - Target staging-equivalent backup restore evidence

Status: DONE
Owner: Codex
Branch: codex/t242-target-backup-restore-evidence
PR: https://github.com/Chinsusu/Billing-V2/pull/516
Risk: operations/database-safety/backup-restore/launch-readiness/secrets
Created: 2026-05-17
Updated: 2026-05-17

## Summary

Capture redacted backup/restore evidence from the approved test server as staging-equivalent launch evidence.

## Scope

- Run the existing guarded backup/restore drill against non-production target-server databases only.
- Use temporary source and restore databases on the approved test server so `dev-db` strict seed smoke has a clean staging-equivalent baseline.
- Record plan, restore, checksum, smoke, cleanup/retention, and `Admin` owner review status without DSNs, passwords, dump files, tokens, provider payloads, or customer data.
- Update launch evidence docs while keeping the overall Go/No-Go decision honest.

## Acceptance Criteria

- Source and target classifications are recorded as non-production/staging-equivalent.
- Restore target overwrite is explicitly bounded to the temporary target database.
- `dev-db` smoke passes against the restored target.
- Backup artifact handling and cleanup/retention owner are recorded.
- Launch docs reflect the new evidence and any remaining blockers without marking GO prematurely.
- Task board validation and whitespace checks pass.

## Notes

- User-provided owner assignment from T241: `Admin` owns Ops, QA, Security, and all launch roles.
- Do not print or commit raw DSNs, passwords, dump files, provider credentials, or service credentials.
- Initial attempt using the long-lived target app DB reached `pg_restore`, but restored `dev-db` smoke failed at `wallet ledger projection` because the app DB includes prior dev/test smoke mutations and is not a clean seed baseline. It is not used as pass evidence.

## Agent Log

- 2026-05-17: Task created and claimed by Codex on branch `codex/t242-target-backup-restore-evidence`.
- 2026-05-17: Target server clean staging-equivalent drill passed with temporary DBs `billing_t242_source_20260517134247` and `billing_t242_restore_20260517134247`: source `dev-db` smoke applied 25 migrations and passed 20 checks, restore applied 0 new migrations and passed 20 checks, checksum `be364dcbd3b434402f89bfbfef941d66e96c04e3d88e4d7ef70b91d9b4f0c0e2`, dump/checksum files deleted, and both temporary DBs dropped. No DSNs, passwords, tokens, provider payloads, dump files, or credentials were recorded.
- 2026-05-17: Verified cleanup on the target server: `billing_t242_%` database count `0` and `/tmp/billing-t242-backup-restore` file count `0`.
- 2026-05-17: Local validation passed: `go run ./cmd/taskguard`, `git diff --check`, and staged secret scan for raw DSNs/tokens/password assignments/private keys.
- 2026-05-17: Opened PR #516 and moved task to `REVIEW`.
- 2026-05-17: PR #516 merged into `main`; task marked `DONE`.
