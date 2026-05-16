#!/usr/bin/env bash
set -euo pipefail

die() {
  printf 'cloudmini mutating pilot preflight: %s\n' "$*" >&2
  exit 1
}

require_command() {
  command -v "$1" >/dev/null 2>&1 || die "$1 is required"
}

require_non_empty() {
  local name="$1"
  local value="$2"
  [[ -n "$value" ]] || die "$name is required"
}

require_filled() {
  local name="$1"
  local value="$2"
  require_non_empty "$name" "$value"
  case "${value,,}" in
    -|todo|tbd|unknown|none|null|placeholder) die "$name must be a real approved value, not a placeholder" ;;
  esac
}

require_one() {
  local name="$1"
  local value="$2"
  [[ "$value" == "1" ]] || die "$name must be 1 for the first controlled pilot"
}

require_exact() {
  local name="$1"
  local value="$2"
  local expected="$3"
  [[ "$value" == "$expected" ]] || die "$name must be $expected"
}

credential_path_is_private() {
  local path="$1"
  [[ "$path" == /* ]] || die "CLOUDMINI_PILOT_CREDENTIAL_PATH must be an absolute path"
  [[ -f "$path" ]] || die "CLOUDMINI_PILOT_CREDENTIAL_PATH must point to a file"
  [[ -r "$path" ]] || die "CLOUDMINI_PILOT_CREDENTIAL_PATH must be readable"

  local repo_root
  repo_root="$(pwd -P)"
  local path_dir
  path_dir="$(cd "$(dirname "$path")" && pwd -P)"
  local path_real="$path_dir/$(basename "$path")"
  case "$path_real" in
    "$repo_root"/*) die "CLOUDMINI_PILOT_CREDENTIAL_PATH must not be inside the repository" ;;
  esac

  local mode
  mode="$(stat -c '%a' "$path" 2>/dev/null)" || die "could not stat CLOUDMINI_PILOT_CREDENTIAL_PATH"
  [[ "$mode" =~ ^[0-7]+$ ]] || die "unexpected credential file mode"
  (( (8#$mode & 8#077) == 0 )) || die "CLOUDMINI_PILOT_CREDENTIAL_PATH must not be group/world accessible"
}

require_command bash

app_env="${APP_ENV:-}"
require_non_empty APP_ENV "$app_env"
case "$app_env" in
  local|dev|staging|sandbox) ;;
  prod|production) die "refusing to run against production" ;;
  *) die "APP_ENV must be local, dev, staging, or sandbox" ;;
esac

require_non_empty DB_DSN "${DB_DSN:-}"
[[ "${BILLING_CLOUDMINI_MUTATING_PREFLIGHT_APPROVED:-}" == "yes" ]] || die "BILLING_CLOUDMINI_MUTATING_PREFLIGHT_APPROVED=yes is required"

required_fields=(
  CLOUDMINI_PILOT_ID
  CLOUDMINI_PILOT_ENVIRONMENT
  CLOUDMINI_SOURCE_ACCOUNT_OWNER
  CLOUDMINI_ENGINEERING_OWNER
  CLOUDMINI_OPS_OWNER
  CLOUDMINI_SECURITY_OWNER
  CLOUDMINI_CLEANUP_OWNER
  CLOUDMINI_FINANCE_QUOTA_OWNER
  CLOUDMINI_REVIEWER_SIGNOFF
  CLOUDMINI_PILOT_CLEANUP_DEADLINE
  CLOUDMINI_PILOT_STOP_CONDITION
  CLOUDMINI_PILOT_READONLY_EVIDENCE_REF
  CLOUDMINI_PILOT_CLEANUP_PROCEDURE_REF
  CLOUDMINI_PILOT_CREDENTIAL_PATH
)

for field in "${required_fields[@]}"; do
  require_filled "$field" "${!field:-}"
done

require_exact CLOUDMINI_PILOT_ENVIRONMENT "$CLOUDMINI_PILOT_ENVIRONMENT" "$app_env"
require_one CLOUDMINI_PILOT_MAX_CREATE_CALLS "${CLOUDMINI_PILOT_MAX_CREATE_CALLS:-}"
require_one CLOUDMINI_PILOT_MAX_ACTIVE_RESOURCES "${CLOUDMINI_PILOT_MAX_ACTIVE_RESOURCES:-}"
require_one CLOUDMINI_PILOT_WORKER_CONCURRENCY "${CLOUDMINI_PILOT_WORKER_CONCURRENCY:-}"
require_exact CLOUDMINI_PILOT_PROVIDER_RATE_LIMIT "${CLOUDMINI_PILOT_PROVIDER_RATE_LIMIT:-}" "no-parallel-mutating-calls"
require_exact CLOUDMINI_PILOT_MAX_SPEND_EXPOSURE "${CLOUDMINI_PILOT_MAX_SPEND_EXPOSURE:-}" "single-dev-resource"
credential_path_is_private "$CLOUDMINI_PILOT_CREDENTIAL_PATH"

[[ -f scripts/cloudmini_mapping_evidence.sh ]] || die "scripts/cloudmini_mapping_evidence.sh is required"

mapping_output="$(
  BILLING_CLOUDMINI_EVIDENCE_APPROVED=yes \
    bash scripts/cloudmini_mapping_evidence.sh
)"

printf '%s\n' "$mapping_output"
printf '%s\n' "$mapping_output" | grep -qx 'result=PASS' || die "mapping evidence did not pass"

printf 'preflight_result=PASS\n'
printf 'pilot_environment=%s\n' "$app_env"
printf 'approval_fields_present=yes\n'
printf 'owner_fields_present=yes\n'
printf 'cleanup_fields_present=yes\n'
printf 'credential_path_private=yes\n'
printf 'max_create_calls=1\n'
printf 'max_active_resources=1\n'
printf 'worker_concurrency=1\n'
printf 'provider_rate_limit=no-parallel-mutating-calls\n'
printf 'maximum_spend_or_quota_exposure=single-dev-resource\n'
printf 'mutating_routes_called=no\n'
