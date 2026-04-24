package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewAccountRoutesReturnsRegistrar(t *testing.T) {
	registrar := newAccountRoutes(newStubDB())
	if registrar == nil {
		t.Fatal("expected account route registrar")
	}

	mux := http.NewServeMux()
	registrar.RegisterRoutes(mux)
	response := httptest.NewRecorder()
	mux.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/admin/tenants", nil))
	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected account route to be registered, got %d", response.Code)
	}
	if !strings.Contains(response.Body.String(), "tenant.context_missing") {
		t.Fatalf("expected tenant validation response, got %s", response.Body.String())
	}
}
