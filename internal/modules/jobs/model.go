package jobs

import (
	"context"
	"errors"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

var (
	ErrJobIDMissing          = errors.New("job id missing")
	ErrJobTypeMissing        = errors.New("job type missing")
	ErrReferenceMissing      = errors.New("job reference missing")
	ErrIdempotencyKeyMissing = errors.New("idempotency key missing")
	ErrCorrelationIDMissing  = errors.New("correlation id missing")
	ErrStatusInvalid         = errors.New("job status invalid")
	ErrWorkerIDMissing       = errors.New("worker id missing")
	ErrClaimLimitInvalid     = errors.New("claim limit invalid")
	ErrLockDurationInvalid   = errors.New("lock duration invalid")
	ErrMaxAttemptsInvalid    = errors.New("max attempts invalid")
)

type ID string
type AttemptID string
type OutboxEventID string
type WorkerID string
type Type string
type ReferenceType string
type ReferenceID string
type SourceID string
type CorrelationID string

type Status string

const (
	StatusQueued          Status = "queued"
	StatusClaimed         Status = "claimed"
	StatusRunning         Status = "running"
	StatusSucceeded       Status = "succeeded"
	StatusFailedRetryable Status = "failed_retryable"
	StatusFailedTerminal  Status = "failed_terminal"
	StatusManualReview    Status = "manual_review"
	StatusCancelled       Status = "cancelled"
)

func (status Status) Valid() bool {
	switch status {
	case StatusQueued,
		StatusClaimed,
		StatusRunning,
		StatusSucceeded,
		StatusFailedRetryable,
		StatusFailedTerminal,
		StatusManualReview,
		StatusCancelled:
		return true
	default:
		return false
	}
}

func (status Status) Claimable() bool {
	return status == StatusQueued || status == StatusFailedRetryable
}

func (status Status) Terminal() bool {
	return status == StatusSucceeded || status == StatusFailedTerminal || status == StatusCancelled
}

type RetrySafety string

const (
	RetrySafetySafeRetry            RetrySafety = "safe_retry"
	RetrySafetyUnsafeRetry          RetrySafety = "unsafe_retry"
	RetrySafetyDoNotRetry           RetrySafety = "do_not_retry"
	RetrySafetyManualReviewRequired RetrySafety = "manual_review_required"
)

type OutboxStatus string

const (
	OutboxStatusPending         OutboxStatus = "pending"
	OutboxStatusProcessing      OutboxStatus = "processing"
	OutboxStatusPublished       OutboxStatus = "published"
	OutboxStatusFailedRetryable OutboxStatus = "failed_retryable"
	OutboxStatusFailedTerminal  OutboxStatus = "failed_terminal"
	OutboxStatusDiscarded       OutboxStatus = "discarded"
)

type OutboxEvent struct {
	ID                       OutboxEventID
	TenantID                 tenant.ID
	AggregateType            string
	AggregateID              string
	EventType                string
	Status                   OutboxStatus
	DedupeKey                string
	AttemptCount             int
	MaxAttempts              int
	NextAttemptAt            time.Time
	LockedBy                 WorkerID
	LockedUntil              time.Time
	LastErrorCode            string
	LastErrorMessageRedacted string
	CorrelationID            CorrelationID
	CreatedAt                time.Time
	PublishedAt              time.Time
}

type Job struct {
	ID                       ID
	TenantID                 tenant.ID
	Type                     Type
	ReferenceType            ReferenceType
	ReferenceID              ReferenceID
	SourceID                 SourceID
	Status                   Status
	Priority                 int
	IdempotencyKey           string
	AttemptCount             int
	MaxAttempts              int
	NextAttemptAt            time.Time
	LockedBy                 WorkerID
	LockedUntil              time.Time
	LastErrorCode            string
	LastErrorMessageRedacted string
	ManualReviewReason       string
	CorrelationID            CorrelationID
	CreatedAt                time.Time
	UpdatedAt                time.Time
	FinishedAt               time.Time
}

func (job Job) Validate() error {
	if job.ID == "" {
		return ErrJobIDMissing
	}
	if job.Type == "" {
		return ErrJobTypeMissing
	}
	if job.ReferenceType == "" || job.ReferenceID == "" {
		return ErrReferenceMissing
	}
	if job.IdempotencyKey == "" {
		return ErrIdempotencyKeyMissing
	}
	if job.CorrelationID == "" {
		return ErrCorrelationIDMissing
	}
	if !job.Status.Valid() {
		return ErrStatusInvalid
	}
	if job.MaxAttempts <= 0 {
		return ErrMaxAttemptsInvalid
	}
	return nil
}

func (job Job) ClaimableAt(now time.Time) bool {
	if !job.Status.Claimable() {
		return false
	}
	if !job.NextAttemptAt.IsZero() && job.NextAttemptAt.After(now) {
		return false
	}
	return job.LockedUntil.IsZero() || !job.LockedUntil.After(now)
}

func (job Job) AttemptsRemaining() bool {
	return job.AttemptCount < job.MaxAttempts
}

type AttemptResult string

const (
	AttemptResultSucceeded       AttemptResult = "succeeded"
	AttemptResultFailedRetryable AttemptResult = "failed_retryable"
	AttemptResultFailedTerminal  AttemptResult = "failed_terminal"
	AttemptResultManualReview    AttemptResult = "manual_review"
	AttemptResultCancelled       AttemptResult = "cancelled"
)

type Attempt struct {
	ID                   AttemptID
	JobID                ID
	WorkerID             WorkerID
	AttemptNumber        int
	StartedAt            time.Time
	FinishedAt           time.Time
	Result               AttemptResult
	ErrorCode            string
	ErrorMessageRedacted string
	Duration             time.Duration
	CorrelationID        CorrelationID
}

type ClaimRequest struct {
	WorkerID WorkerID
	Limit    int
	LockFor  time.Duration
	Now      time.Time
	Types    []Type
}

func (request ClaimRequest) Validate() error {
	if request.WorkerID == "" {
		return ErrWorkerIDMissing
	}
	if request.Limit <= 0 {
		return ErrClaimLimitInvalid
	}
	if request.LockFor <= 0 {
		return ErrLockDurationInvalid
	}
	return nil
}

type Completion struct {
	Status                   Status
	RetrySafety              RetrySafety
	NextAttemptAt            time.Time
	LastErrorCode            string
	LastErrorMessageRedacted string
	ManualReviewReason       string
	FinishedAt               time.Time
}

type Store interface {
	// Claim must use a row lock such as SELECT FOR UPDATE SKIP LOCKED or an equivalent atomic claim.
	Claim(ctx context.Context, request ClaimRequest) ([]Job, error)
	RecordAttempt(ctx context.Context, attempt Attempt) error
	Complete(ctx context.Context, jobID ID, completion Completion) error
}

type BackoffPolicy struct {
	Delays []time.Duration
}

func DefaultBackoffPolicy() BackoffPolicy {
	return BackoffPolicy{Delays: []time.Duration{
		10 * time.Second,
		time.Minute,
		5 * time.Minute,
		15 * time.Minute,
		time.Hour,
	}}
}

func (policy BackoffPolicy) DelayForAttempt(attemptNumber int) time.Duration {
	if attemptNumber <= 0 || len(policy.Delays) == 0 {
		return 0
	}
	index := attemptNumber - 1
	if index >= len(policy.Delays) {
		index = len(policy.Delays) - 1
	}
	return policy.Delays[index]
}
