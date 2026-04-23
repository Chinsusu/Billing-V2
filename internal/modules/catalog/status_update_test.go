package catalog

import (
	"errors"
	"testing"
)

func TestUpdateProductStatusInputNormalizeValidate(t *testing.T) {
	input := UpdateProductStatusInput{
		ID:     " product-1 ",
		Status: " active ",
	}.Normalize()

	if input.ID != ProductID("product-1") {
		t.Fatalf("expected trimmed product id, got %q", input.ID)
	}
	if input.Status != ProductStatusActive {
		t.Fatalf("expected trimmed status, got %q", input.Status)
	}
	if err := input.Validate(); err != nil {
		t.Fatalf("expected valid product status input: %v", err)
	}

	input.Status = ProductStatus("bad")
	if err := input.Validate(); !errors.Is(err, ErrProductStatusInvalid) {
		t.Fatalf("expected product status error, got %v", err)
	}
}

func TestUpdatePlanStatusInputRequiresID(t *testing.T) {
	err := UpdatePlanStatusInput{Status: PlanStatusActive}.Normalize().Validate()
	if !errors.Is(err, ErrPlanIDMissing) {
		t.Fatalf("expected plan id error, got %v", err)
	}
}

func TestUpdateProviderSourceStatusInputRejectsBadStatus(t *testing.T) {
	err := UpdateProviderSourceStatusInput{
		ID:     "source-1",
		Status: ProviderSourceStatus("bad"),
	}.Normalize().Validate()
	if !errors.Is(err, ErrSourceStatusInvalid) {
		t.Fatalf("expected source status error, got %v", err)
	}
}
