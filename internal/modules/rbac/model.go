package rbac

import (
	"context"
	"errors"

	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

var (
	ErrPermissionMissing = errors.New("permission missing")
	ErrPermissionDenied  = errors.New("permission denied")
)

type Permission string

const (
	PermissionTenantView            Permission = "tenant.view"
	PermissionTenantCreate          Permission = "tenant.create"
	PermissionTenantUpdate          Permission = "tenant.update"
	PermissionTenantDomainManage    Permission = "tenant.domain.manage"
	PermissionWalletView            Permission = "wallet.view"
	PermissionWalletTopupApprove    Permission = "wallet.topup.approve"
	PermissionWalletAdjustment      Permission = "wallet.adjustment.create"
	PermissionOrderView             Permission = "order.view"
	PermissionOrderCreate           Permission = "order.create"
	PermissionOrderManage           Permission = "order.manage"
	PermissionServiceView           Permission = "service.view"
	PermissionServiceReveal         Permission = "service.credential.reveal"
	PermissionServiceSuspend        Permission = "service.suspend"
	PermissionServiceUnsuspend      Permission = "service.unsuspend"
	PermissionServiceTerminate      Permission = "service.terminate"
	PermissionProvisioningJobView   Permission = "provisioning.job.view"
	PermissionProvisioningJobRetry  Permission = "provisioning.job.retry"
	PermissionManualReviewResolve   Permission = "provisioning.manual_review.resolve"
	PermissionProviderView          Permission = "provider.view"
	PermissionProviderManage        Permission = "provider.manage"
	PermissionCatalogView           Permission = "catalog.view"
	PermissionCatalogManage         Permission = "catalog.manage"
	PermissionAuditView             Permission = "audit.view"
	PermissionTenantEmergencyAccess Permission = "tenant.emergency_access"
)

type RiskLevel string

const (
	RiskLow      RiskLevel = "low"
	RiskMedium   RiskLevel = "medium"
	RiskHigh     RiskLevel = "high"
	RiskCritical RiskLevel = "critical"
)

type PermissionSet map[Permission]struct{}

func NewPermissionSet(permissions ...Permission) PermissionSet {
	set := make(PermissionSet, len(permissions))
	for _, permission := range permissions {
		if permission == "" {
			continue
		}
		set[permission] = struct{}{}
	}
	return set
}

func (set PermissionSet) Has(permission Permission) bool {
	_, ok := set[permission]
	return ok
}

type CheckRequest struct {
	Actor            identity.Actor
	Permission       Permission
	ResourceTenantID tenant.ID
	Risk             RiskLevel
	Reason           string
}

type Authorizer interface {
	Check(ctx context.Context, request CheckRequest) error
}

type StaticAuthorizer struct {
	Permissions PermissionSet
}

func (authorizer StaticAuthorizer) Check(ctx context.Context, request CheckRequest) error {
	if request.Permission == "" {
		return ErrPermissionMissing
	}
	if err := request.Actor.Validate(); err != nil {
		return err
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
	if !authorizer.Permissions.Has(request.Permission) {
		return ErrPermissionDenied
	}
	if request.Risk == RiskCritical && request.Reason == "" {
		return tenant.ErrEmergencyReasonMissing
	}
	return nil
}
