CREATE TYPE audit_actor_type AS ENUM ('user', 'system', 'worker', 'provider_webhook');

CREATE TABLE audit_logs (
    audit_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID REFERENCES tenants(tenant_id),
    actor_id UUID,
    actor_type audit_actor_type NOT NULL,
    action VARCHAR(255) NOT NULL,
    target_type VARCHAR(255) NOT NULL,
    target_id UUID NOT NULL,
    before_snapshot_redacted JSONB,
    after_snapshot_redacted JSONB,
    metadata_redacted JSONB NOT NULL DEFAULT '{}'::jsonb,
    ip_address VARCHAR(45),
    user_agent TEXT,
    correlation_id UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT audit_logs_action_not_empty CHECK (action <> ''),
    CONSTRAINT audit_logs_target_type_not_empty CHECK (target_type <> '')
);

CREATE INDEX idx_audit_logs_tenant_created ON audit_logs(tenant_id, created_at);
CREATE INDEX idx_audit_logs_actor_created ON audit_logs(actor_id, created_at);
CREATE INDEX idx_audit_logs_target_type_id ON audit_logs(target_type, target_id);
CREATE INDEX idx_audit_logs_correlation_id ON audit_logs(correlation_id);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);
