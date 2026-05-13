DELETE FROM role_permissions
WHERE permission_id IN (
    SELECT permission_id
    FROM permissions
    WHERE permission_key IN ('ticket.manage', 'risk.flag.create', 'abuse.case.manage')
);

DELETE FROM permissions
WHERE permission_key IN ('ticket.manage', 'risk.flag.create', 'abuse.case.manage');
