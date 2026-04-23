-- Manual rollback for T041 wallet schema migration.
-- Use only in a clean/dev environment or after owner approval because it drops wallet balance data.
-- The migration runner intentionally ignores this directory.

BEGIN;

DROP TABLE IF EXISTS wallets;

DROP SEQUENCE IF EXISTS wallets_display_id_seq;

DROP TYPE IF EXISTS wallet_status;
DROP TYPE IF EXISTS wallet_owner_type;

DO $$
BEGIN
    IF to_regclass('public.schema_migrations') IS NOT NULL THEN
        DELETE FROM schema_migrations
        WHERE version = '0010';
    END IF;
END
$$;

COMMIT;
