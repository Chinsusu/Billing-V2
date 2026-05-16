package seed

import (
	"context"
	"database/sql"
	"regexp"
	"strings"
	"testing"

	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
)

func TestDevStatementsAreNamedAndIdempotent(t *testing.T) {
	statements := DevStatements()
	if len(statements) == 0 {
		t.Fatal("expected seed statements")
	}
	seen := map[string]struct{}{}
	for _, statement := range statements {
		if strings.TrimSpace(statement.Name) == "" {
			t.Fatal("statement name is required")
		}
		if strings.TrimSpace(statement.SQL) == "" {
			t.Fatalf("statement %q SQL is required", statement.Name)
		}
		if _, ok := seen[statement.Name]; ok {
			t.Fatalf("duplicate statement name %q", statement.Name)
		}
		seen[statement.Name] = struct{}{}
		if !strings.Contains(strings.ToUpper(statement.SQL), "ON CONFLICT") {
			t.Fatalf("statement %q must be idempotent", statement.Name)
		}
	}
}

func TestDevStatementsIncludeRBACAndCatalogData(t *testing.T) {
	sql := strings.ToLower(joinSeedSQL(DevStatements()))
	required := []string{
		"catalog.view",
		"catalog.manage",
		"tenant.view",
		"order.manage",
		"platform_admin",
		"demo-reseller",
		"master_products",
		"tenant_plans",
		"vps-cx23-40gb-monthly",
	}
	for _, value := range required {
		if !strings.Contains(sql, value) {
			t.Fatalf("expected seed SQL to contain %q", value)
		}
	}
}

func TestDevStatementsIncludeBillingFlowData(t *testing.T) {
	sql := strings.ToLower(joinSeedSQL(DevStatements()))
	required := []string{
		"customer@local.billing",
		"billing_flow",
		"wallet_ledger_entries",
		"payment_transactions",
		"seed-payment-1",
		"seed-provider-provision-1",
		"job.retry",
		"local-vps-405910",
	}
	for _, value := range required {
		if !strings.Contains(sql, value) {
			t.Fatalf("expected seed SQL to contain %q", value)
		}
	}
}

func TestDevUserPasswordHashesAreUsable(t *testing.T) {
	sql := joinSeedSQL(DevStatements())
	normalizedSQL := strings.ToLower(sql)
	if strings.Contains(sql, "dev-only-placeholder-hash") {
		t.Fatal("dev user seed must not use placeholder password hashes")
	}
	if count := strings.Count(normalizedSQL, "password_hash = excluded.password_hash"); count < 3 {
		t.Fatalf("expected dev user seed conflicts to update password hashes, got %d", count)
	}
	matches := regexp.MustCompile(`'\$argon2id[^']+'`).FindAllString(sql, -1)
	if len(matches) == 0 {
		t.Fatal("expected dev user seed to contain argon2id password hashes")
	}
	for _, match := range matches {
		ok, err := identity.VerifyPasswordArgon2id("admin123", strings.Trim(match, "'"))
		if err != nil {
			t.Fatalf("dev user hash is invalid: %v", err)
		}
		if !ok {
			t.Fatal("dev user hash does not verify the local-only password")
		}
	}
}

func TestDevStatementsIncludeProviderReadinessScenarios(t *testing.T) {
	sql := strings.ToLower(joinSeedSQL(DevStatements()))
	required := []string{
		"local fake hetzner ready",
		"local fake hetzner maintenance",
		"vps-maintenance-example-monthly",
		"supportsautoprovision",
		"maintenance",
		"00000000-0000-0000-0000-000000000302",
		"00000000-0000-0000-0000-000000000303",
	}
	for _, value := range required {
		if !strings.Contains(sql, value) {
			t.Fatalf("expected seed SQL to contain readiness scenario %q", value)
		}
	}
}

func TestApplyDevRunsStatementsInOrder(t *testing.T) {
	executor := &fakeSeedExecutor{}
	if err := ApplyDev(context.Background(), executor); err != nil {
		t.Fatalf("ApplyDev returned error: %v", err)
	}
	if len(executor.queries) != len(DevStatements()) {
		t.Fatalf("expected %d query runs, got %d", len(DevStatements()), len(executor.queries))
	}
	if !strings.Contains(executor.queries[0], "INSERT INTO permissions") {
		t.Fatalf("expected permissions to run first, got %s", executor.queries[0])
	}
}

func joinSeedSQL(statements []Statement) string {
	var builder strings.Builder
	for _, statement := range statements {
		builder.WriteString(statement.SQL)
		builder.WriteString("\n")
	}
	return builder.String()
}

type fakeSeedExecutor struct {
	queries []string
}

func (executor *fakeSeedExecutor) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	executor.queries = append(executor.queries, query)
	return fakeSQLResult{}, nil
}

func (executor *fakeSeedExecutor) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return nil, nil
}

func (executor *fakeSeedExecutor) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return nil
}

type fakeSQLResult struct{}

func (fakeSQLResult) LastInsertId() (int64, error) { return 0, nil }
func (fakeSQLResult) RowsAffected() (int64, error) { return 0, nil }
