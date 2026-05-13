package identity

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"
)

func TestAuthServiceRequestPasswordResetCreatesHashedTokenAndDeliversPlainToken(t *testing.T) {
	now := time.Date(2026, 5, 13, 9, 0, 0, 0, time.UTC)
	resetStore := &fakePasswordResetStore{}
	delivery := &fakePasswordResetDelivery{}
	service := newTestAuthService("hash", &fakeSessionStore{}, now)
	service.passwordResets = resetStore
	service.resetDelivery = delivery

	result, err := service.RequestPasswordReset(context.Background(), PasswordResetRequestInput{
		Email:                  " ADMIN@LOCAL.BILLING ",
		LocalTenantID:          "tenant_1",
		AllowLocalTenantHeader: true,
		ClientIP:               "192.0.2.10",
	})
	if err != nil {
		t.Fatalf("RequestPasswordReset returned error: %v", err)
	}
	if !result.Accepted {
		t.Fatal("expected accepted reset request")
	}
	if delivery.input.Token == "" {
		t.Fatal("expected delivery boundary to receive plaintext token")
	}
	if resetStore.created.TokenHash == "" || resetStore.created.TokenHash != HashPasswordResetToken(delivery.input.Token) {
		t.Fatalf("expected stored token hash to match delivered token hash, got stored=%q delivered=%q", resetStore.created.TokenHash, delivery.input.Token)
	}
	if strings.Contains(resetStore.created.TokenHash, delivery.input.Token) {
		t.Fatalf("stored token hash must not contain plaintext token")
	}
	if !resetStore.created.ExpiresAt.Equal(now.Add(30 * time.Minute)) {
		t.Fatalf("expected reset expiry, got %s", resetStore.created.ExpiresAt)
	}
}

func TestAuthServiceRequestPasswordResetRejectsRateLimitedAttempt(t *testing.T) {
	service := newTestAuthService("hash", &fakeSessionStore{}, time.Date(2026, 5, 13, 9, 0, 0, 0, time.UTC))
	delivery := &fakePasswordResetDelivery{}
	service.resetDelivery = delivery
	service.rateLimits = &fakeAuthRateLimitStore{blocked: true}

	_, err := service.RequestPasswordReset(context.Background(), PasswordResetRequestInput{
		Email:                  "admin@local.billing",
		LocalTenantID:          "tenant_1",
		AllowLocalTenantHeader: true,
		ClientIP:               "192.0.2.10",
	})
	if !errors.Is(err, ErrAuthRateLimited) {
		t.Fatalf("expected rate limit error, got %v", err)
	}
	if delivery.input.Token != "" {
		t.Fatalf("rate-limited reset must not deliver token: %+v", delivery.input)
	}
}

func TestAuthServiceConfirmPasswordResetUpdatesPasswordAndUsesToken(t *testing.T) {
	now := time.Date(2026, 5, 13, 9, 0, 0, 0, time.UTC)
	token := "reset-token"
	resetStore := &fakePasswordResetStore{
		token: PasswordResetToken{
			TenantID:  "tenant_1",
			UserID:    "user_1",
			TokenHash: HashPasswordResetToken(token),
			ExpiresAt: now.Add(30 * time.Minute),
		},
	}
	service := newTestAuthService("old-hash", &fakeSessionStore{}, now)
	sessions := &fakeSessionStore{}
	service.sessions = sessions
	service.passwordResets = resetStore

	result, err := service.ConfirmPasswordReset(context.Background(), PasswordResetConfirmInput{
		Token:       token,
		NewPassword: "new-password",
		ClientIP:    "192.0.2.10",
	})
	if err != nil {
		t.Fatalf("ConfirmPasswordReset returned error: %v", err)
	}
	if !result.PasswordUpdated {
		t.Fatal("expected password update result")
	}
	if resetStore.usedHash != HashPasswordResetToken(token) {
		t.Fatalf("expected used token hash, got %q", resetStore.usedHash)
	}
	userStore := service.users.(*fakeAuthUserStore)
	if userStore.user.PasswordHash == "" || userStore.user.PasswordHash == "new-password" || userStore.user.PasswordHash == "old-hash" {
		t.Fatalf("expected new argon2id password hash, got %q", userStore.user.PasswordHash)
	}
	if sessions.revokedUserID != "user_1" {
		t.Fatalf("expected existing sessions to be revoked, got %q", sessions.revokedUserID)
	}
}

func TestAuthServiceConfirmPasswordResetRejectsExpiredToken(t *testing.T) {
	service := newTestAuthService("hash", &fakeSessionStore{}, time.Date(2026, 5, 13, 9, 0, 0, 0, time.UTC))
	service.passwordResets = &fakePasswordResetStore{useErr: ErrPasswordResetTokenExpired}

	_, err := service.ConfirmPasswordReset(context.Background(), PasswordResetConfirmInput{
		Token:       "reset-token",
		NewPassword: "new-password",
		ClientIP:    "192.0.2.10",
	})
	if !errors.Is(err, ErrPasswordResetTokenExpired) {
		t.Fatalf("expected expired token error, got %v", err)
	}
}

func TestAuthServiceConfirmPasswordResetRejectsReplay(t *testing.T) {
	service := newTestAuthService("hash", &fakeSessionStore{}, time.Date(2026, 5, 13, 9, 0, 0, 0, time.UTC))
	service.passwordResets = &fakePasswordResetStore{useErr: ErrPasswordResetTokenUsed}

	_, err := service.ConfirmPasswordReset(context.Background(), PasswordResetConfirmInput{
		Token:       "reset-token",
		NewPassword: "new-password",
		ClientIP:    "192.0.2.10",
	})
	if !errors.Is(err, ErrPasswordResetTokenUsed) {
		t.Fatalf("expected used token error, got %v", err)
	}
}

func TestAuthServiceConfirmPasswordResetRejectsRateLimitedAttempt(t *testing.T) {
	service := newTestAuthService("hash", &fakeSessionStore{}, time.Date(2026, 5, 13, 9, 0, 0, 0, time.UTC))
	service.passwordResets = &fakePasswordResetStore{}
	service.rateLimits = &fakeAuthRateLimitStore{blocked: true}

	_, err := service.ConfirmPasswordReset(context.Background(), PasswordResetConfirmInput{
		Token:       "reset-token",
		NewPassword: "new-password",
		ClientIP:    "192.0.2.10",
	})
	if !errors.Is(err, ErrAuthRateLimited) {
		t.Fatalf("expected rate limit error, got %v", err)
	}
}

type fakePasswordResetStore struct {
	created  CreatePasswordResetTokenInput
	token    PasswordResetToken
	useErr   error
	usedHash string
}

func (store *fakePasswordResetStore) CreatePasswordResetToken(ctx context.Context, input CreatePasswordResetTokenInput) (PasswordResetToken, error) {
	store.created = input
	store.token = PasswordResetToken{
		ID:        "reset_1",
		TenantID:  input.TenantID,
		UserID:    input.UserID,
		TokenHash: input.TokenHash,
		ExpiresAt: input.ExpiresAt,
	}
	return store.token, nil
}

func (store *fakePasswordResetStore) UsePasswordResetToken(ctx context.Context, tokenHash string, now time.Time) (PasswordResetToken, error) {
	store.usedHash = tokenHash
	if store.useErr != nil {
		return PasswordResetToken{}, store.useErr
	}
	if store.token.TokenHash != "" && store.token.TokenHash != tokenHash {
		return PasswordResetToken{}, ErrPasswordResetTokenInvalid
	}
	if !store.token.UsedAt.IsZero() {
		return PasswordResetToken{}, ErrPasswordResetTokenUsed
	}
	if !store.token.ExpiresAt.IsZero() && !store.token.ExpiresAt.After(now) {
		return PasswordResetToken{}, ErrPasswordResetTokenExpired
	}
	store.token.UsedAt = now
	return store.token, nil
}

type fakePasswordResetDelivery struct {
	input PasswordResetDeliveryInput
	err   error
}

func (delivery *fakePasswordResetDelivery) DeliverPasswordReset(ctx context.Context, input PasswordResetDeliveryInput) error {
	delivery.input = input
	return delivery.err
}
