package order

import (
	"errors"
	"testing"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/audit"
	"github.com/Chinsusu/Billing-V2/internal/modules/catalog"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

func TestServiceLifecycleExpireAndGraceValidate(t *testing.T) {
	expire := TransitionServiceLifecycleInput{
		ID:            "service-1",
		TenantID:      tenant.ID("tenant-1"),
		ActorType:     audit.ActorTypeSystem,
		Action:        ServiceLifecycleActionExpire,
		FromStatus:    ServiceStatusActive,
		ToStatus:      ServiceStatusExpired,
		BillingStatus: BillingStatusOverdue,
	}
	if err := expire.Normalize().Validate(); err != nil {
		t.Fatalf("expected expire transition to validate: %v", err)
	}
	grace := TransitionServiceLifecycleInput{
		ID:               "service-1",
		TenantID:         tenant.ID("tenant-1"),
		ActorType:        audit.ActorTypeSystem,
		Action:           ServiceLifecycleActionGrace,
		FromStatus:       ServiceStatusExpired,
		ToStatus:         ServiceStatusSuspended,
		BillingStatus:    BillingStatusOverdue,
		SuspensionReason: SuspensionReasonExpiry,
	}
	if err := grace.Normalize().Validate(); err != nil {
		t.Fatalf("expected grace transition to validate: %v", err)
	}
}

func TestServiceLifecycleRejectsUnsafeManualSuspend(t *testing.T) {
	input := TransitionServiceLifecycleInput{
		ID:               "service-1",
		TenantID:         tenant.ID("tenant-1"),
		ActorID:          "admin-1",
		Action:           ServiceLifecycleActionSuspend,
		FromStatus:       ServiceStatusActive,
		ToStatus:         ServiceStatusSuspended,
		SuspensionReason: SuspensionReasonExpiry,
		Reason:           "manual action",
	}
	if err := input.Normalize().Validate(); !errors.Is(err, ErrSuspensionReasonInvalid) {
		t.Fatalf("expected expiry reason rejection, got %v", err)
	}
}

func TestServiceLifecycleRejectsTerminatedNoop(t *testing.T) {
	input := TransitionServiceLifecycleInput{
		ID:         "service-1",
		TenantID:   tenant.ID("tenant-1"),
		ActorID:    "admin-1",
		Action:     ServiceLifecycleActionTerminate,
		FromStatus: ServiceStatusTerminated,
		ToStatus:   ServiceStatusTerminated,
		Reason:     "duplicate request",
	}
	if err := input.Normalize().Validate(); !errors.Is(err, ErrServiceStatusTransitionInvalid) {
		t.Fatalf("expected terminated no-op rejection, got %v", err)
	}
}

func TestCalculateRenewedTermEndClampsCalendarMonth(t *testing.T) {
	termEnd := time.Date(2026, 1, 31, 8, 30, 0, 0, time.UTC)
	service := ServiceInstance{
		Status:  ServiceStatusActive,
		TermEnd: termEnd,
	}
	newTermEnd, err := CalculateRenewedTermEnd(service, ServiceRenewalCycle{
		Type:  catalog.BillingCycleCalendarMonth,
		Value: 1,
	})
	if err != nil {
		t.Fatalf("expected renewed term: %v", err)
	}
	expected := time.Date(2026, 2, 28, 8, 30, 0, 0, time.UTC)
	if !newTermEnd.Equal(expected) {
		t.Fatalf("expected %s, got %s", expected, newTermEnd)
	}
}

func TestCalculateRenewedTermEndRejectsManualSuspension(t *testing.T) {
	service := ServiceInstance{
		Status:           ServiceStatusSuspended,
		SuspensionReason: SuspensionReasonManualAdmin,
		TermEnd:          time.Date(2026, 1, 31, 0, 0, 0, 0, time.UTC),
	}
	_, err := CalculateRenewedTermEnd(service, ServiceRenewalCycle{
		Type:  catalog.BillingCycleMonth30Days,
		Value: 1,
	})
	if !errors.Is(err, ErrServiceStatusTransitionInvalid) {
		t.Fatalf("expected renewal rejection, got %v", err)
	}
}
