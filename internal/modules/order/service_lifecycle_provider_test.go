package order

import (
	"context"
	"testing"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/catalog"
	"github.com/Chinsusu/Billing-V2/internal/modules/jobs"
	"github.com/Chinsusu/Billing-V2/internal/modules/provider"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

func TestServiceLifecycleHandlerTerminatesProviderBeforeLifecycleTransition(t *testing.T) {
	now := fixedServiceLifecycleSchedulerTime()
	adapter := provider.NewFakeAdapter(provider.TypeCloudminiV3)
	registry, err := provider.NewRegistry(adapter)
	if err != nil {
		t.Fatalf("create registry: %v", err)
	}
	transitioner := &fakeServiceLifecycleTransitioner{}
	handler := &ServiceLifecycleHandler{
		Transitioner:     transitioner,
		ProviderRegistry: registry,
		Now:              func() time.Time { return now },
	}
	payload := mustServiceLifecyclePayload(t, providerBackedTerminatePayload(now))

	completion, err := handler.Handle(context.Background(), serviceLifecycleJob(payload))
	if err != nil {
		t.Fatalf("expected handler success: %v", err)
	}
	if completion.Status != jobs.StatusSucceeded {
		t.Fatalf("expected success completion, got %+v", completion)
	}
	if !sawProviderOperation(adapter.Calls, provider.OperationTerminate) {
		t.Fatalf("expected provider terminate call, got %+v", adapter.Calls)
	}
	if transitioner.input.Action != ServiceLifecycleActionTerminate ||
		transitioner.input.ID != ServiceID("11111111-1111-1111-1111-111111111111") {
		t.Fatalf("expected lifecycle transition after provider cleanup, got %+v", transitioner.input)
	}
}

func TestServiceLifecycleHandlerBlocksTransitionWhenProviderTerminateUnknown(t *testing.T) {
	now := fixedServiceLifecycleSchedulerTime()
	adapter := provider.NewFakeAdapter(provider.TypeCloudminiV3)
	adapter.SetError(provider.OperationTerminate, provider.AdapterError{
		Code:            provider.ErrorTimeoutRequestKnown,
		MessageRedacted: "cloudmini v3 delete status is unknown",
	})
	registry, err := provider.NewRegistry(adapter)
	if err != nil {
		t.Fatalf("create registry: %v", err)
	}
	transitioner := &fakeServiceLifecycleTransitioner{}
	handler := &ServiceLifecycleHandler{
		Transitioner:     transitioner,
		ProviderRegistry: registry,
		Now:              func() time.Time { return now },
	}
	payload := mustServiceLifecyclePayload(t, providerBackedTerminatePayload(now))

	completion, err := handler.Handle(context.Background(), serviceLifecycleJob(payload))
	if err != nil {
		t.Fatalf("expected provider uncertainty to be handled: %v", err)
	}
	if completion.Status != jobs.StatusManualReview ||
		completion.RetrySafety != jobs.RetrySafetyManualReviewRequired ||
		completion.LastErrorCode != string(provider.ErrorTimeoutRequestKnown) {
		t.Fatalf("expected manual review completion, got %+v", completion)
	}
	if transitioner.input.ID != "" {
		t.Fatalf("lifecycle transition must not run after uncertain provider cleanup: %+v", transitioner.input)
	}
}

func providerBackedTerminatePayload(termEnd time.Time) ServiceLifecycleJobPayload {
	return ServiceLifecycleJobPayload{
		ServiceID:                "11111111-1111-1111-1111-111111111111",
		TenantID:                 tenant.ID("22222222-2222-2222-2222-222222222222"),
		ProviderSourceID:         catalog.ProviderSourceID("55555555-5555-5555-5555-555555555555"),
		ProviderType:             provider.TypeCloudminiV3,
		ExternalResourceID:       provider.ExternalResourceID("proxy-1"),
		Action:                   ServiceLifecycleActionTerminate,
		FromStatus:               ServiceStatusSuspended,
		ToStatus:                 ServiceStatusTerminated,
		BillingStatus:            BillingStatusOverdue,
		SuspensionReason:         SuspensionReasonExpiry,
		ExpectedBillingStatus:    BillingStatusOverdue,
		ExpectedSuspensionReason: SuspensionReasonExpiry,
		Reason:                   "service expired beyond grace period",
		TermEnd:                  termEnd,
	}
}

func sawProviderOperation(calls []provider.OperationName, operation provider.OperationName) bool {
	for _, call := range calls {
		if call == operation {
			return true
		}
	}
	return false
}
