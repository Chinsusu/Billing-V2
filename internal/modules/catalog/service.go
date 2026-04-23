package catalog

import (
	"context"

	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

type Service struct {
	store Store
}

func NewService(store Store) *Service {
	return &Service{store: store}
}

func (service *Service) CreateProduct(ctx context.Context, input CreateProductInput) (Product, error) {
	if err := service.ready(); err != nil {
		return Product{}, err
	}
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return Product{}, err
	}
	return service.store.CreateProduct(ctx, input)
}

func (service *Service) CreatePlan(ctx context.Context, input CreatePlanInput) (Plan, error) {
	if err := service.ready(); err != nil {
		return Plan{}, err
	}
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return Plan{}, err
	}
	return service.store.CreatePlan(ctx, input)
}

func (service *Service) CreateProviderSource(ctx context.Context, input CreateProviderSourceInput) (ProviderSource, error) {
	if err := service.ready(); err != nil {
		return ProviderSource{}, err
	}
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return ProviderSource{}, err
	}
	return service.store.CreateProviderSource(ctx, input)
}

func (service *Service) CreatePlanSource(ctx context.Context, input CreatePlanSourceInput) (PlanSource, error) {
	if err := service.ready(); err != nil {
		return PlanSource{}, err
	}
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return PlanSource{}, err
	}
	return service.store.CreatePlanSource(ctx, input)
}

func (service *Service) CloneTenantProduct(ctx context.Context, input CreateTenantProductInput) (TenantProduct, error) {
	if err := service.ready(); err != nil {
		return TenantProduct{}, err
	}
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return TenantProduct{}, err
	}
	return service.store.CreateTenantProduct(ctx, input)
}

func (service *Service) CloneTenantPlan(ctx context.Context, input CreateTenantPlanInput) (TenantPlan, error) {
	if err := service.ready(); err != nil {
		return TenantPlan{}, err
	}
	input = input.Normalize()
	input = applyTenantPlanMarginGuard(input)
	if err := input.Validate(); err != nil {
		return TenantPlan{}, err
	}
	return service.store.CreateTenantPlan(ctx, input)
}

func (service *Service) ListMasterPlans(ctx context.Context, filter MasterPlanFilter) ([]Plan, error) {
	if err := service.ready(); err != nil {
		return nil, err
	}
	return service.store.ListMasterPlans(ctx, filter)
}

func (service *Service) ListTenantCatalog(ctx context.Context, filter TenantCatalogFilter) (TenantCatalog, error) {
	if err := service.ready(); err != nil {
		return TenantCatalog{}, err
	}
	if filter.TenantID.Empty() {
		return TenantCatalog{}, tenant.ErrTenantIDMissing
	}
	return service.store.ListTenantCatalog(ctx, filter)
}

func (service *Service) ready() error {
	if service == nil || service.store == nil {
		return ErrCatalogServiceStoreMissing
	}
	return nil
}

func applyTenantPlanMarginGuard(input CreateTenantPlanInput) CreateTenantPlanInput {
	if input.SellingPriceMinor < input.ResellerCostMinor && input.Status == TenantPlanStatusActive {
		input.Status = TenantPlanStatusMarginRisk
	}
	return input
}
