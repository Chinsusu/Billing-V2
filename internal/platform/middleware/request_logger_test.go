package middleware

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

func TestRequestLoggerWritesStructuredRequestLog(t *testing.T) {
	var logOutput bytes.Buffer
	log := logger.New(&logOutput, config.LogLevelDebug)
	handler := httpserver.WithRequestID(RequestLogger(log)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		httpserver.WriteSuccess(w, r, http.StatusCreated, map[string]string{"status": "created"})
	})))

	request := httptest.NewRequest(http.MethodPost, "/orders", nil)
	request.Header.Set(httpserver.RequestIDHeader, "req_log")
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", response.Code)
	}

	var entry map[string]any
	if err := json.Unmarshal(logOutput.Bytes(), &entry); err != nil {
		t.Fatalf("log entry is not JSON: %v", err)
	}
	if entry["request_id"] != "req_log" {
		t.Fatalf("expected request id, got %v", entry["request_id"])
	}
	if entry["method"] != http.MethodPost {
		t.Fatalf("expected method, got %v", entry["method"])
	}
	if entry["path"] != "/orders" {
		t.Fatalf("expected path, got %v", entry["path"])
	}
	if entry["status"] != float64(http.StatusCreated) {
		t.Fatalf("expected status, got %v", entry["status"])
	}
}

func TestRequestLoggerDefaultsStatusToOKWhenHandlerDoesNotWrite(t *testing.T) {
	var logOutput bytes.Buffer
	log := logger.New(&logOutput, config.LogLevelDebug)
	handler := httpserver.WithRequestID(RequestLogger(log)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})))

	request := httptest.NewRequest(http.MethodGet, "/empty", nil)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	var entry map[string]any
	if err := json.Unmarshal(logOutput.Bytes(), &entry); err != nil {
		t.Fatalf("log entry is not JSON: %v", err)
	}
	if entry["status"] != float64(http.StatusOK) {
		t.Fatalf("expected default status 200, got %v", entry["status"])
	}
}
