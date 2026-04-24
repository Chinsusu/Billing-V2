package jobs

import (
	"errors"
	"testing"

	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

func TestNormalizeJobFilterTrimsAndClampsLimit(t *testing.T) {
	filter := normalizeFilter(Filter{
		TenantID:      " tenant_1 ",
		Type:          " provider.provision ",
		Status:        " queued ",
		ReferenceType: " order ",
		ReferenceID:   " order_1 ",
		SourceID:      " source_1 ",
		Limit:         999,
	})

	if filter.TenantID != tenant.ID("tenant_1") ||
		filter.Type != Type("provider.provision") ||
		filter.Status != StatusQueued ||
		filter.ReferenceType != ReferenceType("order") ||
		filter.ReferenceID != ReferenceID("order_1") ||
		filter.SourceID != SourceID("source_1") ||
		filter.Limit != maxJobListLimit {
		t.Fatalf("unexpected normalized filter: %+v", filter)
	}
}

func TestValidateJobFilterRequiresTenant(t *testing.T) {
	err := validateFilter(Filter{Status: StatusQueued, Limit: 10})

	if !errors.Is(err, tenant.ErrTenantIDMissing) {
		t.Fatalf("expected tenant error, got %v", err)
	}
}

func TestValidateJobFilterRejectsBadStatus(t *testing.T) {
	err := validateFilter(Filter{TenantID: "tenant_1", Status: "lost", Limit: 10})

	if !errors.Is(err, ErrStatusInvalid) {
		t.Fatalf("expected status error, got %v", err)
	}
}

func TestValidateJobLookupRequiresIDAndTenant(t *testing.T) {
	if err := validateLookup(Lookup{TenantID: "tenant_1"}); !errors.Is(err, ErrJobIDMissing) {
		t.Fatalf("expected job id error, got %v", err)
	}
	if err := validateLookup(Lookup{ID: "job_1"}); !errors.Is(err, tenant.ErrTenantIDMissing) {
		t.Fatalf("expected tenant error, got %v", err)
	}
}

func TestNormalizeAttemptFilterTrimsAndDefaultsLimit(t *testing.T) {
	filter := normalizeAttemptFilter(AttemptFilter{JobID: " job_1 ", TenantID: " tenant_1 "})

	if filter.JobID != ID("job_1") || filter.TenantID != tenant.ID("tenant_1") || filter.Limit != defaultJobListLimit {
		t.Fatalf("unexpected attempt filter: %+v", filter)
	}
}

func TestValidateAttemptFilterRequiresJobAndTenant(t *testing.T) {
	if err := validateAttemptFilter(AttemptFilter{TenantID: "tenant_1"}); !errors.Is(err, ErrJobIDMissing) {
		t.Fatalf("expected job id error, got %v", err)
	}
	if err := validateAttemptFilter(AttemptFilter{JobID: "job_1"}); !errors.Is(err, tenant.ErrTenantIDMissing) {
		t.Fatalf("expected tenant error, got %v", err)
	}
}
