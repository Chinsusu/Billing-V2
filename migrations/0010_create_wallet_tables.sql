CREATE TYPE wallet_owner_type AS ENUM (
    'tenant',
    'user',
    'reseller_settlement',
    'platform'
);

CREATE TYPE wallet_status AS ENUM (
    'active',
    'frozen',
    'closed'
);

CREATE SEQUENCE wallets_display_id_seq AS BIGINT START WITH 10000;

CREATE TABLE wallets (
    wallet_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    display_id BIGINT NOT NULL DEFAULT nextval('wallets_display_id_seq'),
    tenant_id UUID NOT NULL REFERENCES tenants(tenant_id),
    owner_type wallet_owner_type NOT NULL,
    owner_id UUID NOT NULL,
    currency TEXT NOT NULL,
    status wallet_status NOT NULL DEFAULT 'active',
    available_balance_minor BIGINT NOT NULL DEFAULT 0,
    locked_balance_minor BIGINT NOT NULL DEFAULT 0,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT wallets_display_id_unique UNIQUE (display_id),
    CONSTRAINT wallets_owner_currency_unique UNIQUE (owner_type, owner_id, currency),
    CONSTRAINT wallets_currency_format CHECK (currency ~ '^[A-Z]{3}$'),
    CONSTRAINT wallets_available_balance_non_negative CHECK (available_balance_minor >= 0),
    CONSTRAINT wallets_locked_balance_non_negative CHECK (locked_balance_minor >= 0)
);

ALTER SEQUENCE wallets_display_id_seq OWNED BY wallets.display_id;

CREATE INDEX idx_wallets_tenant_status ON wallets(tenant_id, status, created_at);
CREATE INDEX idx_wallets_tenant_owner ON wallets(tenant_id, owner_type, owner_id);
