package middleware

import (
	"net/http"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/platform/httpserver"
	"github.com/Chinsusu/Billing-V2/internal/platform/logger"
)

func RequestLogger(log *logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startedAt := time.Now()
			recorder := &statusRecorder{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			next.ServeHTTP(recorder, r)

			if log == nil {
				return
			}
			log.Info("http request completed",
				logger.String("module", "http"),
				logger.String("operation", "request"),
				logger.String("request_id", httpserver.RequestIDFromContext(r.Context())),
				logger.String("method", r.Method),
				logger.String("path", r.URL.Path),
				logger.Int("status", recorder.statusCode),
				logger.Int("duration_ms", int(time.Since(startedAt).Milliseconds())),
			)
		})
	}
}

type statusRecorder struct {
	http.ResponseWriter
	statusCode int
	wrote      bool
}

func (recorder *statusRecorder) WriteHeader(statusCode int) {
	if recorder.wrote {
		return
	}
	recorder.statusCode = statusCode
	recorder.wrote = true
	recorder.ResponseWriter.WriteHeader(statusCode)
}

func (recorder *statusRecorder) Write(body []byte) (int, error) {
	if !recorder.wrote {
		recorder.WriteHeader(http.StatusOK)
	}
	return recorder.ResponseWriter.Write(body)
}
