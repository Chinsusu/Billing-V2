package jobs

import (
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

type jobResponse struct {
	ID                       ID            `json:"id"`
	DisplayID                int64         `json:"display_id"`
	TenantID                 tenant.ID     `json:"tenant_id"`
	Type                     Type          `json:"job_type"`
	ReferenceType            ReferenceType `json:"reference_type"`
	ReferenceID              ReferenceID   `json:"reference_id"`
	SourceID                 SourceID      `json:"source_id,omitempty"`
	Status                   Status        `json:"status"`
	Priority                 int           `json:"priority"`
	AttemptCount             int           `json:"attempt_count"`
	MaxAttempts              int           `json:"max_attempts"`
	NextAttemptAt            time.Time     `json:"next_attempt_at"`
	LockedBy                 WorkerID      `json:"locked_by,omitempty"`
	LockedUntil              time.Time     `json:"locked_until,omitempty"`
	LastErrorCode            string        `json:"last_error_code,omitempty"`
	LastErrorMessageRedacted string        `json:"last_error_message_redacted,omitempty"`
	ManualReviewReason       string        `json:"manual_review_reason,omitempty"`
	CorrelationID            CorrelationID `json:"correlation_id"`
	CreatedAt                time.Time     `json:"created_at"`
	UpdatedAt                time.Time     `json:"updated_at"`
	FinishedAt               time.Time     `json:"finished_at,omitempty"`
}

func newJobResponse(job Job) jobResponse {
	return jobResponse{
		ID:                       job.ID,
		DisplayID:                job.DisplayID,
		TenantID:                 job.TenantID,
		Type:                     job.Type,
		ReferenceType:            job.ReferenceType,
		ReferenceID:              job.ReferenceID,
		SourceID:                 job.SourceID,
		Status:                   job.Status,
		Priority:                 job.Priority,
		AttemptCount:             job.AttemptCount,
		MaxAttempts:              job.MaxAttempts,
		NextAttemptAt:            job.NextAttemptAt,
		LockedBy:                 job.LockedBy,
		LockedUntil:              job.LockedUntil,
		LastErrorCode:            job.LastErrorCode,
		LastErrorMessageRedacted: job.LastErrorMessageRedacted,
		ManualReviewReason:       job.ManualReviewReason,
		CorrelationID:            job.CorrelationID,
		CreatedAt:                job.CreatedAt,
		UpdatedAt:                job.UpdatedAt,
		FinishedAt:               job.FinishedAt,
	}
}

func newJobResponses(jobs []Job) []jobResponse {
	responses := make([]jobResponse, 0, len(jobs))
	for _, job := range jobs {
		responses = append(responses, newJobResponse(job))
	}
	return responses
}

type attemptResponse struct {
	ID                   AttemptID     `json:"id"`
	DisplayID            int64         `json:"display_id"`
	JobID                ID            `json:"job_id"`
	WorkerID             WorkerID      `json:"worker_id"`
	AttemptNumber        int           `json:"attempt_number"`
	StartedAt            time.Time     `json:"started_at"`
	FinishedAt           time.Time     `json:"finished_at,omitempty"`
	Result               AttemptResult `json:"result"`
	ErrorCode            string        `json:"error_code,omitempty"`
	ErrorMessageRedacted string        `json:"error_message_redacted,omitempty"`
	DurationMilliseconds int64         `json:"duration_ms,omitempty"`
	CorrelationID        CorrelationID `json:"correlation_id"`
}

type jobSummaryResponse struct {
	Type                   Type                       `json:"job_type"`
	Total                  int                        `json:"total"`
	AttentionCount         int                        `json:"attention_count"`
	Counts                 jobSummaryCountsResponse   `json:"counts"`
	OldestQueuedAt         *time.Time                 `json:"oldest_queued_at,omitempty"`
	OldestQueuedAgeSeconds int64                      `json:"oldest_queued_age_seconds,omitempty"`
	LatestFailure          *jobSummaryFailureResponse `json:"latest_failure,omitempty"`
	GeneratedAt            time.Time                  `json:"generated_at"`
}

type jobSummaryCountsResponse struct {
	Queued          int `json:"queued"`
	Claimed         int `json:"claimed"`
	Running         int `json:"running"`
	Succeeded       int `json:"succeeded"`
	FailedRetryable int `json:"failed_retryable"`
	FailedTerminal  int `json:"failed_terminal"`
	ManualReview    int `json:"manual_review"`
	Cancelled       int `json:"cancelled"`
}

type jobSummaryFailureResponse struct {
	ID                       ID        `json:"id"`
	DisplayID                int64     `json:"display_id"`
	Status                   Status    `json:"status"`
	LastErrorCode            string    `json:"last_error_code,omitempty"`
	LastErrorMessageRedacted string    `json:"last_error_message_redacted,omitempty"`
	ManualReviewReason       string    `json:"manual_review_reason,omitempty"`
	CreatedAt                time.Time `json:"created_at"`
	UpdatedAt                time.Time `json:"updated_at"`
}

func newJobSummaryResponse(summary JobSummary) jobSummaryResponse {
	var oldestQueuedAt *time.Time
	if !summary.OldestQueuedAt.IsZero() {
		oldestQueuedAt = &summary.OldestQueuedAt
	}
	return jobSummaryResponse{
		Type:           summary.Type,
		Total:          summary.Total,
		AttentionCount: summary.AttentionCount,
		Counts: jobSummaryCountsResponse{
			Queued:          summary.Counts.Queued,
			Claimed:         summary.Counts.Claimed,
			Running:         summary.Counts.Running,
			Succeeded:       summary.Counts.Succeeded,
			FailedRetryable: summary.Counts.FailedRetryable,
			FailedTerminal:  summary.Counts.FailedTerminal,
			ManualReview:    summary.Counts.ManualReview,
			Cancelled:       summary.Counts.Cancelled,
		},
		OldestQueuedAt:         oldestQueuedAt,
		OldestQueuedAgeSeconds: oldestQueuedAgeSeconds(summary),
		LatestFailure:          newJobSummaryFailureResponse(summary.LatestFailure),
		GeneratedAt:            summary.GeneratedAt,
	}
}

func newJobSummaryFailureResponse(failure *JobFailureContext) *jobSummaryFailureResponse {
	if failure == nil {
		return nil
	}
	return &jobSummaryFailureResponse{
		ID:                       failure.ID,
		DisplayID:                failure.DisplayID,
		Status:                   failure.Status,
		LastErrorCode:            failure.LastErrorCode,
		LastErrorMessageRedacted: failure.LastErrorMessageRedacted,
		ManualReviewReason:       failure.ManualReviewReason,
		CreatedAt:                failure.CreatedAt,
		UpdatedAt:                failure.UpdatedAt,
	}
}

func oldestQueuedAgeSeconds(summary JobSummary) int64 {
	if summary.OldestQueuedAt.IsZero() || summary.GeneratedAt.IsZero() {
		return 0
	}
	seconds := int64(summary.GeneratedAt.Sub(summary.OldestQueuedAt).Seconds())
	if seconds < 0 {
		return 0
	}
	return seconds
}

func newAttemptResponse(attempt Attempt) attemptResponse {
	return attemptResponse{
		ID:                   attempt.ID,
		DisplayID:            attempt.DisplayID,
		JobID:                attempt.JobID,
		WorkerID:             attempt.WorkerID,
		AttemptNumber:        attempt.AttemptNumber,
		StartedAt:            attempt.StartedAt,
		FinishedAt:           attempt.FinishedAt,
		Result:               attempt.Result,
		ErrorCode:            attempt.ErrorCode,
		ErrorMessageRedacted: attempt.ErrorMessageRedacted,
		DurationMilliseconds: attempt.Duration.Milliseconds(),
		CorrelationID:        attempt.CorrelationID,
	}
}

func newAttemptResponses(attempts []Attempt) []attemptResponse {
	responses := make([]attemptResponse, 0, len(attempts))
	for _, attempt := range attempts {
		responses = append(responses, newAttemptResponse(attempt))
	}
	return responses
}
