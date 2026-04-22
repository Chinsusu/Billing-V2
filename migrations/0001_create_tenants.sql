CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TYPE tenant_type AS ENUM ('platform', 'reseller', 'direct_store');
CREATE TYPE tenant_status AS ENUM ('active', 'suspended', 'disabled', 'pending_setup');
CREATE TYPE domain_type AS ENUM ('system_subdomain', 'custom_domain');
CREATE TYPE domain_verification_status AS ENUM ('pending', 'verified', 'failed', 'disabled');
CREATE TYPE tls_status AS ENUM ('pending', 'active', 'failed', 'expired');

CREATE TABLE tenants (
    tenant_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    parent_tenant_id UUID REFERENCES tenants(tenant_id),
    tenant_type tenant_type NOT NULL,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL,
    status tenant_status NOT NULL DEFAULT 'pending_setup',
    default_currency VARCHAR(3) NOT NULL,
    timezone VARCHAR(50) NOT NULL DEFAULT 'Asia/Ho_Chi_Minh',
    owner_user_id UUID,
    branding_settings JSONB NOT NULL DEFAULT '{}'::jsonb,
    billing_settings JSONB NOT NULL DEFAULT '{}'::jsonb,
    risk_settings JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT tenants_slug_lowercase CHECK (slug = lower(slug)),
    CONSTRAINT tenants_currency_uppercase CHECK (default_currency = upper(default_currency)),
    CONSTRAINT tenants_unique_slug UNIQUE (slug)
);

CREATE INDEX idx_tenants_parent_tenant_id ON tenants(parent_tenant_id);
CREATE INDEX idx_tenants_tenant_type ON tenants(tenant_type);
CREATE INDEX idx_tenants_status ON tenants(status);
CREATE INDEX idx_tenants_owner_user_id ON tenants(owner_user_id);

CREATE TABLE tenant_domains (
    domain_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(tenant_id),
    domain VARCHAR(255) NOT NULL,
    domain_type domain_type NOT NULL,
    verification_status domain_verification_status NOT NULL DEFAULT 'pending',
    verification_token_hash VARCHAR(255),
    tls_status tls_status NOT NULL DEFAULT 'pending',
    is_primary BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT tenant_domains_domain_lowercase CHECK (domain = lower(domain)),
    CONSTRAINT tenant_domains_unique_domain UNIQUE (domain)
);

CREATE INDEX idx_tenant_domains_tenant_id ON tenant_domains(tenant_id);
CREATE INDEX idx_tenant_domains_status ON tenant_domains(verification_status, tls_status);
CREATE UNIQUE INDEX idx_tenant_domains_primary ON tenant_domains(tenant_id) WHERE is_primary = TRUE;
