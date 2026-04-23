CREATE TYPE payment_transaction_type AS ENUM (
    'charge',
    'refund',
    'adjustment'
);

CREATE TYPE payment_transaction_status AS ENUM (
    'pending',
    'posted',
    'failed',
    'voided'
);

CREATE SEQUENCE payment_transactions_display_id_seq AS BIGINT START WITH 10000;

CREATE TABLE payment_transactions (
    payment_transaction_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    display_id BIGINT NOT NULL DEFAULT nextval('payment_transactions_display_id_seq'),
    tenant_id UUID NOT NULL REFERENCES tenants(tenant_id),
    account_user_id UUID NOT NULL REFERENCES users(user_id),
    order_id UUID REFERENCES orders(order_id),
    invoice_id UUID REFERENCES invoices(invoice_id),
    transaction_type payment_transaction_type NOT NULL,
    status payment_transaction_status NOT NULL DEFAULT 'posted',
    currency TEXT NOT NULL,
    amount_minor BIGINT NOT NULL,
    description TEXT,
    idempotency_key TEXT NOT NULL,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT payment_transactions_display_id_unique UNIQUE (display_id),
    CONSTRAINT payment_transactions_currency_format CHECK (currency ~ '^[A-Z]{3}$'),
    CONSTRAINT payment_transactions_amount_positive CHECK (amount_minor > 0),
    CONSTRAINT payment_transactions_idempotency_not_blank CHECK (btrim(idempotency_key) <> ''),
    CONSTRAINT payment_transactions_tenant_idempotency_unique UNIQUE (tenant_id, idempotency_key)
);

ALTER SEQUENCE payment_transactions_display_id_seq OWNED BY payment_transactions.display_id;

CREATE INDEX idx_payment_transactions_tenant_account_created ON payment_transactions(tenant_id, account_user_id, created_at);
CREATE INDEX idx_payment_transactions_tenant_type_status_created ON payment_transactions(tenant_id, transaction_type, status, created_at);
CREATE INDEX idx_payment_transactions_order ON payment_transactions(order_id) WHERE order_id IS NOT NULL;
CREATE INDEX idx_payment_transactions_invoice ON payment_transactions(invoice_id) WHERE invoice_id IS NOT NULL;
