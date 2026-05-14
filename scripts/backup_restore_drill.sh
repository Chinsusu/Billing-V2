#!/usr/bin/env bash
set -euo pipefail

usage() {
	cat <<'USAGE'
Usage:
  bash scripts/backup_restore_drill.sh --plan
  bash scripts/backup_restore_drill.sh --run

Required for --run:
  BILLING_BACKUP_RESTORE_ENV       local, dev, sandbox, or staging
  BILLING_BACKUP_SOURCE_DSN        PostgreSQL DSN to dump
  BILLING_RESTORE_TARGET_DSN       PostgreSQL DSN to restore into
  BILLING_BACKUP_RESTORE_CONFIRM   restore:<target database name>

Optional:
  BILLING_BACKUP_DRILL_DIR         backup output directory, defaults to /tmp/billing-backup-restore-drill
  BILLING_BACKUP_RESTORE_TIMEOUT   smoke timeout, defaults to 60s
  GO                               Go binary, defaults to go

The target restore is destructive. This script refuses production markers,
source==target, APP_ENV=prod/production, and missing target confirmation.
USAGE
}

log() {
	printf '%s\n' "$*" >&2
}

die() {
	log "error: $*"
	exit 1
}

require_command() {
	command -v "$1" >/dev/null 2>&1 || die "$1 is required on PATH"
}

lowercase() {
	printf '%s' "$1" | tr '[:upper:]' '[:lower:]'
}

ensure_repo_root() {
	[[ -f go.mod && -d migrations && -d cmd/smoke ]] || die "run this script from the repository root"
}

require_non_production_env() {
	local env_name
	env_name="$(lowercase "${BILLING_BACKUP_RESTORE_ENV:-}")"
	case "$env_name" in
	local | dev | sandbox | staging)
		return 0
		;;
	prod | production)
		die "BILLING_BACKUP_RESTORE_ENV must not be prod/production"
		;;
	"")
		die "BILLING_BACKUP_RESTORE_ENV is required"
		;;
	*)
		die "BILLING_BACKUP_RESTORE_ENV must be local, dev, sandbox, or staging"
		;;
	esac
}

reject_production_runtime() {
	local app_env
	app_env="$(lowercase "${APP_ENV:-}")"
	case "$app_env" in
	prod | production)
		die "refusing to run with APP_ENV=${APP_ENV}"
		;;
	esac
}

reject_production_marker() {
	local label="$1"
	local value="$2"
	local lower
	lower="$(lowercase "$value")"
	case "$lower" in
	*production* | *://prod* | */prod* | *-prod* | *_prod* | *.prod* | *prod.* | *prod-* | *prod_* | *live* | *customer* | *rds.amazonaws.com* | *cloudsql*)
		die "$label contains a production-like marker; use an explicit local/sandbox target"
		;;
	esac
}

psql_value() {
	local dsn="$1"
	local sql="$2"
	psql -X -v ON_ERROR_STOP=1 -qAt --dbname="$dsn" --command="$sql" | head -n 1
}

print_plan() {
	ensure_repo_root
	require_command pg_dump
	require_command pg_restore
	require_command psql
	require_command sha256sum
	require_command "${GO:-go}"

	cat <<'PLAN'
Backup/restore drill plan:
1. Create or choose two approved non-production PostgreSQL databases:
   - source: local/dev/sandbox database containing migrated and seeded data
   - target: empty local/dev/sandbox database that can be overwritten
2. Export:
   - BILLING_BACKUP_RESTORE_ENV=local|dev|sandbox|staging
   - BILLING_BACKUP_SOURCE_DSN=<source DSN>
   - BILLING_RESTORE_TARGET_DSN=<target DSN>
3. Run --run once. The script connects to both databases and prints the exact
   BILLING_BACKUP_RESTORE_CONFIRM value required for the target database.
4. Export the confirm value and rerun --run.
5. Capture the backup path, checksum, restore result, and dev-db smoke result
   in docs/03_execution_operations_launch/67_Backup_Restore_Drill_Runbook.md.

The script does not print DSNs and refuses APP_ENV=prod/production, production
markers, identical source/target DSNs, and identical source/target database names.
PLAN
}

run_drill() {
	ensure_repo_root
	require_command pg_dump
	require_command pg_restore
	require_command psql
	require_command sha256sum
	require_command "${GO:-go}"
	require_non_production_env
	reject_production_runtime

	local source_dsn="${BILLING_BACKUP_SOURCE_DSN:-}"
	local target_dsn="${BILLING_RESTORE_TARGET_DSN:-}"
	[[ -n "$source_dsn" ]] || die "BILLING_BACKUP_SOURCE_DSN is required"
	[[ -n "$target_dsn" ]] || die "BILLING_RESTORE_TARGET_DSN is required"
	[[ "$source_dsn" != "$target_dsn" ]] || die "source and target DSNs must be different"
	reject_production_marker "source DSN" "$source_dsn"
	reject_production_marker "target DSN" "$target_dsn"

	local source_db
	local target_db
	source_db="$(psql_value "$source_dsn" "SELECT current_database();")"
	target_db="$(psql_value "$target_dsn" "SELECT current_database();")"
	[[ -n "$source_db" ]] || die "could not resolve source database name"
	[[ -n "$target_db" ]] || die "could not resolve target database name"
	[[ "$source_db" != "$target_db" ]] || die "source and target database names must be different"
	reject_production_marker "source database" "$source_db"
	reject_production_marker "target database" "$target_db"

	local expected_confirm="restore:${target_db}"
	if [[ "${BILLING_BACKUP_RESTORE_CONFIRM:-}" != "$expected_confirm" ]]; then
		die "set BILLING_BACKUP_RESTORE_CONFIRM=${expected_confirm} to overwrite target database ${target_db}"
	fi

	local backup_dir="${BILLING_BACKUP_DRILL_DIR:-/tmp/billing-backup-restore-drill}"
	local timestamp
	timestamp="$(date -u +%Y%m%dT%H%M%SZ)"
	local backup_file="${backup_dir}/billing-${source_db}-${timestamp}.dump"
	local checksum_file="${backup_file}.sha256"

	umask 077
	mkdir -p "$backup_dir"

	log "source database: ${source_db}"
	log "target database: ${target_db}"
	log "backup file: ${backup_file}"
	pg_dump --format=custom --no-owner --no-acl --dbname="$source_dsn" --file="$backup_file"
	sha256sum "$backup_file" >"$checksum_file"
	log "checksum file: ${checksum_file}"

	log "restoring into target database: ${target_db}"
	pg_restore --clean --if-exists --no-owner --no-acl --dbname="$target_dsn" "$backup_file"

	log "running dev-db smoke against restored target"
	"${GO:-go}" run ./cmd/smoke -dsn "$target_dsn" -timeout "${BILLING_BACKUP_RESTORE_TIMEOUT:-60s}" dev-db
	log "backup/restore drill passed"
}

mode="${1:---run}"
case "$mode" in
--help | -h | help)
	usage
	;;
--plan | plan)
	print_plan
	;;
--run | run)
	run_drill
	;;
*)
	usage
	die "unknown mode: $mode"
	;;
esac
