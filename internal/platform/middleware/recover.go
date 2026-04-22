package middleware

import (
	"net/http"

	"github.com/Chinsusu/Billing-V2/internal/platform/httpserver"
	"github.com/Chinsusu/Billing-V2/internal/platform/logger"
)

func Recover(log *logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if recovered := recover(); recovered != nil {
					if log != nil {
						log.Error("http request panic recovered",
							logger.String("module", "http"),
							logger.String("operation", "recover"),
							logger.String("request_id", httpserver.RequestIDFromContext(r.Context())),
							logger.String("method", r.Method),
							logger.String("path", r.URL.Path),
						)
					}
					httpserver.WriteError(w, r, http.StatusInternalServerError, "internal.unexpected_error", "Unexpected server error.")
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
