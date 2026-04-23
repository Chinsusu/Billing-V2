CREATE TYPE invoice_status AS ENUM (
    'draft',
    'issued',
    'paid',
    'partially_paid',
    'overdue',
    'voided'
);

CREATE SEQUENCE invoices_display_id_seq AS BIGINT START WITH 10000;

CREATE TABLE invoices (
    invoice_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    display_id BIGINT NOT NULL DEFAULT nextval('invoices_display_id_seq'),
    tenant_id UUID NOT NULL REFERENCES tenants(tenant_id),
    buyer_user_id UUID NOT NULL REFERENCES users(user_id),
    order_id UUID REFERENCES orders(order_id),
    status invoice_status NOT NULL DEFAULT 'draft',
    currency TEXT NOT NULL,
    subtotal_minor BIGINT NOT NULL,
    tax_minor BIGINT NOT NULL DEFAULT 0,
    discount_minor BIGINT NOT NULL DEFAULT 0,
    total_minor BIGINT NOT NULL,
    issued_at TIMESTAMPTZ,
    due_at TIMESTAMPTZ,
    paid_at TIMESTAMPTZ,
    voided_at TIMESTAMPTZ,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT invoices_display_id_unique UNIQUE (display_id),
    CONSTRAINT invoices_id_tenant_unique UNIQUE (invoice_id, tenant_id),
    CONSTRAINT invoices_currency_format CHECK (currency ~ '^[A-Z]{3}$'),
    CONSTRAINT invoices_subtotal_non_negative CHECK (subtotal_minor >= 0),
    CONSTRAINT invoices_tax_non_negative CHECK (tax_minor >= 0),
    CONSTRAINT invoices_discount_non_negative CHECK (discount_minor >= 0),
    CONSTRAINT invoices_total_non_negative CHECK (total_minor >= 0),
    CONSTRAINT invoices_total_matches_parts CHECK (total_minor = subtotal_minor + tax_minor - discount_minor)
);

ALTER SEQUENCE invoices_display_id_seq OWNED BY invoices.display_id;

CREATE INDEX idx_invoices_tenant_status_created ON invoices(tenant_id, status, created_at);
CREATE INDEX idx_invoices_buyer_created ON invoices(buyer_user_id, created_at);
CREATE INDEX idx_invoices_order ON invoices(order_id) WHERE order_id IS NOT NULL;
CREATE INDEX idx_invoices_due_status ON invoices(status, due_at) WHERE due_at IS NOT NULL;

CREATE TABLE invoice_items (
    invoice_item_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    invoice_id UUID NOT NULL,
    tenant_id UUID NOT NULL,
    order_id UUID REFERENCES orders(order_id),
    order_item_id UUID,
    service_instance_id UUID REFERENCES service_instances(service_instance_id),
    description TEXT NOT NULL,
    quantity INT NOT NULL DEFAULT 1,
    unit_price_minor BIGINT NOT NULL,
    tax_minor BIGINT NOT NULL DEFAULT 0,
    discount_minor BIGINT NOT NULL DEFAULT 0,
    line_total_minor BIGINT NOT NULL,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT invoice_items_invoice_tenant_fk FOREIGN KEY (invoice_id, tenant_id) REFERENCES invoices(invoice_id, tenant_id) ON DELETE CASCADE,
    CONSTRAINT invoice_items_description_not_blank CHECK (btrim(description) <> ''),
    CONSTRAINT invoice_items_quantity_positive CHECK (quantity > 0),
    CONSTRAINT invoice_items_unit_price_non_negative CHECK (unit_price_minor >= 0),
    CONSTRAINT invoice_items_tax_non_negative CHECK (tax_minor >= 0),
    CONSTRAINT invoice_items_discount_non_negative CHECK (discount_minor >= 0),
    CONSTRAINT invoice_items_line_total_non_negative CHECK (line_total_minor >= 0),
    CONSTRAINT invoice_items_total_matches_parts CHECK (line_total_minor = (quantity::BIGINT * unit_price_minor) + tax_minor - discount_minor)
);

CREATE INDEX idx_invoice_items_invoice ON invoice_items(invoice_id);
CREATE INDEX idx_invoice_items_tenant_created ON invoice_items(tenant_id, created_at);
CREATE INDEX idx_invoice_items_order ON invoice_items(order_id) WHERE order_id IS NOT NULL;
CREATE INDEX idx_invoice_items_order_item ON invoice_items(order_item_id) WHERE order_item_id IS NOT NULL;
CREATE INDEX idx_invoice_items_service ON invoice_items(service_instance_id) WHERE service_instance_id IS NOT NULL;
