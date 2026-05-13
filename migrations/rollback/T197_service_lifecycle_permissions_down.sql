DELETE FROM role_permissions
WHERE permission_id IN (
    SELECT permission_id
    FROM permissions
    WHERE permission_key IN ('service.suspend', 'service.unsuspend', 'service.terminate')
);

DELETE FROM permissions
WHERE permission_key IN ('service.suspend', 'service.unsuspend', 'service.terminate');
