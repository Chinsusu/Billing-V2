package identity

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
	platformdb "github.com/Chinsusu/Billing-V2/internal/platform/db"
)

var ErrPasswordResetStoreExecutorMissing = errors.New("password reset store executor missing")

type PostgresPasswordResetStore struct {
	executor platformdb.Executor
}

func NewPostgresPasswordResetStore(executor platformdb.Executor) *PostgresPasswordResetStore {
	return &PostgresPasswordResetStore{executor: executor}
}

const passwordResetTokenColumns = `reset_token_id, tenant_id, user_id, token_hash, expires_at, used_at, created_at, updated_at`

func (store *PostgresPasswordResetStore) CreatePasswordResetToken(ctx context.Context, input CreatePasswordResetTokenInput) (PasswordResetToken, error) {
	if err := store.ready(); err != nil {
		return PasswordResetToken{}, err
	}
	if input.TenantID.Empty() {
		return PasswordResetToken{}, tenant.ErrTenantIDMissing
	}
	if input.UserID == "" {
		return PasswordResetToken{}, ErrUserIDMissing
	}
	if input.TokenHash == "" {
		return PasswordResetToken{}, ErrPasswordResetTokenMissing
	}
	if input.ExpiresAt.IsZero() {
		return PasswordResetToken{}, ErrPasswordResetTokenExpired
	}
	row := store.executor.QueryRowContext(ctx, `
INSERT INTO auth_password_reset_tokens (tenant_id, user_id, token_hash, expires_at)
VALUES ($1, $2, $3, $4)
RETURNING `+passwordResetTokenColumns,
		input.TenantID, input.UserID, input.TokenHash, input.ExpiresAt)
	return scanPasswordResetToken(row)
}

func (store *PostgresPasswordResetStore) UsePasswordResetToken(ctx context.Context, tokenHash string, now time.Time) (PasswordResetToken, error) {
	if err := store.ready(); err != nil {
		return PasswordResetToken{}, err
	}
	if tokenHash == "" {
		return PasswordResetToken{}, ErrPasswordResetTokenMissing
	}
	row := store.executor.QueryRowContext(ctx, `
UPDATE auth_password_reset_tokens
SET used_at = $2, updated_at = $2
WHERE token_hash = $1
  AND used_at IS NULL
  AND expires_at > $2
RETURNING `+passwordResetTokenColumns, tokenHash, now)
	record, err := scanPasswordResetToken(row)
	if err == nil {
		return record, nil
	}
	if !errors.Is(err, ErrPasswordResetTokenInvalid) {
		return PasswordResetToken{}, err
	}
	existing, findErr := store.findPasswordResetTokenByHash(ctx, tokenHash)
	if findErr != nil {
		return PasswordResetToken{}, findErr
	}
	if !existing.UsedAt.IsZero() {
		return PasswordResetToken{}, ErrPasswordResetTokenUsed
	}
	if !existing.ExpiresAt.After(now) {
		return PasswordResetToken{}, ErrPasswordResetTokenExpired
	}
	return PasswordResetToken{}, ErrPasswordResetTokenUsed
}

func (store *PostgresPasswordResetStore) findPasswordResetTokenByHash(ctx context.Context, tokenHash string) (PasswordResetToken, error) {
	row := store.executor.QueryRowContext(ctx, `SELECT `+passwordResetTokenColumns+` FROM auth_password_reset_tokens WHERE token_hash = $1`, tokenHash)
	return scanPasswordResetToken(row)
}

func (store *PostgresPasswordResetStore) ready() error {
	if store == nil || store.executor == nil {
		return ErrPasswordResetStoreExecutorMissing
	}
	return nil
}

func scanPasswordResetToken(row passwordResetScanner) (PasswordResetToken, error) {
	var record PasswordResetToken
	var tenantID, userID string
	var usedAt sql.NullTime
	if err := row.Scan(&record.ID, &tenantID, &userID, &record.TokenHash, &record.ExpiresAt, &usedAt, &record.CreatedAt, &record.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return PasswordResetToken{}, ErrPasswordResetTokenInvalid
		}
		return PasswordResetToken{}, fmt.Errorf("scan password reset token: %w", err)
	}
	record.TenantID = tenant.ID(tenantID)
	record.UserID = UserID(userID)
	record.UsedAt = usedAt.Time
	return record, nil
}

type passwordResetScanner interface {
	Scan(dest ...interface{}) error
}
