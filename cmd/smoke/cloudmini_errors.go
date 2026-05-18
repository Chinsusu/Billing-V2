package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/provider"
)

type cloudminiErrorEvidenceConfig struct {
	AppEnv                    string
	BaseURL                   string
	APIToken                  string
	IncludeCreate             bool
	IncludePermissionDenied   bool
	PermissionKeyManagementOK string
	PermissionKeyMaxCreate    string
}

type cloudminiErrorEvidenceExample struct {
	Name                 string
	Method               string
	Path                 string
	UseValidAuth         bool
	UseInvalidAuth       bool
	UseMalformedJSONBody bool
	UseIdempotencyKey    bool
	MutatingRoute        bool
	ExpectedStatuses     map[int]bool
}

type cloudminiErrorEvidenceResult struct {
	Name                   string
	HTTPStatus             int
	ProviderCode           string
	NormalizedCode         provider.ErrorCode
	RetrySafety            provider.RetrySafety
	ErrorEnvelope          bool
	ErrorMessageField      bool
	ErrorDetailsField      bool
	MutatingRoute          bool
	SideEffectCreated      string
	TemporaryKey           bool
	TemporaryKeyRevoked    bool
	ActiveKeyCountRestored bool
}

type cloudminiErrorEnvelope struct {
	Success bool            `json:"success"`
	Error   json.RawMessage `json:"error"`
}

type cloudminiErrorBody struct {
	Code    string          `json:"code"`
	Message string          `json:"message"`
	Details json.RawMessage `json:"details"`
}

func runCloudminiErrorEvidenceSmoke(timeout time.Duration) error {
	return runCloudminiErrorEvidenceSmokeWithWriter(timeout, os.Stdout)
}

func runCloudminiErrorEvidenceSmokeWithWriter(timeout time.Duration, out io.Writer) error {
	config, err := cloudminiErrorEvidenceConfigFromEnv()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	examples := cloudminiErrorEvidenceExamples(config.IncludeCreate)
	results := make([]cloudminiErrorEvidenceResult, 0, len(examples))
	for _, example := range examples {
		result, err := runCloudminiErrorEvidenceExample(ctx, config, example)
		results = append(results, result)
		if err != nil {
			printCloudminiErrorEvidenceSummary(out, config, results, "FAIL")
			return err
		}
	}
	if config.IncludePermissionDenied {
		result, err := runCloudminiPermissionDeniedEvidence(ctx, config)
		results = append(results, result)
		if err != nil {
			printCloudminiErrorEvidenceSummary(out, config, results, "FAIL")
			return err
		}
	}
	printCloudminiErrorEvidenceSummary(out, config, results, "PASS")
	return nil
}

func cloudminiErrorEvidenceConfigFromEnv() (cloudminiErrorEvidenceConfig, error) {
	appEnv := strings.ToLower(strings.TrimSpace(os.Getenv("APP_ENV")))
	if appEnv == "" {
		return cloudminiErrorEvidenceConfig{}, fmt.Errorf("APP_ENV is required")
	}
	switch appEnv {
	case "local", "dev", "staging", "sandbox":
	case "prod", "production":
		return cloudminiErrorEvidenceConfig{}, fmt.Errorf("refusing to run cloudmini error evidence with APP_ENV=%s", appEnv)
	default:
		return cloudminiErrorEvidenceConfig{}, fmt.Errorf("APP_ENV must be local, dev, staging, or sandbox")
	}
	if os.Getenv("BILLING_CLOUDMINI_ERROR_EVIDENCE_APPROVED") != "yes" {
		return cloudminiErrorEvidenceConfig{}, fmt.Errorf("BILLING_CLOUDMINI_ERROR_EVIDENCE_APPROVED=yes is required")
	}
	for _, key := range []string{
		"CLOUDMINI_SOURCE_ACCOUNT_OWNER",
		"CLOUDMINI_ENGINEERING_OWNER",
		"CLOUDMINI_OPS_OWNER",
		"CLOUDMINI_SECURITY_OWNER",
		"CLOUDMINI_CLEANUP_OWNER",
		"CLOUDMINI_REVIEWER_SIGNOFF",
		"CLOUDMINI_PILOT_STOP_CONDITION",
		"CLOUDMINI_PILOT_READONLY_EVIDENCE_REF",
	} {
		if err := requireCloudminiEvidenceFilled(key); err != nil {
			return cloudminiErrorEvidenceConfig{}, err
		}
	}
	config := cloudminiErrorEvidenceConfig{
		AppEnv:                    appEnv,
		BaseURL:                   strings.TrimSpace(os.Getenv("CLOUDMINI_V3_BASE_URL")),
		APIToken:                  strings.TrimSpace(os.Getenv("CLOUDMINI_V3_API_TOKEN")),
		IncludeCreate:             os.Getenv("CLOUDMINI_ERROR_EVIDENCE_ALLOW_INVALID_CREATE") == "yes",
		IncludePermissionDenied:   os.Getenv("CLOUDMINI_ERROR_EVIDENCE_ALLOW_PERMISSION_DENIED") == "yes",
		PermissionKeyManagementOK: strings.TrimSpace(os.Getenv("CLOUDMINI_ERROR_EVIDENCE_PERMISSION_KEY_MANAGEMENT_APPROVED")),
		PermissionKeyMaxCreate:    strings.TrimSpace(os.Getenv("CLOUDMINI_ERROR_EVIDENCE_PERMISSION_KEY_MAX_CREATE")),
	}
	if config.BaseURL == "" {
		return cloudminiErrorEvidenceConfig{}, fmt.Errorf("CLOUDMINI_V3_BASE_URL is required")
	}
	if config.APIToken == "" {
		return cloudminiErrorEvidenceConfig{}, fmt.Errorf("CLOUDMINI_V3_API_TOKEN is required")
	}
	if _, err := resolveCloudminiErrorEvidenceURL(config.BaseURL, "/api/v3/capabilities"); err != nil {
		return cloudminiErrorEvidenceConfig{}, err
	}
	if config.IncludeCreate {
		if strings.TrimSpace(os.Getenv("CLOUDMINI_ERROR_EVIDENCE_MUTATING_ROUTE_APPROVED")) != "yes" {
			return cloudminiErrorEvidenceConfig{}, fmt.Errorf("CLOUDMINI_ERROR_EVIDENCE_MUTATING_ROUTE_APPROVED=yes is required for malformed create validation evidence")
		}
		if strings.TrimSpace(os.Getenv("CLOUDMINI_ERROR_EVIDENCE_MAX_CREATE_ATTEMPTS")) != "1" {
			return cloudminiErrorEvidenceConfig{}, fmt.Errorf("CLOUDMINI_ERROR_EVIDENCE_MAX_CREATE_ATTEMPTS must be 1")
		}
	}
	if config.IncludePermissionDenied {
		if config.PermissionKeyManagementOK != "yes" {
			return cloudminiErrorEvidenceConfig{}, fmt.Errorf("CLOUDMINI_ERROR_EVIDENCE_PERMISSION_KEY_MANAGEMENT_APPROVED=yes is required for permission-denied evidence")
		}
		if config.PermissionKeyMaxCreate != "1" {
			return cloudminiErrorEvidenceConfig{}, fmt.Errorf("CLOUDMINI_ERROR_EVIDENCE_PERMISSION_KEY_MAX_CREATE must be 1")
		}
	}
	return config, nil
}

func cloudminiErrorEvidenceExamples(includeCreate bool) []cloudminiErrorEvidenceExample {
	examples := []cloudminiErrorEvidenceExample{
		{
			Name:             "auth_missing_capabilities",
			Method:           http.MethodGet,
			Path:             "/api/v3/capabilities",
			ExpectedStatuses: map[int]bool{http.StatusUnauthorized: true},
		},
		{
			Name:             "auth_invalid_capabilities",
			Method:           http.MethodGet,
			Path:             "/api/v3/capabilities",
			UseInvalidAuth:   true,
			ExpectedStatuses: map[int]bool{http.StatusUnauthorized: true, http.StatusForbidden: true},
		},
		{
			Name:             "not_found_proxy",
			Method:           http.MethodGet,
			Path:             "/api/v3/proxies/00000000-0000-4000-8000-000000000000",
			UseValidAuth:     true,
			ExpectedStatuses: map[int]bool{http.StatusNotFound: true},
		},
	}
	if includeCreate {
		examples = append(examples, cloudminiErrorEvidenceExample{
			Name:                 "validation_malformed_create",
			Method:               http.MethodPost,
			Path:                 "/api/v3/proxies",
			UseValidAuth:         true,
			UseMalformedJSONBody: true,
			UseIdempotencyKey:    true,
			MutatingRoute:        true,
			ExpectedStatuses:     map[int]bool{http.StatusBadRequest: true, http.StatusUnprocessableEntity: true},
		})
	}
	return examples
}

func runCloudminiErrorEvidenceExample(ctx context.Context, config cloudminiErrorEvidenceConfig, example cloudminiErrorEvidenceExample) (cloudminiErrorEvidenceResult, error) {
	requestURL, err := resolveCloudminiErrorEvidenceURL(config.BaseURL, example.Path)
	if err != nil {
		return cloudminiErrorEvidenceResult{}, err
	}
	var payload io.Reader
	if example.UseMalformedJSONBody {
		payload = bytes.NewBufferString("{")
	}
	request, err := http.NewRequestWithContext(ctx, example.Method, requestURL, payload)
	if err != nil {
		return cloudminiErrorEvidenceResult{}, err
	}
	request.Header.Set("Accept", "application/json")
	request.Header.Set("User-Agent", "Billing-cloudmini-error-evidence/1")
	request.Header.Set("X-Request-ID", "billing-error-evidence-"+hashCloudminiErrorEvidence(example.Name))
	if example.UseValidAuth {
		request.Header.Set("Authorization", "Bearer "+config.APIToken)
	}
	if example.UseInvalidAuth {
		request.Header.Set("Authorization", "Bearer billing-invalid-token")
	}
	if example.UseMalformedJSONBody {
		request.Header.Set("Content-Type", "application/json")
	}
	if example.UseIdempotencyKey {
		request.Header.Set("Idempotency-Key", "billing-error-evidence-"+hashCloudminiErrorEvidence(example.Name))
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return cloudminiErrorEvidenceResult{}, fmt.Errorf("%s request failed before response", example.Name)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(io.LimitReader(response.Body, 1<<20))
	if err != nil {
		return cloudminiErrorEvidenceResult{}, fmt.Errorf("%s response could not be read", example.Name)
	}
	result := parseCloudminiErrorEvidenceResult(example, response.StatusCode, body)
	if !example.ExpectedStatuses[response.StatusCode] {
		return result, fmt.Errorf("%s returned unexpected HTTP status %d", example.Name, response.StatusCode)
	}
	if response.StatusCode >= 200 && response.StatusCode < 300 {
		return result, fmt.Errorf("%s unexpectedly succeeded", example.Name)
	}
	return result, nil
}

func resolveCloudminiErrorEvidenceURL(baseURL string, path string) (string, error) {
	parsed, err := url.Parse(strings.TrimSpace(baseURL))
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return "", fmt.Errorf("CLOUDMINI_V3_BASE_URL is invalid")
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return "", fmt.Errorf("CLOUDMINI_V3_BASE_URL scheme is invalid")
	}
	resolved := *parsed
	resolved.Path = strings.TrimRight(parsed.Path, "/") + path
	resolved.RawQuery = ""
	return resolved.String(), nil
}

func parseCloudminiErrorEvidenceResult(example cloudminiErrorEvidenceExample, statusCode int, body []byte) cloudminiErrorEvidenceResult {
	apiErr, envelopePresent := parseCloudminiErrorEvidenceBody(body)
	normalized := mapCloudminiErrorEvidenceCode(statusCode, apiErr.Code)
	return cloudminiErrorEvidenceResult{
		Name:              example.Name,
		HTTPStatus:        statusCode,
		ProviderCode:      safeCloudminiProviderErrorCode(apiErr.Code),
		NormalizedCode:    normalized,
		RetrySafety:       provider.DefaultRetrySafety(normalized),
		ErrorEnvelope:     envelopePresent,
		ErrorMessageField: strings.TrimSpace(apiErr.Message) != "",
		ErrorDetailsField: len(apiErr.Details) > 0 && string(apiErr.Details) != "null",
		MutatingRoute:     example.MutatingRoute,
		SideEffectCreated: "not_applicable",
	}
}

func parseCloudminiErrorEvidenceBody(body []byte) (cloudminiErrorBody, bool) {
	var envelope cloudminiErrorEnvelope
	if err := json.Unmarshal(body, &envelope); err == nil && len(envelope.Error) > 0 {
		return parseCloudminiErrorEvidenceError(envelope.Error), true
	}
	var generic map[string]json.RawMessage
	if err := json.Unmarshal(body, &generic); err == nil && len(generic["error"]) > 0 {
		return parseCloudminiErrorEvidenceError(generic["error"]), true
	}
	return cloudminiErrorBody{}, false
}

func parseCloudminiErrorEvidenceError(raw json.RawMessage) cloudminiErrorBody {
	var apiErr cloudminiErrorBody
	if err := json.Unmarshal(raw, &apiErr); err == nil {
		return apiErr
	}
	var message string
	if err := json.Unmarshal(raw, &message); err == nil {
		return cloudminiErrorBody{Message: message}
	}
	return cloudminiErrorBody{}
}

func mapCloudminiErrorEvidenceCode(statusCode int, providerCode string) provider.ErrorCode {
	switch providerCode {
	case "CAPACITY_EXHAUSTED":
		return provider.ErrorOutOfStock
	case "IDEMPOTENCY_CONFLICT", "INVALID_STATE_TRANSITION", "OPERATION_NOT_FOUND", "PROXY_NOT_FOUND":
		return provider.ErrorStateDrift
	case "INVALID_ACTION":
		return provider.ErrorCapabilityNotSupported
	case "INVALID_INPUT", "RESERVATION_NOT_FOUND", "RESERVATION_EXPIRED", "RESERVATION_ALREADY_CONSUMED":
		return provider.ErrorConfigInvalid
	case "INTERNAL_ERROR":
		return provider.ErrorTemporary
	}
	switch statusCode {
	case http.StatusUnauthorized:
		return provider.ErrorAuthFailed
	case http.StatusForbidden:
		return provider.ErrorPermissionDenied
	case http.StatusTooManyRequests:
		return provider.ErrorRateLimited
	case http.StatusConflict:
		return provider.ErrorStateDrift
	case http.StatusBadRequest, http.StatusUnprocessableEntity:
		return provider.ErrorConfigInvalid
	case http.StatusNotFound:
		return provider.ErrorStateDrift
	default:
		if statusCode >= 500 {
			return provider.ErrorTemporary
		}
		return provider.ErrorResponseInvalid
	}
}
