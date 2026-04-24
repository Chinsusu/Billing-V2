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
