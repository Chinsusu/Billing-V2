package catalog

import (
	"testing"

	"github.com/Chinsusu/Billing-V2/internal/modules/provider"
)

func TestProviderSourceReadinessMarksReadySource(t *testing.T) {
	record := ProviderSourceReadiness{
		ProductType:         ProductTypeVPS,
		PlanSourceDisplayID: 10003,
		PlanSourceStatus:    PlanSourceStatusActive,
		SourceDisplayID:     10002,
		SourceType:          provider.TypeHetzner,
		SourceStatus:        ProviderSourceStatusActive,
		capabilityProfile:   provider.DefaultCapabilityProfile(provider.TypeHetzner),
	}.withReadinessState()

	if record.State != ProviderSourceReadinessReady {
		t.Fatalf("expected ready, got %s: %s", record.State, record.Reason)
	}
}

func TestProviderSourceReadinessReportsMissingPlanSource(t *testing.T) {
	record := ProviderSourceReadiness{ProductType: ProductTypeVPS}.withReadinessState()

	if record.State != ProviderSourceReadinessMissingPlanSource {
		t.Fatalf("expected missing plan source, got %s", record.State)
	}
}

func TestProviderSourceReadinessReportsInactiveSource(t *testing.T) {
	record := ProviderSourceReadiness{
		ProductType:         ProductTypeVPS,
		PlanSourceDisplayID: 10003,
		PlanSourceStatus:    PlanSourceStatusActive,
		SourceDisplayID:     10002,
		SourceType:          provider.TypeHetzner,
		SourceStatus:        ProviderSourceStatusMaintenance,
		capabilityProfile:   provider.DefaultCapabilityProfile(provider.TypeHetzner),
	}.withReadinessState()

	if record.State != ProviderSourceReadinessInactiveSource {
		t.Fatalf("expected inactive source, got %s", record.State)
	}
}

func TestProviderSourceReadinessReportsFakeProviderOnly(t *testing.T) {
	record := ProviderSourceReadiness{
		ProductType:         ProductTypeVPS,
		PlanSourceDisplayID: 10003,
		PlanSourceStatus:    PlanSourceStatusActive,
		SourceDisplayID:     10002,
		SourceType:          provider.TypeManual,
		SourceStatus:        ProviderSourceStatusActive,
		capabilityProfile:   provider.DefaultCapabilityProfile(provider.TypeManual),
	}.withReadinessState()

	if record.State != ProviderSourceReadinessFakeProviderOnly {
		t.Fatalf("expected fake-provider-only, got %s", record.State)
	}
}

func TestProviderSourceReadinessReportsUnsupportedCapability(t *testing.T) {
	record := ProviderSourceReadiness{
		ProductType:         ProductTypeProxy,
		PlanSourceDisplayID: 10003,
		PlanSourceStatus:    PlanSourceStatusActive,
		SourceDisplayID:     10002,
		SourceType:          provider.TypeHetzner,
		SourceStatus:        ProviderSourceStatusActive,
		capabilityProfile:   provider.DefaultCapabilityProfile(provider.TypeHetzner),
	}.withReadinessState()

	if record.State != ProviderSourceReadinessUnsupportedCapability {
		t.Fatalf("expected unsupported capability, got %s", record.State)
	}
}
