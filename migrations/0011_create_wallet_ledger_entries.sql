CREATE TYPE wallet_ledger_direction AS ENUM (
    'credit',
    'debit'
);

CREATE TYPE wallet_ledger_entry_type AS ENUM (
    'topup',
    'purchase',
    'reseller_cost',
    'refund',
    'adjustment',
    'reversal',
    'commission',
    'lock',
    'unlock'
);

CREATE TYPE wallet_ledger_status AS ENUM (
    'posted',
    'voided_by_reversal'
);

ALTER TABLE wallets
    ADD CONSTRAINT wallets_id_tenant_unique UNIQUE (wallet_id, tenant_id);

CREATE SEQUENCE wallet_ledger_entries_display_id_seq AS BIGINT START WITH 10000;

CREATE TABLE wallet_ledger_entries (
    ledger_entry_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    display_id BIGINT NOT NULL DEFAULT nextval('wallet_ledger_entries_display_id_seq'),
    wallet_id UUID NOT NULL,
    tenant_id UUID NOT NULL,
    direction wallet_ledger_direction NOT NULL,
    amount_minor BIGINT NOT NULL,
    currency TEXT NOT NULL,
    entry_type wallet_ledger_entry_type NOT NULL,
    status wallet_ledger_status NOT NULL DEFAULT 'posted',
    balance_after_minor BIGINT NOT NULL,
    reference_type TEXT NOT NULL,
    reference_id UUID NOT NULL,
    idempotency_key TEXT NOT NULL,
    created_by UUID REFERENCES users(user_id),
    reason TEXT,
    correlation_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT wallet_ledger_entries_wallet_tenant_fk FOREIGN KEY (wallet_id, tenant_id) REFERENCES wallets(wallet_id, tenant_id),
    CONSTRAINT wallet_ledger_entries_display_id_unique UNIQUE (display_id),
    CONSTRAINT wallet_ledger_entries_amount_positive CHECK (amount_minor > 0),
    CONSTRAINT wallet_ledger_entries_currency_format CHECK (currency ~ '^[A-Z]{3}$'),
    CONSTRAINT wallet_ledger_entries_balance_non_negative CHECK (balance_after_minor >= 0),
    CONSTRAINT wallet_ledger_entries_reference_type_not_blank CHECK (btrim(reference_type) <> ''),
    CONSTRAINT wallet_ledger_entries_idempotency_not_blank CHECK (btrim(idempotency_key) <> ''),
    CONSTRAINT wallet_ledger_entries_adjustment_reason CHECK (entry_type <> 'adjustment' OR btrim(COALESCE(reason, '')) <> ''),
    CONSTRAINT wallet_ledger_entries_wallet_idempotency_unique UNIQUE (wallet_id, idempotency_key)
);

ALTER SEQUENCE wallet_ledger_entries_display_id_seq OWNED BY wallet_ledger_entries.display_id;

CREATE INDEX idx_wallet_ledger_entries_wallet_created ON wallet_ledger_entries(wallet_id, created_at);
CREATE INDEX idx_wallet_ledger_entries_tenant_reference ON wallet_ledger_entries(tenant_id, reference_type, reference_id);
CREATE INDEX idx_wallet_ledger_entries_correlation ON wallet_ledger_entries(correlation_id);
CREATE INDEX idx_wallet_ledger_entries_tenant_type_status ON wallet_ledger_entries(tenant_id, entry_type, status, created_at);
