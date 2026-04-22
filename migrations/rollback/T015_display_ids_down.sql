ALTER TABLE job_attempts
    DROP CONSTRAINT IF EXISTS job_attempts_display_id_unique;
ALTER TABLE job_attempts
    DROP COLUMN IF EXISTS display_id;
DROP SEQUENCE IF EXISTS job_attempts_display_id_seq;

ALTER TABLE jobs
    DROP CONSTRAINT IF EXISTS jobs_display_id_unique;
ALTER TABLE jobs
    DROP COLUMN IF EXISTS display_id;
DROP SEQUENCE IF EXISTS jobs_display_id_seq;

ALTER TABLE outbox_events
    DROP CONSTRAINT IF EXISTS outbox_events_display_id_unique;
ALTER TABLE outbox_events
    DROP COLUMN IF EXISTS display_id;
DROP SEQUENCE IF EXISTS outbox_events_display_id_seq;

ALTER TABLE audit_logs
    DROP CONSTRAINT IF EXISTS audit_logs_display_id_unique;
ALTER TABLE audit_logs
    DROP COLUMN IF EXISTS display_id;
DROP SEQUENCE IF EXISTS audit_logs_display_id_seq;

ALTER TABLE roles
    DROP CONSTRAINT IF EXISTS roles_display_id_unique;
ALTER TABLE roles
    DROP COLUMN IF EXISTS display_id;
DROP SEQUENCE IF EXISTS roles_display_id_seq;

ALTER TABLE users
    DROP CONSTRAINT IF EXISTS users_display_id_unique;
ALTER TABLE users
    DROP COLUMN IF EXISTS display_id;
DROP SEQUENCE IF EXISTS users_display_id_seq;

ALTER TABLE tenant_domains
    DROP CONSTRAINT IF EXISTS tenant_domains_display_id_unique;
ALTER TABLE tenant_domains
    DROP COLUMN IF EXISTS display_id;
DROP SEQUENCE IF EXISTS tenant_domains_display_id_seq;

ALTER TABLE tenants
    DROP CONSTRAINT IF EXISTS tenants_display_id_unique;
ALTER TABLE tenants
    DROP COLUMN IF EXISTS display_id;
DROP SEQUENCE IF EXISTS tenants_display_id_seq;
