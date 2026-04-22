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

func TestRecoverWritesErrorEnvelopeAndLogsRequest(t *testing.T) {
	var logOutput bytes.Buffer
	log := logger.New(&logOutput, config.LogLevelDebug)
	handler := httpserver.WithRequestID(Recover(log)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("boom")
	})))

	request := httptest.NewRequest(http.MethodGet, "/panic", nil)
	request.Header.Set(httpserver.RequestIDHeader, "req_panic")
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", response.Code)
	}

	var body httpserver.ErrorEnvelope
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("response is not JSON: %v", err)
	}
	if body.Error.Code != "internal.unexpected_error" {
		t.Fatalf("expected internal error code, got %q", body.Error.Code)
	}
	if body.RequestID != "req_panic" {
		t.Fatalf("expected request id, got %q", body.RequestID)
	}

	var entry map[string]any
	if err := json.Unmarshal(logOutput.Bytes(), &entry); err != nil {
		t.Fatalf("log entry is not JSON: %v", err)
	}
	if entry["request_id"] != "req_panic" {
		t.Fatalf("expected request id in log, got %v", entry["request_id"])
	}
	if entry["level"] != "error" {
		t.Fatalf("expected error log, got %v", entry["level"])
	}
}
