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
	"github.com/Chinsusu/Billing-V2/internal/modules/wallet"
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

func TestHTTPHandlerCreateClientInvoiceWalletPaymentUsesActorTenantAndIdempotency(t *testing.T) {
	service := &fakePaymentHTTPService{
		paymentResult: WalletInvoicePayment{
			Invoice: invoice.InvoiceDetail{Invoice: invoice.Invoice{
				ID:          "invoice_1",
				DisplayID:   44001,
				TenantID:    "tenant_1",
				BuyerUserID: "account_1",
				Status:      invoice.StatusPaid,
				Currency:    "USD",
				TotalMinor:  2500,
			}},
			Transaction: Transaction{
				ID:            "txn_1",
				DisplayID:     51001,
				TenantID:      "tenant_1",
				AccountUserID: "account_1",
				InvoiceID:     "invoice_1",
				Type:          TransactionTypeCharge,
				Status:        TransactionStatusPosted,
				Currency:      "USD",
				AmountMinor:   2500,
			},
			LedgerEntry: wallet.LedgerEntry{
				ID:                "ledger_1",
				DisplayID:         52001,
				WalletID:          "wallet_1",
				TenantID:          "tenant_1",
				Direction:         wallet.DirectionDebit,
				EntryType:         wallet.EntryTypePurchase,
				Status:            wallet.LedgerStatusPosted,
				Currency:          "USD",
				AmountMinor:       2500,
				BalanceAfterMinor: 7500,
			},
		},
	}
	handler := registerPaymentTestHandler(service)

	request := httptest.NewRequest(http.MethodPost, "/client/invoice-wallet-payments", strings.NewReader(`{"invoice_id":"invoice_1","wallet_id":"wallet_1"}`))
	request.Header.Set(IdempotencyKeyHeader, " pay-key-1 ")
	request = request.WithContext(tenant.WithContext(request.Context(), tenant.NewContext("tenant_1")))
	request = request.WithContext(identity.WithActor(request.Context(), identity.NewActor("account_1", "tenant_1", identity.ActorTypeClient)))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d: %s", response.Code, response.Body.String())
	}
	if service.payCalls != 1 {
		t.Fatalf("expected pay invoice once, got %d", service.payCalls)
	}
	if service.payInput.TenantID != tenant.ID("tenant_1") ||
		service.payInput.ActorID != identity.UserID("account_1") ||
		service.payInput.InvoiceID != invoice.InvoiceID("invoice_1") ||
		service.payInput.WalletID != wallet.WalletID("wallet_1") ||
		service.payInput.IdempotencyKey != IdempotencyKey("pay-key-1") {
		t.Fatalf("unexpected wallet payment input: %+v", service.payInput)
	}
	body := response.Body.String()
	for _, expected := range []string{`"display_id":44001`, `"display_id":51001`, `"display_id":52001`} {
		if !strings.Contains(body, expected) {
			t.Fatalf("expected %s in payment response, got %s", expected, body)
		}
	}
}

func TestHTTPHandlerCreateClientInvoiceWalletPaymentRequiresIdempotencyKey(t *testing.T) {
	service := &fakePaymentHTTPService{}
	handler := registerPaymentTestHandler(service)

	request := httptest.NewRequest(http.MethodPost, "/client/invoice-wallet-payments", strings.NewReader(`{"invoice_id":"invoice_1","wallet_id":"wallet_1"}`))
	request = request.WithContext(tenant.WithContext(request.Context(), tenant.NewContext("tenant_1")))
	request = request.WithContext(identity.WithActor(request.Context(), identity.NewActor("account_1", "tenant_1", identity.ActorTypeClient)))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d: %s", response.Code, response.Body.String())
	}
	if service.payCalls != 0 {
		t.Fatalf("expected no payment call, got %d", service.payCalls)
	}
	if !strings.Contains(response.Body.String(), "payment.idempotency_key_missing") {
		t.Fatalf("expected idempotency validation response, got %s", response.Body.String())
	}
}

func TestHTTPHandlerCreateClientInvoiceWalletPaymentMapsServiceConflict(t *testing.T) {
	service := &fakePaymentHTTPService{payErr: ErrInvoiceNotPayable}
	handler := registerPaymentTestHandler(service)

	request := httptest.NewRequest(http.MethodPost, "/client/invoice-wallet-payments", strings.NewReader(`{"invoice_id":"invoice_1","wallet_id":"wallet_1"}`))
	request.Header.Set(IdempotencyKeyHeader, "pay-key-1")
	request = request.WithContext(tenant.WithContext(request.Context(), tenant.NewContext("tenant_1")))
	request = request.WithContext(identity.WithActor(request.Context(), identity.NewActor("account_1", "tenant_1", identity.ActorTypeClient)))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusConflict {
		t.Fatalf("expected status 409, got %d: %s", response.Code, response.Body.String())
	}
	if service.payCalls != 1 {
		t.Fatalf("expected one payment call, got %d", service.payCalls)
	}
}

func TestHTTPHandlerCreateClientInvoiceWalletPaymentMapsOrderConflict(t *testing.T) {
	service := &fakePaymentHTTPService{payErr: order.ErrOrderStatusConflict}
	handler := registerPaymentTestHandler(service)

	request := httptest.NewRequest(http.MethodPost, "/client/invoice-wallet-payments", strings.NewReader(`{"invoice_id":"invoice_1","wallet_id":"wallet_1"}`))
	request.Header.Set(IdempotencyKeyHeader, "pay-key-1")
	request = request.WithContext(tenant.WithContext(request.Context(), tenant.NewContext("tenant_1")))
	request = request.WithContext(identity.WithActor(request.Context(), identity.NewActor("account_1", "tenant_1", identity.ActorTypeClient)))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusConflict {
		t.Fatalf("expected status 409, got %d: %s", response.Code, response.Body.String())
	}
	if !strings.Contains(response.Body.String(), "order.status_conflict") {
		t.Fatalf("expected order conflict response, got %s", response.Body.String())
	}
}

func TestHTTPHandlerCreateClientInvoiceWalletPaymentMapsProvisioningSourceError(t *testing.T) {
	service := &fakePaymentHTTPService{payErr: order.ErrProvisioningSourceNotFound}
	handler := registerPaymentTestHandler(service)

	request := httptest.NewRequest(http.MethodPost, "/client/invoice-wallet-payments", strings.NewReader(`{"invoice_id":"invoice_1","wallet_id":"wallet_1"}`))
	request.Header.Set(IdempotencyKeyHeader, "pay-key-1")
	request = request.WithContext(tenant.WithContext(request.Context(), tenant.NewContext("tenant_1")))
	request = request.WithContext(identity.WithActor(request.Context(), identity.NewActor("account_1", "tenant_1", identity.ActorTypeClient)))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusConflict {
		t.Fatalf("expected status 409, got %d: %s", response.Code, response.Body.String())
	}
	if !strings.Contains(response.Body.String(), "order.provisioning_source_not_found") {
		t.Fatalf("expected provisioning source response, got %s", response.Body.String())
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

	request := httptest.NewRequest(http.MethodGet, "/admin/transactions?account_user_id=account_2&display_id=51001&status=posted&amount_min=100&amount_max=900", nil)
	request = request.WithContext(tenant.WithContext(request.Context(), tenant.NewContext("tenant_1")))
	request = request.WithContext(identity.WithActor(request.Context(), identity.NewActor("admin_1", "tenant_1", identity.ActorTypeResellerOwner)))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", response.Code, response.Body.String())
	}
	if service.filter.AccountUserID != identity.UserID("account_2") ||
		service.filter.DisplayID != 51001 ||
		service.filter.Status != TransactionStatusPosted ||
		service.filter.AmountMinMinor == nil || *service.filter.AmountMinMinor != 100 ||
		service.filter.AmountMaxMinor == nil || *service.filter.AmountMaxMinor != 900 {
		t.Fatalf("unexpected admin transaction filter: %+v", service.filter)
	}
}

func TestHTTPHandlerListResellerTransactionsUsesTenantAndFilters(t *testing.T) {
	service := &fakePaymentHTTPService{
		transactions: []Transaction{{
			ID:            "txn_3",
			DisplayID:     60003,
			TenantID:      "reseller_tenant",
			AccountUserID: "account_3",
			OrderID:       "order_3",
			InvoiceID:     "invoice_3",
			Type:          TransactionTypeCharge,
			Status:        TransactionStatusPosted,
			Currency:      "USD",
			AmountMinor:   900,
		}},
	}
	handler := registerPaymentTestHandler(service)

	request := httptest.NewRequest(http.MethodGet, "/reseller/transactions?account_user_id=account_3&display_id=60003&type=charge&status=posted&order_id=order_3&invoice_id=invoice_3&amount_min=100&amount_max=900&limit=12", nil)
	request = request.WithContext(tenant.WithContext(request.Context(), tenant.NewContext("reseller_tenant")))
	request = request.WithContext(identity.WithActor(request.Context(), identity.NewActor("reseller_1", "reseller_tenant", identity.ActorTypeResellerOwner)))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", response.Code, response.Body.String())
	}
	if service.filter.TenantID != tenant.ID("reseller_tenant") ||
		service.filter.AccountUserID != identity.UserID("account_3") ||
		service.filter.DisplayID != 60003 ||
		service.filter.OrderID != order.OrderID("order_3") ||
		service.filter.InvoiceID != invoice.InvoiceID("invoice_3") ||
		service.filter.Type != TransactionTypeCharge ||
		service.filter.Status != TransactionStatusPosted ||
		service.filter.AmountMinMinor == nil || *service.filter.AmountMinMinor != 100 ||
		service.filter.AmountMaxMinor == nil || *service.filter.AmountMaxMinor != 900 ||
		service.filter.Limit != 12 {
		t.Fatalf("unexpected reseller transaction filter: %+v", service.filter)
	}
	if !strings.Contains(response.Body.String(), `"display_id":60003`) {
		t.Fatalf("expected transaction response, got %s", response.Body.String())
	}
}

func TestHTTPHandlerResellerTransactionMiddlewareRunsBeforeService(t *testing.T) {
	service := &fakePaymentHTTPService{}
	mux := http.NewServeMux()
	NewHTTPHandlerWithOptions(service, HTTPHandlerOptions{
		ResellerMiddleware: func(next http.HandlerFunc) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusForbidden)
			}
		},
	}).RegisterRoutes(mux)

	request := httptest.NewRequest(http.MethodGet, "/reseller/transactions", nil)
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

func TestHTTPHandlerRejectsBadTransactionDisplayID(t *testing.T) {
	service := &fakePaymentHTTPService{}
	handler := registerPaymentTestHandler(service)

	request := httptest.NewRequest(http.MethodGet, "/admin/transactions?display_id=bad", nil)
	request = request.WithContext(tenant.WithContext(request.Context(), tenant.NewContext("tenant_1")))
	request = request.WithContext(identity.WithActor(request.Context(), identity.NewActor("admin_1", "tenant_1", identity.ActorTypeResellerOwner)))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d: %s", response.Code, response.Body.String())
	}
	if service.listCalls != 0 {
		t.Fatalf("expected no list call, got %d", service.listCalls)
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
	paymentResult           WalletInvoicePayment
	payInput                PayInvoiceFromWalletInput
	payErr                  error
	listCalls               int
	getCalls                int
	reconciliationListCalls int
	reconciliationGetCalls  int
	payCalls                int
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

func (service *fakePaymentHTTPService) PayInvoiceFromWallet(ctx context.Context, input PayInvoiceFromWalletInput) (WalletInvoicePayment, error) {
	service.payCalls++
	service.payInput = input
	if service.payErr != nil {
		return WalletInvoicePayment{}, service.payErr
	}
	return service.paymentResult, nil
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
