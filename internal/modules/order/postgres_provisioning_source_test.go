package order

import (
	"errors"
	"strings"
	"testing"

	"github.com/Chinsusu/Billing-V2/internal/modules/catalog"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

func TestResolveOrderProvisioningSourceArgsNormalizeAndValidate(t *testing.T) {
	args, err := resolveOrderProvisioningSourceArgs(ResolveOrderProvisioningSourceInput{
		TenantID:     tenant.ID(" tenant-1 "),
		TenantPlanID: catalog.TenantPlanID(" tenant-plan-1 "),
	})
	if err != nil {
		t.Fatalf("expected source lookup args: %v", err)
	}
	if len(args) != 2 ||
		args[0] != catalog.TenantPlanID("tenant-plan-1") ||
		args[1] != tenant.ID("tenant-1") {
		t.Fatalf("unexpected source lookup args: %#v", args)
	}
}

func TestResolveOrderProvisioningSourceArgsRejectsMissingTenantPlan(t *testing.T) {
	_, err := resolveOrderProvisioningSourceArgs(ResolveOrderProvisioningSourceInput{
		TenantID: tenant.ID("tenant-1"),
	})
	if !errors.Is(err, ErrTenantPlanIDMissing) {
		t.Fatalf("expected tenant plan error, got %v", err)
	}
}

func TestResolveOrderProvisioningSourceSQLUsesActivePlanSource(t *testing.T) {
	for _, clause := range []string{
		"FROM tenant_plans tenant_plan",
		"JOIN plan_sources plan_source ON plan_source.plan_id = tenant_plan.master_plan_id",
		"JOIN provider_sources source ON source.source_id = plan_source.source_id",
		"tenant_plan.tenant_plan_id = $1",
		"tenant_plan.tenant_id = $2",
		"plan_source.status = 'active'",
		"source.status = 'active'",
		"ORDER BY plan_source.priority ASC",
		"LIMIT 1",
	} {
		if !strings.Contains(resolveOrderProvisioningSourceSQL, clause) {
			t.Fatalf("expected %q in source lookup SQL: %s", clause, resolveOrderProvisioningSourceSQL)
		}
	}
}
