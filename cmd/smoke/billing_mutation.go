package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	platformdb "github.com/Chinsusu/Billing-V2/internal/platform/db"
)

const (
	smokeWalletID      = "00000000-0000-0000-0000-000000000901"
	smokeTenantPlanID  = "00000000-0000-0000-0000-000000000801"
	smokeTopupAmount   = int64(1800)
	smokeOrderAmount   = int64(1400)
	smokeOrderCurrency = "USD"
	auditActionTopup   = "wallet.topup.approved"
	auditActionInvoice = "invoice.wallet_paid"
	auditTargetTopup   = "topup_request"
	auditTargetInvoice = "invoice"
)

var (
	smokeOrderProductSnapshot = json.RawMessage(`{"name":"VPS","product_type":"vps"}`)
	smokeOrderPlanSnapshot    = json.RawMessage(`{"plan_code":"vps-cx23-40gb-monthly","name":"CX23 VPS 40GB"}`)
	smokeOrderPriceSnapshot   = json.RawMessage(`{"selling_price_minor":1400,"currency":"USD"}`)
)

type billingMutationScenario struct {
	RunID string
}

func runDevBillingMutationSmoke(dsn string, baseURL string, timeout time.Duration) error {
	if err := guardDevEnvironment(); err != nil {
		return err
	}
	dsn = strings.TrimSpace(dsn)
	if dsn == "" {
		return fmt.Errorf("DB_DSN or -dsn is required for dev-billing smoke")
	}
	baseURL = strings.TrimSpace(baseURL)
	if baseURL == "" {
		return fmt.Errorf("API_BASE_URL or -base-url is required for dev-billing smoke")
	}
	if _, err := normalizedAPIURL(baseURL, "/healthz"); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	client := &http.Client{Timeout: timeout}
	scenario := newBillingMutationScenario()
	conn, err := platformdb.Open(ctx, platformdb.Config{DriverName: platformdb.DefaultDriverName, DSN: dsn})
	if err != nil {
		return fmt.Errorf("open smoke DB for provisioning checks: %w", err)
	}
	defer conn.Close()

	topup, err := runTopupApprovalSmoke(ctx, client, baseURL, scenario)
	if err != nil {
		return err
	}

	orderRecord, err := runOrderCreateSmoke(ctx, client, baseURL, scenario)
	if err != nil {
		return err
	}

	issuedInvoice, err := runCheckoutInvoiceSmoke(ctx, client, baseURL, scenario, orderRecord.ID)
	if err != nil {
		return err
	}

	if err := verifyIssuedInvoiceVisibleViaAPI(ctx, client, baseURL, orderRecord.ID, issuedInvoice.ID); err != nil {
		return err
	}

	paymentRecord, err := runInvoiceWalletPaymentSmoke(ctx, client, baseURL, scenario, issuedInvoice.ID, orderRecord.ID)
	if err != nil {
		return err
	}
	if err := verifyPaidOrderVisibleViaAPI(ctx, client, baseURL, orderRecord.ID, orderRecord.DisplayID); err != nil {
		return err
	}
	if err := verifyProvisioningJobQueued(ctx, conn, orderRecord.ID, orderRecord.DisplayID); err != nil {
		return err
	}
	serviceRecord, err := runProvisioningFulfillmentSmoke(ctx, conn, client, baseURL, scenario, orderRecord.ID, orderRecord.DisplayID)
	if err != nil {
		return err
	}

	checks := []auditMutationCheck{
		{
			Action:           auditActionTopup,
			TargetType:       auditTargetTopup,
			TargetID:         topup.ID,
			MetadataContains: fmt.Sprintf(`"display_id":%d`, topup.DisplayID),
			AfterContains:    `"status":"approved"`,
		},
		{
			Action:           auditActionInvoice,
			TargetType:       auditTargetInvoice,
			TargetID:         paymentRecord.Invoice.ID,
			MetadataContains: fmt.Sprintf(`"transaction_display_id":%d`, paymentRecord.Transaction.DisplayID),
			AfterContains:    `"status":"paid"`,
		},
	}
	for _, check := range checks {
		if err := verifyAuditMutation(ctx, client, baseURL, check); err != nil {
			return err
		}
	}

	fmt.Printf("dev billing smoke passed: topup=%d order=%d invoice=%d transaction=%d ledger=%d service=%d\n",
		topup.DisplayID,
		orderRecord.DisplayID,
		paymentRecord.Invoice.DisplayID,
		paymentRecord.Transaction.DisplayID,
		paymentLedgerDisplayID(paymentRecord),
		serviceRecord.DisplayID,
	)
	return nil
}

func runTopupApprovalSmoke(
	ctx context.Context,
	client *http.Client,
	baseURL string,
	scenario billingMutationScenario,
) (topupResponse, error) {
	headers := cloneHeaders(clientHeaders())
	headers["Idempotency-Key"] = scenario.topupIdempotencyKey()
	created, err := doJSON[topupResponse](ctx, client, http.MethodPost, baseURL, "/client/topup-requests", headers, topupRequestBody{
		WalletID:         smokeWalletID,
		AmountMinor:      smokeTopupAmount,
		Currency:         smokeOrderCurrency,
		PaymentMethod:    "manual",
		PaymentReference: scenario.topupPaymentReference(),
	}, http.StatusCreated)
	if err != nil {
		return topupResponse{}, err
	}
	if created.DisplayID <= 0 || created.Status != "submitted" {
		return topupResponse{}, fmt.Errorf("expected submitted topup with display id, got %+v", created)
	}

	approved, err := doJSON[topupResponse](ctx, client, http.MethodPost, baseURL, "/reseller/topup-requests/"+created.ID+"/approve", resellerHeaders(), reviewTopupBody{
		ReviewNote: "Smoke approval " + scenario.RunID,
	}, http.StatusOK)
	if err != nil {
		return topupResponse{}, err
	}
	if approved.Status != "approved" || approved.DisplayID <= 0 || approved.LedgerEntryID == "" {
		return topupResponse{}, fmt.Errorf("expected approved topup with ledger entry, got %+v", approved)
	}
	fmt.Printf("billing mutation passed: topup approved %s (%d)\n", approved.ID, approved.DisplayID)
	return approved, nil
}

func runOrderCreateSmoke(
	ctx context.Context,
	client *http.Client,
	baseURL string,
	scenario billingMutationScenario,
) (orderResponse, error) {
	headers := cloneHeaders(clientHeaders())
	headers["Idempotency-Key"] = scenario.orderIdempotencyKey()
	created, err := doJSON[orderResponse](ctx, client, http.MethodPost, baseURL, "/client/orders", headers, createOrderBody{
		TenantPlanID:    smokeTenantPlanID,
		Quantity:        1,
		Currency:        smokeOrderCurrency,
		UnitPriceMinor:  smokeOrderAmount,
		DiscountMinor:   0,
		TotalMinor:      smokeOrderAmount,
		ProductSnapshot: smokeOrderProductSnapshot,
		PlanSnapshot:    smokeOrderPlanSnapshot,
		PriceSnapshot:   smokeOrderPriceSnapshot,
	}, http.StatusCreated)
	if err != nil {
		return orderResponse{}, err
	}
	if created.DisplayID <= 0 || created.OrderStatus != "pending_payment" || created.BillingStatus != "unpaid" {
		return orderResponse{}, fmt.Errorf("expected pending unpaid order, got %+v", created)
	}
	fmt.Printf("billing mutation passed: order created %s (%d)\n", created.ID, created.DisplayID)
	return created, nil
}

func runCheckoutInvoiceSmoke(
	ctx context.Context,
	client *http.Client,
	baseURL string,
	scenario billingMutationScenario,
	orderID string,
) (invoiceSummaryResponse, error) {
	headers := cloneHeaders(clientHeaders())
	headers["Idempotency-Key"] = scenario.checkoutIdempotencyKey()
	created, err := doJSON[invoiceSummaryResponse](ctx, client, http.MethodPost, baseURL, "/client/checkouts", headers, checkoutOrderBody{
		OrderID: orderID,
	}, http.StatusCreated)
	if err != nil {
		return invoiceSummaryResponse{}, err
	}
	if created.DisplayID <= 0 || created.Status != "issued" || created.OrderID != orderID {
		return invoiceSummaryResponse{}, fmt.Errorf("expected issued checkout invoice, got %+v", created)
	}

	duplicate, err := doJSON[invoiceSummaryResponse](ctx, client, http.MethodPost, baseURL, "/client/checkouts", headers, checkoutOrderBody{
		OrderID: orderID,
	}, http.StatusCreated)
	if err != nil {
		return invoiceSummaryResponse{}, err
	}
	if duplicate.ID != created.ID || duplicate.DisplayID != created.DisplayID {
		return invoiceSummaryResponse{}, fmt.Errorf("expected idempotent checkout invoice, got first=%+v duplicate=%+v", created, duplicate)
	}
	fmt.Printf("billing mutation passed: checkout invoice issued %s (%d)\n", created.ID, created.DisplayID)
	return created, nil
}

func verifyIssuedInvoiceVisibleViaAPI(
	ctx context.Context,
	client *http.Client,
	baseURL string,
	orderID string,
	invoiceID string,
) error {
	query := url.Values{}
	query.Set("order_id", orderID)
	query.Set("status", "issued")
	invoices, err := doJSON[[]invoiceSummaryResponse](ctx, client, http.MethodGet, baseURL, "/client/invoices?"+query.Encode(), clientHeaders(), nil, http.StatusOK)
	if err != nil {
		return err
	}
	for _, record := range invoices {
		if record.ID == invoiceID && record.DisplayID > 0 && record.Status == "issued" {
			fmt.Printf("billing mutation passed: invoice visible via API %s (%d)\n", record.ID, record.DisplayID)
			return nil
		}
	}
	return fmt.Errorf("expected issued invoice %s to be visible via API, got %+v", invoiceID, invoices)
}

func runInvoiceWalletPaymentSmoke(
	ctx context.Context,
	client *http.Client,
	baseURL string,
	scenario billingMutationScenario,
	invoiceID string,
	orderID string,
) (invoiceWalletPaymentResponse, error) {
	headers := cloneHeaders(clientHeaders())
	headers["Idempotency-Key"] = scenario.paymentIdempotencyKey()
	record, err := doJSON[invoiceWalletPaymentResponse](ctx, client, http.MethodPost, baseURL, "/client/invoice-wallet-payments", headers, invoiceWalletPaymentBody{
		InvoiceID: invoiceID,
		WalletID:  smokeWalletID,
	}, http.StatusCreated)
	if err != nil {
		return invoiceWalletPaymentResponse{}, err
	}
	if record.Invoice.ID != invoiceID || record.Invoice.DisplayID <= 0 || record.Invoice.Status != "paid" {
		return invoiceWalletPaymentResponse{}, fmt.Errorf("expected paid invoice response, got %+v", record)
	}
	if record.Transaction.DisplayID <= 0 || record.Transaction.Status != "posted" {
		return invoiceWalletPaymentResponse{}, fmt.Errorf("expected posted transaction with display id, got %+v", record.Transaction)
	}
	if record.Ledger == nil || record.Ledger.DisplayID <= 0 || record.Ledger.AmountMinor <= 0 {
		return invoiceWalletPaymentResponse{}, fmt.Errorf("expected debit ledger entry with display id, got %+v", record.Ledger)
	}
	if record.Order == nil ||
		record.Order.ID != orderID ||
		record.Order.DisplayID <= 0 ||
		record.Order.OrderStatus != "paid" ||
		record.Order.BillingStatus != "paid" {
		return invoiceWalletPaymentResponse{}, fmt.Errorf("expected paid order in wallet payment response, got %+v", record.Order)
	}
	fmt.Printf("billing mutation passed: invoice paid %s (%d)\n", record.Invoice.ID, record.Invoice.DisplayID)
	return record, nil
}

func verifyPaidOrderVisibleViaAPI(
	ctx context.Context,
	client *http.Client,
	baseURL string,
	orderID string,
	displayID int64,
) error {
	record, err := doJSON[orderResponse](ctx, client, http.MethodGet, baseURL, "/client/orders/"+url.PathEscape(orderID), clientHeaders(), nil, http.StatusOK)
	if err != nil {
		return err
	}
	if record.ID != orderID || record.DisplayID != displayID || record.OrderStatus != "paid" || record.BillingStatus != "paid" {
		return fmt.Errorf("expected paid order via API, got %+v", record)
	}
	fmt.Printf("billing mutation passed: order finalized %s (%d)\n", record.ID, record.DisplayID)
	return nil
}

func verifyProvisioningJobQueued(ctx context.Context, conn *sql.DB, orderID string, orderDisplayID int64) error {
	var count int
	var displayID int64
	var status string
	err := conn.QueryRowContext(ctx, `
SELECT COUNT(*), COALESCE(MAX(display_id), 0), COALESCE(MAX(status::text), '')
FROM jobs
WHERE tenant_id = $1
  AND job_type = 'provider.provision'
  AND reference_type = 'order'
  AND reference_id = $2::uuid
  AND source_id IS NOT NULL
  AND payload_json->>'order_id' = $2::text
  AND (payload_json->>'order_display_id')::bigint = $3`,
		demoTenantID,
		orderID,
		orderDisplayID,
	).Scan(&count, &displayID, &status)
	if err != nil {
		return fmt.Errorf("query provisioning job for order %s: %w", orderID, err)
	}
	if count != 1 {
		return fmt.Errorf("expected exactly one provider.provision job for paid order %s, got %d", orderID, count)
	}
	if !provisioningJobSmokeStatusOK(status) {
		return fmt.Errorf("expected provider.provision job for order %s to be queued/running/succeeded, got %q", orderID, status)
	}
	fmt.Printf("billing mutation passed: provisioning job %d status=%s order=%d\n", displayID, status, orderDisplayID)
	return nil
}

func provisioningJobSmokeStatusOK(status string) bool {
	switch status {
	case "queued", "claimed", "running", "succeeded":
		return true
	default:
		return false
	}
}

func verifyAuditMutation(
	ctx context.Context,
	client *http.Client,
	baseURL string,
	check auditMutationCheck,
) error {
	query := url.Values{}
	query.Set("action", check.Action)
	query.Set("target_type", check.TargetType)
	query.Set("target_id", check.TargetID)
	records, err := doJSON[[]auditSummaryResponse](ctx, client, http.MethodGet, baseURL, "/admin/audit-logs?"+query.Encode(), adminHeaders(), nil, http.StatusOK)
	if err != nil {
		return err
	}
	if len(records) == 0 {
		return fmt.Errorf("expected audit log for %s %s", check.Action, check.TargetID)
	}
	record := records[0]
	if record.TenantID != demoTenantID || record.Action != check.Action || record.TargetType != check.TargetType || record.TargetID != check.TargetID {
		return fmt.Errorf("unexpected audit list record: %+v", record)
	}
	detail, err := doJSON[auditDetailResponse](ctx, client, http.MethodGet, baseURL, "/admin/audit-logs/"+record.ID, adminHeaders(), nil, http.StatusOK)
	if err != nil {
		return err
	}
	if detail.TenantID != demoTenantID || detail.Action != check.Action || detail.TargetType != check.TargetType || detail.TargetID != check.TargetID || detail.CorrelationID == "" {
		return fmt.Errorf("unexpected audit detail record: %+v", detail)
	}
	if check.MetadataContains != "" && !strings.Contains(string(detail.MetadataRedacted), check.MetadataContains) {
		return fmt.Errorf("audit detail %s missing metadata %q: %s", check.Action, check.MetadataContains, string(detail.MetadataRedacted))
	}
	if check.AfterContains != "" && !strings.Contains(string(detail.AfterSnapshotRedacted), check.AfterContains) {
		return fmt.Errorf("audit detail %s missing after snapshot %q: %s", check.Action, check.AfterContains, string(detail.AfterSnapshotRedacted))
	}
	fmt.Printf("audit check passed: %s %s (%d)\n", check.Action, detail.ID, detail.DisplayID)
	return nil
}
