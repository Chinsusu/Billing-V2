package identity

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

var (
	ErrPasswordResetStoreMissing = errors.New("password reset store missing")
	ErrPasswordResetTTLMissing   = errors.New("password reset ttl missing")
	ErrPasswordResetTokenMissing = errors.New("password reset token missing")
	ErrPasswordResetTokenInvalid = errors.New("password reset token invalid")
	ErrPasswordResetTokenExpired = errors.New("password reset token expired")
	ErrPasswordResetTokenUsed    = errors.New("password reset token used")
)

type PasswordResetRequestInput struct {
	Email                  string
	Domain                 string
	LocalTenantID          tenant.ID
	AllowLocalTenantHeader bool
	ClientIP               string
}

type PasswordResetRequestResult struct {
	Accepted bool
}

type PasswordResetConfirmInput struct {
	Token       string
	NewPassword string
	ClientIP    string
}

type PasswordResetConfirmResult struct {
	PasswordUpdated bool
}

type PasswordResetToken struct {
	ID        string
	TenantID  tenant.ID
	UserID    UserID
	TokenHash string
	ExpiresAt time.Time
	UsedAt    time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreatePasswordResetTokenInput struct {
	TenantID  tenant.ID
	UserID    UserID
	TokenHash string
	ExpiresAt time.Time
}

type PasswordResetStore interface {
	CreatePasswordResetToken(ctx context.Context, input CreatePasswordResetTokenInput) (PasswordResetToken, error)
	UsePasswordResetToken(ctx context.Context, tokenHash string, now time.Time) (PasswordResetToken, error)
}

type PasswordResetDeliveryInput struct {
	TenantID  tenant.ID
	UserID    UserID
	Email     string
	Token     string
	ExpiresAt time.Time
}

type PasswordResetDelivery interface {
	DeliverPasswordReset(ctx context.Context, input PasswordResetDeliveryInput) error
}

type NoopPasswordResetDelivery struct{}

func (NoopPasswordResetDelivery) DeliverPasswordReset(ctx context.Context, input PasswordResetDeliveryInput) error {
	return nil
}

func (service *AuthService) RequestPasswordReset(ctx context.Context, input PasswordResetRequestInput) (PasswordResetRequestResult, error) {
	if service == nil || service.tenants == nil || service.users == nil || service.passwordResets == nil || service.resetDelivery == nil {
		return PasswordResetRequestResult{}, ErrAuthStoreMissing
	}
	if service.passwordResetTTL <= 0 {
		return PasswordResetRequestResult{}, ErrPasswordResetTTLMissing
	}
	input.Email = strings.ToLower(strings.TrimSpace(input.Email))
	if input.Email == "" {
		return PasswordResetRequestResult{}, ErrEmailMissing
	}
	tenantID, err := service.resolvePasswordResetTenant(ctx, input)
	if err != nil {
		return PasswordResetRequestResult{}, err
	}
	if err := service.enforceAuthRateLimit(ctx, AuthActionPasswordReset, tenantID, input.Email, input.ClientIP); err != nil {
		return PasswordResetRequestResult{}, err
	}
	user, err := service.users.FindUserByEmail(ctx, tenantID, input.Email)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return PasswordResetRequestResult{Accepted: true}, nil
		}
		return PasswordResetRequestResult{}, err
	}
	if user.Status != UserStatusActive {
		return PasswordResetRequestResult{Accepted: true}, nil
	}
	token, tokenHash, err := newPasswordResetToken()
	if err != nil {
		return PasswordResetRequestResult{}, err
	}
	expiresAt := service.now().UTC().Add(service.passwordResetTTL)
	record, err := service.passwordResets.CreatePasswordResetToken(ctx, CreatePasswordResetTokenInput{
		TenantID:  user.TenantID,
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: expiresAt,
	})
	if err != nil {
		return PasswordResetRequestResult{}, err
	}
	if err := service.resetDelivery.DeliverPasswordReset(ctx, PasswordResetDeliveryInput{
		TenantID:  user.TenantID,
		UserID:    user.ID,
		Email:     user.Email,
		Token:     token,
		ExpiresAt: record.ExpiresAt,
	}); err != nil {
		return PasswordResetRequestResult{}, err
	}
	return PasswordResetRequestResult{Accepted: true}, nil
}

func (service *AuthService) ConfirmPasswordReset(ctx context.Context, input PasswordResetConfirmInput) (PasswordResetConfirmResult, error) {
	if service == nil || service.users == nil || service.passwordResets == nil {
		return PasswordResetConfirmResult{}, ErrAuthStoreMissing
	}
	input.Token = strings.TrimSpace(input.Token)
	if input.Token == "" {
		return PasswordResetConfirmResult{}, ErrPasswordResetTokenMissing
	}
	if input.NewPassword == "" {
		return PasswordResetConfirmResult{}, ErrPasswordMissing
	}
	if input.ClientIP != "" {
		if err := service.enforceAuthRateLimit(ctx, AuthActionPasswordReset, "", "", input.ClientIP); err != nil {
			return PasswordResetConfirmResult{}, err
		}
	}
	tokenRecord, err := service.passwordResets.UsePasswordResetToken(ctx, HashPasswordResetToken(input.Token), service.now().UTC())
	if err != nil {
		return PasswordResetConfirmResult{}, err
	}
	passwordHash, err := HashPasswordArgon2id(input.NewPassword)
	if err != nil {
		return PasswordResetConfirmResult{}, err
	}
	if service.sessions != nil {
		if err := service.sessions.RevokeUserSessions(ctx, tokenRecord.TenantID, tokenRecord.UserID, service.now().UTC()); err != nil {
			return PasswordResetConfirmResult{}, err
		}
	}
	if err := service.users.UpdatePasswordHash(ctx, tokenRecord.TenantID, tokenRecord.UserID, passwordHash); err != nil {
		return PasswordResetConfirmResult{}, err
	}
	return PasswordResetConfirmResult{PasswordUpdated: true}, nil
}

func (service *AuthService) resolvePasswordResetTenant(ctx context.Context, input PasswordResetRequestInput) (tenant.ID, error) {
	return service.resolveLoginTenant(ctx, LoginInput{
		Email:                  input.Email,
		Password:               "password-reset-tenant-resolution",
		Domain:                 input.Domain,
		LocalTenantID:          input.LocalTenantID,
		AllowLocalTenantHeader: input.AllowLocalTenantHeader,
	})
}

func newPasswordResetToken() (string, string, error) {
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return "", "", fmt.Errorf("generate password reset token: %w", err)
	}
	token := base64.RawURLEncoding.EncodeToString(raw)
	return token, HashPasswordResetToken(token), nil
}

func HashPasswordResetToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
