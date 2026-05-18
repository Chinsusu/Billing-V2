package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestRunCloudminiErrorEvidence(t *testing.T) {
	var malformedCreateCalls int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api/v3/capabilities" && r.Header.Get("Authorization") == "":
			writeCloudminiErrorEvidenceEnvelope(t, w, http.StatusUnauthorized, "AUTH_REQUIRED", "missing token", map[string]string{"secret": "should-not-leak"})
		case r.Method == http.MethodGet && r.URL.Path == "/api/v3/capabilities" && r.Header.Get("Authorization") == "Bearer billing-invalid-token":
			writeCloudminiErrorEvidenceEnvelope(t, w, http.StatusUnauthorized, "AUTH_INVALID", "bad token", map[string]string{"token": "should-not-leak"})
		case r.Method == http.MethodGet && r.URL.Path == "/api/v3/proxies/00000000-0000-4000-8000-000000000000":
			if r.Header.Get("Authorization") != "Bearer secret-token" {
				t.Fatalf("expected valid auth header")
			}
			writeCloudminiErrorEvidenceEnvelope(t, w, http.StatusNotFound, "PROXY_NOT_FOUND", "missing proxy id proxy-raw-secret", map[string]string{"proxy_id": "proxy-raw-secret"})
		case r.Method == http.MethodPost && r.URL.Path == "/api/v3/proxies":
			malformedCreateCalls++
			if r.Header.Get("Authorization") != "Bearer secret-token" {
				t.Fatalf("expected valid auth header")
			}
			if r.Header.Get("Idempotency-Key") == "" {
				t.Fatalf("expected idempotency key")
			}
			writeCloudminiErrorEvidenceEnvelope(t, w, http.StatusBadRequest, "INVALID_INPUT", "bad payload", map[string]string{"payload": "should-not-leak"})
		default:
			t.Fatalf("unexpected request %s %s auth=%q", r.Method, r.URL.Path, r.Header.Get("Authorization"))
		}
	}))
	defer server.Close()

	setCloudminiErrorEvidenceEnv(t, server.URL)
	var out bytes.Buffer
	if err := runCloudminiErrorEvidenceSmokeWithWriter(2*time.Second, &out); err != nil {
		t.Fatalf("expected evidence pass: %v", err)
	}
	if malformedCreateCalls != 1 {
		t.Fatalf("expected one malformed create call, got %d", malformedCreateCalls)
	}
	output := out.String()
	for _, expected := range []string{
		"cloudmini_error_evidence result=PASS",
		"example_count=4",
		"mutating_routes_called=true",
		"example_1_name=auth_missing_capabilities",
		"example_1_normalized_error_code=PROVIDER_AUTH_FAILED",
		"example_3_name=not_found_proxy",
		"example_3_provider_error_code=PROXY_NOT_FOUND",
		"example_3_normalized_error_code=PROVIDER_STATE_DRIFT",
		"example_4_name=validation_malformed_create",
		"example_4_provider_error_code=INVALID_INPUT",
		"example_4_normalized_error_code=PROVIDER_CONFIG_INVALID",
		"raw_response_body_printed=no",
		"sensitive_values_printed=no",
	} {
		if !strings.Contains(output, expected) {
			t.Fatalf("expected output to contain %q, got:\n%s", expected, output)
		}
	}
	for _, leaked := range []string{"secret-token", "should-not-leak", "proxy-raw-secret", "bad payload", "missing proxy id"} {
		if strings.Contains(output, leaked) {
			t.Fatalf("redacted output leaked %q: %s", leaked, output)
		}
	}
}

func TestRunCloudminiErrorEvidenceWithPermissionDenied(t *testing.T) {
	var createCalls int
	var revokeCalls int
	temporaryKeyActive := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api/v3/capabilities" && r.Header.Get("Authorization") == "":
			writeCloudminiErrorEvidenceEnvelope(t, w, http.StatusUnauthorized, "AUTH_REQUIRED", "missing token", map[string]string{"secret": "should-not-leak"})
		case r.Method == http.MethodGet && r.URL.Path == "/api/v3/capabilities" && r.Header.Get("Authorization") == "Bearer billing-invalid-token":
			writeCloudminiErrorEvidenceEnvelope(t, w, http.StatusUnauthorized, "AUTH_INVALID", "bad token", map[string]string{"token": "should-not-leak"})
		case r.Method == http.MethodGet && r.URL.Path == "/api/v3/proxies/00000000-0000-4000-8000-000000000000":
			if r.Header.Get("Authorization") != "Bearer secret-token" {
				t.Fatalf("expected valid auth header")
			}
			writeCloudminiErrorEvidenceEnvelope(t, w, http.StatusNotFound, "PROXY_NOT_FOUND", "missing proxy id proxy-raw-secret", map[string]string{"proxy_id": "proxy-raw-secret"})
		case r.Method == http.MethodGet && r.URL.Path == "/api/v1/api-keys/":
			if r.Header.Get("Authorization") != "Bearer secret-token" {
				t.Fatalf("expected management auth header")
			}
			data := []map[string]interface{}{{"id": "existing-key", "is_active": true}}
			if temporaryKeyActive {
				data = append(data, map[string]interface{}{"id": "temp-key-id", "is_active": true})
			}
			writeCloudminiSuccessEnvelope(t, w, http.StatusOK, data)
		case r.Method == http.MethodPost && r.URL.Path == "/api/v1/api-keys/":
			createCalls++
			if r.Header.Get("Authorization") != "Bearer secret-token" {
				t.Fatalf("expected management auth header")
			}
			temporaryKeyActive = true
			writeCloudminiSuccessEnvelope(t, w, http.StatusCreated, map[string]interface{}{
				"api_key": map[string]interface{}{
					"id":        "temp-key-id",
					"is_active": true,
				},
				"plain_key": "read-token",
			})
		case r.Method == http.MethodGet && r.URL.Path == "/api/v3/proxies":
			if r.Header.Get("Authorization") != "Bearer read-token" {
				t.Fatalf("expected temporary low-scope auth header")
			}
			writeCloudminiErrorString(t, w, http.StatusForbidden, "missing specialized api permissions")
		case r.Method == http.MethodDelete && r.URL.Path == "/api/v1/api-keys/temp-key-id":
			revokeCalls++
			if r.Header.Get("Authorization") != "Bearer secret-token" {
				t.Fatalf("expected management auth header")
			}
			temporaryKeyActive = false
			w.WriteHeader(http.StatusNoContent)
		default:
			t.Fatalf("unexpected request %s %s auth=%q", r.Method, r.URL.Path, r.Header.Get("Authorization"))
		}
	}))
	defer server.Close()

	setCloudminiErrorEvidenceEnv(t, server.URL)
	t.Setenv("CLOUDMINI_ERROR_EVIDENCE_ALLOW_INVALID_CREATE", "")
	t.Setenv("CLOUDMINI_ERROR_EVIDENCE_MUTATING_ROUTE_APPROVED", "")
	t.Setenv("CLOUDMINI_ERROR_EVIDENCE_MAX_CREATE_ATTEMPTS", "")
	t.Setenv("CLOUDMINI_ERROR_EVIDENCE_ALLOW_PERMISSION_DENIED", "yes")
	t.Setenv("CLOUDMINI_ERROR_EVIDENCE_PERMISSION_KEY_MANAGEMENT_APPROVED", "yes")
	t.Setenv("CLOUDMINI_ERROR_EVIDENCE_PERMISSION_KEY_MAX_CREATE", "1")

	var out bytes.Buffer
	if err := runCloudminiErrorEvidenceSmokeWithWriter(2*time.Second, &out); err != nil {
		t.Fatalf("expected evidence pass: %v", err)
	}
	if createCalls != 1 || revokeCalls != 1 {
		t.Fatalf("expected one create and one revoke, got create=%d revoke=%d", createCalls, revokeCalls)
	}
	output := out.String()
	for _, expected := range []string{
		"cloudmini_error_evidence result=PASS",
		"example_count=4",
		"mutating_routes_called=true",
		"example_4_name=permission_denied_proxy_list",
		"example_4_http_status=403",
		"example_4_provider_error_code=none",
		"example_4_normalized_error_code=PROVIDER_PERMISSION_DENIED",
		"example_4_retry_safety=do_not_retry",
		"example_4_temporary_api_key_created=true",
		"example_4_temporary_api_key_revoked=true",
		"example_4_active_key_count_restored=true",
		"remaining_provider_controlled_examples=rate_limited,out_of_capacity,provider_5xx,cancel_rejected",
	} {
		if !strings.Contains(output, expected) {
			t.Fatalf("expected output to contain %q, got:\n%s", expected, output)
		}
	}
	for _, leaked := range []string{"secret-token", "read-token", "temp-key-id", "should-not-leak", "missing specialized api permissions"} {
		if strings.Contains(output, leaked) {
			t.Fatalf("redacted output leaked %q: %s", leaked, output)
		}
	}
}

func TestCloudminiErrorEvidenceRequiresApproval(t *testing.T) {
	t.Setenv("APP_ENV", "local")
	var out bytes.Buffer
	err := runCloudminiErrorEvidenceSmokeWithWriter(time.Second, &out)
	if err == nil || !strings.Contains(err.Error(), "BILLING_CLOUDMINI_ERROR_EVIDENCE_APPROVED") {
		t.Fatalf("expected approval error, got %v", err)
	}
}

func TestCloudminiErrorEvidenceRequiresMutatingApproval(t *testing.T) {
	t.Setenv("APP_ENV", "local")
	t.Setenv("BILLING_CLOUDMINI_ERROR_EVIDENCE_APPROVED", "yes")
	t.Setenv("CLOUDMINI_SOURCE_ACCOUNT_OWNER", "Admin")
	t.Setenv("CLOUDMINI_ENGINEERING_OWNER", "Admin")
	t.Setenv("CLOUDMINI_OPS_OWNER", "Admin")
	t.Setenv("CLOUDMINI_SECURITY_OWNER", "Admin")
	t.Setenv("CLOUDMINI_CLEANUP_OWNER", "Admin")
	t.Setenv("CLOUDMINI_REVIEWER_SIGNOFF", "Admin")
	t.Setenv("CLOUDMINI_PILOT_STOP_CONDITION", "stop-on-unexpected-success")
	t.Setenv("CLOUDMINI_PILOT_READONLY_EVIDENCE_REF", "T249")
	t.Setenv("CLOUDMINI_V3_BASE_URL", "https://example.invalid")
	t.Setenv("CLOUDMINI_V3_API_TOKEN", "secret-token")
	t.Setenv("CLOUDMINI_ERROR_EVIDENCE_ALLOW_INVALID_CREATE", "yes")

	var out bytes.Buffer
	err := runCloudminiErrorEvidenceSmokeWithWriter(time.Second, &out)
	if err == nil || !strings.Contains(err.Error(), "CLOUDMINI_ERROR_EVIDENCE_MUTATING_ROUTE_APPROVED") {
		t.Fatalf("expected mutating approval error, got %v", err)
	}
}

func setCloudminiErrorEvidenceEnv(t *testing.T, baseURL string) {
	t.Helper()
	t.Setenv("APP_ENV", "local")
	t.Setenv("BILLING_CLOUDMINI_ERROR_EVIDENCE_APPROVED", "yes")
	t.Setenv("CLOUDMINI_SOURCE_ACCOUNT_OWNER", "Admin")
	t.Setenv("CLOUDMINI_ENGINEERING_OWNER", "Admin")
	t.Setenv("CLOUDMINI_OPS_OWNER", "Admin")
	t.Setenv("CLOUDMINI_SECURITY_OWNER", "Admin")
	t.Setenv("CLOUDMINI_CLEANUP_OWNER", "Admin")
	t.Setenv("CLOUDMINI_REVIEWER_SIGNOFF", "Admin")
	t.Setenv("CLOUDMINI_PILOT_STOP_CONDITION", "stop-on-unexpected-success")
	t.Setenv("CLOUDMINI_PILOT_READONLY_EVIDENCE_REF", "T249")
	t.Setenv("CLOUDMINI_V3_BASE_URL", baseURL)
	t.Setenv("CLOUDMINI_V3_API_TOKEN", "secret-token")
	t.Setenv("CLOUDMINI_ERROR_EVIDENCE_ALLOW_INVALID_CREATE", "yes")
	t.Setenv("CLOUDMINI_ERROR_EVIDENCE_MUTATING_ROUTE_APPROVED", "yes")
	t.Setenv("CLOUDMINI_ERROR_EVIDENCE_MAX_CREATE_ATTEMPTS", "1")
}

func writeCloudminiErrorEvidenceEnvelope(t *testing.T, w http.ResponseWriter, status int, code string, message string, details interface{}) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"success": false,
		"error": map[string]interface{}{
			"code":    code,
			"message": message,
			"details": details,
		},
	}); err != nil {
		t.Fatalf("write response: %v", err)
	}
}

func writeCloudminiSuccessEnvelope(t *testing.T, w http.ResponseWriter, status int, data interface{}) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    data,
	}); err != nil {
		t.Fatalf("write response: %v", err)
	}
}

func writeCloudminiErrorString(t *testing.T, w http.ResponseWriter, status int, message string) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"error": message,
	}); err != nil {
		t.Fatalf("write response: %v", err)
	}
}
