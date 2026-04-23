package catalog

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Chinsusu/Billing-V2/internal/modules/provider"
)

func TestHTTPHandlerUpdateProductStatusUsesPathAndBody(t *testing.T) {
	service := &fakeCatalogHTTPService{
		product: Product{ID: "product_1", DisplayID: 10001, Type: ProductTypeVPS, Name: "VPS", Status: ProductStatusDisabled},
	}
	handler := registerCatalogTestHandler(service)

	request := httptest.NewRequest(http.MethodPatch, "/admin/catalog/products/product_1", strings.NewReader(`{"status":"disabled"}`))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", response.Code, response.Body.String())
	}
	if service.updateProductStatusCalls != 1 {
		t.Fatalf("expected update product once, got %d", service.updateProductStatusCalls)
	}
	if service.updateProductStatusInput.ID != ProductID("product_1") || service.updateProductStatusInput.Status != ProductStatusDisabled {
		t.Fatalf("unexpected product status input: %+v", service.updateProductStatusInput)
	}
	if !strings.Contains(response.Body.String(), `"display_id":10001`) {
		t.Fatalf("expected updated product response, got %s", response.Body.String())
	}
}

func TestHTTPHandlerUpdatePlanStatusUsesPathAndBody(t *testing.T) {
	service := &fakeCatalogHTTPService{
		plan: Plan{ID: "plan_1", DisplayID: 10002, ProductID: "product_1", Code: "vps", Name: "VPS", Status: PlanStatusActive},
	}
	handler := registerCatalogTestHandler(service)

	request := httptest.NewRequest(http.MethodPatch, "/admin/catalog/plans/plan_1", strings.NewReader(`{"status":"active"}`))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", response.Code, response.Body.String())
	}
	if service.updatePlanStatusCalls != 1 {
		t.Fatalf("expected update plan once, got %d", service.updatePlanStatusCalls)
	}
	if service.updatePlanStatusInput.ID != PlanID("plan_1") || service.updatePlanStatusInput.Status != PlanStatusActive {
		t.Fatalf("unexpected plan status input: %+v", service.updatePlanStatusInput)
	}
}

func TestHTTPHandlerUpdateProviderSourceStatusUsesPathAndBody(t *testing.T) {
	service := &fakeCatalogHTTPService{
		source: ProviderSource{ID: "source_1", DisplayID: 10003, Type: provider.TypeManual, Name: "Manual", Status: ProviderSourceStatusMaintenance},
	}
	handler := registerCatalogTestHandler(service)

	request := httptest.NewRequest(http.MethodPatch, "/admin/catalog/provider-sources/source_1", strings.NewReader(`{"status":"maintenance"}`))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", response.Code, response.Body.String())
	}
	if service.updateProviderSourceStatusCalls != 1 {
		t.Fatalf("expected update source once, got %d", service.updateProviderSourceStatusCalls)
	}
	if service.updateProviderSourceStatusInput.ID != ProviderSourceID("source_1") ||
		service.updateProviderSourceStatusInput.Status != ProviderSourceStatusMaintenance {
		t.Fatalf("unexpected source status input: %+v", service.updateProviderSourceStatusInput)
	}
}

func TestHTTPHandlerUpdateProductStatusRequiresPathID(t *testing.T) {
	service := &fakeCatalogHTTPService{}
	handler := registerCatalogTestHandler(service)

	request := httptest.NewRequest(http.MethodPatch, "/admin/catalog/products/", strings.NewReader(`{"status":"active"}`))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d: %s", response.Code, response.Body.String())
	}
	if service.updateProductStatusCalls != 0 {
		t.Fatalf("expected no update call, got %d", service.updateProductStatusCalls)
	}
}
