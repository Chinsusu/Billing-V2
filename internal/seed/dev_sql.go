package seed

const seedDevUsersSQL = `
INSERT INTO users (user_id, tenant_id, email, password_hash, full_name, user_type, status, two_factor_status)
SELECT
    '00000000-0000-0000-0000-000000000101',
    platform.tenant_id,
    'admin@local.billing',
    'dev-only-placeholder-hash',
    'Platform Admin',
    'platform_staff',
    'active',
    'disabled'
FROM tenants platform
WHERE platform.slug = 'platform'
ON CONFLICT (tenant_id, email) DO UPDATE
SET full_name = EXCLUDED.full_name,
    user_type = EXCLUDED.user_type,
    status = EXCLUDED.status,
    two_factor_status = EXCLUDED.two_factor_status;

INSERT INTO users (user_id, tenant_id, email, password_hash, full_name, user_type, status, two_factor_status)
SELECT
    '00000000-0000-0000-0000-000000000102',
    reseller.tenant_id,
    'reseller@local.billing',
    'dev-only-placeholder-hash',
    'Demo Reseller Owner',
    'reseller_staff',
    'active',
    'disabled'
FROM tenants reseller
WHERE reseller.slug = 'demo-reseller'
ON CONFLICT (tenant_id, email) DO UPDATE
SET full_name = EXCLUDED.full_name,
    user_type = EXCLUDED.user_type,
    status = EXCLUDED.status,
    two_factor_status = EXCLUDED.two_factor_status;
`

const seedSystemRolesSQL = `
INSERT INTO roles (role_id, tenant_id, role_key, name, is_system)
VALUES
    ('00000000-0000-0000-0000-000000000201', NULL, 'platform_admin', 'Platform Admin', TRUE),
    ('00000000-0000-0000-0000-000000000202', NULL, 'catalog_manager', 'Catalog Manager', TRUE),
    ('00000000-0000-0000-0000-000000000203', NULL, 'reseller_catalog_manager', 'Reseller Catalog Manager', TRUE),
    ('00000000-0000-0000-0000-000000000204', NULL, 'customer_catalog_viewer', 'Customer Catalog Viewer', TRUE)
ON CONFLICT (role_key) WHERE tenant_id IS NULL DO UPDATE
SET name = EXCLUDED.name,
    is_system = EXCLUDED.is_system;
`

const seedRolePermissionsSQL = `
INSERT INTO role_permissions (role_id, permission_id)
SELECT platform_admin.role_id, permissions.permission_id
FROM roles platform_admin
CROSS JOIN permissions
WHERE platform_admin.role_key = 'platform_admin'
  AND platform_admin.is_system = TRUE
ON CONFLICT (role_id, permission_id) DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT role.role_id, permission.permission_id
FROM roles role
JOIN permissions permission ON permission.permission_key IN ('catalog.view', 'catalog.manage', 'provider.view')
WHERE role.role_key = 'catalog_manager'
  AND role.is_system = TRUE
ON CONFLICT (role_id, permission_id) DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT role.role_id, permission.permission_id
FROM roles role
JOIN permissions permission ON permission.permission_key IN ('catalog.view', 'catalog.manage')
WHERE role.role_key = 'reseller_catalog_manager'
  AND role.is_system = TRUE
ON CONFLICT (role_id, permission_id) DO NOTHING;

INSERT INTO role_permissions (role_id, permission_id)
SELECT role.role_id, permission.permission_id
FROM roles role
JOIN permissions permission ON permission.permission_key = 'catalog.view'
WHERE role.role_key = 'customer_catalog_viewer'
  AND role.is_system = TRUE
ON CONFLICT (role_id, permission_id) DO NOTHING;
`

const seedUserRolesSQL = `
INSERT INTO user_roles (user_id, tenant_id, role_id)
SELECT admin.user_id, admin.tenant_id, role.role_id
FROM users admin
JOIN tenants platform ON platform.tenant_id = admin.tenant_id
JOIN roles role ON role.role_key = 'platform_admin' AND role.is_system = TRUE
WHERE platform.slug = 'platform'
  AND admin.email = 'admin@local.billing'
ON CONFLICT (user_id, tenant_id, role_id) DO NOTHING;

INSERT INTO user_roles (user_id, tenant_id, role_id)
SELECT reseller_user.user_id, reseller_user.tenant_id, role.role_id
FROM users reseller_user
JOIN tenants reseller ON reseller.tenant_id = reseller_user.tenant_id
JOIN roles role ON role.role_key = 'reseller_catalog_manager' AND role.is_system = TRUE
WHERE reseller.slug = 'demo-reseller'
  AND reseller_user.email = 'reseller@local.billing'
ON CONFLICT (user_id, tenant_id, role_id) DO NOTHING;
`
