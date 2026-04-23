package catalog

import (
	"errors"
	"strings"
	"testing"
)

func TestBuildUpdateProductStatusQuery(t *testing.T) {
	query, args, err := buildUpdateProductStatusQuery(UpdateProductStatusInput{
		ID:     "product-1",
		Status: ProductStatusDisabled,
	})
	if err != nil {
		t.Fatalf("expected query: %v", err)
	}
	for _, clause := range []string{"UPDATE master_products", "status = $2", "WHERE product_id = $1", "RETURNING product_id"} {
		if !strings.Contains(query, clause) {
			t.Fatalf("expected %q in query: %s", clause, query)
		}
	}
	if len(args) != 2 || args[0] != ProductID("product-1") || args[1] != ProductStatusDisabled {
		t.Fatalf("unexpected args: %#v", args)
	}
}

func TestBuildUpdatePlanStatusQueryRejectsBadStatus(t *testing.T) {
	_, _, err := buildUpdatePlanStatusQuery(UpdatePlanStatusInput{
		ID:     "plan-1",
		Status: PlanStatus("bad"),
	})
	if !errors.Is(err, ErrPlanStatusInvalid) {
		t.Fatalf("expected plan status error, got %v", err)
	}
}

func TestBuildUpdateProviderSourceStatusQueryRequiresID(t *testing.T) {
	_, _, err := buildUpdateProviderSourceStatusQuery(UpdateProviderSourceStatusInput{
		Status: ProviderSourceStatusDisabled,
	})
	if !errors.Is(err, ErrSourceIDMissing) {
		t.Fatalf("expected source id error, got %v", err)
	}
}
