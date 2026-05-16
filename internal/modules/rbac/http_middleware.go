package rbac

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
	"github.com/Chinsusu/Billing-V2/internal/platform/httpserver"
)

const HeaderAccessReason = "X-Access-Reason"

type RouteMiddleware func(http.HandlerFunc) http.HandlerFunc

type PermissionMiddlewareOptions struct {
	Authorizer        Authorizer
	Permission        Permission
	Risk              RiskLevel
	AllowedActorTypes []identity.ActorType
	ReasonHeader      string
	ResourceTenantID  func(r *http.Request) tenant.ID
}

func RequirePermission(authorizer Authorizer, permission Permission, risk RiskLevel) RouteMiddleware {
	return RequirePermissionWithOptions(PermissionMiddlewareOptions{
		Authorizer: authorizer,
		Permission: permission,
		Risk:       risk,
	})
}

func RequirePermissionWithOptions(options PermissionMiddlewareOptions) RouteMiddleware {
	reasonHeader := options.ReasonHeader
	if reasonHeader == "" {
		reasonHeader = HeaderAccessReason
	}
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if options.Authorizer == nil {
				httpserver.WriteError(w, r, http.StatusServiceUnavailable, "auth.authorizer_unavailable", "Authorization service is unavailable.")
				return
			}
			actor, err := identity.RequireActor(r.Context())
			if err != nil {
				writeAuthorizationError(w, r, err)
				return
			}
			if err := actor.Validate(); err != nil {
				writeAuthorizationError(w, r, err)
				return
			}
			if !actorTypeAllowed(actor.Type, options.AllowedActorTypes) {
				writeAuthorizationError(w, r, ErrPermissionDenied)
				return
			}
			resourceTenantID := tenant.ID("")
			if options.ResourceTenantID != nil {
				resourceTenantID = options.ResourceTenantID(r)
			} else if tenantContext, ok := tenant.FromContext(r.Context()); ok {
				resourceTenantID = tenantContext.EffectiveTenantID
			}
			err = options.Authorizer.Check(r.Context(), CheckRequest{
				Actor:            actor,
				Permission:       options.Permission,
				ResourceTenantID: resourceTenantID,
				Risk:             options.Risk,
				Reason:           strings.TrimSpace(r.Header.Get(reasonHeader)),
			})
			if err != nil {
				writeAuthorizationError(w, r, err)
				return
			}
			next(w, r)
		}
	}
}

func actorTypeAllowed(actorType identity.ActorType, allowed []identity.ActorType) bool {
	if len(allowed) == 0 {
		return true
	}
	for _, current := range allowed {
		if actorType == current {
			return true
		}
	}
	return false
}

func writeAuthorizationError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, identity.ErrActorContextMissing),
		errors.Is(err, identity.ErrActorIDMissing),
		errors.Is(err, identity.ErrActorTypeMissing),
		errors.Is(err, identity.ErrActorTenantMissing):
		httpserver.WriteError(w, r, http.StatusUnauthorized, "auth.actor_required", "Actor context is required.")
	case errors.Is(err, tenant.ErrContextMissing), errors.Is(err, tenant.ErrTenantIDMissing):
		httpserver.WriteError(w, r, http.StatusBadRequest, "tenant.context_missing", "Tenant context is required.")
	case errors.Is(err, tenant.ErrTenantMismatch):
		httpserver.WriteError(w, r, http.StatusForbidden, "tenant.context_mismatch", "Tenant context does not match actor.")
	case errors.Is(err, tenant.ErrAccessDenied), errors.Is(err, ErrPermissionDenied):
		httpserver.WriteError(w, r, http.StatusForbidden, "auth.permission_denied", "Permission denied.")
	case errors.Is(err, tenant.ErrEmergencyReasonMissing):
		httpserver.WriteError(w, r, http.StatusForbidden, "auth.reason_required", "Access reason is required.")
	case errors.Is(err, ErrPermissionMissing):
		httpserver.WriteError(w, r, http.StatusInternalServerError, "auth.permission_missing", "Permission check is not configured.")
	case errors.Is(err, ErrRBACStoreExecutorMissing):
		httpserver.WriteError(w, r, http.StatusServiceUnavailable, "auth.authorizer_unavailable", "Authorization service is unavailable.")
	default:
		httpserver.WriteError(w, r, http.StatusInternalServerError, "auth.authorization_failed", "Authorization failed.")
	}
}
