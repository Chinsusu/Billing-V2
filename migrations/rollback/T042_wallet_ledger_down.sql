-- Manual rollback for T042 wallet ledger migration.
-- Use only in a clean/dev environment or after owner approval because it drops immutable ledger history.
-- The migration runner intentionally ignores this directory.

BEGIN;

DROP TABLE IF EXISTS wallet_ledger_entries;

DROP SEQUENCE IF EXISTS wallet_ledger_entries_display_id_seq;

ALTER TABLE wallets
    DROP CONSTRAINT IF EXISTS wallets_id_tenant_unique;

DROP TYPE IF EXISTS wallet_ledger_status;
DROP TYPE IF EXISTS wallet_ledger_entry_type;
DROP TYPE IF EXISTS wallet_ledger_direction;

DO $$
BEGIN
    IF to_regclass('public.schema_migrations') IS NOT NULL THEN
        DELETE FROM schema_migrations
        WHERE version = '0011';
    END IF;
END
$$;

COMMIT;
