package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetTargetFinanceJSONUsesReadOnlyFinanceHeaders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET request")
		}
		if r.Header.Get("Idempotency-Key") != "" {
			t.Fatalf("finance reconciliation smoke must not send idempotency key")
		}
		if r.Header.Get("X-Actor-Id") != demoResellerID ||
			r.Header.Get("X-Actor-Type") != targetFinanceSmokeActorType ||
			r.Header.Get("X-Actor-Tenant-Id") != demoTenantID ||
			r.Header.Get("X-Tenant-Id") != demoTenantID {
			t.Fatalf("unexpected finance headers")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"date":"2026-04-23","status":"balanced","wallets":{"checked":1,"balanced":1,"mismatched":0},"invoices":{"checked":1,"mismatched":0},"payments":{"checked":1,"duplicate_reference_count":0}},"request_id":"req_test"}`))
	}))
	defer server.Close()

	report, err := getTargetFinanceJSON[targetFinanceDailyReconciliationResponse](context.Background(), server.Client(), server.URL, "/admin/daily-reconciliation?date=2026-04-23")
	if err != nil {
		t.Fatalf("getTargetFinanceJSON returned error: %v", err)
	}
	if report.Status != "balanced" || report.Payments.DuplicateReferenceCount != 0 {
		t.Fatalf("unexpected report: %+v", report)
	}
}

func TestTargetFinanceStatusErrorDoesNotLeakBody(t *testing.T) {
	err := targetFinanceStatusError(http.StatusInternalServerError, []byte(`{"token_hash":"secret-value"}`))
	if err == nil {
		t.Fatal("expected error")
	}
	if strings.Contains(err.Error(), "secret-value") || strings.Contains(err.Error(), "token_hash") {
		t.Fatalf("error leaked response body: %v", err)
	}
}

func TestValidateTargetFinanceDailyReconciliationAcceptsMismatchedEvidence(t *testing.T) {
	report := targetFinanceDailyReconciliationResponse{
		Date:   "2026-04-23",
		Status: "mismatched",
	}
	report.Wallets.Checked = 2
	report.Wallets.Balanced = 1
	report.Wallets.Mismatched = 1
	report.Invoices.Checked = 1
	report.Payments.Checked = 1
	if err := validateTargetFinanceDailyReconciliation(report, targetFinanceCandidate{ReconciliationDate: "2026-04-23"}); err != nil {
		t.Fatalf("expected mismatched evidence to pass: %v", err)
	}
}

func TestValidateTargetFinanceDailyReconciliationRejectsInconsistentMismatch(t *testing.T) {
	err := validateTargetFinanceDailyReconciliation(targetFinanceDailyReconciliationResponse{
		Date:   "2026-04-23",
		Status: "mismatched",
		Wallets: struct {
			Checked    int `json:"checked"`
			Balanced   int `json:"balanced"`
			Mismatched int `json:"mismatched"`
		}{Checked: 1, Balanced: 1, Mismatched: 0},
		Invoices: struct {
			Checked    int `json:"checked"`
			Mismatched int `json:"mismatched"`
		}{Checked: 1},
		Payments: struct {
			Checked                 int `json:"checked"`
			DuplicateReferenceCount int `json:"duplicate_reference_count"`
		}{Checked: 1},
	}, targetFinanceCandidate{ReconciliationDate: "2026-04-23"})
	if err == nil {
		t.Fatal("expected inconsistent mismatched report to fail")
	}
}
