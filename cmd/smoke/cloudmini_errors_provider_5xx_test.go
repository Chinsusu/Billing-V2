package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestRunCloudminiErrorEvidenceWithProvider5xxFixture(t *testing.T) {
	var fixtureCalls int
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
		case r.Method == http.MethodGet && r.URL.Path == "/api/v3/error-fixtures/internal-error":
			fixtureCalls++
			if r.Header.Get("Authorization") != "Bearer secret-token" {
				t.Fatalf("expected valid auth header")
			}
			if r.Header.Get("X-Cloudmini-Error-Fixture") != "internal_error" {
				t.Fatalf("expected provider 5xx fixture header")
			}
			writeCloudminiErrorEvidenceEnvelope(t, w, http.StatusInternalServerError, "INTERNAL_ERROR", "fixture internal error secret detail", map[string]string{"provider_trace": "should-not-leak"})
		default:
			t.Fatalf("unexpected request %s %s auth=%q", r.Method, r.URL.Path, r.Header.Get("Authorization"))
		}
	}))
	defer server.Close()

	setCloudminiErrorEvidenceEnv(t, server.URL)
	t.Setenv("CLOUDMINI_ERROR_EVIDENCE_ALLOW_INVALID_CREATE", "")
	t.Setenv("CLOUDMINI_ERROR_EVIDENCE_MUTATING_ROUTE_APPROVED", "")
	t.Setenv("CLOUDMINI_ERROR_EVIDENCE_MAX_CREATE_ATTEMPTS", "")
	t.Setenv("CLOUDMINI_ERROR_EVIDENCE_ALLOW_PROVIDER_5XX", "yes")
	t.Setenv("CLOUDMINI_ERROR_EVIDENCE_PROVIDER_5XX_APPROVED", "yes")
	t.Setenv("CLOUDMINI_ERROR_EVIDENCE_PROVIDER_5XX_MAX_REQUESTS", "1")
	t.Setenv("CLOUDMINI_ERROR_EVIDENCE_PROVIDER_5XX_FIXTURE_PATH", "/api/v3/error-fixtures/internal-error")

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
		"example_4_name=provider_5xx_fixture",
		"example_4_http_status=500",
		"example_4_provider_error_code=INTERNAL_ERROR",
		"example_4_normalized_error_code=PROVIDER_TEMPORARY_ERROR",
		"example_4_retry_safety=safe_retry",
		"example_4_side_effect_created=not_applicable",
		"example_4_provider_5xx_fixture_called=true",
		"example_4_provider_5xx_max_requests=1",
		"remaining_provider_controlled_examples=permission_denied,rate_limited,out_of_capacity,cancel_rejected",
	} {
		if !strings.Contains(output, expected) {
			t.Fatalf("expected output to contain %q, got:\n%s", expected, output)
		}
	}
	for _, leaked := range []string{"secret-token", "should-not-leak", "proxy-raw-secret", "fixture internal error secret detail"} {
		if strings.Contains(output, leaked) {
			t.Fatalf("redacted output leaked %q: %s", leaked, output)
		}
	}
}

func TestCloudminiErrorEvidenceRequiresProvider5xxFixtureApproval(t *testing.T) {
	t.Setenv("APP_ENV", "local")
	t.Setenv("BILLING_CLOUDMINI_ERROR_EVIDENCE_APPROVED", "yes")
	t.Setenv("CLOUDMINI_SOURCE_ACCOUNT_OWNER", "Admin")
	t.Setenv("CLOUDMINI_ENGINEERING_OWNER", "Admin")
	t.Setenv("CLOUDMINI_OPS_OWNER", "Admin")
	t.Setenv("CLOUDMINI_SECURITY_OWNER", "Admin")
	t.Setenv("CLOUDMINI_CLEANUP_OWNER", "Admin")
	t.Setenv("CLOUDMINI_REVIEWER_SIGNOFF", "Admin")
	t.Setenv("CLOUDMINI_PILOT_STOP_CONDITION", "stop-on-unexpected-success")
	t.Setenv("CLOUDMINI_PILOT_READONLY_EVIDENCE_REF", "T261")
	t.Setenv("CLOUDMINI_V3_BASE_URL", "https://example.invalid")
	t.Setenv("CLOUDMINI_V3_API_TOKEN", "secret-token")
	t.Setenv("CLOUDMINI_ERROR_EVIDENCE_ALLOW_PROVIDER_5XX", "yes")

	var out bytes.Buffer
	err := runCloudminiErrorEvidenceSmokeWithWriter(time.Second, &out)
	if err == nil || !strings.Contains(err.Error(), "CLOUDMINI_ERROR_EVIDENCE_PROVIDER_5XX_APPROVED") {
		t.Fatalf("expected provider 5xx approval error, got %v", err)
	}
}
