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

func TestRunCloudminiErrorEvidenceWithOutOfCapacity(t *testing.T) {
	var reservationCalls int
	exhaustedGroupID := "00000000-0000-4000-8000-000000000777"
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
		case r.Method == http.MethodGet && r.URL.Path == "/api/v3/inventory/groups":
			if r.Header.Get("Authorization") != "Bearer secret-token" {
				t.Fatalf("expected valid auth header")
			}
			if r.URL.Query().Get("kind") != "residential" {
				t.Fatalf("expected residential inventory kind, got %q", r.URL.Query().Get("kind"))
			}
			writeCloudminiSuccessEnvelope(t, w, http.StatusOK, []map[string]interface{}{{
				"id":                exhaustedGroupID,
				"kind":              "residential",
				"sell_state":        "exhausted",
				"allocatable_units": 0,
			}})
		case r.Method == http.MethodPost && r.URL.Path == "/api/v3/capacity/reservations":
			reservationCalls++
			if r.Header.Get("Authorization") != "Bearer secret-token" {
				t.Fatalf("expected valid auth header")
			}
			var payload map[string]interface{}
			if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
				t.Fatalf("decode reservation payload: %v", err)
			}
			if payload["group_id"] != exhaustedGroupID || payload["kind"] != "residential" || payload["quantity"] != float64(1) || payload["ttl_seconds"] != float64(60) {
				t.Fatalf("unexpected reservation payload: %#v", payload)
			}
			writeCloudminiErrorEvidenceEnvelope(t, w, http.StatusConflict, "CAPACITY_EXHAUSTED", "group has no allocatable units", map[string]string{"group_id": exhaustedGroupID})
		default:
			t.Fatalf("unexpected request %s %s auth=%q", r.Method, r.URL.Path, r.Header.Get("Authorization"))
		}
	}))
	defer server.Close()

	setCloudminiErrorEvidenceEnv(t, server.URL)
	t.Setenv("CLOUDMINI_ERROR_EVIDENCE_ALLOW_INVALID_CREATE", "")
	t.Setenv("CLOUDMINI_ERROR_EVIDENCE_MUTATING_ROUTE_APPROVED", "")
	t.Setenv("CLOUDMINI_ERROR_EVIDENCE_MAX_CREATE_ATTEMPTS", "")
	t.Setenv("CLOUDMINI_ERROR_EVIDENCE_ALLOW_OUT_OF_CAPACITY", "yes")
	t.Setenv("CLOUDMINI_ERROR_EVIDENCE_OUT_OF_CAPACITY_APPROVED", "yes")
	t.Setenv("CLOUDMINI_ERROR_EVIDENCE_OUT_OF_CAPACITY_MAX_RESERVATIONS", "1")
	t.Setenv("CLOUDMINI_ERROR_EVIDENCE_OUT_OF_CAPACITY_KIND", "residential")
	t.Setenv("CLOUDMINI_ERROR_EVIDENCE_OUT_OF_CAPACITY_TTL_SECONDS", "60")

	var out bytes.Buffer
	if err := runCloudminiErrorEvidenceSmokeWithWriter(2*time.Second, &out); err != nil {
		t.Fatalf("expected evidence pass: %v", err)
	}
	if reservationCalls != 1 {
		t.Fatalf("expected one reservation probe, got %d", reservationCalls)
	}
	output := out.String()
	for _, expected := range []string{
		"cloudmini_error_evidence result=PASS",
		"example_count=4",
		"mutating_routes_called=true",
		"example_4_name=out_of_capacity_reservation",
		"example_4_http_status=409",
		"example_4_provider_error_code=CAPACITY_EXHAUSTED",
		"example_4_normalized_error_code=PROVIDER_OUT_OF_STOCK",
		"example_4_retry_safety=do_not_retry",
		"example_4_side_effect_created=no",
		"example_4_reservation_probe_attempted=true",
		"example_4_exhausted_group_selected=true",
		"example_4_reservation_created=false",
		"example_4_reservation_cleaned_up=false",
		"example_4_reservation_max_attempts=1",
		"example_4_reservation_ttl_seconds=60",
		"remaining_provider_controlled_examples=permission_denied,rate_limited,provider_5xx,cancel_rejected",
	} {
		if !strings.Contains(output, expected) {
			t.Fatalf("expected output to contain %q, got:\n%s", expected, output)
		}
	}
	for _, leaked := range []string{"secret-token", "should-not-leak", "proxy-raw-secret", exhaustedGroupID, "group has no allocatable units"} {
		if strings.Contains(output, leaked) {
			t.Fatalf("redacted output leaked %q: %s", leaked, output)
		}
	}
}
