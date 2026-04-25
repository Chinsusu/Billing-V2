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

func TestAPISmokeChecksIncludeAdminPublicIDFilters(t *testing.T) {
	checks := apiSmokeChecks()
	expected := map[string]struct {
		path        string
		contains    []string
		notContains string
	}{
		"admin service public id filter": {
			path:     "/admin/services?display_id=43001&order_display_id=42001&provider_source_display_id=10000",
			contains: []string{`"display_id":43001`, `"order_display_id":42001`, `"provider_source_display_id":10000`},
		},
		"admin invoice public id filter": {
			path:     "/admin/invoices?display_id=44001&buyer_display_id=10002&order_display_id=42001",
			contains: []string{`"display_id":44001`, `"buyer_display_id":10002`, `"order_display_id":42001`},
		},
		"admin audit public id filter": {
			path:     "/admin/audit-logs?actor_display_id=10001&target_type=job&target_display_id=53001",
			contains: []string{`"display_id":70001`, `"actor_display_id":10001`, `"target_display_id":53001`},
		},
		"admin job public id filter": {
			path:     "/admin/jobs?job_type=provider.provision&display_id=53001&source_display_id=10000",
			contains: []string{`"display_id":53001`, `"source_display_id":10000`, `"reference_display_id":42001`},
		},
		"admin invoice public id filter miss": {
			path:        "/admin/invoices?display_id=999999",
			notContains: `"display_id":44001`,
		},
	}
	seen := map[string]bool{}
	for _, check := range checks {
		want, ok := expected[check.Name]
		if !ok {
			continue
		}
		seen[check.Name] = true
		if check.Path != want.path {
			t.Fatalf("unexpected public ID smoke path for %q: %s", check.Name, check.Path)
		}
		if check.Headers["X-Actor-Type"] != "reseller_owner" {
			t.Fatalf("expected admin actor headers for %q, got %+v", check.Name, check.Headers)
		}
		for _, token := range want.contains {
			if !stringSliceContains(check.Contains, token) {
				t.Fatalf("public ID smoke check %q missing contains token %q", check.Name, token)
			}
		}
		if want.notContains != "" && !stringSliceContains(check.NotContains, want.notContains) {
			t.Fatalf("public ID smoke check %q missing not-contains token %q", check.Name, want.notContains)
		}
	}
	for name := range expected {
		if !seen[name] {
			t.Fatalf("missing public ID smoke check %q", name)
		}
	}
}

func TestAPISmokeChecksIncludeRelatedDisplayIDResponses(t *testing.T) {
	checks := apiSmokeChecks()
	expected := map[string][]string{
		"admin wallet list":              {`"owner_display_id":10002`},
		"admin service list":             {`"order_display_id":42001`, `"buyer_display_id":10002`, `"provider_source_display_id":10000`},
		"admin order list":               {`"buyer_display_id":10002`},
		"client wallet list":             {`"owner_display_id":10002`},
		"client order list":              {`"buyer_display_id":10002`},
		"admin topup list":               {`"wallet_display_id":41001`, `"requested_by_display_id":10002`, `"reviewed_by_display_id":10001`},
		"admin transaction list":         {`"account_display_id":10002`, `"order_display_id":42001`, `"invoice_display_id":44001`},
		"admin invoice public id filter": {`"buyer_display_id":10002`, `"order_display_id":42001`},
		"admin job public id filter":     {`"source_display_id":10000`, `"reference_display_id":42001`},
		"admin audit list":               {`"actor_display_id":10001`, `"target_display_id":53001`},
		"admin audit public id filter":   {`"actor_display_id":10001`, `"target_display_id":53001`},
	}
	seen := map[string]bool{}
	for _, check := range checks {
		tokens, ok := expected[check.Name]
		if !ok {
			continue
		}
		seen[check.Name] = true
		for _, token := range tokens {
			if !stringSliceContains(check.Contains, token) {
				t.Fatalf("smoke check %q missing related display ID token %q", check.Name, token)
			}
		}
	}
	for name := range expected {
		if !seen[name] {
			t.Fatalf("missing related display ID smoke check %q", name)
		}
	}
}

func TestAPISmokeChecksIncludeRBACNegativeChecks(t *testing.T) {
	checks := apiRBACNegativeChecks()
	expected := map[string]struct {
		method string
		path   string
	}{
		"deny admin provider readiness": {
			method: http.MethodGet,
			path:   "/admin/catalog/provider-readiness?status=active&limit=20",
		},
		"deny admin job list": {
			method: http.MethodGet,
			path:   "/admin/jobs?job_type=provider.provision&limit=20",
		},
		"deny admin job retry": {
			method: http.MethodPost,
			path:   "/admin/jobs/00000000-0000-0000-0000-000000000999/retry",
		},
	}
	seen := map[string]bool{}
	for _, check := range checks {
		want, ok := expected[check.Name]
		if !ok {
			t.Fatalf("unexpected RBAC negative check %q", check.Name)
		}
		seen[check.Name] = true
		if check.Method != want.method || check.Path != want.path {
			t.Fatalf("unexpected RBAC check route for %q: %s %s", check.Name, check.Method, check.Path)
		}
		if check.WantStatus != http.StatusForbidden || check.WantCode != "auth.permission_denied" {
			t.Fatalf("unexpected RBAC expected error for %q: %d %s", check.Name, check.WantStatus, check.WantCode)
		}
		if check.Headers["X-Actor-Id"] != demoNoPermissionActorID || check.Headers["X-Actor-Type"] != "client" {
			t.Fatalf("expected low-permission actor headers, got %+v", check.Headers)
		}
		for _, blocked := range []string{`"payload_json"`, `"provider_account_id"`, `"raw_response"`, `"provider.provision"`, `"order_display_id"`} {
			if !stringSliceContains(check.NotContains, blocked) {
				t.Fatalf("RBAC check %q missing blocked token %q", check.Name, blocked)
			}
		}
	}
	for name := range expected {
		if !seen[name] {
			t.Fatalf("missing RBAC negative check %q", name)
		}
	}
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

func TestRunAPIRBACNegativeCheckAcceptsPermissionDeniedEnvelope(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte(`{"error":{"code":"auth.permission_denied","message":"Permission denied."},"request_id":"req_test"}`))
	}))
	defer server.Close()

	check := apiRBACNegativeCheck{
		Name:        "deny admin job list",
		Method:      http.MethodGet,
		Path:        "/admin/jobs",
		Headers:     lowPermissionHeaders(),
		WantStatus:  http.StatusForbidden,
		WantCode:    "auth.permission_denied",
		NotContains: sensitiveAPIRedactionTokens(),
	}
	if err := runAPIRBACNegativeCheck(context.Background(), server.Client(), server.URL, check); err != nil {
		t.Fatalf("runAPIRBACNegativeCheck returned error: %v", err)
	}
}

func TestRunAPIRBACNegativeCheckRejectsUnexpectedSuccessWithoutDumpingBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data":[{"payload_json":{"secret":"should-not-leak"}}]}`))
	}))
	defer server.Close()

	check := apiRBACNegativeCheck{
		Name:        "deny admin job list",
		Method:      http.MethodGet,
		Path:        "/admin/jobs",
		Headers:     lowPermissionHeaders(),
		WantStatus:  http.StatusForbidden,
		WantCode:    "auth.permission_denied",
		NotContains: sensitiveAPIRedactionTokens(),
	}
	err := runAPIRBACNegativeCheck(context.Background(), server.Client(), server.URL, check)
	if err == nil {
		t.Fatal("expected unexpected success to fail")
	}
	message := err.Error()
	if !strings.Contains(message, "expected HTTP 403, got 200") {
		t.Fatalf("expected status failure, got %q", message)
	}
	if strings.Contains(message, "should-not-leak") || strings.Contains(message, "payload_json") {
		t.Fatalf("expected error to avoid dumping response body, got %q", message)
	}
}

func TestRunAPIRBACNegativeCheckRejectsDeniedLeakWithoutDumpingSecret(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte(`{"error":{"code":"auth.permission_denied","message":"Permission denied.","details":{"payload_json":{"secret":"should-not-leak"}}},"request_id":"req_test"}`))
	}))
	defer server.Close()

	check := apiRBACNegativeCheck{
		Name:        "deny admin job list",
		Method:      http.MethodGet,
		Path:        "/admin/jobs",
		Headers:     lowPermissionHeaders(),
		WantStatus:  http.StatusForbidden,
		WantCode:    "auth.permission_denied",
		NotContains: sensitiveAPIRedactionTokens(),
	}
	err := runAPIRBACNegativeCheck(context.Background(), server.Client(), server.URL, check)
	if err == nil {
		t.Fatal("expected leaked denied response to fail")
	}
	message := err.Error()
	if !strings.Contains(message, "blocked field") {
		t.Fatalf("expected blocked token failure, got %q", message)
	}
	if strings.Contains(message, "should-not-leak") {
		t.Fatalf("expected error to avoid dumping leaked secret, got %q", message)
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
