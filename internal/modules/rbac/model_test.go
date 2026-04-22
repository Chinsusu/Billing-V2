package rbac

import (
	"context"
	"errors"
	"testing"

	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

func TestStaticAuthorizerRequiresTenantContext(t *testing.T) {
	authorizer := StaticAuthorizer{Permissions: NewPermissionSet(PermissionOrderView)}
	actor := identity.NewActor("user_1", "tenant_a", identity.ActorTypeClient)

	err := authorizer.Check(context.Background(), CheckRequest{
		Actor:      actor,
		Permission: PermissionOrderView,
	})
	if !errors.Is(err, tenant.ErrContextMissing) {
		t.Fatalf("expected tenant context missing, got %v", err)
	}
}

func TestStaticAuthorizerAllowsTenantScopedPermission(t *testing.T) {
	authorizer := StaticAuthorizer{Permissions: NewPermissionSet(PermissionOrderView)}
	actor := identity.NewActor("user_1", "tenant_a", identity.ActorTypeClient)
	ctx := tenant.WithContext(context.Background(), tenant.NewContext("tenant_a"))

	err := authorizer.Check(ctx, CheckRequest{
		Actor:            actor,
		Permission:       PermissionOrderView,
		ResourceTenantID: "tenant_a",
	})
	if err != nil {
		t.Fatalf("expected permission, got %v", err)
	}
}

func TestStaticAuthorizerRejectsMissingPermission(t *testing.T) {
	authorizer := StaticAuthorizer{Permissions: NewPermissionSet(PermissionOrderView)}
	actor := identity.NewActor("user_1", "tenant_a", identity.ActorTypeClient)
	ctx := tenant.WithContext(context.Background(), tenant.NewContext("tenant_a"))

	err := authorizer.Check(ctx, CheckRequest{
		Actor:            actor,
		Permission:       PermissionWalletView,
		ResourceTenantID: "tenant_a",
	})
	if !errors.Is(err, ErrPermissionDenied) {
		t.Fatalf("expected permission denied, got %v", err)
	}
}

func TestStaticAuthorizerRejectsCrossTenantResource(t *testing.T) {
	authorizer := StaticAuthorizer{Permissions: NewPermissionSet(PermissionServiceView)}
	actor := identity.NewActor("user_1", "tenant_a", identity.ActorTypeClient)
	ctx := tenant.WithContext(context.Background(), tenant.NewContext("tenant_a"))

	err := authorizer.Check(ctx, CheckRequest{
		Actor:            actor,
		Permission:       PermissionServiceView,
		ResourceTenantID: "tenant_b",
	})
	if !errors.Is(err, tenant.ErrAccessDenied) {
		t.Fatalf("expected tenant access denied, got %v", err)
	}
}

func TestStaticAuthorizerAllowsPlatformEmergencyAccess(t *testing.T) {
	authorizer := StaticAuthorizer{Permissions: NewPermissionSet(PermissionTenantEmergencyAccess)}
	actor := identity.NewActor("admin_1", "platform", identity.ActorTypePlatformAdmin)
	ctx := tenant.WithContext(context.Background(), tenant.PlatformAdminContext("platform", "tenant_b", "support ticket INC-42"))

	err := authorizer.Check(ctx, CheckRequest{
		Actor:            actor,
		Permission:       PermissionTenantEmergencyAccess,
		ResourceTenantID: "tenant_b",
		Risk:             RiskCritical,
		Reason:           "support ticket INC-42",
	})
	if err != nil {
		t.Fatalf("expected emergency permission, got %v", err)
	}
}
