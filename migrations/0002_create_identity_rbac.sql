CREATE TYPE user_type AS ENUM ('platform_staff', 'reseller_staff', 'client');
CREATE TYPE user_status AS ENUM ('active', 'suspended', 'disabled', 'pending_verification');
CREATE TYPE two_factor_status AS ENUM ('required', 'enabled', 'disabled');
CREATE TYPE permission_risk_level AS ENUM ('low', 'medium', 'high', 'critical');

CREATE TABLE users (
    user_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(tenant_id),
    email VARCHAR(255) NOT NULL,
    email_verified_at TIMESTAMPTZ,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255),
    user_type user_type NOT NULL,
    status user_status NOT NULL DEFAULT 'pending_verification',
    two_factor_status two_factor_status NOT NULL DEFAULT 'disabled',
    last_login_at TIMESTAMPTZ,
    failed_login_count INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT users_email_lowercase CHECK (email = lower(email)),
    CONSTRAINT users_failed_login_count_non_negative CHECK (failed_login_count >= 0),
    CONSTRAINT users_unique_tenant_email UNIQUE (tenant_id, email)
);

CREATE INDEX idx_users_tenant_status ON users(tenant_id, status);
CREATE INDEX idx_users_type ON users(user_type);

ALTER TABLE tenants
    ADD CONSTRAINT tenants_owner_user_fk FOREIGN KEY (owner_user_id) REFERENCES users(user_id);

CREATE TABLE roles (
    role_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID REFERENCES tenants(tenant_id),
    role_key VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    is_system BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT roles_key_lowercase CHECK (role_key = lower(role_key)),
    CONSTRAINT roles_system_tenant_rule CHECK ((is_system = TRUE AND tenant_id IS NULL) OR (is_system = FALSE AND tenant_id IS NOT NULL))
);

CREATE UNIQUE INDEX idx_roles_system_role_key ON roles(role_key) WHERE tenant_id IS NULL;
CREATE UNIQUE INDEX idx_roles_tenant_role_key ON roles(tenant_id, role_key) WHERE tenant_id IS NOT NULL;
CREATE INDEX idx_roles_tenant_id ON roles(tenant_id);

CREATE TABLE permissions (
    permission_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    permission_key VARCHAR(255) NOT NULL,
    module VARCHAR(255) NOT NULL,
    risk_level permission_risk_level NOT NULL,
    CONSTRAINT permissions_key_lowercase CHECK (permission_key = lower(permission_key)),
    CONSTRAINT permissions_unique_key UNIQUE (permission_key)
);

CREATE INDEX idx_permissions_module ON permissions(module);
CREATE INDEX idx_permissions_risk_level ON permissions(risk_level);

CREATE TABLE role_permissions (
    role_id UUID NOT NULL REFERENCES roles(role_id),
    permission_id UUID NOT NULL REFERENCES permissions(permission_id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (role_id, permission_id)
);

CREATE TABLE user_roles (
    user_id UUID NOT NULL REFERENCES users(user_id),
    tenant_id UUID NOT NULL REFERENCES tenants(tenant_id),
    role_id UUID NOT NULL REFERENCES roles(role_id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, tenant_id, role_id)
);

CREATE INDEX idx_user_roles_tenant_id ON user_roles(tenant_id);
