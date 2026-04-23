-- Manual rollback for T028 order schema migration.
-- Use only in a clean/dev environment or after owner approval because it drops order and service data.
-- The migration runner intentionally ignores this directory.

BEGIN;

DROP TABLE IF EXISTS service_instances;
DROP TABLE IF EXISTS order_provisioning_jobs;
DROP TABLE IF EXISTS order_reservations;
DROP TABLE IF EXISTS orders;

DROP SEQUENCE IF EXISTS service_instances_display_id_seq;
DROP SEQUENCE IF EXISTS order_provisioning_jobs_display_id_seq;
DROP SEQUENCE IF EXISTS order_reservations_display_id_seq;
DROP SEQUENCE IF EXISTS orders_display_id_seq;

DROP TYPE IF EXISTS order_service_suspension_reason;
DROP TYPE IF EXISTS order_service_status;
DROP TYPE IF EXISTS order_provisioning_status;
DROP TYPE IF EXISTS order_reservation_status;
DROP TYPE IF EXISTS order_billing_status;
DROP TYPE IF EXISTS order_status;

DO $$
BEGIN
    IF to_regclass('public.schema_migrations') IS NOT NULL THEN
        DELETE FROM schema_migrations
        WHERE version = '0007';
    END IF;
END
$$;

COMMIT;
