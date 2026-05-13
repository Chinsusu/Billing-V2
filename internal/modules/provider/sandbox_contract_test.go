package provider

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestRunSandboxContractPassesFakeAdapter(t *testing.T) {
	adapter := NewFakeAdapter(TypeHetzner)
	adapter.SetResult(OperationProvision, sandboxSuccess("external-1", "running"))
	adapter.SetResult(OperationGetStatus, sandboxSuccess("external-1", "running"))
	adapter.SetResult(OperationTerminate, sandboxSuccess("external-1", "terminated"))

	results := RunSandboxContract(context.Background(), DefaultSandboxContractOptions(adapter))
	if !SandboxContractPassed(results) {
		t.Fatalf("expected sandbox contract to pass, got %#v", results)
	}
	assertSandboxCalls(t, adapter.Calls, []OperationName{
		OperationCheckHealth,
		OperationCheckStock,
		OperationProvision,
		OperationGetStatus,
		OperationTerminate,
		OperationProvision,
		OperationProvision,
	})
}

func TestRunSandboxContractPassesProxyFakeAdapter(t *testing.T) {
	adapter := NewFakeAdapter(TypeProxyUpstream)
	adapter.SetResult(OperationProvision, sandboxSuccess("proxy-external-1", "running"))
	adapter.SetResult(OperationGetStatus, sandboxSuccess("proxy-external-1", "running"))
	adapter.SetResult(OperationTerminate, sandboxSuccess("proxy-external-1", "terminated"))

	results := RunSandboxContract(context.Background(), DefaultSandboxContractOptions(adapter))
	if !SandboxContractPassed(results) {
		t.Fatalf("expected proxy sandbox contract to pass, got %#v", results)
	}
}

func TestRunSandboxContractRequiresIdempotency(t *testing.T) {
	options := DefaultSandboxContractOptions(NewFakeAdapter(TypeHetzner))
	options.Operation.IdempotencyKey = ""

	results := RunSandboxContract(context.Background(), options)
	if len(results) != 1 || !errors.Is(results[0].Err, ErrIdempotencyKeyMissing) {
		t.Fatalf("expected setup idempotency failure, got %#v", results)
	}
}

func TestRunSandboxContractReportsProviderError(t *testing.T) {
	adapter := NewFakeAdapter(TypeHetzner)
	adapter.SetError(OperationProvision, NewError(ErrorTimeoutUnknown, "provider timed out"))

	results := RunSandboxContract(context.Background(), DefaultSandboxContractOptions(adapter))
	result := resultByName(results, SandboxCaseOrder)
	if result.Err == nil {
		t.Fatalf("expected order case error, got %#v", results)
	}
	if result.Status != OperationStatusUnknown || result.Retry != RetrySafetyUnsafeRetry {
		t.Fatalf("expected timeout mapping, got %#v", result)
	}
}

func TestRunSandboxContractMapsRequestKnownTimeoutToManualReview(t *testing.T) {
	adapter := NewFakeAdapter(TypeHetzner)
	adapter.SetError(OperationProvision, NewError(ErrorTimeoutRequestKnown, "provider request status unknown"))

	results := RunSandboxContract(context.Background(), DefaultSandboxContractOptions(adapter))
	result := resultByName(results, SandboxCaseOrder)
	if result.Err == nil {
		t.Fatalf("expected order case error, got %#v", results)
	}
	if result.Status != OperationStatusUnknown || result.Retry != RetrySafetyManualReviewRequired {
		t.Fatalf("expected manual review timeout mapping, got %#v", result)
	}
}

func TestRunSandboxContractRejectsUnredactedMessage(t *testing.T) {
	adapter := NewFakeAdapter(TypeHetzner)
	adapter.SetResult(OperationGetStatus, OperationResult{
		Status:               OperationStatusSuccess,
		ErrorMessageRedacted: "token leaked",
		ObservedAt:           time.Date(2026, 4, 25, 0, 0, 0, 0, time.UTC),
	})

	results := RunSandboxContract(context.Background(), DefaultSandboxContractOptions(adapter))
	result := resultByName(results, SandboxCaseStatus)
	if result.Err == nil || result.Err.Error() != "status_read redacted message contains sensitive text" {
		t.Fatalf("expected redaction failure, got %#v", result)
	}
}

func sandboxSuccess(resourceID ExternalResourceID, status string) OperationResult {
	return OperationResult{
		Status:             OperationStatusSuccess,
		ExternalResourceID: resourceID,
		ProviderStatus:     status,
		RetrySafety:        RetrySafetyDoNotRetry,
		ObservedAt:         time.Date(2026, 4, 25, 0, 0, 0, 0, time.UTC),
	}
}

func resultByName(results []SandboxContractResult, name string) SandboxContractResult {
	for _, result := range results {
		if result.Name == name {
			return result
		}
	}
	return SandboxContractResult{}
}

func assertSandboxCalls(t *testing.T, got []OperationName, want []OperationName) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("expected calls %#v, got %#v", want, got)
	}
	for index := range want {
		if got[index] != want[index] {
			t.Fatalf("expected calls %#v, got %#v", want, got)
		}
	}
}
