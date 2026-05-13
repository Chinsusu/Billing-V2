INSERT INTO permissions (permission_key, module, risk_level)
VALUES
    ('service.suspend', 'service', 'high'),
    ('service.unsuspend', 'service', 'high'),
    ('service.terminate', 'service', 'critical')
ON CONFLICT (permission_key) DO UPDATE
SET module = EXCLUDED.module,
    risk_level = EXCLUDED.risk_level;

INSERT INTO role_permissions (role_id, permission_id)
SELECT role.role_id, permission.permission_id
FROM roles role
JOIN permissions permission ON permission.permission_key IN ('service.suspend', 'service.unsuspend', 'service.terminate')
WHERE role.role_key IN ('platform_admin', 'reseller_admin')
  AND role.is_system = TRUE
ON CONFLICT (role_id, permission_id) DO NOTHING;
