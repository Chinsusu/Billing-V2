package payment

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/invoice"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
	"github.com/Chinsusu/Billing-V2/internal/modules/wallet"
)

func TestHTTPHandlerListAdminPaymentReconciliationUsesFilters(t *testing.T) {
	service := &fakePaymentHTTPService{
		reconciliations: []PaymentReconciliation{testPaymentReconciliation()},
	}
	handler := registerPaymentTestHandler(service)

	request := httptest.NewRequest(http.MethodGet, "/admin/payment-reconciliation?account_user_id=buyer_1&display_id=30001&status=posted&provider=wallet&invoice_id=invoice_1&invoice_display_id=20001&wallet_id=wallet_1&wallet_display_id=40001&amount_min=1000&amount_max=2000&created_from=2026-04-23T00:00:00Z&created_to=2026-04-24T00:00:00Z&limit=20", nil)
	request = request.WithContext(tenant.WithContext(request.Context(), tenant.NewContext("tenant_1")))
	request = request.WithContext(identity.WithActor(request.Context(), identity.NewActor("admin_1", "tenant_1", identity.ActorTypeResellerOwner)))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", response.Code, response.Body.String())
	}
	if service.reconciliationListCalls != 1 {
		t.Fatalf("expected reconciliation list once, got %d", service.reconciliationListCalls)
	}
	if service.reconciliationFilter.TenantID != tenant.ID("tenant_1") ||
		service.reconciliationFilter.AccountUserID != identity.UserID("buyer_1") ||
		service.reconciliationFilter.DisplayID != 30001 ||
		service.reconciliationFilter.Status != TransactionStatusPosted ||
		service.reconciliationFilter.Provider != "wallet" ||
		service.reconciliationFilter.InvoiceID != invoice.InvoiceID("invoice_1") ||
		service.reconciliationFilter.InvoiceDisplayID != 20001 ||
		service.reconciliationFilter.WalletID != wallet.WalletID("wallet_1") ||
		service.reconciliationFilter.WalletDisplayID != 40001 ||
		service.reconciliationFilter.AmountMinMinor == nil || *service.reconciliationFilter.AmountMinMinor != 1000 ||
		service.reconciliationFilter.AmountMaxMinor == nil || *service.reconciliationFilter.AmountMaxMinor != 2000 ||
		service.reconciliationFilter.Limit != 20 {
		t.Fatalf("unexpected reconciliation filter: %+v", service.reconciliationFilter)
	}
	if service.reconciliationFilter.CreatedFrom.IsZero() || service.reconciliationFilter.CreatedTo.IsZero() {
		t.Fatalf("expected created time filters: %+v", service.reconciliationFilter)
	}
	if !strings.Contains(response.Body.String(), `"wallet_display_id":40001`) {
		t.Fatalf("expected reconciliation response, got %s", response.Body.String())
	}
}

func TestHTTPHandlerGetAdminPaymentReconciliationUsesTenantScope(t *testing.T) {
	service := &fakePaymentHTTPService{reconciliation: testPaymentReconciliation()}
	handler := registerPaymentTestHandler(service)

	request := httptest.NewRequest(http.MethodGet, "/admin/payment-reconciliation/txn_1", nil)
	request = request.WithContext(tenant.WithContext(request.Context(), tenant.NewContext("tenant_1")))
	request = request.WithContext(identity.WithActor(request.Context(), identity.NewActor("admin_1", "tenant_1", identity.ActorTypeResellerOwner)))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", response.Code, response.Body.String())
	}
	if service.reconciliationGetCalls != 1 ||
		service.reconciliationLookup.TenantID != tenant.ID("tenant_1") ||
		service.reconciliationLookup.TransactionID != TransactionID("txn_1") {
		t.Fatalf("unexpected reconciliation lookup: %+v", service.reconciliationLookup)
	}
}

func TestHTTPHandlerRejectsBadReconciliationTime(t *testing.T) {
	service := &fakePaymentHTTPService{}
	handler := registerPaymentTestHandler(service)

	request := httptest.NewRequest(http.MethodGet, "/admin/payment-reconciliation?created_from=bad", nil)
	request = request.WithContext(tenant.WithContext(request.Context(), tenant.NewContext("tenant_1")))
	request = request.WithContext(identity.WithActor(request.Context(), identity.NewActor("admin_1", "tenant_1", identity.ActorTypeResellerOwner)))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d: %s", response.Code, response.Body.String())
	}
	if service.reconciliationListCalls != 0 {
		t.Fatalf("expected no service call, got %d", service.reconciliationListCalls)
	}
}

func testPaymentReconciliation() PaymentReconciliation {
	paidAt := time.Date(2026, 4, 23, 2, 0, 0, 0, time.UTC)
	return PaymentReconciliation{
		Provider: "wallet",
		Transaction: Transaction{
			ID:            TransactionID("txn_1"),
			DisplayID:     30001,
			TenantID:      tenant.ID("tenant_1"),
			AccountUserID: identity.UserID("buyer_1"),
			InvoiceID:     invoice.InvoiceID("invoice_1"),
			Type:          TransactionTypeCharge,
			Status:        TransactionStatusPosted,
			Currency:      "USD",
			AmountMinor:   1800,
		},
		Invoice: ReconciliationInvoice{
			ID:         invoice.InvoiceID("invoice_1"),
			DisplayID:  20001,
			Status:     invoice.StatusPaid,
			TotalMinor: 1800,
			PaidAt:     paidAt,
		},
		Ledger: ReconciliationLedger{
			ID:                wallet.LedgerEntryID("ledger_1"),
			DisplayID:         50001,
			WalletID:          wallet.WalletID("wallet_1"),
			WalletDisplayID:   40001,
			Direction:         wallet.DirectionDebit,
			EntryType:         wallet.EntryTypePurchase,
			Status:            wallet.LedgerStatusPosted,
			BalanceAfterMinor: 700,
		},
	}
}
