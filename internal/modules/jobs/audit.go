package jobs

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/audit"
	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
)

const (
	jobAuditActionRetry        = "job.retry"
	jobAuditActionManualReview = "job.manual_review"
	jobAuditActionCancel       = "job.cancel"
)

type AuditAppender interface {
	Append(ctx context.Context, input audit.AppendInput) (audit.Log, error)
}

func (service *Service) appendJobRecoveryAudit(
	ctx context.Context,
	action string,
	actorID identity.UserID,
	before Job,
	after Job,
	reason string,
) error {
	if service.audit == nil {
		return nil
	}
	_, err := service.audit.Append(ctx, audit.AppendInput{
		TenantID:               after.TenantID,
		ActorID:                audit.ActorID(actorID),
		ActorType:              audit.ActorTypeUser,
		Action:                 action,
		TargetType:             "job",
		TargetID:               audit.TargetID(after.ID),
		BeforeSnapshotRedacted: jobAuditJSON(jobAuditStateFromJob(before)),
		AfterSnapshotRedacted:  jobAuditJSON(jobAuditStateFromJob(after)),
		MetadataRedacted: jobAuditJSON(jobAuditMetadata{
			DisplayID:     after.DisplayID,
			JobType:       after.Type,
			ReferenceType: after.ReferenceType,
			ReferenceID:   after.ReferenceID,
			SourceID:      after.SourceID,
			Reason:        reason,
		}),
		CorrelationID: audit.CorrelationID(after.CorrelationID),
	})
	return err
}

type jobAuditState struct {
	Status                   Status `json:"status"`
	NextAttemptAt            string `json:"next_attempt_at,omitempty"`
	LastErrorCode            string `json:"last_error_code,omitempty"`
	LastErrorMessageRedacted string `json:"last_error_message_redacted,omitempty"`
	ManualReviewReason       string `json:"manual_review_reason,omitempty"`
	FinishedAt               string `json:"finished_at,omitempty"`
}

type jobAuditMetadata struct {
	DisplayID     int64         `json:"display_id"`
	JobType       Type          `json:"job_type"`
	ReferenceType ReferenceType `json:"reference_type"`
	ReferenceID   ReferenceID   `json:"reference_id"`
	SourceID      SourceID      `json:"source_id,omitempty"`
	Reason        string        `json:"reason,omitempty"`
}

func jobAuditStateFromJob(job Job) jobAuditState {
	return jobAuditState{
		Status:                   job.Status,
		NextAttemptAt:            auditTime(job.NextAttemptAt),
		LastErrorCode:            job.LastErrorCode,
		LastErrorMessageRedacted: job.LastErrorMessageRedacted,
		ManualReviewReason:       job.ManualReviewReason,
		FinishedAt:               auditTime(job.FinishedAt),
	}
}

func auditTime(value time.Time) string {
	if value.IsZero() {
		return ""
	}
	return value.UTC().Format(time.RFC3339)
}

func jobAuditJSON(value interface{}) json.RawMessage {
	data, err := json.Marshal(value)
	if err != nil {
		return json.RawMessage(`{}`)
	}
	return data
}
