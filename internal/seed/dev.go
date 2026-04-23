package seed

import (
	"context"
	"fmt"

	platformdb "github.com/Chinsusu/Billing-V2/internal/platform/db"
)

type Statement struct {
	Name string
	SQL  string
}

func DevStatements() []Statement {
	return []Statement{
		{Name: "rbac_permissions", SQL: seedRBACPermissionsSQL},
		{Name: "platform_tenant", SQL: seedPlatformTenantSQL},
		{Name: "demo_reseller_tenant", SQL: seedDemoResellerTenantSQL},
		{Name: "dev_users", SQL: seedDevUsersSQL},
		{Name: "system_roles", SQL: seedSystemRolesSQL},
		{Name: "role_permissions", SQL: seedRolePermissionsSQL},
		{Name: "user_roles", SQL: seedUserRolesSQL},
		{Name: "provider_sources", SQL: seedProviderSourcesSQL},
		{Name: "master_products", SQL: seedMasterProductsSQL},
		{Name: "master_plans", SQL: seedMasterPlansSQL},
		{Name: "plan_sources", SQL: seedPlanSourcesSQL},
		{Name: "demo_reseller_catalog", SQL: seedDemoResellerCatalogSQL},
		{Name: "billing_flow", SQL: seedBillingFlowSQL},
	}
}

func ApplyDev(ctx context.Context, executor platformdb.Executor) error {
	if executor == nil {
		return fmt.Errorf("seed executor is required")
	}
	for _, statement := range DevStatements() {
		if _, err := executor.ExecContext(ctx, statement.SQL); err != nil {
			return fmt.Errorf("apply seed %s: %w", statement.Name, err)
		}
	}
	return nil
}

const seedRBACPermissionsSQL = `
INSERT INTO permissions (permission_key, module, risk_level)
VALUES
    ('tenant.view', 'tenant', 'low'),
    ('tenant.create', 'tenant', 'high'),
    ('tenant.update', 'tenant', 'high'),
    ('tenant.domain.manage', 'tenant', 'high'),
    ('wallet.view', 'wallet', 'low'),
    ('wallet.topup.approve', 'wallet', 'critical'),
    ('wallet.adjustment.create', 'wallet', 'critical'),
    ('order.view', 'order', 'low'),
    ('order.create', 'order', 'medium'),
    ('order.manage', 'order', 'high'),
    ('service.view', 'service', 'low'),
    ('service.credential.reveal', 'service', 'critical'),
    ('provisioning.job.retry', 'provisioning', 'high'),
    ('provisioning.manual_review.resolve', 'provisioning', 'high'),
    ('provider.view', 'provider', 'low'),
    ('provider.manage', 'provider', 'high'),
    ('catalog.view', 'catalog', 'low'),
    ('catalog.manage', 'catalog', 'high'),
    ('audit.view', 'audit', 'high'),
    ('tenant.emergency_access', 'tenant', 'critical')
ON CONFLICT (permission_key) DO UPDATE
SET module = EXCLUDED.module,
    risk_level = EXCLUDED.risk_level;
`

const seedPlatformTenantSQL = `
INSERT INTO tenants (tenant_id, tenant_type, name, slug, status, default_currency, timezone)
VALUES ('00000000-0000-0000-0000-000000000001', 'platform', 'Billing Platform', 'platform', 'active', 'USD', 'Asia/Ho_Chi_Minh')
ON CONFLICT (slug) DO UPDATE
SET tenant_type = EXCLUDED.tenant_type,
    name = EXCLUDED.name,
    status = EXCLUDED.status,
    default_currency = EXCLUDED.default_currency,
    timezone = EXCLUDED.timezone;
`

const seedDemoResellerTenantSQL = `
INSERT INTO tenants (tenant_id, parent_tenant_id, tenant_type, name, slug, status, default_currency, timezone)
SELECT
    '00000000-0000-0000-0000-000000000010',
    platform.tenant_id,
    'reseller',
    'Demo Reseller',
    'demo-reseller',
    'active',
    'USD',
    'Asia/Ho_Chi_Minh'
FROM tenants platform
WHERE platform.slug = 'platform'
ON CONFLICT (slug) DO UPDATE
SET parent_tenant_id = EXCLUDED.parent_tenant_id,
    tenant_type = EXCLUDED.tenant_type,
    name = EXCLUDED.name,
    status = EXCLUDED.status,
    default_currency = EXCLUDED.default_currency,
    timezone = EXCLUDED.timezone;
`
