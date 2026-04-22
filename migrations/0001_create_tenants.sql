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
    status tenant_status NOT NULL,
    default_currency VARCHAR(3) NOT NULL,
    timezone VARCHAR(50) NOT NULL DEFAULT 'Asia/Ho_Chi_Minh',
    owner_user_id UUID,
    branding_settings JSONB,
    billing_settings JSONB,
    risk_settings JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT unique_tenant_slug UNIQUE (slug)
);

CREATE INDEX idx_tenants_tenant_type ON tenants(tenant_type);
CREATE INDEX idx_tenants_status ON tenants(status);
CREATE INDEX idx_tenants_owner_user_id ON tenants(owner_user_id);

CREATE TABLE tenant_domains (
    domain_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(tenant_id),
    domain VARCHAR(255) NOT NULL,
    domain_type domain_type NOT NULL,
    verification_status domain_verification_status NOT NULL,
    verification_token_hash VARCHAR(255),
    tls_status tls_status NOT NULL,
    is_primary BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT unique_tenant_domain UNIQUE (domain)
);

CREATE UNIQUE INDEX idx_unique_primary_domain ON tenant_domains(tenant_id) WHERE is_primary = TRUE;
