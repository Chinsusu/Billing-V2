package order

import (
	"context"
	"errors"

	"github.com/Chinsusu/Billing-V2/internal/modules/jobs"
	"github.com/Chinsusu/Billing-V2/internal/modules/provider"
)

func (handler *ServiceLifecycleHandler) terminateProviderBeforeLifecycle(
	ctx context.Context,
	job jobs.Job,
	payload ServiceLifecycleJobPayload,
) (jobs.Completion, bool, error) {
	if payload.Action != ServiceLifecycleActionTerminate || handler.ProviderRegistry == nil {
		return jobs.Completion{}, false, nil
	}
	result, err := handler.terminateProviderResource(ctx, job, payload)
	if err != nil {
		var adapterErr provider.AdapterError
		if errors.As(err, &adapterErr) && result.Status == "" {
			result = provider.ResultFromError(adapterErr, handler.now())
		}
		if result.Status == "" {
			result = provider.ResultFromError(provider.NewError(provider.ErrorTemporary, "provider terminate failed"), handler.now())
		}
	}
	if result.Status == provider.OperationStatusSuccess {
		return jobs.Completion{}, false, nil
	}
	return completionFromProviderResult(result, handler.now()), true, nil
}

func (handler *ServiceLifecycleHandler) terminateProviderResource(
	ctx context.Context,
	job jobs.Job,
	payload ServiceLifecycleJobPayload,
) (provider.OperationResult, error) {
	if payload.ProviderSourceID.Empty() || payload.ProviderType == "" || payload.ExternalResourceID == "" {
		return provider.ResultFromError(
			provider.NewError(provider.ErrorConfigInvalid, "service lifecycle provider cleanup metadata is missing"),
			handler.now(),
		), nil
	}
	operation := provider.OperationContext{
		OperationID:        provider.OperationID("terminate:" + string(payload.ServiceID) + ":" + string(payload.ProviderSourceID)),
		TenantID:           payload.TenantID,
		SourceID:           provider.SourceID(payload.ProviderSourceID),
		IdempotencyKey:     provider.IdempotencyKey(job.IdempotencyKey),
		CorrelationID:      provider.CorrelationID(job.CorrelationID),
		AttemptNumber:      job.AttemptCount + 1,
		ActorOrSystemID:    provider.ActorID("service.lifecycle"),
		CapabilitySnapshot: provider.CapabilityProfile{},
		ProviderSourceSnapshot: provider.SourceSnapshot{
			ProviderType: payload.ProviderType,
		},
	}
	if err := operation.Validate(); err != nil {
		return provider.ResultFromError(
			provider.NewError(provider.ErrorConfigInvalid, "service lifecycle provider cleanup context is invalid"),
			handler.now(),
		), nil
	}
	adapter, err := handler.ProviderRegistry.AdapterForOperation(operation)
	if err != nil {
		return provider.ResultFromError(
			provider.NewError(provider.ErrorAdapterNotFound, "provider adapter was not found for service cleanup"),
			handler.now(),
		), nil
	}
	return adapter.Terminate(ctx, operation, provider.ResourceRequest{
		ExternalResourceID: payload.ExternalResourceID,
		Reason:             payload.Reason,
	})
}
