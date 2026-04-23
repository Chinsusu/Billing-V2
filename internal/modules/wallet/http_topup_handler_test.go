package wallet

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

func TestHTTPHandlerCreateClientTopupRequestUsesContextScope(t *testing.T) {
	service := &fakeWalletHTTPService{
		wallet: Wallet{ID: "wallet_1", TenantID: "tenant_1", OwnerType: OwnerTypeUser, OwnerID: "account_1"},
		topup:  TopupRequest{ID: "topup_1", DisplayID: 90001, TenantID: "tenant_1", WalletID: "wallet_1"},
	}
	handler := registerWalletTestHandler(service)

	body := `{"wallet_id":"wallet_1","amount_minor":5000,"currency":"usd","payment_method":"bank_transfer","payment_reference":" bank ref "}`
	request := httptest.NewRequest(http.MethodPost, "/client/topup-requests", strings.NewReader(body))
	request.Header.Set("Idempotency-Key", " topup-key-1 ")
	request = request.WithContext(tenant.WithContext(request.Context(), tenant.NewContext("tenant_1")))
	request = request.WithContext(identity.WithActor(request.Context(), identity.NewActor("account_1", "tenant_1", identity.ActorTypeClient)))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d: %s", response.Code, response.Body.String())
	}
	if service.getWalletCalls != 1 || service.createTopupCalls != 1 {
		t.Fatalf("expected wallet check and create, got wallet=%d create=%d", service.getWalletCalls, service.createTopupCalls)
	}
	if service.walletLookup.OwnerType != OwnerTypeUser || service.walletLookup.OwnerID != OwnerID("account_1") {
		t.Fatalf("unexpected wallet lookup: %+v", service.walletLookup)
	}
	if service.topupInput.TenantID != tenant.ID("tenant_1") ||
		service.topupInput.RequestedBy != identity.UserID("account_1") ||
		service.topupInput.IdempotencyKey != IdempotencyKey("topup-key-1") {
		t.Fatalf("unexpected topup input scope: %+v", service.topupInput)
	}
	if !strings.Contains(response.Body.String(), `"display_id":90001`) {
		t.Fatalf("expected topup response, got %s", response.Body.String())
	}
}

func TestHTTPHandlerListClientTopupRequestsUsesActorScope(t *testing.T) {
	service := &fakeWalletHTTPService{topups: []TopupRequest{{ID: "topup_1", DisplayID: 90002, TenantID: "tenant_1"}}}
	handler := registerWalletTestHandler(service)

	request := httptest.NewRequest(http.MethodGet, "/client/topup-requests?requested_by=other&wallet_id=wallet_1&payment_method=bank_transfer&status=submitted&limit=10", nil)
	request = request.WithContext(tenant.WithContext(request.Context(), tenant.NewContext("tenant_1")))
	request = request.WithContext(identity.WithActor(request.Context(), identity.NewActor("account_1", "tenant_1", identity.ActorTypeClient)))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", response.Code, response.Body.String())
	}
	if service.topupFilter.TenantID != tenant.ID("tenant_1") ||
		service.topupFilter.RequestedBy != identity.UserID("account_1") ||
		service.topupFilter.WalletID != WalletID("wallet_1") ||
		service.topupFilter.PaymentMethod != PaymentMethodBankTransfer ||
		service.topupFilter.Status != TopupStatusSubmitted ||
		service.topupFilter.Limit != 10 {
		t.Fatalf("unexpected topup filter: %+v", service.topupFilter)
	}
}

func TestHTTPHandlerListAdminTopupRequestsUsesReviewFilters(t *testing.T) {
	service := &fakeWalletHTTPService{}
	handler := registerWalletTestHandler(service)

	request := httptest.NewRequest(http.MethodGet, "/admin/topup-requests?requested_by=account_2&status=under_review", nil)
	request = request.WithContext(tenant.WithContext(request.Context(), tenant.NewContext("tenant_1")))
	request = request.WithContext(identity.WithActor(request.Context(), identity.NewActor("admin_1", "tenant_1", identity.ActorTypeResellerOwner)))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", response.Code, response.Body.String())
	}
	if service.topupFilter.TenantID != tenant.ID("tenant_1") ||
		service.topupFilter.RequestedBy != identity.UserID("account_2") ||
		service.topupFilter.Status != TopupStatusUnderReview {
		t.Fatalf("unexpected admin topup filter: %+v", service.topupFilter)
	}
}

func TestHTTPHandlerRejectsBadTopupPaymentMethod(t *testing.T) {
	service := &fakeWalletHTTPService{}
	handler := registerWalletTestHandler(service)

	request := httptest.NewRequest(http.MethodGet, "/client/topup-requests?payment_method=bad", nil)
	request = request.WithContext(tenant.WithContext(request.Context(), tenant.NewContext("tenant_1")))
	request = request.WithContext(identity.WithActor(request.Context(), identity.NewActor("account_1", "tenant_1", identity.ActorTypeClient)))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d: %s", response.Code, response.Body.String())
	}
	if service.listTopupCalls != 0 {
		t.Fatalf("expected no topup list call, got %d", service.listTopupCalls)
	}
}
