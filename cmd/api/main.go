package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Chinsusu/Billing-V2/internal/app"
	"github.com/Chinsusu/Billing-V2/internal/modules/catalog"
	"github.com/Chinsusu/Billing-V2/internal/platform/config"
	platformdb "github.com/Chinsusu/Billing-V2/internal/platform/db"
	"github.com/Chinsusu/Billing-V2/internal/platform/logger"
)

type databaseOpener func(ctx context.Context, cfg platformdb.Config) (*sql.DB, error)

type apiRuntime struct {
	api     *app.API
	cleanup func() error
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "api exited: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.LoadFromEnv()
	if err != nil {
		return err
	}

	log := logger.New(os.Stdout, cfg.LogLevel)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	runtime, err := newRuntime(ctx, cfg, log, platformdb.Open)
	if err != nil {
		return err
	}
	defer func() {
		if err := runtime.close(); err != nil {
			log.Error("api cleanup failed", logger.String("module", "api"), logger.String("error", err.Error()))
		}
	}()

	return runtime.api.Run(ctx)
}

func newRuntime(ctx context.Context, cfg config.Config, log *logger.Logger, openDatabase databaseOpener) (*apiRuntime, error) {
	options := app.APIOptions{}
	cleanup := func() error { return nil }

	if cfg.DatabaseDSN != "" {
		conn, err := openDatabase(ctx, platformdb.Config{
			DriverName: platformdb.DefaultDriverName,
			DSN:        cfg.DatabaseDSN,
		})
		if err != nil {
			return nil, fmt.Errorf("open api database: %w", err)
		}
		cleanup = conn.Close
		options.CatalogRoutes = newCatalogRoutes(conn)
	}

	api, err := app.NewAPIWithOptions(cfg, log, options)
	if err != nil {
		_ = cleanup()
		return nil, err
	}
	return &apiRuntime{
		api:     api,
		cleanup: cleanup,
	}, nil
}

func newCatalogRoutes(executor platformdb.Executor) app.RouteRegistrar {
	store := catalog.NewPostgresStore(executor)
	service := catalog.NewService(store)
	return catalog.NewHTTPHandler(service)
}

func (runtime *apiRuntime) close() error {
	if runtime == nil || runtime.cleanup == nil {
		return nil
	}
	return runtime.cleanup()
}
