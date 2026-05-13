package identity

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
	platformdb "github.com/Chinsusu/Billing-V2/internal/platform/db"
	"github.com/lib/pq"
)

var ErrSessionStoreExecutorMissing = errors.New("session store executor missing")

type PostgresSessionStore struct {
	executor platformdb.Executor
}

func NewPostgresSessionStore(executor platformdb.Executor) *PostgresSessionStore {
	return &PostgresSessionStore{executor: executor}
}

const sessionColumns = `session_id, tenant_id, user_id, token_hash, user_agent_hash, expires_at, revoked_at, last_seen_at, two_factor_satisfied_at, created_at, updated_at`

func (store *PostgresSessionStore) CreateSession(ctx context.Context, input CreateSessionInput) (Session, error) {
	if err := store.ready(); err != nil {
		return Session{}, err
	}
	if input.TenantID.Empty() {
		return Session{}, tenant.ErrTenantIDMissing
	}
	if input.UserID == "" {
		return Session{}, ErrUserIDMissing
	}
	if input.TokenHash == "" {
		return Session{}, ErrSessionTokenMissing
	}
	if input.ExpiresAt.IsZero() {
		return Session{}, ErrSessionExpired
	}
	row := store.executor.QueryRowContext(ctx, `
INSERT INTO auth_sessions (tenant_id, user_id, token_hash, user_agent_hash, expires_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING `+sessionColumns,
		input.TenantID, input.UserID, input.TokenHash, nullableString(input.UserAgentHash), input.ExpiresAt)
	return scanSession(row)
}

func (store *PostgresSessionStore) FindSessionIdentityByTokenHash(ctx context.Context, tokenHash string, now time.Time) (SessionIdentity, error) {
	if err := store.ready(); err != nil {
		return SessionIdentity{}, err
	}
	if tokenHash == "" {
		return SessionIdentity{}, ErrSessionTokenMissing
	}
	row := store.executor.QueryRowContext(ctx, `
SELECT `+sessionSelectColumns("s")+`, `+userSelectColumns("u")+`, COALESCE(array_agg(ur.role_id ORDER BY ur.role_id) FILTER (WHERE ur.role_id IS NOT NULL), '{}') AS role_ids
FROM auth_sessions s
JOIN users u ON u.user_id = s.user_id AND u.tenant_id = s.tenant_id
LEFT JOIN user_roles ur ON ur.user_id = u.user_id AND ur.tenant_id = u.tenant_id
WHERE s.token_hash = $1
  AND s.revoked_at IS NULL
  AND s.expires_at > $2
  AND u.status = 'active'
GROUP BY s.session_id, u.user_id`, tokenHash, now)
	identity, err := scanSessionIdentity(row)
	if err != nil {
		return SessionIdentity{}, err
	}
	if _, err := store.executor.ExecContext(ctx, `UPDATE auth_sessions SET last_seen_at = $2, updated_at = $2 WHERE session_id = $1`, identity.Session.ID, now); err != nil {
		return SessionIdentity{}, fmt.Errorf("touch session: %w", err)
	}
	return identity, nil
}

func (store *PostgresSessionStore) RevokeSessionByTokenHash(ctx context.Context, tokenHash string, now time.Time) error {
	if err := store.ready(); err != nil {
		return err
	}
	if tokenHash == "" {
		return nil
	}
	if _, err := store.executor.ExecContext(ctx, `
UPDATE auth_sessions
SET revoked_at = COALESCE(revoked_at, $2), updated_at = $2
WHERE token_hash = $1`, tokenHash, now); err != nil {
		return fmt.Errorf("revoke session: %w", err)
	}
	return nil
}

func (store *PostgresSessionStore) RevokeUserSessions(ctx context.Context, tenantID tenant.ID, userID UserID, now time.Time) error {
	if err := store.ready(); err != nil {
		return err
	}
	if tenantID.Empty() {
		return tenant.ErrTenantIDMissing
	}
	if userID == "" {
		return ErrUserIDMissing
	}
	if _, err := store.executor.ExecContext(ctx, `
UPDATE auth_sessions
SET revoked_at = COALESCE(revoked_at, $3), updated_at = $3
WHERE tenant_id = $1 AND user_id = $2`, tenantID, userID, now); err != nil {
		return fmt.Errorf("revoke user sessions: %w", err)
	}
	return nil
}

func (store *PostgresSessionStore) MarkSessionTwoFactorSatisfied(ctx context.Context, tokenHash string, now time.Time) (Session, error) {
	if err := store.ready(); err != nil {
		return Session{}, err
	}
	if tokenHash == "" {
		return Session{}, ErrSessionTokenMissing
	}
	row := store.executor.QueryRowContext(ctx, `
UPDATE auth_sessions
SET two_factor_satisfied_at = $2, updated_at = $2
WHERE token_hash = $1
  AND revoked_at IS NULL
  AND expires_at > $2
RETURNING `+sessionColumns, tokenHash, now)
	return scanSession(row)
}

func (store *PostgresSessionStore) ready() error {
	if store == nil || store.executor == nil {
		return ErrSessionStoreExecutorMissing
	}
	return nil
}

type sessionScanner interface {
	Scan(dest ...interface{}) error
}

func sessionSelectColumns(alias string) string {
	columns := []string{
		"session_id",
		"tenant_id",
		"user_id",
		"token_hash",
		"user_agent_hash",
		"expires_at",
		"revoked_at",
		"last_seen_at",
		"two_factor_satisfied_at",
		"created_at",
		"updated_at",
	}
	for index, column := range columns {
		columns[index] = alias + "." + column
	}
	return joinSessionColumns(columns)
}

func joinSessionColumns(columns []string) string {
	value := ""
	for index, column := range columns {
		if index > 0 {
			value += ", "
		}
		value += column
	}
	return value
}

func scanSession(row sessionScanner) (Session, error) {
	var record Session
	var tenantID, userID string
	var userAgentHash sql.NullString
	var revokedAt, lastSeenAt, twoFactorSatisfiedAt sql.NullTime
	if err := row.Scan(
		&record.ID, &tenantID, &userID, &record.TokenHash, &userAgentHash, &record.ExpiresAt,
		&revokedAt, &lastSeenAt, &twoFactorSatisfiedAt, &record.CreatedAt, &record.UpdatedAt,
	); err != nil {
		return Session{}, mapSessionError(err)
	}
	record.TenantID = tenant.ID(tenantID)
	record.UserID = UserID(userID)
	record.UserAgentHash = userAgentHash.String
	record.RevokedAt = revokedAt.Time
	record.LastSeenAt = lastSeenAt.Time
	record.TwoFactorSatisfiedAt = twoFactorSatisfiedAt.Time
	return record, nil
}

func scanSessionIdentity(row sessionScanner) (SessionIdentity, error) {
	var session Session
	var user User
	var sessionTenantID, sessionUserID, userID, userTenantID, userType, status, twoFactorStatus string
	var userAgentHash, fullName sql.NullString
	var revokedAt, lastSeenAt, twoFactorSatisfiedAt, emailVerifiedAt, userLastLoginAt sql.NullTime
	var roleIDs []string
	if err := row.Scan(
		&session.ID, &sessionTenantID, &sessionUserID, &session.TokenHash, &userAgentHash, &session.ExpiresAt,
		&revokedAt, &lastSeenAt, &twoFactorSatisfiedAt, &session.CreatedAt, &session.UpdatedAt,
		&userID, &user.DisplayID, &userTenantID, &user.Email, &emailVerifiedAt, &user.PasswordHash, &fullName, &userType,
		&status, &twoFactorStatus, &userLastLoginAt, &user.FailedLoginCount, &user.CreatedAt, &user.UpdatedAt,
		pq.Array(&roleIDs),
	); err != nil {
		return SessionIdentity{}, mapSessionError(err)
	}
	session.TenantID = tenant.ID(sessionTenantID)
	session.UserID = UserID(sessionUserID)
	session.UserAgentHash = userAgentHash.String
	session.RevokedAt = revokedAt.Time
	session.LastSeenAt = lastSeenAt.Time
	session.TwoFactorSatisfiedAt = twoFactorSatisfiedAt.Time
	user.ID = UserID(userID)
	user.TenantID = tenant.ID(userTenantID)
	user.EmailVerifiedAt = emailVerifiedAt.Time
	user.FullName = fullName.String
	user.Type = UserType(userType)
	user.Status = UserStatus(status)
	user.TwoFactorStatus = TwoFactorStatus(twoFactorStatus)
	user.LastLoginAt = userLastLoginAt.Time
	return SessionIdentity{
		Session: session,
		User:    user,
		RoleIDs: roleIDsFromStrings(roleIDs),
	}, nil
}

func mapSessionError(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return ErrSessionInvalid
	}
	return fmt.Errorf("scan session: %w", err)
}

func roleIDsFromStrings(values []string) []RoleID {
	roleIDs := make([]RoleID, 0, len(values))
	for _, value := range values {
		if value == "" {
			continue
		}
		roleIDs = append(roleIDs, RoleID(value))
	}
	return roleIDs
}
