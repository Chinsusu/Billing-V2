package order

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
	"github.com/Chinsusu/Billing-V2/internal/modules/wallet"
)

func TestHTTPHandlerRenewClientServiceUsesTenantActorWalletAndIdempotency(t *testing.T) {
	service := &fakeOrderHTTPService{}
	handler := registerOrderTestHandler(service)

	request := httptest.NewRequest(http.MethodPost, "/client/services/service_1/renew", strings.NewReader(`{
		"wallet_id": "wallet_1",
		"from_status": "suspended",
		"reason": "renew from portal"
	}`))
	request.Header.Set(IdempotencyKeyHeader, "renew-key-1")
	request = request.WithContext(tenant.WithContext(request.Context(), tenant.NewContext("tenant_1")))
	request = request.WithContext(identity.WithActor(request.Context(), identity.NewActor("buyer_1", "tenant_1", identity.ActorTypeClient)))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", response.Code, response.Body.String())
	}
	if service.renewClientServiceCalls != 1 {
		t.Fatalf("expected renew once, got %d", service.renewClientServiceCalls)
	}
	input := service.renewClientServiceInput
	if input.TenantID != tenant.ID("tenant_1") ||
		input.BuyerUserID != identity.UserID("buyer_1") ||
		input.ActorID != identity.UserID("buyer_1") ||
		input.ServiceID != ServiceID("service_1") ||
		input.WalletID != wallet.WalletID("wallet_1") ||
		input.FromStatus != ServiceStatusSuspended ||
		input.IdempotencyKey != IdempotencyKey("renew-key-1") {
		t.Fatalf("unexpected renewal input: %+v", input)
	}
	body := response.Body.String()
	for _, expected := range []string{`"invoice"`, `"status":"paid"`, `"ledger"`, `"renewed":true`} {
		if !strings.Contains(body, expected) {
			t.Fatalf("expected %s in response, got %s", expected, body)
		}
	}
}

func TestHTTPHandlerRenewClientServiceRejectsMissingIdempotencyKey(t *testing.T) {
	service := &fakeOrderHTTPService{}
	handler := registerOrderTestHandler(service)

	request := httptest.NewRequest(http.MethodPost, "/client/services/service_1/renew", strings.NewReader(`{
		"wallet_id": "wallet_1",
		"from_status": "active"
	}`))
	request = request.WithContext(tenant.WithContext(request.Context(), tenant.NewContext("tenant_1")))
	request = request.WithContext(identity.WithActor(request.Context(), identity.NewActor("buyer_1", "tenant_1", identity.ActorTypeClient)))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d: %s", response.Code, response.Body.String())
	}
	if service.renewClientServiceCalls != 0 {
		t.Fatalf("expected no renew call, got %d", service.renewClientServiceCalls)
	}
	if !strings.Contains(response.Body.String(), "order.idempotency_key_missing") {
		t.Fatalf("expected idempotency validation error, got %s", response.Body.String())
	}
}

func TestHTTPHandlerRenewClientServiceMapsInsufficientBalance(t *testing.T) {
	service := &fakeOrderHTTPService{renewClientServiceError: wallet.ErrInsufficientBalance}
	handler := registerOrderTestHandler(service)

	request := httptest.NewRequest(http.MethodPost, "/client/services/service_1/renew", strings.NewReader(`{
		"wallet_id": "wallet_1",
		"from_status": "active"
	}`))
	request.Header.Set(IdempotencyKeyHeader, "renew-key-1")
	request = request.WithContext(tenant.WithContext(request.Context(), tenant.NewContext("tenant_1")))
	request = request.WithContext(identity.WithActor(request.Context(), identity.NewActor("buyer_1", "tenant_1", identity.ActorTypeClient)))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusConflict {
		t.Fatalf("expected status 409, got %d: %s", response.Code, response.Body.String())
	}
	if !strings.Contains(response.Body.String(), "wallet.insufficient_balance") {
		t.Fatalf("expected wallet error, got %s", response.Body.String())
	}
}

func TestHTTPHandlerRenewClientServiceMiddlewareRunsBeforeService(t *testing.T) {
	service := &fakeOrderHTTPService{}
	mux := http.NewServeMux()
	NewHTTPHandlerWithOptions(service, HTTPHandlerOptions{
		ClientServiceRenewMiddleware: func(next http.HandlerFunc) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusForbidden)
			}
		},
	}).RegisterRoutes(mux)

	request := httptest.NewRequest(http.MethodPost, "/client/services/service_1/renew", strings.NewReader(`{}`))
	request = request.WithContext(tenant.WithContext(request.Context(), tenant.NewContext("tenant_1")))
	response := httptest.NewRecorder()

	mux.ServeHTTP(response, request)

	if response.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d", response.Code)
	}
	if service.renewClientServiceCalls != 0 {
		t.Fatalf("expected service not to run, got %d calls", service.renewClientServiceCalls)
	}
}
