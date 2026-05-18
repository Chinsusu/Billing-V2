package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/Chinsusu/Billing-V2/internal/modules/provider"
)

const cloudminiPermissionDeniedExampleName = "permission_denied_proxy_list"

type cloudminiAPIKeyListResponse struct {
	Success bool `json:"success"`
	Data    []struct {
		ID       string `json:"id"`
		IsActive bool   `json:"is_active"`
	} `json:"data"`
}

type cloudminiAPIKeyCreateResponse struct {
	Success bool `json:"success"`
	Data    struct {
		APIKey struct {
			ID       string `json:"id"`
			IsActive bool   `json:"is_active"`
		} `json:"api_key"`
		PlainKey string `json:"plain_key"`
	} `json:"data"`
}

func runCloudminiPermissionDeniedEvidence(ctx context.Context, config cloudminiErrorEvidenceConfig) (cloudminiErrorEvidenceResult, error) {
	result := cloudminiErrorEvidenceResult{
		Name:              cloudminiPermissionDeniedExampleName,
		ProviderCode:      "none",
		NormalizedCode:    provider.ErrorPermissionDenied,
		RetrySafety:       provider.RetrySafetyDoNotRetry,
		SideEffectCreated: "no",
		TemporaryKey:      true,
	}

	activeBefore, err := countCloudminiActiveAPIKeys(ctx, config)
	if err != nil {
		return result, err
	}

	keyID, plainKey, err := createCloudminiReadOnlyAPIKey(ctx, config)
	if err != nil {
		return result, err
	}

	deniedResult, deniedErr := callCloudminiProxyListWithTemporaryKey(ctx, config, plainKey)
	result = deniedResult

	revokeErr := revokeCloudminiAPIKey(ctx, config, keyID)
	if revokeErr == nil {
		result.TemporaryKeyRevoked = true
		result.SideEffectCreated = "cleaned_up"
	}
	result.MutatingRoute = true

	activeAfter, countErr := countCloudminiActiveAPIKeys(ctx, config)
	if countErr == nil && activeAfter == activeBefore {
		result.ActiveKeyCountRestored = true
	}

	switch {
	case revokeErr != nil:
		return result, revokeErr
	case countErr != nil:
		return result, countErr
	case !result.ActiveKeyCountRestored:
		return result, fmt.Errorf("permission-denied temporary api key active count was not restored")
	case deniedErr != nil:
		return result, deniedErr
	default:
		return result, nil
	}
}

func countCloudminiActiveAPIKeys(ctx context.Context, config cloudminiErrorEvidenceConfig) (int, error) {
	request, err := newCloudminiManagementRequest(ctx, config, http.MethodGet, "/api/v1/api-keys/", nil)
	if err != nil {
		return 0, err
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return 0, fmt.Errorf("permission-denied api key list request failed before response")
	}
	defer response.Body.Close()
	body, err := io.ReadAll(io.LimitReader(response.Body, 1<<20))
	if err != nil {
		return 0, fmt.Errorf("permission-denied api key list response could not be read")
	}
	if response.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("permission-denied api key list returned unexpected HTTP status %d", response.StatusCode)
	}
	var parsed cloudminiAPIKeyListResponse
	if err := json.Unmarshal(body, &parsed); err != nil || !parsed.Success {
		return 0, fmt.Errorf("permission-denied api key list response envelope was invalid")
	}
	active := 0
	for _, item := range parsed.Data {
		if item.IsActive {
			active++
		}
	}
	return active, nil
}

func createCloudminiReadOnlyAPIKey(ctx context.Context, config cloudminiErrorEvidenceConfig) (string, string, error) {
	payload := bytes.NewBufferString(`{"name":"billing-t256-permission-denied","permissions":["read"]}`)
	request, err := newCloudminiManagementRequest(ctx, config, http.MethodPost, "/api/v1/api-keys/", payload)
	if err != nil {
		return "", "", err
	}
	request.Header.Set("Content-Type", "application/json")
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return "", "", fmt.Errorf("permission-denied api key create request failed before response")
	}
	defer response.Body.Close()
	body, err := io.ReadAll(io.LimitReader(response.Body, 1<<20))
	if err != nil {
		return "", "", fmt.Errorf("permission-denied api key create response could not be read")
	}
	if response.StatusCode != http.StatusCreated {
		return "", "", fmt.Errorf("permission-denied api key create returned unexpected HTTP status %d", response.StatusCode)
	}
	var parsed cloudminiAPIKeyCreateResponse
	if err := json.Unmarshal(body, &parsed); err != nil || !parsed.Success {
		return "", "", fmt.Errorf("permission-denied api key create response envelope was invalid")
	}
	keyID := strings.TrimSpace(parsed.Data.APIKey.ID)
	plainKey := strings.TrimSpace(parsed.Data.PlainKey)
	if keyID == "" || plainKey == "" || !parsed.Data.APIKey.IsActive {
		return "", "", fmt.Errorf("permission-denied api key create response was incomplete")
	}
	return keyID, plainKey, nil
}

func callCloudminiProxyListWithTemporaryKey(ctx context.Context, config cloudminiErrorEvidenceConfig, plainKey string) (cloudminiErrorEvidenceResult, error) {
	result := cloudminiErrorEvidenceResult{
		Name:              cloudminiPermissionDeniedExampleName,
		ProviderCode:      "none",
		NormalizedCode:    provider.ErrorPermissionDenied,
		RetrySafety:       provider.RetrySafetyDoNotRetry,
		SideEffectCreated: "no",
		TemporaryKey:      true,
	}
	requestURL, err := resolveCloudminiErrorEvidenceURL(config.BaseURL, "/api/v3/proxies")
	if err != nil {
		return result, err
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return result, err
	}
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", "Bearer "+plainKey)
	request.Header.Set("User-Agent", "Billing-cloudmini-error-evidence/1")
	request.Header.Set("X-Request-ID", "billing-error-evidence-"+hashCloudminiErrorEvidence(cloudminiPermissionDeniedExampleName))

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return result, fmt.Errorf("permission-denied proxy list request failed before response")
	}
	defer response.Body.Close()
	body, err := io.ReadAll(io.LimitReader(response.Body, 1<<20))
	if err != nil {
		return result, fmt.Errorf("permission-denied proxy list response could not be read")
	}

	apiErr, envelopePresent := parseCloudminiErrorEvidenceBody(body)
	normalized := mapCloudminiErrorEvidenceCode(response.StatusCode, apiErr.Code)
	result.HTTPStatus = response.StatusCode
	result.ProviderCode = safeCloudminiProviderErrorCode(apiErr.Code)
	result.NormalizedCode = normalized
	result.RetrySafety = provider.DefaultRetrySafety(normalized)
	result.ErrorEnvelope = envelopePresent
	result.ErrorMessageField = strings.TrimSpace(apiErr.Message) != ""
	result.ErrorDetailsField = len(apiErr.Details) > 0 && string(apiErr.Details) != "null"
	if response.StatusCode != http.StatusForbidden {
		return result, fmt.Errorf("permission-denied proxy list returned unexpected HTTP status %d", response.StatusCode)
	}
	if normalized != provider.ErrorPermissionDenied {
		return result, fmt.Errorf("permission-denied proxy list mapped to unexpected provider code")
	}
	return result, nil
}

func revokeCloudminiAPIKey(ctx context.Context, config cloudminiErrorEvidenceConfig, keyID string) error {
	escapedID := url.PathEscape(strings.TrimSpace(keyID))
	if escapedID == "" {
		return fmt.Errorf("permission-denied api key cleanup id was empty")
	}
	request, err := newCloudminiManagementRequest(ctx, config, http.MethodDelete, "/api/v1/api-keys/"+escapedID, nil)
	if err != nil {
		return err
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return fmt.Errorf("permission-denied api key cleanup request failed before response")
	}
	defer response.Body.Close()
	_, _ = io.Copy(io.Discard, io.LimitReader(response.Body, 1<<20))
	if response.StatusCode != http.StatusNoContent {
		return fmt.Errorf("permission-denied api key cleanup returned unexpected HTTP status %d", response.StatusCode)
	}
	return nil
}

func newCloudminiManagementRequest(ctx context.Context, config cloudminiErrorEvidenceConfig, method string, path string, body io.Reader) (*http.Request, error) {
	requestURL, err := resolveCloudminiErrorEvidenceURL(config.BaseURL, path)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequestWithContext(ctx, method, requestURL, body)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Authorization", "Bearer "+config.APIToken)
	request.Header.Set("User-Agent", "Billing-cloudmini-error-evidence/1")
	request.Header.Set("X-Request-ID", "billing-error-evidence-"+hashCloudminiErrorEvidence(method+path))
	return request, nil
}
