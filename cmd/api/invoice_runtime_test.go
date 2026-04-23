package main

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	platformdb "github.com/Chinsusu/Billing-V2/internal/platform/db"
)

func TestNewRuntimeWithDSNRegistersInvoiceRoutes(t *testing.T) {
	runtime, err := newRuntime(context.Background(), testRuntimeConfig("postgres://billing@localhost/billing"), testRuntimeLogger(), func(ctx context.Context, cfg platformdb.Config) (*sql.DB, error) {
		return newStubDB(), nil
	})
	if err != nil {
		t.Fatalf("newRuntime returned error: %v", err)
	}
	defer closeRuntime(t, runtime)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/client/invoices", nil)
	runtime.api.Handler().ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected registered invoice route to validate tenant context, got %d: %s", response.Code, response.Body.String())
	}
	if !strings.Contains(response.Body.String(), "tenant.context_missing") {
		t.Fatalf("expected tenant validation response, got %s", response.Body.String())
	}
}

func TestNewRuntimeWithDSNProtectsClientInvoiceRoutes(t *testing.T) {
	runtime, err := newRuntime(context.Background(), testRuntimeConfig("postgres://billing@localhost/billing"), testRuntimeLogger(), func(ctx context.Context, cfg platformdb.Config) (*sql.DB, error) {
		return newStubDB(), nil
	})
	if err != nil {
		t.Fatalf("newRuntime returned error: %v", err)
	}
	defer closeRuntime(t, runtime)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/client/invoices", nil)
	request.Header.Set("X-Tenant-Id", "tenant_1")
	runtime.api.Handler().ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("expected missing actor to be rejected, got %d: %s", response.Code, response.Body.String())
	}
}

func TestNewInvoiceRoutesReturnsRegistrar(t *testing.T) {
	registrar := newInvoiceRoutes(newStubDB())
	if registrar == nil {
		t.Fatal("expected invoice route registrar")
	}

	mux := http.NewServeMux()
	registrar.RegisterRoutes(mux)
	response := httptest.NewRecorder()
	mux.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/client/invoices", nil))
	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected invoice route to be registered, got %d", response.Code)
	}
}
