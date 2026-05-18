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

func TestRunCloudminiErrorEvidenceWithCancelDeleteFixture(t *testing.T) {
	var fixtureCalls int
	invalidAuth := "Bearer " + "billing-invalid-token"
	validAuth := "Bearer " + "secret-token"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api/v3/capabilities" && r.Header.Get("Authorization") == "":
			writeCloudminiErrorEvidenceEnvelope(t, w, http.StatusUnauthorized, "AUTH_REQUIRED", "missing token", map[string]string{"secret": "should-not-leak"})
		case r.Method == http.MethodGet && r.URL.Path == "/api/v3/capabilities" && r.Header.Get("Authorization") == invalidAuth:
			writeCloudminiErrorEvidenceEnvelope(t, w, http.StatusUnauthorized, "AUTH_INVALID", "bad token", map[string]string{"token": "should-not-leak"})
		case r.Method == http.MethodGet && r.URL.Path == "/api/v3/proxies/00000000-0000-4000-8000-000000000000":
			if r.Header.Get("Authorization") != validAuth {
				t.Fatalf("expected valid auth header")
			}
			writeCloudminiErrorEvidenceEnvelope(t, w, http.StatusNotFound, "PROXY_NOT_FOUND", "missing proxy id proxy-raw-secret", map[string]string{"proxy_id": "proxy-raw-secret"})
		case r.Method == http.MethodGet && r.URL.Path == "/api/v3/error-fixtures/delete-rejected":
			fixtureCalls++
			if r.Header.Get("Authorization") != validAuth {
				t.Fatalf("expected valid auth header")
			}
			if r.Header.Get("X-Cloudmini-Error-Fixture") != "delete_rejected" {
				t.Fatalf("expected cancel/delete fixture header")
			}
			writeCloudminiOperationEvidenceEnvelope(t, w, http.StatusOK, map[string]string{
				"state":         "failed",
				"error_code":    "DELETE_FAILED",
				"error_message": "delete rejected secret detail",
			})
		default:
			t.Fatalf("unexpected request %s %s auth=%q", r.Method, r.URL.Path, r.Header.Get("Authorization"))
		}
	}))
	defer server.Close()

	setCloudminiErrorEvidenceEnv(t, server.URL)
	t.Setenv("CLOUDMINI_ERROR_EVIDENCE_ALLOW_INVALID_CREATE", "")
	t.Setenv("CLOUDMINI_ERROR_EVIDENCE_MUTATING_ROUTE_APPROVED", "")
	t.Setenv("CLOUDMINI_ERROR_EVIDENCE_MAX_CREATE_ATTEMPTS", "")
	t.Setenv("CLOUDMINI_ERROR_EVIDENCE_ALLOW_CANCEL_DELETE_REJECTED", "yes")
	t.Setenv("CLOUDMINI_ERROR_EVIDENCE_CANCEL_DELETE_APPROVED", "yes")
	t.Setenv("CLOUDMINI_ERROR_EVIDENCE_CANCEL_DELETE_MAX_REQUESTS", "1")
	t.Setenv("CLOUDMINI_ERROR_EVIDENCE_CANCEL_DELETE_FIXTURE_PATH", "/api/v3/error-fixtures/delete-rejected")

	var out bytes.Buffer
	if err := runCloudminiErrorEvidenceSmokeWithWriter(2*time.Second, &out); err != nil {
		t.Fatalf("expected evidence pass: %v", err)
	}
	if fixtureCalls != 1 {
		t.Fatalf("expected one fixture call, got %d", fixtureCalls)
	}
	output := out.String()
	for _, expected := range []string{
		"cloudmini_error_evidence result=PASS",
		"example_count=4",
		"mutating_routes_called=false",
		"example_4_name=cancel_delete_rejected_fixture",
		"example_4_http_status=200",
		"example_4_provider_error_code=DELETE_FAILED",
		"example_4_normalized_error_code=PROVIDER_PARTIAL_SUCCESS",
		"example_4_retry_safety=manual_review_required",
		"example_4_side_effect_created=not_applicable",
		"example_4_cancel_delete_fixture_called=true",
		"example_4_cancel_delete_max_requests=1",
		"example_4_provider_operation_state=failed",
		"remaining_provider_controlled_examples=permission_denied,rate_limited,out_of_capacity,provider_5xx",
	} {
		if !strings.Contains(output, expected) {
			t.Fatalf("expected output to contain %q, got:\n%s", expected, output)
		}
	}
	for _, leaked := range []string{"secret-token", "should-not-leak", "proxy-raw-secret", "delete rejected secret detail"} {
		if strings.Contains(output, leaked) {
			t.Fatalf("redacted output leaked %q: %s", leaked, output)
		}
	}
}

func TestCloudminiErrorEvidenceRequiresCancelDeleteFixtureApproval(t *testing.T) {
	t.Setenv("APP_ENV", "local")
	t.Setenv("BILLING_CLOUDMINI_ERROR_EVIDENCE_APPROVED", "yes")
	t.Setenv("CLOUDMINI_SOURCE_ACCOUNT_OWNER", "Admin")
	t.Setenv("CLOUDMINI_ENGINEERING_OWNER", "Admin")
	t.Setenv("CLOUDMINI_OPS_OWNER", "Admin")
	t.Setenv("CLOUDMINI_SECURITY_OWNER", "Admin")
	t.Setenv("CLOUDMINI_CLEANUP_OWNER", "Admin")
	t.Setenv("CLOUDMINI_REVIEWER_SIGNOFF", "Admin")
	t.Setenv("CLOUDMINI_PILOT_STOP_CONDITION", "stop-on-unexpected-success")
	t.Setenv("CLOUDMINI_PILOT_READONLY_EVIDENCE_REF", "T262")
	t.Setenv("CLOUDMINI_V3_BASE_URL", "https://example.invalid")
	t.Setenv("CLOUDMINI_V3_API_TOKEN", "ok")
	t.Setenv("CLOUDMINI_ERROR_EVIDENCE_ALLOW_CANCEL_DELETE_REJECTED", "yes")

	var out bytes.Buffer
	err := runCloudminiErrorEvidenceSmokeWithWriter(time.Second, &out)
	if err == nil || !strings.Contains(err.Error(), "CLOUDMINI_ERROR_EVIDENCE_CANCEL_DELETE_APPROVED") {
		t.Fatalf("expected cancel/delete approval error, got %v", err)
	}
}

func writeCloudminiOperationEvidenceEnvelope(t *testing.T, w http.ResponseWriter, status int, data interface{}) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    data,
	}); err != nil {
		t.Fatalf("write operation envelope: %v", err)
	}
}
