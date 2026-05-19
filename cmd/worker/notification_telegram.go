package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/jobs"
	"github.com/Chinsusu/Billing-V2/internal/modules/notification"
	platformdb "github.com/Chinsusu/Billing-V2/internal/platform/db"
)

func runNotificationTelegramOnce(cfg workerConfig, deps workerDependencies) error {
	if err := validateNotificationTelegramWorkerConfig(cfg, commandNotificationTelegramOnce); err != nil {
		return err
	}

	ctx, cancel := workerCommandContext(cfg.Timeout)
	defer cancel()

	runner, cleanup, err := deps.newTelegramRunner(ctx, cfg)
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
	writeSummary(deps.stdout, commandNotificationTelegramOnce, summary)
	return nil
}

func runNotificationTelegramLoop(cfg workerConfig, deps workerDependencies) error {
	if err := validateNotificationTelegramWorkerConfig(cfg, commandNotificationTelegramLoop); err != nil {
		return err
	}
	if cfg.Interval <= 0 {
		return fmt.Errorf("worker interval must be positive")
	}

	ctx, cancel := workerCommandContext(cfg.Timeout)
	defer cancel()

	runner, cleanup, err := deps.newTelegramRunner(ctx, cfg)
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
		writeLoopSummary(deps.stdout, commandNotificationTelegramLoop, pass, summary)
		if summary.Claimed == 0 {
			if err := waitWorkerInterval(ctx, cfg.Interval); err != nil {
				return nil
			}
		}
	}
}

func runNotificationTelegramPreflight(cfg workerConfig, deps workerDependencies) error {
	if err := validateNotificationTelegramEnvironment(); err != nil {
		return err
	}
	telegramCfg := readTelegramConfigFromEnv()
	client := &http.Client{Timeout: cfg.Timeout}
	sender, err := notification.NewTelegramSender(telegramCfg, client)
	if err != nil {
		return err
	}

	ctx, cancel := workerCommandContext(cfg.Timeout)
	defer cancel()

	message := fmt.Sprintf(
		"Billing Telegram delivery preflight\nEnvironment: %s\nTimestamp UTC: %s\nPayload: redacted test message",
		safeEnvLabel(os.Getenv("APP_ENV")),
		time.Now().UTC().Format(time.RFC3339),
	)
	if err := sender.Send(ctx, message); err != nil {
		return fmt.Errorf("send telegram preflight: %w", err)
	}
	fmt.Fprintln(deps.stdout, "telegram-preflight result=PASS")
	fmt.Fprintln(deps.stdout, "telegram_api_called=yes")
	fmt.Fprintln(deps.stdout, "message_payload_redacted=yes")
	fmt.Fprintln(deps.stdout, "secrets_printed=no")
	return nil
}

func validateNotificationTelegramWorkerConfig(cfg workerConfig, command string) error {
	if err := validateNotificationTelegramEnvironment(); err != nil {
		return err
	}
	if cfg.DSN == "" {
		return fmt.Errorf("DB_DSN or -dsn is required for %s", command)
	}
	if cfg.WorkerID == "" {
		return fmt.Errorf("worker id is required")
	}
	return readTelegramConfigFromEnv().Normalize().Validate()
}

func newNotificationTelegramRunner(ctx context.Context, cfg workerConfig) (provisionRunner, func() error, error) {
	conn, err := platformdb.Open(ctx, platformdb.Config{
		DriverName: platformdb.DefaultDriverName,
		DSN:        cfg.DSN,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("open worker database: %w", err)
	}

	jobStore := jobs.NewPostgresStore(conn)
	deliveryStore := notification.NewPostgresStore(conn)
	runner, err := notification.NewTelegramDeliveryRunner(
		jobStore,
		deliveryStore,
		jobs.WorkerID(cfg.WorkerID),
		readTelegramConfigFromEnv(),
		&http.Client{Timeout: cfg.Timeout},
	)
	if err != nil {
		_ = conn.Close()
		return nil, nil, err
	}
	runner.BatchSize = cfg.BatchSize
	runner.LockFor = cfg.LockFor
	return runner, closeDatabase(conn), nil
}

func readTelegramConfigFromEnv() notification.TelegramConfig {
	return notification.TelegramConfig{
		BotToken:   os.Getenv("TELEGRAM_BOT_TOKEN"),
		ChatID:     os.Getenv("TELEGRAM_CHAT_ID"),
		APIBaseURL: os.Getenv("TELEGRAM_API_BASE_URL"),
	}
}

func validateNotificationTelegramEnvironment() error {
	environment := strings.ToLower(strings.TrimSpace(os.Getenv("APP_ENV")))
	switch environment {
	case "local", "dev", "staging":
		return nil
	case "prod", "production":
		if os.Getenv("BILLING_TELEGRAM_DELIVERY_PRODUCTION_APPROVED") == "yes" {
			return nil
		}
		return fmt.Errorf("refusing to run telegram notification delivery with APP_ENV=%s without BILLING_TELEGRAM_DELIVERY_PRODUCTION_APPROVED=yes", os.Getenv("APP_ENV"))
	default:
		return fmt.Errorf("APP_ENV must be local, dev, staging, or production for telegram notification delivery")
	}
}

func safeEnvLabel(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "unknown"
	}
	return value
}
