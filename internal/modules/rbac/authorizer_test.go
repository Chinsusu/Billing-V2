package rbac

import (
	"context"
	"errors"
	"testing"

	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

func TestStoreAuthorizerAllowsPlatformAdminForGlobalRequest(t *testing.T) {
	authorizer := NewStoreAuthorizer(nil)
	actor := identity.NewActor("admin_1", "platform", identity.ActorTypePlatformAdmin)

	err := authorizer.Check(context.Background(), CheckRequest{
		Actor:      actor,
		Permission: PermissionCatalogManage,
		Risk:       RiskHigh,
	})

	if err != nil {
		t.Fatalf("expected platform admin global access, got %v", err)
	}
}

func TestStoreAuthorizerChecksStorePermissionsForTenantActor(t *testing.T) {
	store := &fakeRBACStore{permissions: NewPermissionSet(PermissionCatalogView)}
	authorizer := NewStoreAuthorizer(store)
	actor := identity.NewActor("user_1", "tenant_a", identity.ActorTypeClient)
	ctx := tenant.WithContext(context.Background(), tenant.NewContext("tenant_a"))

	err := authorizer.Check(ctx, CheckRequest{
		Actor:      actor,
		Permission: PermissionCatalogView,
		Risk:       RiskLow,
	})

	if err != nil {
		t.Fatalf("expected permission, got %v", err)
	}
	if store.userID != "user_1" || store.tenantID != "tenant_a" {
		t.Fatalf("expected tenant/user lookup, got tenant=%q user=%q", store.tenantID, store.userID)
	}
}

func TestStoreAuthorizerDeniesMissingPermission(t *testing.T) {
	store := &fakeRBACStore{permissions: NewPermissionSet(PermissionCatalogView)}
	authorizer := NewStoreAuthorizer(store)
	actor := identity.NewActor("user_1", "tenant_a", identity.ActorTypeClient)
	ctx := tenant.WithContext(context.Background(), tenant.NewContext("tenant_a"))

	err := authorizer.Check(ctx, CheckRequest{
		Actor:      actor,
		Permission: PermissionCatalogManage,
		Risk:       RiskMedium,
	})

	if !errors.Is(err, ErrPermissionDenied) {
		t.Fatalf("expected permission denied, got %v", err)
	}
}

type fakeRBACStore struct {
	tenantID    tenant.ID
	userID      identity.UserID
	permissions PermissionSet
}

func (store *fakeRBACStore) ListRoleIDsForUser(ctx context.Context, tenantID tenant.ID, userID identity.UserID) ([]identity.RoleID, error) {
	return nil, nil
}

func (store *fakeRBACStore) ListPermissionsForUser(ctx context.Context, tenantID tenant.ID, userID identity.UserID) (PermissionSet, error) {
	store.tenantID = tenantID
	store.userID = userID
	return store.permissions, nil
}

func (store *fakeRBACStore) ListPermissionsForRoles(ctx context.Context, roleIDs []identity.RoleID) (PermissionSet, error) {
	return nil, nil
}
