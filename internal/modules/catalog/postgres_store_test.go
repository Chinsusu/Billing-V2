package catalog

import (
	"errors"
	"strings"
	"testing"

	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

func TestBuildListMasterPlansQueryAddsOptionalFilters(t *testing.T) {
	query, args, err := buildListMasterPlansQuery(MasterPlanFilter{
		ProductType: ProductTypeVPS,
		Status:      PlanStatusActive,
		Limit:       25,
	})
	if err != nil {
		t.Fatalf("expected query: %v", err)
	}
	if !hasSQLClause(query, "p.product_type = $1") {
		t.Fatalf("expected product type filter in query: %s", query)
	}
	if !hasSQLClause(query, "mp.status = $2") {
		t.Fatalf("expected status filter in query: %s", query)
	}
	if !hasSQLClause(query, "LIMIT $3") {
		t.Fatalf("expected limit placeholder in query: %s", query)
	}
	if len(args) != 3 || args[2] != 25 {
		t.Fatalf("unexpected args: %#v", args)
	}
}

func TestBuildListMasterPlansQueryDefaultsLimit(t *testing.T) {
	query, args, err := buildListMasterPlansQuery(MasterPlanFilter{})
	if err != nil {
		t.Fatalf("expected query: %v", err)
	}
	if !hasSQLClause(query, "LIMIT $1") {
		t.Fatalf("expected default limit placeholder in query: %s", query)
	}
	if len(args) != 1 || args[0] != defaultCatalogListLimit {
		t.Fatalf("unexpected args: %#v", args)
	}
}

func TestBuildListTenantPlansQueryRequiresTenant(t *testing.T) {
	_, _, err := buildListTenantPlansQuery(TenantCatalogFilter{})
	if !errors.Is(err, tenant.ErrTenantIDMissing) {
		t.Fatalf("expected tenant id error, got %v", err)
	}
}

func TestBuildListTenantPlansQueryAddsFiltersAndClampsLimit(t *testing.T) {
	query, args, err := buildListTenantPlansQuery(TenantCatalogFilter{
		TenantID:    tenant.ID("tenant-1"),
		ProductType: ProductTypeProxy,
		Visibility:  TenantPlanVisibilityPublic,
		Status:      TenantPlanStatusActive,
		Limit:       999,
	})
	if err != nil {
		t.Fatalf("expected query: %v", err)
	}
	for _, clause := range []string{"mp.product_type = $2", "tpl.visibility = $3", "tpl.status = $4", "LIMIT $5"} {
		if !hasSQLClause(query, clause) {
			t.Fatalf("expected %q in query: %s", clause, query)
		}
	}
	if len(args) != 5 || args[4] != maxCatalogListLimit {
		t.Fatalf("unexpected args: %#v", args)
	}
}

func TestValidateTenantCatalogFilterRejectsBadVisibility(t *testing.T) {
	err := validateTenantCatalogFilter(TenantCatalogFilter{
		TenantID:   tenant.ID("tenant-1"),
		Visibility: TenantPlanVisibility("bad"),
	})
	if !errors.Is(err, ErrTenantPlanVisibility) {
		t.Fatalf("expected visibility error, got %v", err)
	}
}

func hasSQLClause(query string, clause string) bool {
	return strings.Contains(query, clause)
}
