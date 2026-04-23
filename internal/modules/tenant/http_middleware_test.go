package tenant

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHeaderContextMiddlewareAddsTenantContext(t *testing.T) {
	var captured Context
	var ok bool
	handler := HeaderContextMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		captured, ok = FromContext(r.Context())
	}))

	request := httptest.NewRequest(http.MethodGet, "/tenant", nil)
	request.Header.Set(HeaderTenantID, " tenant_a ")
	handler.ServeHTTP(httptest.NewRecorder(), request)

	if !ok {
		t.Fatal("expected tenant context")
	}
	if captured.ActorTenantID != "tenant_a" || captured.EffectiveTenantID != "tenant_a" {
		t.Fatalf("expected tenant_a context, got %+v", captured)
	}
}

func TestHeaderContextMiddlewareLeavesMissingTenantUnset(t *testing.T) {
	var ok bool
	handler := HeaderContextMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, ok = FromContext(r.Context())
	}))

	handler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/tenant", nil))

	if ok {
		t.Fatal("expected missing tenant context to stay unset")
	}
}

func TestHeaderContextMiddlewareSupportsCustomHeader(t *testing.T) {
	var captured Context
	handler := HeaderContextMiddlewareWithOptions(HeaderContextOptions{TenantHeader: "X-Test-Tenant"}, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		captured, _ = RequireContext(r.Context())
	}))

	request := httptest.NewRequest(http.MethodGet, "/tenant", nil)
	request.Header.Set("X-Test-Tenant", "tenant_custom")
	handler.ServeHTTP(httptest.NewRecorder(), request)

	if captured.EffectiveTenantID != "tenant_custom" {
		t.Fatalf("expected custom tenant, got %q", captured.EffectiveTenantID)
	}
}
