package payment

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

func TestHTTPHandlerGetAdminDailyReconciliationUsesTenantAndDate(t *testing.T) {
	service := &fakeDailyReconciliationHTTPService{report: DailyReconciliationReport{
		TenantID:    "tenant-1",
		Date:        "2026-05-13",
		WindowFrom:  time.Date(2026, 5, 13, 0, 0, 0, 0, time.UTC),
		WindowTo:    time.Date(2026, 5, 14, 0, 0, 0, 0, time.UTC),
		GeneratedAt: time.Date(2026, 5, 13, 9, 0, 0, 0, time.UTC),
		Status:      ReconciliationReportStatusMismatched,
		Wallets: WalletReconciliationSummary{
			Checked:    1,
			Mismatched: 1,
			Mismatches: []WalletBalanceMismatch{{
				WalletID:            "wallet-1",
				WalletDisplayID:     70001,
				LastLedgerEntryID:   "ledger-1",
				LastLedgerDisplayID: 71001,
			}},
		},
	}}
	handler := registerDailyReconciliationTestHandler(service)

	request := httptest.NewRequest(http.MethodGet, "/admin/daily-reconciliation?date=2026-05-13", nil)
	request = request.WithContext(tenant.WithContext(request.Context(), tenant.NewContext("tenant-1")))
	request = request.WithContext(identity.WithActor(request.Context(), identity.NewActor("admin-1", "tenant-1", identity.ActorTypeResellerOwner)))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", response.Code, response.Body.String())
	}
	if service.calls != 1 ||
		service.input.TenantID != tenant.ID("tenant-1") ||
		!service.input.Date.Equal(time.Date(2026, 5, 13, 0, 0, 0, 0, time.UTC)) {
		t.Fatalf("unexpected reconciliation input: %+v calls=%d", service.input, service.calls)
	}
	if !strings.Contains(response.Body.String(), `"status":"mismatched"`) ||
		!strings.Contains(response.Body.String(), `"wallet_display_id":70001`) {
		t.Fatalf("expected reconciliation response, got %s", response.Body.String())
	}
	if strings.Contains(response.Body.String(), `"wallet_id"`) ||
		strings.Contains(response.Body.String(), `"last_ledger_entry_id"`) {
		t.Fatalf("response should expose public display IDs only, got %s", response.Body.String())
	}
}

func TestHTTPHandlerGetAdminDailyReconciliationRequiresDate(t *testing.T) {
	service := &fakeDailyReconciliationHTTPService{}
	handler := registerDailyReconciliationTestHandler(service)

	request := httptest.NewRequest(http.MethodGet, "/admin/daily-reconciliation", nil)
	request = request.WithContext(tenant.WithContext(request.Context(), tenant.NewContext("tenant-1")))
	request = request.WithContext(identity.WithActor(request.Context(), identity.NewActor("admin-1", "tenant-1", identity.ActorTypeResellerOwner)))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d: %s", response.Code, response.Body.String())
	}
	if service.calls != 0 {
		t.Fatalf("expected no service call, got %d", service.calls)
	}
}

func registerDailyReconciliationTestHandler(service HTTPService) http.Handler {
	mux := http.NewServeMux()
	NewHTTPHandler(service).RegisterRoutes(mux)
	return mux
}

type fakeDailyReconciliationHTTPService struct {
	report DailyReconciliationReport
	input  DailyReconciliationInput
	calls  int
}

func (service *fakeDailyReconciliationHTTPService) ListTransactions(ctx context.Context, filter TransactionFilter) ([]Transaction, error) {
	return nil, nil
}

func (service *fakeDailyReconciliationHTTPService) GetTransaction(ctx context.Context, lookup TransactionLookup) (Transaction, error) {
	return Transaction{}, nil
}

func (service *fakeDailyReconciliationHTTPService) PayInvoiceFromWallet(ctx context.Context, input PayInvoiceFromWalletInput) (WalletInvoicePayment, error) {
	return WalletInvoicePayment{}, nil
}

func (service *fakeDailyReconciliationHTTPService) ListPaymentReconciliations(ctx context.Context, filter ReconciliationFilter) ([]PaymentReconciliation, error) {
	return nil, nil
}

func (service *fakeDailyReconciliationHTTPService) GetPaymentReconciliation(ctx context.Context, lookup ReconciliationLookup) (PaymentReconciliation, error) {
	return PaymentReconciliation{}, nil
}

func (service *fakeDailyReconciliationHTTPService) BuildDailyReconciliationReport(ctx context.Context, input DailyReconciliationInput) (DailyReconciliationReport, error) {
	service.calls++
	service.input = input
	return service.report, nil
}
