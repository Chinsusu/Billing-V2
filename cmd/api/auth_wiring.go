package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/app"
	"github.com/Chinsusu/Billing-V2/internal/modules/audit"
	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/rbac"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
	"github.com/Chinsusu/Billing-V2/internal/platform/config"
	platformdb "github.com/Chinsusu/Billing-V2/internal/platform/db"
)

type authAuditAdapter struct {
	service *audit.Service
}

func (adapter authAuditAdapter) AppendAuthAudit(ctx context.Context, input identity.AuthAuditInput) error {
	if adapter.service == nil {
		return nil
	}
	_, err := adapter.service.Append(ctx, audit.AppendInput{
		TenantID:         input.TenantID,
		ActorID:          audit.ActorID(input.ActorID),
		ActorType:        audit.ActorTypeUser,
		Action:           input.Action,
		TargetType:       "user",
		TargetID:         audit.TargetID(input.TargetUserID),
		MetadataRedacted: json.RawMessage(`{"method":"totp"}`),
		CorrelationID:    audit.CorrelationID(input.CorrelationID),
	})
	return err
}

func newAuthService(executor platformdb.Executor, cfg config.Config) (*identity.AuthService, error) {
	var cipher identity.SecretCipher
	if cfg.EncryptionKey != "" {
		secretCipher, err := identity.NewAESGCMSecretCipher(cfg.EncryptionKey)
		if err != nil {
			return nil, fmt.Errorf("configure auth secret cipher: %w", err)
		}
		cipher = secretCipher
	}
	return identity.NewAuthService(identity.AuthServiceOptions{
		Tenants:          tenant.NewPostgresStore(executor),
		Users:            identity.NewPostgresUserStore(executor),
		Sessions:         identity.NewPostgresSessionStore(executor),
		TwoFactor:        identity.NewPostgresTwoFactorStore(executor),
		RateLimits:       identity.NewPostgresAuthRateLimitStore(executor),
		PasswordResets:   identity.NewPostgresPasswordResetStore(executor),
		ResetDelivery:    identity.NoopPasswordResetDelivery{},
		Roles:            rbac.NewPostgresStore(executor),
		Cipher:           cipher,
		Audit:            authAuditAdapter{service: audit.NewService(audit.NewPostgresStore(executor))},
		SessionTTL:       cfg.SessionTokenTTL,
		PasswordResetTTL: cfg.PasswordResetTTL,
		Now:              time.Now,
	}), nil
}

func newAuthRoutes(service *identity.AuthService, cfg config.Config) app.RouteRegistrar {
	return identity.NewAuthHTTPHandlerWithOptions(service, identity.AuthHTTPHandlerOptions{
		CookieName:             cfg.SessionCookieName,
		CookieSecure:           cfg.SessionCookieSecure,
		AllowLocalTenantHeader: cfg.AllowDevActorHeaders(),
	})
}
