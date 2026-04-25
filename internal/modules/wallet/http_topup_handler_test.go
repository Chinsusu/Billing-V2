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
	service := &fakeWalletHTTPService{topups: []TopupRequest{{
		ID:                   "topup_1",
		DisplayID:            90004,
		TenantID:             "tenant_1",
		WalletID:             "wallet_1",
		WalletDisplayID:      70004,
		RequestedBy:          "account_2",
		RequestedByDisplayID: 10002,
		AmountMinor:          5000,
		Currency:             "USD",
		PaymentMethod:        PaymentMethodBankTransfer,
		Status:               TopupStatusUnderReview,
	}}}
	handler := registerWalletTestHandler(service)

	request := httptest.NewRequest(http.MethodGet, "/admin/topup-requests?requested_by=account_2&requested_by_display_id=10002&wallet_display_id=70004&display_id=90004&status=under_review&amount_min=100&amount_max=5000", nil)
	request = request.WithContext(tenant.WithContext(request.Context(), tenant.NewContext("tenant_1")))
	request = request.WithContext(identity.WithActor(request.Context(), identity.NewActor("admin_1", "tenant_1", identity.ActorTypeResellerOwner)))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", response.Code, response.Body.String())
	}
	if service.topupFilter.TenantID != tenant.ID("tenant_1") ||
		service.topupFilter.RequestedBy != identity.UserID("account_2") ||
		service.topupFilter.RequestedByDisplayID != 10002 ||
		service.topupFilter.WalletDisplayID != 70004 ||
		service.topupFilter.DisplayID != 90004 ||
		service.topupFilter.Status != TopupStatusUnderReview ||
		service.topupFilter.AmountMinMinor == nil || *service.topupFilter.AmountMinMinor != 100 ||
		service.topupFilter.AmountMaxMinor == nil || *service.topupFilter.AmountMaxMinor != 5000 {
		t.Fatalf("unexpected admin topup filter: %+v", service.topupFilter)
	}
	body := response.Body.String()
	for _, expected := range []string{`"wallet_display_id":70004`, `"requested_by_display_id":10002`} {
		if !strings.Contains(body, expected) {
			t.Fatalf("expected %s in top-up response, got %s", expected, body)
		}
	}
}

func TestHTTPHandlerRejectsBadTopupWalletDisplayID(t *testing.T) {
	service := &fakeWalletHTTPService{}
	handler := registerWalletTestHandler(service)

	request := httptest.NewRequest(http.MethodGet, "/admin/topup-requests?wallet_display_id=bad", nil)
	request = request.WithContext(tenant.WithContext(request.Context(), tenant.NewContext("tenant_1")))
	request = request.WithContext(identity.WithActor(request.Context(), identity.NewActor("admin_1", "tenant_1", identity.ActorTypeResellerOwner)))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d: %s", response.Code, response.Body.String())
	}
	if service.listTopupCalls != 0 {
		t.Fatalf("expected no topup list call, got %d", service.listTopupCalls)
	}
}

func TestHTTPHandlerListResellerTopupRequestsUsesTenantAndFilters(t *testing.T) {
	service := &fakeWalletHTTPService{topups: []TopupRequest{{
		ID:            "topup_3",
		DisplayID:     90005,
		TenantID:      "reseller_tenant",
		WalletID:      "wallet_3",
		RequestedBy:   "account_3",
		AmountMinor:   5000,
		Currency:      "USD",
		PaymentMethod: PaymentMethodBankTransfer,
		Status:        TopupStatusSubmitted,
	}}}
	handler := registerWalletTestHandler(service)

	request := httptest.NewRequest(http.MethodGet, "/reseller/topup-requests?requested_by=account_3&wallet_id=wallet_3&payment_method=bank_transfer&status=submitted&display_id=90005&amount_min=100&amount_max=5000&limit=10", nil)
	request = request.WithContext(tenant.WithContext(request.Context(), tenant.NewContext("reseller_tenant")))
	request = request.WithContext(identity.WithActor(request.Context(), identity.NewActor("reseller_1", "reseller_tenant", identity.ActorTypeResellerOwner)))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", response.Code, response.Body.String())
	}
	if service.topupFilter.TenantID != tenant.ID("reseller_tenant") ||
		service.topupFilter.RequestedBy != identity.UserID("account_3") ||
		service.topupFilter.WalletID != WalletID("wallet_3") ||
		service.topupFilter.PaymentMethod != PaymentMethodBankTransfer ||
		service.topupFilter.Status != TopupStatusSubmitted ||
		service.topupFilter.DisplayID != 90005 ||
		service.topupFilter.AmountMinMinor == nil || *service.topupFilter.AmountMinMinor != 100 ||
		service.topupFilter.AmountMaxMinor == nil || *service.topupFilter.AmountMaxMinor != 5000 ||
		service.topupFilter.Limit != 10 {
		t.Fatalf("unexpected reseller top-up filter: %+v", service.topupFilter)
	}
	if !strings.Contains(response.Body.String(), `"display_id":90005`) {
		t.Fatalf("expected top-up response, got %s", response.Body.String())
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

func TestHTTPHandlerApproveAdminTopupRequestUsesReviewer(t *testing.T) {
	service := &fakeWalletHTTPService{topup: TopupRequest{ID: "topup_1", DisplayID: 90003, Status: TopupStatusApproved}}
	handler := registerWalletTestHandler(service)

	request := httptest.NewRequest(http.MethodPost, "/admin/topup-requests/topup_1/approve", strings.NewReader(`{"review_note":"verified bank transfer"}`))
	request = request.WithContext(tenant.WithContext(request.Context(), tenant.NewContext("tenant_1")))
	request = request.WithContext(identity.WithActor(request.Context(), identity.NewActor("admin_1", "tenant_1", identity.ActorTypeResellerOwner)))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", response.Code, response.Body.String())
	}
	if service.approveCalls != 1 {
		t.Fatalf("expected approve once, got %d", service.approveCalls)
	}
	if service.approveInput.ID != TopupRequestID("topup_1") ||
		service.approveInput.TenantID != tenant.ID("tenant_1") ||
		service.approveInput.ReviewedBy != identity.UserID("admin_1") ||
		service.approveInput.ReviewNote != "verified bank transfer" {
		t.Fatalf("unexpected approve input: %+v", service.approveInput)
	}
}

func TestHTTPHandlerRejectAdminTopupRequestRequiresNote(t *testing.T) {
	service := &fakeWalletHTTPService{}
	handler := registerWalletTestHandler(service)

	request := httptest.NewRequest(http.MethodPost, "/admin/topup-requests/topup_1/reject", strings.NewReader(`{"review_note":""}`))
	request = request.WithContext(tenant.WithContext(request.Context(), tenant.NewContext("tenant_1")))
	request = request.WithContext(identity.WithActor(request.Context(), identity.NewActor("admin_1", "tenant_1", identity.ActorTypeResellerOwner)))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d: %s", response.Code, response.Body.String())
	}
	if service.rejectCalls != 1 {
		t.Fatalf("expected service validation call, got %d", service.rejectCalls)
	}
}
