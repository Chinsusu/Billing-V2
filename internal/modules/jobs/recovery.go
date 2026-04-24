package jobs

import (
	"context"
	"strings"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

type RecoveryStore interface {
	RetryJob(ctx context.Context, input RetryJobInput) (Job, error)
	MarkManualReview(ctx context.Context, input ManualReviewJobInput) (Job, error)
	CancelJob(ctx context.Context, input CancelJobInput) (Job, error)
}

type RetryJobInput struct {
	ID            ID
	TenantID      tenant.ID
	ActorID       identity.UserID
	NextAttemptAt time.Time
	Now           time.Time
}

type ManualReviewJobInput struct {
	ID       ID
	TenantID tenant.ID
	ActorID  identity.UserID
	Reason   string
	Now      time.Time
}

type CancelJobInput struct {
	ID       ID
	TenantID tenant.ID
	ActorID  identity.UserID
	Reason   string
	Now      time.Time
}

func (service *Service) RetryJob(ctx context.Context, input RetryJobInput) (Job, error) {
	if err := service.readyRecovery(); err != nil {
		return Job{}, err
	}
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return Job{}, err
	}
	before, err := service.store.GetJob(ctx, Lookup{ID: input.ID, TenantID: input.TenantID})
	if err != nil {
		return Job{}, err
	}
	after, err := service.recovery.RetryJob(ctx, input)
	if err != nil {
		return Job{}, err
	}
	if err := service.appendJobRecoveryAudit(ctx, jobAuditActionRetry, input.ActorID, before, after, ""); err != nil {
		return Job{}, err
	}
	return after, nil
}

func (service *Service) MarkManualReview(ctx context.Context, input ManualReviewJobInput) (Job, error) {
	if err := service.readyRecovery(); err != nil {
		return Job{}, err
	}
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return Job{}, err
	}
	before, err := service.store.GetJob(ctx, Lookup{ID: input.ID, TenantID: input.TenantID})
	if err != nil {
		return Job{}, err
	}
	after, err := service.recovery.MarkManualReview(ctx, input)
	if err != nil {
		return Job{}, err
	}
	if err := service.appendJobRecoveryAudit(ctx, jobAuditActionManualReview, input.ActorID, before, after, input.Reason); err != nil {
		return Job{}, err
	}
	return after, nil
}

func (service *Service) CancelJob(ctx context.Context, input CancelJobInput) (Job, error) {
	if err := service.readyRecovery(); err != nil {
		return Job{}, err
	}
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return Job{}, err
	}
	before, err := service.store.GetJob(ctx, Lookup{ID: input.ID, TenantID: input.TenantID})
	if err != nil {
		return Job{}, err
	}
	after, err := service.recovery.CancelJob(ctx, input)
	if err != nil {
		return Job{}, err
	}
	if err := service.appendJobRecoveryAudit(ctx, jobAuditActionCancel, input.ActorID, before, after, input.Reason); err != nil {
		return Job{}, err
	}
	return after, nil
}

func (input RetryJobInput) Normalize() RetryJobInput {
	input.ID = ID(trimJobString(string(input.ID)))
	input.TenantID = tenant.ID(trimJobString(string(input.TenantID)))
	input.ActorID = identity.UserID(trimJobString(string(input.ActorID)))
	return input
}

func (input RetryJobInput) Validate() error {
	return validateRecoveryBase(input.ID, input.TenantID, input.ActorID)
}

func (input ManualReviewJobInput) Normalize() ManualReviewJobInput {
	input.ID = ID(trimJobString(string(input.ID)))
	input.TenantID = tenant.ID(trimJobString(string(input.TenantID)))
	input.ActorID = identity.UserID(trimJobString(string(input.ActorID)))
	input.Reason = strings.TrimSpace(input.Reason)
	return input
}

func (input ManualReviewJobInput) Validate() error {
	if err := validateRecoveryBase(input.ID, input.TenantID, input.ActorID); err != nil {
		return err
	}
	if input.Reason == "" {
		return ErrManualReviewReasonMissing
	}
	return nil
}

func (input CancelJobInput) Normalize() CancelJobInput {
	input.ID = ID(trimJobString(string(input.ID)))
	input.TenantID = tenant.ID(trimJobString(string(input.TenantID)))
	input.ActorID = identity.UserID(trimJobString(string(input.ActorID)))
	input.Reason = strings.TrimSpace(input.Reason)
	return input
}

func (input CancelJobInput) Validate() error {
	return validateRecoveryBase(input.ID, input.TenantID, input.ActorID)
}

func validateRecoveryBase(id ID, tenantID tenant.ID, actorID identity.UserID) error {
	if id == "" {
		return ErrJobIDMissing
	}
	if tenantID.Empty() {
		return tenant.ErrTenantIDMissing
	}
	if actorID == "" {
		return identity.ErrActorIDMissing
	}
	return nil
}
