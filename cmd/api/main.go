package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Chinsusu/Billing-V2/internal/app"
	"github.com/Chinsusu/Billing-V2/internal/platform/config"
	"github.com/Chinsusu/Billing-V2/internal/platform/logger"
)

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
	api, err := app.NewAPI(cfg, log)
	if err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	return api.Run(ctx)
}
