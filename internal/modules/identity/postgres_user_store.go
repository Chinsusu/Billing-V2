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

var userColumnNames = []string{
	"user_id",
	"display_id",
	"tenant_id",
	"email",
	"email_verified_at",
	"password_hash",
	"full_name",
	"user_type",
	"status",
	"two_factor_status",
	"last_login_at",
	"failed_login_count",
	"created_at",
	"updated_at",
}

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

func (store *PostgresUserStore) ListUsers(ctx context.Context, filter UserListFilter) ([]UserSummary, error) {
	if err := store.ready(); err != nil {
		return nil, err
	}
	query, args := buildListUsersQuery(filter)
	rows, err := store.executor.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}
	defer rows.Close()

	records := []UserSummary{}
	for rows.Next() {
		record, err := scanUserSummary(rows)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("read users: %w", err)
	}
	return records, nil
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

func buildListUsersQuery(filter UserListFilter) (string, []interface{}) {
	args := []interface{}{}
	conditions := []string{}
	if !filter.TenantID.Empty() {
		args = append(args, filter.TenantID)
		conditions = append(conditions, fmt.Sprintf("u.tenant_id = $%d", len(args)))
	}
	if filter.Type != "" {
		args = append(args, filter.Type)
		conditions = append(conditions, fmt.Sprintf("u.user_type = $%d", len(args)))
	}
	if filter.Status != "" {
		args = append(args, filter.Status)
		conditions = append(conditions, fmt.Sprintf("u.status = $%d", len(args)))
	}
	if filter.DisplayID > 0 {
		args = append(args, filter.DisplayID)
		conditions = append(conditions, fmt.Sprintf("u.display_id = $%d", len(args)))
	}
	if filter.Email != "" {
		args = append(args, strings.ToLower(strings.TrimSpace(filter.Email)))
		conditions = append(conditions, fmt.Sprintf("u.email = $%d", len(args)))
	}

	limit := filter.Limit
	if limit <= 0 {
		limit = 20
	}
	args = append(args, limit)
	limitPlaceholder := fmt.Sprintf("$%d", len(args))

	var builder strings.Builder
	builder.WriteString("SELECT ")
	builder.WriteString(userSelectColumns("u"))
	builder.WriteString(`,
       t.name AS tenant_name,
       t.slug AS tenant_slug
FROM users u
JOIN tenants t ON t.tenant_id = u.tenant_id`)
	if len(conditions) > 0 {
		builder.WriteString("\nWHERE ")
		builder.WriteString(strings.Join(conditions, "\n  AND "))
	}
	builder.WriteString("\nORDER BY u.created_at DESC, u.display_id DESC")
	builder.WriteString("\nLIMIT ")
	builder.WriteString(limitPlaceholder)
	return builder.String(), args
}

func userSelectColumns(alias string) string {
	columns := make([]string, 0, len(userColumnNames))
	for _, column := range userColumnNames {
		if alias == "" {
			columns = append(columns, column)
			continue
		}
		columns = append(columns, alias+"."+column)
	}
	return strings.Join(columns, ", ")
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

func scanUserSummary(row userScanner) (UserSummary, error) {
	var record User
	var userID, tenantID, userType, status, twoFactorStatus string
	var emailVerifiedAt, lastLoginAt sql.NullTime
	var fullName sql.NullString
	var tenantName, tenantSlug string

	if err := row.Scan(
		&userID, &record.DisplayID, &tenantID, &record.Email, &emailVerifiedAt, &record.PasswordHash, &fullName, &userType,
		&status, &twoFactorStatus, &lastLoginAt, &record.FailedLoginCount, &record.CreatedAt, &record.UpdatedAt,
		&tenantName, &tenantSlug,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return UserSummary{}, ErrUserNotFound
		}
		return UserSummary{}, fmt.Errorf("scan user summary: %w", err)
	}
	record.ID = UserID(userID)
	record.TenantID = tenant.ID(tenantID)
	record.EmailVerifiedAt = emailVerifiedAt.Time
	record.FullName = fullName.String
	record.Type = UserType(userType)
	record.Status = UserStatus(status)
	record.TwoFactorStatus = TwoFactorStatus(twoFactorStatus)
	record.LastLoginAt = lastLoginAt.Time
	return UserSummary{
		User:       record,
		TenantName: tenantName,
		TenantSlug: tenantSlug,
	}, nil
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
