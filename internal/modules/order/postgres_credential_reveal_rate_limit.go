package order

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
	platformdb "github.com/Chinsusu/Billing-V2/internal/platform/db"
)

type PostgresCredentialRevealRateLimiter struct {
	executor platformdb.Executor
}

const incrementCredentialRevealRateLimitSQL = `
INSERT INTO service_credential_reveal_rate_limits (tenant_id, actor_id, service_instance_id, window_start, attempt_count)
VALUES ($1, $2, $3, $4, 1)
ON CONFLICT (tenant_id, actor_id, service_instance_id, window_start) DO UPDATE
SET attempt_count = service_credential_reveal_rate_limits.attempt_count + 1,
    updated_at = NOW()
RETURNING tenant_id, actor_id, service_instance_id, window_start, attempt_count, created_at, updated_at`

func NewPostgresCredentialRevealRateLimiter(executor platformdb.Executor) *PostgresCredentialRevealRateLimiter {
	return &PostgresCredentialRevealRateLimiter{executor: executor}
}

func (limiter *PostgresCredentialRevealRateLimiter) IncrementCredentialRevealRateLimit(
	ctx context.Context,
	input CredentialRevealRateLimitInput,
) (CredentialRevealRateLimitCounter, error) {
	if limiter == nil || limiter.executor == nil {
		return CredentialRevealRateLimitCounter{}, ErrCredentialRevealLimiterMissing
	}
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return CredentialRevealRateLimitCounter{}, err
	}
	row := limiter.executor.QueryRowContext(ctx, incrementCredentialRevealRateLimitSQL,
		input.TenantID, input.ActorID, input.ServiceID, input.WindowStart)
	return scanCredentialRevealRateLimitCounter(row)
}

func scanCredentialRevealRateLimitCounter(row credentialRevealRateLimitScanner) (CredentialRevealRateLimitCounter, error) {
	var record CredentialRevealRateLimitCounter
	var tenantID, actorID, serviceID string
	if err := row.Scan(
		&tenantID,
		&actorID,
		&serviceID,
		&record.WindowStart,
		&record.AttemptCount,
		&record.CreatedAt,
		&record.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return CredentialRevealRateLimitCounter{}, ErrCredentialRevealRateLimitWindowMissing
		}
		return CredentialRevealRateLimitCounter{}, fmt.Errorf("scan credential reveal rate limit counter: %w", err)
	}
	record.TenantID = tenant.ID(tenantID)
	record.ActorID = identity.UserID(actorID)
	record.ServiceID = ServiceID(serviceID)
	return record, nil
}

type credentialRevealRateLimitScanner interface {
	Scan(dest ...interface{}) error
}
