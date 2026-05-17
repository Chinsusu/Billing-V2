package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	platformdb "github.com/Chinsusu/Billing-V2/internal/platform/db"
)

const targetFinanceSmokeActorType = "platform_staff"

type targetFinanceCandidate struct {
	TransactionID        string
	TransactionDisplayID int64
	InvoiceDisplayID     int64
	WalletDisplayID      int64
	LedgerDisplayID      int64
	AmountMinor          int64
	ReconciliationDate   string
}

type targetFinanceBaseline struct {
	TopupCount       int
	TransactionCount int
	LedgerCount      int
	WalletBalanceSum int64
	OrderCount       int
	ProviderJobCount int
	ServiceCount     int
}

type targetFinanceReconciliationResponse struct {
	Transaction struct {
		ID               string          `json:"id"`
		DisplayID        int64           `json:"display_id"`
		Status           string          `json:"status"`
		AmountMinor      int64           `json:"amount_minor"`
		InvoiceDisplayID int64           `json:"invoice_display_id"`
		Metadata         json.RawMessage `json:"metadata"`
	} `json:"transaction"`
	Provider string `json:"provider"`
	Invoice  *struct {
		DisplayID  int64  `json:"display_id"`
		Status     string `json:"status"`
		TotalMinor int64  `json:"total_minor"`
	} `json:"invoice"`
	Ledger *struct {
		DisplayID       int64  `json:"display_id"`
		WalletDisplayID int64  `json:"wallet_display_id"`
		Direction       string `json:"direction"`
		EntryType       string `json:"entry_type"`
		Status          string `json:"status"`
	} `json:"ledger"`
}

type targetFinanceDailyReconciliationResponse struct {
	Date    string `json:"date"`
	Status  string `json:"status"`
	Wallets struct {
		Checked    int `json:"checked"`
		Balanced   int `json:"balanced"`
		Mismatched int `json:"mismatched"`
	} `json:"wallets"`
	Invoices struct {
		Checked    int `json:"checked"`
		Mismatched int `json:"mismatched"`
	} `json:"invoices"`
	Payments struct {
		Checked                 int `json:"checked"`
		DuplicateReferenceCount int `json:"duplicate_reference_count"`
	} `json:"payments"`
}

func runDevTargetFinanceReconciliationSmoke(dsn string, baseURL string, timeout time.Duration) error {
	if err := guardDevEnvironment(); err != nil {
		return err
	}
	dsn = strings.TrimSpace(dsn)
	if dsn == "" {
		return fmt.Errorf("DB_DSN or -dsn is required for dev-target-finance-reconciliation smoke")
	}
	baseURL = strings.TrimSpace(baseURL)
	if baseURL == "" {
		return fmt.Errorf("API_BASE_URL or -base-url is required for dev-target-finance-reconciliation smoke")
	}
	if _, err := normalizedAPIURL(baseURL, "/healthz"); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn, err := platformdb.Open(ctx, platformdb.Config{DriverName: platformdb.DefaultDriverName, DSN: dsn})
	if err != nil {
		return fmt.Errorf("open smoke DB for target finance reconciliation checks: %w", err)
	}
	defer conn.Close()

	candidate, err := loadTargetFinanceCandidate(ctx, conn)
	if err != nil {
		return err
	}
	baseline, err := readTargetFinanceBaseline(ctx, conn)
	if err != nil {
		return err
	}

	client := &http.Client{Timeout: timeout}
	list, err := runTargetFinanceReconciliationList(ctx, client, baseURL, candidate)
	if err != nil {
		return err
	}
	if err := validateTargetFinanceReconciliationRecord(list, candidate); err != nil {
		return err
	}
	detail, err := getTargetFinanceJSON[targetFinanceReconciliationResponse](ctx, client, baseURL, "/admin/payment-reconciliation/"+candidate.TransactionID)
	if err != nil {
		return err
	}
	if err := validateTargetFinanceReconciliationRecord(detail, candidate); err != nil {
		return err
	}
	daily, err := getTargetFinanceJSON[targetFinanceDailyReconciliationResponse](ctx, client, baseURL, "/admin/daily-reconciliation?date="+url.QueryEscape(candidate.ReconciliationDate))
	if err != nil {
		return err
	}
	if err := validateTargetFinanceDailyReconciliation(daily, candidate); err != nil {
		return err
	}

	after, err := readTargetFinanceBaseline(ctx, conn)
	if err != nil {
		return err
	}
	if baseline != after {
		return fmt.Errorf("target finance reconciliation smoke changed read-only baseline")
	}

	fmt.Printf("target finance reconciliation smoke passed: transaction_display_id=%d invoice_display_id=%d wallet_display_id=%d ledger_display_id=%d daily_date=%s daily_status=%s wallets_checked=%d invoices_checked=%d payments_checked=%d wallet_mismatches=%d invoice_mismatches=%d duplicate_payment_references=%d money_mutation_routes_called=no provider_mutation_routes_called=no\n",
		candidate.TransactionDisplayID,
		candidate.InvoiceDisplayID,
		candidate.WalletDisplayID,
		candidate.LedgerDisplayID,
		candidate.ReconciliationDate,
		daily.Status,
		daily.Wallets.Checked,
		daily.Invoices.Checked,
		daily.Payments.Checked,
		daily.Wallets.Mismatched,
		daily.Invoices.Mismatched,
		daily.Payments.DuplicateReferenceCount,
	)
	fmt.Println("Target finance reconciliation smoke output intentionally excludes raw transaction IDs, invoice IDs, wallet IDs, ledger IDs, actor IDs, session tokens, cookies, DSNs, provider payloads, and credentials.")
	return nil
}

func loadTargetFinanceCandidate(ctx context.Context, conn *sql.DB) (targetFinanceCandidate, error) {
	var candidate targetFinanceCandidate
	err := conn.QueryRowContext(ctx, `
SELECT
  txn.payment_transaction_id::text,
  txn.display_id,
  invoice.display_id,
  linked_wallet.display_id,
  ledger.display_id,
  txn.amount_minor,
  txn.created_at::date::text
FROM payment_transactions txn
JOIN invoices invoice
  ON invoice.invoice_id = txn.invoice_id
 AND invoice.tenant_id = txn.tenant_id
JOIN LATERAL (
  SELECT ledger_entry_id, display_id, wallet_id
  FROM wallet_ledger_entries ledger
  WHERE ledger.tenant_id = txn.tenant_id
    AND ledger.reference_type = 'invoice'
    AND ledger.reference_id = txn.invoice_id
    AND ledger.entry_type = 'purchase'
    AND ledger.status = 'posted'
  ORDER BY ledger.created_at DESC
  LIMIT 1
) ledger ON TRUE
JOIN wallets linked_wallet
  ON linked_wallet.wallet_id = ledger.wallet_id
 AND linked_wallet.tenant_id = txn.tenant_id
WHERE txn.tenant_id = $1::uuid
  AND txn.transaction_type = 'charge'
  AND txn.status = 'posted'
  AND COALESCE(txn.metadata->>'provider', '') = 'wallet'
ORDER BY txn.created_at DESC, txn.display_id DESC
LIMIT 1`,
		demoTenantID,
	).Scan(
		&candidate.TransactionID,
		&candidate.TransactionDisplayID,
		&candidate.InvoiceDisplayID,
		&candidate.WalletDisplayID,
		&candidate.LedgerDisplayID,
		&candidate.AmountMinor,
		&candidate.ReconciliationDate,
	)
	if err != nil {
		return targetFinanceCandidate{}, fmt.Errorf("load target finance reconciliation candidate: %w", err)
	}
	if candidate.TransactionDisplayID <= 0 || candidate.InvoiceDisplayID <= 0 || candidate.WalletDisplayID <= 0 || candidate.LedgerDisplayID <= 0 || candidate.AmountMinor <= 0 || candidate.ReconciliationDate == "" {
		return targetFinanceCandidate{}, fmt.Errorf("target finance reconciliation candidate missing public evidence fields")
	}
	return candidate, nil
}

func readTargetFinanceBaseline(ctx context.Context, conn *sql.DB) (targetFinanceBaseline, error) {
	var baseline targetFinanceBaseline
	err := conn.QueryRowContext(ctx, `
SELECT
  (SELECT COUNT(*) FROM topup_requests WHERE tenant_id = $1::uuid),
  (SELECT COUNT(*) FROM payment_transactions WHERE tenant_id = $1::uuid),
  (SELECT COUNT(*) FROM wallet_ledger_entries WHERE tenant_id = $1::uuid),
  (SELECT COALESCE(SUM(available_balance_minor), 0) FROM wallets WHERE tenant_id = $1::uuid),
  (SELECT COUNT(*) FROM orders WHERE tenant_id = $1::uuid),
  (SELECT COUNT(*) FROM jobs WHERE tenant_id = $1::uuid AND job_type = 'provider.provision'),
  (SELECT COUNT(*) FROM service_instances WHERE tenant_id = $1::uuid)`,
		demoTenantID,
	).Scan(
		&baseline.TopupCount,
		&baseline.TransactionCount,
		&baseline.LedgerCount,
		&baseline.WalletBalanceSum,
		&baseline.OrderCount,
		&baseline.ProviderJobCount,
		&baseline.ServiceCount,
	)
	if err != nil {
		return targetFinanceBaseline{}, fmt.Errorf("read target finance reconciliation baseline: %w", err)
	}
	return baseline, nil
}

func runTargetFinanceReconciliationList(
	ctx context.Context,
	client *http.Client,
	baseURL string,
	candidate targetFinanceCandidate,
) (targetFinanceReconciliationResponse, error) {
	query := url.Values{}
	query.Set("display_id", fmt.Sprintf("%d", candidate.TransactionDisplayID))
	query.Set("status", "posted")
	query.Set("provider", "wallet")
	query.Set("invoice_display_id", fmt.Sprintf("%d", candidate.InvoiceDisplayID))
	query.Set("wallet_display_id", fmt.Sprintf("%d", candidate.WalletDisplayID))
	query.Set("limit", "5")
	records, err := getTargetFinanceJSON[[]targetFinanceReconciliationResponse](ctx, client, baseURL, "/admin/payment-reconciliation?"+query.Encode())
	if err != nil {
		return targetFinanceReconciliationResponse{}, err
	}
	if len(records) != 1 {
		return targetFinanceReconciliationResponse{}, fmt.Errorf("target finance reconciliation list expected one matching record, got %d", len(records))
	}
	return records[0], nil
}

func getTargetFinanceJSON[T any](ctx context.Context, client *http.Client, baseURL string, path string) (T, error) {
	var zero T
	fullURL, err := normalizedAPIURL(baseURL, path)
	if err != nil {
		return zero, err
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)
	if err != nil {
		return zero, fmt.Errorf("build target finance reconciliation request")
	}
	for key, value := range targetFinanceReadHeaders() {
		request.Header.Set(key, value)
	}

	response, err := client.Do(request)
	if err != nil {
		return zero, fmt.Errorf("request target finance reconciliation failed")
	}
	defer response.Body.Close()
	body, err := io.ReadAll(io.LimitReader(response.Body, 1<<20))
	if err != nil {
		return zero, fmt.Errorf("read target finance reconciliation response: %w", err)
	}
	if response.StatusCode != http.StatusOK {
		return zero, targetFinanceStatusError(response.StatusCode, body)
	}
	if err := assertTargetFinanceResponseRedaction(body); err != nil {
		return zero, err
	}
	var envelope struct {
		Data      T      `json:"data"`
		RequestID string `json:"request_id"`
	}
	if err := json.Unmarshal(body, &envelope); err != nil {
		return zero, fmt.Errorf("decode target finance reconciliation response")
	}
	if strings.TrimSpace(envelope.RequestID) == "" {
		return zero, fmt.Errorf("target finance reconciliation response missing request_id")
	}
	return envelope.Data, nil
}

func targetFinanceReadHeaders() map[string]string {
	headers := actorHeaders(demoResellerID, targetFinanceSmokeActorType)
	headers["X-Actor-Tenant-Id"] = demoTenantID
	return headers
}

func validateTargetFinanceReconciliationRecord(record targetFinanceReconciliationResponse, candidate targetFinanceCandidate) error {
	if record.Transaction.ID == "" || record.Transaction.DisplayID != candidate.TransactionDisplayID {
		return fmt.Errorf("target finance reconciliation transaction mismatch")
	}
	if record.Transaction.Status != "posted" || record.Transaction.AmountMinor != candidate.AmountMinor || record.Provider != "wallet" {
		return fmt.Errorf("target finance reconciliation transaction status or provider mismatch")
	}
	if record.Invoice == nil || record.Invoice.DisplayID != candidate.InvoiceDisplayID || record.Invoice.Status != "paid" || record.Invoice.TotalMinor != candidate.AmountMinor {
		return fmt.Errorf("target finance reconciliation invoice mismatch")
	}
	if record.Ledger == nil ||
		record.Ledger.DisplayID != candidate.LedgerDisplayID ||
		record.Ledger.WalletDisplayID != candidate.WalletDisplayID ||
		record.Ledger.Direction != "debit" ||
		record.Ledger.EntryType != "purchase" ||
		record.Ledger.Status != "posted" {
		return fmt.Errorf("target finance reconciliation ledger mismatch")
	}
	return nil
}

func validateTargetFinanceDailyReconciliation(report targetFinanceDailyReconciliationResponse, candidate targetFinanceCandidate) error {
	if report.Date != candidate.ReconciliationDate {
		return fmt.Errorf("target daily reconciliation report date mismatch")
	}
	if report.Status != "balanced" && report.Status != "mismatched" {
		return fmt.Errorf("target daily reconciliation report status mismatch")
	}
	if report.Wallets.Checked <= 0 || report.Wallets.Balanced+report.Wallets.Mismatched != report.Wallets.Checked {
		return fmt.Errorf("target daily reconciliation wallet summary mismatch")
	}
	if report.Invoices.Checked <= 0 {
		return fmt.Errorf("target daily reconciliation invoice summary mismatch")
	}
	if report.Payments.Checked <= 0 {
		return fmt.Errorf("target daily reconciliation payment summary mismatch")
	}
	if report.Status == "balanced" && (report.Wallets.Mismatched != 0 || report.Invoices.Mismatched != 0 || report.Payments.DuplicateReferenceCount != 0) {
		return fmt.Errorf("target daily reconciliation balanced report contains mismatches")
	}
	if report.Status == "mismatched" && report.Wallets.Mismatched+report.Invoices.Mismatched+report.Payments.DuplicateReferenceCount == 0 {
		return fmt.Errorf("target daily reconciliation mismatched report lacks mismatch counts")
	}
	return nil
}

func assertTargetFinanceResponseRedaction(body []byte) error {
	bodyLower := strings.ToLower(string(body))
	blocked := []string{
		`"idempotency_key"`,
		`"encrypted_payload"`,
		`"raw_response"`,
		`"raw_payload"`,
		`"provider_credentials"`,
		`"api_key"`,
		`"token_hash"`,
		`"session_token"`,
		`"password"`,
		`"credential"`,
	}
	for _, token := range blocked {
		if strings.Contains(bodyLower, token) {
			return fmt.Errorf("target finance reconciliation response exposed blocked field")
		}
	}
	return nil
}

func targetFinanceStatusError(gotStatus int, body []byte) error {
	var apiError errorEnvelope
	if err := json.Unmarshal(body, &apiError); err == nil && apiError.Error.Code != "" {
		return fmt.Errorf("target finance reconciliation expected HTTP 200, got %d (%s)", gotStatus, apiError.Error.Code)
	}
	return fmt.Errorf("target finance reconciliation expected HTTP 200, got %d", gotStatus)
}
