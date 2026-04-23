package order

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Chinsusu/Billing-V2/internal/modules/catalog"
	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

func TestHTTPHandlerCreateClientOrderUsesContextAndHeaders(t *testing.T) {
	service := &fakeOrderHTTPService{
		order: Order{
			ID:             "order_1",
			DisplayID:      30001,
			TenantID:       "tenant_1",
			BuyerUserID:    "buyer_1",
			TenantPlanID:   "tenant_plan_1",
			Quantity:       2,
			Currency:       "USD",
			UnitPriceMinor: 1000,
			TotalMinor:     2000,
			OrderStatus:    OrderStatusPendingPayment,
			BillingStatus:  BillingStatusUnpaid,
		},
	}
	handler := registerOrderTestHandler(service)

	request := httptest.NewRequest(http.MethodPost, "/client/orders", strings.NewReader(`{
		"tenant_plan_id": "tenant_plan_1",
		"quantity": 2,
		"currency": "usd",
		"unit_price_minor": 1000,
		"total_minor": 2000
	}`))
	request.Header.Set(IdempotencyKeyHeader, " order-key-1 ")
	request = request.WithContext(tenant.WithContext(request.Context(), tenant.NewContext("tenant_1")))
	request = request.WithContext(identity.WithActor(request.Context(), identity.NewActor("buyer_1", "tenant_1", identity.ActorTypeClient)))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d: %s", response.Code, response.Body.String())
	}
	if service.createOrderCalls != 1 {
		t.Fatalf("expected create order once, got %d", service.createOrderCalls)
	}
	if service.createOrderInput.TenantID != tenant.ID("tenant_1") || service.createOrderInput.BuyerUserID != identity.UserID("buyer_1") {
		t.Fatalf("unexpected tenant/buyer input: %+v", service.createOrderInput)
	}
	if service.createOrderInput.IdempotencyKey != IdempotencyKey("order-key-1") {
		t.Fatalf("expected idempotency key from header, got %q", service.createOrderInput.IdempotencyKey)
	}
	if strings.Contains(response.Body.String(), "idempotency") {
		t.Fatalf("response should not expose idempotency key: %s", response.Body.String())
	}
	if !strings.Contains(response.Body.String(), `"display_id":30001`) {
		t.Fatalf("expected order response, got %s", response.Body.String())
	}
}

func TestHTTPHandlerClientOrderRequiresTenant(t *testing.T) {
	service := &fakeOrderHTTPService{}
	handler := registerOrderTestHandler(service)

	request := httptest.NewRequest(http.MethodPost, "/client/orders", strings.NewReader(`{}`))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d: %s", response.Code, response.Body.String())
	}
	if service.createOrderCalls != 0 {
		t.Fatalf("expected no service call, got %d", service.createOrderCalls)
	}
}

func TestHTTPHandlerClientOrderRequiresActor(t *testing.T) {
	service := &fakeOrderHTTPService{}
	handler := registerOrderTestHandler(service)

	request := httptest.NewRequest(http.MethodPost, "/client/orders", strings.NewReader(`{}`))
	request = request.WithContext(tenant.WithContext(request.Context(), tenant.NewContext("tenant_1")))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d: %s", response.Code, response.Body.String())
	}
	if service.createOrderCalls != 0 {
		t.Fatalf("expected no service call, got %d", service.createOrderCalls)
	}
}

func TestHTTPHandlerClientOrderRequiresIdempotencyKey(t *testing.T) {
	service := &fakeOrderHTTPService{}
	handler := registerOrderTestHandler(service)

	request := httptest.NewRequest(http.MethodPost, "/client/orders", strings.NewReader(`{
		"tenant_plan_id": "tenant_plan_1",
		"quantity": 1,
		"currency": "USD",
		"unit_price_minor": 1000,
		"total_minor": 1000
	}`))
	request = request.WithContext(tenant.WithContext(request.Context(), tenant.NewContext("tenant_1")))
	request = request.WithContext(identity.WithActor(request.Context(), identity.NewActor("buyer_1", "tenant_1", identity.ActorTypeClient)))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d: %s", response.Code, response.Body.String())
	}
	if !strings.Contains(response.Body.String(), "order.idempotency_key_missing") {
		t.Fatalf("expected idempotency validation error, got %s", response.Body.String())
	}
}

func registerOrderTestHandler(service HTTPService) http.Handler {
	mux := http.NewServeMux()
	NewHTTPHandler(service).RegisterRoutes(mux)
	return mux
}

type fakeOrderHTTPService struct {
	createOrderCalls int
	createOrderInput CreateOrderInput
	order            Order
}

func (service *fakeOrderHTTPService) CreateOrder(ctx context.Context, input CreateOrderInput) (Order, error) {
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return Order{}, err
	}
	service.createOrderCalls++
	service.createOrderInput = input
	if service.order.ID != "" {
		return service.order, nil
	}
	return Order{
		ID:             "order_1",
		TenantID:       input.TenantID,
		BuyerUserID:    input.BuyerUserID,
		TenantPlanID:   catalog.TenantPlanID(input.TenantPlanID),
		Quantity:       input.Quantity,
		Currency:       input.Currency,
		UnitPriceMinor: input.UnitPriceMinor,
		DiscountMinor:  input.DiscountMinor,
		TotalMinor:     input.TotalMinor,
		OrderStatus:    input.OrderStatus,
		BillingStatus:  input.BillingStatus,
	}, nil
}
