package httpserver

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWriteSuccessIncludesRequestID(t *testing.T) {
	handler := WithRequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		WriteSuccess(w, r, http.StatusOK, map[string]string{"status": "ok"})
	}))

	request := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	request.Header.Set(RequestIDHeader, "req_test")
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	var body SuccessEnvelope
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("response is not JSON: %v", err)
	}
	if body.RequestID != "req_test" {
		t.Fatalf("expected request id, got %q", body.RequestID)
	}
	if response.Header().Get(RequestIDHeader) != "req_test" {
		t.Fatalf("expected response request id header")
	}
}

func TestWriteErrorIncludesStableCode(t *testing.T) {
	handler := WithRequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		WriteError(w, r, http.StatusServiceUnavailable, "service.unavailable", "Service is unavailable.")
	}))

	request := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	var body ErrorEnvelope
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("response is not JSON: %v", err)
	}
	if body.Error.Code != "service.unavailable" {
		t.Fatalf("expected stable error code, got %q", body.Error.Code)
	}
}
