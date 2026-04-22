package middleware

import (
	"net/http"

	"github.com/Chinsusu/Billing-V2/internal/platform/httpserver"
)

func RequireMethod(method string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			w.Header().Set("Allow", method)
			httpserver.WriteError(w, r, http.StatusMethodNotAllowed, "request.method_not_allowed", "Method is not allowed.")
			return
		}
		next(w, r)
	}
}
