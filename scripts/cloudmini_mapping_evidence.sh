#!/usr/bin/env bash
set -euo pipefail

die() {
  printf 'cloudmini mapping evidence: %s\n' "$*" >&2
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

require_positive_integer_if_set() {
  local name="$1"
  local value="$2"
  [[ -z "$value" || "$value" =~ ^[1-9][0-9]*$ ]] || die "$name must be a positive integer when set"
}

app_env="${APP_ENV:-}"
require_non_empty APP_ENV "$app_env"
case "$app_env" in
  local|dev|staging|sandbox) ;;
  prod|production) die "refusing to run against production" ;;
  *) die "APP_ENV must be local, dev, staging, or sandbox" ;;
esac

require_non_empty DB_DSN "${DB_DSN:-}"
[[ "${BILLING_CLOUDMINI_EVIDENCE_APPROVED:-}" == "yes" ]] || die "BILLING_CLOUDMINI_EVIDENCE_APPROVED=yes is required"

require_command psql

plan_code="${CLOUDMINI_V3_PLAN_CODE:-proxy-static-10gb-monthly}"
source_display_id="${CLOUDMINI_V3_SOURCE_DISPLAY_ID:-}"
plan_source_display_id="${CLOUDMINI_V3_PLAN_SOURCE_DISPLAY_ID:-}"

require_non_empty CLOUDMINI_V3_PLAN_CODE "$plan_code"
require_positive_integer_if_set CLOUDMINI_V3_SOURCE_DISPLAY_ID "$source_display_id"
require_positive_integer_if_set CLOUDMINI_V3_PLAN_SOURCE_DISPLAY_ID "$plan_source_display_id"

evidence="$(
  psql "$DB_DSN" -X -qAt -v ON_ERROR_STOP=1 \
    -v plan_code="$plan_code" \
    -v source_display_id="$source_display_id" \
    -v plan_source_display_id="$plan_source_display_id" <<'SQL'
BEGIN READ ONLY;

WITH requested AS (
  SELECT
    :'plan_code'::text AS plan_code,
    NULLIF(:'source_display_id', '')::bigint AS expected_source_display_id,
    NULLIF(:'plan_source_display_id', '')::bigint AS expected_plan_source_display_id
),
enum_state AS (
  SELECT EXISTS (
    SELECT 1
    FROM pg_enum enum_value
    JOIN pg_type enum_type ON enum_type.oid = enum_value.enumtypid
    WHERE enum_type.typname = 'catalog_provider_type'
      AND enum_value.enumlabel = 'cloudmini_v3'
  ) AS has_cloudmini_v3
),
target_plan AS (
  SELECT
    plan.plan_id,
    plan.display_id,
    plan.plan_code,
    plan.status::text AS plan_status,
    product.product_type::text AS product_type
  FROM master_plans plan
  JOIN master_products product ON product.product_id = plan.product_id
  JOIN requested ON requested.plan_code = plan.plan_code
),
ranked_plan_sources AS (
  SELECT
    plan_source.plan_id,
    plan_source.display_id AS plan_source_display_id,
    plan_source.status::text AS plan_source_status,
    plan_source.priority,
    plan_source.capacity_policy,
    source.display_id AS source_display_id,
    source.source_type::text AS source_type,
    source.name AS source_name,
    source.status::text AS source_status,
    source.inventory_mode::text AS inventory_mode,
    source.capability_profile,
    ROW_NUMBER() OVER (
      PARTITION BY plan_source.plan_id
      ORDER BY CASE WHEN plan_source.status = 'active' AND source.status = 'active' THEN 0 ELSE 1 END,
        plan_source.priority ASC,
        plan_source.created_at ASC
    ) AS row_number
  FROM plan_sources plan_source
  JOIN provider_sources source ON source.source_id = plan_source.source_id
  JOIN target_plan ON target_plan.plan_id = plan_source.plan_id
),
selected_source AS (
  SELECT *
  FROM ranked_plan_sources
  WHERE row_number = 1
),
evidence AS (
  SELECT
    requested.plan_code AS requested_plan_code,
    enum_state.has_cloudmini_v3,
    target_plan.display_id AS plan_display_id,
    target_plan.plan_code,
    target_plan.plan_status,
    target_plan.product_type,
    selected_source.plan_source_display_id,
    selected_source.plan_source_status,
    selected_source.priority,
    selected_source.source_display_id,
    selected_source.source_type,
    selected_source.source_name,
    selected_source.source_status,
    selected_source.inventory_mode,
    CASE
      WHEN selected_source.capacity_policy->>'pilotMaxCreate' ~ '^[0-9]+$'
        THEN (selected_source.capacity_policy->>'pilotMaxCreate')::int
    END AS pilot_max_create,
    CASE
      WHEN selected_source.capacity_policy->>'pilotMaxActiveResources' ~ '^[0-9]+$'
        THEN (selected_source.capacity_policy->>'pilotMaxActiveResources')::int
    END AS pilot_max_active_resources,
    CASE
      WHEN selected_source.capacity_policy->>'workerConcurrency' ~ '^[0-9]+$'
        THEN (selected_source.capacity_policy->>'workerConcurrency')::int
    END AS worker_concurrency,
    selected_source.capacity_policy->>'providerRateLimit' AS provider_rate_limit,
    selected_source.capacity_policy->>'maximumSpendOrQuotaExposure' AS maximum_spend_or_quota_exposure,
    selected_source.capacity_policy->>'providerKind' AS provider_kind,
    selected_source.capacity_policy->>'providerGroupRef' AS provider_group_ref,
    selected_source.capacity_policy->>'protocol' AS protocol,
    requested.expected_source_display_id,
    requested.expected_plan_source_display_id,
    CASE
      WHEN selected_source.plan_source_display_id IS NULL OR selected_source.source_display_id IS NULL
        THEN 'missing_plan_source'
      WHEN selected_source.plan_source_status <> 'active' OR selected_source.source_status <> 'active'
        THEN 'inactive_source'
      WHEN selected_source.source_type = 'manual'
        THEN 'fake_provider_only'
      WHEN NOT (
        selected_source.capability_profile->>'supportsAutoProvision' = 'true'
        AND CASE target_plan.product_type
          WHEN 'proxy' THEN (
            selected_source.capability_profile->'proxy'->>'supportsHTTPProtocol' = 'true'
            OR selected_source.capability_profile->'proxy'->>'supportsSOCKS5Protocol' = 'true'
            OR selected_source.capability_profile->'proxy'->>'supportsRotatingProxy' = 'true'
            OR selected_source.capability_profile->'proxy'->>'supportsStaticProxy' = 'true'
          )
          WHEN 'vps' THEN (
            selected_source.capability_profile->'vps'->>'supportsOSTemplateSelection' = 'true'
            OR selected_source.capability_profile->'vps'->>'supportsCustomHostname' = 'true'
            OR selected_source.capability_profile->'vps'->>'supportsIPv6' = 'true'
            OR selected_source.capability_profile->'vps'->>'supportsResize' = 'true'
            OR selected_source.capability_profile->'vps'->>'supportsVNCConsole' = 'true'
          )
          WHEN 'service_addon' THEN true
          ELSE false
        END
      )
        THEN 'unsupported_capability'
      ELSE 'ready'
    END AS readiness_state
  FROM requested
  CROSS JOIN enum_state
  LEFT JOIN target_plan ON true
  LEFT JOIN selected_source ON true
),
checks AS (
  SELECT 'enum_has_cloudmini_v3' AS check_name, has_cloudmini_v3 AS passed FROM evidence
  UNION ALL SELECT 'plan_exists', plan_display_id IS NOT NULL FROM evidence
  UNION ALL SELECT 'plan_active', plan_status = 'active' FROM evidence
  UNION ALL SELECT 'selected_source_exists', source_display_id IS NOT NULL FROM evidence
  UNION ALL SELECT 'expected_source_display_id', expected_source_display_id IS NULL OR source_display_id = expected_source_display_id FROM evidence
  UNION ALL SELECT 'expected_plan_source_display_id', expected_plan_source_display_id IS NULL OR plan_source_display_id = expected_plan_source_display_id FROM evidence
  UNION ALL SELECT 'source_type_cloudmini_v3', source_type = 'cloudmini_v3' FROM evidence
  UNION ALL SELECT 'readiness_ready', readiness_state = 'ready' FROM evidence
  UNION ALL SELECT 'plan_source_priority_one', priority = 1 FROM evidence
  UNION ALL SELECT 'pilot_max_create_one', pilot_max_create = 1 FROM evidence
  UNION ALL SELECT 'pilot_max_active_resources_one', pilot_max_active_resources = 1 FROM evidence
  UNION ALL SELECT 'worker_concurrency_one', worker_concurrency = 1 FROM evidence
  UNION ALL SELECT 'provider_group_ref_redacted', provider_group_ref LIKE 'redacted:%' FROM evidence
  UNION ALL SELECT 'provider_kind_allowed', provider_kind IN ('ipv4_dc', 'residential') FROM evidence
  UNION ALL SELECT 'protocol_allowed', protocol IN ('http', 'socks5') FROM evidence
),
summary AS (
  SELECT
    evidence.*,
    COALESCE(NULLIF(string_agg(checks.check_name, ',' ORDER BY checks.check_name) FILTER (WHERE checks.passed IS NOT TRUE), ''), 'none') AS failed_checks
  FROM evidence
  CROSS JOIN checks
  GROUP BY
    evidence.requested_plan_code,
    evidence.has_cloudmini_v3,
    evidence.plan_display_id,
    evidence.plan_code,
    evidence.plan_status,
    evidence.product_type,
    evidence.plan_source_display_id,
    evidence.plan_source_status,
    evidence.priority,
    evidence.source_display_id,
    evidence.source_type,
    evidence.source_name,
    evidence.source_status,
    evidence.inventory_mode,
    evidence.pilot_max_create,
    evidence.pilot_max_active_resources,
    evidence.worker_concurrency,
    evidence.provider_rate_limit,
    evidence.maximum_spend_or_quota_exposure,
    evidence.provider_kind,
    evidence.provider_group_ref,
    evidence.protocol,
    evidence.expected_source_display_id,
    evidence.expected_plan_source_display_id,
    evidence.readiness_state
)
SELECT output.line
FROM summary
CROSS JOIN LATERAL (
  VALUES
    (1, 'result=' || CASE WHEN failed_checks = 'none' THEN 'PASS' ELSE 'FAIL' END),
    (2, 'plan_code=' || requested_plan_code),
    (3, 'plan_display_id=' || COALESCE(plan_display_id::text, 'not_shown')),
    (4, 'product_type=' || COALESCE(product_type, 'not_shown')),
    (5, 'readiness_state=' || readiness_state),
    (6, 'plan_source_display_id=' || COALESCE(plan_source_display_id::text, 'not_shown')),
    (7, 'source_display_id=' || COALESCE(source_display_id::text, 'not_shown')),
    (8, 'source_type=' || COALESCE(source_type, 'not_shown')),
    (9, 'source_status=' || COALESCE(source_status, 'not_shown')),
    (10, 'inventory_mode=' || COALESCE(inventory_mode, 'not_shown')),
    (11, 'priority=' || COALESCE(priority::text, 'not_shown')),
    (12, 'pilot_max_create=' || COALESCE(pilot_max_create::text, 'not_shown')),
    (13, 'pilot_max_active_resources=' || COALESCE(pilot_max_active_resources::text, 'not_shown')),
    (14, 'worker_concurrency=' || COALESCE(worker_concurrency::text, 'not_shown')),
    (15, 'provider_rate_limit=' || COALESCE(provider_rate_limit, 'not_shown')),
    (16, 'maximum_spend_or_quota_exposure=' || COALESCE(maximum_spend_or_quota_exposure, 'not_shown')),
    (17, 'provider_kind=' || COALESCE(provider_kind, 'not_shown')),
    (18, 'provider_group_ref=' || COALESCE(provider_group_ref, 'not_shown')),
    (19, 'protocol=' || COALESCE(protocol, 'not_shown')),
    (20, 'failed_checks=' || failed_checks)
) AS output(sort_order, line)
ORDER BY output.sort_order;

COMMIT;
SQL
)"

printf '%s\n' "$evidence"
printf '%s\n' "$evidence" | grep -qx 'result=PASS' || die "evidence checks failed"

printf 'Cloudmini mapping evidence passed. Output intentionally excludes DB_DSN, tokens, raw provider group IDs, raw provider payloads, and proxy credentials.\n'
