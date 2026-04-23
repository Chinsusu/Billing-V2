-- Manual rollback for T040 payment transaction schema migration.
-- Use only in a clean/dev environment or after owner approval because it drops payment history.
-- The migration runner intentionally ignores this directory.

BEGIN;

DROP TABLE IF EXISTS payment_transactions;

DROP SEQUENCE IF EXISTS payment_transactions_display_id_seq;

DROP TYPE IF EXISTS payment_transaction_status;
DROP TYPE IF EXISTS payment_transaction_type;

DO $$
BEGIN
    IF to_regclass('public.schema_migrations') IS NOT NULL THEN
        DELETE FROM schema_migrations
        WHERE version = '0009';
    END IF;
END
$$;

COMMIT;
