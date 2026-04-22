CREATE SEQUENCE tenants_display_id_seq AS BIGINT START WITH 10000;

ALTER TABLE tenants
    ADD COLUMN display_id BIGINT;

WITH ordered AS (
    SELECT tenant_id, 9999 + ROW_NUMBER() OVER (ORDER BY created_at, tenant_id) AS display_id
    FROM tenants
)
UPDATE tenants target
SET display_id = ordered.display_id
FROM ordered
WHERE target.tenant_id = ordered.tenant_id;

SELECT setval('tenants_display_id_seq', COALESCE((SELECT MAX(display_id) FROM tenants), 9999), TRUE);

ALTER TABLE tenants
    ALTER COLUMN display_id SET DEFAULT nextval('tenants_display_id_seq'),
    ALTER COLUMN display_id SET NOT NULL,
    ADD CONSTRAINT tenants_display_id_unique UNIQUE (display_id);

ALTER SEQUENCE tenants_display_id_seq OWNED BY tenants.display_id;

CREATE SEQUENCE tenant_domains_display_id_seq AS BIGINT START WITH 10000;

ALTER TABLE tenant_domains
    ADD COLUMN display_id BIGINT;

WITH ordered AS (
    SELECT domain_id, 9999 + ROW_NUMBER() OVER (ORDER BY created_at, domain_id) AS display_id
    FROM tenant_domains
)
UPDATE tenant_domains target
SET display_id = ordered.display_id
FROM ordered
WHERE target.domain_id = ordered.domain_id;

SELECT setval('tenant_domains_display_id_seq', COALESCE((SELECT MAX(display_id) FROM tenant_domains), 9999), TRUE);

ALTER TABLE tenant_domains
    ALTER COLUMN display_id SET DEFAULT nextval('tenant_domains_display_id_seq'),
    ALTER COLUMN display_id SET NOT NULL,
    ADD CONSTRAINT tenant_domains_display_id_unique UNIQUE (display_id);

ALTER SEQUENCE tenant_domains_display_id_seq OWNED BY tenant_domains.display_id;

CREATE SEQUENCE users_display_id_seq AS BIGINT START WITH 10000;

ALTER TABLE users
    ADD COLUMN display_id BIGINT;

WITH ordered AS (
    SELECT user_id, 9999 + ROW_NUMBER() OVER (ORDER BY created_at, user_id) AS display_id
    FROM users
)
UPDATE users target
SET display_id = ordered.display_id
FROM ordered
WHERE target.user_id = ordered.user_id;

SELECT setval('users_display_id_seq', COALESCE((SELECT MAX(display_id) FROM users), 9999), TRUE);

ALTER TABLE users
    ALTER COLUMN display_id SET DEFAULT nextval('users_display_id_seq'),
    ALTER COLUMN display_id SET NOT NULL,
    ADD CONSTRAINT users_display_id_unique UNIQUE (display_id);

ALTER SEQUENCE users_display_id_seq OWNED BY users.display_id;

CREATE SEQUENCE roles_display_id_seq AS BIGINT START WITH 10000;

ALTER TABLE roles
    ADD COLUMN display_id BIGINT;

WITH ordered AS (
    SELECT role_id, 9999 + ROW_NUMBER() OVER (ORDER BY created_at, role_id) AS display_id
    FROM roles
)
UPDATE roles target
SET display_id = ordered.display_id
FROM ordered
WHERE target.role_id = ordered.role_id;

SELECT setval('roles_display_id_seq', COALESCE((SELECT MAX(display_id) FROM roles), 9999), TRUE);

ALTER TABLE roles
    ALTER COLUMN display_id SET DEFAULT nextval('roles_display_id_seq'),
    ALTER COLUMN display_id SET NOT NULL,
    ADD CONSTRAINT roles_display_id_unique UNIQUE (display_id);

ALTER SEQUENCE roles_display_id_seq OWNED BY roles.display_id;

CREATE SEQUENCE audit_logs_display_id_seq AS BIGINT START WITH 10000;

ALTER TABLE audit_logs
    ADD COLUMN display_id BIGINT;

WITH ordered AS (
    SELECT audit_id, 9999 + ROW_NUMBER() OVER (ORDER BY created_at, audit_id) AS display_id
    FROM audit_logs
)
UPDATE audit_logs target
SET display_id = ordered.display_id
FROM ordered
WHERE target.audit_id = ordered.audit_id;

SELECT setval('audit_logs_display_id_seq', COALESCE((SELECT MAX(display_id) FROM audit_logs), 9999), TRUE);

ALTER TABLE audit_logs
    ALTER COLUMN display_id SET DEFAULT nextval('audit_logs_display_id_seq'),
    ALTER COLUMN display_id SET NOT NULL,
    ADD CONSTRAINT audit_logs_display_id_unique UNIQUE (display_id);

ALTER SEQUENCE audit_logs_display_id_seq OWNED BY audit_logs.display_id;

CREATE SEQUENCE outbox_events_display_id_seq AS BIGINT START WITH 10000;

ALTER TABLE outbox_events
    ADD COLUMN display_id BIGINT;

WITH ordered AS (
    SELECT outbox_event_id, 9999 + ROW_NUMBER() OVER (ORDER BY created_at, outbox_event_id) AS display_id
    FROM outbox_events
)
UPDATE outbox_events target
SET display_id = ordered.display_id
FROM ordered
WHERE target.outbox_event_id = ordered.outbox_event_id;

SELECT setval('outbox_events_display_id_seq', COALESCE((SELECT MAX(display_id) FROM outbox_events), 9999), TRUE);

ALTER TABLE outbox_events
    ALTER COLUMN display_id SET DEFAULT nextval('outbox_events_display_id_seq'),
    ALTER COLUMN display_id SET NOT NULL,
    ADD CONSTRAINT outbox_events_display_id_unique UNIQUE (display_id);

ALTER SEQUENCE outbox_events_display_id_seq OWNED BY outbox_events.display_id;

CREATE SEQUENCE jobs_display_id_seq AS BIGINT START WITH 10000;

ALTER TABLE jobs
    ADD COLUMN display_id BIGINT;

WITH ordered AS (
    SELECT job_id, 9999 + ROW_NUMBER() OVER (ORDER BY created_at, job_id) AS display_id
    FROM jobs
)
UPDATE jobs target
SET display_id = ordered.display_id
FROM ordered
WHERE target.job_id = ordered.job_id;

SELECT setval('jobs_display_id_seq', COALESCE((SELECT MAX(display_id) FROM jobs), 9999), TRUE);

ALTER TABLE jobs
    ALTER COLUMN display_id SET DEFAULT nextval('jobs_display_id_seq'),
    ALTER COLUMN display_id SET NOT NULL,
    ADD CONSTRAINT jobs_display_id_unique UNIQUE (display_id);

ALTER SEQUENCE jobs_display_id_seq OWNED BY jobs.display_id;

CREATE SEQUENCE job_attempts_display_id_seq AS BIGINT START WITH 10000;

ALTER TABLE job_attempts
    ADD COLUMN display_id BIGINT;

WITH ordered AS (
    SELECT job_attempt_id, 9999 + ROW_NUMBER() OVER (ORDER BY started_at, job_attempt_id) AS display_id
    FROM job_attempts
)
UPDATE job_attempts target
SET display_id = ordered.display_id
FROM ordered
WHERE target.job_attempt_id = ordered.job_attempt_id;

SELECT setval('job_attempts_display_id_seq', COALESCE((SELECT MAX(display_id) FROM job_attempts), 9999), TRUE);

ALTER TABLE job_attempts
    ALTER COLUMN display_id SET DEFAULT nextval('job_attempts_display_id_seq'),
    ALTER COLUMN display_id SET NOT NULL,
    ADD CONSTRAINT job_attempts_display_id_unique UNIQUE (display_id);

ALTER SEQUENCE job_attempts_display_id_seq OWNED BY job_attempts.display_id;
