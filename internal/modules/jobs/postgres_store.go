package jobs

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
	platformdb "github.com/Chinsusu/Billing-V2/internal/platform/db"
	"github.com/lib/pq"
)

var ErrStoreExecutorMissing = errors.New("jobs store executor missing")

type PostgresStore struct {
	executor platformdb.Executor
}

func NewPostgresStore(executor platformdb.Executor) *PostgresStore {
	return &PostgresStore{executor: executor}
}

const jobColumns = `job_id, tenant_id, job_type, reference_type, reference_id, source_id, payload_json, status, priority, idempotency_key, attempt_count, max_attempts, next_attempt_at, locked_by, locked_until, last_error_code, last_error_message_redacted, manual_review_reason, correlation_id, created_at, updated_at, finished_at`
const outboxColumns = `outbox_event_id, tenant_id, aggregate_type, aggregate_id, event_type, payload_json, status, dedupe_key, attempt_count, max_attempts, next_attempt_at, locked_by, locked_until, last_error_code, last_error_message_redacted, correlation_id, created_at, published_at`

func (store *PostgresStore) Claim(ctx context.Context, request ClaimRequest) ([]Job, error) {
	if err := store.ready(); err != nil {
		return nil, err
	}
	if err := request.Validate(); err != nil {
		return nil, err
	}
	now := requestTime(request.Now)
	query, args := claimJobsQuery(request, now)
	rows, err := store.executor.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("claim jobs: %w", err)
	}
	defer rows.Close()
	return scanJobs(rows)
}

func (store *PostgresStore) RecordAttempt(ctx context.Context, attempt Attempt) error {
	if err := store.ready(); err != nil {
		return err
	}
	if err := validateAttempt(attempt); err != nil {
		return err
	}
	_, err := store.executor.ExecContext(ctx, `
INSERT INTO job_attempts (job_id, worker_id, attempt_number, started_at, finished_at, result, error_code, error_message_redacted, duration_ms, correlation_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		attempt.JobID, attempt.WorkerID, attempt.AttemptNumber, attempt.StartedAt, nullableTime(attempt.FinishedAt), attempt.Result,
		nullableString(attempt.ErrorCode), nullableString(attempt.ErrorMessageRedacted), nullableDurationMillis(attempt.Duration), attempt.CorrelationID)
	if err != nil {
		return fmt.Errorf("record job attempt: %w", err)
	}
	return nil
}

func (store *PostgresStore) Complete(ctx context.Context, jobID ID, completion Completion) error {
	if err := store.ready(); err != nil {
		return err
	}
	if jobID == "" {
		return ErrJobIDMissing
	}
	if err := completion.Validate(); err != nil {
		return err
	}
	finishedAt := nullableCompletionTime(completion.Status, completion.FinishedAt)
	result, err := store.executor.ExecContext(ctx, `
UPDATE jobs
SET status = $2,
    attempt_count = attempt_count + 1,
    next_attempt_at = COALESCE($3, next_attempt_at),
    locked_by = NULL,
    locked_until = NULL,
    last_error_code = $4,
    last_error_message_redacted = $5,
    manual_review_reason = $6,
    finished_at = $7,
    updated_at = NOW()
WHERE job_id = $1`, jobID, completion.Status, nullableTime(completion.NextAttemptAt), nullableString(completion.LastErrorCode), nullableString(completion.LastErrorMessageRedacted), nullableString(completion.ManualReviewReason), finishedAt)
	if err != nil {
		return fmt.Errorf("complete job: %w", err)
	}
	return ensureChanged(result, ErrJobNotFound)
}

func (store *PostgresStore) ClaimOutbox(ctx context.Context, request OutboxClaimRequest) ([]OutboxEvent, error) {
	if err := store.ready(); err != nil {
		return nil, err
	}
	if err := request.Validate(); err != nil {
		return nil, err
	}
	now := requestTime(request.Now)
	rows, err := store.executor.QueryContext(ctx, `
WITH candidates AS (
    SELECT outbox_event_id
    FROM outbox_events
    WHERE status IN ('pending', 'failed_retryable')
      AND next_attempt_at <= $1
      AND (locked_until IS NULL OR locked_until <= $1)
    ORDER BY created_at
    LIMIT $2
    FOR UPDATE SKIP LOCKED
)
UPDATE outbox_events event
SET status = 'processing', locked_by = $3, locked_until = $4
FROM candidates
WHERE event.outbox_event_id = candidates.outbox_event_id
RETURNING `+outboxColumns, now, request.Limit, request.WorkerID, now.Add(request.LockFor))
	if err != nil {
		return nil, fmt.Errorf("claim outbox events: %w", err)
	}
	defer rows.Close()
	return scanOutboxEvents(rows)
}

func (store *PostgresStore) CompleteOutbox(ctx context.Context, eventID OutboxEventID, completion OutboxCompletion) error {
	if err := store.ready(); err != nil {
		return err
	}
	if eventID == "" {
		return ErrOutboxEventIDMissing
	}
	if err := completion.Validate(); err != nil {
		return err
	}
	publishedAt := nullableOutboxPublishedAt(completion.Status, completion.PublishedAt)
	result, err := store.executor.ExecContext(ctx, `
UPDATE outbox_events
SET status = $2,
    attempt_count = attempt_count + 1,
    next_attempt_at = COALESCE($3, next_attempt_at),
    locked_by = NULL,
    locked_until = NULL,
    last_error_code = $4,
    last_error_message_redacted = $5,
    published_at = $6
WHERE outbox_event_id = $1`, eventID, completion.Status, nullableTime(completion.NextAttemptAt), nullableString(completion.LastErrorCode), nullableString(completion.LastErrorMessageRedacted), publishedAt)
	if err != nil {
		return fmt.Errorf("complete outbox event: %w", err)
	}
	return ensureChanged(result, ErrOutboxEventNotFound)
}

func (store *PostgresStore) ready() error {
	if store == nil || store.executor == nil {
		return ErrStoreExecutorMissing
	}
	return nil
}

func claimJobsQuery(request ClaimRequest, now time.Time) (string, []interface{}) {
	types := jobTypeStrings(request.Types)
	query := `
WITH candidates AS (
    SELECT job_id
    FROM jobs
    WHERE status IN ('queued', 'failed_retryable')
      AND next_attempt_at <= $1
      AND (locked_until IS NULL OR locked_until <= $1)`
	args := []interface{}{now, request.Limit, request.WorkerID, now.Add(request.LockFor)}
	if len(types) > 0 {
		query += `
      AND job_type = ANY($5)`
		args = append(args, pq.Array(types))
	}
	query += `
    ORDER BY priority ASC, created_at ASC
    LIMIT $2
    FOR UPDATE SKIP LOCKED
)
UPDATE jobs job
SET status = 'claimed', locked_by = $3, locked_until = $4, updated_at = $1
FROM candidates
WHERE job.job_id = candidates.job_id
RETURNING ` + jobColumns
	return query, args
}

func scanJobs(rows *sql.Rows) ([]Job, error) {
	jobs := make([]Job, 0)
	for rows.Next() {
		job, err := scanJob(rows)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, job)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("read jobs: %w", err)
	}
	return jobs, nil
}

func scanOutboxEvents(rows *sql.Rows) ([]OutboxEvent, error) {
	events := make([]OutboxEvent, 0)
	for rows.Next() {
		event, err := scanOutboxEvent(rows)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("read outbox events: %w", err)
	}
	return events, nil
}

type rowScanner interface {
	Scan(dest ...interface{}) error
}

func scanJob(row rowScanner) (Job, error) {
	var job Job
	var id, jobType, referenceType, referenceID, status, correlationID string
	var tenantID, sourceID, lockedBy, lastErrorCode, lastErrorMessage, manualReviewReason sql.NullString
	var lockedUntil, finishedAt sql.NullTime
	var payload []byte
	if err := row.Scan(
		&id, &tenantID, &jobType, &referenceType, &referenceID, &sourceID, &payload, &status, &job.Priority,
		&job.IdempotencyKey, &job.AttemptCount, &job.MaxAttempts, &job.NextAttemptAt, &lockedBy, &lockedUntil,
		&lastErrorCode, &lastErrorMessage, &manualReviewReason, &correlationID, &job.CreatedAt, &job.UpdatedAt, &finishedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Job{}, ErrJobNotFound
		}
		return Job{}, fmt.Errorf("scan job: %w", err)
	}
	job.ID = ID(id)
	job.TenantID = tenantIDFromNull(tenantID)
	job.Type = Type(jobType)
	job.ReferenceType = ReferenceType(referenceType)
	job.ReferenceID = ReferenceID(referenceID)
	job.SourceID = SourceID(sourceID.String)
	job.PayloadJSON = append(job.PayloadJSON, payload...)
	job.Status = Status(status)
	job.LockedBy = WorkerID(lockedBy.String)
	job.LockedUntil = lockedUntil.Time
	job.LastErrorCode = lastErrorCode.String
	job.LastErrorMessageRedacted = lastErrorMessage.String
	job.ManualReviewReason = manualReviewReason.String
	job.CorrelationID = CorrelationID(correlationID)
	job.FinishedAt = finishedAt.Time
	return job, nil
}

func scanOutboxEvent(row rowScanner) (OutboxEvent, error) {
	var event OutboxEvent
	var id, aggregateID, status, correlationID string
	var tenantID, lockedBy, lastErrorCode, lastErrorMessage sql.NullString
	var lockedUntil, publishedAt sql.NullTime
	var payload []byte
	if err := row.Scan(
		&id, &tenantID, &event.AggregateType, &aggregateID, &event.EventType, &payload, &status,
		&event.DedupeKey, &event.AttemptCount, &event.MaxAttempts, &event.NextAttemptAt, &lockedBy,
		&lockedUntil, &lastErrorCode, &lastErrorMessage, &correlationID, &event.CreatedAt, &publishedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return OutboxEvent{}, ErrOutboxEventNotFound
		}
		return OutboxEvent{}, fmt.Errorf("scan outbox event: %w", err)
	}
	event.ID = OutboxEventID(id)
	event.TenantID = tenantIDFromNull(tenantID)
	event.AggregateID = aggregateID
	event.PayloadJSON = append(event.PayloadJSON, payload...)
	event.Status = OutboxStatus(status)
	event.LockedBy = WorkerID(lockedBy.String)
	event.LockedUntil = lockedUntil.Time
	event.LastErrorCode = lastErrorCode.String
	event.LastErrorMessageRedacted = lastErrorMessage.String
	event.CorrelationID = CorrelationID(correlationID)
	event.PublishedAt = publishedAt.Time
	return event, nil
}

func validateAttempt(attempt Attempt) error {
	if attempt.JobID == "" {
		return ErrJobIDMissing
	}
	if attempt.WorkerID == "" {
		return ErrWorkerIDMissing
	}
	if attempt.AttemptNumber <= 0 {
		return ErrAttemptNumberInvalid
	}
	if !attempt.Result.Valid() {
		return ErrAttemptResultInvalid
	}
	if attempt.CorrelationID == "" {
		return ErrCorrelationIDMissing
	}
	return nil
}

func ensureChanged(result sql.Result, notFound error) error {
	changed, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("read rows affected: %w", err)
	}
	if changed == 0 {
		return notFound
	}
	return nil
}

func requestTime(value time.Time) time.Time {
	if value.IsZero() {
		return time.Now().UTC()
	}
	return value
}

func nullableTime(value time.Time) sql.NullTime {
	if value.IsZero() {
		return sql.NullTime{}
	}
	return sql.NullTime{Time: value, Valid: true}
}

func nullableCompletionTime(status Status, value time.Time) sql.NullTime {
	if status.Terminal() || status == StatusManualReview {
		return nullableTime(requestTime(value))
	}
	return sql.NullTime{}
}

func nullableOutboxPublishedAt(status OutboxStatus, value time.Time) sql.NullTime {
	if status == OutboxStatusPublished {
		return nullableTime(requestTime(value))
	}
	return sql.NullTime{}
}

func nullableString(value string) sql.NullString {
	if value == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: value, Valid: true}
}

func nullableDurationMillis(value time.Duration) sql.NullInt64 {
	if value <= 0 {
		return sql.NullInt64{}
	}
	return sql.NullInt64{Int64: value.Milliseconds(), Valid: true}
}

func jobTypeStrings(types []Type) []string {
	values := make([]string, 0, len(types))
	for _, jobType := range types {
		if jobType == "" {
			continue
		}
		values = append(values, string(jobType))
	}
	return values
}

func tenantIDFromNull(value sql.NullString) tenant.ID {
	return tenant.ID(value.String)
}
