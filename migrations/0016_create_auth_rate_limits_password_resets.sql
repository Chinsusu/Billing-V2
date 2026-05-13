CREATE TABLE auth_rate_limit_counters (
    action VARCHAR(80) NOT NULL,
    key_hash CHAR(64) NOT NULL,
    window_start TIMESTAMPTZ NOT NULL,
    attempt_count INT NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (action, key_hash, window_start),
    CONSTRAINT auth_rate_limit_counters_attempt_positive CHECK (attempt_count > 0)
);

CREATE INDEX idx_auth_rate_limit_counters_updated ON auth_rate_limit_counters(updated_at);

CREATE TABLE auth_password_reset_tokens (
    reset_token_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(tenant_id),
    user_id UUID NOT NULL REFERENCES users(user_id),
    token_hash CHAR(64) NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    used_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT auth_password_reset_tokens_expiry_after_created CHECK (expires_at > created_at)
);

CREATE INDEX idx_auth_password_reset_tokens_user ON auth_password_reset_tokens(tenant_id, user_id, created_at DESC);
CREATE INDEX idx_auth_password_reset_tokens_active_expiry ON auth_password_reset_tokens(expires_at)
    WHERE used_at IS NULL;
