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

func TestNewRuntimeWithDSNProtectsAdminCatalogRoutes(t *testing.T) {
	assertRuntimeRejectsMissingActor(t, http.MethodPost, "/admin/catalog/products", `{"product_type":"vps","name":"VPS"}`)
}

func TestNewRuntimeWithDSNProtectsAdminAccountRoutes(t *testing.T) {
	assertRuntimeRejectsMissingActor(t, http.MethodGet, "/admin/tenants", "")
}

func TestNewRuntimeWithDSNProtectsResellerCustomerRoutes(t *testing.T) {
	assertRuntimeRejectsMissingActor(t, http.MethodGet, "/reseller/customers", "")
}

func TestNewRuntimeWithDSNProtectsClientOrderRoutes(t *testing.T) {
	assertRuntimeRejectsMissingActor(t, http.MethodPost, "/client/orders", `{"tenant_plan_id":"tenant_plan_1"}`)
}

func TestNewRuntimeWithDSNProtectsClientCheckoutRoutes(t *testing.T) {
	assertRuntimeRejectsMissingActor(t, http.MethodPost, "/client/checkouts", `{"order_id":"order_1"}`)
}

func TestNewRuntimeWithDSNProtectsAdminOrderRoutes(t *testing.T) {
	assertRuntimeRejectsMissingActor(t, http.MethodGet, "/admin/orders", "")
}

func TestNewRuntimeWithDSNProtectsResellerOrderRoutes(t *testing.T) {
	assertRuntimeRejectsMissingActor(t, http.MethodGet, "/reseller/orders", "")
}

func TestNewRuntimeWithDSNProtectsAdminOrderStatusRoute(t *testing.T) {
	assertRuntimeRejectsMissingActor(t, http.MethodPatch, "/admin/orders/order_1/status", `{"from_status":"pending_payment","to_status":"paid","billing_status":"paid"}`)
}

func TestNewRuntimeWithDSNProtectsClientServiceRoutes(t *testing.T) {
	assertRuntimeRejectsMissingActor(t, http.MethodGet, "/client/services", "")
}

func TestNewRuntimeWithDSNProtectsResellerServiceRoutes(t *testing.T) {
	assertRuntimeRejectsMissingActor(t, http.MethodGet, "/reseller/services", "")
}

func TestNewRuntimeWithDSNProtectsAdminServiceRoutes(t *testing.T) {
	assertRuntimeRejectsMissingActor(t, http.MethodGet, "/admin/services", "")
}

func TestNewRuntimeWithDSNProtectsClientCredentialRevealRoutes(t *testing.T) {
	assertRuntimeRejectsMissingActor(t, http.MethodPost, "/client/services/service_1/credentials/credential_1/reveal", `{}`)
}

func TestNewRuntimeWithDSNProtectsResellerCredentialRevealRoutes(t *testing.T) {
	assertRuntimeRejectsMissingActor(t, http.MethodPost, "/reseller/services/service_1/credentials/credential_1/reveal", `{}`)
}

func TestNewRuntimeWithDSNProtectsAdminCredentialRevealRoutes(t *testing.T) {
	assertRuntimeRejectsMissingActor(t, http.MethodPost, "/admin/services/service_1/credentials/credential_1/reveal", `{}`)
}

func TestNewRuntimeWithDSNProtectsClientPaymentRoutes(t *testing.T) {
	assertRuntimeRejectsMissingActor(t, http.MethodGet, "/client/transactions", "")
}

func TestNewRuntimeWithDSNProtectsAdminPaymentRoutes(t *testing.T) {
	assertRuntimeRejectsMissingActor(t, http.MethodGet, "/admin/transactions", "")
}

func TestNewRuntimeWithDSNProtectsResellerPaymentRoutes(t *testing.T) {
	assertRuntimeRejectsMissingActor(t, http.MethodGet, "/reseller/transactions", "")
}

func TestNewRuntimeWithDSNProtectsAdminPaymentReconciliationRoutes(t *testing.T) {
	assertRuntimeRejectsMissingActor(t, http.MethodGet, "/admin/payment-reconciliation", "")
}

func TestNewRuntimeWithDSNProtectsResellerInvoiceRoutes(t *testing.T) {
	assertRuntimeRejectsMissingActor(t, http.MethodGet, "/reseller/invoices", "")
}

func TestNewRuntimeWithDSNProtectsClientWalletRoutes(t *testing.T) {
	assertRuntimeRejectsMissingActor(t, http.MethodGet, "/client/wallets", "")
}

func TestNewRuntimeWithDSNProtectsResellerWalletRoutes(t *testing.T) {
	assertRuntimeRejectsMissingActor(t, http.MethodGet, "/reseller/wallets", "")
}

func TestNewRuntimeWithDSNProtectsResellerTopupRoutes(t *testing.T) {
	assertRuntimeRejectsMissingActor(t, http.MethodGet, "/reseller/topup-requests", "")
}

func assertRuntimeRejectsMissingActor(t *testing.T, method string, path string, body string) {
	t.Helper()
	runtime, err := newRuntime(context.Background(), testRuntimeConfig("postgres://billing@localhost/billing"), testRuntimeLogger(), func(ctx context.Context, cfg platformdb.Config) (*sql.DB, error) {
		return newStubDB(), nil
	})
	if err != nil {
		t.Fatalf("newRuntime returned error: %v", err)
	}
	defer closeRuntime(t, runtime)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(method, path, strings.NewReader(body))
	request.Header.Set("X-Tenant-Id", "tenant_1")
	runtime.api.Handler().ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("expected missing actor to be rejected, got %d: %s", response.Code, response.Body.String())
	}
}
