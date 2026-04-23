package rbac

import (
	"context"

	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

type StoreAuthorizer struct {
	Store Store
}

func NewStoreAuthorizer(store Store) StoreAuthorizer {
	return StoreAuthorizer{Store: store}
}

func (authorizer StoreAuthorizer) Check(ctx context.Context, request CheckRequest) error {
	if request.Permission == "" {
		return ErrPermissionMissing
	}
	if err := request.Actor.Validate(); err != nil {
		return err
	}
	if request.Actor.IsSystem {
		return nil
	}
	if request.Actor.IsPlatformAdmin && request.ResourceTenantID.Empty() {
		return nil
	}
	tenantContext, err := tenant.RequireContext(ctx)
	if err != nil {
		return err
	}
	if request.Actor.TenantID != tenantContext.ActorTenantID {
		return tenant.ErrTenantMismatch
	}
	if request.ResourceTenantID.Empty() {
		request.ResourceTenantID = tenantContext.EffectiveTenantID
	}
	if !tenantContext.CanAccessTenant(request.ResourceTenantID) {
		return tenant.ErrAccessDenied
	}
	if request.Risk == RiskCritical && request.Reason == "" {
		return tenant.ErrEmergencyReasonMissing
	}
	if authorizer.Store == nil {
		return ErrRBACStoreExecutorMissing
	}
	permissions, err := authorizer.Store.ListPermissionsForUser(ctx, request.Actor.TenantID, request.Actor.ID)
	if err != nil {
		return err
	}
	if !permissions.Has(request.Permission) {
		return ErrPermissionDenied
	}
	return nil
}
