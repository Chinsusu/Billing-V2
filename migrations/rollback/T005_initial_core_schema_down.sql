-- Manual rollback for T005 initial core schema migrations.
-- Use only in a clean/dev environment or after owner approval because it drops data.
-- The migration runner intentionally ignores this directory.

BEGIN;

DROP TABLE IF EXISTS audit_logs;
DROP TYPE IF EXISTS audit_actor_type;

DROP TABLE IF EXISTS user_roles;
DROP TABLE IF EXISTS role_permissions;
DROP TABLE IF EXISTS permissions;
DROP TABLE IF EXISTS roles;

ALTER TABLE IF EXISTS tenants DROP CONSTRAINT IF EXISTS tenants_owner_user_fk;
DROP TABLE IF EXISTS users;

DROP TABLE IF EXISTS tenant_domains;
DROP TABLE IF EXISTS tenants;

DROP TYPE IF EXISTS permission_risk_level;
DROP TYPE IF EXISTS two_factor_status;
DROP TYPE IF EXISTS user_status;
DROP TYPE IF EXISTS user_type;
DROP TYPE IF EXISTS tls_status;
DROP TYPE IF EXISTS domain_verification_status;
DROP TYPE IF EXISTS domain_type;
DROP TYPE IF EXISTS tenant_status;
DROP TYPE IF EXISTS tenant_type;

DO $$
BEGIN
    IF to_regclass('public.schema_migrations') IS NOT NULL THEN
        DELETE FROM schema_migrations
        WHERE version IN ('0001', '0002', '0003');
    END IF;
END
$$;

COMMIT;
