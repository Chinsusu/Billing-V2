-- Manual rollback for T039 invoice schema migration.
-- Use only in a clean/dev environment or after owner approval because it drops invoice data.
-- The migration runner intentionally ignores this directory.

BEGIN;

DROP TABLE IF EXISTS invoice_items;
DROP TABLE IF EXISTS invoices;

DROP SEQUENCE IF EXISTS invoices_display_id_seq;

DROP TYPE IF EXISTS invoice_status;

DO $$
BEGIN
    IF to_regclass('public.schema_migrations') IS NOT NULL THEN
        DELETE FROM schema_migrations
        WHERE version = '0008';
    END IF;
END
$$;

COMMIT;
