package invoice

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/order"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

func TestHTTPHandlerListClientInvoicesUsesActorScope(t *testing.T) {
	service := &fakeInvoiceHTTPService{invoices: []Invoice{{
		ID:          "invoice_1",
		DisplayID:   80001,
		TenantID:    "tenant_1",
		BuyerUserID: "account_1",
		Status:      StatusIssued,
		Currency:    "USD",
		TotalMinor:  1200,
	}}}
	handler := registerInvoiceTestHandler(service)

	request := httptest.NewRequest(http.MethodGet, "/client/invoices?buyer_user_id=other&status=issued&limit=10", nil)
	request = request.WithContext(tenant.WithContext(request.Context(), tenant.NewContext("tenant_1")))
	request = request.WithContext(identity.WithActor(request.Context(), identity.NewActor("account_1", "tenant_1", identity.ActorTypeClient)))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", response.Code, response.Body.String())
	}
	if service.invoiceFilter.TenantID != tenant.ID("tenant_1") ||
		service.invoiceFilter.BuyerUserID != identity.UserID("account_1") ||
		service.invoiceFilter.Status != StatusIssued ||
		service.invoiceFilter.Limit != 10 {
		t.Fatalf("unexpected invoice filter: %+v", service.invoiceFilter)
	}
	if !strings.Contains(response.Body.String(), `"display_id":80001`) {
		t.Fatalf("expected invoice response, got %s", response.Body.String())
	}
}

func TestHTTPHandlerGetClientInvoiceIncludesItems(t *testing.T) {
	service := &fakeInvoiceHTTPService{detail: InvoiceDetail{
		Invoice: Invoice{ID: "invoice_1", DisplayID: 80002, TenantID: "tenant_1", BuyerUserID: "account_1", Currency: "USD"},
		Items: []Item{{
			ID:             "item_1",
			InvoiceID:      "invoice_1",
			TenantID:       "tenant_1",
			OrderID:        order.OrderID("order_1"),
			Description:    "VPS service",
			Quantity:       1,
			UnitPriceMinor: 1000,
			LineTotalMinor: 1000,
		}},
	}}
	handler := registerInvoiceTestHandler(service)

	request := httptest.NewRequest(http.MethodGet, "/client/invoices/invoice_1", nil)
	request = request.WithContext(tenant.WithContext(request.Context(), tenant.NewContext("tenant_1")))
	request = request.WithContext(identity.WithActor(request.Context(), identity.NewActor("account_1", "tenant_1", identity.ActorTypeClient)))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", response.Code, response.Body.String())
	}
	if service.invoiceLookup.ID != InvoiceID("invoice_1") ||
		service.invoiceLookup.TenantID != tenant.ID("tenant_1") ||
		service.invoiceLookup.BuyerUserID != identity.UserID("account_1") {
		t.Fatalf("unexpected invoice lookup: %+v", service.invoiceLookup)
	}
	if !strings.Contains(response.Body.String(), `"items":[`) || !strings.Contains(response.Body.String(), `"order_id":"order_1"`) {
		t.Fatalf("expected item detail response, got %s", response.Body.String())
	}
}

func TestHTTPHandlerListAdminInvoicesUsesFilters(t *testing.T) {
	service := &fakeInvoiceHTTPService{}
	handler := registerInvoiceTestHandler(service)

	request := httptest.NewRequest(http.MethodGet, "/admin/invoices?buyer_user_id=account_2&display_id=44001&order_id=order_2&status=paid&amount_min=100&amount_max=5000", nil)
	request = request.WithContext(tenant.WithContext(request.Context(), tenant.NewContext("tenant_1")))
	request = request.WithContext(identity.WithActor(request.Context(), identity.NewActor("admin_1", "tenant_1", identity.ActorTypeResellerOwner)))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", response.Code, response.Body.String())
	}
	if service.invoiceFilter.TenantID != tenant.ID("tenant_1") ||
		service.invoiceFilter.BuyerUserID != identity.UserID("account_2") ||
		service.invoiceFilter.DisplayID != 44001 ||
		service.invoiceFilter.OrderID != order.OrderID("order_2") ||
		service.invoiceFilter.Status != StatusPaid ||
		service.invoiceFilter.AmountMinMinor == nil || *service.invoiceFilter.AmountMinMinor != 100 ||
		service.invoiceFilter.AmountMaxMinor == nil || *service.invoiceFilter.AmountMaxMinor != 5000 {
		t.Fatalf("unexpected admin invoice filter: %+v", service.invoiceFilter)
	}
}

func TestHTTPHandlerListResellerInvoicesUsesTenantAndFilters(t *testing.T) {
	service := &fakeInvoiceHTTPService{invoices: []Invoice{{
		ID:          "invoice_3",
		DisplayID:   80003,
		TenantID:    "reseller_tenant",
		BuyerUserID: "account_3",
		Status:      StatusPaid,
		Currency:    "USD",
		TotalMinor:  5000,
	}}}
	handler := registerInvoiceTestHandler(service)

	request := httptest.NewRequest(http.MethodGet, "/reseller/invoices?buyer_user_id=account_3&display_id=80003&order_id=order_3&status=paid&amount_min=100&amount_max=5000&limit=12", nil)
	request = request.WithContext(tenant.WithContext(request.Context(), tenant.NewContext("reseller_tenant")))
	request = request.WithContext(identity.WithActor(request.Context(), identity.NewActor("reseller_1", "reseller_tenant", identity.ActorTypeResellerOwner)))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", response.Code, response.Body.String())
	}
	if service.invoiceFilter.TenantID != tenant.ID("reseller_tenant") ||
		service.invoiceFilter.BuyerUserID != identity.UserID("account_3") ||
		service.invoiceFilter.DisplayID != 80003 ||
		service.invoiceFilter.OrderID != order.OrderID("order_3") ||
		service.invoiceFilter.Status != StatusPaid ||
		service.invoiceFilter.AmountMinMinor == nil || *service.invoiceFilter.AmountMinMinor != 100 ||
		service.invoiceFilter.AmountMaxMinor == nil || *service.invoiceFilter.AmountMaxMinor != 5000 ||
		service.invoiceFilter.Limit != 12 {
		t.Fatalf("unexpected reseller invoice filter: %+v", service.invoiceFilter)
	}
	if !strings.Contains(response.Body.String(), `"display_id":80003`) {
		t.Fatalf("expected invoice response, got %s", response.Body.String())
	}
}

func TestHTTPHandlerResellerInvoiceMiddlewareRunsBeforeService(t *testing.T) {
	service := &fakeInvoiceHTTPService{}
	mux := http.NewServeMux()
	NewHTTPHandlerWithOptions(service, HTTPHandlerOptions{
		ResellerMiddleware: func(next http.HandlerFunc) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusForbidden)
			}
		},
	}).RegisterRoutes(mux)

	request := httptest.NewRequest(http.MethodGet, "/reseller/invoices", nil)
	request = request.WithContext(tenant.WithContext(request.Context(), tenant.NewContext("tenant_1")))
	response := httptest.NewRecorder()

	mux.ServeHTTP(response, request)

	if response.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d", response.Code)
	}
	if service.listCalls != 0 {
		t.Fatalf("expected service not to run, got %d calls", service.listCalls)
	}
}

func TestHTTPHandlerRejectsBadInvoiceAmountRange(t *testing.T) {
	service := &fakeInvoiceHTTPService{}
	handler := registerInvoiceTestHandler(service)

	request := httptest.NewRequest(http.MethodGet, "/admin/invoices?amount_min=200&amount_max=100", nil)
	request = request.WithContext(tenant.WithContext(request.Context(), tenant.NewContext("tenant_1")))
	request = request.WithContext(identity.WithActor(request.Context(), identity.NewActor("admin_1", "tenant_1", identity.ActorTypeResellerOwner)))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d: %s", response.Code, response.Body.String())
	}
	if service.listCalls != 0 {
		t.Fatalf("expected no invoice list call, got %d", service.listCalls)
	}
}

func TestHTTPHandlerRejectsBadInvoiceStatus(t *testing.T) {
	service := &fakeInvoiceHTTPService{}
	handler := registerInvoiceTestHandler(service)

	request := httptest.NewRequest(http.MethodGet, "/client/invoices?status=bad", nil)
	request = request.WithContext(tenant.WithContext(request.Context(), tenant.NewContext("tenant_1")))
	request = request.WithContext(identity.WithActor(request.Context(), identity.NewActor("account_1", "tenant_1", identity.ActorTypeClient)))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d: %s", response.Code, response.Body.String())
	}
	if service.listCalls != 0 {
		t.Fatalf("expected no invoice list call, got %d", service.listCalls)
	}
}

func registerInvoiceTestHandler(service HTTPService) http.Handler {
	mux := http.NewServeMux()
	NewHTTPHandler(service).RegisterRoutes(mux)
	return mux
}

type fakeInvoiceHTTPService struct {
	invoices      []Invoice
	detail        InvoiceDetail
	invoiceFilter InvoiceFilter
	invoiceLookup InvoiceLookup
	listCalls     int
	getCalls      int
}

func (service *fakeInvoiceHTTPService) ListInvoices(ctx context.Context, filter InvoiceFilter) ([]Invoice, error) {
	service.listCalls++
	service.invoiceFilter = filter
	return service.invoices, nil
}

func (service *fakeInvoiceHTTPService) GetInvoice(ctx context.Context, lookup InvoiceLookup) (InvoiceDetail, error) {
	service.getCalls++
	service.invoiceLookup = lookup
	return service.detail, nil
}
