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
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/provider"
)

const (
	cloudminiOutOfCapacityExampleName    = "out_of_capacity_reservation"
	cloudminiOutOfCapacityCleanupTimeout = 10 * time.Second
)

type cloudminiGroupInventoryListResponse struct {
	Success bool `json:"success"`
	Data    []struct {
		ID               string `json:"id"`
		Kind             string `json:"kind"`
		SellState        string `json:"sell_state"`
		AllocatableUnits int    `json:"allocatable_units"`
	} `json:"data"`
}

type cloudminiReservationCreateResponse struct {
	Success bool `json:"success"`
	Data    struct {
		ID string `json:"id"`
	} `json:"data"`
}

func runCloudminiOutOfCapacityEvidence(ctx context.Context, config cloudminiErrorEvidenceConfig) (cloudminiErrorEvidenceResult, error) {
	result := cloudminiOutOfCapacityBaselineResult(config)

	groupID, err := findCloudminiExhaustedGroup(ctx, config)
	if err != nil {
		return result, err
	}
	result.ExhaustedGroupSelected = true

	probeResult, reservationID, err := callCloudminiOutOfCapacityReservationProbe(ctx, config, groupID)
	result = probeResult
	if reservationID != "" {
		result.ReservationCreated = true
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), cloudminiOutOfCapacityCleanupTimeout)
		defer cleanupCancel()
		if cleanupErr := deleteCloudminiReservation(cleanupCtx, config, reservationID); cleanupErr == nil {
			result.ReservationCleanedUp = true
			result.SideEffectCreated = "cleaned_up"
		}
	}
	if err != nil {
		return result, err
	}
	if result.NormalizedCode != provider.ErrorOutOfStock {
		return result, fmt.Errorf("out-of-capacity reservation mapped to unexpected provider code")
	}
	return result, nil
}

func cloudminiOutOfCapacityBaselineResult(config cloudminiErrorEvidenceConfig) cloudminiErrorEvidenceResult {
	return cloudminiErrorEvidenceResult{
		Name:                   cloudminiOutOfCapacityExampleName,
		ProviderCode:           "CAPACITY_EXHAUSTED",
		NormalizedCode:         provider.ErrorOutOfStock,
		RetrySafety:            provider.RetrySafetyDoNotRetry,
		MutatingRoute:          true,
		SideEffectCreated:      "no",
		ReservationProbe:       true,
		ReservationMaxAttempts: 1,
		ReservationTTLSeconds:  config.OutOfCapacityTTLSeconds,
	}
}

func findCloudminiExhaustedGroup(ctx context.Context, config cloudminiErrorEvidenceConfig) (string, error) {
	request, err := newCloudminiErrorEvidenceV3Request(ctx, config, http.MethodGet, "/api/v3/inventory/groups", nil)
	if err != nil {
		return "", err
	}
	query := request.URL.Query()
	query.Set("kind", config.OutOfCapacityKind)
	request.URL.RawQuery = query.Encode()
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return "", fmt.Errorf("out-of-capacity inventory request failed before response")
	}
	defer response.Body.Close()
	body, err := io.ReadAll(io.LimitReader(response.Body, 1<<20))
	if err != nil {
		return "", fmt.Errorf("out-of-capacity inventory response could not be read")
	}
	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("out-of-capacity inventory returned unexpected HTTP status %d", response.StatusCode)
	}
	var parsed cloudminiGroupInventoryListResponse
	if err := json.Unmarshal(body, &parsed); err != nil || !parsed.Success {
		return "", fmt.Errorf("out-of-capacity inventory response envelope was invalid")
	}
	for _, item := range parsed.Data {
		if item.Kind == config.OutOfCapacityKind && item.AllocatableUnits <= 0 {
			id := strings.TrimSpace(item.ID)
			if id != "" {
				return id, nil
			}
		}
	}
	return "", fmt.Errorf("out-of-capacity inventory found no exhausted group for requested kind")
}

func callCloudminiOutOfCapacityReservationProbe(ctx context.Context, config cloudminiErrorEvidenceConfig, groupID string) (cloudminiErrorEvidenceResult, string, error) {
	result := cloudminiOutOfCapacityBaselineResult(config)
	result.ExhaustedGroupSelected = true
	payload := fmt.Sprintf(
		`{"kind":%q,"group_id":%q,"quantity":1,"ttl_seconds":%d,"external_ref":%q}`,
		config.OutOfCapacityKind,
		groupID,
		config.OutOfCapacityTTLSeconds,
		"billing-t257-out-of-capacity",
	)
	request, err := newCloudminiErrorEvidenceV3Request(ctx, config, http.MethodPost, "/api/v3/capacity/reservations", bytes.NewBufferString(payload))
	if err != nil {
		return result, "", err
	}
	request.Header.Set("Content-Type", "application/json")
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return result, "", fmt.Errorf("out-of-capacity reservation request failed before response")
	}
	defer response.Body.Close()
	body, err := io.ReadAll(io.LimitReader(response.Body, 1<<20))
	if err != nil {
		return result, "", fmt.Errorf("out-of-capacity reservation response could not be read")
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
	if response.StatusCode == http.StatusCreated {
		result.ReservationCreated = true
		result.SideEffectCreated = "yes"
		reservationID := cloudminiReservationIDFromBody(body)
		return result, reservationID, fmt.Errorf("out-of-capacity reservation unexpectedly succeeded")
	}
	if response.StatusCode != http.StatusConflict {
		return result, "", fmt.Errorf("out-of-capacity reservation returned unexpected HTTP status %d", response.StatusCode)
	}
	if normalized != provider.ErrorOutOfStock {
		return result, "", fmt.Errorf("out-of-capacity reservation mapped to unexpected provider code")
	}
	return result, "", nil
}

func cloudminiReservationIDFromBody(body []byte) string {
	var parsed cloudminiReservationCreateResponse
	if err := json.Unmarshal(body, &parsed); err != nil || !parsed.Success {
		return ""
	}
	return strings.TrimSpace(parsed.Data.ID)
}

func deleteCloudminiReservation(ctx context.Context, config cloudminiErrorEvidenceConfig, reservationID string) error {
	escapedID := url.PathEscape(strings.TrimSpace(reservationID))
	if escapedID == "" {
		return fmt.Errorf("out-of-capacity reservation cleanup id was empty")
	}
	request, err := newCloudminiErrorEvidenceV3Request(ctx, config, http.MethodDelete, "/api/v3/capacity/reservations/"+escapedID, nil)
	if err != nil {
		return err
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return fmt.Errorf("out-of-capacity reservation cleanup request failed before response")
	}
	defer response.Body.Close()
	_, _ = io.Copy(io.Discard, io.LimitReader(response.Body, 1<<20))
	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusNoContent {
		return fmt.Errorf("out-of-capacity reservation cleanup returned unexpected HTTP status %d", response.StatusCode)
	}
	return nil
}

func newCloudminiErrorEvidenceV3Request(ctx context.Context, config cloudminiErrorEvidenceConfig, method string, path string, body io.Reader) (*http.Request, error) {
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
