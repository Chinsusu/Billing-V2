CREATE TABLE service_credential_reveal_rate_limits (
    tenant_id UUID NOT NULL REFERENCES tenants(tenant_id),
    actor_id UUID NOT NULL REFERENCES users(user_id),
    service_instance_id UUID NOT NULL,
    window_start TIMESTAMPTZ NOT NULL,
    attempt_count INT NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (tenant_id, actor_id, service_instance_id, window_start),
    CONSTRAINT service_credential_reveal_rate_limits_attempt_positive CHECK (attempt_count > 0),
    CONSTRAINT service_credential_reveal_rate_limits_service_fk FOREIGN KEY (service_instance_id, tenant_id)
        REFERENCES service_instances(service_instance_id, tenant_id) ON DELETE CASCADE
);

CREATE INDEX idx_service_credential_reveal_rate_limits_updated
    ON service_credential_reveal_rate_limits(updated_at);
