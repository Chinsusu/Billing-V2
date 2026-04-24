package catalog

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Chinsusu/Billing-V2/internal/modules/provider"
)

func TestHTTPHandlerListProductsUsesQueryFilters(t *testing.T) {
	service := &fakeCatalogHTTPService{
		products: []Product{{ID: "product_1", DisplayID: 10001, Type: ProductTypeVPS, Name: "VPS", Status: ProductStatusActive}},
	}
	handler := registerCatalogTestHandler(service)

	request := httptest.NewRequest(http.MethodGet, "/admin/catalog/products?product_type=vps&status=active&limit=15", nil)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", response.Code, response.Body.String())
	}
	if service.listProductCalls != 1 {
		t.Fatalf("expected list products once, got %d", service.listProductCalls)
	}
	if service.productFilter.Type != ProductTypeVPS || service.productFilter.Status != ProductStatusActive || service.productFilter.Limit != 15 {
		t.Fatalf("unexpected product filter: %+v", service.productFilter)
	}
	if !strings.Contains(response.Body.String(), `"display_id":10001`) {
		t.Fatalf("expected product response, got %s", response.Body.String())
	}
}

func TestHTTPHandlerListProviderSourcesUsesQueryFilters(t *testing.T) {
	service := &fakeCatalogHTTPService{
		providerSources: []ProviderSource{{ID: "source_1", DisplayID: 10002, Type: provider.TypeManual, Name: "Manual", Status: ProviderSourceStatusActive}},
	}
	handler := registerCatalogTestHandler(service)

	request := httptest.NewRequest(http.MethodGet, "/admin/catalog/provider-sources?source_type=manual&status=active&limit=9", nil)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", response.Code, response.Body.String())
	}
	if service.listSourceCalls != 1 {
		t.Fatalf("expected list provider sources once, got %d", service.listSourceCalls)
	}
	if service.sourceFilter.Type != provider.TypeManual || service.sourceFilter.Status != ProviderSourceStatusActive || service.sourceFilter.Limit != 9 {
		t.Fatalf("unexpected source filter: %+v", service.sourceFilter)
	}
}

func TestHTTPHandlerListProviderSourceReadinessUsesQueryFilters(t *testing.T) {
	service := &fakeCatalogHTTPService{
		sourceReadiness: []ProviderSourceReadiness{{
			PlanDisplayID:       10001,
			PlanCode:            "vps-s",
			PlanName:            "VPS Small",
			ProductType:         ProductTypeVPS,
			PlanStatus:          PlanStatusActive,
			PlanSourceDisplayID: 10003,
			PlanSourceStatus:    PlanSourceStatusActive,
			SourceDisplayID:     10002,
			SourceName:          "Hetzner Falkenstein",
			SourceType:          provider.TypeHetzner,
			SourceStatus:        ProviderSourceStatusActive,
			InventoryMode:       InventoryModeProviderLive,
			State:               ProviderSourceReadinessReady,
			Reason:              "Source is active and supports automatic provisioning.",
		}},
	}
	handler := registerCatalogTestHandler(service)

	request := httptest.NewRequest(http.MethodGet, "/admin/catalog/provider-readiness?product_type=vps&status=active&limit=7", nil)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", response.Code, response.Body.String())
	}
	if service.listReadinessCalls != 1 {
		t.Fatalf("expected list readiness once, got %d", service.listReadinessCalls)
	}
	if service.readinessFilter.ProductType != ProductTypeVPS || service.readinessFilter.PlanStatus != PlanStatusActive || service.readinessFilter.Limit != 7 {
		t.Fatalf("unexpected readiness filter: %+v", service.readinessFilter)
	}
	body := response.Body.String()
	for _, expected := range []string{`"plan_display_id":10001`, `"source_display_id":10002`, `"state":"ready"`} {
		if !strings.Contains(body, expected) {
			t.Fatalf("expected %s in response, got %s", expected, body)
		}
	}
	for _, blocked := range []string{"capability_profile", "provider_account_id"} {
		if strings.Contains(body, blocked) {
			t.Fatalf("response should not expose %s: %s", blocked, body)
		}
	}
}

func TestHTTPHandlerAdminProductsStillSupportsPost(t *testing.T) {
	service := &fakeCatalogHTTPService{
		product: Product{ID: "product_1", Type: ProductTypeVPS, Name: "VPS", Status: ProductStatusDraft, CreatedBy: "actor_1"},
	}
	handler := registerCatalogTestHandler(service)

	request := httptest.NewRequest(http.MethodPost, "/admin/catalog/products", strings.NewReader(`{"product_type":"vps","name":"VPS"}`))
	request.Header.Set(ActorHeader, "actor_1")
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d: %s", response.Code, response.Body.String())
	}
	if service.createProductCalls != 1 {
		t.Fatalf("expected create product once, got %d", service.createProductCalls)
	}
}
