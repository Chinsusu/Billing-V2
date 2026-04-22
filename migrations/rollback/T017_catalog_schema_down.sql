-- Manual rollback for T017 catalog schema migration.
-- Use only in a clean/dev environment or after owner approval because it drops catalog data.
-- The migration runner intentionally ignores this directory.

BEGIN;

DROP TABLE IF EXISTS tenant_plans;
DROP TABLE IF EXISTS tenant_products;
DROP TABLE IF EXISTS plan_sources;
DROP TABLE IF EXISTS provider_sources;
DROP TABLE IF EXISTS master_plans;
DROP TABLE IF EXISTS master_products;

DROP SEQUENCE IF EXISTS tenant_plans_display_id_seq;
DROP SEQUENCE IF EXISTS tenant_products_display_id_seq;
DROP SEQUENCE IF EXISTS plan_sources_display_id_seq;
DROP SEQUENCE IF EXISTS provider_sources_display_id_seq;
DROP SEQUENCE IF EXISTS master_plans_display_id_seq;
DROP SEQUENCE IF EXISTS master_products_display_id_seq;

DROP TYPE IF EXISTS catalog_tenant_plan_status;
DROP TYPE IF EXISTS catalog_tenant_plan_visibility;
DROP TYPE IF EXISTS catalog_tenant_product_status;
DROP TYPE IF EXISTS catalog_plan_source_status;
DROP TYPE IF EXISTS catalog_risk_level;
DROP TYPE IF EXISTS catalog_inventory_mode;
DROP TYPE IF EXISTS catalog_provider_source_status;
DROP TYPE IF EXISTS catalog_provider_type;
DROP TYPE IF EXISTS catalog_plan_status;
DROP TYPE IF EXISTS catalog_billing_cycle_type;
DROP TYPE IF EXISTS catalog_product_status;
DROP TYPE IF EXISTS catalog_product_type;

DO $$
BEGIN
    IF to_regclass('public.schema_migrations') IS NOT NULL THEN
        DELETE FROM schema_migrations
        WHERE version = '0006';
    END IF;
END
$$;

COMMIT;
