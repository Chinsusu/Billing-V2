#!/usr/bin/env bash
set -euo pipefail

usage() {
	cat <<'USAGE'
Usage:
  bash scripts/full_e2e_quality_gate.sh

Required:
  DB_DSN                         local/dev PostgreSQL DSN

Optional:
  APP_ENV                        local or dev, defaults to local
  BILLING_E2E_API_ADDR           local API listen address, defaults to 127.0.0.1:18080
  API_BASE_URL                   local API URL, defaults to http://127.0.0.1:18080
  BILLING_E2E_API_LOG            API log path, defaults to /tmp/billing-full-e2e-api.log
  BILLING_E2E_TIMEOUT            smoke timeout, defaults to 90s
  BILLING_E2E_SKIP_NPM_CI        set to 1 to skip npm ci when node_modules is already trusted
  GO                             Go binary, defaults to go

This gate is local/dev only. Do not run it with production services, production
DSNs, real provider credentials, or unmasked customer data.
USAGE
}

log() {
	printf '\n==> %s\n' "$*" >&2
}

die() {
	printf 'error: %s\n' "$*" >&2
	exit 1
}

require_command() {
	command -v "$1" >/dev/null 2>&1 || die "$1 is required on PATH"
}

lowercase() {
	printf '%s' "$1" | tr '[:upper:]' '[:lower:]'
}

ensure_repo_root() {
	[[ -f go.mod && -d frontend && -d cmd/smoke ]] || die "run this script from the repository root"
}

reject_production_runtime() {
	local app_env
	app_env="$(lowercase "${APP_ENV:-local}")"
	case "$app_env" in
	local | dev)
		return 0
		;;
	prod | production | staging)
		die "full E2E smoke uses dev actor headers and must run with APP_ENV=local or APP_ENV=dev"
		;;
	*)
		die "APP_ENV must be local or dev for this gate"
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
		die "$label contains a production-like marker"
		;;
	esac
}

run_step() {
	log "$*"
	"$@"
}

run_without_db_dsn() {
	log "$*"
	env -u DB_DSN "$@"
}

run_dev_db_smoke() {
	local timeout="$1"
	log "${GO:-go} run ./cmd/smoke -dsn <redacted> -timeout ${timeout} dev-db"
	"${GO:-go}" run ./cmd/smoke -dsn "$DB_DSN" -timeout "$timeout" dev-db
}

run_api_smoke() {
	local api_base="$1"
	local timeout="$2"
	log "${GO:-go} run ./cmd/smoke -dsn <redacted> -base-url ${api_base} -timeout ${timeout} dev-api"
	"${GO:-go}" run ./cmd/smoke -dsn "$DB_DSN" -base-url "$api_base" -timeout "$timeout" dev-api
}

run_billing_smoke() {
	local api_base="$1"
	local timeout="$2"
	log "${GO:-go} run ./cmd/smoke -dsn <redacted> -base-url ${api_base} -timeout ${timeout} dev-billing"
	"${GO:-go}" run ./cmd/smoke -dsn "$DB_DSN" -base-url "$api_base" -timeout "$timeout" dev-billing
}

wait_for_url() {
	local url="$1"
	local attempts="${2:-30}"
	local delay="${3:-1}"
	local i
	for ((i = 1; i <= attempts; i++)); do
		if curl -fsS "$url" >/dev/null 2>&1; then
			return 0
		fi
		sleep "$delay"
	done
	return 1
}

start_api() {
	local api_addr="$1"
	local api_log="$2"

	: >"$api_log"
	log "starting local API at ${api_addr}"
	APP_ENV="${APP_ENV:-local}" \
	APP_HTTP_ADDR="$api_addr" \
	DB_DSN="$DB_DSN" \
	LOG_LEVEL="${LOG_LEVEL:-info}" \
	AUTH_SESSION_COOKIE_SECURE="${AUTH_SESSION_COOKIE_SECURE:-false}" \
	"${GO:-go}" run ./cmd/api >"$api_log" 2>&1 &
	api_pid="$!"
}

cleanup() {
	if [[ -n "${api_pid:-}" ]]; then
		kill "$api_pid" >/dev/null 2>&1 || true
		wait "$api_pid" >/dev/null 2>&1 || true
	fi
}

main() {
	if [[ "${1:-}" == "--help" || "${1:-}" == "-h" ]]; then
		usage
		return 0
	fi
	[[ $# -eq 0 ]] || die "unknown argument: $1"

	ensure_repo_root
	reject_production_runtime
	require_command "${GO:-go}"
	require_command make
	require_command npm
	require_command curl
	require_command git

	[[ -n "${DB_DSN:-}" ]] || die "DB_DSN is required"
	reject_production_marker "DB_DSN" "$DB_DSN"

	local api_addr="${BILLING_E2E_API_ADDR:-127.0.0.1:18080}"
	local api_base="${API_BASE_URL:-http://${api_addr}}"
	local api_log="${BILLING_E2E_API_LOG:-/tmp/billing-full-e2e-api.log}"
	local timeout="${BILLING_E2E_TIMEOUT:-90s}"
	reject_production_marker "API_BASE_URL" "$api_base"

	trap cleanup EXIT

	run_without_db_dsn make task-guard
	run_without_db_dsn make test
	run_without_db_dsn make contract-guard
	run_without_db_dsn make error-code-guard
	run_without_db_dsn make build
	run_step git diff --check
	run_dev_db_smoke "$timeout"

	start_api "$api_addr" "$api_log"
	if ! wait_for_url "${api_base}/healthz" 45 1; then
		tail -n 80 "$api_log" >&2 || true
		die "API health check did not pass"
	fi
	if ! wait_for_url "${api_base}/readyz" 45 1; then
		tail -n 80 "$api_log" >&2 || true
		die "API readiness check did not pass"
	fi

	run_api_smoke "$api_base" "$timeout"
	run_billing_smoke "$api_base" "$timeout"

	if [[ "${BILLING_E2E_SKIP_NPM_CI:-}" != "1" ]]; then
		run_step npm --prefix frontend ci
	fi
	run_step npm --prefix frontend audit --omit=dev
	run_step npm --prefix frontend run check:sensitive-text
	run_step npm --prefix frontend run lint
	run_step npm --prefix frontend run build
	run_step npm --prefix frontend run smoke:admin:ci

	log "full E2E quality gate passed"
}

api_pid=""
main "$@"
