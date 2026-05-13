package main

import (
	"time"

	"github.com/Chinsusu/Billing-V2/internal/app"
	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/rbac"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
	"github.com/Chinsusu/Billing-V2/internal/platform/config"
	platformdb "github.com/Chinsusu/Billing-V2/internal/platform/db"
)

func newAuthService(executor platformdb.Executor, cfg config.Config) *identity.AuthService {
	return identity.NewAuthService(identity.AuthServiceOptions{
		Tenants:    tenant.NewPostgresStore(executor),
		Users:      identity.NewPostgresUserStore(executor),
		Sessions:   identity.NewPostgresSessionStore(executor),
		Roles:      rbac.NewPostgresStore(executor),
		SessionTTL: cfg.SessionTokenTTL,
		Now:        time.Now,
	})
}

func newAuthRoutes(service *identity.AuthService, cfg config.Config) app.RouteRegistrar {
	return identity.NewAuthHTTPHandlerWithOptions(service, identity.AuthHTTPHandlerOptions{
		CookieName:             cfg.SessionCookieName,
		CookieSecure:           cfg.SessionCookieSecure,
		AllowLocalTenantHeader: cfg.AllowDevActorHeaders(),
	})
}
