package main

import "testing"

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
