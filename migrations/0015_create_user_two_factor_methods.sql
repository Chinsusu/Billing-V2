CREATE TABLE user_two_factor_methods (
    tenant_id UUID NOT NULL REFERENCES tenants(tenant_id),
    user_id UUID NOT NULL REFERENCES users(user_id),
    method VARCHAR(20) NOT NULL DEFAULT 'totp',
    secret_ciphertext TEXT NOT NULL,
    enabled_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (tenant_id, user_id, method),
    CONSTRAINT user_two_factor_methods_method_totp CHECK (method = 'totp')
);

ALTER TABLE auth_sessions
    ADD COLUMN two_factor_satisfied_at TIMESTAMPTZ;

CREATE INDEX idx_auth_sessions_two_factor_satisfied ON auth_sessions(two_factor_satisfied_at)
    WHERE two_factor_satisfied_at IS NOT NULL;
