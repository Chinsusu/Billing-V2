package app

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/platform/config"
	"github.com/Chinsusu/Billing-V2/internal/platform/httpserver"
	"github.com/Chinsusu/Billing-V2/internal/platform/logger"
)

type API struct {
	cfg     config.Config
	log     *logger.Logger
	handler http.Handler
}

type HealthResponse struct {
	Status      string `json:"status"`
	Service     string `json:"service"`
	Environment string `json:"environment"`
}

func NewAPI(cfg config.Config, log *logger.Logger) (*API, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	if log == nil {
		log = logger.New(io.Discard, cfg.LogLevel)
	}

	mux := http.NewServeMux()
	api := &API{
		cfg: cfg,
		log: log,
	}
	mux.HandleFunc("/healthz", api.handleHealth)
	mux.HandleFunc("/readyz", api.handleReady)
	api.handler = httpserver.WithRequestID(mux)
	return api, nil
}

func (api *API) Handler() http.Handler {
	return api.handler
}

func (api *API) Run(ctx context.Context) error {
	server, err := httpserver.New(httpserver.Config{
		Addr:              api.cfg.HTTPAddr,
		Handler:           api.handler,
		ReadHeaderTimeout: 5 * time.Second,
		ShutdownTimeout:   10 * time.Second,
	})
	if err != nil {
		return err
	}

	api.log.Info("api server starting",
		logger.String("module", "api"),
		logger.String("addr", api.cfg.HTTPAddr),
		logger.String("environment", string(api.cfg.AppEnvironment)),
	)
	return server.Start(ctx)
}

func (api *API) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httpserver.WriteError(w, r, http.StatusMethodNotAllowed, "request.method_not_allowed", "Method is not allowed.")
		return
	}

	httpserver.WriteSuccess(w, r, http.StatusOK, HealthResponse{
		Status:      "ok",
		Service:     api.cfg.AppName,
		Environment: string(api.cfg.AppEnvironment),
	})
}

func (api *API) handleReady(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		httpserver.WriteError(w, r, http.StatusMethodNotAllowed, "request.method_not_allowed", "Method is not allowed.")
		return
	}

	httpserver.WriteSuccess(w, r, http.StatusOK, HealthResponse{
		Status:      "ready",
		Service:     api.cfg.AppName,
		Environment: string(api.cfg.AppEnvironment),
	})
}
