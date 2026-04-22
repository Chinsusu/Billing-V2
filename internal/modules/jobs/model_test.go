package jobs

import (
	"errors"
	"testing"
	"time"
)

func validJob() Job {
	return Job{
		ID:             "job_1",
		Type:           "notification.send",
		ReferenceType:  "order",
		ReferenceID:    "order_1",
		Status:         StatusQueued,
		IdempotencyKey: "tenant_a:notification.send:order_1",
		MaxAttempts:    5,
		CorrelationID:  "11111111-1111-1111-1111-111111111111",
	}
}

func TestJobValidateAcceptsRequiredFields(t *testing.T) {
	if err := validJob().Validate(); err != nil {
		t.Fatalf("expected valid job, got %v", err)
	}
}

func TestJobValidateRequiresKnownStatus(t *testing.T) {
	job := validJob()
	job.Status = "lost"

	if err := job.Validate(); !errors.Is(err, ErrStatusInvalid) {
		t.Fatalf("expected status error, got %v", err)
	}
}

func TestJobValidateRequiresMaxAttempts(t *testing.T) {
	job := validJob()
	job.MaxAttempts = 0

	if err := job.Validate(); !errors.Is(err, ErrMaxAttemptsInvalid) {
		t.Fatalf("expected max attempts error, got %v", err)
	}
}

func TestJobClaimableAtRequiresClaimableStatus(t *testing.T) {
	now := time.Date(2026, 4, 22, 10, 0, 0, 0, time.UTC)
	job := validJob()
	job.Status = StatusRunning

	if job.ClaimableAt(now) {
		t.Fatal("running job must not be claimable")
	}
}

func TestJobClaimableAtRejectsFutureRetry(t *testing.T) {
	now := time.Date(2026, 4, 22, 10, 0, 0, 0, time.UTC)
	job := validJob()
	job.Status = StatusFailedRetryable
	job.NextAttemptAt = now.Add(time.Minute)

	if job.ClaimableAt(now) {
		t.Fatal("future retry must not be claimable")
	}
}

func TestJobClaimableAtAllowsExpiredLock(t *testing.T) {
	now := time.Date(2026, 4, 22, 10, 0, 0, 0, time.UTC)
	job := validJob()
	job.LockedBy = "worker-a"
	job.LockedUntil = now.Add(-time.Second)

	if !job.ClaimableAt(now) {
		t.Fatal("expired lock should be claimable")
	}
}

func TestJobClaimableAtRejectsActiveLock(t *testing.T) {
	now := time.Date(2026, 4, 22, 10, 0, 0, 0, time.UTC)
	job := validJob()
	job.LockedBy = "worker-a"
	job.LockedUntil = now.Add(time.Minute)

	if job.ClaimableAt(now) {
		t.Fatal("active lock must not be claimable")
	}
}

func TestClaimRequestValidate(t *testing.T) {
	request := ClaimRequest{WorkerID: "worker-a", Limit: 10, LockFor: time.Minute}

	if err := request.Validate(); err != nil {
		t.Fatalf("expected valid claim request, got %v", err)
	}
}

func TestClaimRequestValidateRejectsMissingWorker(t *testing.T) {
	request := ClaimRequest{Limit: 10, LockFor: time.Minute}

	if err := request.Validate(); !errors.Is(err, ErrWorkerIDMissing) {
		t.Fatalf("expected worker id error, got %v", err)
	}
}

func TestDefaultBackoffPolicyCapsAtLastDelay(t *testing.T) {
	policy := DefaultBackoffPolicy()

	if delay := policy.DelayForAttempt(99); delay != time.Hour {
		t.Fatalf("expected capped one hour delay, got %v", delay)
	}
}
