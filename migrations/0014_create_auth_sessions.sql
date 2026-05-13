CREATE TABLE auth_sessions (
    session_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(tenant_id),
    user_id UUID NOT NULL REFERENCES users(user_id),
    token_hash CHAR(64) NOT NULL UNIQUE,
    user_agent_hash CHAR(64),
    expires_at TIMESTAMPTZ NOT NULL,
    revoked_at TIMESTAMPTZ,
    last_seen_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT auth_sessions_expiry_after_created CHECK (expires_at > created_at)
);

CREATE INDEX idx_auth_sessions_user_tenant ON auth_sessions(tenant_id, user_id);
CREATE INDEX idx_auth_sessions_active_expiry ON auth_sessions(expires_at) WHERE revoked_at IS NULL;
