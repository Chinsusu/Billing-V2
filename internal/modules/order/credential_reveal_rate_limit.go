package order

import (
	"context"
	"errors"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

const (
	credentialRevealRateLimitMaxAttempts = 5
	credentialRevealRateLimitWindow      = 15 * time.Minute
)

var (
	ErrCredentialRevealRateLimitWindowMissing = errors.New("credential reveal rate limit window missing")
)

type CredentialRevealRateLimitInput struct {
	TenantID    tenant.ID
	ActorID     identity.UserID
	ServiceID   ServiceID
	WindowStart time.Time
}

type CredentialRevealRateLimitCounter struct {
	TenantID     tenant.ID
	ActorID      identity.UserID
	ServiceID    ServiceID
	WindowStart  time.Time
	AttemptCount int
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type CredentialRevealRateLimiter interface {
	IncrementCredentialRevealRateLimit(ctx context.Context, input CredentialRevealRateLimitInput) (CredentialRevealRateLimitCounter, error)
}

func (input CredentialRevealRateLimitInput) Normalize() CredentialRevealRateLimitInput {
	return CredentialRevealRateLimitInput{
		TenantID:    tenant.ID(trim(string(input.TenantID))),
		ActorID:     identity.UserID(trim(string(input.ActorID))),
		ServiceID:   ServiceID(trim(string(input.ServiceID))),
		WindowStart: input.WindowStart,
	}
}

func (input CredentialRevealRateLimitInput) Validate() error {
	if input.TenantID.Empty() {
		return tenant.ErrTenantIDMissing
	}
	if trim(string(input.ActorID)) == "" {
		return identity.ErrActorIDMissing
	}
	if input.ServiceID.Empty() {
		return ErrServiceIDMissing
	}
	if input.WindowStart.IsZero() {
		return ErrCredentialRevealRateLimitWindowMissing
	}
	return nil
}
