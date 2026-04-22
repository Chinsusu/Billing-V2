CREATE TYPE catalog_product_type AS ENUM (
    'vps',
    'proxy',
    'service_addon'
);

CREATE TYPE catalog_product_status AS ENUM (
    'draft',
    'active',
    'disabled',
    'archived'
);

CREATE TYPE catalog_billing_cycle_type AS ENUM (
    'day',
    'month_30d',
    'calendar_month',
    'custom'
);

CREATE TYPE catalog_plan_status AS ENUM (
    'draft',
    'active',
    'disabled',
    'archived'
);

CREATE TYPE catalog_provider_type AS ENUM (
    'manual',
    'proxmox',
    'ovh',
    'hetzner',
    'proxy_upstream',
    'preloaded_proxy_pool',
    'custom_api'
);

CREATE TYPE catalog_provider_source_status AS ENUM (
    'active',
    'disabled',
    'maintenance',
    'out_of_stock'
);

CREATE TYPE catalog_inventory_mode AS ENUM (
    'finite_stock',
    'provider_live',
    'manual_unlimited',
    'preloaded_list'
);

CREATE TYPE catalog_risk_level AS ENUM (
    'low',
    'medium',
    'high'
);

CREATE TYPE catalog_plan_source_status AS ENUM (
    'active',
    'disabled'
);

CREATE TYPE catalog_tenant_product_status AS ENUM (
    'active',
    'hidden',
    'disabled'
);

CREATE TYPE catalog_tenant_plan_visibility AS ENUM (
    'public',
    'hidden',
    'private'
);

CREATE TYPE catalog_tenant_plan_status AS ENUM (
    'active',
    'disabled',
    'margin_risk',
    'archived'
);

CREATE SEQUENCE master_products_display_id_seq AS BIGINT START WITH 10000;
CREATE SEQUENCE master_plans_display_id_seq AS BIGINT START WITH 10000;
CREATE SEQUENCE provider_sources_display_id_seq AS BIGINT START WITH 10000;
CREATE SEQUENCE plan_sources_display_id_seq AS BIGINT START WITH 10000;
CREATE SEQUENCE tenant_products_display_id_seq AS BIGINT START WITH 10000;
CREATE SEQUENCE tenant_plans_display_id_seq AS BIGINT START WITH 10000;

CREATE TABLE master_products (
    product_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    display_id BIGINT NOT NULL DEFAULT nextval('master_products_display_id_seq'),
    product_type catalog_product_type NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    status catalog_product_status NOT NULL DEFAULT 'draft',
    display_order INT NOT NULL DEFAULT 0,
    created_by UUID NOT NULL REFERENCES users(user_id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT master_products_display_id_unique UNIQUE (display_id),
    CONSTRAINT master_products_name_not_blank CHECK (btrim(name) <> ''),
    CONSTRAINT master_products_display_order_non_negative CHECK (display_order >= 0)
);

ALTER SEQUENCE master_products_display_id_seq OWNED BY master_products.display_id;

CREATE INDEX idx_master_products_type_status_order ON master_products(product_type, status, display_order, created_at);
CREATE INDEX idx_master_products_created_by ON master_products(created_by);

CREATE TABLE master_plans (
    plan_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    display_id BIGINT NOT NULL DEFAULT nextval('master_plans_display_id_seq'),
    product_id UUID NOT NULL REFERENCES master_products(product_id),
    plan_code TEXT NOT NULL,
    name TEXT NOT NULL,
    specs JSONB NOT NULL DEFAULT '{}'::jsonb,
    billing_cycle_type catalog_billing_cycle_type NOT NULL,
    billing_cycle_value INT NOT NULL,
    base_cost_minor BIGINT NOT NULL,
    suggested_price_minor BIGINT NOT NULL,
    reseller_min_price_minor BIGINT NOT NULL DEFAULT 0,
    currency TEXT NOT NULL,
    status catalog_plan_status NOT NULL DEFAULT 'draft',
    version INT NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT master_plans_display_id_unique UNIQUE (display_id),
    CONSTRAINT master_plans_plan_code_not_blank CHECK (btrim(plan_code) <> ''),
    CONSTRAINT master_plans_name_not_blank CHECK (btrim(name) <> ''),
    CONSTRAINT master_plans_billing_cycle_value_positive CHECK (billing_cycle_value > 0),
    CONSTRAINT master_plans_base_cost_non_negative CHECK (base_cost_minor >= 0),
    CONSTRAINT master_plans_suggested_price_non_negative CHECK (suggested_price_minor >= 0),
    CONSTRAINT master_plans_reseller_min_price_non_negative CHECK (reseller_min_price_minor >= 0),
    CONSTRAINT master_plans_currency_format CHECK (currency ~ '^[A-Z]{3}$'),
    CONSTRAINT master_plans_version_positive CHECK (version > 0),
    CONSTRAINT master_plans_unique_code_version UNIQUE (plan_code, version)
);

ALTER SEQUENCE master_plans_display_id_seq OWNED BY master_plans.display_id;

CREATE INDEX idx_master_plans_product_status ON master_plans(product_id, status, version DESC);
CREATE INDEX idx_master_plans_status_created ON master_plans(status, created_at);

CREATE TABLE provider_sources (
    source_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    display_id BIGINT NOT NULL DEFAULT nextval('provider_sources_display_id_seq'),
    source_type catalog_provider_type NOT NULL,
    name TEXT NOT NULL,
    provider_account_id UUID,
    location TEXT,
    status catalog_provider_source_status NOT NULL DEFAULT 'disabled',
    capability_profile JSONB NOT NULL DEFAULT '{}'::jsonb,
    inventory_mode catalog_inventory_mode NOT NULL,
    risk_level catalog_risk_level NOT NULL DEFAULT 'medium',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT provider_sources_display_id_unique UNIQUE (display_id),
    CONSTRAINT provider_sources_name_not_blank CHECK (btrim(name) <> '')
);

ALTER SEQUENCE provider_sources_display_id_seq OWNED BY provider_sources.display_id;

CREATE INDEX idx_provider_sources_type_status ON provider_sources(source_type, status);
CREATE INDEX idx_provider_sources_location_status ON provider_sources(location, status) WHERE location IS NOT NULL;
CREATE INDEX idx_provider_sources_account ON provider_sources(provider_account_id) WHERE provider_account_id IS NOT NULL;

CREATE TABLE plan_sources (
    plan_source_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    display_id BIGINT NOT NULL DEFAULT nextval('plan_sources_display_id_seq'),
    plan_id UUID NOT NULL REFERENCES master_plans(plan_id) ON DELETE CASCADE,
    source_id UUID NOT NULL REFERENCES provider_sources(source_id),
    priority INT NOT NULL,
    cost_override_minor BIGINT NOT NULL DEFAULT 0,
    capacity_policy JSONB NOT NULL DEFAULT '{}'::jsonb,
    capability_override JSONB NOT NULL DEFAULT '{}'::jsonb,
    status catalog_plan_source_status NOT NULL DEFAULT 'disabled',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT plan_sources_display_id_unique UNIQUE (display_id),
    CONSTRAINT plan_sources_priority_positive CHECK (priority > 0),
    CONSTRAINT plan_sources_cost_override_non_negative CHECK (cost_override_minor >= 0),
    CONSTRAINT plan_sources_unique_plan_source UNIQUE (plan_id, source_id)
);

ALTER SEQUENCE plan_sources_display_id_seq OWNED BY plan_sources.display_id;

CREATE INDEX idx_plan_sources_plan_status_priority ON plan_sources(plan_id, status, priority, created_at);
CREATE INDEX idx_plan_sources_source_status ON plan_sources(source_id, status);

CREATE TABLE tenant_products (
    tenant_product_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    display_id BIGINT NOT NULL DEFAULT nextval('tenant_products_display_id_seq'),
    tenant_id UUID NOT NULL REFERENCES tenants(tenant_id),
    master_product_id UUID NOT NULL REFERENCES master_products(product_id),
    name_override TEXT,
    description_override TEXT,
    status catalog_tenant_product_status NOT NULL DEFAULT 'hidden',
    clone_version INT NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT tenant_products_display_id_unique UNIQUE (display_id),
    CONSTRAINT tenant_products_clone_version_positive CHECK (clone_version > 0),
    CONSTRAINT tenant_products_unique_master_product UNIQUE (tenant_id, master_product_id)
);

ALTER SEQUENCE tenant_products_display_id_seq OWNED BY tenant_products.display_id;

CREATE INDEX idx_tenant_products_tenant_status ON tenant_products(tenant_id, status, created_at);
CREATE INDEX idx_tenant_products_master_product ON tenant_products(master_product_id);

CREATE TABLE tenant_plans (
    tenant_plan_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    display_id BIGINT NOT NULL DEFAULT nextval('tenant_plans_display_id_seq'),
    tenant_id UUID NOT NULL REFERENCES tenants(tenant_id),
    tenant_product_id UUID NOT NULL REFERENCES tenant_products(tenant_product_id),
    master_plan_id UUID NOT NULL REFERENCES master_plans(plan_id),
    selling_price_minor BIGINT NOT NULL,
    reseller_cost_minor BIGINT NOT NULL,
    currency TEXT NOT NULL,
    margin_policy JSONB NOT NULL DEFAULT '{}'::jsonb,
    visibility catalog_tenant_plan_visibility NOT NULL DEFAULT 'hidden',
    status catalog_tenant_plan_status NOT NULL DEFAULT 'disabled',
    clone_version INT NOT NULL DEFAULT 1,
    product_snapshot JSONB NOT NULL DEFAULT '{}'::jsonb,
    plan_snapshot JSONB NOT NULL DEFAULT '{}'::jsonb,
    price_snapshot JSONB NOT NULL DEFAULT '{}'::jsonb,
    capability_snapshot JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT tenant_plans_display_id_unique UNIQUE (display_id),
    CONSTRAINT tenant_plans_selling_price_non_negative CHECK (selling_price_minor >= 0),
    CONSTRAINT tenant_plans_reseller_cost_non_negative CHECK (reseller_cost_minor >= 0),
    CONSTRAINT tenant_plans_currency_format CHECK (currency ~ '^[A-Z]{3}$'),
    CONSTRAINT tenant_plans_clone_version_positive CHECK (clone_version > 0),
    CONSTRAINT tenant_plans_unique_master_plan UNIQUE (tenant_id, tenant_product_id, master_plan_id)
);

ALTER SEQUENCE tenant_plans_display_id_seq OWNED BY tenant_plans.display_id;

CREATE INDEX idx_tenant_plans_tenant_status_visibility ON tenant_plans(tenant_id, status, visibility, created_at);
CREATE INDEX idx_tenant_plans_tenant_product ON tenant_plans(tenant_product_id, status);
CREATE INDEX idx_tenant_plans_master_plan ON tenant_plans(master_plan_id);
