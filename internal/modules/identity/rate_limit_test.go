package identity

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestAuthServiceLoginRejectsRateLimitedAttempt(t *testing.T) {
	sessions := &fakeSessionStore{}
	service := newTestAuthService("invalid-hash", sessions, time.Date(2026, 5, 13, 9, 0, 0, 0, time.UTC))
	service.rateLimits = &fakeAuthRateLimitStore{blocked: true}

	_, err := service.Login(context.Background(), LoginInput{
		Email:                  "admin@local.billing",
		Password:               "wrong",
		LocalTenantID:          "tenant_1",
		AllowLocalTenantHeader: true,
		ClientIP:               "192.0.2.10",
	})
	if !errors.Is(err, ErrAuthRateLimited) {
		t.Fatalf("expected rate limit error, got %v", err)
	}
	if sessions.created.TokenHash != "" {
		t.Fatalf("rate-limited login must not create a session: %+v", sessions.created)
	}
}

func TestAuthRateLimitHashesEmailAndIPKeys(t *testing.T) {
	store := &fakeAuthRateLimitStore{}
	service := newTestAuthService("invalid-hash", &fakeSessionStore{}, time.Date(2026, 5, 13, 9, 0, 0, 0, time.UTC))
	service.rateLimits = store

	err := service.enforceAuthRateLimit(context.Background(), AuthActionLogin, "tenant_1", "Client@Example.com", "192.0.2.10")
	if err != nil {
		t.Fatalf("enforceAuthRateLimit returned error: %v", err)
	}
	if len(store.inputs) != 3 {
		t.Fatalf("expected tenant/email, ip, and tenant/ip keys, got %d", len(store.inputs))
	}
	for _, input := range store.inputs {
		if len(input.KeyHash) != 64 || input.KeyHash == "Client@Example.com" || input.KeyHash == "192.0.2.10" {
			t.Fatalf("expected hashed limiter key, got %+v", input)
		}
	}
}
