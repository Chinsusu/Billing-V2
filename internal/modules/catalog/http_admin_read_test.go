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
