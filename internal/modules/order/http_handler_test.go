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

func TestHTTPHandlerListClientOrdersUsesContextAndFilters(t *testing.T) {
	service := &fakeOrderHTTPService{
		orders: []Order{{
			ID:            "order_1",
			DisplayID:     30002,
			TenantID:      "tenant_1",
			BuyerUserID:   "buyer_1",
			TenantPlanID:  "tenant_plan_1",
			OrderStatus:   OrderStatusPendingPayment,
			BillingStatus: BillingStatusUnpaid,
		}},
	}
	handler := registerOrderTestHandler(service)

	request := httptest.NewRequest(http.MethodGet, "/client/orders?status=pending_payment&billing_status=unpaid&limit=15", nil)
	request = request.WithContext(tenant.WithContext(request.Context(), tenant.NewContext("tenant_1")))
	request = request.WithContext(identity.WithActor(request.Context(), identity.NewActor("buyer_1", "tenant_1", identity.ActorTypeClient)))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", response.Code, response.Body.String())
	}
	if service.listOrderCalls != 1 {
		t.Fatalf("expected list orders once, got %d", service.listOrderCalls)
	}
	if service.orderFilter.TenantID != tenant.ID("tenant_1") || service.orderFilter.BuyerUserID != identity.UserID("buyer_1") {
		t.Fatalf("unexpected tenant/buyer filter: %+v", service.orderFilter)
	}
	if service.orderFilter.OrderStatus != OrderStatusPendingPayment ||
		service.orderFilter.BillingStatus != BillingStatusUnpaid ||
		service.orderFilter.Limit != 15 {
		t.Fatalf("unexpected order filter: %+v", service.orderFilter)
	}
	if !strings.Contains(response.Body.String(), `"display_id":30002`) {
		t.Fatalf("expected order response, got %s", response.Body.String())
	}
}

func TestHTTPHandlerListAdminOrdersUsesTenantAndBuyerFilter(t *testing.T) {
	service := &fakeOrderHTTPService{
		orders: []Order{{
			ID:            "order_1",
			DisplayID:     30004,
			TenantID:      "tenant_1",
			BuyerUserID:   "buyer_2",
			TenantPlanID:  "tenant_plan_1",
			OrderStatus:   OrderStatusPaid,
			BillingStatus: BillingStatusPaid,
		}},
	}
	handler := registerOrderTestHandler(service)

	request := httptest.NewRequest(http.MethodGet, "/admin/orders?buyer_user_id=buyer_2&buyer_display_id=10002&display_id=30004&status=paid&billing_status=paid&amount_min=1000&amount_max=3000&limit=20", nil)
	request = request.WithContext(tenant.WithContext(request.Context(), tenant.NewContext("tenant_1")))
	request = request.WithContext(identity.WithActor(request.Context(), identity.NewActor("admin_1", "tenant_1", identity.ActorTypeResellerOwner)))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", response.Code, response.Body.String())
	}
	if service.listOrderCalls != 1 {
		t.Fatalf("expected list orders once, got %d", service.listOrderCalls)
	}
	if service.orderFilter.TenantID != tenant.ID("tenant_1") || service.orderFilter.BuyerUserID != identity.UserID("buyer_2") {
		t.Fatalf("unexpected admin order filter: %+v", service.orderFilter)
	}
	if service.orderFilter.DisplayID != 30004 ||
		service.orderFilter.BuyerDisplayID != 10002 ||
		service.orderFilter.OrderStatus != OrderStatusPaid ||
		service.orderFilter.BillingStatus != BillingStatusPaid ||
		service.orderFilter.AmountMinMinor == nil || *service.orderFilter.AmountMinMinor != 1000 ||
		service.orderFilter.AmountMaxMinor == nil || *service.orderFilter.AmountMaxMinor != 3000 {
		t.Fatalf("unexpected admin status filters: %+v", service.orderFilter)
	}
}

func TestHTTPHandlerListAdminOrdersRejectsBadBuyerDisplayID(t *testing.T) {
	service := &fakeOrderHTTPService{}
	handler := registerOrderTestHandler(service)

	request := httptest.NewRequest(http.MethodGet, "/admin/orders?buyer_display_id=bad", nil)
	request = request.WithContext(tenant.WithContext(request.Context(), tenant.NewContext("tenant_1")))
	request = request.WithContext(identity.WithActor(request.Context(), identity.NewActor("admin_1", "tenant_1", identity.ActorTypeResellerOwner)))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d: %s", response.Code, response.Body.String())
	}
	if service.listOrderCalls != 0 {
		t.Fatalf("expected no list call, got %d", service.listOrderCalls)
	}
}

func TestHTTPHandlerGetClientOrderUsesPathAndContext(t *testing.T) {
	service := &fakeOrderHTTPService{
		order: Order{
			ID:            "order_1",
			DisplayID:     30003,
			TenantID:      "tenant_1",
			BuyerUserID:   "buyer_1",
			TenantPlanID:  "tenant_plan_1",
			OrderStatus:   OrderStatusPaid,
			BillingStatus: BillingStatusPaid,
		},
	}
	handler := registerOrderTestHandler(service)

	request := httptest.NewRequest(http.MethodGet, "/client/orders/order_1", nil)
	request = request.WithContext(tenant.WithContext(request.Context(), tenant.NewContext("tenant_1")))
	request = request.WithContext(identity.WithActor(request.Context(), identity.NewActor("buyer_1", "tenant_1", identity.ActorTypeClient)))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", response.Code, response.Body.String())
	}
	if service.getOrderCalls != 1 {
		t.Fatalf("expected get order once, got %d", service.getOrderCalls)
	}
	if service.orderLookup.ID != OrderID("order_1") ||
		service.orderLookup.TenantID != tenant.ID("tenant_1") ||
		service.orderLookup.BuyerUserID != identity.UserID("buyer_1") {
		t.Fatalf("unexpected order lookup: %+v", service.orderLookup)
	}
}

func TestHTTPHandlerGetAdminOrderUsesTenantScopeOnly(t *testing.T) {
	service := &fakeOrderHTTPService{
		order: Order{
			ID:            "order_1",
			DisplayID:     30005,
			TenantID:      "tenant_1",
			BuyerUserID:   "buyer_2",
			TenantPlanID:  "tenant_plan_1",
			OrderStatus:   OrderStatusPaid,
			BillingStatus: BillingStatusPaid,
		},
	}
	handler := registerOrderTestHandler(service)

	request := httptest.NewRequest(http.MethodGet, "/admin/orders/order_1", nil)
	request = request.WithContext(tenant.WithContext(request.Context(), tenant.NewContext("tenant_1")))
	request = request.WithContext(identity.WithActor(request.Context(), identity.NewActor("admin_1", "tenant_1", identity.ActorTypeResellerOwner)))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", response.Code, response.Body.String())
	}
	if service.orderLookup.ID != OrderID("order_1") ||
		service.orderLookup.TenantID != tenant.ID("tenant_1") ||
		service.orderLookup.BuyerUserID != "" {
		t.Fatalf("unexpected admin order lookup: %+v", service.orderLookup)
	}
}

func TestHTTPHandlerTransitionAdminOrderStatusUsesTenantAndPath(t *testing.T) {
	service := &fakeOrderHTTPService{
		order: Order{
			ID:            "order_1",
			DisplayID:     30006,
			TenantID:      "tenant_1",
			BuyerUserID:   "buyer_2",
			TenantPlanID:  "tenant_plan_1",
			OrderStatus:   OrderStatusPaid,
			BillingStatus: BillingStatusPaid,
		},
	}
	handler := registerOrderTestHandler(service)

	request := httptest.NewRequest(http.MethodPatch, "/admin/orders/order_1/status", strings.NewReader(`{
		"from_status": "pending_payment",
		"to_status": "paid",
		"billing_status": "paid"
	}`))
	request = request.WithContext(tenant.WithContext(request.Context(), tenant.NewContext("tenant_1")))
	request = request.WithContext(identity.WithActor(request.Context(), identity.NewActor("admin_1", "tenant_1", identity.ActorTypeResellerOwner)))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", response.Code, response.Body.String())
	}
	if service.transitionOrderStatusCalls != 1 {
		t.Fatalf("expected status transition once, got %d", service.transitionOrderStatusCalls)
	}
	if service.transitionOrderStatusInput.ID != OrderID("order_1") ||
		service.transitionOrderStatusInput.TenantID != tenant.ID("tenant_1") ||
		service.transitionOrderStatusInput.ActorID != identity.UserID("admin_1") ||
		service.transitionOrderStatusInput.FromStatus != OrderStatusPendingPayment ||
		service.transitionOrderStatusInput.ToStatus != OrderStatusPaid ||
		service.transitionOrderStatusInput.BillingStatus != BillingStatusPaid {
		t.Fatalf("unexpected transition input: %+v", service.transitionOrderStatusInput)
	}
	if !strings.Contains(response.Body.String(), `"display_id":30006`) {
		t.Fatalf("expected order response, got %s", response.Body.String())
	}
}

func TestHTTPHandlerTransitionAdminOrderStatusRejectsBadChange(t *testing.T) {
	service := &fakeOrderHTTPService{}
	handler := registerOrderTestHandler(service)

	request := httptest.NewRequest(http.MethodPatch, "/admin/orders/order_1/status", strings.NewReader(`{
		"from_status": "pending_payment",
		"to_status": "refunded",
		"billing_status": "refunded"
	}`))
	request = request.WithContext(tenant.WithContext(request.Context(), tenant.NewContext("tenant_1")))
	request = request.WithContext(identity.WithActor(request.Context(), identity.NewActor("admin_1", "tenant_1", identity.ActorTypeResellerOwner)))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d: %s", response.Code, response.Body.String())
	}
	if service.transitionOrderStatusCalls != 0 {
		t.Fatalf("expected no transition call, got %d", service.transitionOrderStatusCalls)
	}
	if !strings.Contains(response.Body.String(), "order.status_transition_invalid") {
		t.Fatalf("expected transition validation error, got %s", response.Body.String())
	}
}

func TestHTTPHandlerTransitionAdminOrderStatusConflict(t *testing.T) {
	service := &fakeOrderHTTPService{transitionOrderStatusError: ErrOrderStatusConflict}
	handler := registerOrderTestHandler(service)

	request := httptest.NewRequest(http.MethodPatch, "/admin/orders/order_1/status", strings.NewReader(`{
		"from_status": "pending_payment",
		"to_status": "paid",
		"billing_status": "paid"
	}`))
	request = request.WithContext(tenant.WithContext(request.Context(), tenant.NewContext("tenant_1")))
	request = request.WithContext(identity.WithActor(request.Context(), identity.NewActor("admin_1", "tenant_1", identity.ActorTypeResellerOwner)))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusConflict {
		t.Fatalf("expected status 409, got %d: %s", response.Code, response.Body.String())
	}
	if !strings.Contains(response.Body.String(), "order.status_conflict") {
		t.Fatalf("expected conflict response, got %s", response.Body.String())
	}
}

func TestHTTPHandlerListClientOrdersRejectsBadStatus(t *testing.T) {
	service := &fakeOrderHTTPService{}
	handler := registerOrderTestHandler(service)

	request := httptest.NewRequest(http.MethodGet, "/client/orders?status=bad", nil)
	request = request.WithContext(tenant.WithContext(request.Context(), tenant.NewContext("tenant_1")))
	request = request.WithContext(identity.WithActor(request.Context(), identity.NewActor("buyer_1", "tenant_1", identity.ActorTypeClient)))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d: %s", response.Code, response.Body.String())
	}
	if service.listOrderCalls != 0 {
		t.Fatalf("expected no list call, got %d", service.listOrderCalls)
	}
}

func registerOrderTestHandler(service HTTPService) http.Handler {
	mux := http.NewServeMux()
	NewHTTPHandler(service).RegisterRoutes(mux)
	return mux
}

type fakeOrderHTTPService struct {
	createOrderCalls           int
	createOrderInput           CreateOrderInput
	listOrderCalls             int
	orderFilter                OrderFilter
	getOrderCalls              int
	orderLookup                OrderLookup
	transitionOrderStatusCalls int
	transitionOrderStatusInput TransitionOrderStatusInput
	transitionOrderStatusError error
	listServiceCalls           int
	serviceFilter              ServiceInstanceFilter
	getServiceCalls            int
	serviceLookup              ServiceInstanceLookup
	order                      Order
	orders                     []Order
	service                    ServiceInstance
	services                   []ServiceInstance
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

func (service *fakeOrderHTTPService) ListOrders(ctx context.Context, filter OrderFilter) ([]Order, error) {
	service.listOrderCalls++
	service.orderFilter = filter
	return service.orders, nil
}

func (service *fakeOrderHTTPService) GetOrder(ctx context.Context, lookup OrderLookup) (Order, error) {
	service.getOrderCalls++
	service.orderLookup = lookup
	return service.order, nil
}

func (service *fakeOrderHTTPService) TransitionOrderStatus(ctx context.Context, input TransitionOrderStatusInput) (Order, error) {
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return Order{}, err
	}
	service.transitionOrderStatusCalls++
	service.transitionOrderStatusInput = input
	if service.transitionOrderStatusError != nil {
		return Order{}, service.transitionOrderStatusError
	}
	if service.order.ID != "" {
		return service.order, nil
	}
	return Order{
		ID:            input.ID,
		TenantID:      input.TenantID,
		OrderStatus:   input.ToStatus,
		BillingStatus: input.BillingStatus,
	}, nil
}

func (service *fakeOrderHTTPService) ListServiceInstances(ctx context.Context, filter ServiceInstanceFilter) ([]ServiceInstance, error) {
	service.listServiceCalls++
	service.serviceFilter = filter
	return service.services, nil
}

func (service *fakeOrderHTTPService) GetServiceInstance(ctx context.Context, lookup ServiceInstanceLookup) (ServiceInstance, error) {
	service.getServiceCalls++
	service.serviceLookup = lookup
	return service.service, nil
}
