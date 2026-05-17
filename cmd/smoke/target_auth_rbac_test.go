package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLoginForTargetAuthSmokeExtractsHttpOnlyCookie(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/auth/login" || r.Method != http.MethodPost {
			t.Fatalf("unexpected request %s %s", r.Method, r.URL.Path)
		}
		if r.Header.Get("X-Tenant-Id") != demoTenantID {
			t.Fatalf("expected tenant header")
		}
		http.SetCookie(w, &http.Cookie{Name: "billing_session", Value: "session-token", HttpOnly: true})
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"session_id":"session_1","user_id":"user_1","tenant_id":"tenant_1","actor_type":"client","two_factor_required":false,"two_factor_satisfied":false},"request_id":"req_test"}`))
	}))
	defer server.Close()

	data, cookie, err := loginForTargetAuthSmoke(context.Background(), server.Client(), server.URL, "billing_session", demoTenantID, targetAuthSmokeClientEmail)
	if err != nil {
		t.Fatalf("loginForTargetAuthSmoke returned error: %v", err)
	}
	if data.ActorType != "client" || cookie == nil || cookie.Name != "billing_session" || cookie.Value != "session-token" || !cookie.HttpOnly {
		t.Fatalf("unexpected login result: data=%+v cookie=%+v", data, cookie)
	}
}

func TestRunTargetSessionGETCheckUsesCookieOnlyAuth(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Actor-Id") != "" {
			t.Fatalf("session check should not send dev actor headers")
		}
		cookie, err := r.Cookie("billing_session")
		if err != nil || cookie.Value != "session-token" {
			t.Fatalf("expected session cookie")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"status":"ok"},"request_id":"req_test"}`))
	}))
	defer server.Close()

	cookie := &http.Cookie{Name: "billing_session", Value: "session-token"}
	if err := runTargetSessionGETCheck(context.Background(), server.Client(), server.URL, "client session", "/client/catalog", cookie); err != nil {
		t.Fatalf("runTargetSessionGETCheck returned error: %v", err)
	}
}

func TestRunTargetSessionGETCheckRejectsUnexpectedStatusWithoutLeakingBody(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"session_token":"secret-value"}`))
	}))
	defer server.Close()

	err := runTargetSessionGETCheck(context.Background(), server.Client(), server.URL, "client session", "/client/catalog", &http.Cookie{Name: "billing_session", Value: "session-token"})
	if err == nil {
		t.Fatal("expected error")
	}
	if strings.Contains(err.Error(), "secret-value") || strings.Contains(err.Error(), "session_token") {
		t.Fatalf("error leaked response body: %v", err)
	}
}

func TestTargetAuthRBACNegativeChecksIncludeSessionAndTenantCases(t *testing.T) {
	checks := targetAuthRBACNegativeChecks("billing_session")
	want := map[string]string{
		"deny invalid session cookie":      "auth.session_invalid",
		"deny missing actor context":       "auth.actor_required",
		"deny cross-tenant actor mismatch": "tenant.context_mismatch",
	}
	for _, check := range checks {
		delete(want, check.Name)
		if check.WantStatus == 0 || check.WantCode == "" {
			t.Fatalf("check %q missing expected status/code", check.Name)
		}
	}
	if len(want) != 0 {
		t.Fatalf("missing target auth RBAC checks: %+v", want)
	}
}
