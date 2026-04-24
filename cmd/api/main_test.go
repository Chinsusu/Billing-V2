package main

import (
	"bytes"
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Chinsusu/Billing-V2/internal/platform/config"
	platformdb "github.com/Chinsusu/Billing-V2/internal/platform/db"
	"github.com/Chinsusu/Billing-V2/internal/platform/logger"
)

func TestNewRuntimeWithoutDSNLeavesDomainRoutesDisabled(t *testing.T) {
	runtime, err := newRuntime(context.Background(), testRuntimeConfig(""), testRuntimeLogger(), func(ctx context.Context, cfg platformdb.Config) (*sql.DB, error) {
		t.Fatal("database opener should not be called without DB_DSN")
		return nil, nil
	})
	if err != nil {
		t.Fatalf("newRuntime returned error: %v", err)
	}
	defer closeRuntime(t, runtime)

	health := httptest.NewRecorder()
	runtime.api.Handler().ServeHTTP(health, httptest.NewRequest(http.MethodGet, "/healthz", nil))
	if health.Code != http.StatusOK {
		t.Fatalf("expected health route to stay enabled, got %d", health.Code)
	}

	catalogResponse := httptest.NewRecorder()
	runtime.api.Handler().ServeHTTP(catalogResponse, httptest.NewRequest(http.MethodGet, "/client/catalog", nil))
	if catalogResponse.Code != http.StatusNotFound {
		t.Fatalf("expected catalog route to be disabled without DB_DSN, got %d", catalogResponse.Code)
	}

	accountResponse := httptest.NewRecorder()
	runtime.api.Handler().ServeHTTP(accountResponse, httptest.NewRequest(http.MethodGet, "/admin/tenants", nil))
	if accountResponse.Code != http.StatusNotFound {
		t.Fatalf("expected account route to be disabled without DB_DSN, got %d", accountResponse.Code)
	}

	orderResponse := httptest.NewRecorder()
	runtime.api.Handler().ServeHTTP(orderResponse, httptest.NewRequest(http.MethodPost, "/client/orders", strings.NewReader(`{}`)))
	if orderResponse.Code != http.StatusNotFound {
		t.Fatalf("expected order route to be disabled without DB_DSN, got %d", orderResponse.Code)
	}

	checkoutResponse := httptest.NewRecorder()
	runtime.api.Handler().ServeHTTP(checkoutResponse, httptest.NewRequest(http.MethodPost, "/client/checkouts", strings.NewReader(`{}`)))
	if checkoutResponse.Code != http.StatusNotFound {
		t.Fatalf("expected checkout route to be disabled without DB_DSN, got %d", checkoutResponse.Code)
	}

	paymentResponse := httptest.NewRecorder()
	runtime.api.Handler().ServeHTTP(paymentResponse, httptest.NewRequest(http.MethodGet, "/client/transactions", nil))
	if paymentResponse.Code != http.StatusNotFound {
		t.Fatalf("expected payment route to be disabled without DB_DSN, got %d", paymentResponse.Code)
	}

	walletResponse := httptest.NewRecorder()
	runtime.api.Handler().ServeHTTP(walletResponse, httptest.NewRequest(http.MethodGet, "/client/wallets", nil))
	if walletResponse.Code != http.StatusNotFound {
		t.Fatalf("expected wallet route to be disabled without DB_DSN, got %d", walletResponse.Code)
	}
}

func TestNewRuntimeWithDSNRegistersCatalogRoutes(t *testing.T) {
	var opened platformdb.Config
	runtime, err := newRuntime(context.Background(), testRuntimeConfig("postgres://billing@localhost/billing"), testRuntimeLogger(), func(ctx context.Context, cfg platformdb.Config) (*sql.DB, error) {
		opened = cfg
		return newStubDB(), nil
	})
	if err != nil {
		t.Fatalf("newRuntime returned error: %v", err)
	}
	defer closeRuntime(t, runtime)

	if opened.DriverName != platformdb.DefaultDriverName {
		t.Fatalf("expected default driver, got %q", opened.DriverName)
	}
	if opened.DSN == "" {
		t.Fatal("expected DSN to be passed to opener")
	}

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/client/catalog", nil)
	runtime.api.Handler().ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected registered catalog route to validate tenant context, got %d: %s", response.Code, response.Body.String())
	}
	if !strings.Contains(response.Body.String(), "tenant.context_missing") {
		t.Fatalf("expected tenant validation response, got %s", response.Body.String())
	}
}

func TestNewRuntimeWithDSNRegistersOrderRoutes(t *testing.T) {
	runtime, err := newRuntime(context.Background(), testRuntimeConfig("postgres://billing@localhost/billing"), testRuntimeLogger(), func(ctx context.Context, cfg platformdb.Config) (*sql.DB, error) {
		return newStubDB(), nil
	})
	if err != nil {
		t.Fatalf("newRuntime returned error: %v", err)
	}
	defer closeRuntime(t, runtime)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/client/orders", strings.NewReader(`{}`))
	runtime.api.Handler().ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected registered order route to validate tenant context, got %d: %s", response.Code, response.Body.String())
	}
	if !strings.Contains(response.Body.String(), "tenant.context_missing") {
		t.Fatalf("expected tenant validation response, got %s", response.Body.String())
	}
}

func TestNewRuntimeWithDSNRegistersCheckoutRoutes(t *testing.T) {
	runtime, err := newRuntime(context.Background(), testRuntimeConfig("postgres://billing@localhost/billing"), testRuntimeLogger(), func(ctx context.Context, cfg platformdb.Config) (*sql.DB, error) {
		return newStubDB(), nil
	})
	if err != nil {
		t.Fatalf("newRuntime returned error: %v", err)
	}
	defer closeRuntime(t, runtime)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/client/checkouts", strings.NewReader(`{}`))
	runtime.api.Handler().ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected registered checkout route to validate tenant context, got %d: %s", response.Code, response.Body.String())
	}
	if !strings.Contains(response.Body.String(), "tenant.context_missing") {
		t.Fatalf("expected tenant validation response, got %s", response.Body.String())
	}
}

func TestNewRuntimeWithDSNRegistersPaymentRoutes(t *testing.T) {
	runtime, err := newRuntime(context.Background(), testRuntimeConfig("postgres://billing@localhost/billing"), testRuntimeLogger(), func(ctx context.Context, cfg platformdb.Config) (*sql.DB, error) {
		return newStubDB(), nil
	})
	if err != nil {
		t.Fatalf("newRuntime returned error: %v", err)
	}
	defer closeRuntime(t, runtime)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/client/transactions", nil)
	runtime.api.Handler().ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected registered payment route to validate tenant context, got %d: %s", response.Code, response.Body.String())
	}
	if !strings.Contains(response.Body.String(), "tenant.context_missing") {
		t.Fatalf("expected tenant validation response, got %s", response.Body.String())
	}
}

func TestNewRuntimeWithDSNRegistersWalletRoutes(t *testing.T) {
	runtime, err := newRuntime(context.Background(), testRuntimeConfig("postgres://billing@localhost/billing"), testRuntimeLogger(), func(ctx context.Context, cfg platformdb.Config) (*sql.DB, error) {
		return newStubDB(), nil
	})
	if err != nil {
		t.Fatalf("newRuntime returned error: %v", err)
	}
	defer closeRuntime(t, runtime)

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/client/wallets", nil)
	runtime.api.Handler().ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected registered wallet route to validate tenant context, got %d: %s", response.Code, response.Body.String())
	}
	if !strings.Contains(response.Body.String(), "tenant.context_missing") {
		t.Fatalf("expected tenant validation response, got %s", response.Body.String())
	}
}

func TestNewRuntimeReturnsDatabaseOpenError(t *testing.T) {
	expected := errors.New("dial failed")
	_, err := newRuntime(context.Background(), testRuntimeConfig("postgres://billing@localhost/billing"), testRuntimeLogger(), func(ctx context.Context, cfg platformdb.Config) (*sql.DB, error) {
		return nil, expected
	})

	if !errors.Is(err, expected) {
		t.Fatalf("expected wrapped database error, got %v", err)
	}
	if err == nil || !strings.Contains(err.Error(), "open api database") {
		t.Fatalf("expected clear database open error, got %v", err)
	}
}

func TestNewCatalogRoutesReturnsRegistrar(t *testing.T) {
	registrar := newCatalogRoutes(newStubDB())
	if registrar == nil {
		t.Fatal("expected catalog route registrar")
	}

	mux := http.NewServeMux()
	registrar.RegisterRoutes(mux)
	response := httptest.NewRecorder()
	mux.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/client/catalog", nil))
	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected catalog route to be registered, got %d", response.Code)
	}
}

func TestNewOrderRoutesReturnsRegistrar(t *testing.T) {
	registrar := newOrderRoutes(newStubDB())
	if registrar == nil {
		t.Fatal("expected order route registrar")
	}

	mux := http.NewServeMux()
	registrar.RegisterRoutes(mux)
	response := httptest.NewRecorder()
	mux.ServeHTTP(response, httptest.NewRequest(http.MethodPost, "/client/orders", nil))
	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected order route to be registered, got %d", response.Code)
	}
}

func TestNewCheckoutRoutesReturnsRegistrar(t *testing.T) {
	registrar := newCheckoutRoutes(newStubDB())
	if registrar == nil {
		t.Fatal("expected checkout route registrar")
	}

	mux := http.NewServeMux()
	registrar.RegisterRoutes(mux)
	response := httptest.NewRecorder()
	mux.ServeHTTP(response, httptest.NewRequest(http.MethodPost, "/client/checkouts", nil))
	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected checkout route to be registered, got %d", response.Code)
	}
}

func TestNewPaymentRoutesReturnsRegistrar(t *testing.T) {
	registrar := newPaymentRoutes(newStubDB())
	if registrar == nil {
		t.Fatal("expected payment route registrar")
	}

	mux := http.NewServeMux()
	registrar.RegisterRoutes(mux)
	response := httptest.NewRecorder()
	mux.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/client/transactions", nil))
	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected payment route to be registered, got %d", response.Code)
	}

	reconciliationResponse := httptest.NewRecorder()
	mux.ServeHTTP(reconciliationResponse, httptest.NewRequest(http.MethodGet, "/admin/payment-reconciliation", nil))
	if reconciliationResponse.Code != http.StatusBadRequest {
		t.Fatalf("expected payment reconciliation route to be registered, got %d", reconciliationResponse.Code)
	}
}

func TestNewWalletRoutesReturnsRegistrar(t *testing.T) {
	registrar := newWalletRoutes(newStubDB())
	if registrar == nil {
		t.Fatal("expected wallet route registrar")
	}

	mux := http.NewServeMux()
	registrar.RegisterRoutes(mux)
	response := httptest.NewRecorder()
	mux.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/client/wallets", nil))
	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected wallet route to be registered, got %d", response.Code)
	}
}

func closeRuntime(t *testing.T, runtime *apiRuntime) {
	t.Helper()
	if err := runtime.close(); err != nil {
		t.Fatalf("runtime close returned error: %v", err)
	}
}

func testRuntimeConfig(dsn string) config.Config {
	return config.Config{
		AppEnvironment: config.EnvironmentLocal,
		AppName:        "billing-v2",
		HTTPAddr:       ":8080",
		LogLevel:       config.LogLevelDebug,
		DatabaseDSN:    dsn,
	}
}

func testRuntimeLogger() *logger.Logger {
	return logger.New(&bytes.Buffer{}, config.LogLevelDebug)
}

func newStubDB() *sql.DB {
	return sql.OpenDB(stubConnector{})
}

type stubConnector struct{}

func (stubConnector) Connect(ctx context.Context) (driver.Conn, error) {
	return stubConn{}, nil
}

func (stubConnector) Driver() driver.Driver {
	return stubDriver{}
}

type stubDriver struct{}

func (stubDriver) Open(name string) (driver.Conn, error) {
	return stubConn{}, nil
}

type stubConn struct{}

func (stubConn) Prepare(query string) (driver.Stmt, error) {
	return nil, errors.New("prepare is not implemented")
}

func (stubConn) Close() error {
	return nil
}

func (stubConn) Begin() (driver.Tx, error) {
	return nil, errors.New("begin is not implemented")
}
