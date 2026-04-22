package catalog

import (
	"context"

	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

type MasterPlanFilter struct {
	ProductType ProductType
	Status      PlanStatus
	Limit       int
}

type TenantCatalogFilter struct {
	TenantID    tenant.ID
	ProductType ProductType
	Visibility  TenantPlanVisibility
	Status      TenantPlanStatus
	Limit       int
}

type TenantCatalog struct {
	Products []TenantProduct
	Plans    []TenantPlan
}

type Store interface {
	CreateProduct(ctx context.Context, input CreateProductInput) (Product, error)
	CreatePlan(ctx context.Context, input CreatePlanInput) (Plan, error)
	CreateProviderSource(ctx context.Context, input CreateProviderSourceInput) (ProviderSource, error)
	CreatePlanSource(ctx context.Context, input CreatePlanSourceInput) (PlanSource, error)
	CreateTenantProduct(ctx context.Context, input CreateTenantProductInput) (TenantProduct, error)
	CreateTenantPlan(ctx context.Context, input CreateTenantPlanInput) (TenantPlan, error)
	ListMasterPlans(ctx context.Context, filter MasterPlanFilter) ([]Plan, error)
	ListTenantCatalog(ctx context.Context, filter TenantCatalogFilter) (TenantCatalog, error)
}
