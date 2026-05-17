package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	platformTenantID             = "00000000-0000-0000-0000-000000000001"
	targetAuthSmokeSeedPassword  = "admin123"
	targetAuthSmokeClientEmail   = "customer@local.billing"
	targetAuthSmokePlatformEmail = "admin@local.billing"
)

type targetAuthLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type targetAuthLoginResponse struct {
	SessionID              string `json:"session_id"`
	UserID                 string `json:"user_id"`
	TenantID               string `json:"tenant_id"`
	ActorType              string `json:"actor_type"`
	TwoFactorRequired      bool   `json:"two_factor_required"`
	TwoFactorSatisfied     bool   `json:"two_factor_satisfied"`
	TwoFactorSetupRequired bool   `json:"two_factor_setup_required"`
}

func runDevTargetAuthRBACSmoke(baseURL string, timeout time.Duration) error {
	if err := guardDevEnvironment(); err != nil {
		return err
	}
	baseURL = strings.TrimSpace(baseURL)
	if baseURL == "" {
		return fmt.Errorf("API_BASE_URL or -base-url is required for dev-target-auth-rbac smoke")
	}
	if _, err := normalizedAPIURL(baseURL, "/healthz"); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	client := &http.Client{Timeout: timeout}
	cookieName := targetAuthSessionCookieName()

	clientLogin, clientCookie, err := loginForTargetAuthSmoke(ctx, client, baseURL, cookieName, demoTenantID, targetAuthSmokeClientEmail)
	if err != nil {
		return err
	}
	defer func() { _ = logoutTargetAuthSmoke(context.Background(), client, baseURL, clientCookie) }()
	if clientLogin.ActorType != "client" || clientLogin.TwoFactorRequired || clientLogin.TwoFactorSatisfied {
		return fmt.Errorf("expected client login without 2FA requirement")
	}
	if err := runTargetSessionGETCheck(ctx, client, baseURL, "client cookie-only catalog", "/client/catalog", clientCookie); err != nil {
		return err
	}

	adminLogin, adminCookie, err := loginForTargetAuthSmoke(ctx, client, baseURL, cookieName, platformTenantID, targetAuthSmokePlatformEmail)
	if err != nil {
		return err
	}
	defer func() { _ = logoutTargetAuthSmoke(context.Background(), client, baseURL, adminCookie) }()
	if adminLogin.ActorType != "platform_staff" || !adminLogin.TwoFactorRequired || adminLogin.TwoFactorSatisfied {
		return fmt.Errorf("expected platform staff login to require unsatisfied 2FA")
	}
	if err := runAPIRBACNegativeCheck(ctx, client, baseURL, apiRBACNegativeCheck{
		Name:        "deny unsatisfied admin 2FA session",
		Method:      http.MethodGet,
		Path:        "/admin/catalog/provider-readiness?status=active&limit=1",
		Headers:     cookieHeader(adminCookie),
		WantStatus:  http.StatusForbidden,
		WantCode:    "auth.2fa_required",
		NotContains: targetAuthSensitiveTokens(),
	}); err != nil {
		return err
	}

	negativeChecks := targetAuthRBACNegativeChecks(cookieName)
	for _, check := range negativeChecks {
		if err := runAPIRBACNegativeCheck(ctx, client, baseURL, check); err != nil {
			return err
		}
	}
	rbacChecks := apiRBACNegativeChecks()
	for _, check := range rbacChecks {
		if err := runAPIRBACNegativeCheck(ctx, client, baseURL, check); err != nil {
			return err
		}
	}

	fmt.Printf("target auth RBAC smoke passed: client_session_cookie_only=pass admin_2fa_gate=pass invalid_session_denied=pass actor_required_denied=pass tenant_mismatch_denied=pass rbac_denials=%d provider_mutation_routes_called=no money_mutation_routes_called=no\n",
		len(rbacChecks),
	)
	fmt.Println("Target auth RBAC smoke output intentionally excludes raw session tokens, cookies, passwords, DSNs, provider payloads, and credentials.")
	return nil
}

func loginForTargetAuthSmoke(ctx context.Context, client *http.Client, baseURL string, cookieName string, tenantID string, email string) (targetAuthLoginResponse, *http.Cookie, error) {
	var zero targetAuthLoginResponse
	fullURL, err := normalizedAPIURL(baseURL, "/auth/login")
	if err != nil {
		return zero, nil, err
	}
	payload, err := json.Marshal(targetAuthLoginRequest{Email: email, Password: targetAuthSmokeSeedPassword})
	if err != nil {
		return zero, nil, fmt.Errorf("marshal target auth login request: %w", err)
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, fullURL, bytes.NewReader(payload))
	if err != nil {
		return zero, nil, fmt.Errorf("build target auth login request: %w", err)
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Tenant-Id", tenantID)

	response, err := client.Do(request)
	if err != nil {
		return zero, nil, fmt.Errorf("request target auth login: %w", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(io.LimitReader(response.Body, 1<<20))
	if err != nil {
		return zero, nil, fmt.Errorf("read target auth login response: %w", err)
	}
	if response.StatusCode != http.StatusOK {
		return zero, nil, targetAuthStatusError("target auth login", http.StatusOK, response.StatusCode, body)
	}

	var envelope struct {
		Data      targetAuthLoginResponse `json:"data"`
		RequestID string                  `json:"request_id"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		return zero, nil, fmt.Errorf("decode target auth login response: %w", err)
	}
	if strings.TrimSpace(envelope.RequestID) == "" {
		return zero, nil, fmt.Errorf("target auth login response missing request_id")
	}
	if envelope.Data.SessionID == "" || envelope.Data.UserID == "" || envelope.Data.TenantID == "" || envelope.Data.ActorType == "" {
		return zero, nil, fmt.Errorf("target auth login response missing session identity fields")
	}

	cookie := findResponseCookie(response.Cookies(), cookieName)
	if cookie == nil || strings.TrimSpace(cookie.Value) == "" {
		return zero, nil, fmt.Errorf("target auth login did not set session cookie")
	}
	if !cookie.HttpOnly {
		return zero, nil, fmt.Errorf("target auth login session cookie is not HttpOnly")
	}
	return envelope.Data, cookie, nil
}

func runTargetSessionGETCheck(ctx context.Context, client *http.Client, baseURL string, name string, path string, cookie *http.Cookie) error {
	fullURL, err := normalizedAPIURL(baseURL, path)
	if err != nil {
		return err
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)
	if err != nil {
		return fmt.Errorf("build target session check %q: %w", name, err)
	}
	request.AddCookie(&http.Cookie{Name: cookie.Name, Value: cookie.Value})

	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("request target session check %q: %w", name, err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(io.LimitReader(response.Body, 1<<20))
	if err != nil {
		return fmt.Errorf("read target session check %q: %w", name, err)
	}
	if response.StatusCode != http.StatusOK {
		return targetAuthStatusError(name, http.StatusOK, response.StatusCode, body)
	}
	if err := assertResponseOmitsTokens(name, string(body), targetAuthSensitiveTokens()); err != nil {
		return err
	}
	var envelope struct {
		Data      json.RawMessage `json:"data"`
		RequestID string          `json:"request_id"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		return fmt.Errorf("decode target session check %q response: %w", name, err)
	}
	if len(envelope.Data) == 0 || strings.TrimSpace(envelope.RequestID) == "" {
		return fmt.Errorf("target session check %q response missing data or request_id", name)
	}
	return nil
}

func logoutTargetAuthSmoke(ctx context.Context, client *http.Client, baseURL string, cookie *http.Cookie) error {
	if cookie == nil {
		return nil
	}
	fullURL, err := normalizedAPIURL(baseURL, "/auth/logout")
	if err != nil {
		return err
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, fullURL, nil)
	if err != nil {
		return fmt.Errorf("build target auth logout request: %w", err)
	}
	request.AddCookie(&http.Cookie{Name: cookie.Name, Value: cookie.Value})
	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("request target auth logout: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("target auth logout expected HTTP 200, got %d", response.StatusCode)
	}
	return nil
}

func targetAuthRBACNegativeChecks(cookieName string) []apiRBACNegativeCheck {
	crossTenantHeaders := actorHeaders(demoCustomerID, "client")
	crossTenantHeaders["X-Actor-Tenant-Id"] = platformTenantID
	return []apiRBACNegativeCheck{
		{
			Name:        "deny invalid session cookie",
			Method:      http.MethodGet,
			Path:        "/client/catalog",
			Headers:     map[string]string{"Cookie": cookieName + "=invalid-target-auth-smoke"},
			WantStatus:  http.StatusUnauthorized,
			WantCode:    "auth.session_invalid",
			NotContains: targetAuthSensitiveTokens(),
		},
		{
			Name:        "deny missing actor context",
			Method:      http.MethodGet,
			Path:        "/client/catalog",
			Headers:     map[string]string{"X-Tenant-Id": demoTenantID},
			WantStatus:  http.StatusUnauthorized,
			WantCode:    "auth.actor_required",
			NotContains: targetAuthSensitiveTokens(),
		},
		{
			Name:        "deny cross-tenant actor mismatch",
			Method:      http.MethodGet,
			Path:        "/client/catalog",
			Headers:     crossTenantHeaders,
			WantStatus:  http.StatusForbidden,
			WantCode:    "tenant.context_mismatch",
			NotContains: targetAuthSensitiveTokens(),
		},
	}
}

func targetAuthStatusError(checkName string, wantStatus int, gotStatus int, body []byte) error {
	var apiError errorEnvelope
	if err := json.Unmarshal(body, &apiError); err == nil && apiError.Error.Code != "" {
		return fmt.Errorf("%s expected HTTP %d, got %d (%s)", checkName, wantStatus, gotStatus, apiError.Error.Code)
	}
	return fmt.Errorf("%s expected HTTP %d, got %d", checkName, wantStatus, gotStatus)
}

func targetAuthSessionCookieName() string {
	if value := strings.TrimSpace(os.Getenv("AUTH_SESSION_COOKIE_NAME")); value != "" {
		return value
	}
	return "billing_session"
}

func findResponseCookie(cookies []*http.Cookie, name string) *http.Cookie {
	for _, cookie := range cookies {
		if cookie != nil && cookie.Name == name {
			return cookie
		}
	}
	return nil
}

func cookieHeader(cookie *http.Cookie) map[string]string {
	if cookie == nil {
		return map[string]string{}
	}
	return map[string]string{"Cookie": cookie.Name + "=" + cookie.Value}
}

func targetAuthSensitiveTokens() []string {
	return []string{
		"billing_session",
		"password",
		"session_token",
		"token_hash",
		"reset_token",
		"cookie",
	}
}
