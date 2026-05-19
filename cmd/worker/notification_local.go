package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/Chinsusu/Billing-V2/internal/modules/jobs"
	"github.com/Chinsusu/Billing-V2/internal/modules/notification"
	platformdb "github.com/Chinsusu/Billing-V2/internal/platform/db"
)

func runNotificationLocalOnce(cfg workerConfig, deps workerDependencies) error {
	if err := validateNotificationLocalWorkerConfig(cfg, commandNotificationLocalOnce); err != nil {
		return err
	}

	ctx, cancel := workerCommandContext(cfg.Timeout)
	defer cancel()

	runner, cleanup, err := deps.newNotificationRunner(ctx, cfg)
	if err != nil {
		return err
	}
	if cleanup != nil {
		defer cleanup()
	}

	summary, err := runner.RunOnce(ctx)
	if err != nil {
		return err
	}
	writeSummary(deps.stdout, commandNotificationLocalOnce, summary)
	return nil
}

func runNotificationLocalLoop(cfg workerConfig, deps workerDependencies) error {
	if err := validateNotificationLocalWorkerConfig(cfg, commandNotificationLocalLoop); err != nil {
		return err
	}
	if cfg.Interval <= 0 {
		return fmt.Errorf("worker interval must be positive")
	}

	ctx, cancel := workerCommandContext(cfg.Timeout)
	defer cancel()

	runner, cleanup, err := deps.newNotificationRunner(ctx, cfg)
	if err != nil {
		return err
	}
	if cleanup != nil {
		defer cleanup()
	}

	for pass := 1; ; pass++ {
		if err := ctx.Err(); err != nil {
			return nil
		}
		summary, err := runner.RunOnce(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			return err
		}
		writeLoopSummary(deps.stdout, commandNotificationLocalLoop, pass, summary)
		if summary.Claimed == 0 {
			if err := waitWorkerInterval(ctx, cfg.Interval); err != nil {
				return nil
			}
		}
	}
}

func validateNotificationLocalWorkerConfig(cfg workerConfig, command string) error {
	if err := guardNotificationLocalEnvironment(); err != nil {
		return err
	}
	if cfg.DSN == "" {
		return fmt.Errorf("DB_DSN or -dsn is required for %s", command)
	}
	if cfg.WorkerID == "" {
		return fmt.Errorf("worker id is required")
	}
	return nil
}

func newNotificationLocalRunner(ctx context.Context, cfg workerConfig) (provisionRunner, func() error, error) {
	conn, err := platformdb.Open(ctx, platformdb.Config{
		DriverName: platformdb.DefaultDriverName,
		DSN:        cfg.DSN,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("open worker database: %w", err)
	}

	jobStore := jobs.NewPostgresStore(conn)
	deliveryStore := notification.NewPostgresStore(conn)
	runner := notification.NewLocalDeliveryRunner(jobStore, deliveryStore, jobs.WorkerID(cfg.WorkerID))
	runner.BatchSize = cfg.BatchSize
	runner.LockFor = cfg.LockFor
	return runner, closeDatabase(conn), nil
}

func guardNotificationLocalEnvironment() error {
	switch strings.ToLower(strings.TrimSpace(os.Getenv("APP_ENV"))) {
	case "local", "dev":
		return nil
	default:
		return fmt.Errorf("refusing to run local notification worker with APP_ENV=%s; use APP_ENV=local or APP_ENV=dev", os.Getenv("APP_ENV"))
	}
}
