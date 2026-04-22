-- Manual rollback for T006 outbox/jobs migration.
-- Use only in a clean/dev environment or after owner approval because it drops data.
-- The migration runner intentionally ignores this directory.

BEGIN;

DROP TABLE IF EXISTS job_attempts;
DROP TABLE IF EXISTS jobs;
DROP TABLE IF EXISTS outbox_events;

DROP TYPE IF EXISTS job_attempt_result;
DROP TYPE IF EXISTS job_status;
DROP TYPE IF EXISTS outbox_event_status;

DO $$
BEGIN
    IF to_regclass('public.schema_migrations') IS NOT NULL THEN
        DELETE FROM schema_migrations
        WHERE version = '0004';
    END IF;
END
$$;

COMMIT;
