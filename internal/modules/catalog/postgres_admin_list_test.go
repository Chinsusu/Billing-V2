package catalog

import (
	"errors"
	"strings"
	"testing"

	"github.com/Chinsusu/Billing-V2/internal/modules/provider"
)

func TestBuildListProductsQueryAddsFilters(t *testing.T) {
	query, args, err := buildListProductsQuery(ProductFilter{
		Type:   ProductTypeVPS,
		Status: ProductStatusActive,
		Limit:  25,
	})
	if err != nil {
		t.Fatalf("expected query: %v", err)
	}
	for _, clause := range []string{"product_type = $1", "status = $2", "LIMIT $3"} {
		if !strings.Contains(query, clause) {
			t.Fatalf("expected %q in query: %s", clause, query)
		}
	}
	if len(args) != 3 || args[2] != 25 {
		t.Fatalf("unexpected args: %#v", args)
	}
}

func TestBuildListProductsQueryRejectsBadStatus(t *testing.T) {
	_, _, err := buildListProductsQuery(ProductFilter{Status: ProductStatus("bad")})
	if !errors.Is(err, ErrProductStatusInvalid) {
		t.Fatalf("expected product status error, got %v", err)
	}
}

func TestBuildListProviderSourcesQueryAddsFilters(t *testing.T) {
	query, args, err := buildListProviderSourcesQuery(ProviderSourceFilter{
		Type:   provider.TypeManual,
		Status: ProviderSourceStatusActive,
		Limit:  10,
	})
	if err != nil {
		t.Fatalf("expected query: %v", err)
	}
	for _, clause := range []string{"source_type = $1", "status = $2", "LIMIT $3"} {
		if !strings.Contains(query, clause) {
			t.Fatalf("expected %q in query: %s", clause, query)
		}
	}
	if len(args) != 3 || args[2] != 10 {
		t.Fatalf("unexpected args: %#v", args)
	}
}

func TestBuildListProviderSourcesQueryRejectsBadType(t *testing.T) {
	_, _, err := buildListProviderSourcesQuery(ProviderSourceFilter{Type: provider.Type("bad")})
	if !errors.Is(err, ErrSourceTypeInvalid) {
		t.Fatalf("expected source type error, got %v", err)
	}
}
