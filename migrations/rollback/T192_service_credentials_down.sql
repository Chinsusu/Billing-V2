DROP TABLE IF EXISTS service_credentials;
ALTER TABLE service_instances DROP CONSTRAINT IF EXISTS service_instances_id_tenant_unique;
DROP TYPE IF EXISTS service_credential_status;
DROP TYPE IF EXISTS service_credential_type;
