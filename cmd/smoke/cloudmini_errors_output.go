package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"strings"
)

func safeCloudminiProviderErrorCode(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "none"
	}
	for _, r := range value {
		if (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' {
			continue
		}
		return "redacted:" + hashCloudminiErrorEvidence(value)
	}
	return value
}

func printCloudminiErrorEvidenceSummary(out io.Writer, config cloudminiErrorEvidenceConfig, results []cloudminiErrorEvidenceResult, result string) {
	fmt.Fprintf(out, "cloudmini_error_evidence result=%s\n", result)
	fmt.Fprintf(out, "pilot_environment=%s\n", config.AppEnv)
	fmt.Fprintln(out, "approval_fields_present=yes")
	fmt.Fprintln(out, "owner_fields_present=yes")
	fmt.Fprintf(out, "example_count=%d\n", len(results))
	mutatingCalled := false
	for _, item := range results {
		if item.MutatingRoute {
			mutatingCalled = true
		}
	}
	fmt.Fprintf(out, "mutating_routes_called=%t\n", mutatingCalled)
	for index, item := range results {
		fmt.Fprintf(out, "example_%d_name=%s\n", index+1, item.Name)
		fmt.Fprintf(out, "example_%d_http_status=%d\n", index+1, item.HTTPStatus)
		fmt.Fprintf(out, "example_%d_provider_error_code=%s\n", index+1, item.ProviderCode)
		fmt.Fprintf(out, "example_%d_normalized_error_code=%s\n", index+1, item.NormalizedCode)
		fmt.Fprintf(out, "example_%d_retry_safety=%s\n", index+1, item.RetrySafety)
		fmt.Fprintf(out, "example_%d_error_envelope_present=%t\n", index+1, item.ErrorEnvelope)
		fmt.Fprintf(out, "example_%d_error_message_field_present=%t\n", index+1, item.ErrorMessageField)
		fmt.Fprintf(out, "example_%d_error_details_field_present=%t\n", index+1, item.ErrorDetailsField)
		fmt.Fprintf(out, "example_%d_side_effect_created=%s\n", index+1, item.SideEffectCreated)
		if item.TemporaryKey {
			fmt.Fprintf(out, "example_%d_temporary_api_key_created=true\n", index+1)
			fmt.Fprintf(out, "example_%d_temporary_api_key_revoked=%t\n", index+1, item.TemporaryKeyRevoked)
			fmt.Fprintf(out, "example_%d_active_key_count_restored=%t\n", index+1, item.ActiveKeyCountRestored)
		}
		if item.ReservationProbe {
			fmt.Fprintf(out, "example_%d_reservation_probe_attempted=true\n", index+1)
			fmt.Fprintf(out, "example_%d_exhausted_group_selected=%t\n", index+1, item.ExhaustedGroupSelected)
			fmt.Fprintf(out, "example_%d_reservation_created=%t\n", index+1, item.ReservationCreated)
			fmt.Fprintf(out, "example_%d_reservation_cleaned_up=%t\n", index+1, item.ReservationCleanedUp)
			fmt.Fprintf(out, "example_%d_reservation_max_attempts=%d\n", index+1, item.ReservationMaxAttempts)
			fmt.Fprintf(out, "example_%d_reservation_ttl_seconds=%d\n", index+1, item.ReservationTTLSeconds)
		}
	}
	fmt.Fprintln(out, "raw_response_body_printed=no")
	fmt.Fprintln(out, "sensitive_values_printed=no")
	fmt.Fprintln(out, "raw_provider_ids_printed=no")
	fmt.Fprintln(out, "provider_payloads_printed=no")
	fmt.Fprintf(out, "remaining_provider_controlled_examples=%s\n", remainingCloudminiProviderControlledExamples(results))
}

func hashCloudminiErrorEvidence(value string) string {
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:])[:12]
}

func remainingCloudminiProviderControlledExamples(results []cloudminiErrorEvidenceResult) string {
	permissionDeniedClosed := false
	outOfCapacityClosed := false
	for _, item := range results {
		if item.Name == cloudminiPermissionDeniedExampleName {
			permissionDeniedClosed = true
		}
		if item.Name == cloudminiOutOfCapacityExampleName {
			outOfCapacityClosed = true
		}
	}
	remaining := make([]string, 0, 5)
	if !permissionDeniedClosed {
		remaining = append(remaining, "permission_denied")
	}
	remaining = append(remaining, "rate_limited")
	if !outOfCapacityClosed {
		remaining = append(remaining, "out_of_capacity")
	}
	remaining = append(remaining, "provider_5xx", "cancel_rejected")
	return strings.Join(remaining, ",")
}
