package identity

import (
	"context"
	"errors"
	"net/http"

	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
	"github.com/Chinsusu/Billing-V2/internal/platform/httpserver"
)

type SessionResolver interface {
	ResolveSession(ctx context.Context, token string) (SessionIdentity, error)
}

type SessionMiddlewareOptions struct {
	CookieName string
	Resolver   SessionResolver
}

func SessionMiddleware(options SessionMiddlewareOptions) func(http.Handler) http.Handler {
	cookieName := options.CookieName
	if cookieName == "" {
		cookieName = "billing_session"
	}
	return func(next http.Handler) http.Handler {
		if next == nil {
			next = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
		}
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if isAuthRoute(r.URL.Path) {
				next.ServeHTTP(w, r)
				return
			}
			cookie, err := r.Cookie(cookieName)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}
			if options.Resolver == nil {
				httpserver.WriteError(w, r, http.StatusServiceUnavailable, "auth.session_unavailable", "Session service is unavailable.")
				return
			}
			identity, err := options.Resolver.ResolveSession(r.Context(), cookie.Value)
			if err != nil {
				writeSessionMiddlewareError(w, r, err)
				return
			}
			actor := identity.Actor()
			ctx := WithActor(r.Context(), actor)
			if _, ok := tenant.FromContext(ctx); !ok {
				ctx = tenant.WithContext(ctx, tenant.NewContext(identity.User.TenantID))
			}
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func isAuthRoute(path string) bool {
	return path == "/auth/login" || path == "/auth/logout"
}

func writeSessionMiddlewareError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, ErrSessionTokenMissing), errors.Is(err, ErrSessionInvalid), errors.Is(err, ErrSessionExpired):
		httpserver.WriteError(w, r, http.StatusUnauthorized, "auth.session_invalid", "Session is invalid or expired.")
	case errors.Is(err, ErrAuthStoreMissing), errors.Is(err, ErrSessionStoreExecutorMissing):
		httpserver.WriteError(w, r, http.StatusServiceUnavailable, "auth.session_unavailable", "Session service is unavailable.")
	default:
		httpserver.WriteError(w, r, http.StatusUnauthorized, "auth.session_invalid", "Session is invalid or expired.")
	}
}
