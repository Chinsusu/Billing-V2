package tenant

import (
	"context"
	"errors"
)

var (
	ErrContextMissing         = errors.New("tenant context missing")
	ErrTenantIDMissing        = errors.New("tenant id missing")
	ErrTenantMismatch         = errors.New("tenant context mismatch")
	ErrAccessDenied           = errors.New("tenant access denied")
	ErrEmergencyReasonMissing = errors.New("emergency access reason missing")
)

type ID string

func (id ID) Empty() bool {
	return id == ""
}

type Type string

const (
	TypePlatform Type = "platform"
	TypeAdmin    Type = "admin"
	TypeReseller Type = "reseller"
)

type Context struct {
	Domain            string
	DomainTenantID    ID
	ActorTenantID     ID
	EffectiveTenantID ID
	TargetTenantID    ID
	IsPlatformAdmin   bool
	IsEmergencyAccess bool
	EmergencyReason   string
}

func NewContext(actorTenantID ID) Context {
	return Context{
		ActorTenantID:     actorTenantID,
		EffectiveTenantID: actorTenantID,
	}
}

func SystemContext(jobTenantID ID) Context {
	return Context{
		ActorTenantID:     jobTenantID,
		EffectiveTenantID: jobTenantID,
		TargetTenantID:    jobTenantID,
	}
}

func PlatformAdminContext(actorTenantID ID, targetTenantID ID, reason string) Context {
	return Context{
		ActorTenantID:     actorTenantID,
		EffectiveTenantID: targetTenantID,
		TargetTenantID:    targetTenantID,
		IsPlatformAdmin:   true,
		IsEmergencyAccess: targetTenantID != "" && actorTenantID != targetTenantID,
		EmergencyReason:   reason,
	}
}

func (tenantContext Context) Validate() error {
	if tenantContext.ActorTenantID.Empty() {
		return ErrTenantIDMissing
	}
	if tenantContext.EffectiveTenantID.Empty() {
		return ErrTenantIDMissing
	}
	if !tenantContext.DomainTenantID.Empty() && tenantContext.DomainTenantID != tenantContext.EffectiveTenantID {
		return ErrTenantMismatch
	}
	if tenantContext.IsPlatformAdmin {
		if tenantContext.IsEmergencyAccess && tenantContext.EmergencyReason == "" {
			return ErrEmergencyReasonMissing
		}
		return nil
	}
	if tenantContext.ActorTenantID != tenantContext.EffectiveTenantID {
		return ErrTenantMismatch
	}
	return nil
}

func (tenantContext Context) CanAccessTenant(targetTenantID ID) bool {
	if targetTenantID.Empty() {
		return false
	}
	if tenantContext.EffectiveTenantID == targetTenantID {
		return true
	}
	return tenantContext.IsPlatformAdmin &&
		tenantContext.IsEmergencyAccess &&
		tenantContext.TargetTenantID == targetTenantID &&
		tenantContext.EmergencyReason != ""
}

type contextKey struct{}

func WithContext(ctx context.Context, tenantContext Context) context.Context {
	return context.WithValue(ctx, contextKey{}, tenantContext)
}

func FromContext(ctx context.Context) (Context, bool) {
	tenantContext, ok := ctx.Value(contextKey{}).(Context)
	return tenantContext, ok
}

func RequireContext(ctx context.Context) (Context, error) {
	tenantContext, ok := FromContext(ctx)
	if !ok {
		return Context{}, ErrContextMissing
	}
	if err := tenantContext.Validate(); err != nil {
		return Context{}, err
	}
	return tenantContext, nil
}

func RequireAccess(ctx context.Context, targetTenantID ID) (Context, error) {
	tenantContext, err := RequireContext(ctx)
	if err != nil {
		return Context{}, err
	}
	if !tenantContext.CanAccessTenant(targetTenantID) {
		return Context{}, ErrAccessDenied
	}
	return tenantContext, nil
}
