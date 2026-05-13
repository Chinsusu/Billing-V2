package identity

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

const (
	AuthActionLogin         = "auth.login"
	AuthActionPasswordReset = "auth.password_reset"
)

var (
	ErrAuthRateLimited       = errors.New("auth rate limited")
	ErrAuthRateLimitMissing  = errors.New("auth rate limit missing")
	ErrAuthRateLimitKeyEmpty = errors.New("auth rate limit key empty")
)

type AuthRateLimitPolicy struct {
	Limit  int
	Window time.Duration
}

type AuthRateLimitIncrementInput struct {
	Action      string
	KeyHash     string
	WindowStart time.Time
}

type AuthRateLimitCounter struct {
	Action       string
	KeyHash      string
	WindowStart  time.Time
	AttemptCount int
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type AuthRateLimitStore interface {
	IncrementAuthRateLimit(ctx context.Context, input AuthRateLimitIncrementInput) (AuthRateLimitCounter, error)
}

func (service *AuthService) enforceAuthRateLimit(ctx context.Context, action string, tenantID tenant.ID, email string, clientIP string) error {
	if service == nil || service.rateLimits == nil {
		return ErrAuthStoreMissing
	}
	policy, ok := authRateLimitPolicy(action)
	if !ok {
		return ErrAuthRateLimitMissing
	}
	keys := authRateLimitKeys(action, tenantID, email, clientIP)
	if len(keys) == 0 {
		return ErrAuthRateLimitKeyEmpty
	}
	windowStart := service.now().UTC().Truncate(policy.Window)
	for _, key := range keys {
		counter, err := service.rateLimits.IncrementAuthRateLimit(ctx, AuthRateLimitIncrementInput{
			Action:      action,
			KeyHash:     hashAuthRateLimitKey(key),
			WindowStart: windowStart,
		})
		if err != nil {
			return err
		}
		if counter.AttemptCount > policy.Limit {
			return ErrAuthRateLimited
		}
	}
	return nil
}

func authRateLimitPolicy(action string) (AuthRateLimitPolicy, bool) {
	switch action {
	case AuthActionLogin:
		return AuthRateLimitPolicy{Limit: 5, Window: 15 * time.Minute}, true
	case AuthActionPasswordReset:
		return AuthRateLimitPolicy{Limit: 3, Window: time.Hour}, true
	default:
		return AuthRateLimitPolicy{}, false
	}
}

func authRateLimitKeys(action string, tenantID tenant.ID, email string, clientIP string) []string {
	email = strings.ToLower(strings.TrimSpace(email))
	clientIP = strings.TrimSpace(clientIP)
	keys := []string{}
	if !tenantID.Empty() && email != "" {
		keys = append(keys, action+"|tenant:"+string(tenantID)+"|email:"+email)
	}
	if clientIP != "" {
		keys = append(keys, action+"|ip:"+clientIP)
	}
	if !tenantID.Empty() && clientIP != "" {
		keys = append(keys, action+"|tenant:"+string(tenantID)+"|ip:"+clientIP)
	}
	return keys
}

func hashAuthRateLimitKey(value string) string {
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:])
}
