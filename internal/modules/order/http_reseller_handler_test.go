package order

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

func TestHTTPHandlerListResellerOrdersUsesTenantAndFilters(t *testing.T) {
	service := &fakeOrderHTTPService{
		orders: []Order{{
			ID:            "order_2",
			DisplayID:     30007,
			TenantID:      "reseller_tenant",
			BuyerUserID:   "buyer_3",
			TenantPlanID:  "tenant_plan_2",
			OrderStatus:   OrderStatusPaid,
			BillingStatus: BillingStatusPaid,
		}},
	}
	handler := registerOrderTestHandler(service)

	request := httptest.NewRequest(http.MethodGet, "/reseller/orders?buyer_user_id=buyer_3&display_id=30007&status=paid&billing_status=paid&amount_min=1000&amount_max=3000&limit=20", nil)
	request = request.WithContext(tenant.WithContext(request.Context(), tenant.NewContext("reseller_tenant")))
	request = request.WithContext(identity.WithActor(request.Context(), identity.NewActor("reseller_1", "reseller_tenant", identity.ActorTypeResellerOwner)))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", response.Code, response.Body.String())
	}
	if service.listOrderCalls != 1 {
		t.Fatalf("expected list orders once, got %d", service.listOrderCalls)
	}
	if service.orderFilter.TenantID != tenant.ID("reseller_tenant") || service.orderFilter.BuyerUserID != identity.UserID("buyer_3") {
		t.Fatalf("unexpected reseller order filter: %+v", service.orderFilter)
	}
	if service.orderFilter.DisplayID != 30007 ||
		service.orderFilter.OrderStatus != OrderStatusPaid ||
		service.orderFilter.BillingStatus != BillingStatusPaid ||
		service.orderFilter.AmountMinMinor == nil || *service.orderFilter.AmountMinMinor != 1000 ||
		service.orderFilter.AmountMaxMinor == nil || *service.orderFilter.AmountMaxMinor != 3000 ||
		service.orderFilter.Limit != 20 {
		t.Fatalf("unexpected reseller status filters: %+v", service.orderFilter)
	}
	if !strings.Contains(response.Body.String(), `"display_id":30007`) {
		t.Fatalf("expected order response, got %s", response.Body.String())
	}
}

func TestHTTPHandlerResellerOrderMiddlewareRunsBeforeService(t *testing.T) {
	service := &fakeOrderHTTPService{}
	mux := http.NewServeMux()
	NewHTTPHandlerWithOptions(service, HTTPHandlerOptions{
		ResellerMiddleware: func(next http.HandlerFunc) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusForbidden)
			}
		},
	}).RegisterRoutes(mux)

	request := httptest.NewRequest(http.MethodGet, "/reseller/orders", nil)
	request = request.WithContext(tenant.WithContext(request.Context(), tenant.NewContext("tenant_1")))
	response := httptest.NewRecorder()

	mux.ServeHTTP(response, request)

	if response.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d", response.Code)
	}
	if service.listOrderCalls != 0 {
		t.Fatalf("expected service not to run, got %d calls", service.listOrderCalls)
	}
}
