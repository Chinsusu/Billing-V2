-- Manual rollback for T044 top-up request migration.
-- Use only in a clean/dev environment or after owner approval because it drops funding request history.
-- The migration runner intentionally ignores this directory.

BEGIN;

DROP TABLE IF EXISTS topup_requests;

DROP SEQUENCE IF EXISTS topup_requests_display_id_seq;

DROP TYPE IF EXISTS wallet_topup_status;
DROP TYPE IF EXISTS wallet_topup_payment_method;

DO $$
BEGIN
    IF to_regclass('public.schema_migrations') IS NOT NULL THEN
        DELETE FROM schema_migrations
        WHERE version = '0012';
    END IF;
END
$$;

COMMIT;
