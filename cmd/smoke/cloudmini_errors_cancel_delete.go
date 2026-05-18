package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Chinsusu/Billing-V2/internal/modules/provider"
)

const cloudminiCancelDeleteExampleName = "cancel_delete_rejected_fixture"

type cloudminiErrorEvidenceOperation struct {
	State        string `json:"state"`
	ErrorCode    string `json:"error_code"`
	ErrorMessage string `json:"error_message"`
}

func validateCloudminiCancelDeleteFixturePath(path string) error {
	return validateCloudminiErrorFixturePath("CLOUDMINI_ERROR_EVIDENCE_CANCEL_DELETE_FIXTURE_PATH", path, "fixture", "delete", "rejected")
}

func runCloudminiCancelDeleteEvidence(ctx context.Context, config cloudminiErrorEvidenceConfig) (cloudminiErrorEvidenceResult, error) {
	result := cloudminiErrorEvidenceResult{
		Name:                   cloudminiCancelDeleteExampleName,
		ProviderCode:           "DELETE_FAILED",
		NormalizedCode:         provider.ErrorPartialSuccess,
		RetrySafety:            provider.RetrySafetyManualReviewRequired,
		SideEffectCreated:      "not_applicable",
		CancelDeleteFixture:    true,
		CancelDeleteMaxRequest: 1,
	}

	request, err := newCloudminiErrorEvidenceV3Request(ctx, config, http.MethodGet, config.CancelDeleteFixturePath, nil)
	if err != nil {
		return result, err
	}
	request.Header.Set(cloudminiErrorFixtureHeader, "delete_rejected")
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return result, fmt.Errorf("cancel/delete rejected fixture request failed before response")
	}
	defer response.Body.Close()
	body, err := io.ReadAll(io.LimitReader(response.Body, 1<<20))
	if err != nil {
		return result, fmt.Errorf("cancel/delete rejected fixture response could not be read")
	}

	operation, err := parseCloudminiErrorEvidenceOperation(body)
	if err != nil {
		return result, err
	}
	normalized := mapCloudminiErrorEvidenceCode(response.StatusCode, operation.ErrorCode)
	result.HTTPStatus = response.StatusCode
	result.ProviderCode = safeCloudminiProviderErrorCode(operation.ErrorCode)
	result.NormalizedCode = normalized
	result.RetrySafety = provider.DefaultRetrySafety(normalized)
	result.ErrorEnvelope = false
	result.ErrorMessageField = strings.TrimSpace(operation.ErrorMessage) != ""
	result.ErrorDetailsField = false
	result.ProviderOperationState = strings.TrimSpace(operation.State)

	if response.StatusCode != http.StatusOK {
		return result, fmt.Errorf("cancel/delete rejected fixture returned unexpected HTTP status %d", response.StatusCode)
	}
	if operation.State != "failed" {
		return result, fmt.Errorf("cancel/delete rejected fixture returned unexpected operation state")
	}
	if normalized != provider.ErrorPartialSuccess {
		return result, fmt.Errorf("cancel/delete rejected fixture mapped to unexpected provider code")
	}
	return result, nil
}

func parseCloudminiErrorEvidenceOperation(body []byte) (cloudminiErrorEvidenceOperation, error) {
	var envelope cloudminiErrorEnvelope
	if err := json.Unmarshal(body, &envelope); err != nil || !envelope.Success || len(envelope.Data) == 0 {
		return cloudminiErrorEvidenceOperation{}, fmt.Errorf("cancel/delete rejected fixture response envelope is invalid")
	}
	var operation cloudminiErrorEvidenceOperation
	if err := json.Unmarshal(envelope.Data, &operation); err == nil && operation.State != "" {
		return operation, nil
	}
	var wrapped struct {
		Operation cloudminiErrorEvidenceOperation `json:"operation"`
	}
	if err := json.Unmarshal(envelope.Data, &wrapped); err == nil && wrapped.Operation.State != "" {
		return wrapped.Operation, nil
	}
	return cloudminiErrorEvidenceOperation{}, fmt.Errorf("cancel/delete rejected fixture operation data is invalid")
}
