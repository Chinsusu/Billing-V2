package identity

import (
	"context"

	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

const (
	AuditActionTwoFactorSetup   = "auth.2fa.setup"
	AuditActionTwoFactorSuccess = "auth.2fa.success"
	AuditActionTwoFactorFailure = "auth.2fa.failure"
)

type AuthAuditInput struct {
	TenantID      tenant.ID
	ActorID       UserID
	Action        string
	TargetUserID  UserID
	CorrelationID string
}

type AuthAuditAppender interface {
	AppendAuthAudit(ctx context.Context, input AuthAuditInput) error
}

func (service *AuthService) appendTwoFactorAudit(ctx context.Context, identity SessionIdentity, action string) error {
	if service == nil || service.audit == nil {
		return nil
	}
	return service.audit.AppendAuthAudit(ctx, AuthAuditInput{
		TenantID:      identity.User.TenantID,
		ActorID:       identity.User.ID,
		Action:        action,
		TargetUserID:  identity.User.ID,
		CorrelationID: identity.Session.ID,
	})
}
