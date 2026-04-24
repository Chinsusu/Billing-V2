package order

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Chinsusu/Billing-V2/internal/modules/catalog"
	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

func TestHTTPHandlerListClientServicesUsesAccountScope(t *testing.T) {
	service := &fakeOrderHTTPService{
		services: []ServiceInstance{{
			ID:               "service_1",
			DisplayID:        50001,
			TenantID:         "tenant_1",
			OrderID:          "order_1",
			TenantPlanID:     catalog.TenantPlanID("tenant_plan_1"),
			ProviderSourceID: catalog.ProviderSourceID("source_1"),
			Status:           ServiceStatusActive,
			BillingStatus:    BillingStatusPaid,
		}},
	}
	handler := registerOrderTestHandler(service)

	request := httptest.NewRequest(http.MethodGet, "/client/services?status=active&order_id=order_1&limit=10", nil)
	request = request.WithContext(tenant.WithContext(request.Context(), tenant.NewContext("tenant_1")))
	request = request.WithContext(identity.WithActor(request.Context(), identity.NewActor("buyer_1", "tenant_1", identity.ActorTypeClient)))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", response.Code, response.Body.String())
	}
	if service.listServiceCalls != 1 {
		t.Fatalf("expected list services once, got %d", service.listServiceCalls)
	}
	if service.serviceFilter.TenantID != tenant.ID("tenant_1") ||
		service.serviceFilter.BuyerUserID != identity.UserID("buyer_1") ||
		service.serviceFilter.OrderID != OrderID("order_1") ||
		service.serviceFilter.Status != ServiceStatusActive ||
		service.serviceFilter.Limit != 10 {
		t.Fatalf("unexpected service filter: %+v", service.serviceFilter)
	}
	if !strings.Contains(response.Body.String(), `"display_id":50001`) {
		t.Fatalf("expected service response, got %s", response.Body.String())
	}
}

func TestHTTPHandlerGetAdminServiceUsesTenantScopeOnly(t *testing.T) {
	service := &fakeOrderHTTPService{
		service: ServiceInstance{
			ID:               "service_1",
			DisplayID:        50002,
			TenantID:         "tenant_1",
			OrderID:          "order_1",
			TenantPlanID:     catalog.TenantPlanID("tenant_plan_1"),
			ProviderSourceID: catalog.ProviderSourceID("source_1"),
			Status:           ServiceStatusActive,
			BillingStatus:    BillingStatusPaid,
		},
	}
	handler := registerOrderTestHandler(service)

	request := httptest.NewRequest(http.MethodGet, "/admin/services/service_1", nil)
	request = request.WithContext(tenant.WithContext(request.Context(), tenant.NewContext("tenant_1")))
	request = request.WithContext(identity.WithActor(request.Context(), identity.NewActor("admin_1", "tenant_1", identity.ActorTypeResellerOwner)))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", response.Code, response.Body.String())
	}
	if service.serviceLookup.ID != ServiceID("service_1") ||
		service.serviceLookup.TenantID != tenant.ID("tenant_1") ||
		service.serviceLookup.BuyerUserID != "" {
		t.Fatalf("unexpected service lookup: %+v", service.serviceLookup)
	}
}

func TestHTTPHandlerListAdminServicesUsesSearchFilters(t *testing.T) {
	service := &fakeOrderHTTPService{}
	handler := registerOrderTestHandler(service)

	request := httptest.NewRequest(http.MethodGet, "/admin/services?buyer_user_id=buyer_2&display_id=50002&order_display_id=30005&status=active", nil)
	request = request.WithContext(tenant.WithContext(request.Context(), tenant.NewContext("tenant_1")))
	request = request.WithContext(identity.WithActor(request.Context(), identity.NewActor("admin_1", "tenant_1", identity.ActorTypeResellerOwner)))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", response.Code, response.Body.String())
	}
	if service.serviceFilter.TenantID != tenant.ID("tenant_1") ||
		service.serviceFilter.BuyerUserID != identity.UserID("buyer_2") ||
		service.serviceFilter.DisplayID != 50002 ||
		service.serviceFilter.OrderDisplayID != 30005 ||
		service.serviceFilter.Status != ServiceStatusActive {
		t.Fatalf("unexpected admin service filter: %+v", service.serviceFilter)
	}
}

func TestHTTPHandlerListResellerServicesUsesTenantAndFilters(t *testing.T) {
	service := &fakeOrderHTTPService{
		services: []ServiceInstance{{
			ID:               "service_3",
			DisplayID:        50003,
			TenantID:         "reseller_tenant",
			OrderID:          "order_3",
			TenantPlanID:     catalog.TenantPlanID("tenant_plan_3"),
			ProviderSourceID: catalog.ProviderSourceID("source_3"),
			Status:           ServiceStatusActive,
			BillingStatus:    BillingStatusPaid,
		}},
	}
	handler := registerOrderTestHandler(service)

	request := httptest.NewRequest(http.MethodGet, "/reseller/services?buyer_user_id=buyer_3&display_id=50003&order_display_id=30007&status=active&limit=10", nil)
	request = request.WithContext(tenant.WithContext(request.Context(), tenant.NewContext("reseller_tenant")))
	request = request.WithContext(identity.WithActor(request.Context(), identity.NewActor("reseller_1", "reseller_tenant", identity.ActorTypeResellerOwner)))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", response.Code, response.Body.String())
	}
	if service.listServiceCalls != 1 {
		t.Fatalf("expected list services once, got %d", service.listServiceCalls)
	}
	if service.serviceFilter.TenantID != tenant.ID("reseller_tenant") ||
		service.serviceFilter.BuyerUserID != identity.UserID("buyer_3") ||
		service.serviceFilter.DisplayID != 50003 ||
		service.serviceFilter.OrderDisplayID != 30007 ||
		service.serviceFilter.Status != ServiceStatusActive ||
		service.serviceFilter.Limit != 10 {
		t.Fatalf("unexpected reseller service filter: %+v", service.serviceFilter)
	}
	if !strings.Contains(response.Body.String(), `"display_id":50003`) {
		t.Fatalf("expected service response, got %s", response.Body.String())
	}
}

func TestHTTPHandlerListServicesRejectsBadStatus(t *testing.T) {
	service := &fakeOrderHTTPService{}
	handler := registerOrderTestHandler(service)

	request := httptest.NewRequest(http.MethodGet, "/client/services?status=bad", nil)
	request = request.WithContext(tenant.WithContext(request.Context(), tenant.NewContext("tenant_1")))
	request = request.WithContext(identity.WithActor(request.Context(), identity.NewActor("buyer_1", "tenant_1", identity.ActorTypeClient)))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d: %s", response.Code, response.Body.String())
	}
	if service.listServiceCalls != 0 {
		t.Fatalf("expected no list service call, got %d", service.listServiceCalls)
	}
}

func TestHTTPHandlerResellerServiceMiddlewareRunsBeforeService(t *testing.T) {
	service := &fakeOrderHTTPService{}
	mux := http.NewServeMux()
	NewHTTPHandlerWithOptions(service, HTTPHandlerOptions{
		ResellerServiceMiddleware: func(next http.HandlerFunc) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusForbidden)
			}
		},
	}).RegisterRoutes(mux)

	request := httptest.NewRequest(http.MethodGet, "/reseller/services", nil)
	request = request.WithContext(tenant.WithContext(request.Context(), tenant.NewContext("tenant_1")))
	response := httptest.NewRecorder()

	mux.ServeHTTP(response, request)

	if response.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d", response.Code)
	}
	if service.listServiceCalls != 0 {
		t.Fatalf("expected service not to run, got %d calls", service.listServiceCalls)
	}
}
