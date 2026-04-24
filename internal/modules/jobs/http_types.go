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
