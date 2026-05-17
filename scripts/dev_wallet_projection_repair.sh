#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat <<'USAGE'
Usage:
  scripts/dev_wallet_projection_repair.sh [--apply]

Repairs a non-production dev/test wallet projection so wallet.available_balance_minor
matches the posted ledger sum. This is not a production finance adjustment path.

Required env:
  APP_ENV=dev|test|local
  DB_DSN=<postgres dsn>
  BILLING_DEV_WALLET_PROJECTION_REPAIR_APPROVED=yes

Optional env:
  BILLING_DEV_REPAIR_TENANT_SLUG=demo-reseller
  BILLING_DEV_REPAIR_WALLET_DISPLAY_ID=41001
USAGE
}

apply=no
if [ "${1:-}" = "--apply" ]; then
  apply=yes
elif [ "${1:-}" = "-h" ] || [ "${1:-}" = "--help" ]; then
  usage
  exit 0
elif [ "${1:-}" != "" ]; then
  usage >&2
  exit 2
fi

app_env="${APP_ENV:-}"
case "$app_env" in
  local|dev|test|development)
    ;;
  "")
    echo "dev wallet projection repair: APP_ENV is required" >&2
    exit 1
    ;;
  *)
    echo "dev wallet projection repair: refusing APP_ENV=$app_env" >&2
    exit 1
    ;;
esac

if [ -z "${DB_DSN:-}" ]; then
  echo "dev wallet projection repair: DB_DSN is required" >&2
  exit 1
fi

if [ "${BILLING_DEV_WALLET_PROJECTION_REPAIR_APPROVED:-}" != "yes" ]; then
  echo "dev wallet projection repair: BILLING_DEV_WALLET_PROJECTION_REPAIR_APPROVED=yes is required" >&2
  exit 1
fi

tenant_slug="${BILLING_DEV_REPAIR_TENANT_SLUG:-demo-reseller}"
wallet_display_id="${BILLING_DEV_REPAIR_WALLET_DISPLAY_ID:-41001}"

case "$wallet_display_id" in
  ''|*[!0-9]*)
    echo "dev wallet projection repair: BILLING_DEV_REPAIR_WALLET_DISPLAY_ID must be a positive integer" >&2
    exit 1
    ;;
esac

if [ "$apply" != "yes" ]; then
  psql "$DB_DSN" -X -q -v ON_ERROR_STOP=1 -v tenant_slug="$tenant_slug" -v wallet_display_id="$wallet_display_id" -P pager=off <<'SQL'
WITH ledger_balance AS (
  SELECT wallet_id,
         COALESCE(SUM(CASE WHEN direction = 'credit' THEN amount_minor ELSE -amount_minor END), 0) AS ledger_balance_minor,
         COUNT(*) AS posted_ledger_entries,
         MAX(display_id) AS last_ledger_display_id
  FROM wallet_ledger_entries
  WHERE status = 'posted'
  GROUP BY wallet_id
)
SELECT tenant.slug AS tenant_slug,
       wallet.display_id AS wallet_display_id,
       wallet.currency,
       wallet.available_balance_minor,
       COALESCE(ledger_balance.ledger_balance_minor, 0) AS ledger_balance_minor,
       wallet.available_balance_minor - COALESCE(ledger_balance.ledger_balance_minor, 0) AS difference_minor,
       COALESCE(ledger_balance.posted_ledger_entries, 0) AS posted_ledger_entries,
       ledger_balance.last_ledger_display_id
FROM wallets wallet
JOIN tenants tenant ON tenant.tenant_id = wallet.tenant_id
LEFT JOIN ledger_balance ON ledger_balance.wallet_id = wallet.wallet_id
WHERE tenant.slug = :'tenant_slug'
  AND wallet.display_id = :'wallet_display_id'::bigint;
SQL
  echo "result=PLAN"
  echo "apply_required=yes"
  echo "mutating_routes_called=no"
  echo "secrets_printed=no"
  exit 0
fi

psql "$DB_DSN" -X -q -v ON_ERROR_STOP=1 -v tenant_slug="$tenant_slug" -v wallet_display_id="$wallet_display_id" -P pager=off <<'SQL'
WITH ledger_balance AS (
  SELECT wallet_id,
         COALESCE(SUM(CASE WHEN direction = 'credit' THEN amount_minor ELSE -amount_minor END), 0) AS ledger_balance_minor,
         COUNT(*) AS posted_ledger_entries,
         MAX(display_id) AS last_ledger_display_id
  FROM wallet_ledger_entries
  WHERE status = 'posted'
  GROUP BY wallet_id
),
mismatch AS (
  SELECT wallet.wallet_id,
         wallet.tenant_id,
         wallet.display_id AS wallet_display_id,
         wallet.currency,
         wallet.available_balance_minor AS before_available_balance_minor,
         COALESCE(ledger_balance.ledger_balance_minor, 0) AS ledger_balance_minor,
         wallet.available_balance_minor - COALESCE(ledger_balance.ledger_balance_minor, 0) AS difference_minor,
         COALESCE(ledger_balance.posted_ledger_entries, 0) AS posted_ledger_entries,
         ledger_balance.last_ledger_display_id
  FROM wallets wallet
  JOIN tenants tenant ON tenant.tenant_id = wallet.tenant_id
  LEFT JOIN ledger_balance ON ledger_balance.wallet_id = wallet.wallet_id
  WHERE tenant.slug = :'tenant_slug'
    AND wallet.display_id = :'wallet_display_id'::bigint
    AND wallet.available_balance_minor <> COALESCE(ledger_balance.ledger_balance_minor, 0)
),
repaired AS (
  UPDATE wallets wallet
  SET available_balance_minor = mismatch.ledger_balance_minor,
      updated_at = NOW()
  FROM mismatch
  WHERE wallet.wallet_id = mismatch.wallet_id
    AND wallet.tenant_id = mismatch.tenant_id
  RETURNING wallet.wallet_id,
            wallet.tenant_id,
            mismatch.wallet_display_id,
            mismatch.currency,
            mismatch.before_available_balance_minor,
            wallet.available_balance_minor AS after_available_balance_minor,
            mismatch.ledger_balance_minor,
            mismatch.difference_minor,
            mismatch.posted_ledger_entries,
            mismatch.last_ledger_display_id
),
audit_insert AS (
  INSERT INTO audit_logs (
    tenant_id,
    actor_id,
    actor_type,
    action,
    target_type,
    target_id,
    before_snapshot_redacted,
    after_snapshot_redacted,
    metadata_redacted,
    correlation_id
  )
  SELECT repaired.tenant_id,
         NULL,
         'system',
         'wallet.projection.repaired',
         'wallet',
         repaired.wallet_id,
         jsonb_build_object('available_balance_minor', repaired.before_available_balance_minor),
         jsonb_build_object('available_balance_minor', repaired.after_available_balance_minor),
         jsonb_build_object(
           'wallet_display_id', repaired.wallet_display_id,
           'currency', repaired.currency,
           'ledger_balance_minor', repaired.ledger_balance_minor,
           'difference_minor', repaired.difference_minor,
           'posted_ledger_entries', repaired.posted_ledger_entries,
           'last_ledger_display_id', repaired.last_ledger_display_id,
           'task', 'T239',
           'reason', 'dev/test projection repair from posted ledger source of truth'
         ),
         gen_random_uuid()
  FROM repaired
  RETURNING display_id AS audit_display_id
)
SELECT repaired.wallet_display_id,
       repaired.currency,
       repaired.before_available_balance_minor,
       repaired.after_available_balance_minor,
       repaired.ledger_balance_minor,
       repaired.difference_minor,
       repaired.posted_ledger_entries,
       repaired.last_ledger_display_id,
       audit_insert.audit_display_id
FROM repaired
CROSS JOIN audit_insert;
SQL

echo "result=APPLIED"
echo "mutating_routes_called=no"
echo "ledger_rows_inserted=0"
echo "posted_ledger_rows_updated=0"
echo "secrets_printed=no"
