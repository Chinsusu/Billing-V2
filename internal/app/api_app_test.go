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

func newTestAPI(t *testing.T) *API {
	t.Helper()

	cfg := config.Config{
		AppEnvironment: config.EnvironmentLocal,
		AppName:        "billing-v2",
		HTTPAddr:       ":8080",
		LogLevel:       config.LogLevelDebug,
	}
	api, err := NewAPI(cfg, logger.New(&bytes.Buffer{}, config.LogLevelDebug))
	if err != nil {
		t.Fatalf("NewAPI returned error: %v", err)
	}
	return api
}
