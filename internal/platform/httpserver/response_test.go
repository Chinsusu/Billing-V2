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

func TestWriteListIncludesPagination(t *testing.T) {
	handler := WithRequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		WriteList(w, r, http.StatusOK, []string{"order_1"}, NewPage(20, "cursor_2"))
	}))

	request := httptest.NewRequest(http.MethodGet, "/orders", nil)
	request.Header.Set(RequestIDHeader, "req_list")
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	var body SuccessEnvelope
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("response is not JSON: %v", err)
	}
	if body.Page == nil || body.Page.Limit != 20 {
		t.Fatalf("expected page info, got %#v", body.Page)
	}
	if body.Page.NextCursor == nil || *body.Page.NextCursor != "cursor_2" {
		t.Fatalf("expected next cursor, got %#v", body.Page.NextCursor)
	}
}

func TestWriteErrorWithDetails(t *testing.T) {
	handler := WithRequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		WriteErrorWithDetails(w, r, http.StatusConflict, "idempotency.conflict", "Idempotency key already exists.", map[string]string{"idempotency_key": "idem_1"})
	}))

	request := httptest.NewRequest(http.MethodPost, "/orders", nil)
	request.Header.Set(RequestIDHeader, "req_error")
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	var body ErrorEnvelope
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("response is not JSON: %v", err)
	}
	if body.Error.Code != "idempotency.conflict" {
		t.Fatalf("expected code, got %q", body.Error.Code)
	}
	if body.Error.Details == nil {
		t.Fatal("expected safe details")
	}
}

func TestWriteValidationError(t *testing.T) {
	handler := WithRequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		WriteValidationError(w, r, []ValidationField{{Field: "email", Code: "email.invalid", Message: "Email is invalid."}})
	}))

	request := httptest.NewRequest(http.MethodPost, "/auth/login", nil)
	request.Header.Set(RequestIDHeader, "req_validation")
	response := httptest.NewRecorder()

	handler.ServeHTTP(response, request)

	var body ErrorEnvelope
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("response is not JSON: %v", err)
	}
	if response.Code != http.StatusBadRequest {
		t.Fatalf("expected bad request, got %d", response.Code)
	}
	if body.Error.Code != "validation.failed" {
		t.Fatalf("expected validation code, got %q", body.Error.Code)
	}
	if len(body.Error.Fields) != 1 || body.Error.Fields[0].Field != "email" {
		t.Fatalf("expected field error, got %#v", body.Error.Fields)
	}
}

func TestParseCursorPageDefaults(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/orders", nil)

	page, err := ParseCursorPage(request)
	if err != nil {
		t.Fatalf("expected default page, got %v", err)
	}
	if page.Limit != DefaultPageLimit || page.Cursor != "" {
		t.Fatalf("unexpected page request: %#v", page)
	}
}

func TestParseCursorPageRejectsTooLargeLimit(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "/orders?limit=101", nil)

	if _, err := ParseCursorPage(request); err != ErrPageLimitTooLarge {
		t.Fatalf("expected limit error, got %v", err)
	}
}
