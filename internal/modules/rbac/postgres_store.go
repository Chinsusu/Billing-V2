package rbac

import (
	"context"
	"errors"
	"fmt"

	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
	platformdb "github.com/Chinsusu/Billing-V2/internal/platform/db"
	"github.com/lib/pq"
)

var ErrRBACStoreExecutorMissing = errors.New("rbac store executor missing")

type PostgresStore struct {
	executor platformdb.Executor
}

func NewPostgresStore(executor platformdb.Executor) *PostgresStore {
	return &PostgresStore{executor: executor}
}

func (store *PostgresStore) ListRoleIDsForUser(ctx context.Context, tenantID tenant.ID, userID identity.UserID) ([]identity.RoleID, error) {
	if err := store.ready(); err != nil {
		return nil, err
	}
	if err := validateTenantUserLookup(tenantID, userID); err != nil {
		return nil, err
	}
	rows, err := store.executor.QueryContext(ctx, `
SELECT role_id
FROM user_roles
WHERE tenant_id = $1 AND user_id = $2
ORDER BY role_id`, tenantID, userID)
	if err != nil {
		return nil, fmt.Errorf("list user roles: %w", err)
	}
	defer rows.Close()

	roleIDs := make([]identity.RoleID, 0)
	for rows.Next() {
		var roleID string
		if err := rows.Scan(&roleID); err != nil {
			return nil, fmt.Errorf("scan user role: %w", err)
		}
		roleIDs = append(roleIDs, identity.RoleID(roleID))
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("read user roles: %w", err)
	}
	return roleIDs, nil
}

func (store *PostgresStore) ListPermissionsForUser(ctx context.Context, tenantID tenant.ID, userID identity.UserID) (PermissionSet, error) {
	if err := store.ready(); err != nil {
		return nil, err
	}
	if err := validateTenantUserLookup(tenantID, userID); err != nil {
		return nil, err
	}
	rows, err := store.executor.QueryContext(ctx, `
SELECT DISTINCT p.permission_key
FROM user_roles ur
JOIN role_permissions rp ON rp.role_id = ur.role_id
JOIN permissions p ON p.permission_id = rp.permission_id
WHERE ur.tenant_id = $1 AND ur.user_id = $2
ORDER BY p.permission_key`, tenantID, userID)
	if err != nil {
		return nil, fmt.Errorf("list user permissions: %w", err)
	}
	defer rows.Close()
	return scanPermissionSet(rows)
}

func (store *PostgresStore) ListPermissionsForRoles(ctx context.Context, roleIDs []identity.RoleID) (PermissionSet, error) {
	if err := store.ready(); err != nil {
		return nil, err
	}
	if len(roleIDs) == 0 {
		return NewPermissionSet(), nil
	}
	rows, err := store.executor.QueryContext(ctx, `
SELECT DISTINCT p.permission_key
FROM role_permissions rp
JOIN permissions p ON p.permission_id = rp.permission_id
WHERE rp.role_id = ANY($1)
ORDER BY p.permission_key`, pq.Array(roleIDStrings(roleIDs)))
	if err != nil {
		return nil, fmt.Errorf("list role permissions: %w", err)
	}
	defer rows.Close()
	return scanPermissionSet(rows)
}

func (store *PostgresStore) ready() error {
	if store == nil || store.executor == nil {
		return ErrRBACStoreExecutorMissing
	}
	return nil
}

type permissionRows interface {
	Next() bool
	Scan(dest ...interface{}) error
	Err() error
}

func scanPermissionSet(rows permissionRows) (PermissionSet, error) {
	permissions := NewPermissionSet()
	for rows.Next() {
		var permission string
		if err := rows.Scan(&permission); err != nil {
			return nil, fmt.Errorf("scan permission: %w", err)
		}
		permissions[Permission(permission)] = struct{}{}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("read permissions: %w", err)
	}
	return permissions, nil
}

func roleIDStrings(roleIDs []identity.RoleID) []string {
	values := make([]string, 0, len(roleIDs))
	for _, roleID := range roleIDs {
		if roleID == "" {
			continue
		}
		values = append(values, string(roleID))
	}
	return values
}
