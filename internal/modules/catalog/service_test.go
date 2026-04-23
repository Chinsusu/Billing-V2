package catalog

import (
	"context"
	"errors"
	"testing"

	"github.com/Chinsusu/Billing-V2/internal/modules/provider"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

type fakeCatalogStore struct {
	createProductInput        CreateProductInput
	createPlanInput           CreatePlanInput
	createProviderSourceInput CreateProviderSourceInput
	createPlanSourceInput     CreatePlanSourceInput
	createTenantProductInput  CreateTenantProductInput
	createTenantPlanInput     CreateTenantPlanInput
	listProductsFilter        ProductFilter
	listProviderSourcesFilter ProviderSourceFilter
	listTenantCatalogCalled   bool
}

func (store *fakeCatalogStore) CreateProduct(_ context.Context, input CreateProductInput) (Product, error) {
	store.createProductInput = input
	return Product{Name: input.Name, Status: input.Status, CreatedBy: input.CreatedBy}, nil
}

func (store *fakeCatalogStore) CreatePlan(_ context.Context, input CreatePlanInput) (Plan, error) {
	store.createPlanInput = input
	return Plan{Code: input.Code, Currency: input.Currency, Version: input.Version}, nil
}

func (store *fakeCatalogStore) CreateProviderSource(_ context.Context, input CreateProviderSourceInput) (ProviderSource, error) {
	store.createProviderSourceInput = input
	return ProviderSource{Name: input.Name, Status: input.Status, RiskLevel: input.RiskLevel}, nil
}

func (store *fakeCatalogStore) CreatePlanSource(_ context.Context, input CreatePlanSourceInput) (PlanSource, error) {
	store.createPlanSourceInput = input
	return PlanSource{Priority: input.Priority, Status: input.Status}, nil
}

func (store *fakeCatalogStore) CreateTenantProduct(_ context.Context, input CreateTenantProductInput) (TenantProduct, error) {
	store.createTenantProductInput = input
	return TenantProduct{TenantID: input.TenantID, Status: input.Status}, nil
}

func (store *fakeCatalogStore) CreateTenantPlan(_ context.Context, input CreateTenantPlanInput) (TenantPlan, error) {
	store.createTenantPlanInput = input
	return TenantPlan{TenantID: input.TenantID, Status: input.Status, Currency: input.Currency}, nil
}

func (store *fakeCatalogStore) ListMasterPlans(_ context.Context, _ MasterPlanFilter) ([]Plan, error) {
	return []Plan{{ID: PlanID("plan-1")}}, nil
}

func (store *fakeCatalogStore) ListProducts(_ context.Context, filter ProductFilter) ([]Product, error) {
	store.listProductsFilter = filter
	return []Product{{ID: ProductID("product-1")}}, nil
}

func (store *fakeCatalogStore) ListProviderSources(_ context.Context, filter ProviderSourceFilter) ([]ProviderSource, error) {
	store.listProviderSourcesFilter = filter
	return []ProviderSource{{ID: ProviderSourceID("source-1")}}, nil
}

func (store *fakeCatalogStore) ListTenantCatalog(_ context.Context, _ TenantCatalogFilter) (TenantCatalog, error) {
	store.listTenantCatalogCalled = true
	return TenantCatalog{}, nil
}

func TestServiceRejectsMissingStore(t *testing.T) {
	_, err := NewService(nil).CreateProduct(context.Background(), CreateProductInput{})
	if !errors.Is(err, ErrCatalogServiceStoreMissing) {
		t.Fatalf("expected missing store error, got %v", err)
	}
}

func TestCreateProductNormalizesAndValidatesBeforeStore(t *testing.T) {
	store := &fakeCatalogStore{}
	service := NewService(store)

	_, err := service.CreateProduct(context.Background(), CreateProductInput{
		Type:      ProductTypeVPS,
		Name:      "  VPS SG  ",
		CreatedBy: " admin-1 ",
	})
	if err != nil {
		t.Fatalf("expected product create: %v", err)
	}
	if store.createProductInput.Name != "VPS SG" {
		t.Fatalf("expected trimmed product name, got %q", store.createProductInput.Name)
	}
	if store.createProductInput.Status != ProductStatusDraft {
		t.Fatalf("expected draft status, got %q", store.createProductInput.Status)
	}
	if store.createProductInput.CreatedBy != "admin-1" {
		t.Fatalf("expected trimmed creator, got %q", store.createProductInput.CreatedBy)
	}
}

func TestCreateProviderSourceAppliesDefaultsBeforeStore(t *testing.T) {
	store := &fakeCatalogStore{}
	service := NewService(store)

	_, err := service.CreateProviderSource(context.Background(), CreateProviderSourceInput{
		Type:          provider.TypeProxmox,
		Name:          " Proxmox SG ",
		InventoryMode: InventoryModeProviderLive,
	})
	if err != nil {
		t.Fatalf("expected source create: %v", err)
	}
	if store.createProviderSourceInput.Status != ProviderSourceStatusDisabled {
		t.Fatalf("expected default disabled status, got %q", store.createProviderSourceInput.Status)
	}
	if store.createProviderSourceInput.RiskLevel != RiskLevelMedium {
		t.Fatalf("expected default medium risk, got %q", store.createProviderSourceInput.RiskLevel)
	}
	if !store.createProviderSourceInput.CapabilityProfile.SupportsAutoProvision {
		t.Fatal("expected default provider capability profile")
	}
}

func TestCloneTenantPlanAppliesMarginGuard(t *testing.T) {
	store := &fakeCatalogStore{}
	service := NewService(store)

	_, err := service.CloneTenantPlan(context.Background(), CreateTenantPlanInput{
		TenantID:          tenant.ID("tenant-1"),
		TenantProductID:   TenantProductID("tenant-product-1"),
		MasterPlanID:      PlanID("plan-1"),
		SellingPriceMinor: 700,
		ResellerCostMinor: 900,
		Currency:          "usd",
		Status:            TenantPlanStatusActive,
	})
	if err != nil {
		t.Fatalf("expected tenant plan clone: %v", err)
	}
	if store.createTenantPlanInput.Status != TenantPlanStatusMarginRisk {
		t.Fatalf("expected margin risk status, got %q", store.createTenantPlanInput.Status)
	}
	if store.createTenantPlanInput.Currency != "USD" {
		t.Fatalf("expected normalized currency, got %q", store.createTenantPlanInput.Currency)
	}
}

func TestListTenantCatalogRequiresTenantBeforeStore(t *testing.T) {
	store := &fakeCatalogStore{}
	service := NewService(store)

	_, err := service.ListTenantCatalog(context.Background(), TenantCatalogFilter{})
	if !errors.Is(err, tenant.ErrTenantIDMissing) {
		t.Fatalf("expected tenant id error, got %v", err)
	}
	if store.listTenantCatalogCalled {
		t.Fatal("store should not be called without tenant scope")
	}
}

func TestListProductsDelegatesToStore(t *testing.T) {
	store := &fakeCatalogStore{}
	service := NewService(store)

	_, err := service.ListProducts(context.Background(), ProductFilter{Type: ProductTypeVPS, Status: ProductStatusActive, Limit: 25})
	if err != nil {
		t.Fatalf("expected list products: %v", err)
	}
	if store.listProductsFilter.Type != ProductTypeVPS || store.listProductsFilter.Status != ProductStatusActive || store.listProductsFilter.Limit != 25 {
		t.Fatalf("unexpected product filter: %+v", store.listProductsFilter)
	}
}

func TestListProviderSourcesDelegatesToStore(t *testing.T) {
	store := &fakeCatalogStore{}
	service := NewService(store)

	_, err := service.ListProviderSources(context.Background(), ProviderSourceFilter{Type: provider.TypeManual, Status: ProviderSourceStatusActive, Limit: 10})
	if err != nil {
		t.Fatalf("expected list provider sources: %v", err)
	}
	if store.listProviderSourcesFilter.Type != provider.TypeManual || store.listProviderSourcesFilter.Status != ProviderSourceStatusActive || store.listProviderSourcesFilter.Limit != 10 {
		t.Fatalf("unexpected provider source filter: %+v", store.listProviderSourcesFilter)
	}
}
