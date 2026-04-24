package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNormalizedAPIURL(t *testing.T) {
	got, err := normalizedAPIURL("http://localhost:8080/", "healthz")
	if err != nil {
		t.Fatalf("normalizedAPIURL returned error: %v", err)
	}
	if got != "http://localhost:8080/healthz" {
		t.Fatalf("expected normalized URL, got %q", got)
	}
}

func TestNormalizedAPIURLRequiresHost(t *testing.T) {
	if _, err := normalizedAPIURL("localhost:8080", "/healthz"); err == nil {
		t.Fatal("expected missing scheme and host to fail")
	}
}

func TestAPISmokeChecksIncludeAdminAudit(t *testing.T) {
	checks := apiSmokeChecks()
	for _, check := range checks {
		if check.Name == "admin audit list" && check.Headers["X-Actor-Type"] == "reseller_owner" {
			return
		}
	}
	t.Fatal("expected admin audit smoke check with reseller actor")
}

func TestAPISmokeChecksIncludeAdminProviderReadiness(t *testing.T) {
	checks := apiSmokeChecks()
	for _, check := range checks {
		if check.Name != "admin provider readiness" {
			continue
		}
		if check.Path != "/admin/catalog/provider-readiness?status=active&limit=20" {
			t.Fatalf("unexpected readiness path %q", check.Path)
		}
		if check.Headers["X-Actor-Type"] != "reseller_owner" {
			t.Fatalf("expected admin actor headers, got %+v", check.Headers)
		}
		for _, expected := range []string{`"plan_display_id":`, `"source_display_id":`, `"state":`, `"reason":`} {
			if !stringSliceContains(check.Contains, expected) {
				t.Fatalf("readiness check missing required token %q", expected)
			}
		}
		for _, blocked := range []string{`"capability_profile"`, `"provider_account_id"`, `"raw_payload"`, `"credentials"`, `"encrypted_payload_ref"`} {
			if !stringSliceContains(check.NotContains, blocked) {
				t.Fatalf("readiness check missing blocked token %q", blocked)
			}
		}
		if !check.RedactBodyOnFailure {
			t.Fatal("expected readiness check to redact failure bodies")
		}
		for _, field := range []string{"plan_display_id", "source_display_id"} {
			if !stringSliceContains(check.SummaryFields, field) {
				t.Fatalf("readiness check missing summary field %q", field)
			}
		}
		return
	}
	t.Fatal("expected admin provider readiness smoke check")
}

func TestRunAPICheckRejectsBlockedFieldsWithoutDumpingBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"plan_display_id":10001,"source_display_id":10002,"state":"ready","reason":"ok","provider_account_id":"acct-secret"}`))
	}))
	defer server.Close()

	check := apiSmokeCheck{
		Name:                "admin provider readiness",
		Path:                "/admin/catalog/provider-readiness",
		Contains:            []string{`"plan_display_id":`, `"source_display_id":`, `"state":`, `"reason":`},
		NotContains:         []string{`"provider_account_id"`},
		RedactBodyOnFailure: true,
	}
	_, err := runAPICheck(context.Background(), server.Client(), server.URL, check)
	if err == nil {
		t.Fatal("expected blocked field to fail the smoke check")
	}
	message := err.Error()
	if !strings.Contains(message, "blocked field") {
		t.Fatalf("expected actionable blocked field error, got %q", message)
	}
	if strings.Contains(message, "acct-secret") {
		t.Fatalf("expected error to avoid dumping response body, got %q", message)
	}
}

func TestRunAPICheckRedactsSensitiveMissingFieldBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"secret":"should-not-leak"}`))
	}))
	defer server.Close()

	check := apiSmokeCheck{
		Name:                "admin provider readiness",
		Path:                "/admin/catalog/provider-readiness",
		Contains:            []string{`"plan_display_id":`},
		RedactBodyOnFailure: true,
	}
	_, err := runAPICheck(context.Background(), server.Client(), server.URL, check)
	if err == nil {
		t.Fatal("expected missing field to fail the smoke check")
	}
	message := err.Error()
	if !strings.Contains(message, "response body omitted") {
		t.Fatalf("expected redacted response body message, got %q", message)
	}
	if strings.Contains(message, "should-not-leak") {
		t.Fatalf("expected error to redact response body, got %q", message)
	}
}

func TestRunAPICheckReturnsDisplayIDSummary(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data":[{"plan_display_id":10001,"source_display_id":10002,"state":"ready","reason":"ok"}]}`))
	}))
	defer server.Close()

	check := apiSmokeCheck{
		Name:          "admin provider readiness",
		Path:          "/admin/catalog/provider-readiness",
		Contains:      []string{`"plan_display_id":`, `"source_display_id":`, `"state":`, `"reason":`},
		SummaryFields: []string{"plan_display_id", "source_display_id"},
	}
	summary, err := runAPICheck(context.Background(), server.Client(), server.URL, check)
	if err != nil {
		t.Fatalf("runAPICheck returned error: %v", err)
	}
	for _, expected := range []string{"display_ids", "plan_display_id=10001", "source_display_id=10002"} {
		if !strings.Contains(summary, expected) {
			t.Fatalf("expected summary to include %q, got %q", expected, summary)
		}
	}
}

func stringSliceContains(values []string, expected string) bool {
	for _, value := range values {
		if value == expected {
			return true
		}
	}
	return false
}

func TestBillingMutationScenarioKeysIncludeRunID(t *testing.T) {
	scenario := billingMutationScenario{RunID: "12345"}

	for _, value := range []string{
		scenario.topupIdempotencyKey(),
		scenario.orderIdempotencyKey(),
		scenario.checkoutIdempotencyKey(),
		scenario.paymentIdempotencyKey(),
		scenario.topupPaymentReference(),
	} {
		if value == "" || value == "12345" {
			t.Fatalf("expected derived billing smoke value, got %q", value)
		}
		if value[len(value)-5:] != "12345" {
			t.Fatalf("expected run id suffix in %q", value)
		}
	}
}

func TestProvisioningJobSmokeStatusOK(t *testing.T) {
	for _, status := range []string{"queued", "claimed", "running", "succeeded"} {
		if !provisioningJobSmokeStatusOK(status) {
			t.Fatalf("expected %q to be accepted", status)
		}
	}
	for _, status := range []string{"", "failed_retryable", "failed_terminal", "manual_review", "cancelled"} {
		if provisioningJobSmokeStatusOK(status) {
			t.Fatalf("expected %q to be rejected", status)
		}
	}
}
