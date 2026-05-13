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
	service := &fakeOrderHTTPService{services: []ServiceInstance{{
		ID:                      "service_1",
		DisplayID:               50002,
		TenantID:                "tenant_1",
		OrderID:                 "order_1",
		OrderDisplayID:          30005,
		BuyerDisplayID:          10002,
		TenantPlanID:            catalog.TenantPlanID("tenant_plan_1"),
		ProviderSourceID:        catalog.ProviderSourceID("source_1"),
		ProviderSourceDisplayID: 10003,
		Status:                  ServiceStatusActive,
		BillingStatus:           BillingStatusPaid,
	}}}
	handler := registerOrderTestHandler(service)

	request := httptest.NewRequest(http.MethodGet, "/admin/services?buyer_user_id=buyer_2&buyer_display_id=10002&display_id=50002&order_display_id=30005&provider_source_display_id=10003&status=active", nil)
	request = request.WithContext(tenant.WithContext(request.Context(), tenant.NewContext("tenant_1")))
	request = request.WithContext(identity.WithActor(request.Context(), identity.NewActor("admin_1", "tenant_1", identity.ActorTypeResellerOwner)))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", response.Code, response.Body.String())
	}
	if service.serviceFilter.TenantID != tenant.ID("tenant_1") ||
		service.serviceFilter.BuyerUserID != identity.UserID("buyer_2") ||
		service.serviceFilter.BuyerDisplayID != 10002 ||
		service.serviceFilter.DisplayID != 50002 ||
		service.serviceFilter.OrderDisplayID != 30005 ||
		service.serviceFilter.ProviderSourceDisplayID != 10003 ||
		service.serviceFilter.Status != ServiceStatusActive {
		t.Fatalf("unexpected admin service filter: %+v", service.serviceFilter)
	}
	body := response.Body.String()
	for _, expected := range []string{`"order_display_id":30005`, `"buyer_display_id":10002`, `"provider_source_display_id":10003`} {
		if !strings.Contains(body, expected) {
			t.Fatalf("expected %s in service response, got %s", expected, body)
		}
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

func TestHTTPHandlerSuspendAdminServiceUsesTenantActorAndReason(t *testing.T) {
	service := &fakeOrderHTTPService{
		service: ServiceInstance{
			ID:               "service_1",
			DisplayID:        50004,
			TenantID:         "tenant_1",
			OrderID:          "order_1",
			TenantPlanID:     catalog.TenantPlanID("tenant_plan_1"),
			ProviderSourceID: catalog.ProviderSourceID("source_1"),
			Status:           ServiceStatusSuspended,
			BillingStatus:    BillingStatusPaid,
		},
	}
	handler := registerOrderTestHandler(service)

	request := httptest.NewRequest(http.MethodPost, "/admin/services/service_1/suspend", strings.NewReader(`{
		"from_status": "active",
		"reason": "abuse ticket AB-1",
		"notify_client": true
	}`))
	request = request.WithContext(tenant.WithContext(request.Context(), tenant.NewContext("tenant_1")))
	request = request.WithContext(identity.WithActor(request.Context(), identity.NewActor("admin_1", "tenant_1", identity.ActorTypeResellerOwner)))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", response.Code, response.Body.String())
	}
	if service.transitionServiceLifecycleCalls != 1 {
		t.Fatalf("expected lifecycle transition once, got %d", service.transitionServiceLifecycleCalls)
	}
	input := service.transitionServiceLifecycleInput
	if input.ID != ServiceID("service_1") ||
		input.TenantID != tenant.ID("tenant_1") ||
		input.ActorID != "admin_1" ||
		input.Action != ServiceLifecycleActionSuspend ||
		input.FromStatus != ServiceStatusActive ||
		input.ToStatus != ServiceStatusSuspended ||
		input.SuspensionReason != SuspensionReasonManualAdmin ||
		input.Reason != "abuse ticket AB-1" {
		t.Fatalf("unexpected lifecycle input: %+v", input)
	}
}

func TestHTTPHandlerUnsuspendResellerServiceSetsPaidBilling(t *testing.T) {
	service := &fakeOrderHTTPService{}
	handler := registerOrderTestHandler(service)

	request := httptest.NewRequest(http.MethodPost, "/reseller/services/service_2/unsuspend", strings.NewReader(`{
		"from_status": "suspended",
		"reason": "abuse cleared"
	}`))
	request = request.WithContext(tenant.WithContext(request.Context(), tenant.NewContext("reseller_tenant")))
	request = request.WithContext(identity.WithActor(request.Context(), identity.NewActor("reseller_1", "reseller_tenant", identity.ActorTypeResellerOwner)))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", response.Code, response.Body.String())
	}
	input := service.transitionServiceLifecycleInput
	if input.Action != ServiceLifecycleActionUnsuspend ||
		input.ToStatus != ServiceStatusActive ||
		input.BillingStatus != BillingStatusPaid ||
		input.SuspensionReason != "" {
		t.Fatalf("unexpected unsuspend input: %+v", input)
	}
}

func TestHTTPHandlerTerminateAdminServiceRejectsMissingReason(t *testing.T) {
	service := &fakeOrderHTTPService{}
	handler := registerOrderTestHandler(service)

	request := httptest.NewRequest(http.MethodPost, "/admin/services/service_1/terminate", strings.NewReader(`{
		"from_status": "suspended"
	}`))
	request = request.WithContext(tenant.WithContext(request.Context(), tenant.NewContext("tenant_1")))
	request = request.WithContext(identity.WithActor(request.Context(), identity.NewActor("admin_1", "tenant_1", identity.ActorTypeResellerOwner)))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d: %s", response.Code, response.Body.String())
	}
	if service.transitionServiceLifecycleCalls != 0 {
		t.Fatalf("expected no lifecycle transition, got %d", service.transitionServiceLifecycleCalls)
	}
	if !strings.Contains(response.Body.String(), "service.reason_missing") {
		t.Fatalf("expected reason validation error, got %s", response.Body.String())
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
