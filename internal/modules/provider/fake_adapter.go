package provider

import (
	"context"
	"time"
)

// FakeAdapter is a test double that returns pre-configured results.
// It covers success, every failure category, and timeout scenarios.
// Never use FakeAdapter with real credentials.
type FakeAdapter struct {
	providerType string
	capability   CapabilityProfile

	HealthResult    HealthResult
	StockResult     StockResult
	ProvisionResult OperationResult
	StatusResult    OperationResult
	SuspendResult   OperationResult
	UnsuspendResult OperationResult
	TerminateResult OperationResult
	RenewResult     OperationResult
	ResetPwdResult  OperationResult
}

// NewFakeAdapter returns a FakeAdapter with all operations returning success
// by default so tests only override what they care about.
func NewFakeAdapter(providerType string) *FakeAdapter {
	now := time.Now().UTC()
	ok := OperationResult{
		Status:      StatusSuccess,
		RetrySafety: RetrySafe,
		ObservedAt:  now,
	}
	return &FakeAdapter{
		providerType: providerType,
		capability: CapabilityProfile{
			HealthCheck:      true,
			LiveStockCheck:   true,
			AutoProvision:    true,
			StatusSync:       true,
			Suspend:          true,
			Unsuspend:        true,
			Terminate:        true,
			Renew:            true,
			ResetPassword:    true,
			CredentialFetch:  true,
			CredentialRotate: true,
		},
		HealthResult:    HealthResult{Status: HealthHealthy},
		StockResult:     StockResult{Status: StockAvailable},
		ProvisionResult: ok,
		StatusResult:    ok,
		SuspendResult:   ok,
		UnsuspendResult: ok,
		TerminateResult: ok,
		RenewResult:     ok,
		ResetPwdResult:  ok,
	}
}

func (f *FakeAdapter) ProviderType() string { return f.providerType }

func (f *FakeAdapter) CapabilityProfile(_ context.Context, _ string) (CapabilityProfile, error) {
	return f.capability, nil
}

func (f *FakeAdapter) CheckHealth(_ context.Context, _ OperationContext) (HealthResult, error) {
	return f.HealthResult, nil
}

func (f *FakeAdapter) CheckStock(_ context.Context, _ OperationContext) (StockResult, error) {
	return f.StockResult, nil
}

func (f *FakeAdapter) Provision(_ context.Context, _ ProvisionInput) (OperationResult, error) {
	return f.ProvisionResult, nil
}

func (f *FakeAdapter) GetStatus(_ context.Context, _ StatusInput) (OperationResult, error) {
	return f.StatusResult, nil
}

func (f *FakeAdapter) Suspend(_ context.Context, _ SuspendInput) (OperationResult, error) {
	return f.SuspendResult, nil
}

func (f *FakeAdapter) Unsuspend(_ context.Context, _ SuspendInput) (OperationResult, error) {
	return f.UnsuspendResult, nil
}

func (f *FakeAdapter) Terminate(_ context.Context, _ TerminateInput) (OperationResult, error) {
	return f.TerminateResult, nil
}

func (f *FakeAdapter) Renew(_ context.Context, _ RenewInput) (OperationResult, error) {
	return f.RenewResult, nil
}

func (f *FakeAdapter) ResetPassword(_ context.Context, _ ResetPasswordInput) (OperationResult, error) {
	return f.ResetPwdResult, nil
}
