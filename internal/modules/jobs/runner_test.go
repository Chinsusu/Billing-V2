package jobs

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

func TestRunnerRunOnceCompletesSucceededJob(t *testing.T) {
	now := fixedRunnerTime()
	store := &fakeJobStore{claimed: []Job{runnerJob(0, 5)}}
	runner := Runner{
		Store: store,
		Handler: HandlerFunc(func(ctx context.Context, job Job) (Completion, error) {
			return Completion{Status: StatusSucceeded}, nil
		}),
		WorkerID:  "worker_1",
		BatchSize: 2,
		LockFor:   time.Minute,
		Now:       func() time.Time { return now },
	}

	summary, err := runner.RunOnce(context.Background())
	if err != nil {
		t.Fatalf("expected run success, got %v", err)
	}
	if summary.Claimed != 1 || summary.Succeeded != 1 {
		t.Fatalf("unexpected summary: %#v", summary)
	}
	if store.claimRequest.WorkerID != "worker_1" || store.claimRequest.Limit != 2 {
		t.Fatalf("unexpected claim request: %#v", store.claimRequest)
	}
	if store.attempts[0].Result != AttemptResultSucceeded {
		t.Fatalf("expected succeeded attempt, got %q", store.attempts[0].Result)
	}
	if store.completions[0].Status != StatusSucceeded {
		t.Fatalf("expected succeeded completion, got %q", store.completions[0].Status)
	}
}

func TestRunnerRunOnceRetriesHandlerError(t *testing.T) {
	now := fixedRunnerTime()
	store := &fakeJobStore{claimed: []Job{runnerJob(1, 5)}}
	runner := Runner{
		Store: store,
		Handler: HandlerFunc(func(ctx context.Context, job Job) (Completion, error) {
			return Completion{}, errors.New("provider timeout")
		}),
		WorkerID: "worker_1",
		Backoff:  BackoffPolicy{Delays: []time.Duration{2 * time.Minute}},
		Now:      func() time.Time { return now },
	}

	summary, err := runner.RunOnce(context.Background())
	if err != nil {
		t.Fatalf("expected handled retry, got %v", err)
	}
	if summary.Retried != 1 {
		t.Fatalf("expected retry summary, got %#v", summary)
	}
	completion := store.completions[0]
	if completion.Status != StatusFailedRetryable {
		t.Fatalf("expected retryable completion, got %q", completion.Status)
	}
	if !completion.NextAttemptAt.Equal(now.Add(2 * time.Minute)) {
		t.Fatalf("unexpected retry time: %v", completion.NextAttemptAt)
	}
	if completion.LastErrorMessageRedacted != "job handler failed" {
		t.Fatalf("expected redacted error, got %q", completion.LastErrorMessageRedacted)
	}
}

func TestRunnerRunOnceMovesExhaustedErrorToManualReview(t *testing.T) {
	now := fixedRunnerTime()
	store := &fakeJobStore{claimed: []Job{runnerJob(4, 5)}}
	runner := Runner{
		Store: store,
		Handler: HandlerFunc(func(ctx context.Context, job Job) (Completion, error) {
			return Completion{}, errors.New("still failing")
		}),
		WorkerID: "worker_1",
		Now:      func() time.Time { return now },
	}

	summary, err := runner.RunOnce(context.Background())
	if err != nil {
		t.Fatalf("expected manual review completion, got %v", err)
	}
	if summary.ManualReview != 1 {
		t.Fatalf("expected manual review summary, got %#v", summary)
	}
	completion := store.completions[0]
	if completion.Status != StatusManualReview {
		t.Fatalf("expected manual review, got %q", completion.Status)
	}
	if completion.ManualReviewReason == "" {
		t.Fatal("expected manual review reason")
	}
}

func TestRunnerValidateRequiresHandler(t *testing.T) {
	runner := Runner{Store: &fakeJobStore{}, WorkerID: "worker_1"}

	if err := runner.Validate(); !errors.Is(err, ErrRunnerHandlerMissing) {
		t.Fatalf("expected handler error, got %v", err)
	}
}

type fakeJobStore struct {
	claimed      []Job
	claimRequest ClaimRequest
	attempts     []Attempt
	completions  []Completion
}

func (store *fakeJobStore) Claim(ctx context.Context, request ClaimRequest) ([]Job, error) {
	store.claimRequest = request
	return append([]Job(nil), store.claimed...), nil
}

func (store *fakeJobStore) RecordAttempt(ctx context.Context, attempt Attempt) error {
	store.attempts = append(store.attempts, attempt)
	return nil
}

func (store *fakeJobStore) Complete(ctx context.Context, jobID ID, completion Completion) error {
	store.completions = append(store.completions, completion)
	return nil
}

func runnerJob(attemptCount int, maxAttempts int) Job {
	return Job{
		ID:             "11111111-1111-1111-1111-111111111111",
		TenantID:       tenant.ID("22222222-2222-2222-2222-222222222222"),
		Type:           "provider.provision",
		ReferenceType:  "order_item",
		ReferenceID:    "33333333-3333-3333-3333-333333333333",
		Status:         StatusClaimed,
		IdempotencyKey: "idem_1",
		AttemptCount:   attemptCount,
		MaxAttempts:    maxAttempts,
		CorrelationID:  "44444444-4444-4444-4444-444444444444",
	}
}

func fixedRunnerTime() time.Time {
	return time.Date(2026, 4, 23, 9, 0, 0, 0, time.UTC)
}
