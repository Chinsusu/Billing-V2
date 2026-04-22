package tenant

import (
	"context"
	"errors"
	"testing"
)

func TestRequireContextRejectsMissingTenantContext(t *testing.T) {
	_, err := RequireContext(context.Background())
	if !errors.Is(err, ErrContextMissing) {
		t.Fatalf("expected tenant context missing, got %v", err)
	}
}

func TestRequireContextRejectsActorTenantMismatch(t *testing.T) {
	ctx := WithContext(context.Background(), Context{
		ActorTenantID:     "tenant_a",
		EffectiveTenantID: "tenant_b",
	})

	_, err := RequireContext(ctx)
	if !errors.Is(err, ErrTenantMismatch) {
		t.Fatalf("expected tenant mismatch, got %v", err)
	}
}

func TestRequireAccessAllowsCurrentTenant(t *testing.T) {
	ctx := WithContext(context.Background(), NewContext("tenant_a"))

	tenantContext, err := RequireAccess(ctx, "tenant_a")
	if err != nil {
		t.Fatalf("expected access, got %v", err)
	}
	if tenantContext.EffectiveTenantID != "tenant_a" {
		t.Fatalf("expected effective tenant, got %q", tenantContext.EffectiveTenantID)
	}
}

func TestRequireAccessRejectsOtherTenant(t *testing.T) {
	ctx := WithContext(context.Background(), NewContext("tenant_a"))

	_, err := RequireAccess(ctx, "tenant_b")
	if !errors.Is(err, ErrAccessDenied) {
		t.Fatalf("expected access denied, got %v", err)
	}
}

func TestPlatformAdminEmergencyRequiresReason(t *testing.T) {
	ctx := WithContext(context.Background(), PlatformAdminContext("platform", "tenant_b", ""))

	_, err := RequireContext(ctx)
	if !errors.Is(err, ErrEmergencyReasonMissing) {
		t.Fatalf("expected emergency reason error, got %v", err)
	}
}

func TestPlatformAdminEmergencyCanAccessTargetTenant(t *testing.T) {
	ctx := WithContext(context.Background(), PlatformAdminContext("platform", "tenant_b", "support ticket INC-42"))

	_, err := RequireAccess(ctx, "tenant_b")
	if err != nil {
		t.Fatalf("expected emergency tenant access, got %v", err)
	}
}
