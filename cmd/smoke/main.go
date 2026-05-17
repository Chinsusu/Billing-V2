package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/platform/db"
	"github.com/Chinsusu/Billing-V2/internal/seed"
)

type checkMode string

const (
	checkExact checkMode = "exact"
	checkMin   checkMode = "min"
)

type smokeCheck struct {
	Name  string
	Query string
	Mode  checkMode
	Want  int
}

type checkResult struct {
	Name  string
	Count int
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "smoke failed: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	flags := flag.NewFlagSet("smoke", flag.ContinueOnError)
	dir := flags.String("dir", "migrations", "migration directory")
	dsn := flags.String("dsn", os.Getenv("DB_DSN"), "PostgreSQL DSN")
	baseURL := flags.String("base-url", envOrDefault("API_BASE_URL", "http://localhost:8080"), "API base URL")
	timeout := flags.Duration("timeout", 60*time.Second, "smoke command timeout")
	if err := flags.Parse(args); err != nil {
		return err
	}

	command := "dev-db"
	if flags.NArg() > 0 {
		command = flags.Arg(0)
	}

	switch command {
	case "dev-db":
		return runDevDBSmoke(*dsn, *dir, *timeout)
	case "dev-api":
		return runDevAPISmoke(*baseURL, *timeout)
	case "dev-billing":
		return runDevBillingMutationSmoke(*dsn, *baseURL, *timeout)
	case "dev-topup-review":
		return runDevTopupReviewSmoke(*dsn, *baseURL, *timeout)
	default:
		return fmt.Errorf("unknown command %q; use dev-db, dev-api, dev-billing, or dev-topup-review", command)
	}
}

func envOrDefault(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func runDevDBSmoke(dsn string, migrationDir string, timeout time.Duration) error {
	if err := guardDevEnvironment(); err != nil {
		return err
	}
	if strings.TrimSpace(dsn) == "" {
		return fmt.Errorf("DB_DSN or -dsn is required for dev-db smoke")
	}

	migrations, err := db.LoadMigrations(os.DirFS(migrationDir))
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn, err := db.Open(ctx, db.Config{DriverName: db.DefaultDriverName, DSN: dsn})
	if err != nil {
		return err
	}
	defer conn.Close()

	migrator, err := db.NewMigrator(conn, migrations)
	if err != nil {
		return err
	}
	applied, err := migrator.ApplyAll(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("applied migration(s): %d\n", len(applied))

	if err := seed.ApplyDev(ctx, conn); err != nil {
		return err
	}
	fmt.Printf("applied seed statement(s): %d\n", len(seed.DevStatements()))

	if err := seed.ApplyDev(ctx, conn); err != nil {
		return fmt.Errorf("reapply dev seed for idempotency: %w", err)
	}
	fmt.Printf("reapplied seed statement(s): %d\n", len(seed.DevStatements()))

	results, err := runChecks(ctx, conn, migrations)
	if err != nil {
		return err
	}
	for _, result := range results {
		fmt.Printf("check passed: %s (%d row(s))\n", result.Name, result.Count)
	}
	fmt.Printf("dev DB smoke passed: %d check(s)\n", len(results))
	return nil
}

func guardDevEnvironment() error {
	switch strings.ToLower(strings.TrimSpace(os.Getenv("APP_ENV"))) {
	case "prod", "production":
		return fmt.Errorf("refusing to run dev-db smoke with APP_ENV=%s", os.Getenv("APP_ENV"))
	default:
		return nil
	}
}

func runChecks(ctx context.Context, conn *sql.DB, migrations []db.Migration) ([]checkResult, error) {
	checks := append([]smokeCheck{{
		Name:  "schema migrations applied",
		Query: "SELECT COUNT(*) FROM schema_migrations",
		Mode:  checkMin,
		Want:  len(migrations),
	}}, seededBillingChecks()...)

	results := make([]checkResult, 0, len(checks))
	for _, check := range checks {
		count, err := countRows(ctx, conn, check.Query)
		if err != nil {
			return nil, fmt.Errorf("run check %q: %w", check.Name, err)
		}
		if err := check.validate(count); err != nil {
			return nil, err
		}
		results = append(results, checkResult{Name: check.Name, Count: count})
	}
	return results, nil
}

func seededBillingChecks() []smokeCheck {
	return []smokeCheck{
		exactCheck("audit table exists", "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'audit_logs'", 1),
		exactCheck("platform tenant", "SELECT COUNT(*) FROM tenants WHERE tenant_id = '00000000-0000-0000-0000-000000000001' AND slug = 'platform'", 1),
		exactCheck("demo reseller tenant", "SELECT COUNT(*) FROM tenants WHERE tenant_id = '00000000-0000-0000-0000-000000000010' AND slug = 'demo-reseller'", 1),
		exactCheck("demo customer user", "SELECT COUNT(*) FROM users WHERE user_id = '00000000-0000-0000-0000-000000000103' AND email = 'customer@local.billing'", 1),
		minCheck("rbac permissions", "SELECT COUNT(*) FROM permissions WHERE permission_key IN ('wallet.view', 'order.view', 'audit.view')", 3),
		minCheck("provider source", "SELECT COUNT(*) FROM provider_sources WHERE display_id >= 10000", 1),
		exactCheck("tenant plan", "SELECT COUNT(*) FROM tenant_plans WHERE tenant_plan_id = '00000000-0000-0000-0000-000000000801'", 1),
		exactCheck("demo wallet", "SELECT COUNT(*) FROM wallets WHERE wallet_id = '00000000-0000-0000-0000-000000000901' AND display_id = 41001 AND available_balance_minor = 3200", 1),
		exactCheck("topup request", "SELECT COUNT(*) FROM topup_requests WHERE topup_request_id = '00000000-0000-0000-0000-000000000908' AND display_id = 52001 AND status = 'approved'", 1),
		exactCheck("topup idempotency", "SELECT COUNT(*) FROM topup_requests WHERE tenant_id = '00000000-0000-0000-0000-000000000010' AND idempotency_key = 'seed-topup-request-1'", 1),
		exactCheck("demo order", "SELECT COUNT(*) FROM orders WHERE order_id = '00000000-0000-0000-0000-000000000903' AND display_id = 42001 AND order_status = 'paid' AND billing_status = 'paid'", 1),
		exactCheck("order idempotency", "SELECT COUNT(*) FROM orders WHERE tenant_id = '00000000-0000-0000-0000-000000000010' AND idempotency_key = 'seed-order-billing-flow-1'", 1),
		exactCheck("service instance", "SELECT COUNT(*) FROM service_instances WHERE service_instance_id = '00000000-0000-0000-0000-000000000909' AND display_id = 43001 AND status = 'active'", 1),
		exactCheck("paid invoice", "SELECT COUNT(*) FROM invoices WHERE invoice_id = '00000000-0000-0000-0000-000000000904' AND display_id = 44001 AND status = 'paid'", 1),
		exactCheck("invoice item", "SELECT COUNT(*) FROM invoice_items WHERE invoice_item_id = '00000000-0000-0000-0000-000000000905'", 1),
		exactCheck("wallet ledger entries", "SELECT COUNT(*) FROM wallet_ledger_entries WHERE wallet_id = '00000000-0000-0000-0000-000000000901' AND display_id IN (50001, 50002) AND status = 'posted'", 2),
		exactCheck("payment transaction", "SELECT COUNT(*) FROM payment_transactions WHERE payment_transaction_id = '00000000-0000-0000-0000-000000000907' AND display_id = 51001 AND status = 'posted'", 1),
		exactCheck("payment idempotency", "SELECT COUNT(*) FROM payment_transactions WHERE tenant_id = '00000000-0000-0000-0000-000000000010' AND idempotency_key = 'seed-payment-1'", 1),
	}
}

func exactCheck(name string, query string, want int) smokeCheck {
	return smokeCheck{Name: name, Query: query, Mode: checkExact, Want: want}
}

func minCheck(name string, query string, want int) smokeCheck {
	return smokeCheck{Name: name, Query: query, Mode: checkMin, Want: want}
}

func countRows(ctx context.Context, conn *sql.DB, query string) (int, error) {
	var count int
	if err := conn.QueryRowContext(ctx, query).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (check smokeCheck) validate(count int) error {
	switch check.Mode {
	case checkExact:
		if count != check.Want {
			return fmt.Errorf("check %q expected exactly %d row(s), got %d", check.Name, check.Want, count)
		}
	case checkMin:
		if count < check.Want {
			return fmt.Errorf("check %q expected at least %d row(s), got %d", check.Name, check.Want, count)
		}
	default:
		return fmt.Errorf("check %q has unknown mode %q", check.Name, check.Mode)
	}
	return nil
}
