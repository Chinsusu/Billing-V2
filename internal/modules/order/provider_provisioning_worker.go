package order

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/jobs"
	"github.com/Chinsusu/Billing-V2/internal/modules/provider"
)

var (
	ErrProvisioningRegistryMissing = errors.New("provisioning registry missing")
	ErrProvisioningRecorderMissing = errors.New("provisioning recorder missing")
	ErrProvisioningPayloadInvalid  = errors.New("provisioning job payload invalid")
)

type ProvisioningResultRecorder interface {
	RecordProvisioningResult(ctx context.Context, input RecordProvisioningResultInput) (ProvisioningJob, error)
}

type ProviderProvisioningHandler struct {
	Registry *provider.Registry
	Recorder ProvisioningResultRecorder
	Now      func() time.Time
}

func NewProviderProvisioningHandler(registry *provider.Registry, recorder ProvisioningResultRecorder) *ProviderProvisioningHandler {
	return &ProviderProvisioningHandler{Registry: registry, Recorder: recorder}
}

func NewProviderProvisioningRunner(store jobs.Store, registry *provider.Registry, recorder ProvisioningResultRecorder, workerID jobs.WorkerID) jobs.Runner {
	return jobs.Runner{
		Store:     store,
		Handler:   NewProviderProvisioningHandler(registry, recorder),
		WorkerID:  workerID,
		BatchSize: 1,
		Types:     []jobs.Type{ProvisioningJobType},
	}
}

func (handler *ProviderProvisioningHandler) Handle(ctx context.Context, job jobs.Job) (jobs.Completion, error) {
	if err := handler.ready(); err != nil {
		return jobs.Completion{}, err
	}
	payload, err := decodeProvisioningPayload(job.PayloadJSON)
	if err != nil {
		return jobs.Completion{
			Status:                   jobs.StatusFailedTerminal,
			RetrySafety:              jobs.RetrySafetyDoNotRetry,
			LastErrorCode:            "provisioning_payload_invalid",
			LastErrorMessageRedacted: "provisioning job payload is invalid",
			FinishedAt:               handler.now(),
		}, nil
	}
	operation := providerOperationFromJob(job, payload)
	adapter, err := handler.Registry.AdapterForOperation(operation)
	if err != nil {
		result := provider.ResultFromError(provider.NewError(provider.ErrorAdapterNotFound, "provider adapter was not found"), handler.now())
		return handler.record(ctx, job, payload, operation, result)
	}
	result, err := adapter.Provision(ctx, operation, provider.ProvisionRequest{
		PlanKey: string(payload.TenantPlanID),
		InputSnapshot: map[string]string{
			"order_id":           string(payload.OrderID),
			"order_display_id":   strconv.FormatInt(payload.OrderDisplayID, 10),
			"provider_source_id": string(payload.ProviderSourceID),
		},
	})
	if err != nil {
		var adapterErr provider.AdapterError
		if errors.As(err, &adapterErr) {
			result = provider.ResultFromError(adapterErr, handler.now())
		} else if result.Status == "" {
			result = provider.ResultFromError(provider.NewError(provider.ErrorTemporary, "provider provisioning failed"), handler.now())
		}
	}
	return handler.record(ctx, job, payload, operation, result)
}

func (handler *ProviderProvisioningHandler) ready() error {
	if handler == nil || handler.Registry == nil {
		return ErrProvisioningRegistryMissing
	}
	if handler.Recorder == nil {
		return ErrProvisioningRecorderMissing
	}
	return nil
}

func (handler *ProviderProvisioningHandler) record(ctx context.Context, job jobs.Job, payload ProvisioningQueuePayload, operation provider.OperationContext, result provider.OperationResult) (jobs.Completion, error) {
	status := provisioningStatusFromProviderResult(result)
	_, err := handler.Recorder.RecordProvisioningResult(ctx, RecordProvisioningResultInput{
		OrderID:             payload.OrderID,
		TenantID:            payload.TenantID,
		ProviderSourceID:    payload.ProviderSourceID,
		ProviderOperationID: ProviderOperationID(operation.OperationID),
		Status:              status,
		IdempotencyKey:      IdempotencyKey(job.IdempotencyKey),
		AttemptNumber:       job.AttemptCount + 1,
		LastErrorCode:       string(result.ErrorCode),
		LastErrorMessage:    result.ErrorMessageRedacted,
	})
	if err != nil {
		return jobs.Completion{}, err
	}
	return completionFromProviderResult(result, handler.now()), nil
}

func decodeProvisioningPayload(body json.RawMessage) (ProvisioningQueuePayload, error) {
	var payload ProvisioningQueuePayload
	if err := json.Unmarshal(body, &payload); err != nil {
		return ProvisioningQueuePayload{}, ErrProvisioningPayloadInvalid
	}
	if payload.OrderID.Empty() || payload.TenantID.Empty() || payload.ProviderSourceID.Empty() || payload.ProviderType == "" {
		return ProvisioningQueuePayload{}, ErrProvisioningPayloadInvalid
	}
	return payload, nil
}

func providerOperationFromJob(job jobs.Job, payload ProvisioningQueuePayload) provider.OperationContext {
	attemptNumber := job.AttemptCount + 1
	return provider.OperationContext{
		OperationID:    provider.OperationID("provision:" + string(payload.OrderID) + ":" + string(payload.ProviderSourceID)),
		TenantID:       payload.TenantID,
		SourceID:       provider.SourceID(payload.ProviderSourceID),
		IdempotencyKey: provider.IdempotencyKey(job.IdempotencyKey),
		CorrelationID:  provider.CorrelationID(job.CorrelationID),
		AttemptNumber:  attemptNumber,
		ProviderSourceSnapshot: provider.SourceSnapshot{
			ProviderType: payload.ProviderType,
			PlanKey:      string(payload.TenantPlanID),
		},
	}
}

func provisioningStatusFromProviderResult(result provider.OperationResult) ProvisioningStatus {
	if result.Status == provider.OperationStatusSuccess {
		return ProvisioningStatusProvisioned
	}
	if result.RetrySafety == provider.RetrySafetyManualReviewRequired || result.RetrySafety == provider.RetrySafetyUnsafeRetry {
		return ProvisioningStatusManualReview
	}
	return ProvisioningStatusFailed
}

func completionFromProviderResult(result provider.OperationResult, now time.Time) jobs.Completion {
	if result.Status == provider.OperationStatusSuccess {
		return jobs.Completion{Status: jobs.StatusSucceeded, FinishedAt: now}
	}
	completion := jobs.Completion{
		LastErrorCode:            string(result.ErrorCode),
		LastErrorMessageRedacted: result.ErrorMessageRedacted,
		RetrySafety:              jobRetrySafety(result.RetrySafety),
	}
	switch result.RetrySafety {
	case provider.RetrySafetySafeRetry:
		completion.Status = jobs.StatusFailedRetryable
	case provider.RetrySafetyDoNotRetry, provider.RetrySafetyManualReviewRequired, provider.RetrySafetyUnsafeRetry:
		completion.Status = jobs.StatusManualReview
		completion.ManualReviewReason = result.ErrorMessageRedacted
		completion.FinishedAt = now
	default:
		completion.Status = jobs.StatusManualReview
		completion.ManualReviewReason = result.ErrorMessageRedacted
		completion.FinishedAt = now
	}
	return completion
}

func jobRetrySafety(safety provider.RetrySafety) jobs.RetrySafety {
	switch safety {
	case provider.RetrySafetySafeRetry:
		return jobs.RetrySafetySafeRetry
	case provider.RetrySafetyUnsafeRetry:
		return jobs.RetrySafetyUnsafeRetry
	case provider.RetrySafetyDoNotRetry:
		return jobs.RetrySafetyDoNotRetry
	default:
		return jobs.RetrySafetyManualReviewRequired
	}
}

func (handler *ProviderProvisioningHandler) now() time.Time {
	if handler.Now == nil {
		return time.Now().UTC()
	}
	return handler.Now()
}
