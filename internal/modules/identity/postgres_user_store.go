package identity

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
	platformdb "github.com/Chinsusu/Billing-V2/internal/platform/db"
)

var ErrUserStoreExecutorMissing = errors.New("user store executor missing")

type PostgresUserStore struct {
	executor platformdb.Executor
}

func NewPostgresUserStore(executor platformdb.Executor) *PostgresUserStore {
	return &PostgresUserStore{executor: executor}
}

const userColumns = `user_id, display_id, tenant_id, email, email_verified_at, password_hash, full_name, user_type, status, two_factor_status, last_login_at, failed_login_count, created_at, updated_at`

func (store *PostgresUserStore) CreateUser(ctx context.Context, input CreateUserInput) (User, error) {
	if err := store.ready(); err != nil {
		return User{}, err
	}
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return User{}, err
	}

	row := store.executor.QueryRowContext(ctx, `
INSERT INTO users (tenant_id, email, email_verified_at, password_hash, full_name, user_type, status, two_factor_status)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING `+userColumns,
		input.TenantID, input.Email, nullableTime(input.EmailVerifiedAt), input.PasswordHash, nullableString(input.FullName),
		input.Type, input.Status, input.TwoFactorStatus)
	return scanUser(row)
}

func (store *PostgresUserStore) GetUserByID(ctx context.Context, tenantID tenant.ID, userID UserID) (User, error) {
	if err := store.ready(); err != nil {
		return User{}, err
	}
	if tenantID.Empty() {
		return User{}, tenant.ErrTenantIDMissing
	}
	if userID == "" {
		return User{}, ErrUserIDMissing
	}
	row := store.executor.QueryRowContext(ctx, `SELECT `+userColumns+` FROM users WHERE tenant_id = $1 AND user_id = $2`, tenantID, userID)
	return scanUser(row)
}

func (store *PostgresUserStore) FindUserByEmail(ctx context.Context, tenantID tenant.ID, email string) (User, error) {
	if err := store.ready(); err != nil {
		return User{}, err
	}
	if tenantID.Empty() {
		return User{}, tenant.ErrTenantIDMissing
	}
	email = strings.ToLower(strings.TrimSpace(email))
	if email == "" {
		return User{}, ErrEmailMissing
	}
	row := store.executor.QueryRowContext(ctx, `SELECT `+userColumns+` FROM users WHERE tenant_id = $1 AND email = $2`, tenantID, email)
	return scanUser(row)
}

func (store *PostgresUserStore) ready() error {
	if store == nil || store.executor == nil {
		return ErrUserStoreExecutorMissing
	}
	return nil
}

type userScanner interface {
	Scan(dest ...interface{}) error
}

func scanUser(row userScanner) (User, error) {
	var record User
	var userID, tenantID, userType, status, twoFactorStatus string
	var emailVerifiedAt, lastLoginAt sql.NullTime
	var fullName sql.NullString

	if err := row.Scan(
		&userID, &record.DisplayID, &tenantID, &record.Email, &emailVerifiedAt, &record.PasswordHash, &fullName, &userType,
		&status, &twoFactorStatus, &lastLoginAt, &record.FailedLoginCount, &record.CreatedAt, &record.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, ErrUserNotFound
		}
		return User{}, fmt.Errorf("scan user: %w", err)
	}
	record.ID = UserID(userID)
	record.TenantID = tenant.ID(tenantID)
	record.EmailVerifiedAt = emailVerifiedAt.Time
	record.FullName = fullName.String
	record.Type = UserType(userType)
	record.Status = UserStatus(status)
	record.TwoFactorStatus = TwoFactorStatus(twoFactorStatus)
	record.LastLoginAt = lastLoginAt.Time
	return record, nil
}

func nullableTime(value time.Time) sql.NullTime {
	if value.IsZero() {
		return sql.NullTime{}
	}
	return sql.NullTime{Time: value, Valid: true}
}

func nullableString(value string) sql.NullString {
	if value == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: value, Valid: true}
}
