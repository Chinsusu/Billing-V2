package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Chinsusu/Billing-V2/internal/platform/httpserver"
)

func TestRequireMethodAllowsExpectedMethod(t *testing.T) {
	handler := httpserver.WithRequestID(http.HandlerFunc(RequireMethod(http.MethodGet, func(w http.ResponseWriter, r *http.Request) {
		httpserver.WriteSuccess(w, r, http.StatusOK, map[string]string{"status": "ok"})
	})))

	request := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.Code)
	}
}

func TestRequireMethodRejectsUnexpectedMethod(t *testing.T) {
	handler := httpserver.WithRequestID(http.HandlerFunc(RequireMethod(http.MethodGet, func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("next handler should not be called")
	})))

	request := httptest.NewRequest(http.MethodPost, "/healthz", nil)
	request.Header.Set(httpserver.RequestIDHeader, "req_method")
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	if response.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405, got %d", response.Code)
	}
	if response.Header().Get("Allow") != http.MethodGet {
		t.Fatalf("expected Allow header to be %q, got %q", http.MethodGet, response.Header().Get("Allow"))
	}

	var body httpserver.ErrorEnvelope
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("response is not JSON: %v", err)
	}
	if body.Error.Code != "request.method_not_allowed" {
		t.Fatalf("expected method error code, got %q", body.Error.Code)
	}
	if body.RequestID != "req_method" {
		t.Fatalf("expected request id, got %q", body.RequestID)
	}
}
