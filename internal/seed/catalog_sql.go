package seed

const seedProviderSourcesSQL = `
INSERT INTO provider_sources (source_id, source_type, name, provider_account_id, location, status, capability_profile, inventory_mode, risk_level)
VALUES
    (
        '00000000-0000-0000-0000-000000000301',
        'manual',
        'Manual Local Pool',
        NULL,
        'local',
        'active',
        '{"supportsHealthCheck":true,"supportsManualProvision":true,"supportsStatusSync":true}'::jsonb,
        'manual_unlimited',
        'low'
    )
ON CONFLICT (source_id) DO UPDATE
SET source_type = EXCLUDED.source_type,
    name = EXCLUDED.name,
    location = EXCLUDED.location,
    status = EXCLUDED.status,
    capability_profile = EXCLUDED.capability_profile,
    inventory_mode = EXCLUDED.inventory_mode,
    risk_level = EXCLUDED.risk_level;
`

const seedMasterProductsSQL = `
INSERT INTO master_products (product_id, product_type, name, description, status, display_order, created_by)
SELECT
    '00000000-0000-0000-0000-000000000401',
    'vps',
    'VPS',
    'Virtual private server plans for local development.',
    'active',
    10,
    admin.user_id
FROM users admin
JOIN tenants platform ON platform.tenant_id = admin.tenant_id
WHERE platform.slug = 'platform'
  AND admin.email = 'admin@local.billing'
ON CONFLICT (product_id) DO UPDATE
SET product_type = EXCLUDED.product_type,
    name = EXCLUDED.name,
    description = EXCLUDED.description,
    status = EXCLUDED.status,
    display_order = EXCLUDED.display_order,
    created_by = EXCLUDED.created_by;

INSERT INTO master_products (product_id, product_type, name, description, status, display_order, created_by)
SELECT
    '00000000-0000-0000-0000-000000000402',
    'proxy',
    'Proxy',
    'Proxy service plans for local development.',
    'active',
    20,
    admin.user_id
FROM users admin
JOIN tenants platform ON platform.tenant_id = admin.tenant_id
WHERE platform.slug = 'platform'
  AND admin.email = 'admin@local.billing'
ON CONFLICT (product_id) DO UPDATE
SET product_type = EXCLUDED.product_type,
    name = EXCLUDED.name,
    description = EXCLUDED.description,
    status = EXCLUDED.status,
    display_order = EXCLUDED.display_order,
    created_by = EXCLUDED.created_by;
`

const seedMasterPlansSQL = `
INSERT INTO master_plans (plan_id, product_id, plan_code, name, specs, billing_cycle_type, billing_cycle_value, base_cost_minor, suggested_price_minor, reseller_min_price_minor, currency, status, version)
VALUES
    (
        '00000000-0000-0000-0000-000000000501',
        '00000000-0000-0000-0000-000000000401',
        'vps-cx23-40gb-monthly',
        'CX23 VPS 40GB',
        '{"cpu":"2 vCPU","memory":"4 GB","disk":"40 GB","region":"local"}'::jsonb,
        'month_30d',
        1,
        700,
        1200,
        900,
        'USD',
        'active',
        1
    ),
    (
        '00000000-0000-0000-0000-000000000502',
        '00000000-0000-0000-0000-000000000401',
        'vps-cx33-80gb-monthly',
        'CX33 VPS 80GB',
        '{"cpu":"4 vCPU","memory":"8 GB","disk":"80 GB","region":"local"}'::jsonb,
        'month_30d',
        1,
        1300,
        2200,
        1700,
        'USD',
        'active',
        1
    ),
    (
        '00000000-0000-0000-0000-000000000503',
        '00000000-0000-0000-0000-000000000402',
        'proxy-static-10gb-monthly',
        'Static Proxy 10GB',
        '{"traffic":"10 GB","protocols":["http","socks5"],"region":"local"}'::jsonb,
        'month_30d',
        1,
        300,
        600,
        450,
        'USD',
        'active',
        1
    )
ON CONFLICT (plan_code, version) DO UPDATE
SET product_id = EXCLUDED.product_id,
    name = EXCLUDED.name,
    specs = EXCLUDED.specs,
    billing_cycle_type = EXCLUDED.billing_cycle_type,
    billing_cycle_value = EXCLUDED.billing_cycle_value,
    base_cost_minor = EXCLUDED.base_cost_minor,
    suggested_price_minor = EXCLUDED.suggested_price_minor,
    reseller_min_price_minor = EXCLUDED.reseller_min_price_minor,
    currency = EXCLUDED.currency,
    status = EXCLUDED.status;
`
