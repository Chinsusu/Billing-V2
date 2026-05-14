DELETE FROM role_permissions
USING permissions
WHERE role_permissions.permission_id = permissions.permission_id
  AND permissions.permission_key = 'service.renew';

DELETE FROM permissions
WHERE permission_key = 'service.renew';
