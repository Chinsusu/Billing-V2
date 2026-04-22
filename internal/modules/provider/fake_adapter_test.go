package provider_test

import (
	"context"
	"testing"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/provider"
)

func opCtx() provider.OperationContext {
	return provider.OperationContext{
		OperationID:    "op-1",
		TenantID:       "tenant-1",
		SourceID:       "src-1",
		IdempotencyKey: "idem-1",
		CorrelationID:  "corr-1",
		AttemptNumber:  1,
		RequestTimeout: 10 * time.Second,
	}
}

func TestFakeAdapter_ProviderType(t *testing.T) {
	a := provider.NewFakeAdapter("fake")
	if a.ProviderType() != "fake" {
		t.Fatalf("expected 'fake', got %q", a.ProviderType())
	}
}

func TestFakeAdapter_CapabilityProfile(t *testing.T) {
	a := provider.NewFakeAdapter("fake")
	cap, err := a.CapabilityProfile(context.Background(), "src-1")
	if err != nil {
		t.Fatal(err)
	}
	if !cap.AutoProvision {
		t.Error("expected AutoProvision = true in default fake")
	}
}

func TestFakeAdapter_HealthCheck_Healthy(t *testing.T) {
	a := provider.NewFakeAdapter("fake")
	res, err := a.CheckHealth(context.Background(), opCtx())
	if err != nil {
		t.Fatal(err)
	}
	if res.Status != provider.HealthHealthy {
		t.Errorf("expected %s, got %s", provider.HealthHealthy, res.Status)
	}
}

func TestFakeAdapter_HealthCheck_Down(t *testing.T) {
	a := provider.NewFakeAdapter("fake")
	a.HealthResult = provider.HealthResult{Status: provider.HealthDown, Message: "provider unreachable"}
	res, _ := a.CheckHealth(context.Background(), opCtx())
	if res.Status != provider.HealthDown {
		t.Errorf("expected %s, got %s", provider.HealthDown, res.Status)
	}
}

func TestFakeAdapter_StockCheck_Available(t *testing.T) {
	a := provider.NewFakeAdapter("fake")
	res, err := a.CheckStock(context.Background(), opCtx())
	if err != nil {
		t.Fatal(err)
	}
	if res.Status != provider.StockAvailable {
		t.Errorf("expected %s, got %s", provider.StockAvailable, res.Status)
	}
}

func TestFakeAdapter_StockCheck_OutOfStock(t *testing.T) {
	a := provider.NewFakeAdapter("fake")
	a.StockResult = provider.StockResult{Status: provider.StockOutOfStock}
	res, _ := a.CheckStock(context.Background(), opCtx())
	if res.Status != provider.StockOutOfStock {
		t.Errorf("expected %s, got %s", provider.StockOutOfStock, res.Status)
	}
}

func TestFakeAdapter_Provision_Success(t *testing.T) {
	a := provider.NewFakeAdapter("fake")
	a.ProvisionResult = provider.OperationResult{
		Status:             provider.StatusSuccess,
		ExternalResourceID: "vm-123",
		ServiceIdentifier:  "10.0.0.1",
		RetrySafety:        provider.RetrySafe,
		ObservedAt:         time.Now().UTC(),
	}
	res, err := a.Provision(context.Background(), provider.ProvisionInput{OpCtx: opCtx()})
	if err != nil {
		t.Fatal(err)
	}
	if res.Status != provider.StatusSuccess {
		t.Errorf("expected success, got %s", res.Status)
	}
	if res.ExternalResourceID != "vm-123" {
		t.Errorf("unexpected resource id: %s", res.ExternalResourceID)
	}
}

func TestFakeAdapter_Provision_AuthFailed(t *testing.T) {
	a := provider.NewFakeAdapter("fake")
	a.ProvisionResult = provider.OperationResult{
		Status:       provider.StatusFailed,
		ErrorCode:    provider.ErrAuthFailed,
		ErrorMessage: "invalid api key",
		RetrySafety:  provider.RetryDoNot,
		ObservedAt:   time.Now().UTC(),
	}
	res, _ := a.Provision(context.Background(), provider.ProvisionInput{OpCtx: opCtx()})
	if res.ErrorCode != provider.ErrAuthFailed {
		t.Errorf("expected ErrAuthFailed, got %s", res.ErrorCode)
	}
	if res.RetrySafety != provider.RetryDoNot {
		t.Errorf("expected do_not_retry, got %s", res.RetrySafety)
	}
}

func TestFakeAdapter_Provision_TimeoutUnknown(t *testing.T) {
	a := provider.NewFakeAdapter("fake")
	a.ProvisionResult = provider.OperationResult{
		Status:       provider.StatusUnknown,
		ErrorCode:    provider.ErrTimeoutUnknown,
		ErrorMessage: "request may have reached provider",
		RetrySafety:  provider.RetryUnsafe,
		ObservedAt:   time.Now().UTC(),
	}
	res, _ := a.Provision(context.Background(), provider.ProvisionInput{OpCtx: opCtx()})
	if res.Status != provider.StatusUnknown {
		t.Errorf("expected unknown, got %s", res.Status)
	}
	if res.RetrySafety != provider.RetryUnsafe {
		t.Errorf("expected unsafe_retry, got %s", res.RetrySafety)
	}
}

func TestFakeAdapter_Provision_PartialSuccess(t *testing.T) {
	a := provider.NewFakeAdapter("fake")
	a.ProvisionResult = provider.OperationResult{
		Status:             provider.StatusPartialSuccess,
		ExternalResourceID: "vm-456",
		ErrorCode:          provider.ErrPartialSuccess,
		ErrorMessage:       "resource created but credential missing",
		RetrySafety:        provider.RetryManualReview,
		ObservedAt:         time.Now().UTC(),
	}
	res, _ := a.Provision(context.Background(), provider.ProvisionInput{OpCtx: opCtx()})
	if res.Status != provider.StatusPartialSuccess {
		t.Errorf("expected partial_success, got %s", res.Status)
	}
	if res.RetrySafety != provider.RetryManualReview {
		t.Errorf("expected manual_review_required, got %s", res.RetrySafety)
	}
}

func TestFakeAdapter_Provision_OutOfStock(t *testing.T) {
	a := provider.NewFakeAdapter("fake")
	a.ProvisionResult = provider.OperationResult{
		Status:       provider.StatusFailed,
		ErrorCode:    provider.ErrOutOfStock,
		ErrorMessage: "no capacity available",
		RetrySafety:  provider.RetryDoNot,
		ObservedAt:   time.Now().UTC(),
	}
	res, _ := a.Provision(context.Background(), provider.ProvisionInput{OpCtx: opCtx()})
	if res.ErrorCode != provider.ErrOutOfStock {
		t.Errorf("expected ErrOutOfStock, got %s", res.ErrorCode)
	}
	if res.RetrySafety != provider.RetryDoNot {
		t.Errorf("expected do_not_retry, got %s", res.RetrySafety)
	}
}

func TestFakeAdapter_Provision_RateLimited(t *testing.T) {
	a := provider.NewFakeAdapter("fake")
	a.ProvisionResult = provider.OperationResult{
		Status:       provider.StatusFailed,
		ErrorCode:    provider.ErrRateLimited,
		ErrorMessage: "too many requests",
		RetrySafety:  provider.RetrySafe,
		ObservedAt:   time.Now().UTC(),
	}
	res, _ := a.Provision(context.Background(), provider.ProvisionInput{OpCtx: opCtx()})
	if res.RetrySafety != provider.RetrySafe {
		t.Errorf("expected safe_retry, got %s", res.RetrySafety)
	}
}

func TestFakeAdapter_Terminate_Success(t *testing.T) {
	a := provider.NewFakeAdapter("fake")
	res, err := a.Terminate(context.Background(), provider.TerminateInput{OpCtx: opCtx(), ExternalResourceID: "vm-123"})
	if err != nil {
		t.Fatal(err)
	}
	if res.Status != provider.StatusSuccess {
		t.Errorf("expected success, got %s", res.Status)
	}
}

func TestFakeAdapter_ResetPassword_Success(t *testing.T) {
	a := provider.NewFakeAdapter("fake")
	// Simulate credential returned encrypted.
	a.ResetPwdResult = provider.OperationResult{
		Status:            provider.StatusSuccess,
		CredentialPayload: []byte("encrypted-credential-bytes"),
		RetrySafety:       provider.RetrySafe,
		ObservedAt:        time.Now().UTC(),
	}
	res, err := a.ResetPassword(context.Background(), provider.ResetPasswordInput{OpCtx: opCtx()})
	if err != nil {
		t.Fatal(err)
	}
	if len(res.CredentialPayload) == 0 {
		t.Error("expected non-empty credential payload")
	}
}

func TestFakeAdapter_CapabilityNotSupported(t *testing.T) {
	a := provider.NewFakeAdapter("fake")
	a.SuspendResult = provider.OperationResult{
		Status:       provider.StatusCapabilityNotSupported,
		ErrorCode:    provider.ErrCapabilityNotSupported,
		ErrorMessage: "provider does not support suspend",
		RetrySafety:  provider.RetryDoNot,
	}
	res, _ := a.Suspend(context.Background(), provider.SuspendInput{OpCtx: opCtx()})
	if res.Status != provider.StatusCapabilityNotSupported {
		t.Errorf("expected capability_not_supported, got %s", res.Status)
	}
}

// TestRetrySafetyFor verifies that every error code maps to a retry safety value.
func TestRetrySafetyFor(t *testing.T) {
	cases := []struct {
		code     provider.ErrorCode
		expected provider.RetrySafety
	}{
		{provider.ErrAuthFailed, provider.RetryDoNot},
		{provider.ErrPermissionDenied, provider.RetryDoNot},
		{provider.ErrOutOfStock, provider.RetryDoNot},
		{provider.ErrCapabilityNotSupported, provider.RetryDoNot},
		{provider.ErrRateLimited, provider.RetrySafe},
		{provider.ErrTemporaryError, provider.RetrySafe},
		{provider.ErrNetworkBeforeSend, provider.RetrySafe},
		{provider.ErrTimeoutUnknown, provider.RetryUnsafe},
		{provider.ErrTimeoutRequestKnown, provider.RetryManualReview},
		{provider.ErrPartialSuccess, provider.RetryManualReview},
		{provider.ErrCredentialMissing, provider.RetryManualReview},
	}
	for _, c := range cases {
		err := &provider.AdapterError{Code: c.code, Message: "test"}
		if err.RetrySafety() != c.expected {
			t.Errorf("code %s: expected %s, got %s", c.code, c.expected, err.RetrySafety())
		}
	}
}
