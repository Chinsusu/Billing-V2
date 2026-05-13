package wallet

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

func TestHTTPHandlerCreateAdminWalletAdjustmentUsesActorAndIdempotency(t *testing.T) {
	service := &fakeWalletHTTPService{adjustmentEntry: LedgerEntry{
		ID:          "ledger_1",
		WalletID:    "wallet_1",
		TenantID:    "tenant_1",
		Direction:   DirectionDebit,
		AmountMinor: 700,
		Currency:    "USD",
		EntryType:   EntryTypeAdjustment,
		Status:      LedgerStatusPosted,
	}}
	handler := registerWalletTestHandler(service)

	body := `{"wallet_id":"wallet_1","direction":"debit","amount_minor":700,"currency":"usd","reference_type":"manual_adjustment","reference_id":"00000000-0000-0000-0000-000000000102","reason":" finance correction ","correlation_id":"00000000-0000-0000-0000-000000000202"}`
	request := httptest.NewRequest(http.MethodPost, "/admin/wallet-adjustments", strings.NewReader(body))
	request.Header.Set(walletIdempotencyKeyHeader, " adjustment-key-1 ")
	request = request.WithContext(tenant.WithContext(request.Context(), tenant.NewContext("tenant_1")))
	request = request.WithContext(identity.WithActor(request.Context(), identity.NewActor("admin_1", "tenant_1", identity.ActorTypeResellerOwner)))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d: %s", response.Code, response.Body.String())
	}
	if service.createAdjustmentCalls != 1 {
		t.Fatalf("expected one adjustment call, got %d", service.createAdjustmentCalls)
	}
	if service.adjustmentInput.CreatedBy != identity.UserID("admin_1") ||
		service.adjustmentInput.IdempotencyKey != IdempotencyKey("adjustment-key-1") ||
		service.adjustmentInput.Currency != "USD" ||
		service.adjustmentInput.Reason != "finance correction" {
		t.Fatalf("unexpected adjustment input: %+v", service.adjustmentInput)
	}
	if !strings.Contains(response.Body.String(), `"entry_type":"adjustment"`) {
		t.Fatalf("expected ledger response, got %s", response.Body.String())
	}
}

func TestHTTPHandlerCreateAdminWalletRefundRequiresIdempotency(t *testing.T) {
	service := &fakeWalletHTTPService{}
	handler := registerWalletTestHandler(service)

	body := `{"wallet_id":"wallet_1","amount_minor":1200,"currency":"USD","reference_type":"invoice","reference_id":"00000000-0000-0000-0000-000000000101","reason":"refund request","correlation_id":"00000000-0000-0000-0000-000000000201"}`
	request := httptest.NewRequest(http.MethodPost, "/admin/wallet-refunds", strings.NewReader(body))
	request = request.WithContext(tenant.WithContext(request.Context(), tenant.NewContext("tenant_1")))
	request = request.WithContext(identity.WithActor(request.Context(), identity.NewActor("admin_1", "tenant_1", identity.ActorTypeResellerOwner)))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d: %s", response.Code, response.Body.String())
	}
	if !strings.Contains(response.Body.String(), "wallet.idempotency_key_missing") {
		t.Fatalf("expected idempotency validation error, got %s", response.Body.String())
	}
}
