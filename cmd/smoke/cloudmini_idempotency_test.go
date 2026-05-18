package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestRunCloudminiIdempotencyEvidenceDuplicateCreate(t *testing.T) {
	var createCalls, cleanupCalls int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer secret-token" {
			t.Fatalf("unexpected auth header")
		}
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/api/v3/proxies":
			createCalls++
			if r.Header.Get("Idempotency-Key") != "billing-t248-test-duplicate-create" {
				t.Fatalf("unexpected idempotency key")
			}
			writeCloudminiSmokeSuccess(t, w, http.StatusAccepted, map[string]interface{}{
				"resource":  map[string]interface{}{"id": "proxy-raw-1", "kind": "ipv4_dc", "status": "creating"},
				"operation": map[string]interface{}{"id": "op-create-1", "state": "accepted"},
			})
		case r.Method == http.MethodGet && r.URL.Path == "/api/v3/operations/op-create-1":
			writeCloudminiSmokeSuccess(t, w, http.StatusOK, map[string]interface{}{
				"id":          "op-create-1",
				"resource_id": "proxy-raw-1",
				"state":       "succeeded",
				"resource_snapshot": map[string]interface{}{
					"id":          "proxy-raw-1",
					"kind":        "ipv4_dc",
					"status":      "running",
					"host":        "203.0.113.50",
					"outbound_ip": "203.0.113.50",
					"port_socks":  1080,
					"username":    "proxy-user",
					"password":    "proxy-pass",
				},
			})
		case r.Method == http.MethodDelete && r.URL.Path == "/api/v3/proxies/proxy-raw-1":
			cleanupCalls++
			writeCloudminiSmokeSuccess(t, w, http.StatusAccepted, map[string]interface{}{
				"resource":  map[string]interface{}{"id": "proxy-raw-1", "kind": "ipv4_dc", "status": "deleting"},
				"operation": map[string]interface{}{"id": "op-delete-1", "state": "accepted"},
			})
		case r.Method == http.MethodGet && r.URL.Path == "/api/v3/operations/op-delete-1":
			writeCloudminiSmokeSuccess(t, w, http.StatusOK, map[string]interface{}{
				"id":          "op-delete-1",
				"resource_id": "proxy-raw-1",
				"state":       "succeeded",
			})
		default:
			t.Fatalf("unexpected request %s %s", r.Method, r.URL.Path)
		}
	}))
	defer server.Close()

	rawPath := filepath.Join(t.TempDir(), "raw.json")
	setCloudminiEvidenceEnv(t, server.URL, rawPath, cloudminiScenarioDuplicateCreate, "2")

	var out bytes.Buffer
	if err := runCloudminiIdempotencyEvidenceSmokeWithWriter(2*time.Second, &out); err != nil {
		t.Fatalf("expected evidence pass: %v", err)
	}
	if createCalls != 2 || cleanupCalls != 1 {
		t.Fatalf("unexpected call counts: create=%d cleanup=%d", createCalls, cleanupCalls)
	}
	output := out.String()
	for _, expected := range []string{
		"cloudmini_idempotency_evidence result=PASS",
		"scenario=duplicate-create",
		"create_attempts=2",
		"distinct_resource_count=1",
		"duplicate_same_resource=true",
		"cleanup_attempts=1",
		"sensitive_values_printed=no",
		"raw_provider_ids_printed=no",
	} {
		if !strings.Contains(output, expected) {
			t.Fatalf("expected output to contain %q, got:\n%s", expected, output)
		}
	}
	for _, leaked := range []string{"secret-token", "proxy-raw-1", "op-create-1", "proxy-pass"} {
		if strings.Contains(output, leaked) {
			t.Fatalf("redacted output leaked %q: %s", leaked, output)
		}
	}
	info, err := os.Stat(rawPath)
	if err != nil {
		t.Fatalf("expected raw evidence file: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Fatalf("raw evidence file must be 0600, got %o", info.Mode().Perm())
	}
}

func TestRunCloudminiIdempotencyEvidenceTimeoutAfterSend(t *testing.T) {
	var cleanupCalls int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/api/v3/proxies":
			writeCloudminiSmokeSuccess(t, w, http.StatusAccepted, map[string]interface{}{
				"resource":  map[string]interface{}{"id": "proxy-timeout-1", "kind": "ipv4_dc", "status": "creating"},
				"operation": map[string]interface{}{"id": "op-timeout-1", "state": "accepted"},
			})
		case r.Method == http.MethodGet && r.URL.Path == "/api/v3/operations/op-timeout-1":
			writeCloudminiSmokeSuccess(t, w, http.StatusOK, map[string]interface{}{
				"id":          "op-timeout-1",
				"resource_id": "proxy-timeout-1",
				"state":       "running",
			})
		case r.Method == http.MethodDelete && r.URL.Path == "/api/v3/proxies/proxy-timeout-1":
			cleanupCalls++
			writeCloudminiSmokeSuccess(t, w, http.StatusAccepted, map[string]interface{}{
				"resource":  map[string]interface{}{"id": "proxy-timeout-1", "kind": "ipv4_dc", "status": "deleting"},
				"operation": map[string]interface{}{"id": "op-timeout-delete-1", "state": "accepted"},
			})
		case r.Method == http.MethodGet && r.URL.Path == "/api/v3/operations/op-timeout-delete-1":
			writeCloudminiSmokeSuccess(t, w, http.StatusOK, map[string]interface{}{
				"id":          "op-timeout-delete-1",
				"resource_id": "proxy-timeout-1",
				"state":       "succeeded",
			})
		default:
			t.Fatalf("unexpected request %s %s", r.Method, r.URL.Path)
		}
	}))
	defer server.Close()

	rawPath := filepath.Join(t.TempDir(), "raw.json")
	setCloudminiEvidenceEnv(t, server.URL, rawPath, cloudminiScenarioTimeoutAfterSend, "1")
	t.Setenv("CLOUDMINI_V3_POLL_INTERVAL", "1ms")
	t.Setenv("CLOUDMINI_V3_POLL_TIMEOUT", "2ms")

	var out bytes.Buffer
	if err := runCloudminiIdempotencyEvidenceSmokeWithWriter(2*time.Second, &out); err != nil {
		t.Fatalf("expected timeout evidence pass with cleanup: %v", err)
	}
	if cleanupCalls != 1 {
		t.Fatalf("expected one cleanup call, got %d", cleanupCalls)
	}
	output := out.String()
	for _, expected := range []string{
		"cloudmini_idempotency_evidence result=PASS",
		"scenario=timeout-after-send",
		"create_attempts=1",
		"create_1_error_code=PROVIDER_TIMEOUT_REQUEST_KNOWN",
		"create_1_retry_safety=manual_review_required",
		"cleanup_attempts=1",
	} {
		if !strings.Contains(output, expected) {
			t.Fatalf("expected output to contain %q, got:\n%s", expected, output)
		}
	}
}

func TestCloudminiIdempotencyEvidenceRequiresApproval(t *testing.T) {
	t.Setenv("APP_ENV", "local")
	var out bytes.Buffer
	err := runCloudminiIdempotencyEvidenceSmokeWithWriter(time.Second, &out)
	if err == nil || !strings.Contains(err.Error(), "BILLING_CLOUDMINI_IDEMPOTENCY_EVIDENCE_APPROVED") {
		t.Fatalf("expected approval error, got %v", err)
	}
}

func setCloudminiEvidenceEnv(t *testing.T, baseURL string, rawPath string, scenario string, maxCreates string) {
	t.Helper()
	t.Setenv("APP_ENV", "local")
	t.Setenv("BILLING_CLOUDMINI_IDEMPOTENCY_EVIDENCE_APPROVED", "yes")
	t.Setenv("CLOUDMINI_SOURCE_ACCOUNT_OWNER", "Admin")
	t.Setenv("CLOUDMINI_ENGINEERING_OWNER", "Admin")
	t.Setenv("CLOUDMINI_OPS_OWNER", "Admin")
	t.Setenv("CLOUDMINI_SECURITY_OWNER", "Admin")
	t.Setenv("CLOUDMINI_CLEANUP_OWNER", "Admin")
	t.Setenv("CLOUDMINI_FINANCE_QUOTA_OWNER", "Admin")
	t.Setenv("CLOUDMINI_REVIEWER_SIGNOFF", "Admin")
	t.Setenv("CLOUDMINI_PILOT_CLEANUP_DEADLINE", "same-session")
	t.Setenv("CLOUDMINI_PILOT_STOP_CONDITION", "cleanup-failure-or-duplicate-resource")
	t.Setenv("CLOUDMINI_PILOT_READONLY_EVIDENCE_REF", "T216/T221")
	t.Setenv("CLOUDMINI_PILOT_CLEANUP_PROCEDURE_REF", "docs/03_execution_operations_launch/71_Cloudmini_Controlled_Pilot_Runbook.md")
	t.Setenv("CLOUDMINI_IDEMPOTENCY_SCENARIO", scenario)
	t.Setenv("CLOUDMINI_IDEMPOTENCY_PILOT_ID", "t248-test")
	t.Setenv("CLOUDMINI_IDEMPOTENCY_MAX_CREATE_ATTEMPTS", maxCreates)
	t.Setenv("CLOUDMINI_IDEMPOTENCY_MAX_ACTIVE_RESOURCES", "1")
	t.Setenv("CLOUDMINI_IDEMPOTENCY_PROVIDER_RATE_LIMIT", "no-parallel-mutating-calls")
	t.Setenv("CLOUDMINI_IDEMPOTENCY_MAX_SPEND_EXPOSURE", "single-dev-resource")
	t.Setenv("CLOUDMINI_IDEMPOTENCY_RAW_EVIDENCE_PATH", rawPath)
	t.Setenv("CLOUDMINI_V3_BASE_URL", baseURL)
	t.Setenv("CLOUDMINI_V3_API_TOKEN", "secret-token")
	t.Setenv("CLOUDMINI_V3_SOURCE_ID", "source-1")
	t.Setenv("CLOUDMINI_V3_KIND", "ipv4_dc")
	t.Setenv("CLOUDMINI_V3_GROUP_ID", "group-1")
	t.Setenv("CLOUDMINI_V3_PROTOCOL", "socks5")
	t.Setenv("ENCRYPTION_KEY", "12345678901234567890123456789012")
}

func writeCloudminiSmokeSuccess(t *testing.T, w http.ResponseWriter, status int, data interface{}) {
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
