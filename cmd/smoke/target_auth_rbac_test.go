package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestLoginForTargetAuthSmokeExtractsHttpOnlyCookie(t *testing.T) {
	const (
		expectedEmail    = "client@example.test"
		expectedPassword = "client-pw"
	)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/auth/login" || r.Method != http.MethodPost {
			t.Fatalf("unexpected request %s %s", r.Method, r.URL.Path)
		}
		if r.Header.Get("X-Tenant-Id") != demoTenantID {
			t.Fatalf("expected tenant header")
		}
		var request targetAuthLoginRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Fatalf("decode request body: %v", err)
		}
		if request.Email != expectedEmail || request.Password != expectedPassword {
			t.Fatalf("unexpected login request fields")
		}
		http.SetCookie(w, &http.Cookie{Name: "billing_session", Value: "session-token", HttpOnly: true})
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"data":{"session_id":"session_1","user_id":"user_1","tenant_id":"tenant_1","actor_type":"client","two_factor_required":false,"two_factor_satisfied":false},"request_id":"req_test"}`))
	}))
	defer server.Close()

	data, cookie, err := loginForTargetAuthSmoke(context.Background(), server.Client(), server.URL, "billing_session", demoTenantID, expectedEmail, expectedPassword)
	if err != nil {
		t.Fatalf("loginForTargetAuthSmoke returned error: %v", err)
	}
	if data.ActorType != "client" || cookie == nil || cookie.Name != "billing_session" || cookie.Value != "session-token" || !cookie.HttpOnly {
		t.Fatalf("unexpected login result: data=%+v cookie=%+v", data, cookie)
	}
}

func TestTargetAuthSmokeCredentialsFromEnvUsesOverrides(t *testing.T) {
	t.Setenv(targetAuthSmokeClientEmailEnvName, "client@example.test")
	t.Setenv(targetAuthSmokeClientPasswordEnvName, "client-pw")
	t.Setenv(targetAuthSmokeAdminEmailEnvName, "admin@example.test")
	t.Setenv(targetAuthSmokeAdminPasswordEnvName, "admin-pw")

	credentials := targetAuthSmokeCredentialsFromEnv()
	if credentials.ClientEmail != "client@example.test" || credentials.ClientPassword != "client-pw" {
		t.Fatalf("unexpected client credential overrides")
	}
	if credentials.AdminEmail != "admin@example.test" || credentials.AdminPassword != "admin-pw" {
		t.Fatalf("unexpected admin credential overrides")
	}
}

func TestTargetAuthSmokeBaseURLsFromInputsDefaultsToBaseURL(t *testing.T) {
	baseURLs := targetAuthSmokeBaseURLsFromInputs(" http://api.test ", "", "")
	if baseURLs.Client != "http://api.test" || baseURLs.Admin != "http://api.test" {
		t.Fatalf("expected client and admin base URLs to default to base URL")
	}
}

func TestRunDevTargetAuthRBACSmokeUsesSeparateBaseURLs(t *testing.T) {
	const (
		clientEmail    = "client@example.test"
		clientPassword = "client-pw"
		adminEmail     = "admin@example.test"
		adminPassword  = "admin-pw"
	)
	t.Setenv(targetAuthSmokeClientEmailEnvName, clientEmail)
	t.Setenv(targetAuthSmokeClientPasswordEnvName, clientPassword)
	t.Setenv(targetAuthSmokeAdminEmailEnvName, adminEmail)
	t.Setenv(targetAuthSmokeAdminPasswordEnvName, adminPassword)

	clientServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/auth/login":
			if r.Method != http.MethodPost || r.Header.Get("X-Tenant-Id") != demoTenantID {
				t.Fatalf("unexpected client login request")
			}
			var request targetAuthLoginRequest
			if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
				t.Fatalf("decode client login request: %v", err)
			}
			if request.Email != clientEmail || request.Password != clientPassword {
				t.Fatalf("unexpected client login request fields")
			}
			http.SetCookie(w, &http.Cookie{Name: "billing_session", Value: "client-session", HttpOnly: true})
			writeTargetAuthTestSuccess(w, `{"session_id":"session_client","user_id":"user_client","tenant_id":"tenant_client","actor_type":"client","two_factor_required":false,"two_factor_satisfied":false}`)
		case "/client/catalog":
			if cookie, err := r.Cookie("billing_session"); err == nil {
				if cookie.Value == "client-session" {
					writeTargetAuthTestSuccess(w, `{"catalog":"ok"}`)
					return
				}
				writeTargetAuthTestError(w, http.StatusUnauthorized, "auth.session_invalid")
				return
			}
			if r.Header.Get("X-Actor-Id") == "" {
				writeTargetAuthTestError(w, http.StatusUnauthorized, "auth.actor_required")
				return
			}
			if r.Header.Get("X-Actor-Tenant-Id") == platformTenantID {
				writeTargetAuthTestError(w, http.StatusForbidden, "tenant.context_mismatch")
				return
			}
			writeTargetAuthTestError(w, http.StatusForbidden, "auth.permission_denied")
		case "/auth/logout":
			if r.Method != http.MethodPost {
				t.Fatalf("unexpected client logout method")
			}
			writeTargetAuthTestSuccess(w, `{"status":"logged_out"}`)
		default:
			t.Fatalf("unexpected client server path %s", r.URL.Path)
		}
	}))
	defer clientServer.Close()

	adminServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/auth/login":
			if r.Method != http.MethodPost || r.Header.Get("X-Tenant-Id") != platformTenantID {
				t.Fatalf("unexpected admin login request")
			}
			var request targetAuthLoginRequest
			if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
				t.Fatalf("decode admin login request: %v", err)
			}
			if request.Email != adminEmail || request.Password != adminPassword {
				t.Fatalf("unexpected admin login request fields")
			}
			http.SetCookie(w, &http.Cookie{Name: "billing_session", Value: "admin-session", HttpOnly: true})
			writeTargetAuthTestSuccess(w, `{"session_id":"session_admin","user_id":"user_admin","tenant_id":"tenant_admin","actor_type":"platform_staff","two_factor_required":true,"two_factor_satisfied":false}`)
		case "/admin/catalog/provider-readiness":
			if cookie, err := r.Cookie("billing_session"); err == nil && cookie.Value == "admin-session" {
				writeTargetAuthTestError(w, http.StatusForbidden, "auth.2fa_required")
				return
			}
			writeTargetAuthTestError(w, http.StatusForbidden, "auth.permission_denied")
		case "/admin/jobs":
			writeTargetAuthTestError(w, http.StatusForbidden, "auth.permission_denied")
		case "/admin/jobs/00000000-0000-0000-0000-000000000999/retry":
			if r.Method != http.MethodPost {
				t.Fatalf("unexpected admin retry method")
			}
			writeTargetAuthTestError(w, http.StatusForbidden, "auth.permission_denied")
		case "/auth/logout":
			if r.Method != http.MethodPost {
				t.Fatalf("unexpected admin logout method")
			}
			writeTargetAuthTestSuccess(w, `{"status":"logged_out"}`)
		default:
			t.Fatalf("unexpected admin server path %s", r.URL.Path)
		}
	}))
	defer adminServer.Close()

	err := runDevTargetAuthRBACSmoke("http://unused.invalid", clientServer.URL, adminServer.URL, 5*time.Second)
	if err != nil {
		t.Fatalf("runDevTargetAuthRBACSmoke returned error: %v", err)
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

func writeTargetAuthTestSuccess(w http.ResponseWriter, data string) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"data":` + data + `,"request_id":"req_test"}`))
}

func writeTargetAuthTestError(w http.ResponseWriter, status int, code string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write([]byte(`{"error":{"code":"` + code + `","message":"denied"},"request_id":"req_test"}`))
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
