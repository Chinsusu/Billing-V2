CREATE TYPE service_credential_type AS ENUM (
    'vps_root',
    'proxy_auth',
    'ssh_key',
    'console_url',
    'api_token',
    'recovery_code'
);

CREATE TYPE service_credential_status AS ENUM (
    'active',
    'rotated',
    'revoked',
    'expired'
);

ALTER TABLE service_instances
    ADD CONSTRAINT service_instances_id_tenant_unique UNIQUE (service_instance_id, tenant_id);

CREATE TABLE service_credentials (
    credential_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    service_instance_id UUID NOT NULL,
    credential_type service_credential_type NOT NULL,
    encrypted_payload TEXT NOT NULL,
    encryption_key_version TEXT NOT NULL DEFAULT 'v1',
    encryption_algorithm TEXT NOT NULL DEFAULT 'aes-256-gcm',
    secret_version TEXT,
    masked_hint TEXT NOT NULL,
    status service_credential_status NOT NULL DEFAULT 'active',
    last_revealed_at TIMESTAMPTZ,
    last_revealed_by UUID REFERENCES users(user_id),
    rotated_at TIMESTAMPTZ,
    rotated_by UUID REFERENCES users(user_id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT service_credentials_payload_not_blank CHECK (btrim(encrypted_payload) <> ''),
    CONSTRAINT service_credentials_key_version_not_blank CHECK (btrim(encryption_key_version) <> ''),
    CONSTRAINT service_credentials_algorithm_not_blank CHECK (btrim(encryption_algorithm) <> ''),
    CONSTRAINT service_credentials_masked_hint_not_blank CHECK (btrim(masked_hint) <> ''),
    CONSTRAINT service_credentials_tenant_fk FOREIGN KEY (tenant_id) REFERENCES tenants(tenant_id),
    CONSTRAINT service_credentials_service_tenant_fk FOREIGN KEY (service_instance_id, tenant_id)
        REFERENCES service_instances(service_instance_id, tenant_id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX idx_service_credentials_active_type
    ON service_credentials(tenant_id, service_instance_id, credential_type)
    WHERE status = 'active';

CREATE INDEX idx_service_credentials_service_status ON service_credentials(service_instance_id, status, created_at);
CREATE INDEX idx_service_credentials_tenant_status ON service_credentials(tenant_id, status, created_at);
CREATE INDEX idx_service_credentials_last_revealed ON service_credentials(last_revealed_at) WHERE last_revealed_at IS NOT NULL;
