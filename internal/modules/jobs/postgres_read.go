package jobs

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

const attemptColumns = `attempt.job_attempt_id, attempt.display_id, attempt.job_id, attempt.worker_id, attempt.attempt_number, attempt.started_at, attempt.finished_at, attempt.result, attempt.error_code, attempt.error_message_redacted, attempt.duration_ms, attempt.correlation_id`

func (store *PostgresStore) ListJobs(ctx context.Context, filter Filter) ([]Job, error) {
	if err := store.ready(); err != nil {
		return nil, err
	}
	query, args, err := buildListJobsQuery(filter)
	if err != nil {
		return nil, err
	}
	rows, err := store.executor.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list jobs: %w", err)
	}
	defer rows.Close()
	return scanJobs(rows)
}

func (store *PostgresStore) GetJob(ctx context.Context, lookup Lookup) (Job, error) {
	if err := store.ready(); err != nil {
		return Job{}, err
	}
	query, args, err := buildGetJobQuery(lookup)
	if err != nil {
		return Job{}, err
	}
	return scanJob(store.executor.QueryRowContext(ctx, query, args...))
}

func (store *PostgresStore) ListAttempts(ctx context.Context, filter AttemptFilter) ([]Attempt, error) {
	if err := store.ready(); err != nil {
		return nil, err
	}
	if err := store.ensureAttemptJobVisible(ctx, filter); err != nil {
		return nil, err
	}
	query, args, err := buildListAttemptsQuery(filter)
	if err != nil {
		return nil, err
	}
	rows, err := store.executor.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list job attempts: %w", err)
	}
	defer rows.Close()
	return scanAttempts(rows)
}

func (store *PostgresStore) SummarizeJobs(ctx context.Context, filter SummaryFilter) (JobSummary, error) {
	if err := store.ready(); err != nil {
		return JobSummary{}, err
	}
	filter = normalizeSummaryFilter(filter)
	query, args, err := buildJobSummaryQuery(filter)
	if err != nil {
		return JobSummary{}, err
	}
	summary, err := scanJobSummary(store.executor.QueryRowContext(ctx, query, args...))
	if err != nil {
		return JobSummary{}, err
	}
	summary.TenantID = filter.TenantID
	summary.Type = filter.Type
	return summary, nil
}

func (store *PostgresStore) ensureAttemptJobVisible(ctx context.Context, filter AttemptFilter) error {
	query, args, err := buildAttemptJobVisibleQuery(filter)
	if err != nil {
		return err
	}
	var exists int
	if err := store.executor.QueryRowContext(ctx, query, args...).Scan(&exists); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrJobNotFound
		}
		return fmt.Errorf("read job for attempts: %w", err)
	}
	return nil
}

func buildJobSummaryQuery(filter SummaryFilter) (string, []interface{}, error) {
	filter = normalizeSummaryFilter(filter)
	if err := validateSummaryFilter(filter); err != nil {
		return "", nil, err
	}
	return `WITH scoped AS (
    SELECT job_id, display_id, status, last_error_code, last_error_message_redacted, manual_review_reason, created_at, updated_at
    FROM jobs
    WHERE tenant_id = $1
      AND job_type = $2
),
counts AS (
    SELECT COUNT(*) AS total,
           COUNT(*) FILTER (WHERE status = 'queued') AS queued,
           COUNT(*) FILTER (WHERE status = 'claimed') AS claimed,
           COUNT(*) FILTER (WHERE status = 'running') AS running,
           COUNT(*) FILTER (WHERE status = 'succeeded') AS succeeded,
           COUNT(*) FILTER (WHERE status = 'failed_retryable') AS failed_retryable,
           COUNT(*) FILTER (WHERE status = 'failed_terminal') AS failed_terminal,
           COUNT(*) FILTER (WHERE status = 'manual_review') AS manual_review,
           COUNT(*) FILTER (WHERE status = 'cancelled') AS cancelled,
           MIN(created_at) FILTER (WHERE status = 'queued') AS oldest_queued_at,
           NOW() AS generated_at
    FROM scoped
),
latest_failure AS (
    SELECT job_id, display_id, status, last_error_code, last_error_message_redacted, manual_review_reason, created_at, updated_at
    FROM scoped
    WHERE status IN ('failed_retryable', 'failed_terminal', 'manual_review')
    ORDER BY updated_at DESC, created_at DESC
    LIMIT 1
)
SELECT counts.total, counts.queued, counts.claimed, counts.running, counts.succeeded,
       counts.failed_retryable, counts.failed_terminal, counts.manual_review, counts.cancelled,
       counts.oldest_queued_at, counts.generated_at,
       latest_failure.job_id, latest_failure.display_id, latest_failure.status,
       latest_failure.last_error_code, latest_failure.last_error_message_redacted,
       latest_failure.manual_review_reason, latest_failure.created_at, latest_failure.updated_at
FROM counts
LEFT JOIN latest_failure ON TRUE`, []interface{}{filter.TenantID, filter.Type}, nil
}

func buildListJobsQuery(filter Filter) (string, []interface{}, error) {
	filter = normalizeFilter(filter)
	if err := validateFilter(filter); err != nil {
		return "", nil, err
	}
	query := `SELECT ` + jobColumns + `
FROM jobs
WHERE tenant_id = $1`
	args := []interface{}{filter.TenantID}
	if filter.DisplayID > 0 {
		args = append(args, filter.DisplayID)
		query += fmt.Sprintf("\n  AND display_id = $%d", len(args))
	}
	if filter.Type != "" {
		args = append(args, filter.Type)
		query += fmt.Sprintf("\n  AND job_type = $%d", len(args))
	}
	if filter.Status != "" {
		args = append(args, filter.Status)
		query += fmt.Sprintf("\n  AND status = $%d", len(args))
	}
	if filter.ReferenceType != "" {
		args = append(args, filter.ReferenceType)
		query += fmt.Sprintf("\n  AND reference_type = $%d", len(args))
	}
	if filter.ReferenceID != "" {
		args = append(args, filter.ReferenceID)
		query += fmt.Sprintf("\n  AND reference_id = $%d", len(args))
	}
	if filter.SourceID != "" {
		args = append(args, filter.SourceID)
		query += fmt.Sprintf("\n  AND source_id = $%d", len(args))
	}
	args = append(args, filter.Limit)
	query += fmt.Sprintf("\nORDER BY created_at DESC\nLIMIT $%d", len(args))
	return query, args, nil
}

func buildGetJobQuery(lookup Lookup) (string, []interface{}, error) {
	lookup = normalizeLookup(lookup)
	if err := validateLookup(lookup); err != nil {
		return "", nil, err
	}
	return `SELECT ` + jobColumns + `
FROM jobs
WHERE job_id = $1
  AND tenant_id = $2`, []interface{}{lookup.ID, lookup.TenantID}, nil
}

func buildAttemptJobVisibleQuery(filter AttemptFilter) (string, []interface{}, error) {
	filter = normalizeAttemptFilter(filter)
	if err := validateAttemptFilter(filter); err != nil {
		return "", nil, err
	}
	return `SELECT 1
FROM jobs
WHERE job_id = $1
  AND tenant_id = $2`, []interface{}{filter.JobID, filter.TenantID}, nil
}

func buildListAttemptsQuery(filter AttemptFilter) (string, []interface{}, error) {
	filter = normalizeAttemptFilter(filter)
	if err := validateAttemptFilter(filter); err != nil {
		return "", nil, err
	}
	return `SELECT ` + attemptColumns + `
FROM job_attempts attempt
JOIN jobs job ON job.job_id = attempt.job_id
WHERE attempt.job_id = $1
  AND job.tenant_id = $2
ORDER BY attempt.attempt_number DESC
LIMIT $3`, []interface{}{filter.JobID, filter.TenantID, filter.Limit}, nil
}

func scanJobSummary(row rowScanner) (JobSummary, error) {
	var summary JobSummary
	var oldestQueuedAt, generatedAt, failureCreatedAt, failureUpdatedAt sql.NullTime
	var failureID, failureStatus, failureErrorCode, failureErrorMessage, failureManualReason sql.NullString
	var failureDisplayID sql.NullInt64
	if err := row.Scan(
		&summary.Total,
		&summary.Counts.Queued,
		&summary.Counts.Claimed,
		&summary.Counts.Running,
		&summary.Counts.Succeeded,
		&summary.Counts.FailedRetryable,
		&summary.Counts.FailedTerminal,
		&summary.Counts.ManualReview,
		&summary.Counts.Cancelled,
		&oldestQueuedAt,
		&generatedAt,
		&failureID,
		&failureDisplayID,
		&failureStatus,
		&failureErrorCode,
		&failureErrorMessage,
		&failureManualReason,
		&failureCreatedAt,
		&failureUpdatedAt,
	); err != nil {
		return JobSummary{}, fmt.Errorf("scan job summary: %w", err)
	}
	summary.AttentionCount = summary.Counts.AttentionCount()
	summary.OldestQueuedAt = oldestQueuedAt.Time
	summary.GeneratedAt = generatedAt.Time
	if failureID.Valid {
		summary.LatestFailure = &JobFailureContext{
			ID:                       ID(failureID.String),
			DisplayID:                failureDisplayID.Int64,
			Status:                   Status(failureStatus.String),
			LastErrorCode:            failureErrorCode.String,
			LastErrorMessageRedacted: failureErrorMessage.String,
			ManualReviewReason:       failureManualReason.String,
			CreatedAt:                failureCreatedAt.Time,
			UpdatedAt:                failureUpdatedAt.Time,
		}
	}
	return summary, nil
}

func scanAttempts(rows *sql.Rows) ([]Attempt, error) {
	attempts := make([]Attempt, 0)
	for rows.Next() {
		attempt, err := scanAttempt(rows)
		if err != nil {
			return nil, err
		}
		attempts = append(attempts, attempt)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("read job attempts: %w", err)
	}
	return attempts, nil
}

func scanAttempt(row rowScanner) (Attempt, error) {
	var attempt Attempt
	var id, jobID, workerID, result, correlationID string
	var finishedAt sql.NullTime
	var errorCode, errorMessage sql.NullString
	var durationMillis sql.NullInt64
	if err := row.Scan(
		&id, &attempt.DisplayID, &jobID, &workerID, &attempt.AttemptNumber, &attempt.StartedAt, &finishedAt, &result,
		&errorCode, &errorMessage, &durationMillis, &correlationID,
	); err != nil {
		return Attempt{}, fmt.Errorf("scan job attempt: %w", err)
	}
	attempt.ID = AttemptID(id)
	attempt.JobID = ID(jobID)
	attempt.WorkerID = WorkerID(workerID)
	attempt.FinishedAt = finishedAt.Time
	attempt.Result = AttemptResult(result)
	attempt.ErrorCode = errorCode.String
	attempt.ErrorMessageRedacted = errorMessage.String
	if durationMillis.Valid {
		attempt.Duration = time.Duration(durationMillis.Int64) * time.Millisecond
	}
	attempt.CorrelationID = CorrelationID(correlationID)
	return attempt, nil
}
