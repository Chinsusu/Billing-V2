package tenant

import (
	"net/http"
	"strings"
)

const HeaderTenantID = "X-Tenant-Id"

type HeaderContextOptions struct {
	TenantHeader string
}

func HeaderContextMiddleware(next http.Handler) http.Handler {
	return HeaderContextMiddlewareWithOptions(HeaderContextOptions{}, next)
}

func HeaderContextMiddlewareWithOptions(options HeaderContextOptions, next http.Handler) http.Handler {
	if next == nil {
		next = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	}
	tenantHeader := options.TenantHeader
	if tenantHeader == "" {
		tenantHeader = HeaderTenantID
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tenantID := ID(strings.TrimSpace(r.Header.Get(tenantHeader)))
		if tenantID.Empty() {
			next.ServeHTTP(w, r)
			return
		}
		next.ServeHTTP(w, r.WithContext(WithContext(r.Context(), NewContext(tenantID))))
	})
}
