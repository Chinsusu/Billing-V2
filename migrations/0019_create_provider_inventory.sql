CREATE TYPE provider_inventory_status AS ENUM (
    'active',
    'out_of_stock',
    'sync_failed',
    'disabled'
);

CREATE SEQUENCE provider_inventory_display_id_seq AS BIGINT START WITH 10000;

CREATE TABLE provider_inventory (
    inventory_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    display_id BIGINT NOT NULL DEFAULT nextval('provider_inventory_display_id_seq'),
    source_id UUID NOT NULL REFERENCES provider_sources(source_id),
    capacity_total INT,
    reserved_count INT NOT NULL DEFAULT 0,
    allocated_count INT NOT NULL DEFAULT 0,
    available_count_cache INT NOT NULL DEFAULT 0,
    last_synced_at TIMESTAMPTZ,
    status provider_inventory_status NOT NULL DEFAULT 'active',
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT provider_inventory_display_id_unique UNIQUE (display_id),
    CONSTRAINT provider_inventory_source_unique UNIQUE (source_id),
    CONSTRAINT provider_inventory_capacity_non_negative CHECK (capacity_total IS NULL OR capacity_total >= 0),
    CONSTRAINT provider_inventory_reserved_non_negative CHECK (reserved_count >= 0),
    CONSTRAINT provider_inventory_allocated_non_negative CHECK (allocated_count >= 0),
    CONSTRAINT provider_inventory_available_non_negative CHECK (available_count_cache >= 0),
    CONSTRAINT provider_inventory_counts_within_capacity CHECK (
        capacity_total IS NULL OR capacity_total >= reserved_count + allocated_count
    )
);

ALTER SEQUENCE provider_inventory_display_id_seq OWNED BY provider_inventory.display_id;

CREATE INDEX idx_provider_inventory_status ON provider_inventory(status);
CREATE INDEX idx_provider_inventory_source_status ON provider_inventory(source_id, status);

ALTER TABLE order_reservations
    ADD COLUMN quantity INT NOT NULL DEFAULT 1,
    ADD CONSTRAINT order_reservations_quantity_positive CHECK (quantity > 0);
