package order

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/catalog"
	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/jobs"
	"github.com/Chinsusu/Billing-V2/internal/modules/provider"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

func TestProviderProvisioningRunnerClaimsOneProvisioningJob(t *testing.T) {
	registry, err := provider.NewFakeRegistry(provider.TypeManual)
	if err != nil {
		t.Fatalf("expected registry: %v", err)
	}
	runner := NewProviderProvisioningRunner(&fakeWorkerJobStore{}, registry, &fakeProvisioningResultRecorder{}, "worker-1")

	if runner.BatchSize != 1 || runner.WorkerID != jobs.WorkerID("worker-1") || len(runner.Types) != 1 || runner.Types[0] != ProvisioningJobType {
		t.Fatalf("unexpected provisioning runner config: %+v", runner)
	}
}

func TestProviderProvisioningHandlerRecordsSuccess(t *testing.T) {
	now := fixedProvisioningWorkerTime()
	adapter := provider.NewFakeAdapter(provider.TypeManual)
	adapter.SetResult(provider.OperationProvision, provider.OperationResult{
		Status:             provider.OperationStatusSuccess,
		ExternalResourceID: "external-1",
		RetrySafety:        provider.RetrySafetyDoNotRetry,
		ObservedAt:         now,
	})
	registry, err := provider.NewRegistry(adapter)
	if err != nil {
		t.Fatalf("expected registry: %v", err)
	}
	recorder := &fakeProvisioningResultRecorder{}
	handler := &ProviderProvisioningHandler{Registry: registry, Recorder: recorder, Now: func() time.Time { return now }}

	completion, err := handler.Handle(context.Background(), provisioningWorkerJob())
	if err != nil {
		t.Fatalf("expected handler success: %v", err)
	}
	if completion.Status != jobs.StatusSucceeded {
		t.Fatalf("expected succeeded completion, got %+v", completion)
	}
	if len(adapter.Calls) != 1 || adapter.Calls[0] != provider.OperationProvision {
		t.Fatalf("expected one provider provision call, got %#v", adapter.Calls)
	}
	if recorder.input.Status != ProvisioningStatusProvisioned || recorder.input.AttemptNumber != 2 {
		t.Fatalf("unexpected recorded provisioning result: %+v", recorder.input)
	}
	if recorder.serviceInput.ExternalResourceID != provider.ExternalResourceID("external-1") ||
		recorder.serviceInput.TenantPlanID != catalog.TenantPlanID("44444444-4444-4444-4444-444444444444") ||
		!recorder.serviceInput.TermEnd.After(recorder.serviceInput.TermStart) {
		t.Fatalf("unexpected service create input: %+v", recorder.serviceInput)
	}
}

func TestProviderProvisioningHandlerStoresEncryptedCredential(t *testing.T) {
	now := fixedProvisioningWorkerTime()
	adapter := provider.NewFakeAdapter(provider.TypeManual)
	adapter.SetResult(provider.OperationProvision, provider.OperationResult{
		Status:             provider.OperationStatusSuccess,
		ExternalResourceID: "external-1",
		Credential: provider.CredentialEnvelope{
			Type:                 provider.CredentialTypeVPSRoot,
			EncryptedPayload:     "encrypted-fixture",
			EncryptionKeyVersion: "v1",
			MaskedHint:           "root / ****",
		},
		RetrySafety: provider.RetrySafetyDoNotRetry,
		ObservedAt:  now,
	})
	registry, err := provider.NewRegistry(adapter)
	if err != nil {
		t.Fatalf("expected registry: %v", err)
	}
	recorder := &fakeProvisioningResultRecorder{}
	handler := &ProviderProvisioningHandler{Registry: registry, Recorder: recorder, Now: func() time.Time { return now }}

	completion, err := handler.Handle(context.Background(), provisioningWorkerJob())
	if err != nil {
		t.Fatalf("expected handler success: %v", err)
	}
	if completion.Status != jobs.StatusSucceeded {
		t.Fatalf("expected succeeded completion, got %+v", completion)
	}
	if !recorder.credentialCalled {
		t.Fatal("expected encrypted credential to be stored")
	}
	if recorder.credentialInput.ServiceID != ServiceID("77777777-7777-7777-7777-777777777777") ||
		recorder.credentialInput.Type != CredentialTypeVPSRoot ||
		recorder.credentialInput.EncryptedPayload != "encrypted-fixture" ||
		recorder.credentialInput.MaskedHint != "root / ****" {
		t.Fatalf("unexpected credential input: %+v", recorder.credentialInput)
	}
}

func TestProviderProvisioningHandlerRecordsRetryableProviderError(t *testing.T) {
	now := fixedProvisioningWorkerTime()
	adapter := provider.NewFakeAdapter(provider.TypeManual)
	adapter.SetError(provider.OperationProvision, provider.NewError(provider.ErrorTemporary, "temporary provider outage"))
	registry, err := provider.NewRegistry(adapter)
	if err != nil {
		t.Fatalf("expected registry: %v", err)
	}
	recorder := &fakeProvisioningResultRecorder{}
	handler := &ProviderProvisioningHandler{Registry: registry, Recorder: recorder, Now: func() time.Time { return now }}

	completion, err := handler.Handle(context.Background(), provisioningWorkerJob())
	if err != nil {
		t.Fatalf("expected handled provider error: %v", err)
	}
	if completion.Status != jobs.StatusFailedRetryable || completion.RetrySafety != jobs.RetrySafetySafeRetry {
		t.Fatalf("expected retryable completion, got %+v", completion)
	}
	if recorder.input.Status != ProvisioningStatusFailed || recorder.input.LastErrorCode != string(provider.ErrorTemporary) {
		t.Fatalf("unexpected recorded retryable result: %+v", recorder.input)
	}
}

func TestProviderProvisioningHandlerRecordsPermanentProviderErrorForReview(t *testing.T) {
	now := fixedProvisioningWorkerTime()
	adapter := provider.NewFakeAdapter(provider.TypeManual)
	adapter.SetError(provider.OperationProvision, provider.NewError(provider.ErrorCapabilityNotSupported, "plan is not supported"))
	registry, err := provider.NewRegistry(adapter)
	if err != nil {
		t.Fatalf("expected registry: %v", err)
	}
	recorder := &fakeProvisioningResultRecorder{}
	handler := &ProviderProvisioningHandler{Registry: registry, Recorder: recorder, Now: func() time.Time { return now }}

	completion, err := handler.Handle(context.Background(), provisioningWorkerJob())
	if err != nil {
		t.Fatalf("expected handled provider error: %v", err)
	}
	if completion.Status != jobs.StatusManualReview || completion.RetrySafety != jobs.RetrySafetyDoNotRetry {
		t.Fatalf("expected manual review completion, got %+v", completion)
	}
	if recorder.input.Status != ProvisioningStatusFailed || recorder.input.LastErrorCode != string(provider.ErrorCapabilityNotSupported) {
		t.Fatalf("unexpected recorded permanent result: %+v", recorder.input)
	}
}

func TestProviderProvisioningHandlerMovesUnknownTimeoutToManualReview(t *testing.T) {
	now := fixedProvisioningWorkerTime()
	adapter := provider.NewFakeAdapter(provider.TypeManual)
	adapter.SetError(provider.OperationProvision, provider.NewError(provider.ErrorTimeoutRequestKnown, "provider request status unknown"))
	registry, err := provider.NewRegistry(adapter)
	if err != nil {
		t.Fatalf("expected registry: %v", err)
	}
	recorder := &fakeProvisioningResultRecorder{}
	handler := &ProviderProvisioningHandler{Registry: registry, Recorder: recorder, Now: func() time.Time { return now }}

	completion, err := handler.Handle(context.Background(), provisioningWorkerJob())
	if err != nil {
		t.Fatalf("expected handled timeout error: %v", err)
	}
	if completion.Status != jobs.StatusManualReview || completion.RetrySafety != jobs.RetrySafetyManualReviewRequired {
		t.Fatalf("expected manual review completion, got %+v", completion)
	}
	if recorder.input.Status != ProvisioningStatusManualReview || recorder.input.LastErrorCode != string(provider.ErrorTimeoutRequestKnown) {
		t.Fatalf("unexpected recorded timeout result: %+v", recorder.input)
	}
}

func TestProviderProvisioningHandlerRejectsInvalidPayload(t *testing.T) {
	registry, err := provider.NewFakeRegistry(provider.TypeManual)
	if err != nil {
		t.Fatalf("expected registry: %v", err)
	}
	handler := NewProviderProvisioningHandler(registry, &fakeProvisioningResultRecorder{})
	job := provisioningWorkerJob()
	job.PayloadJSON = []byte(`{"order_id": ""}`)

	completion, err := handler.Handle(context.Background(), job)
	if err != nil {
		t.Fatalf("expected invalid payload to be handled, got %v", err)
	}
	if completion.Status != jobs.StatusFailedTerminal || completion.RetrySafety != jobs.RetrySafetyDoNotRetry {
		t.Fatalf("expected terminal payload completion, got %+v", completion)
	}
}

func TestProviderProvisioningHandlerCreatesDeterministicLocalServiceID(t *testing.T) {
	now := fixedProvisioningWorkerTime()
	registry, err := provider.NewFakeRegistry(provider.TypeManual)
	if err != nil {
		t.Fatalf("expected registry: %v", err)
	}
	recorder := &fakeProvisioningResultRecorder{}
	handler := &ProviderProvisioningHandler{Registry: registry, Recorder: recorder, Now: func() time.Time { return now }}

	completion, err := handler.Handle(context.Background(), provisioningWorkerJob())
	if err != nil {
		t.Fatalf("expected handler success: %v", err)
	}
	if completion.Status != jobs.StatusSucceeded {
		t.Fatalf("expected succeeded completion, got %+v", completion)
	}
	if recorder.serviceInput.ExternalResourceID != provider.ExternalResourceID("local-11111111-1111-1111-1111-111111111111") {
		t.Fatalf("unexpected local resource id: %+v", recorder.serviceInput)
	}
}

func TestProviderProvisioningHandlerRequiresRecorder(t *testing.T) {
	registry, err := provider.NewFakeRegistry(provider.TypeManual)
	if err != nil {
		t.Fatalf("expected registry: %v", err)
	}
	handler := NewProviderProvisioningHandler(registry, nil)

	_, err = handler.Handle(context.Background(), provisioningWorkerJob())
	if !errors.Is(err, ErrProvisioningRecorderMissing) {
		t.Fatalf("expected recorder error, got %v", err)
	}
}

func provisioningWorkerJob() jobs.Job {
	payload, err := provisioningQueuePayloadJSON(Order{
		ID:            "11111111-1111-1111-1111-111111111111",
		DisplayID:     30001,
		TenantID:      tenant.ID("22222222-2222-2222-2222-222222222222"),
		BuyerUserID:   identity.UserID("33333333-3333-3333-3333-333333333333"),
		TenantPlanID:  catalog.TenantPlanID("44444444-4444-4444-4444-444444444444"),
		Currency:      "USD",
		TotalMinor:    2500,
		OrderStatus:   OrderStatusPaid,
		BillingStatus: BillingStatusPaid,
	}, catalog.ProviderSourceID("55555555-5555-5555-5555-555555555555"), provider.TypeManual)
	if err != nil {
		panic(err)
	}
	return jobs.Job{
		ID:             "66666666-6666-6666-6666-666666666666",
		TenantID:       tenant.ID("22222222-2222-2222-2222-222222222222"),
		Type:           ProvisioningJobType,
		ReferenceType:  ProvisioningReferenceType,
		ReferenceID:    jobs.ReferenceID("11111111-1111-1111-1111-111111111111"),
		SourceID:       jobs.SourceID("55555555-5555-5555-5555-555555555555"),
		PayloadJSON:    payload,
		Status:         jobs.StatusClaimed,
		IdempotencyKey: "provisioning:22222222-2222-2222-2222-222222222222:11111111-1111-1111-1111-111111111111:55555555-5555-5555-5555-555555555555",
		AttemptCount:   1,
		MaxAttempts:    5,
		CorrelationID:  jobs.CorrelationID("11111111-1111-1111-1111-111111111111"),
	}
}

type fakeProvisioningResultRecorder struct {
	input            RecordProvisioningResultInput
	serviceInput     CreateServiceInstanceInput
	credentialInput  CreateServiceCredentialInput
	credentialCalled bool
	err              error
}

func (recorder *fakeProvisioningResultRecorder) RecordProvisioningResult(_ context.Context, input RecordProvisioningResultInput) (ProvisioningJob, error) {
	recorder.input = input
	if recorder.err != nil {
		return ProvisioningJob{}, recorder.err
	}
	return ProvisioningJob{OrderID: input.OrderID, Status: input.Status}, nil
}

func (recorder *fakeProvisioningResultRecorder) CreateServiceInstance(_ context.Context, input CreateServiceInstanceInput) (ServiceInstance, error) {
	recorder.serviceInput = input
	if recorder.err != nil {
		return ServiceInstance{}, recorder.err
	}
	return ServiceInstance{
		ID:            "77777777-7777-7777-7777-777777777777",
		TenantID:      input.TenantID,
		OrderID:       input.OrderID,
		Status:        input.Status,
		BillingStatus: input.BillingStatus,
	}, nil
}

func (recorder *fakeProvisioningResultRecorder) CreateServiceCredential(_ context.Context, input CreateServiceCredentialInput) (ServiceCredential, error) {
	recorder.credentialInput = input
	recorder.credentialCalled = true
	if recorder.err != nil {
		return ServiceCredential{}, recorder.err
	}
	return ServiceCredential{ServiceID: input.ServiceID, Type: input.Type, Status: input.Status}, nil
}

type fakeWorkerJobStore struct{}

func (store *fakeWorkerJobStore) Claim(ctx context.Context, request jobs.ClaimRequest) ([]jobs.Job, error) {
	return nil, nil
}

func (store *fakeWorkerJobStore) RecordAttempt(ctx context.Context, attempt jobs.Attempt) error {
	return nil
}

func (store *fakeWorkerJobStore) Complete(ctx context.Context, jobID jobs.ID, completion jobs.Completion) error {
	return nil
}

func fixedProvisioningWorkerTime() time.Time {
	return time.Date(2026, 4, 23, 12, 0, 0, 0, time.UTC)
}
