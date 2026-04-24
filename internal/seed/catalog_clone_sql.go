package seed

const seedPlanSourcesSQL = `
INSERT INTO plan_sources (plan_source_id, plan_id, source_id, priority, cost_override_minor, capacity_policy, capability_override, status)
VALUES
    ('00000000-0000-0000-0000-000000000604', '00000000-0000-0000-0000-000000000501', '00000000-0000-0000-0000-000000000302', 1, 0, '{}'::jsonb, '{}'::jsonb, 'active'),
    ('00000000-0000-0000-0000-000000000601', '00000000-0000-0000-0000-000000000501', '00000000-0000-0000-0000-000000000301', 2, 0, '{}'::jsonb, '{}'::jsonb, 'active'),
    ('00000000-0000-0000-0000-000000000602', '00000000-0000-0000-0000-000000000502', '00000000-0000-0000-0000-000000000301', 1, 0, '{}'::jsonb, '{}'::jsonb, 'active'),
    ('00000000-0000-0000-0000-000000000605', '00000000-0000-0000-0000-000000000503', '00000000-0000-0000-0000-000000000302', 1, 0, '{}'::jsonb, '{}'::jsonb, 'active'),
    ('00000000-0000-0000-0000-000000000603', '00000000-0000-0000-0000-000000000503', '00000000-0000-0000-0000-000000000301', 2, 0, '{}'::jsonb, '{}'::jsonb, 'active'),
    ('00000000-0000-0000-0000-000000000606', '00000000-0000-0000-0000-000000000504', '00000000-0000-0000-0000-000000000303', 1, 0, '{}'::jsonb, '{}'::jsonb, 'active')
ON CONFLICT (plan_id, source_id) DO UPDATE
SET priority = EXCLUDED.priority,
    cost_override_minor = EXCLUDED.cost_override_minor,
    capacity_policy = EXCLUDED.capacity_policy,
    capability_override = EXCLUDED.capability_override,
    status = EXCLUDED.status;
`

const seedDemoResellerCatalogSQL = `
INSERT INTO tenant_products (tenant_product_id, tenant_id, master_product_id, name_override, description_override, status, clone_version)
SELECT
    '00000000-0000-0000-0000-000000000701',
    reseller.tenant_id,
    '00000000-0000-0000-0000-000000000401',
    NULL,
    NULL,
    'active',
    1
FROM tenants reseller
WHERE reseller.slug = 'demo-reseller'
ON CONFLICT (tenant_id, master_product_id) DO UPDATE
SET name_override = EXCLUDED.name_override,
    description_override = EXCLUDED.description_override,
    status = EXCLUDED.status,
    clone_version = EXCLUDED.clone_version;

INSERT INTO tenant_products (tenant_product_id, tenant_id, master_product_id, name_override, description_override, status, clone_version)
SELECT
    '00000000-0000-0000-0000-000000000702',
    reseller.tenant_id,
    '00000000-0000-0000-0000-000000000402',
    NULL,
    NULL,
    'active',
    1
FROM tenants reseller
WHERE reseller.slug = 'demo-reseller'
ON CONFLICT (tenant_id, master_product_id) DO UPDATE
SET name_override = EXCLUDED.name_override,
    description_override = EXCLUDED.description_override,
    status = EXCLUDED.status,
    clone_version = EXCLUDED.clone_version;

INSERT INTO tenant_plans (
    tenant_plan_id,
    tenant_id,
    tenant_product_id,
    master_plan_id,
    selling_price_minor,
    reseller_cost_minor,
    currency,
    margin_policy,
    visibility,
    status,
    clone_version,
    product_snapshot,
    plan_snapshot,
    price_snapshot,
    capability_snapshot
)
SELECT
    seed_plan.tenant_plan_id::uuid,
    reseller.tenant_id,
    tenant_product.tenant_product_id,
    master_plan.plan_id,
    seed_plan.selling_price_minor,
    seed_plan.reseller_cost_minor,
    master_plan.currency,
    '{"min_margin_minor":100}'::jsonb,
    'public',
    'active',
    1,
    jsonb_build_object('product_id', master_product.product_id, 'name', master_product.name, 'product_type', master_product.product_type),
    jsonb_build_object('plan_id', master_plan.plan_id, 'plan_code', master_plan.plan_code, 'name', master_plan.name, 'specs', master_plan.specs),
    jsonb_build_object('selling_price_minor', seed_plan.selling_price_minor, 'currency', master_plan.currency),
    '{}'::jsonb
FROM (
    VALUES
        ('00000000-0000-0000-0000-000000000801', '00000000-0000-0000-0000-000000000501', 1400, 900),
        ('00000000-0000-0000-0000-000000000802', '00000000-0000-0000-0000-000000000502', 2500, 1700),
        ('00000000-0000-0000-0000-000000000803', '00000000-0000-0000-0000-000000000503', 750, 450)
) AS seed_plan(tenant_plan_id, master_plan_id, selling_price_minor, reseller_cost_minor)
JOIN tenants reseller ON reseller.slug = 'demo-reseller'
JOIN master_plans master_plan ON master_plan.plan_id = seed_plan.master_plan_id::uuid
JOIN master_products master_product ON master_product.product_id = master_plan.product_id
JOIN tenant_products tenant_product ON tenant_product.tenant_id = reseller.tenant_id
    AND tenant_product.master_product_id = master_product.product_id
ON CONFLICT (tenant_id, tenant_product_id, master_plan_id) DO UPDATE
SET selling_price_minor = EXCLUDED.selling_price_minor,
    reseller_cost_minor = EXCLUDED.reseller_cost_minor,
    currency = EXCLUDED.currency,
    margin_policy = EXCLUDED.margin_policy,
    visibility = EXCLUDED.visibility,
    status = EXCLUDED.status,
    clone_version = EXCLUDED.clone_version,
    product_snapshot = EXCLUDED.product_snapshot,
    plan_snapshot = EXCLUDED.plan_snapshot,
    price_snapshot = EXCLUDED.price_snapshot,
    capability_snapshot = EXCLUDED.capability_snapshot;
`
