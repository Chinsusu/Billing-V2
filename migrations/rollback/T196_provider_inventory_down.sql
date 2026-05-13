ALTER TABLE order_reservations
    DROP CONSTRAINT IF EXISTS order_reservations_quantity_positive,
    DROP COLUMN IF EXISTS quantity;

DROP TABLE IF EXISTS provider_inventory;
DROP SEQUENCE IF EXISTS provider_inventory_display_id_seq;
DROP TYPE IF EXISTS provider_inventory_status;
