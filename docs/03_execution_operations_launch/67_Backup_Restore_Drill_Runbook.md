# 67 - Backup Restore Drill Runbook

**Date:** 2026-05-14  
**Scope:** Repeatable local/sandbox/staging-equivalent PostgreSQL backup and restore drill for launch readiness evidence.

## Purpose

This runbook proves the platform can restore a database backup into an approved non-production target and still pass the local seeded DB smoke. It does not approve production restore, production data handling, or production secret restore.

## Safety Boundary

Allowed targets:

```text
local
dev
sandbox
staging
```

Forbidden:

```text
production DSNs
production data copied into local/dev without masking approval
source and target pointing at the same database
shared staging target without owner approval
committed dump files, raw DSNs, passwords, provider credentials, or customer data
```

The drill script refuses `APP_ENV=prod`, `APP_ENV=production`, production-like DSN markers, identical source/target DSNs, and identical source/target database names. Backup files are written with `umask 077` and default to `/tmp/billing-backup-restore-drill`.

## Prerequisites

```text
pg_dump
pg_restore
psql
sha256sum
go
two approved non-production PostgreSQL databases
```

The source database should already pass:

```bash
go run ./cmd/smoke -dsn "$BILLING_BACKUP_SOURCE_DSN" dev-db
```

## Plan Command

Use the non-destructive plan first:

```bash
make backup-restore-drill-plan
```

## Run Command

Export non-production DSNs without printing them to shared logs:

```bash
export BILLING_BACKUP_RESTORE_ENV=local
export BILLING_BACKUP_SOURCE_DSN='<source local/sandbox DSN>'
export BILLING_RESTORE_TARGET_DSN='<target local/sandbox DSN>'
```

Run once to resolve the target database name and get the exact confirmation value:

```bash
make backup-restore-drill
```

Then confirm the destructive target restore and rerun:

```bash
export BILLING_BACKUP_RESTORE_CONFIRM='restore:<target database name>'
make backup-restore-drill
```

The script performs:

```text
pg_dump custom-format backup of source
sha256 checksum capture
pg_restore --clean --if-exists into target
go run ./cmd/smoke -dsn "$BILLING_RESTORE_TARGET_DSN" dev-db
```

## Acceptance Criteria

The drill is valid only when all items are true:

- Source and target are approved non-production databases.
- Backup file and checksum are captured outside the repository.
- Restore completes without `pg_restore` errors.
- Restored target passes `dev-db` smoke.
- Smoke verifies migrations, seed idempotency, tenant/user/catalog/wallet/order/service/invoice/ledger/payment records.
- Operator records evidence without secrets, raw DSNs, provider credentials, or customer data.

## Shared Staging Evidence Status

T242 records a target-server staging-equivalent restore drill on 2026-05-17 using temporary non-production source and restore databases on the approved test server. This proves the target runner, PostgreSQL tooling, backup artifact handling, restore path, and `dev-db` smoke path for a clean staging-equivalent database. It does not prove restore of the long-lived target app database, because that database contains prior dev/test smoke mutations and failed the strict seed-baseline `dev-db` smoke after restore. Use the T242 result as launch evidence only if Ops/QA/Security accept the temporary staging-equivalent scope; otherwise run an additional drill against an approved clean shared staging source.

| Evidence area | Required proof | Current repo status |
|---|---|---|
| Environment approval | Ops owner confirms source and target are non-production and target may be overwritten. | T242 pass for temporary target-server staging-equivalent source/restore DBs; final scope acceptance pending Admin/Ops/QA/Security review. |
| Data classification | Source data classification, masking status if any copied data exists, and confirmation that no production customer data is used without approval. | T242 source and target were temporary seeded dev/test DBs with no production data. |
| Tooling prerequisites | `pg_dump`, `pg_restore`, `psql`, `sha256sum`, and `go` versions or runner image. | T242 target check found all required tools present on the test server. |
| Plan run | `make backup-restore-drill-plan` or equivalent dry-run reviewed before destructive restore. | T242 ran `bash scripts/backup_restore_drill.sh --plan` before restore. |
| Backup artifact | Redacted backup path outside the repository and checksum path or checksum value. | T242 wrote the dump/checksum under `/tmp/billing-t242-backup-restore`, recorded checksum `be364dcbd3b434402f89bfbfef941d66e96c04e3d88e4d7ef70b91d9b4f0c0e2`, then deleted both files. |
| Restore confirmation | Target database name and `BILLING_BACKUP_RESTORE_CONFIRM` value captured without DSN, password, or host secret. | T242 confirm probe required `restore:billing_t242_restore_20260517134247`; no DSN/password was recorded. |
| Restore result | Restore completed without `pg_restore` errors against the approved target. | T242 restore completed and reported `backup/restore drill passed`. |
| Smoke result | `dev-db` smoke passed against the restored target with migration/check counts recorded. | T242 restored target applied 0 new migrations and passed 20 `dev-db` smoke checks after source passed 25 migrations and 20 checks. |
| Cleanup/retention | Backup artifact retention or deletion decision, target cleanup owner, and follow-up issues. | T242 deleted dump/checksum files and dropped both temporary databases after evidence capture; cleanup owner is Admin by T241 assignment. |
| Sign-off | Operator, Ops owner, QA owner, and date/time UTC. | Operator: Codex. Assigned Ops/QA/Security owner: Admin by T241. Final launch sign-off remains pending evidence packet review. |

## Shared Staging Evidence Packet Template

Use this packet only after the source and target are approved non-production databases. Store redacted values in git and keep DSNs/passwords in the approved secret channel only.

```text
Drill ID:
Date/time UTC:
Operator:
Ops owner approval:
QA reviewer:
Environment: staging/sandbox/staging-equivalent
Source classification:
Target classification:
Target overwrite approval:
Production data present: no/approved masked exception
Masking approval reference:
Tooling runner:
pg_dump version:
pg_restore version:
psql version:
Go version:
Plan command:
Plan result:
Restore command:
Target database name:
Confirm value used:
Backup artifact path: redacted non-repo path
Backup checksum:
Restore result:
Smoke command:
Smoke result:
Migration count:
Smoke check count:
Backup artifact retention/deletion:
Target cleanup owner:
Issues:
Follow-up:
Sign-off decision: pass/fail
```

## Shared Staging Pass Criteria

Do not use this gate as launch evidence unless all items are true:

- Source and target approvals are recorded before the destructive restore.
- Source and target are different databases and both are non-production or approved staging-equivalent.
- No raw DSN, password, token, credential, provider secret, dump file, or customer data is committed.
- The backup artifact is stored outside the repository with restricted access.
- Checksum is captured and tied to the drill ID.
- Restore completes without `pg_restore` errors.
- Restored target passes `go run ./cmd/smoke -dsn "$BILLING_RESTORE_TARGET_DSN" dev-db`.
- Operator records migration and smoke check counts.
- Cleanup/retention owner is recorded.
- Ops and QA sign off on the evidence packet.

If any item is missing, keep the backup/restore launch gate blocked or partial.

## Evidence Template

```text
Drill ID:
Date/time UTC:
Operator:
Environment: local/dev/sandbox/staging
Source classification: local/dev/sandbox/staging, no production data
Target classification: local/dev/sandbox/staging, approved to overwrite
Backup file location: /tmp/... (no repository path)
Backup checksum file: /tmp/...
Restore command: make backup-restore-drill
Smoke command: go run ./cmd/smoke -dsn "$BILLING_RESTORE_TARGET_DSN" dev-db
Result: pass/fail
Issues:
Follow-up:
```

## T203 Repository Evidence

```text
Drill ID: T203-local-20260514T021629Z
Date/time UTC: 2026-05-14 02:16:29
Operator: Codex
Environment: local
Source classification: local temporary seed database, no production data
Target classification: local temporary restore database, approved to overwrite
Source database: billing_t203_src_20260514015919
Target database: billing_t203_dst_20260514015919
Backup file location: /tmp/billing-backup-restore-drill/billing-billing_t203_src_20260514015919-20260514T021629Z.dump
Backup checksum: fd93d5684bbe9c880397eef80c1cbc042ef99d20a6c7dc2bd95751b0e17edb34
Source smoke: dev-db passed, 23 migrations, 19 checks
Restore command: bash scripts/backup_restore_drill.sh --run with redacted local DSNs
Restored target smoke: dev-db passed, 0 new migrations, 19 checks
Result: pass
Cleanup: temporary source/target databases and local dump/checksum files removed after evidence capture
Issues: none
Follow-up: repeat against approved shared staging before final T205 launch sign-off
```

## T242 Target Staging-Equivalent Evidence

```text
Drill ID: T242-target-20260517T134247Z
Date/time UTC: 2026-05-17 13:42:47
Operator: Codex
Environment: staging-equivalent test server
Ops owner approval: Admin assigned by T241; final scope acceptance pending evidence review
QA reviewer: Admin assigned by T241; final scope acceptance pending evidence review
Security reviewer: Admin assigned by T241; final scope acceptance pending evidence review
Source classification: temporary target-server staging-equivalent seed DB, no production data
Target classification: temporary target-server staging-equivalent restore DB, approved to overwrite
Source database: billing_t242_source_20260517134247
Target database: billing_t242_restore_20260517134247
Production data present: no
Plan command: bash scripts/backup_restore_drill.sh --plan
Plan result: pass
Restore command: bash scripts/backup_restore_drill.sh --run with redacted target-server DSNs
Confirm value used: restore:billing_t242_restore_20260517134247
Backup artifact path: /tmp/billing-t242-backup-restore/billing-billing_t242_source_20260517134247-20260517T134248Z.dump
Backup checksum: be364dcbd3b434402f89bfbfef941d66e96c04e3d88e4d7ef70b91d9b4f0c0e2
Source smoke: dev-db passed, 25 migrations, 20 checks
Restore result: pass; pg_restore completed without error
Restored target smoke: dev-db passed, 0 new migrations, 20 checks
Backup artifact retention/deletion: dump and checksum files deleted after evidence capture
Target cleanup owner: Admin
Cleanup result: source and target temporary databases dropped
Issues: the long-lived target app DB was not used as pass evidence because strict seed-baseline dev-db smoke failed after restore due prior dev/test smoke mutations.
Follow-up: if launch requires proof against the long-lived app DB snapshot, run an approved app-DB restore drill with a smoke designed for long-lived staging data or record an owner-approved staging-equivalent exception.
Sign-off decision: evidence collected; final Admin/Ops/QA/Security launch sign-off pending
```
