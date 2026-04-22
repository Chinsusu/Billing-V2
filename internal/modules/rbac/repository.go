package rbac

import (
	"context"
	"errors"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

var ErrRoleIDMissing = errors.New("role id missing")

type Role struct {
	ID        identity.RoleID
	DisplayID int64
	TenantID  tenant.ID
	Key       string
	Name      string
	IsSystem  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type PermissionRecord struct {
	ID         string
	Permission Permission
	Module     string
	Risk       RiskLevel
}

type Store interface {
	ListRoleIDsForUser(ctx context.Context, tenantID tenant.ID, userID identity.UserID) ([]identity.RoleID, error)
	ListPermissionsForUser(ctx context.Context, tenantID tenant.ID, userID identity.UserID) (PermissionSet, error)
	ListPermissionsForRoles(ctx context.Context, roleIDs []identity.RoleID) (PermissionSet, error)
}

func validateTenantUserLookup(tenantID tenant.ID, userID identity.UserID) error {
	if tenantID.Empty() {
		return tenant.ErrTenantIDMissing
	}
	if userID == "" {
		return identity.ErrUserIDMissing
	}
	return nil
}
