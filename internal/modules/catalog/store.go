package catalog

import (
	"context"

	"github.com/Chinsusu/Billing-V2/internal/modules/provider"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

type ProductFilter struct {
	Type   ProductType
	Status ProductStatus
	Limit  int
}

type MasterPlanFilter struct {
	ProductType ProductType
	Status      PlanStatus
	Limit       int
}

type ProviderSourceFilter struct {
	Type   provider.Type
	Status ProviderSourceStatus
	Limit  int
}

type ProviderSourceReadinessFilter struct {
	ProductType ProductType
	PlanStatus  PlanStatus
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
	UpdateProductStatus(ctx context.Context, input UpdateProductStatusInput) (Product, error)
	UpdatePlanStatus(ctx context.Context, input UpdatePlanStatusInput) (Plan, error)
	UpdateProviderSourceStatus(ctx context.Context, input UpdateProviderSourceStatusInput) (ProviderSource, error)
	ListProducts(ctx context.Context, filter ProductFilter) ([]Product, error)
	ListMasterPlans(ctx context.Context, filter MasterPlanFilter) ([]Plan, error)
	ListProviderSources(ctx context.Context, filter ProviderSourceFilter) ([]ProviderSource, error)
	ListProviderSourceReadiness(ctx context.Context, filter ProviderSourceReadinessFilter) ([]ProviderSourceReadiness, error)
	ListTenantCatalog(ctx context.Context, filter TenantCatalogFilter) (TenantCatalog, error)
}
