-- Manual rollback for T046 invoice generation idempotency index.
-- Use only after owner approval because dropping this index allows duplicate invoices for one order.
-- The migration runner intentionally ignores this directory.

BEGIN;

DROP INDEX IF EXISTS idx_invoices_tenant_order_unique;

DO $$
BEGIN
    IF to_regclass('public.schema_migrations') IS NOT NULL THEN
        DELETE FROM schema_migrations
        WHERE version = '0013';
    END IF;
END
$$;

COMMIT;
