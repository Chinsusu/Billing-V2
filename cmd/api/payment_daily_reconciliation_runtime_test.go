package main

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	platformdb "github.com/Chinsusu/Billing-V2/internal/platform/db"
)

func TestNewRuntimeWithDSNProtectsAdminDailyReconciliationRoute(t *testing.T) {
	runtime, err := newRuntime(context.Background(), testRuntimeConfig("postgres://billing@localhost/billing"), testRuntimeLogger(), func(ctx context.Context, cfg platformdb.Config) (*sql.DB, error) {
		return newStubDB(), nil
	})
	if err != nil {
		t.Fatalf("newRuntime returned error: %v", err)
	}
	defer closeRuntime(t, runtime)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/admin/daily-reconciliation?date=2026-05-13", nil)
	request.Header.Set("X-Tenant-Id", "tenant_1")
	runtime.api.Handler().ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("expected missing actor to be rejected, got %d: %s", response.Code, response.Body.String())
	}
}
