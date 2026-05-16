#!/usr/bin/env bash
set -euo pipefail

die() {
  printf 'cloudmini pilot mapping: %s\n' "$*" >&2
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

require_uuid() {
  local name="$1"
  local value="$2"
  [[ "$value" =~ ^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$ ]] || die "$name must be a UUID"
}

require_positive_integer() {
  local name="$1"
  local value="$2"
  [[ "$value" =~ ^[1-9][0-9]*$ ]] || die "$name must be a positive integer"
}

require_one() {
  local name="$1"
  local value="$2"
  [[ "$value" == "1" ]] || die "$name must be 1 for the first controlled pilot"
}

require_command psql

app_env="${APP_ENV:-}"
require_non_empty APP_ENV "$app_env"
case "$app_env" in
  local|dev|staging|sandbox) ;;
  prod|production) die "refusing to run against production" ;;
  *) die "APP_ENV must be local, dev, staging, or sandbox" ;;
esac

require_non_empty DB_DSN "${DB_DSN:-}"
[[ "${BILLING_CLOUDMINI_PILOT_APPROVED:-}" == "yes" ]] || die "BILLING_CLOUDMINI_PILOT_APPROVED=yes is required"

plan_code="${CLOUDMINI_V3_PLAN_CODE:-proxy-static-10gb-monthly}"
source_id="${CLOUDMINI_V3_SOURCE_ID:-00000000-0000-0000-0000-000000000304}"
plan_source_id="${CLOUDMINI_V3_PLAN_SOURCE_ID:-00000000-0000-0000-0000-000000000607}"
source_name="${CLOUDMINI_V3_SOURCE_NAME:-Cloudmini V3 Dev Pilot}"
source_location="${CLOUDMINI_V3_LOCATION:-cloudmini-dev}"
provider_kind="${CLOUDMINI_V3_KIND:-ipv4_dc}"
provider_protocol="${CLOUDMINI_V3_PROTOCOL:-socks5}"
provider_group_ref="${CLOUDMINI_V3_GROUP_REF:-redacted:c6a7189f0a}"
max_create="${CLOUDMINI_V3_PILOT_MAX_CREATE:-1}"
max_active="${CLOUDMINI_V3_PILOT_MAX_ACTIVE_RESOURCES:-1}"
worker_concurrency="${CLOUDMINI_V3_WORKER_CONCURRENCY:-1}"
provider_rate_limit="${CLOUDMINI_V3_PROVIDER_RATE_LIMIT:-no-parallel-mutating-calls}"
spend_exposure="${CLOUDMINI_V3_MAX_SPEND_EXPOSURE:-single-dev-resource}"

require_non_empty CLOUDMINI_V3_PLAN_CODE "$plan_code"
require_uuid CLOUDMINI_V3_SOURCE_ID "$source_id"
require_uuid CLOUDMINI_V3_PLAN_SOURCE_ID "$plan_source_id"
require_non_empty CLOUDMINI_V3_SOURCE_NAME "$source_name"
require_non_empty CLOUDMINI_V3_LOCATION "$source_location"
[[ "$provider_kind" == "ipv4_dc" || "$provider_kind" == "residential" ]] || die "CLOUDMINI_V3_KIND must be ipv4_dc or residential"
[[ "$provider_protocol" == "http" || "$provider_protocol" == "socks5" ]] || die "CLOUDMINI_V3_PROTOCOL must be http or socks5"
[[ "$provider_group_ref" == redacted:* ]] || die "CLOUDMINI_V3_GROUP_REF must be a redacted reference, not a raw provider group id"
require_positive_integer CLOUDMINI_V3_PILOT_MAX_CREATE "$max_create"
require_positive_integer CLOUDMINI_V3_PILOT_MAX_ACTIVE_RESOURCES "$max_active"
require_positive_integer CLOUDMINI_V3_WORKER_CONCURRENCY "$worker_concurrency"
require_one CLOUDMINI_V3_PILOT_MAX_CREATE "$max_create"
require_one CLOUDMINI_V3_PILOT_MAX_ACTIVE_RESOURCES "$max_active"
require_one CLOUDMINI_V3_WORKER_CONCURRENCY "$worker_concurrency"

enum_exists="$(
  psql "$DB_DSN" -X -qAt -v ON_ERROR_STOP=1 <<'SQL'
SELECT EXISTS (
  SELECT 1
  FROM pg_enum e
  JOIN pg_type t ON t.oid = e.enumtypid
  WHERE t.typname = 'catalog_provider_type'
    AND e.enumlabel = 'cloudmini_v3'
);
SQL
)"
[[ "$enum_exists" == "t" ]] || die "catalog_provider_type lacks cloudmini_v3; run migrations first"

plan_exists="$(
  psql "$DB_DSN" -X -qAt -v ON_ERROR_STOP=1 -v plan_code="$plan_code" <<'SQL'
SELECT EXISTS (
  SELECT 1
  FROM master_plans
  WHERE plan_code = :'plan_code'
);
SQL
)"
[[ "$plan_exists" == "t" ]] || die "target plan code was not found"

psql "$DB_DSN" -X -v ON_ERROR_STOP=1 \
  -v plan_code="$plan_code" \
  -v source_id="$source_id" \
  -v plan_source_id="$plan_source_id" \
  -v source_name="$source_name" \
  -v source_location="$source_location" \
  -v provider_kind="$provider_kind" \
  -v provider_protocol="$provider_protocol" \
  -v provider_group_ref="$provider_group_ref" \
  -v max_create="$max_create" \
  -v max_active="$max_active" \
  -v worker_concurrency="$worker_concurrency" \
  -v provider_rate_limit="$provider_rate_limit" \
  -v spend_exposure="$spend_exposure" <<'SQL'
BEGIN;

INSERT INTO provider_sources (
  source_id,
  source_type,
  name,
  provider_account_id,
  location,
  status,
  capability_profile,
  inventory_mode,
  risk_level
)
VALUES (
  :'source_id'::uuid,
  'cloudmini_v3',
  :'source_name',
  NULL,
  :'source_location',
  'active',
  jsonb_build_object(
    'supportsHealthCheck', true,
    'supportsLiveStockCheck', true,
    'supportsAutoProvision', true,
    'supportsStatusSync', true,
    'supportsSuspend', true,
    'supportsUnsuspend', true,
    'supportsTerminate', true,
    'supportsRenew', false,
    'supportsResetPassword', false,
    'supportsCredentialFetch', true,
    'proxy', jsonb_build_object(
      'supportsHTTPProtocol', true,
      'supportsSOCKS5Protocol', true,
      'supportsRotatingProxy', false,
      'supportsStaticProxy', true,
      'supportsUserPassAuth', true,
      'supportsBandwidthQuota', true
    )
  ),
  'provider_live',
  'high'
)
ON CONFLICT (source_id) DO UPDATE
SET source_type = EXCLUDED.source_type,
    name = EXCLUDED.name,
    provider_account_id = EXCLUDED.provider_account_id,
    location = EXCLUDED.location,
    status = EXCLUDED.status,
    capability_profile = EXCLUDED.capability_profile,
    inventory_mode = EXCLUDED.inventory_mode,
    risk_level = EXCLUDED.risk_level,
    updated_at = NOW();

WITH target_plan AS (
  SELECT plan_id
  FROM master_plans
  WHERE plan_code = :'plan_code'
)
INSERT INTO plan_sources (
  plan_source_id,
  plan_id,
  source_id,
  priority,
  cost_override_minor,
  capacity_policy,
  capability_override,
  status
)
SELECT
  :'plan_source_id'::uuid,
  target_plan.plan_id,
  :'source_id'::uuid,
  1,
  0,
  jsonb_build_object(
    'pilotMaxCreate', :'max_create'::int,
    'pilotMaxActiveResources', :'max_active'::int,
    'workerConcurrency', :'worker_concurrency'::int,
    'providerRateLimit', :'provider_rate_limit',
    'maximumSpendOrQuotaExposure', :'spend_exposure',
    'providerKind', :'provider_kind',
    'providerGroupRef', :'provider_group_ref',
    'protocol', :'provider_protocol'
  ),
  '{}'::jsonb,
  'active'
FROM target_plan
ON CONFLICT (plan_id, source_id) DO UPDATE
SET priority = EXCLUDED.priority,
    cost_override_minor = EXCLUDED.cost_override_minor,
    capacity_policy = EXCLUDED.capacity_policy,
    capability_override = EXCLUDED.capability_override,
    status = EXCLUDED.status,
    updated_at = NOW();

UPDATE plan_sources existing
SET priority = 5,
    updated_at = NOW()
FROM master_plans plan
WHERE existing.plan_id = plan.plan_id
  AND plan.plan_code = :'plan_code'
  AND existing.source_id = '00000000-0000-0000-0000-000000000302'::uuid
  AND existing.source_id <> :'source_id'::uuid
  AND existing.priority <= 1;

COMMIT;

SELECT
  'plan_code=' || plan.plan_code,
  'plan_source_display_id=' || plan_source.display_id,
  'source_display_id=' || source.display_id,
  'source_type=' || source.source_type,
  'source_status=' || source.status,
  'inventory_mode=' || source.inventory_mode,
  'plan_source_status=' || plan_source.status,
  'priority=' || plan_source.priority
FROM master_plans plan
JOIN plan_sources plan_source ON plan_source.plan_id = plan.plan_id
JOIN provider_sources source ON source.source_id = plan_source.source_id
WHERE plan.plan_code = :'plan_code'
  AND source.source_id = :'source_id'::uuid;
SQL

printf 'Cloudmini pilot mapping applied. Recheck provider readiness via the admin readiness API before running any mutating pilot.\n'
