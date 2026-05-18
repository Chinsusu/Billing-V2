package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Chinsusu/Billing-V2/internal/modules/provider"
)

const cloudminiProvider5xxExampleName = "provider_5xx_fixture"

func validateCloudminiProvider5xxFixturePath(path string) error {
	return validateCloudminiErrorFixturePath("CLOUDMINI_ERROR_EVIDENCE_PROVIDER_5XX_FIXTURE_PATH", path, "fixture", "internal", "error")
}

func runCloudminiProvider5xxEvidence(ctx context.Context, config cloudminiErrorEvidenceConfig) (cloudminiErrorEvidenceResult, error) {
	result := cloudminiErrorEvidenceResult{
		Name:                   cloudminiProvider5xxExampleName,
		ProviderCode:           "INTERNAL_ERROR",
		NormalizedCode:         provider.ErrorTemporary,
		RetrySafety:            provider.RetrySafetySafeRetry,
		SideEffectCreated:      "not_applicable",
		Provider5xxFixture:     true,
		Provider5xxMaxRequests: 1,
	}

	request, err := newCloudminiErrorEvidenceV3Request(ctx, config, http.MethodGet, config.Provider5xxFixturePath, nil)
	if err != nil {
		return result, err
	}
	request.Header.Set(cloudminiErrorFixtureHeader, "internal_error")
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return result, fmt.Errorf("provider 5xx fixture request failed before response")
	}
	defer response.Body.Close()
	body, err := io.ReadAll(io.LimitReader(response.Body, 1<<20))
	if err != nil {
		return result, fmt.Errorf("provider 5xx fixture response could not be read")
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

	if response.StatusCode != http.StatusInternalServerError {
		return result, fmt.Errorf("provider 5xx fixture returned unexpected HTTP status %d", response.StatusCode)
	}
	if normalized != provider.ErrorTemporary {
		return result, fmt.Errorf("provider 5xx fixture mapped to unexpected provider code")
	}
	return result, nil
}
