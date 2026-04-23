package app

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Chinsusu/Billing-V2/internal/platform/config"
	"github.com/Chinsusu/Billing-V2/internal/platform/httpserver"
	"github.com/Chinsusu/Billing-V2/internal/platform/logger"
)

func TestHealthEndpointReturnsSuccessEnvelope(t *testing.T) {
	api := newTestAPI(t)

	request := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	request.Header.Set(httpserver.RequestIDHeader, "req_health")
	response := httptest.NewRecorder()

	api.Handler().ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.Code)
	}

	var body struct {
		Data      HealthResponse `json:"data"`
		RequestID string         `json:"request_id"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("response is not JSON: %v", err)
	}
	if body.RequestID != "req_health" {
		t.Fatalf("expected request id, got %q", body.RequestID)
	}
	if body.Data.Status != "ok" {
		t.Fatalf("expected ok status, got %q", body.Data.Status)
	}
}

func TestReadyEndpointReturnsReadyStatus(t *testing.T) {
	api := newTestAPI(t)

	request := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	response := httptest.NewRecorder()

	api.Handler().ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.Code)
	}
}

func TestHealthEndpointRejectsUnsupportedMethod(t *testing.T) {
	api := newTestAPI(t)

	request := httptest.NewRequest(http.MethodPost, "/healthz", nil)
	response := httptest.NewRecorder()

	api.Handler().ServeHTTP(response, request)

	if response.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405, got %d", response.Code)
	}
}

func TestNewAPIWithOptionsRegistersOptionalRoutes(t *testing.T) {
	api, err := NewAPIWithOptions(testAPIConfig(), logger.New(&bytes.Buffer{}, config.LogLevelDebug), APIOptions{
		CatalogRoutes: testRouteRegistrar{},
		OrderRoutes:   testOrderRouteRegistrar{},
		PaymentRoutes: testPaymentRouteRegistrar{},
		WalletRoutes:  testWalletRouteRegistrar{},
	})
	if err != nil {
		t.Fatalf("NewAPIWithOptions returned error: %v", err)
	}

	request := httptest.NewRequest(http.MethodGet, "/catalog-test", nil)
	response := httptest.NewRecorder()

	api.Handler().ServeHTTP(response, request)

	if response.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", response.Code)
	}

	paymentResponse := httptest.NewRecorder()
	api.Handler().ServeHTTP(paymentResponse, httptest.NewRequest(http.MethodGet, "/payment-test", nil))
	if paymentResponse.Code != http.StatusNoContent {
		t.Fatalf("expected payment route status 204, got %d", paymentResponse.Code)
	}

	walletResponse := httptest.NewRecorder()
	api.Handler().ServeHTTP(walletResponse, httptest.NewRequest(http.MethodGet, "/wallet-test", nil))
	if walletResponse.Code != http.StatusNoContent {
		t.Fatalf("expected wallet route status 204, got %d", walletResponse.Code)
	}
}

type testRouteRegistrar struct{}

func (testRouteRegistrar) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/catalog-test", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNoContent) })
}

type testOrderRouteRegistrar struct{}

func (testOrderRouteRegistrar) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/order-test", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNoContent) })
}

type testPaymentRouteRegistrar struct{}

func (testPaymentRouteRegistrar) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/payment-test", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNoContent) })
}

type testWalletRouteRegistrar struct{}

func (testWalletRouteRegistrar) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/wallet-test", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNoContent) })
}

func newTestAPI(t *testing.T) *API {
	t.Helper()

	api, err := NewAPI(testAPIConfig(), logger.New(&bytes.Buffer{}, config.LogLevelDebug))
	if err != nil {
		t.Fatalf("NewAPI returned error: %v", err)
	}
	return api
}

func testAPIConfig() config.Config {
	return config.Config{
		AppEnvironment: config.EnvironmentLocal,
		AppName:        "billing-v2",
		HTTPAddr:       ":8080",
		LogLevel:       config.LogLevelDebug,
	}
}
