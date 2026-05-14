# 67 - Backup Restore Drill Runbook

**Date:** 2026-05-14  
**Scope:** Repeatable local/sandbox PostgreSQL backup and restore drill for launch readiness evidence.

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
