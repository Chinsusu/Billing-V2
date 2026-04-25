package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const demoNoPermissionActorID = "00000000-0000-0000-0000-000000000199"

type apiRBACNegativeCheck struct {
	Name        string
	Method      string
	Path        string
	Headers     map[string]string
	Body        string
	WantStatus  int
	WantCode    string
	NotContains []string
}

func runAPIRBACNegativeCheck(ctx context.Context, client *http.Client, baseURL string, check apiRBACNegativeCheck) error {
	method := strings.TrimSpace(check.Method)
	if method == "" {
		method = http.MethodGet
	}
	wantStatus := check.WantStatus
	if wantStatus == 0 {
		wantStatus = http.StatusForbidden
	}
	wantCode := strings.TrimSpace(check.WantCode)
	if wantCode == "" {
		wantCode = "auth.permission_denied"
	}

	fullURL, err := normalizedAPIURL(baseURL, check.Path)
	if err != nil {
		return err
	}
	var body io.Reader
	if strings.TrimSpace(check.Body) != "" {
		body = strings.NewReader(check.Body)
	}
	request, err := http.NewRequestWithContext(ctx, method, fullURL, body)
	if err != nil {
		return fmt.Errorf("build RBAC request %q: %w", check.Name, err)
	}
	if body != nil {
		request.Header.Set("Content-Type", "application/json")
	}
	for key, value := range check.Headers {
		request.Header.Set(key, value)
	}

	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("request RBAC check %q: %w", check.Name, err)
	}
	defer response.Body.Close()

	payload, err := io.ReadAll(io.LimitReader(response.Body, 1<<20))
	if err != nil {
		return fmt.Errorf("read RBAC response %q: %w", check.Name, err)
	}
	bodyText := string(payload)
	if response.StatusCode != wantStatus {
		return fmt.Errorf("RBAC check %q expected HTTP %d, got %d: response body omitted for RBAC smoke check", check.Name, wantStatus, response.StatusCode)
	}
	if err := assertResponseOmitsTokens(check.Name, bodyText, check.NotContains); err != nil {
		return err
	}

	var envelope errorEnvelope
	if err := json.Unmarshal(payload, &envelope); err != nil {
		return fmt.Errorf("RBAC check %q response is not valid JSON error envelope: %w", check.Name, err)
	}
	if envelope.Error.Code != wantCode {
		return fmt.Errorf("RBAC check %q expected error code %q, got %q", check.Name, wantCode, envelope.Error.Code)
	}
	if strings.TrimSpace(envelope.Error.Message) == "" {
		return fmt.Errorf("RBAC check %q error envelope missing message", check.Name)
	}
	if strings.TrimSpace(envelope.RequestID) == "" {
		return fmt.Errorf("RBAC check %q error envelope missing request_id", check.Name)
	}
	return nil
}

func apiRBACNegativeChecks() []apiRBACNegativeCheck {
	lowPermission := lowPermissionHeaders()
	blocked := sensitiveAPIRedactionTokens()
	return []apiRBACNegativeCheck{
		{
			Name:        "deny admin provider readiness",
			Method:      http.MethodGet,
			Path:        "/admin/catalog/provider-readiness?status=active&limit=20",
			Headers:     lowPermission,
			WantStatus:  http.StatusForbidden,
			WantCode:    "auth.permission_denied",
			NotContains: blocked,
		},
		{
			Name:        "deny admin job list",
			Method:      http.MethodGet,
			Path:        "/admin/jobs?job_type=provider.provision&limit=20",
			Headers:     lowPermission,
			WantStatus:  http.StatusForbidden,
			WantCode:    "auth.permission_denied",
			NotContains: blocked,
		},
		{
			Name:        "deny admin job retry",
			Method:      http.MethodPost,
			Path:        "/admin/jobs/00000000-0000-0000-0000-000000000999/retry",
			Headers:     lowPermission,
			Body:        `{}`,
			WantStatus:  http.StatusForbidden,
			WantCode:    "auth.permission_denied",
			NotContains: blocked,
		},
	}
}

func lowPermissionHeaders() map[string]string {
	return actorHeaders(demoNoPermissionActorID, "client")
}

func sensitiveAPIRedactionTokens() []string {
	return []string{
		`"capabilities"`,
		`"capability_override"`,
		`"capability_json"`,
		`"capability_profile"`,
		`"capability_snapshot"`,
		`"access_token"`,
		`"api_key"`,
		`"credential"`,
		`"credentials"`,
		`"encrypted_payload_ref"`,
		`"idempotency_key"`,
		`"order_display_id"`,
		`"order_id"`,
		`"payload"`,
		`"payload_json"`,
		`"provider.provision"`,
		`"provider_account_id"`,
		`"provider_credential"`,
		`"provider_credentials"`,
		`"raw_payload"`,
		`"raw_response"`,
		`"secret"`,
		`"secret_version"`,
		`"source_id"`,
		`"token"`,
	}
}

func assertResponseOmitsTokens(checkName string, bodyText string, blockedTokens []string) error {
	bodyLower := strings.ToLower(bodyText)
	for _, blocked := range blockedTokens {
		if strings.Contains(bodyLower, strings.ToLower(blocked)) {
			return fmt.Errorf("check %q response exposed blocked field %q", checkName, blocked)
		}
	}
	return nil
}
