package wallet

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

func TestHTTPHandlerListClientWalletsUsesActorScope(t *testing.T) {
	service := &fakeWalletHTTPService{
		wallets: []Wallet{{
			ID:                    "wallet_1",
			DisplayID:             70001,
			TenantID:              "tenant_1",
			OwnerType:             OwnerTypeUser,
			OwnerID:               OwnerID("account_1"),
			Currency:              "USD",
			Status:                StatusActive,
			AvailableBalanceMinor: 1000,
		}},
	}
	handler := registerWalletTestHandler(service)

	request := httptest.NewRequest(http.MethodGet, "/client/wallets?owner_id=other&status=active&limit=10", nil)
	request = request.WithContext(tenant.WithContext(request.Context(), tenant.NewContext("tenant_1")))
	request = request.WithContext(identity.WithActor(request.Context(), identity.NewActor("account_1", "tenant_1", identity.ActorTypeClient)))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", response.Code, response.Body.String())
	}
	if service.listWalletCalls != 1 {
		t.Fatalf("expected list wallets once, got %d", service.listWalletCalls)
	}
	if service.walletFilter.TenantID != tenant.ID("tenant_1") ||
		service.walletFilter.OwnerType != OwnerTypeUser ||
		service.walletFilter.OwnerID != OwnerID("account_1") ||
		service.walletFilter.Status != StatusActive ||
		service.walletFilter.Limit != 10 {
		t.Fatalf("unexpected wallet filter: %+v", service.walletFilter)
	}
	if !strings.Contains(response.Body.String(), `"display_id":70001`) {
		t.Fatalf("expected wallet response, got %s", response.Body.String())
	}
}

func TestHTTPHandlerGetClientWalletUsesOwnerScope(t *testing.T) {
	service := &fakeWalletHTTPService{wallet: Wallet{ID: "wallet_1", DisplayID: 70002, TenantID: "tenant_1"}}
	handler := registerWalletTestHandler(service)

	request := httptest.NewRequest(http.MethodGet, "/client/wallets/wallet_1", nil)
	request = request.WithContext(tenant.WithContext(request.Context(), tenant.NewContext("tenant_1")))
	request = request.WithContext(identity.WithActor(request.Context(), identity.NewActor("account_1", "tenant_1", identity.ActorTypeClient)))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", response.Code, response.Body.String())
	}
	if service.walletLookup.ID != WalletID("wallet_1") ||
		service.walletLookup.OwnerType != OwnerTypeUser ||
		service.walletLookup.OwnerID != OwnerID("account_1") {
		t.Fatalf("unexpected wallet lookup: %+v", service.walletLookup)
	}
}

func TestHTTPHandlerListClientLedgerVerifiesWalletOwner(t *testing.T) {
	service := &fakeWalletHTTPService{
		wallet: Wallet{ID: "wallet_1", DisplayID: 70003, TenantID: "tenant_1"},
		entries: []LedgerEntry{{
			ID:          "entry_1",
			DisplayID:   71001,
			WalletID:    "wallet_1",
			TenantID:    "tenant_1",
			Direction:   DirectionCredit,
			AmountMinor: 1000,
			Currency:    "USD",
			EntryType:   EntryTypeTopup,
			Status:      LedgerStatusPosted,
		}},
	}
	handler := registerWalletTestHandler(service)

	request := httptest.NewRequest(http.MethodGet, "/client/wallets/wallet_1/ledger?direction=credit&entry_type=topup&limit=8", nil)
	request = request.WithContext(tenant.WithContext(request.Context(), tenant.NewContext("tenant_1")))
	request = request.WithContext(identity.WithActor(request.Context(), identity.NewActor("account_1", "tenant_1", identity.ActorTypeClient)))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", response.Code, response.Body.String())
	}
	if service.getWalletCalls != 1 || service.listLedgerCalls != 1 {
		t.Fatalf("expected wallet verification and ledger list, got wallet=%d ledger=%d", service.getWalletCalls, service.listLedgerCalls)
	}
	if service.ledgerFilter.WalletID != WalletID("wallet_1") ||
		service.ledgerFilter.Direction != DirectionCredit ||
		service.ledgerFilter.EntryType != EntryTypeTopup ||
		service.ledgerFilter.Limit != 8 {
		t.Fatalf("unexpected ledger filter: %+v", service.ledgerFilter)
	}
}

func TestHTTPHandlerListAdminWalletsUsesOwnerFilter(t *testing.T) {
	service := &fakeWalletHTTPService{}
	handler := registerWalletTestHandler(service)

	request := httptest.NewRequest(http.MethodGet, "/admin/wallets?owner_type=user&owner_id=account_2", nil)
	request = request.WithContext(tenant.WithContext(request.Context(), tenant.NewContext("tenant_1")))
	request = request.WithContext(identity.WithActor(request.Context(), identity.NewActor("admin_1", "tenant_1", identity.ActorTypeResellerOwner)))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", response.Code, response.Body.String())
	}
	if service.walletFilter.OwnerType != OwnerTypeUser || service.walletFilter.OwnerID != OwnerID("account_2") {
		t.Fatalf("unexpected admin wallet filter: %+v", service.walletFilter)
	}
}

func TestHTTPHandlerRejectsBadLedgerStatus(t *testing.T) {
	service := &fakeWalletHTTPService{wallet: Wallet{ID: "wallet_1"}}
	handler := registerWalletTestHandler(service)

	request := httptest.NewRequest(http.MethodGet, "/client/wallets/wallet_1/ledger?status=bad", nil)
	request = request.WithContext(tenant.WithContext(request.Context(), tenant.NewContext("tenant_1")))
	request = request.WithContext(identity.WithActor(request.Context(), identity.NewActor("account_1", "tenant_1", identity.ActorTypeClient)))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d: %s", response.Code, response.Body.String())
	}
	if service.listLedgerCalls != 0 {
		t.Fatalf("expected no ledger call, got %d", service.listLedgerCalls)
	}
}

func registerWalletTestHandler(service HTTPService) http.Handler {
	mux := http.NewServeMux()
	NewHTTPHandler(service).RegisterRoutes(mux)
	return mux
}

type fakeWalletHTTPService struct {
	wallets          []Wallet
	wallet           Wallet
	entries          []LedgerEntry
	topup            TopupRequest
	topups           []TopupRequest
	walletFilter     WalletFilter
	walletLookup     WalletLookup
	ledgerFilter     LedgerEntryFilter
	topupInput       CreateTopupRequestInput
	topupFilter      TopupRequestFilter
	topupLookup      TopupRequestLookup
	listWalletCalls  int
	getWalletCalls   int
	listLedgerCalls  int
	createTopupCalls int
	listTopupCalls   int
	getTopupCalls    int
}

func (service *fakeWalletHTTPService) ListWallets(ctx context.Context, filter WalletFilter) ([]Wallet, error) {
	service.listWalletCalls++
	service.walletFilter = filter
	return service.wallets, nil
}

func (service *fakeWalletHTTPService) GetWallet(ctx context.Context, lookup WalletLookup) (Wallet, error) {
	service.getWalletCalls++
	service.walletLookup = lookup
	return service.wallet, nil
}

func (service *fakeWalletHTTPService) ListLedgerEntries(ctx context.Context, filter LedgerEntryFilter) ([]LedgerEntry, error) {
	service.listLedgerCalls++
	service.ledgerFilter = filter
	return service.entries, nil
}

func (service *fakeWalletHTTPService) CreateTopupRequest(ctx context.Context, input CreateTopupRequestInput) (TopupRequest, error) {
	service.createTopupCalls++
	service.topupInput = input
	return service.topup, nil
}

func (service *fakeWalletHTTPService) ListTopupRequests(ctx context.Context, filter TopupRequestFilter) ([]TopupRequest, error) {
	service.listTopupCalls++
	service.topupFilter = filter
	return service.topups, nil
}

func (service *fakeWalletHTTPService) GetTopupRequest(ctx context.Context, lookup TopupRequestLookup) (TopupRequest, error) {
	service.getTopupCalls++
	service.topupLookup = lookup
	return service.topup, nil
}
