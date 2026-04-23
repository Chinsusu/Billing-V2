CREATE TYPE wallet_topup_payment_method AS ENUM (
    'bank_transfer',
    'crypto',
    'manual',
    'other'
);

CREATE TYPE wallet_topup_status AS ENUM (
    'draft',
    'submitted',
    'under_review',
    'approved',
    'rejected',
    'expired',
    'cancelled'
);

CREATE SEQUENCE topup_requests_display_id_seq AS BIGINT START WITH 10000;

CREATE TABLE topup_requests (
    topup_request_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    display_id BIGINT NOT NULL DEFAULT nextval('topup_requests_display_id_seq'),
    tenant_id UUID NOT NULL REFERENCES tenants(tenant_id),
    wallet_id UUID NOT NULL,
    requested_by UUID NOT NULL REFERENCES users(user_id),
    amount_minor BIGINT NOT NULL,
    currency TEXT NOT NULL,
    payment_method wallet_topup_payment_method NOT NULL,
    payment_reference TEXT,
    status wallet_topup_status NOT NULL DEFAULT 'submitted',
    reviewed_by UUID REFERENCES users(user_id),
    reviewed_at TIMESTAMPTZ,
    review_note TEXT,
    ledger_entry_id UUID REFERENCES wallet_ledger_entries(ledger_entry_id),
    idempotency_key TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT topup_requests_wallet_tenant_fk FOREIGN KEY (wallet_id, tenant_id) REFERENCES wallets(wallet_id, tenant_id),
    CONSTRAINT topup_requests_display_id_unique UNIQUE (display_id),
    CONSTRAINT topup_requests_amount_positive CHECK (amount_minor > 0),
    CONSTRAINT topup_requests_currency_format CHECK (currency ~ '^[A-Z]{3}$'),
    CONSTRAINT topup_requests_idempotency_not_blank CHECK (btrim(idempotency_key) <> ''),
    CONSTRAINT topup_requests_requester_idempotency_unique UNIQUE (tenant_id, requested_by, idempotency_key)
);

ALTER SEQUENCE topup_requests_display_id_seq OWNED BY topup_requests.display_id;

CREATE INDEX idx_topup_requests_tenant_status_created ON topup_requests(tenant_id, status, created_at);
CREATE INDEX idx_topup_requests_requested_by_created ON topup_requests(requested_by, created_at);
CREATE INDEX idx_topup_requests_wallet_created ON topup_requests(wallet_id, created_at);
