package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/audit"
	"github.com/Chinsusu/Billing-V2/internal/modules/jobs"
	"github.com/Chinsusu/Billing-V2/internal/modules/order"
	"github.com/Chinsusu/Billing-V2/internal/modules/provider"
	platformdb "github.com/Chinsusu/Billing-V2/internal/platform/db"
)

const (
	commandProvisionOnce         = "provision-once"
	commandProvisionLoop         = "provision-loop"
	commandLifecycleScheduleOnce = "lifecycle-schedule-once"
	commandLifecycleOnce         = "lifecycle-once"
	commandLifecycleLoop         = "lifecycle-loop"
	defaultWorkerCommand         = commandProvisionOnce
)

type workerConfig struct {
	DSN         string
	WorkerID    string
	Timeout     time.Duration
	BatchSize   int
	LockFor     time.Duration
	Interval    time.Duration
	GracePeriod time.Duration
}

type provisionRunner interface {
	RunOnce(ctx context.Context) (jobs.RunSummary, error)
}

type provisionRunnerFactory func(ctx context.Context, cfg workerConfig) (provisionRunner, func() error, error)

type lifecycleScheduler interface {
	ScheduleDue(ctx context.Context, input order.ListDueServiceLifecycleActionsInput) (order.ServiceLifecycleScheduleSummary, error)
}

type lifecycleSchedulerFactory func(ctx context.Context, cfg workerConfig) (lifecycleScheduler, func() error, error)

type workerDependencies struct {
	stdout                io.Writer
	newRunner             provisionRunnerFactory
	newLifecycleRunner    provisionRunnerFactory
	newLifecycleScheduler lifecycleSchedulerFactory
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
	if deps.newLifecycleRunner == nil {
		deps.newLifecycleRunner = newLifecycleRunner
	}
	if deps.newLifecycleScheduler == nil {
		deps.newLifecycleScheduler = newLifecycleScheduler
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
	case commandProvisionOnce:
		return runProvisionOnce(cfg, deps)
	case commandProvisionLoop:
		return runProvisionLoop(cfg, deps)
	case commandLifecycleScheduleOnce:
		return runLifecycleScheduleOnce(cfg, deps)
	case commandLifecycleOnce:
		return runLifecycleOnce(cfg, deps)
	case commandLifecycleLoop:
		return runLifecycleLoop(cfg, deps)
	default:
		return fmt.Errorf(
			"unknown command %q; use %s, %s, %s, %s, or %s",
			command,
			commandProvisionOnce,
			commandProvisionLoop,
			commandLifecycleScheduleOnce,
			commandLifecycleOnce,
			commandLifecycleLoop,
		)
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
	flags.DurationVar(&cfg.Interval, "interval", 5*time.Second, "idle wait interval for provision-loop")
	flags.DurationVar(&cfg.GracePeriod, "grace-period", order.DefaultServiceLifecycleGracePeriod, "service lifecycle grace period before termination")
	if err := flags.Parse(args); err != nil {
		return workerConfig{}, nil, err
	}
	cfg.DSN = strings.TrimSpace(cfg.DSN)
	cfg.WorkerID = strings.TrimSpace(cfg.WorkerID)
	return cfg, flags.Args(), nil
}

func runProvisionOnce(cfg workerConfig, deps workerDependencies) error {
	if err := validateWorkerConfig(cfg, commandProvisionOnce); err != nil {
		return err
	}

	ctx, cancel := workerCommandContext(cfg.Timeout)
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
	writeSummary(deps.stdout, commandProvisionOnce, summary)
	return nil
}

func runProvisionLoop(cfg workerConfig, deps workerDependencies) error {
	if err := validateWorkerConfig(cfg, commandProvisionLoop); err != nil {
		return err
	}
	if cfg.Interval <= 0 {
		return fmt.Errorf("worker interval must be positive")
	}

	ctx, cancel := workerCommandContext(cfg.Timeout)
	defer cancel()

	runner, cleanup, err := deps.newRunner(ctx, cfg)
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
		writeLoopSummary(deps.stdout, commandProvisionLoop, pass, summary)
		if summary.Claimed == 0 {
			if err := waitWorkerInterval(ctx, cfg.Interval); err != nil {
				return nil
			}
		}
	}
}

func runLifecycleScheduleOnce(cfg workerConfig, deps workerDependencies) error {
	if err := validateWorkerConfig(cfg, commandLifecycleScheduleOnce); err != nil {
		return err
	}

	ctx, cancel := workerCommandContext(cfg.Timeout)
	defer cancel()

	scheduler, cleanup, err := deps.newLifecycleScheduler(ctx, cfg)
	if err != nil {
		return err
	}
	if cleanup != nil {
		defer cleanup()
	}

	summary, err := scheduler.ScheduleDue(ctx, order.ListDueServiceLifecycleActionsInput{
		Limit:       cfg.BatchSize,
		GracePeriod: cfg.GracePeriod,
	})
	if err != nil {
		return err
	}
	writeLifecycleScheduleSummary(deps.stdout, summary)
	return nil
}

func runLifecycleOnce(cfg workerConfig, deps workerDependencies) error {
	if err := validateWorkerConfig(cfg, commandLifecycleOnce); err != nil {
		return err
	}

	ctx, cancel := workerCommandContext(cfg.Timeout)
	defer cancel()

	runner, cleanup, err := deps.newLifecycleRunner(ctx, cfg)
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
	writeSummary(deps.stdout, commandLifecycleOnce, summary)
	return nil
}

func runLifecycleLoop(cfg workerConfig, deps workerDependencies) error {
	if err := validateWorkerConfig(cfg, commandLifecycleLoop); err != nil {
		return err
	}
	if cfg.Interval <= 0 {
		return fmt.Errorf("worker interval must be positive")
	}

	ctx, cancel := workerCommandContext(cfg.Timeout)
	defer cancel()

	runner, cleanup, err := deps.newLifecycleRunner(ctx, cfg)
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
		writeLoopSummary(deps.stdout, commandLifecycleLoop, pass, summary)
		if summary.Claimed == 0 {
			if err := waitWorkerInterval(ctx, cfg.Interval); err != nil {
				return nil
			}
		}
	}
}

func validateWorkerConfig(cfg workerConfig, command string) error {
	if err := guardLocalWorkerEnvironment(); err != nil {
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

func workerCommandContext(timeout time.Duration) (context.Context, context.CancelFunc) {
	signalCtx, stopSignals := signal.NotifyContext(context.Background(), os.Interrupt)
	ctx, cancelTimeout := context.WithTimeout(signalCtx, timeout)
	return ctx, func() {
		cancelTimeout()
		stopSignals()
	}
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

func newLifecycleRunner(ctx context.Context, cfg workerConfig) (provisionRunner, func() error, error) {
	conn, err := platformdb.Open(ctx, platformdb.Config{
		DriverName: platformdb.DefaultDriverName,
		DSN:        cfg.DSN,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("open worker database: %w", err)
	}

	jobStore := jobs.NewPostgresStore(conn)
	orderStore := order.NewPostgresStore(conn)
	auditService := audit.NewService(audit.NewPostgresStore(conn))
	orderService := order.NewServiceWithOptions(order.ServiceOptions{Store: orderStore, Audit: auditService})
	runner := order.NewServiceLifecycleRunner(jobStore, orderService, jobs.WorkerID(cfg.WorkerID))
	runner.BatchSize = cfg.BatchSize
	runner.LockFor = cfg.LockFor
	return runner, closeDatabase(conn), nil
}

func newLifecycleScheduler(ctx context.Context, cfg workerConfig) (lifecycleScheduler, func() error, error) {
	conn, err := platformdb.Open(ctx, platformdb.Config{
		DriverName: platformdb.DefaultDriverName,
		DSN:        cfg.DSN,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("open worker database: %w", err)
	}

	jobStore := jobs.NewPostgresStore(conn)
	orderStore := order.NewPostgresStore(conn)
	scheduler := order.NewServiceLifecycleScheduler(orderStore, jobStore)
	scheduler.Limit = cfg.BatchSize
	scheduler.GracePeriod = cfg.GracePeriod
	return scheduler, closeDatabase(conn), nil
}

func closeDatabase(conn *sql.DB) func() error {
	return func() error {
		return conn.Close()
	}
}

func writeSummary(w io.Writer, label string, summary jobs.RunSummary) {
	fmt.Fprintf(
		w,
		"%s claimed=%d succeeded=%d retried=%d manual_review=%d terminal_failed=%d cancelled=%d\n",
		label,
		summary.Claimed,
		summary.Succeeded,
		summary.Retried,
		summary.ManualReview,
		summary.TerminalFailed,
		summary.Cancelled,
	)
}

func writeLoopSummary(w io.Writer, label string, pass int, summary jobs.RunSummary) {
	fmt.Fprintf(
		w,
		"%s pass=%d claimed=%d succeeded=%d retried=%d manual_review=%d terminal_failed=%d cancelled=%d\n",
		label,
		pass,
		summary.Claimed,
		summary.Succeeded,
		summary.Retried,
		summary.ManualReview,
		summary.TerminalFailed,
		summary.Cancelled,
	)
}

func writeLifecycleScheduleSummary(w io.Writer, summary order.ServiceLifecycleScheduleSummary) {
	fmt.Fprintf(w, "lifecycle-schedule-once due=%d scheduled=%d\n", summary.Due, summary.Scheduled)
}

func waitWorkerInterval(ctx context.Context, interval time.Duration) error {
	timer := time.NewTimer(interval)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
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
