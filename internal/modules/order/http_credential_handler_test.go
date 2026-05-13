package order

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/catalog"
	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

func TestHTTPHandlerGetClientServiceIncludesMaskedCredentialMetadata(t *testing.T) {
	service := &fakeCredentialOrderHTTPService{
		fakeOrderHTTPService: fakeOrderHTTPService{
			service: ServiceInstance{
				ID:               "service_1",
				DisplayID:        50001,
				TenantID:         "tenant_1",
				OrderID:          "order_1",
				TenantPlanID:     catalog.TenantPlanID("tenant_plan_1"),
				ProviderSourceID: catalog.ProviderSourceID("source_1"),
				Status:           ServiceStatusActive,
				BillingStatus:    BillingStatusPaid,
			},
		},
		credentials: []ServiceCredential{{
			ID:               "credential_1",
			TenantID:         "tenant_1",
			ServiceID:        "service_1",
			Type:             CredentialTypeVPSRoot,
			EncryptedPayload: "encrypted-fixture",
			MaskedHint:       "root / ****",
			Status:           CredentialStatusActive,
		}},
	}
	handler := registerOrderTestHandler(service)

	request := httptest.NewRequest(http.MethodGet, "/client/services/service_1", nil)
	request = request.WithContext(tenant.WithContext(request.Context(), tenant.NewContext("tenant_1")))
	request = request.WithContext(identity.WithActor(request.Context(), identity.NewActor("buyer_1", "tenant_1", identity.ActorTypeClient)))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", response.Code, response.Body.String())
	}
	if service.serviceLookup.BuyerUserID != identity.UserID("buyer_1") {
		t.Fatalf("expected client owner scope, got %+v", service.serviceLookup)
	}
	if service.credentialFilter.TenantID != tenant.ID("tenant_1") ||
		service.credentialFilter.ServiceID != ServiceID("service_1") ||
		service.credentialFilter.Status != CredentialStatusActive {
		t.Fatalf("unexpected credential filter: %+v", service.credentialFilter)
	}
	body := response.Body.String()
	if !strings.Contains(body, `"masked_hint":"root / ****"`) || !strings.Contains(body, `"credential_1"`) {
		t.Fatalf("expected masked credential metadata, got %s", body)
	}
	if strings.Contains(body, "encrypted-fixture") {
		t.Fatalf("detail response must not expose encrypted payload: %s", body)
	}
}

func TestHTTPHandlerRevealClientServiceCredentialUsesOwnerScope(t *testing.T) {
	revealedAt := time.Date(2026, 5, 13, 9, 30, 0, 0, time.UTC)
	service := &fakeCredentialOrderHTTPService{
		revealResult: RevealServiceCredentialResult{
			Credential: ServiceCredential{
				ID:         "credential_1",
				Type:       CredentialTypeVPSRoot,
				MaskedHint: "root / ****",
				Status:     CredentialStatusActive,
			},
			Payload:              []byte(`{"username":"root","password":"fixture-access"}`),
			RevealedAt:           revealedAt,
			RevealExpiresMessage: "shown once",
		},
	}
	handler := registerOrderTestHandler(service)

	request := httptest.NewRequest(http.MethodPost, "/client/services/service_1/credentials/credential_1/reveal", strings.NewReader(`{"reason":"support handoff"}`))
	request.Header.Set("X-Forwarded-For", "203.0.113.10, 198.51.100.1")
	request.Header.Set("User-Agent", "billing-test")
	request = request.WithContext(tenant.WithContext(request.Context(), tenant.NewContext("tenant_1")))
	request = request.WithContext(identity.WithActor(request.Context(), identity.NewActor("buyer_1", "tenant_1", identity.ActorTypeClient)))
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d: %s", response.Code, response.Body.String())
	}
	if response.Header().Get("Cache-Control") != "no-store" || response.Header().Get("Pragma") != "no-cache" {
		t.Fatalf("expected no-store reveal response headers, got %v", response.Header())
	}
	if service.revealInput.TenantID != tenant.ID("tenant_1") ||
		service.revealInput.ServiceID != ServiceID("service_1") ||
		service.revealInput.CredentialID != CredentialID("credential_1") ||
		service.revealInput.ActorID != identity.UserID("buyer_1") ||
		service.revealInput.BuyerUserID != identity.UserID("buyer_1") ||
		service.revealInput.ClientIP != "203.0.113.10" ||
		service.revealInput.UserAgent != "billing-test" ||
		service.revealInput.Reason != "support handoff" {
		t.Fatalf("unexpected reveal input: %+v", service.revealInput)
	}
	if !strings.Contains(response.Body.String(), `"reveal_expires_message":"shown once"`) {
		t.Fatalf("expected reveal response, got %s", response.Body.String())
	}
}

func TestHTTPHandlerResellerServiceCredentialMiddlewareRunsBeforeReveal(t *testing.T) {
	service := &fakeCredentialOrderHTTPService{}
	mux := http.NewServeMux()
	NewHTTPHandlerWithOptions(service, HTTPHandlerOptions{
		ResellerCredentialMiddleware: func(next http.HandlerFunc) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusForbidden)
			}
		},
	}).RegisterRoutes(mux)

	request := httptest.NewRequest(http.MethodPost, "/reseller/services/service_1/credentials/credential_1/reveal", strings.NewReader(`{}`))
	request = request.WithContext(tenant.WithContext(request.Context(), tenant.NewContext("tenant_1")))
	response := httptest.NewRecorder()

	mux.ServeHTTP(response, request)

	if response.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d", response.Code)
	}
	if service.revealCalls != 0 {
		t.Fatalf("expected reveal not to run, got %d calls", service.revealCalls)
	}
}

type fakeCredentialOrderHTTPService struct {
	fakeOrderHTTPService

	credentials      []ServiceCredential
	credentialFilter ServiceCredentialFilter
	revealCalls      int
	revealInput      RevealServiceCredentialInput
	revealResult     RevealServiceCredentialResult
	revealError      error
}

func (service *fakeCredentialOrderHTTPService) ListServiceCredentials(_ context.Context, filter ServiceCredentialFilter) ([]ServiceCredential, error) {
	service.credentialFilter = filter
	return service.credentials, nil
}

func (service *fakeCredentialOrderHTTPService) RevealServiceCredential(_ context.Context, input RevealServiceCredentialInput) (RevealServiceCredentialResult, error) {
	service.revealCalls++
	service.revealInput = input
	if service.revealError != nil {
		return RevealServiceCredentialResult{}, service.revealError
	}
	return service.revealResult, nil
}
