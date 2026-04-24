package checkout

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/invoice"
	"github.com/Chinsusu/Billing-V2/internal/modules/order"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

func TestHTTPHandlerCheckoutUsesContextScopeAndHeader(t *testing.T) {
	service := &fakeCheckoutHTTPService{detail: testInvoiceDetail()}
	handler := registerCheckoutTestHandler(service)

	request := httptest.NewRequest(http.MethodPost, "/client/checkouts", strings.NewReader(`{"order_id":"order_1"}`))
	request.Header.Set(IdempotencyKeyHeader, " checkout-key-1 ")
	request = request.WithContext(tenant.WithContext(request.Context(), tenant.NewContext("tenant_1")))
	request = request.WithContext(identity.WithActor(request.Context(), identity.NewActor("buyer_1", "tenant_1", identity.ActorTypeClient)))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d: %s", response.Code, response.Body.String())
	}
	if service.calls != 1 {
		t.Fatalf("expected checkout once, got %d", service.calls)
	}
	if service.input.TenantID != tenant.ID("tenant_1") ||
		service.input.BuyerUserID != identity.UserID("buyer_1") ||
		service.input.OrderID != order.OrderID("order_1") ||
		service.input.IdempotencyKey != invoice.IdempotencyKey("checkout-key-1") {
		t.Fatalf("unexpected checkout input: %+v", service.input)
	}
	if !strings.Contains(response.Body.String(), `"display_id":70001`) {
		t.Fatalf("expected invoice response, got %s", response.Body.String())
	}
}

func TestHTTPHandlerCheckoutRequiresIdempotencyKey(t *testing.T) {
	service := &fakeCheckoutHTTPService{}
	handler := registerCheckoutTestHandler(service)

	request := httptest.NewRequest(http.MethodPost, "/client/checkouts", strings.NewReader(`{"order_id":"order_1"}`))
	request = request.WithContext(tenant.WithContext(request.Context(), tenant.NewContext("tenant_1")))
	request = request.WithContext(identity.WithActor(request.Context(), identity.NewActor("buyer_1", "tenant_1", identity.ActorTypeClient)))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d: %s", response.Code, response.Body.String())
	}
	if service.calls != 1 {
		t.Fatalf("expected service validation call, got %d", service.calls)
	}
	if !strings.Contains(response.Body.String(), "checkout.idempotency_key_missing") {
		t.Fatalf("expected idempotency error, got %s", response.Body.String())
	}
}

func TestHTTPHandlerCheckoutRejectsOtherTenantOrder(t *testing.T) {
	service := &fakeCheckoutHTTPService{err: tenant.ErrAccessDenied}
	handler := registerCheckoutTestHandler(service)

	request := httptest.NewRequest(http.MethodPost, "/client/checkouts", strings.NewReader(`{"order_id":"order_2"}`))
	request.Header.Set(IdempotencyKeyHeader, "checkout-key-2")
	request = request.WithContext(tenant.WithContext(request.Context(), tenant.NewContext("tenant_1")))
	request = request.WithContext(identity.WithActor(request.Context(), identity.NewActor("buyer_1", "tenant_1", identity.ActorTypeClient)))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d: %s", response.Code, response.Body.String())
	}
	if !strings.Contains(response.Body.String(), "tenant.context_invalid") {
		t.Fatalf("expected tenant scope error, got %s", response.Body.String())
	}
}

func TestHTTPHandlerCheckoutRejectsNonCheckoutableOrder(t *testing.T) {
	service := &fakeCheckoutHTTPService{err: invoice.ErrOrderNotCheckoutable}
	handler := registerCheckoutTestHandler(service)

	request := httptest.NewRequest(http.MethodPost, "/client/checkouts", strings.NewReader(`{"order_id":"order_1"}`))
	request.Header.Set(IdempotencyKeyHeader, "checkout-key-3")
	request = request.WithContext(tenant.WithContext(request.Context(), tenant.NewContext("tenant_1")))
	request = request.WithContext(identity.WithActor(request.Context(), identity.NewActor("buyer_1", "tenant_1", identity.ActorTypeClient)))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusConflict {
		t.Fatalf("expected status 409, got %d: %s", response.Code, response.Body.String())
	}
	if !strings.Contains(response.Body.String(), "checkout.order_not_checkoutable") {
		t.Fatalf("expected checkout conflict, got %s", response.Body.String())
	}
}

func registerCheckoutTestHandler(service HTTPService) http.Handler {
	mux := http.NewServeMux()
	NewHTTPHandler(service).RegisterRoutes(mux)
	return mux
}

type fakeCheckoutHTTPService struct {
	calls  int
	input  CheckoutOrderInput
	detail invoice.InvoiceDetail
	err    error
}

func (service *fakeCheckoutHTTPService) CheckoutOrder(ctx context.Context, input CheckoutOrderInput) (invoice.InvoiceDetail, error) {
	service.calls++
	service.input = input.Normalize()
	if service.err != nil {
		return invoice.InvoiceDetail{}, service.err
	}
	if err := service.input.Validate(); err != nil {
		return invoice.InvoiceDetail{}, err
	}
	return service.detail, nil
}

func testInvoiceDetail() invoice.InvoiceDetail {
	now := time.Date(2026, 4, 24, 10, 0, 0, 0, time.UTC)
	return invoice.InvoiceDetail{
		Invoice: invoice.Invoice{
			ID:            "invoice_1",
			DisplayID:     70001,
			TenantID:      "tenant_1",
			BuyerUserID:   "buyer_1",
			OrderID:       "order_1",
			Status:        invoice.StatusIssued,
			Currency:      "USD",
			SubtotalMinor: 1400,
			TotalMinor:    1400,
			IssuedAt:      now,
			CreatedAt:     now,
			UpdatedAt:     now,
			Metadata:      []byte(`{"source":"checkout"}`),
		},
		Items: []invoice.Item{{
			ID:             "item_1",
			Description:    "Checkout item",
			Quantity:       1,
			UnitPriceMinor: 1400,
			LineTotalMinor: 1400,
		}},
	}
}
