package payment

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/invoice"
	"github.com/Chinsusu/Billing-V2/internal/modules/order"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

func TestHTTPHandlerListClientTransactionsUsesAccountScope(t *testing.T) {
	service := &fakePaymentHTTPService{
		transactions: []Transaction{{
			ID:            "txn_1",
			DisplayID:     60001,
			TenantID:      "tenant_1",
			AccountUserID: "account_1",
			OrderID:       "order_1",
			InvoiceID:     "invoice_1",
			Type:          TransactionTypeCharge,
			Status:        TransactionStatusPosted,
			Currency:      "USD",
			AmountMinor:   2500,
		}},
	}
	handler := registerPaymentTestHandler(service)

	request := httptest.NewRequest(http.MethodGet, "/client/transactions?type=charge&status=posted&order_id=order_1&invoice_id=invoice_1&limit=10", nil)
	request = request.WithContext(tenant.WithContext(request.Context(), tenant.NewContext("tenant_1")))
	request = request.WithContext(identity.WithActor(request.Context(), identity.NewActor("account_1", "tenant_1", identity.ActorTypeClient)))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", response.Code, response.Body.String())
	}
	if service.listCalls != 1 {
		t.Fatalf("expected list transactions once, got %d", service.listCalls)
	}
	if service.filter.TenantID != tenant.ID("tenant_1") ||
		service.filter.AccountUserID != identity.UserID("account_1") ||
		service.filter.OrderID != order.OrderID("order_1") ||
		service.filter.InvoiceID != invoice.InvoiceID("invoice_1") ||
		service.filter.Type != TransactionTypeCharge ||
		service.filter.Status != TransactionStatusPosted ||
		service.filter.Limit != 10 {
		t.Fatalf("unexpected transaction filter: %+v", service.filter)
	}
	if !strings.Contains(response.Body.String(), `"display_id":60001`) {
		t.Fatalf("expected transaction response, got %s", response.Body.String())
	}
}

func TestHTTPHandlerGetAdminTransactionUsesTenantScopeOnly(t *testing.T) {
	service := &fakePaymentHTTPService{
		transaction: Transaction{
			ID:            "txn_2",
			DisplayID:     60002,
			TenantID:      "tenant_1",
			AccountUserID: "account_2",
			Type:          TransactionTypeRefund,
			Status:        TransactionStatusPosted,
			Currency:      "USD",
			AmountMinor:   500,
		},
	}
	handler := registerPaymentTestHandler(service)

	request := httptest.NewRequest(http.MethodGet, "/admin/transactions/txn_2", nil)
	request = request.WithContext(tenant.WithContext(request.Context(), tenant.NewContext("tenant_1")))
	request = request.WithContext(identity.WithActor(request.Context(), identity.NewActor("admin_1", "tenant_1", identity.ActorTypeResellerOwner)))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", response.Code, response.Body.String())
	}
	if service.lookup.ID != TransactionID("txn_2") ||
		service.lookup.TenantID != tenant.ID("tenant_1") ||
		service.lookup.AccountUserID != "" {
		t.Fatalf("unexpected transaction lookup: %+v", service.lookup)
	}
}

func TestHTTPHandlerListAdminTransactionsUsesAccountFilter(t *testing.T) {
	service := &fakePaymentHTTPService{}
	handler := registerPaymentTestHandler(service)

	request := httptest.NewRequest(http.MethodGet, "/admin/transactions?account_user_id=account_2&status=posted", nil)
	request = request.WithContext(tenant.WithContext(request.Context(), tenant.NewContext("tenant_1")))
	request = request.WithContext(identity.WithActor(request.Context(), identity.NewActor("admin_1", "tenant_1", identity.ActorTypeResellerOwner)))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", response.Code, response.Body.String())
	}
	if service.filter.AccountUserID != identity.UserID("account_2") || service.filter.Status != TransactionStatusPosted {
		t.Fatalf("unexpected admin transaction filter: %+v", service.filter)
	}
}

func TestHTTPHandlerRejectsBadTransactionType(t *testing.T) {
	service := &fakePaymentHTTPService{}
	handler := registerPaymentTestHandler(service)

	request := httptest.NewRequest(http.MethodGet, "/client/transactions?type=bad", nil)
	request = request.WithContext(tenant.WithContext(request.Context(), tenant.NewContext("tenant_1")))
	request = request.WithContext(identity.WithActor(request.Context(), identity.NewActor("account_1", "tenant_1", identity.ActorTypeClient)))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d: %s", response.Code, response.Body.String())
	}
	if service.listCalls != 0 {
		t.Fatalf("expected no list call, got %d", service.listCalls)
	}
}

func TestHTTPHandlerClientTransactionRequiresActor(t *testing.T) {
	service := &fakePaymentHTTPService{}
	handler := registerPaymentTestHandler(service)

	request := httptest.NewRequest(http.MethodGet, "/client/transactions", nil)
	request = request.WithContext(tenant.WithContext(request.Context(), tenant.NewContext("tenant_1")))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d: %s", response.Code, response.Body.String())
	}
	if service.listCalls != 0 {
		t.Fatalf("expected no list call, got %d", service.listCalls)
	}
}

func registerPaymentTestHandler(service HTTPService) http.Handler {
	mux := http.NewServeMux()
	NewHTTPHandler(service).RegisterRoutes(mux)
	return mux
}

type fakePaymentHTTPService struct {
	transactions            []Transaction
	transaction             Transaction
	reconciliations         []PaymentReconciliation
	reconciliation          PaymentReconciliation
	filter                  TransactionFilter
	lookup                  TransactionLookup
	reconciliationFilter    ReconciliationFilter
	reconciliationLookup    ReconciliationLookup
	listCalls               int
	getCalls                int
	reconciliationListCalls int
	reconciliationGetCalls  int
}

func (service *fakePaymentHTTPService) ListTransactions(ctx context.Context, filter TransactionFilter) ([]Transaction, error) {
	service.listCalls++
	service.filter = filter
	return service.transactions, nil
}

func (service *fakePaymentHTTPService) GetTransaction(ctx context.Context, lookup TransactionLookup) (Transaction, error) {
	service.getCalls++
	service.lookup = lookup
	return service.transaction, nil
}

func (service *fakePaymentHTTPService) ListPaymentReconciliations(ctx context.Context, filter ReconciliationFilter) ([]PaymentReconciliation, error) {
	service.reconciliationListCalls++
	service.reconciliationFilter = filter
	return service.reconciliations, nil
}

func (service *fakePaymentHTTPService) GetPaymentReconciliation(ctx context.Context, lookup ReconciliationLookup) (PaymentReconciliation, error) {
	service.reconciliationGetCalls++
	service.reconciliationLookup = lookup
	return service.reconciliation, nil
}
