package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Chinsusu/Billing-V2/internal/modules/provider"
)

const cloudminiRateLimitExampleName = "rate_limited_fixture"

func validateCloudminiRateLimitFixturePath(path string) error {
	path = strings.TrimSpace(path)
	if path == "" {
		return fmt.Errorf("CLOUDMINI_ERROR_EVIDENCE_RATE_LIMIT_FIXTURE_PATH is required for rate-limit evidence")
	}
	if !strings.HasPrefix(path, "/api/v3/") {
		return fmt.Errorf("CLOUDMINI_ERROR_EVIDENCE_RATE_LIMIT_FIXTURE_PATH must start with /api/v3/")
	}
	if strings.ContainsAny(path, "?#") {
		return fmt.Errorf("CLOUDMINI_ERROR_EVIDENCE_RATE_LIMIT_FIXTURE_PATH must not contain query or fragment")
	}
	if !strings.Contains(path, "fixture") || !strings.Contains(path, "rate") {
		return fmt.Errorf("CLOUDMINI_ERROR_EVIDENCE_RATE_LIMIT_FIXTURE_PATH must be a rate-limit fixture path")
	}
	if path == "/api/v3/proxies" || strings.HasPrefix(path, "/api/v3/capacity/") {
		return fmt.Errorf("CLOUDMINI_ERROR_EVIDENCE_RATE_LIMIT_FIXTURE_PATH must be a side-effect-free fixture path")
	}
	return nil
}

func runCloudminiRateLimitEvidence(ctx context.Context, config cloudminiErrorEvidenceConfig) (cloudminiErrorEvidenceResult, error) {
	result := cloudminiErrorEvidenceResult{
		Name:                 cloudminiRateLimitExampleName,
		ProviderCode:         "RATE_LIMITED",
		NormalizedCode:       provider.ErrorRateLimited,
		RetrySafety:          provider.RetrySafetySafeRetry,
		SideEffectCreated:    "not_applicable",
		RateLimitFixture:     true,
		RateLimitMaxRequests: 1,
	}

	request, err := newCloudminiErrorEvidenceV3Request(ctx, config, http.MethodGet, config.RateLimitFixturePath, nil)
	if err != nil {
		return result, err
	}
	request.Header.Set("X-Cloudmini-Error-Fixture", "rate_limited")
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return result, fmt.Errorf("rate-limit fixture request failed before response")
	}
	defer response.Body.Close()
	body, err := io.ReadAll(io.LimitReader(response.Body, 1<<20))
	if err != nil {
		return result, fmt.Errorf("rate-limit fixture response could not be read")
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

	if response.StatusCode != http.StatusTooManyRequests {
		return result, fmt.Errorf("rate-limit fixture returned unexpected HTTP status %d", response.StatusCode)
	}
	if normalized != provider.ErrorRateLimited {
		return result, fmt.Errorf("rate-limit fixture mapped to unexpected provider code")
	}
	return result, nil
}
