CREATE TYPE order_status AS ENUM (
    'draft',
    'pending_payment',
    'paid',
    'cancelled',
    'failed',
    'refunded'
);

CREATE TYPE order_billing_status AS ENUM (
    'unpaid',
    'paid',
    'overdue',
    'refunded',
    'partially_refunded'
);

CREATE TYPE order_reservation_status AS ENUM (
    'pending_reserve',
    'reserved',
    'reservation_expired',
    'reservation_released',
    'allocated'
);

CREATE TYPE order_provisioning_status AS ENUM (
    'queued',
    'provisioning',
    'provisioned',
    'failed',
    'manual_review'
);

CREATE TYPE order_service_status AS ENUM (
    'active',
    'suspended',
    'expired',
    'cancelled',
    'terminated'
);

CREATE TYPE order_service_suspension_reason AS ENUM (
    'expiry',
    'manual_admin',
    'manual_reseller',
    'abuse',
    'system_issue'
);

CREATE SEQUENCE orders_display_id_seq AS BIGINT START WITH 10000;
CREATE SEQUENCE order_reservations_display_id_seq AS BIGINT START WITH 10000;
CREATE SEQUENCE order_provisioning_jobs_display_id_seq AS BIGINT START WITH 10000;
CREATE SEQUENCE service_instances_display_id_seq AS BIGINT START WITH 10000;

CREATE TABLE orders (
    order_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    display_id BIGINT NOT NULL DEFAULT nextval('orders_display_id_seq'),
    tenant_id UUID NOT NULL REFERENCES tenants(tenant_id),
    buyer_user_id UUID NOT NULL REFERENCES users(user_id),
    tenant_plan_id UUID NOT NULL REFERENCES tenant_plans(tenant_plan_id),
    quantity INT NOT NULL DEFAULT 1,
    currency TEXT NOT NULL,
    unit_price_minor BIGINT NOT NULL,
    discount_minor BIGINT NOT NULL DEFAULT 0,
    total_minor BIGINT NOT NULL,
    order_status order_status NOT NULL DEFAULT 'pending_payment',
    billing_status order_billing_status NOT NULL DEFAULT 'unpaid',
    idempotency_key TEXT NOT NULL,
    product_snapshot JSONB NOT NULL DEFAULT '{}'::jsonb,
    plan_snapshot JSONB NOT NULL DEFAULT '{}'::jsonb,
    price_snapshot JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT orders_display_id_unique UNIQUE (display_id),
    CONSTRAINT orders_quantity_positive CHECK (quantity > 0),
    CONSTRAINT orders_unit_price_non_negative CHECK (unit_price_minor >= 0),
    CONSTRAINT orders_discount_non_negative CHECK (discount_minor >= 0),
    CONSTRAINT orders_total_non_negative CHECK (total_minor >= 0),
    CONSTRAINT orders_currency_format CHECK (currency ~ '^[A-Z]{3}$'),
    CONSTRAINT orders_idempotency_key_not_blank CHECK (btrim(idempotency_key) <> ''),
    CONSTRAINT orders_tenant_idempotency_unique UNIQUE (tenant_id, idempotency_key)
);

ALTER SEQUENCE orders_display_id_seq OWNED BY orders.display_id;

CREATE INDEX idx_orders_tenant_status_created ON orders(tenant_id, order_status, created_at);
CREATE INDEX idx_orders_tenant_billing_created ON orders(tenant_id, billing_status, created_at);
CREATE INDEX idx_orders_buyer_created ON orders(buyer_user_id, created_at);
CREATE INDEX idx_orders_tenant_plan ON orders(tenant_plan_id);

CREATE TABLE order_reservations (
    reservation_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    display_id BIGINT NOT NULL DEFAULT nextval('order_reservations_display_id_seq'),
    order_id UUID NOT NULL REFERENCES orders(order_id) ON DELETE CASCADE,
    tenant_id UUID NOT NULL REFERENCES tenants(tenant_id),
    provider_source_id UUID NOT NULL REFERENCES provider_sources(source_id),
    status order_reservation_status NOT NULL DEFAULT 'pending_reserve',
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT order_reservations_display_id_unique UNIQUE (display_id)
);

ALTER SEQUENCE order_reservations_display_id_seq OWNED BY order_reservations.display_id;

CREATE INDEX idx_order_reservations_order ON order_reservations(order_id);
CREATE INDEX idx_order_reservations_tenant_status ON order_reservations(tenant_id, status, created_at);
CREATE INDEX idx_order_reservations_status_expiry ON order_reservations(status, expires_at);
CREATE INDEX idx_order_reservations_source_status ON order_reservations(provider_source_id, status);

CREATE TABLE order_provisioning_jobs (
    provisioning_job_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    display_id BIGINT NOT NULL DEFAULT nextval('order_provisioning_jobs_display_id_seq'),
    order_id UUID NOT NULL REFERENCES orders(order_id) ON DELETE CASCADE,
    tenant_id UUID NOT NULL REFERENCES tenants(tenant_id),
    provider_source_id UUID NOT NULL REFERENCES provider_sources(source_id),
    provider_operation_id TEXT NOT NULL,
    status order_provisioning_status NOT NULL DEFAULT 'queued',
    idempotency_key TEXT NOT NULL,
    attempt_number INT NOT NULL DEFAULT 1,
    last_error_code TEXT,
    last_error_message TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT order_provisioning_jobs_display_id_unique UNIQUE (display_id),
    CONSTRAINT order_provisioning_jobs_operation_not_blank CHECK (btrim(provider_operation_id) <> ''),
    CONSTRAINT order_provisioning_jobs_idempotency_not_blank CHECK (btrim(idempotency_key) <> ''),
    CONSTRAINT order_provisioning_jobs_attempt_positive CHECK (attempt_number > 0),
    CONSTRAINT order_provisioning_jobs_operation_unique UNIQUE (provider_operation_id),
    CONSTRAINT order_provisioning_jobs_tenant_idempotency_unique UNIQUE (tenant_id, idempotency_key)
);

ALTER SEQUENCE order_provisioning_jobs_display_id_seq OWNED BY order_provisioning_jobs.display_id;

CREATE INDEX idx_order_provisioning_jobs_order ON order_provisioning_jobs(order_id);
CREATE INDEX idx_order_provisioning_jobs_tenant_status ON order_provisioning_jobs(tenant_id, status, created_at);
CREATE INDEX idx_order_provisioning_jobs_source_status ON order_provisioning_jobs(provider_source_id, status);

CREATE TABLE service_instances (
    service_instance_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    display_id BIGINT NOT NULL DEFAULT nextval('service_instances_display_id_seq'),
    tenant_id UUID NOT NULL REFERENCES tenants(tenant_id),
    order_id UUID NOT NULL REFERENCES orders(order_id),
    tenant_plan_id UUID NOT NULL REFERENCES tenant_plans(tenant_plan_id),
    provider_source_id UUID NOT NULL REFERENCES provider_sources(source_id),
    external_resource_id TEXT NOT NULL,
    status order_service_status NOT NULL DEFAULT 'active',
    billing_status order_billing_status NOT NULL DEFAULT 'paid',
    suspension_reason order_service_suspension_reason,
    term_start TIMESTAMPTZ NOT NULL,
    term_end TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT service_instances_display_id_unique UNIQUE (display_id),
    CONSTRAINT service_instances_external_resource_not_blank CHECK (btrim(external_resource_id) <> ''),
    CONSTRAINT service_instances_term_window CHECK (term_end > term_start),
    CONSTRAINT service_instances_suspended_reason CHECK (status <> 'suspended' OR suspension_reason IS NOT NULL),
    CONSTRAINT service_instances_order_unique UNIQUE (order_id),
    CONSTRAINT service_instances_source_resource_unique UNIQUE (provider_source_id, external_resource_id)
);

ALTER SEQUENCE service_instances_display_id_seq OWNED BY service_instances.display_id;

CREATE INDEX idx_service_instances_tenant_status ON service_instances(tenant_id, status, created_at);
CREATE INDEX idx_service_instances_tenant_billing ON service_instances(tenant_id, billing_status, created_at);
CREATE INDEX idx_service_instances_tenant_plan ON service_instances(tenant_plan_id);
CREATE INDEX idx_service_instances_term_end ON service_instances(status, term_end);
CREATE INDEX idx_service_instances_source ON service_instances(provider_source_id);
