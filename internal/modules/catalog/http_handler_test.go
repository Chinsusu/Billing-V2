package catalog

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

func TestHTTPHandlerCreateProductUsesActorHeader(t *testing.T) {
	service := &fakeCatalogHTTPService{
		product: Product{
			ID:        "product_1",
			DisplayID: 10001,
			Type:      ProductTypeVPS,
			Name:      "VPS",
			Status:    ProductStatusDraft,
			CreatedBy: "actor_1",
		},
	}
	handler := registerCatalogTestHandler(service)

	request := httptest.NewRequest(http.MethodPost, "/admin/catalog/products", strings.NewReader(`{
		"product_type": "vps",
		"name": "VPS",
		"description": "Compute",
		"status": "draft",
		"display_order": 10
	}`))
	request.Header.Set(ActorHeader, " actor_1 ")
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d: %s", response.Code, response.Body.String())
	}
	if service.createProductCalls != 1 {
		t.Fatalf("expected create product once, got %d", service.createProductCalls)
	}
	if service.createProductInput.CreatedBy != "actor_1" {
		t.Fatalf("expected actor from header, got %q", service.createProductInput.CreatedBy)
	}
	if service.createProductInput.Type != ProductTypeVPS {
		t.Fatalf("expected VPS product type, got %q", service.createProductInput.Type)
	}
}

func TestHTTPHandlerAdminMiddlewareRunsBeforeService(t *testing.T) {
	service := &fakeCatalogHTTPService{}
	mux := http.NewServeMux()
	NewHTTPHandlerWithOptions(service, HTTPHandlerOptions{
		AdminMiddleware: func(next http.HandlerFunc) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusForbidden)
			}
		},
	}).RegisterRoutes(mux)

	request := httptest.NewRequest(http.MethodPost, "/admin/catalog/products", strings.NewReader(`{"product_type":"vps","name":"VPS"}`))
	response := httptest.NewRecorder()

	mux.ServeHTTP(response, request)

	if response.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d", response.Code)
	}
	if service.createProductCalls != 0 {
		t.Fatalf("expected service not to run, got %d calls", service.createProductCalls)
	}
}

func TestHTTPHandlerClientCatalogRequiresTenant(t *testing.T) {
	service := &fakeCatalogHTTPService{}
	handler := registerCatalogTestHandler(service)

	request := httptest.NewRequest(http.MethodGet, "/client/catalog", nil)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d: %s", response.Code, response.Body.String())
	}
	if service.listTenantCatalogCalls != 0 {
		t.Fatalf("expected no service call, got %d", service.listTenantCatalogCalls)
	}
	if !strings.Contains(response.Body.String(), "tenant.context_missing") {
		t.Fatalf("expected tenant validation error, got %s", response.Body.String())
	}
}

func TestHTTPHandlerClientCatalogUsesTenantContext(t *testing.T) {
	service := &fakeCatalogHTTPService{}
	handler := registerCatalogTestHandler(service)

	request := httptest.NewRequest(http.MethodGet, "/client/catalog", nil)
	request = request.WithContext(tenant.WithContext(request.Context(), tenant.NewContext("tenant_context")))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", response.Code, response.Body.String())
	}
	if service.tenantCatalogFilter.TenantID != tenant.ID("tenant_context") {
		t.Fatalf("expected tenant from context, got %q", service.tenantCatalogFilter.TenantID)
	}
}

func TestHTTPHandlerClientCatalogOmitsResellerCost(t *testing.T) {
	service := &fakeCatalogHTTPService{
		tenantCatalog: TenantCatalog{
			Products: []TenantProduct{{
				ID:              "tenant_product_1",
				DisplayID:       11001,
				TenantID:        "tenant_a",
				MasterProductID: "product_1",
				Status:          TenantProductStatusActive,
				CloneVersion:    1,
			}},
			Plans: []TenantPlan{{
				ID:                "tenant_plan_1",
				DisplayID:         12001,
				TenantID:          "tenant_a",
				TenantProductID:   "tenant_product_1",
				MasterPlanID:      "plan_1",
				SellingPriceMinor: 15000,
				ResellerCostMinor: 9000,
				Currency:          "USD",
				Visibility:        TenantPlanVisibilityPublic,
				Status:            TenantPlanStatusActive,
				CloneVersion:      1,
			}},
		},
	}
	handler := registerCatalogTestHandler(service)

	request := httptest.NewRequest(http.MethodGet, "/client/catalog?limit=10", nil)
	request.Header.Set(TenantHeader, " tenant_a ")
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", response.Code, response.Body.String())
	}
	if service.listTenantCatalogCalls != 1 {
		t.Fatalf("expected list tenant catalog once, got %d", service.listTenantCatalogCalls)
	}
	if service.tenantCatalogFilter.TenantID != tenant.ID("tenant_a") {
		t.Fatalf("expected tenant filter, got %q", service.tenantCatalogFilter.TenantID)
	}
	body := response.Body.String()
	if strings.Contains(body, "reseller_cost") {
		t.Fatalf("client response exposed reseller cost: %s", body)
	}
	if !strings.Contains(body, `"selling_price_minor":15000`) {
		t.Fatalf("expected selling price in client response, got %s", body)
	}
}

func registerCatalogTestHandler(service HTTPService) http.Handler {
	mux := http.NewServeMux()
	NewHTTPHandler(service).RegisterRoutes(mux)
	return mux
}

type fakeCatalogHTTPService struct {
	createProductCalls int
	createProductInput CreateProductInput
	product            Product

	tenantCatalog          TenantCatalog
	tenantCatalogFilter    TenantCatalogFilter
	listTenantCatalogCalls int
}

func (service *fakeCatalogHTTPService) CreateProduct(ctx context.Context, input CreateProductInput) (Product, error) {
	service.createProductCalls++
	service.createProductInput = input
	return service.product, nil
}

func (service *fakeCatalogHTTPService) CreatePlan(ctx context.Context, input CreatePlanInput) (Plan, error) {
	return Plan{}, nil
}

func (service *fakeCatalogHTTPService) CreateProviderSource(ctx context.Context, input CreateProviderSourceInput) (ProviderSource, error) {
	return ProviderSource{}, nil
}

func (service *fakeCatalogHTTPService) CreatePlanSource(ctx context.Context, input CreatePlanSourceInput) (PlanSource, error) {
	return PlanSource{}, nil
}

func (service *fakeCatalogHTTPService) CloneTenantProduct(ctx context.Context, input CreateTenantProductInput) (TenantProduct, error) {
	return TenantProduct{}, nil
}

func (service *fakeCatalogHTTPService) CloneTenantPlan(ctx context.Context, input CreateTenantPlanInput) (TenantPlan, error) {
	return TenantPlan{}, nil
}

func (service *fakeCatalogHTTPService) ListMasterPlans(ctx context.Context, filter MasterPlanFilter) ([]Plan, error) {
	return nil, nil
}

func (service *fakeCatalogHTTPService) ListTenantCatalog(ctx context.Context, filter TenantCatalogFilter) (TenantCatalog, error) {
	service.listTenantCatalogCalls++
	service.tenantCatalogFilter = filter
	return service.tenantCatalog, nil
}
