package order

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/jobs"
	"github.com/Chinsusu/Billing-V2/internal/modules/provider"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

func TestServiceLifecycleSchedulerSchedulesDueActionsWithIdempotentKeys(t *testing.T) {
	now := fixedServiceLifecycleSchedulerTime()
	termEnd := now.Add(-time.Minute)
	dueStore := &fakeServiceLifecycleDueStore{
		actions: []ServiceLifecycleDueAction{{
			ServiceID:             "11111111-1111-1111-1111-111111111111",
			TenantID:              tenant.ID("22222222-2222-2222-2222-222222222222"),
			ProviderSourceID:      "55555555-5555-5555-5555-555555555555",
			ProviderType:          provider.TypeCloudminiV3,
			ExternalResourceID:    "proxy-1",
			Action:                ServiceLifecycleActionExpire,
			FromStatus:            ServiceStatusActive,
			ToStatus:              ServiceStatusExpired,
			BillingStatus:         BillingStatusOverdue,
			ExpectedBillingStatus: BillingStatusPaid,
			Reason:                "service term expired",
			TermEnd:               termEnd,
		}},
	}
	queue := &fakeServiceLifecycleQueue{}
	scheduler := &ServiceLifecycleScheduler{
		Store:       dueStore,
		Queue:       queue,
		Now:         func() time.Time { return now },
		Limit:       3,
		GracePeriod: time.Hour,
	}

	summary, err := scheduler.ScheduleDue(context.Background(), ListDueServiceLifecycleActionsInput{})
	if err != nil {
		t.Fatalf("expected schedule success: %v", err)
	}
	if summary.Due != 1 || summary.Scheduled != 1 {
		t.Fatalf("unexpected schedule summary: %+v", summary)
	}
	if dueStore.input.Limit != 3 || dueStore.input.GracePeriod != time.Hour || !dueStore.input.Now.Equal(now) {
		t.Fatalf("unexpected due input: %+v", dueStore.input)
	}
	if len(queue.inputs) != 1 {
		t.Fatalf("expected one queued job, got %d", len(queue.inputs))
	}
	input := queue.inputs[0]
	if input.Type != ServiceLifecycleJobType || input.ReferenceType != ServiceLifecycleReferenceType ||
		input.ReferenceID != jobs.ReferenceID("11111111-1111-1111-1111-111111111111") {
		t.Fatalf("unexpected queued job input: %+v", input)
	}
	if input.SourceID != jobs.SourceID("55555555-5555-5555-5555-555555555555") {
		t.Fatalf("unexpected provider source id: %+v", input)
	}
	if !strings.Contains(input.IdempotencyKey, "service_lifecycle:22222222-2222-2222-2222-222222222222:11111111-1111-1111-1111-111111111111:expire:active:expired") {
		t.Fatalf("unexpected idempotency key: %s", input.IdempotencyKey)
	}

	payload, err := DecodeServiceLifecycleJobPayload(input.PayloadJSON)
	if err != nil {
		t.Fatalf("expected valid payload: %v", err)
	}
	if payload.Action != ServiceLifecycleActionExpire ||
		payload.ProviderSourceID != "55555555-5555-5555-5555-555555555555" ||
		payload.ProviderType != provider.TypeCloudminiV3 ||
		payload.ExternalResourceID != "proxy-1" ||
		payload.ExpectedBillingStatus != BillingStatusPaid ||
		!payload.TermEnd.Equal(termEnd) {
		t.Fatalf("unexpected payload: %+v", payload)
	}
}

func TestServiceLifecycleHandlerTransitionsWithExpectedTermEnd(t *testing.T) {
	termEnd := fixedServiceLifecycleSchedulerTime()
	payload := mustServiceLifecyclePayload(t, ServiceLifecycleJobPayload{
		ServiceID:             "11111111-1111-1111-1111-111111111111",
		TenantID:              tenant.ID("22222222-2222-2222-2222-222222222222"),
		Action:                ServiceLifecycleActionExpire,
		FromStatus:            ServiceStatusActive,
		ToStatus:              ServiceStatusExpired,
		BillingStatus:         BillingStatusOverdue,
		ExpectedBillingStatus: BillingStatusPaid,
		TermEnd:               termEnd,
	})
	transitioner := &fakeServiceLifecycleTransitioner{}
	handler := &ServiceLifecycleHandler{Transitioner: transitioner, Now: func() time.Time { return termEnd }}

	completion, err := handler.Handle(context.Background(), serviceLifecycleJob(payload))
	if err != nil {
		t.Fatalf("expected handler success: %v", err)
	}
	if completion.Status != jobs.StatusSucceeded {
		t.Fatalf("expected success completion, got %+v", completion)
	}
	if transitioner.input.Action != ServiceLifecycleActionExpire ||
		!transitioner.input.TermEnd.Equal(termEnd) ||
		!transitioner.input.ExpectedTermEnd.Equal(termEnd) ||
		transitioner.input.ExpectedBillingStatus != BillingStatusPaid {
		t.Fatalf("unexpected transition input: %+v", transitioner.input)
	}
}

func TestServiceLifecycleHandlerTreatsStaleTransitionAsNoop(t *testing.T) {
	termEnd := fixedServiceLifecycleSchedulerTime()
	payload := mustServiceLifecyclePayload(t, ServiceLifecycleJobPayload{
		ServiceID:             "11111111-1111-1111-1111-111111111111",
		TenantID:              tenant.ID("22222222-2222-2222-2222-222222222222"),
		Action:                ServiceLifecycleActionExpire,
		FromStatus:            ServiceStatusActive,
		ToStatus:              ServiceStatusExpired,
		BillingStatus:         BillingStatusOverdue,
		ExpectedBillingStatus: BillingStatusPaid,
		TermEnd:               termEnd,
	})
	handler := &ServiceLifecycleHandler{
		Transitioner: &fakeServiceLifecycleTransitioner{err: ErrServiceStatusConflict},
		Now:          func() time.Time { return termEnd },
	}

	completion, err := handler.Handle(context.Background(), serviceLifecycleJob(payload))
	if err != nil {
		t.Fatalf("expected conflict no-op, got %v", err)
	}
	if completion.Status != jobs.StatusSucceeded {
		t.Fatalf("expected no-op success, got %+v", completion)
	}
}

func TestServiceLifecycleHandlerRejectsInvalidPayload(t *testing.T) {
	handler := &ServiceLifecycleHandler{
		Transitioner: &fakeServiceLifecycleTransitioner{},
		Now:          fixedServiceLifecycleSchedulerTime,
	}
	job := serviceLifecycleJob(json.RawMessage(`{"service_id": ""}`))

	completion, err := handler.Handle(context.Background(), job)
	if err != nil {
		t.Fatalf("expected invalid payload to be handled: %v", err)
	}
	if completion.Status != jobs.StatusFailedTerminal || completion.RetrySafety != jobs.RetrySafetyDoNotRetry {
		t.Fatalf("expected terminal invalid payload completion, got %+v", completion)
	}
}

func TestServiceLifecycleRunnerMovesRepeatedFailureToManualReview(t *testing.T) {
	now := fixedServiceLifecycleSchedulerTime()
	payload := mustServiceLifecyclePayload(t, ServiceLifecycleJobPayload{
		ServiceID:             "11111111-1111-1111-1111-111111111111",
		TenantID:              tenant.ID("22222222-2222-2222-2222-222222222222"),
		Action:                ServiceLifecycleActionExpire,
		FromStatus:            ServiceStatusActive,
		ToStatus:              ServiceStatusExpired,
		BillingStatus:         BillingStatusOverdue,
		ExpectedBillingStatus: BillingStatusPaid,
		TermEnd:               now,
	})
	store := &fakeServiceLifecycleJobStore{
		claimed: []jobs.Job{serviceLifecycleJob(payload)},
	}
	store.claimed[0].AttemptCount = 4
	store.claimed[0].MaxAttempts = 5
	runner := NewServiceLifecycleRunner(store, &fakeServiceLifecycleTransitioner{err: errors.New("database unavailable")}, "worker-1")
	runner.Now = func() time.Time { return now }

	summary, err := runner.RunOnce(context.Background())
	if err != nil {
		t.Fatalf("expected runner to complete manual review: %v", err)
	}
	if summary.ManualReview != 1 {
		t.Fatalf("expected manual review summary, got %+v", summary)
	}
	if store.completions[0].Status != jobs.StatusManualReview {
		t.Fatalf("expected manual review completion, got %+v", store.completions[0])
	}
}

func mustServiceLifecyclePayload(t *testing.T, payload ServiceLifecycleJobPayload) json.RawMessage {
	t.Helper()
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}
	return body
}

func serviceLifecycleJob(payload json.RawMessage) jobs.Job {
	return jobs.Job{
		ID:             "33333333-3333-3333-3333-333333333333",
		TenantID:       tenant.ID("22222222-2222-2222-2222-222222222222"),
		Type:           ServiceLifecycleJobType,
		ReferenceType:  ServiceLifecycleReferenceType,
		ReferenceID:    jobs.ReferenceID("11111111-1111-1111-1111-111111111111"),
		PayloadJSON:    payload,
		Status:         jobs.StatusClaimed,
		IdempotencyKey: "service_lifecycle:22222222-2222-2222-2222-222222222222:11111111-1111-1111-1111-111111111111:expire:active:expired:1",
		AttemptCount:   1,
		MaxAttempts:    5,
		CorrelationID:  "11111111-1111-1111-1111-111111111111",
	}
}

type fakeServiceLifecycleDueStore struct {
	input   ListDueServiceLifecycleActionsInput
	actions []ServiceLifecycleDueAction
	err     error
}

func (store *fakeServiceLifecycleDueStore) ListDueServiceLifecycleActions(_ context.Context, input ListDueServiceLifecycleActionsInput) ([]ServiceLifecycleDueAction, error) {
	store.input = input
	if store.err != nil {
		return nil, store.err
	}
	return append([]ServiceLifecycleDueAction(nil), store.actions...), nil
}

type fakeServiceLifecycleQueue struct {
	inputs []jobs.CreateJobInput
	err    error
}

func (queue *fakeServiceLifecycleQueue) CreateJob(_ context.Context, input jobs.CreateJobInput) (jobs.Job, error) {
	queue.inputs = append(queue.inputs, input)
	if queue.err != nil {
		return jobs.Job{}, queue.err
	}
	return jobs.Job{
		TenantID:       input.TenantID,
		Type:           input.Type,
		ReferenceType:  input.ReferenceType,
		ReferenceID:    input.ReferenceID,
		IdempotencyKey: input.IdempotencyKey,
		CorrelationID:  input.CorrelationID,
		Status:         jobs.StatusQueued,
		MaxAttempts:    input.MaxAttempts,
	}, nil
}

type fakeServiceLifecycleTransitioner struct {
	input TransitionServiceLifecycleInput
	err   error
}

func (transitioner *fakeServiceLifecycleTransitioner) TransitionServiceLifecycle(_ context.Context, input TransitionServiceLifecycleInput) (ServiceInstance, error) {
	transitioner.input = input
	if transitioner.err != nil {
		return ServiceInstance{}, transitioner.err
	}
	return ServiceInstance{ID: input.ID, TenantID: input.TenantID, Status: input.ToStatus}, nil
}

type fakeServiceLifecycleJobStore struct {
	claimed      []jobs.Job
	claimRequest jobs.ClaimRequest
	attempts     []jobs.Attempt
	completions  []jobs.Completion
}

func (store *fakeServiceLifecycleJobStore) Claim(_ context.Context, request jobs.ClaimRequest) ([]jobs.Job, error) {
	store.claimRequest = request
	return append([]jobs.Job(nil), store.claimed...), nil
}

func (store *fakeServiceLifecycleJobStore) RecordAttempt(_ context.Context, attempt jobs.Attempt) error {
	store.attempts = append(store.attempts, attempt)
	return nil
}

func (store *fakeServiceLifecycleJobStore) Complete(_ context.Context, _ jobs.ID, completion jobs.Completion) error {
	store.completions = append(store.completions, completion)
	return nil
}

func fixedServiceLifecycleSchedulerTime() time.Time {
	return time.Date(2026, 5, 13, 10, 0, 0, 0, time.UTC)
}
