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

var ErrTwoFactorStoreExecutorMissing = errors.New("two factor store executor missing")

type PostgresTwoFactorStore struct {
	executor platformdb.Executor
}

func NewPostgresTwoFactorStore(executor platformdb.Executor) *PostgresTwoFactorStore {
	return &PostgresTwoFactorStore{executor: executor}
}

const twoFactorColumns = `tenant_id, user_id, method, secret_ciphertext, enabled_at, created_at, updated_at`

func (store *PostgresTwoFactorStore) UpsertTOTPSecret(ctx context.Context, input UpsertTOTPSecretInput) (TwoFactorMethod, error) {
	if err := store.ready(); err != nil {
		return TwoFactorMethod{}, err
	}
	if input.TenantID.Empty() {
		return TwoFactorMethod{}, tenant.ErrTenantIDMissing
	}
	if input.UserID == "" {
		return TwoFactorMethod{}, ErrUserIDMissing
	}
	if input.SecretCiphertext == "" {
		return TwoFactorMethod{}, ErrSecretCipherMissing
	}
	row := store.executor.QueryRowContext(ctx, `
INSERT INTO user_two_factor_methods (tenant_id, user_id, method, secret_ciphertext)
VALUES ($1, $2, 'totp', $3)
ON CONFLICT (tenant_id, user_id, method) DO UPDATE
SET secret_ciphertext = EXCLUDED.secret_ciphertext,
    enabled_at = NULL,
    updated_at = NOW()
RETURNING `+twoFactorColumns,
		input.TenantID, input.UserID, input.SecretCiphertext)
	return scanTwoFactorMethod(row)
}

func (store *PostgresTwoFactorStore) GetTOTPMethod(ctx context.Context, tenantID tenant.ID, userID UserID) (TwoFactorMethod, error) {
	if err := store.ready(); err != nil {
		return TwoFactorMethod{}, err
	}
	if tenantID.Empty() {
		return TwoFactorMethod{}, tenant.ErrTenantIDMissing
	}
	if userID == "" {
		return TwoFactorMethod{}, ErrUserIDMissing
	}
	row := store.executor.QueryRowContext(ctx, `SELECT `+twoFactorColumns+` FROM user_two_factor_methods WHERE tenant_id = $1 AND user_id = $2 AND method = 'totp'`, tenantID, userID)
	return scanTwoFactorMethod(row)
}

func (store *PostgresTwoFactorStore) MarkTOTPEnabled(ctx context.Context, tenantID tenant.ID, userID UserID, now time.Time) error {
	if err := store.ready(); err != nil {
		return err
	}
	result, err := store.executor.ExecContext(ctx, `
UPDATE user_two_factor_methods
SET enabled_at = COALESCE(enabled_at, $3), updated_at = $3
WHERE tenant_id = $1 AND user_id = $2 AND method = 'totp'`, tenantID, userID, now)
	if err != nil {
		return fmt.Errorf("mark totp enabled: %w", err)
	}
	count, err := result.RowsAffected()
	if err == nil && count == 0 {
		return ErrTwoFactorMethodNotFound
	}
	return nil
}

func (store *PostgresTwoFactorStore) SetUserTwoFactorStatus(ctx context.Context, tenantID tenant.ID, userID UserID, status TwoFactorStatus) error {
	if err := store.ready(); err != nil {
		return err
	}
	if tenantID.Empty() {
		return tenant.ErrTenantIDMissing
	}
	if userID == "" {
		return ErrUserIDMissing
	}
	if !status.Valid() {
		return ErrTwoFactorStatusInvalid
	}
	_, err := store.executor.ExecContext(ctx, `
UPDATE users
SET two_factor_status = $3, updated_at = NOW()
WHERE tenant_id = $1 AND user_id = $2`, tenantID, userID, status)
	if err != nil {
		return fmt.Errorf("set user two factor status: %w", err)
	}
	return nil
}

func (store *PostgresTwoFactorStore) ready() error {
	if store == nil || store.executor == nil {
		return ErrTwoFactorStoreExecutorMissing
	}
	return nil
}

func scanTwoFactorMethod(row sessionScanner) (TwoFactorMethod, error) {
	var record TwoFactorMethod
	var tenantID, userID string
	var enabledAt sql.NullTime
	if err := row.Scan(&tenantID, &userID, &record.Method, &record.SecretCiphertext, &enabledAt, &record.CreatedAt, &record.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return TwoFactorMethod{}, ErrTwoFactorMethodNotFound
		}
		return TwoFactorMethod{}, fmt.Errorf("scan two factor method: %w", err)
	}
	record.TenantID = tenant.ID(tenantID)
	record.UserID = UserID(userID)
	record.EnabledAt = enabledAt.Time
	return record, nil
}
