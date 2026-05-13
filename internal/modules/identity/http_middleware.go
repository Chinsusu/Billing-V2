package identity

import (
	"net/http"
	"os"
	"strings"

	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

const (
	HeaderActorID       = "X-Actor-Id"
	HeaderActorType     = "X-Actor-Type"
	HeaderActorTenantID = "X-Actor-Tenant-Id"
	HeaderActorRoleIDs  = "X-Actor-Role-Ids"
)

type HeaderActorOptions struct {
	ActorIDHeader       string
	ActorTypeHeader     string
	ActorTenantIDHeader string
	ActorRoleIDsHeader  string
}

func HeaderActorMiddleware(next http.Handler) http.Handler {
	return HeaderActorMiddlewareWithOptions(HeaderActorOptions{}, next)
}

func HeaderActorMiddlewareWithOptions(options HeaderActorOptions, next http.Handler) http.Handler {
	if next == nil {
		next = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	}
	headers := normalizeActorHeaders(options)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !headerActorAuthEnabled() {
			next.ServeHTTP(w, r)
			return
		}
		actorID := UserID(strings.TrimSpace(r.Header.Get(headers.ActorIDHeader)))
		if actorID == "" {
			next.ServeHTTP(w, r)
			return
		}
		actorType := ActorType(strings.TrimSpace(r.Header.Get(headers.ActorTypeHeader)))
		actorTenantID := tenant.ID(strings.TrimSpace(r.Header.Get(headers.ActorTenantIDHeader)))
		if actorTenantID.Empty() {
			if tenantContext, ok := tenant.FromContext(r.Context()); ok {
				actorTenantID = tenantContext.ActorTenantID
			}
		}
		actor := NewActor(actorID, actorTenantID, actorType, parseRoleIDs(r.Header.Get(headers.ActorRoleIDsHeader))...)
		next.ServeHTTP(w, r.WithContext(WithActor(r.Context(), actor)))
	})
}

func headerActorAuthEnabled() bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv("APP_ENV"))) {
	case "", "local", "dev":
		return true
	default:
		return false
	}
}

func normalizeActorHeaders(options HeaderActorOptions) HeaderActorOptions {
	if options.ActorIDHeader == "" {
		options.ActorIDHeader = HeaderActorID
	}
	if options.ActorTypeHeader == "" {
		options.ActorTypeHeader = HeaderActorType
	}
	if options.ActorTenantIDHeader == "" {
		options.ActorTenantIDHeader = HeaderActorTenantID
	}
	if options.ActorRoleIDsHeader == "" {
		options.ActorRoleIDsHeader = HeaderActorRoleIDs
	}
	return options
}

func parseRoleIDs(value string) []RoleID {
	parts := strings.Split(value, ",")
	roleIDs := make([]RoleID, 0, len(parts))
	for _, part := range parts {
		roleID := RoleID(strings.TrimSpace(part))
		if roleID == "" {
			continue
		}
		roleIDs = append(roleIDs, roleID)
	}
	return roleIDs
}
