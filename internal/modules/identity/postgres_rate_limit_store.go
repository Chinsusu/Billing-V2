package identity

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	platformdb "github.com/Chinsusu/Billing-V2/internal/platform/db"
)

var ErrAuthRateLimitStoreExecutorMissing = errors.New("auth rate limit store executor missing")

type PostgresAuthRateLimitStore struct {
	executor platformdb.Executor
}

func NewPostgresAuthRateLimitStore(executor platformdb.Executor) *PostgresAuthRateLimitStore {
	return &PostgresAuthRateLimitStore{executor: executor}
}

func (store *PostgresAuthRateLimitStore) IncrementAuthRateLimit(ctx context.Context, input AuthRateLimitIncrementInput) (AuthRateLimitCounter, error) {
	if err := store.ready(); err != nil {
		return AuthRateLimitCounter{}, err
	}
	if input.Action == "" {
		return AuthRateLimitCounter{}, ErrAuthRateLimitMissing
	}
	if input.KeyHash == "" {
		return AuthRateLimitCounter{}, ErrAuthRateLimitKeyEmpty
	}
	if input.WindowStart.IsZero() {
		return AuthRateLimitCounter{}, ErrAuthRateLimitMissing
	}
	row := store.executor.QueryRowContext(ctx, `
INSERT INTO auth_rate_limit_counters (action, key_hash, window_start, attempt_count)
VALUES ($1, $2, $3, 1)
ON CONFLICT (action, key_hash, window_start) DO UPDATE
SET attempt_count = auth_rate_limit_counters.attempt_count + 1,
    updated_at = NOW()
RETURNING action, key_hash, window_start, attempt_count, created_at, updated_at`,
		input.Action, input.KeyHash, input.WindowStart)
	return scanAuthRateLimitCounter(row)
}

func (store *PostgresAuthRateLimitStore) ready() error {
	if store == nil || store.executor == nil {
		return ErrAuthRateLimitStoreExecutorMissing
	}
	return nil
}

func scanAuthRateLimitCounter(row rateLimitScanner) (AuthRateLimitCounter, error) {
	var record AuthRateLimitCounter
	if err := row.Scan(&record.Action, &record.KeyHash, &record.WindowStart, &record.AttemptCount, &record.CreatedAt, &record.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return AuthRateLimitCounter{}, ErrAuthRateLimitMissing
		}
		return AuthRateLimitCounter{}, fmt.Errorf("scan auth rate limit counter: %w", err)
	}
	return record, nil
}

type rateLimitScanner interface {
	Scan(dest ...interface{}) error
}
