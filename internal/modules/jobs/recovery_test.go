package jobs

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/audit"
	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

func TestServiceRetryJobAuditsRecovery(t *testing.T) {
	before := testRecoveryJob(StatusFailedRetryable)
	after := testRecoveryJob(StatusQueued)
	after.NextAttemptAt = time.Date(2026, 4, 24, 2, 0, 0, 0, time.UTC)
	store := &fakeRecoveryStore{before: before, after: after}
	audits := &fakeJobsAuditAppender{}
	service := NewServiceWithAudit(store, audits)

	result, err := service.RetryJob(context.Background(), RetryJobInput{
		ID:            " job_1 ",
		TenantID:      " tenant_1 ",
		ActorID:       " admin_1 ",
		NextAttemptAt: after.NextAttemptAt,
	})
	if err != nil {
		t.Fatalf("expected retry: %v", err)
	}
	if result.Status != StatusQueued {
		t.Fatalf("expected queued job, got %+v", result)
	}
	if store.retryInput.ID != ID("job_1") ||
		store.retryInput.TenantID != tenant.ID("tenant_1") ||
		store.retryInput.ActorID != identity.UserID("admin_1") {
		t.Fatalf("unexpected retry input: %+v", store.retryInput)
	}
	if audits.calls != 1 ||
		audits.input.Action != jobAuditActionRetry ||
		audits.input.TargetID != audit.TargetID("job_1") ||
		audits.input.ActorID != audit.ActorID("admin_1") {
		t.Fatalf("unexpected audit input: %+v", audits.input)
	}
}

func TestManualReviewInputRequiresReason(t *testing.T) {
	input := ManualReviewJobInput{ID: "job_1", TenantID: "tenant_1", ActorID: "admin_1"}

	if err := input.Normalize().Validate(); !errors.Is(err, ErrManualReviewReasonMissing) {
		t.Fatalf("expected reason error, got %v", err)
	}
}

func TestRetryInputRequiresActor(t *testing.T) {
	input := RetryJobInput{ID: "job_1", TenantID: "tenant_1"}

	if err := input.Validate(); !errors.Is(err, identity.ErrActorIDMissing) {
		t.Fatalf("expected actor error, got %v", err)
	}
}

func testRecoveryJob(status Status) Job {
	return Job{
		ID:            "job_1",
		DisplayID:     81001,
		TenantID:      "tenant_1",
		Type:          "provider.provision",
		ReferenceType: "order",
		ReferenceID:   "order_1",
		Status:        status,
		Priority:      50,
		AttemptCount:  2,
		MaxAttempts:   5,
		NextAttemptAt: time.Date(2026, 4, 24, 1, 0, 0, 0, time.UTC),
		CorrelationID: "correlation_1",
		CreatedAt:     time.Date(2026, 4, 24, 0, 0, 0, 0, time.UTC),
		UpdatedAt:     time.Date(2026, 4, 24, 0, 30, 0, 0, time.UTC),
		FinishedAt:    time.Date(2026, 4, 24, 1, 30, 0, 0, time.UTC),
	}
}

type fakeRecoveryStore struct {
	before     Job
	after      Job
	retryInput RetryJobInput
}

func (store *fakeRecoveryStore) ListJobs(ctx context.Context, filter Filter) ([]Job, error) {
	return []Job{store.before}, nil
}

func (store *fakeRecoveryStore) GetJob(ctx context.Context, lookup Lookup) (Job, error) {
	return store.before, nil
}

func (store *fakeRecoveryStore) ListAttempts(ctx context.Context, filter AttemptFilter) ([]Attempt, error) {
	return nil, nil
}

func (store *fakeRecoveryStore) RetryJob(ctx context.Context, input RetryJobInput) (Job, error) {
	store.retryInput = input
	return store.after, nil
}

func (store *fakeRecoveryStore) MarkManualReview(ctx context.Context, input ManualReviewJobInput) (Job, error) {
	return store.after, nil
}

func (store *fakeRecoveryStore) CancelJob(ctx context.Context, input CancelJobInput) (Job, error) {
	return store.after, nil
}

type fakeJobsAuditAppender struct {
	calls int
	input audit.AppendInput
}

func (appender *fakeJobsAuditAppender) Append(ctx context.Context, input audit.AppendInput) (audit.Log, error) {
	appender.calls++
	appender.input = input
	return audit.Log{ID: "audit_1"}, nil
}
