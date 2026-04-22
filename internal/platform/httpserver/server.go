package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type Config struct {
	Addr              string
	Handler           http.Handler
	ReadHeaderTimeout time.Duration
	ShutdownTimeout   time.Duration
}

type Server struct {
	server          *http.Server
	shutdownTimeout time.Duration
}

func New(cfg Config) (*Server, error) {
	if cfg.Addr == "" {
		return nil, fmt.Errorf("http server address is required")
	}
	if cfg.Handler == nil {
		return nil, fmt.Errorf("http server handler is required")
	}
	if cfg.ReadHeaderTimeout == 0 {
		cfg.ReadHeaderTimeout = 5 * time.Second
	}
	if cfg.ShutdownTimeout == 0 {
		cfg.ShutdownTimeout = 10 * time.Second
	}

	return &Server{
		server: &http.Server{
			Addr:              cfg.Addr,
			Handler:           cfg.Handler,
			ReadHeaderTimeout: cfg.ReadHeaderTimeout,
		},
		shutdownTimeout: cfg.ShutdownTimeout,
	}, nil
}

func (srv *Server) Start(ctx context.Context) error {
	errCh := make(chan error, 1)
	go func() {
		if err := srv.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
			return
		}
		errCh <- nil
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), srv.shutdownTimeout)
		defer cancel()
		return srv.server.Shutdown(shutdownCtx)
	case err := <-errCh:
		return err
	}
}
