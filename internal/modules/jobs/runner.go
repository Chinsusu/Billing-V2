package jobs

import (
	"context"
	"errors"
	"time"
)

var (
	ErrRunnerStoreMissing   = errors.New("runner store missing")
	ErrRunnerHandlerMissing = errors.New("runner handler missing")
)

type Handler interface {
	Handle(ctx context.Context, job Job) (Completion, error)
}

type HandlerFunc func(ctx context.Context, job Job) (Completion, error)

func (fn HandlerFunc) Handle(ctx context.Context, job Job) (Completion, error) {
	return fn(ctx, job)
}

type Runner struct {
	Store     Store
	Handler   Handler
	WorkerID  WorkerID
	BatchSize int
	LockFor   time.Duration
	Types     []Type
	Backoff   BackoffPolicy
	Now       func() time.Time
}

type RunSummary struct {
	Claimed        int
	Succeeded      int
	Retried        int
	ManualReview   int
	TerminalFailed int
	Cancelled      int
}

func (runner Runner) RunOnce(ctx context.Context) (RunSummary, error) {
	if err := runner.Validate(); err != nil {
		return RunSummary{}, err
	}
	now := runner.now()
	jobs, err := runner.Store.Claim(ctx, ClaimRequest{
		WorkerID: runner.WorkerID,
		Limit:    runner.batchSize(),
		LockFor:  runner.lockFor(),
		Now:      now,
		Types:    append([]Type(nil), runner.Types...),
	})
	if err != nil {
		return RunSummary{}, err
	}

	summary := RunSummary{Claimed: len(jobs)}
	for _, job := range jobs {
		completion, err := runner.handleJob(ctx, job)
		if err != nil {
			return summary, err
		}
		summary.record(completion.Status)
	}
	return summary, nil
}

func (runner Runner) Validate() error {
	if runner.Store == nil {
		return ErrRunnerStoreMissing
	}
	if runner.Handler == nil {
		return ErrRunnerHandlerMissing
	}
	request := ClaimRequest{WorkerID: runner.WorkerID, Limit: runner.batchSize(), LockFor: runner.lockFor()}
	return request.Validate()
}

func (runner Runner) handleJob(ctx context.Context, job Job) (Completion, error) {
	startedAt := runner.now()
	completion, handlerErr := runner.Handler.Handle(ctx, job)
	completion = runner.normalizeCompletion(job, completion, handlerErr, startedAt)
	finishedAt := completionFinishedAt(completion, runner.now())
	if completion.FinishedAt.IsZero() {
		completion.FinishedAt = finishedAt
	}

	attempt := Attempt{
		JobID:                job.ID,
		WorkerID:             runner.WorkerID,
		AttemptNumber:        job.AttemptCount + 1,
		StartedAt:            startedAt,
		FinishedAt:           finishedAt,
		Result:               attemptResultFromStatus(completion.Status),
		ErrorCode:            completion.LastErrorCode,
		ErrorMessageRedacted: completion.LastErrorMessageRedacted,
		Duration:             finishedAt.Sub(startedAt),
		CorrelationID:        job.CorrelationID,
	}
	if err := runner.Store.RecordAttempt(ctx, attempt); err != nil {
		return completion, err
	}
	if err := runner.Store.Complete(ctx, job.ID, completion); err != nil {
		return completion, err
	}
	return completion, nil
}

func (runner Runner) normalizeCompletion(job Job, completion Completion, handlerErr error, startedAt time.Time) Completion {
	attemptNumber := job.AttemptCount + 1
	if handlerErr != nil {
		return runner.failureCompletion(job, attemptNumber, startedAt)
	}
	if completion.Status == "" {
		completion.Status = StatusSucceeded
	}
	if completion.Status == StatusFailedRetryable {
		if attemptNumber >= job.MaxAttempts {
			return runner.manualReviewCompletion(startedAt)
		}
		if completion.NextAttemptAt.IsZero() {
			completion.NextAttemptAt = startedAt.Add(runner.backoff().DelayForAttempt(attemptNumber))
		}
	}
	return completion
}

func (runner Runner) failureCompletion(job Job, attemptNumber int, startedAt time.Time) Completion {
	if attemptNumber >= job.MaxAttempts {
		return runner.manualReviewCompletion(startedAt)
	}
	return Completion{
		Status:                   StatusFailedRetryable,
		RetrySafety:              RetrySafetySafeRetry,
		NextAttemptAt:            startedAt.Add(runner.backoff().DelayForAttempt(attemptNumber)),
		LastErrorCode:            "job_handler_error",
		LastErrorMessageRedacted: "job handler failed",
	}
}

func (runner Runner) manualReviewCompletion(startedAt time.Time) Completion {
	return Completion{
		Status:                   StatusManualReview,
		RetrySafety:              RetrySafetyManualReviewRequired,
		LastErrorCode:            "job_attempts_exhausted",
		LastErrorMessageRedacted: "job attempts exhausted",
		ManualReviewReason:       "job attempts exhausted",
		FinishedAt:               startedAt,
	}
}

func (summary *RunSummary) record(status Status) {
	switch status {
	case StatusSucceeded:
		summary.Succeeded++
	case StatusFailedRetryable:
		summary.Retried++
	case StatusManualReview:
		summary.ManualReview++
	case StatusFailedTerminal:
		summary.TerminalFailed++
	case StatusCancelled:
		summary.Cancelled++
	}
}

func (runner Runner) now() time.Time {
	if runner.Now == nil {
		return time.Now().UTC()
	}
	return runner.Now()
}

func (runner Runner) batchSize() int {
	if runner.BatchSize <= 0 {
		return 10
	}
	return runner.BatchSize
}

func (runner Runner) lockFor() time.Duration {
	if runner.LockFor <= 0 {
		return time.Minute
	}
	return runner.LockFor
}

func (runner Runner) backoff() BackoffPolicy {
	if len(runner.Backoff.Delays) == 0 {
		return DefaultBackoffPolicy()
	}
	return runner.Backoff
}

func completionFinishedAt(completion Completion, now time.Time) time.Time {
	if !completion.FinishedAt.IsZero() {
		return completion.FinishedAt
	}
	if completion.Status == StatusFailedRetryable {
		return now
	}
	if completion.Status.Terminal() || completion.Status == StatusManualReview {
		return now
	}
	return now
}

func attemptResultFromStatus(status Status) AttemptResult {
	switch status {
	case StatusSucceeded:
		return AttemptResultSucceeded
	case StatusFailedRetryable:
		return AttemptResultFailedRetryable
	case StatusManualReview:
		return AttemptResultManualReview
	case StatusCancelled:
		return AttemptResultCancelled
	default:
		return AttemptResultFailedTerminal
	}
}
