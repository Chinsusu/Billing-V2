package identity

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

func TestHeaderActorMiddlewareAddsActorContext(t *testing.T) {
	var captured Actor
	var ok bool
	handler := HeaderActorMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		captured, ok = FromContext(r.Context())
	}))

	request := httptest.NewRequest(http.MethodGet, "/auth", nil)
	request.Header.Set(HeaderActorID, " user_1 ")
	request.Header.Set(HeaderActorType, string(ActorTypeResellerOwner))
	request.Header.Set(HeaderActorTenantID, "tenant_a")
	request.Header.Set(HeaderActorRoleIDs, "role_a, role_b")
	handler.ServeHTTP(httptest.NewRecorder(), request)

	if !ok {
		t.Fatal("expected actor context")
	}
	if captured.ID != "user_1" || captured.TenantID != "tenant_a" || captured.Type != ActorTypeResellerOwner {
		t.Fatalf("unexpected actor: %+v", captured)
	}
	if !captured.HasRole("role_a") || !captured.HasRole("role_b") {
		t.Fatalf("expected parsed role ids, got %+v", captured.RoleIDs)
	}
}

func TestHeaderActorMiddlewareFallsBackToTenantContext(t *testing.T) {
	var captured Actor
	handler := HeaderActorMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		captured, _ = RequireActor(r.Context())
	}))

	request := httptest.NewRequest(http.MethodGet, "/auth", nil)
	request = request.WithContext(tenant.WithContext(request.Context(), tenant.NewContext("tenant_from_context")))
	request.Header.Set(HeaderActorID, "user_1")
	request.Header.Set(HeaderActorType, string(ActorTypeClient))
	handler.ServeHTTP(httptest.NewRecorder(), request)

	if captured.TenantID != "tenant_from_context" {
		t.Fatalf("expected tenant fallback, got %q", captured.TenantID)
	}
}

func TestRequireActorRejectsMissingContext(t *testing.T) {
	_, err := RequireActor(context.Background())
	if !errors.Is(err, ErrActorContextMissing) {
		t.Fatalf("expected actor context missing, got %v", err)
	}
}
