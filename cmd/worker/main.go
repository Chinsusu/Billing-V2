package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/jobs"
	"github.com/Chinsusu/Billing-V2/internal/modules/order"
	"github.com/Chinsusu/Billing-V2/internal/modules/provider"
	platformdb "github.com/Chinsusu/Billing-V2/internal/platform/db"
)

const defaultWorkerCommand = "provision-once"

type workerConfig struct {
	DSN       string
	WorkerID  string
	Timeout   time.Duration
	BatchSize int
	LockFor   time.Duration
}

type provisionRunner interface {
	RunOnce(ctx context.Context) (jobs.RunSummary, error)
}

type provisionRunnerFactory func(ctx context.Context, cfg workerConfig) (provisionRunner, func() error, error)

type workerDependencies struct {
	stdout    io.Writer
	newRunner provisionRunnerFactory
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "worker failed: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	return runWithDependencies(args, workerDependencies{
		stdout:    os.Stdout,
		newRunner: newProvisioningRunner,
	})
}

func runWithDependencies(args []string, deps workerDependencies) error {
	if deps.stdout == nil {
		deps.stdout = io.Discard
	}
	if deps.newRunner == nil {
		deps.newRunner = newProvisioningRunner
	}

	command, flagArgs := splitCommand(args)
	cfg, remaining, err := parseConfig(flagArgs)
	if err != nil {
		return err
	}
	if len(remaining) > 0 {
		if command != defaultWorkerCommand || len(remaining) > 1 {
			return fmt.Errorf("unexpected argument %q", remaining[0])
		}
		command = remaining[0]
	}
	if command == "" {
		command = defaultWorkerCommand
	}
	switch command {
	case defaultWorkerCommand:
		return runProvisionOnce(cfg, deps)
	default:
		return fmt.Errorf("unknown command %q; use %s", command, defaultWorkerCommand)
	}
}

func splitCommand(args []string) (string, []string) {
	if len(args) == 0 {
		return defaultWorkerCommand, args
	}
	if !strings.HasPrefix(args[0], "-") {
		return args[0], args[1:]
	}
	return defaultWorkerCommand, args
}

func parseConfig(args []string) (workerConfig, []string, error) {
	flags := flag.NewFlagSet("worker", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	cfg := workerConfig{}
	flags.StringVar(&cfg.DSN, "dsn", os.Getenv("DB_DSN"), "PostgreSQL DSN")
	flags.StringVar(&cfg.WorkerID, "worker-id", envOrDefault("WORKER_ID", "local-provisioner-1"), "worker identifier")
	flags.DurationVar(&cfg.Timeout, "timeout", 60*time.Second, "worker command timeout")
	flags.IntVar(&cfg.BatchSize, "batch-size", 1, "maximum jobs to claim once")
	flags.DurationVar(&cfg.LockFor, "lock-for", time.Minute, "job lock duration")
	if err := flags.Parse(args); err != nil {
		return workerConfig{}, nil, err
	}
	cfg.DSN = strings.TrimSpace(cfg.DSN)
	cfg.WorkerID = strings.TrimSpace(cfg.WorkerID)
	return cfg, flags.Args(), nil
}

func runProvisionOnce(cfg workerConfig, deps workerDependencies) error {
	if err := guardLocalWorkerEnvironment(); err != nil {
		return err
	}
	if cfg.DSN == "" {
		return fmt.Errorf("DB_DSN or -dsn is required for %s", defaultWorkerCommand)
	}
	if cfg.WorkerID == "" {
		return fmt.Errorf("worker id is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
	defer cancel()

	runner, cleanup, err := deps.newRunner(ctx, cfg)
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
	writeSummary(deps.stdout, summary)
	return nil
}

func newProvisioningRunner(ctx context.Context, cfg workerConfig) (provisionRunner, func() error, error) {
	conn, err := platformdb.Open(ctx, platformdb.Config{
		DriverName: platformdb.DefaultDriverName,
		DSN:        cfg.DSN,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("open worker database: %w", err)
	}

	registry, err := provider.NewFakeRegistry()
	if err != nil {
		_ = conn.Close()
		return nil, nil, err
	}

	jobStore := jobs.NewPostgresStore(conn)
	orderStore := order.NewPostgresStore(conn)
	runner := order.NewProviderProvisioningRunner(jobStore, registry, orderStore, jobs.WorkerID(cfg.WorkerID))
	runner.BatchSize = cfg.BatchSize
	runner.LockFor = cfg.LockFor
	return runner, closeDatabase(conn), nil
}

func closeDatabase(conn *sql.DB) func() error {
	return func() error {
		return conn.Close()
	}
}

func writeSummary(w io.Writer, summary jobs.RunSummary) {
	fmt.Fprintf(
		w,
		"provision-once claimed=%d succeeded=%d retried=%d manual_review=%d terminal_failed=%d cancelled=%d\n",
		summary.Claimed,
		summary.Succeeded,
		summary.Retried,
		summary.ManualReview,
		summary.TerminalFailed,
		summary.Cancelled,
	)
}

func guardLocalWorkerEnvironment() error {
	switch strings.ToLower(strings.TrimSpace(os.Getenv("APP_ENV"))) {
	case "prod", "production":
		return fmt.Errorf("refusing to run worker with APP_ENV=%s", os.Getenv("APP_ENV"))
	default:
		return nil
	}
}

func envOrDefault(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
