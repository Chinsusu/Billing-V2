package catalog

import (
	"errors"
	"strings"
	"testing"
)

func TestBuildListProviderSourceReadinessQueryDefaultsToActivePlans(t *testing.T) {
	query, args, err := buildListProviderSourceReadinessQuery(ProviderSourceReadinessFilter{Limit: 25})
	if err != nil {
		t.Fatalf("expected query: %v", err)
	}
	for _, clause := range []string{"WITH ranked_plan_sources", "LEFT JOIN ranked_plan_sources", "mp.status = $1", "LIMIT $2"} {
		if !strings.Contains(query, clause) {
			t.Fatalf("expected %q in query: %s", clause, query)
		}
	}
	if len(args) != 2 || args[0] != PlanStatusActive || args[1] != 25 {
		t.Fatalf("unexpected args: %#v", args)
	}
}

func TestBuildListProviderSourceReadinessQueryAddsProductFilter(t *testing.T) {
	query, args, err := buildListProviderSourceReadinessQuery(ProviderSourceReadinessFilter{
		ProductType: ProductTypeProxy,
		PlanStatus:  PlanStatusActive,
		Limit:       10,
	})
	if err != nil {
		t.Fatalf("expected query: %v", err)
	}
	for _, clause := range []string{"mp.status = $1", "product.product_type = $2", "LIMIT $3"} {
		if !strings.Contains(query, clause) {
			t.Fatalf("expected %q in query: %s", clause, query)
		}
	}
	if len(args) != 3 || args[0] != PlanStatusActive || args[1] != ProductTypeProxy || args[2] != 10 {
		t.Fatalf("unexpected args: %#v", args)
	}
}

func TestBuildListProviderSourceReadinessQueryRejectsBadFilters(t *testing.T) {
	_, _, err := buildListProviderSourceReadinessQuery(ProviderSourceReadinessFilter{ProductType: ProductType("bad")})
	if !errors.Is(err, ErrProductTypeInvalid) {
		t.Fatalf("expected product type error, got %v", err)
	}

	_, _, err = buildListProviderSourceReadinessQuery(ProviderSourceReadinessFilter{PlanStatus: PlanStatus("bad")})
	if !errors.Is(err, ErrPlanStatusInvalid) {
		t.Fatalf("expected plan status error, got %v", err)
	}
}
