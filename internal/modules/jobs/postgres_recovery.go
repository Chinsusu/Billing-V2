package jobs

import (
	"context"
	"errors"
	"fmt"
	"time"
)

const retryJobSQL = `
UPDATE jobs
SET status = 'queued',
    next_attempt_at = $3,
    locked_by = NULL,
    locked_until = NULL,
    last_error_code = NULL,
    last_error_message_redacted = NULL,
    manual_review_reason = NULL,
    finished_at = NULL,
    updated_at = NOW()
WHERE job_id = $1
  AND tenant_id = $2
  AND status IN ('failed_retryable', 'manual_review')
RETURNING ` + jobColumns

const markManualReviewJobSQL = `
UPDATE jobs
SET status = 'manual_review',
    locked_by = NULL,
    locked_until = NULL,
    last_error_code = 'manual_review',
    last_error_message_redacted = $3,
    manual_review_reason = $3,
    finished_at = $4,
    updated_at = NOW()
WHERE job_id = $1
  AND tenant_id = $2
  AND status IN ('queued', 'failed_retryable', 'failed_terminal', 'manual_review')
RETURNING ` + jobColumns

const cancelJobSQL = `
UPDATE jobs
SET status = 'cancelled',
    locked_by = NULL,
    locked_until = NULL,
    last_error_code = 'job_cancelled',
    last_error_message_redacted = COALESCE($3, last_error_message_redacted),
    manual_review_reason = COALESCE($3, manual_review_reason),
    finished_at = COALESCE(finished_at, $4),
    updated_at = NOW()
WHERE job_id = $1
  AND tenant_id = $2
  AND status IN ('queued', 'failed_retryable', 'failed_terminal', 'manual_review', 'cancelled')
RETURNING ` + jobColumns

func (store *PostgresStore) RetryJob(ctx context.Context, input RetryJobInput) (Job, error) {
	if err := store.ready(); err != nil {
		return Job{}, err
	}
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return Job{}, err
	}
	lookup := Lookup{ID: input.ID, TenantID: input.TenantID}
	job, err := scanJob(store.executor.QueryRowContext(ctx, retryJobSQL, input.ID, input.TenantID, retryNextAttemptAt(input)))
	if errors.Is(err, ErrJobNotFound) {
		return Job{}, store.recoveryConflict(ctx, lookup)
	}
	if err != nil {
		return Job{}, fmt.Errorf("retry job: %w", err)
	}
	return job, nil
}

func (store *PostgresStore) MarkManualReview(ctx context.Context, input ManualReviewJobInput) (Job, error) {
	if err := store.ready(); err != nil {
		return Job{}, err
	}
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return Job{}, err
	}
	lookup := Lookup{ID: input.ID, TenantID: input.TenantID}
	job, err := scanJob(store.executor.QueryRowContext(ctx, markManualReviewJobSQL, input.ID, input.TenantID, input.Reason, requestTime(input.Now)))
	if errors.Is(err, ErrJobNotFound) {
		return Job{}, store.recoveryConflict(ctx, lookup)
	}
	if err != nil {
		return Job{}, fmt.Errorf("mark job manual review: %w", err)
	}
	return job, nil
}

func (store *PostgresStore) CancelJob(ctx context.Context, input CancelJobInput) (Job, error) {
	if err := store.ready(); err != nil {
		return Job{}, err
	}
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return Job{}, err
	}
	lookup := Lookup{ID: input.ID, TenantID: input.TenantID}
	job, err := scanJob(store.executor.QueryRowContext(ctx, cancelJobSQL, input.ID, input.TenantID, nullableString(input.Reason), requestTime(input.Now)))
	if errors.Is(err, ErrJobNotFound) {
		return Job{}, store.recoveryConflict(ctx, lookup)
	}
	if err != nil {
		return Job{}, fmt.Errorf("cancel job: %w", err)
	}
	return job, nil
}

func (store *PostgresStore) recoveryConflict(ctx context.Context, lookup Lookup) error {
	if _, err := store.GetJob(ctx, lookup); err != nil {
		if errors.Is(err, ErrJobNotFound) {
			return ErrJobNotFound
		}
		return err
	}
	return ErrJobStatusConflict
}

func retryNextAttemptAt(input RetryJobInput) time.Time {
	if input.NextAttemptAt.IsZero() {
		return requestTime(input.Now)
	}
	return input.NextAttemptAt
}
