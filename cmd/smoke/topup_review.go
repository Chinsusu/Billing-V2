package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	platformdb "github.com/Chinsusu/Billing-V2/internal/platform/db"
)

const (
	topupReviewApproveAmount = int64(111)
	topupReviewRejectAmount  = int64(222)
)

type topupReviewBaseline struct {
	WalletBalanceMinor int64
	OrderCount         int
	ProviderJobCount   int
	ServiceCount       int
}

type topupReviewFixture struct {
	ClientID        string
	WalletID        string
	WalletDisplayID int64
}

type topupLedgerEvidence struct {
	Count        int
	DisplayID    int64
	AmountMinor  int64
	EntryIDMatch bool
}

type topupAuditEvidence struct {
	Count     int
	DisplayID int64
}

func runDevTopupReviewSmoke(dsn string, baseURL string, timeout time.Duration) error {
	if err := guardDevEnvironment(); err != nil {
		return err
	}
	if dsn == "" {
		return fmt.Errorf("DB_DSN or -dsn is required for dev-topup-review smoke")
	}
	if baseURL == "" {
		return fmt.Errorf("API_BASE_URL or -base-url is required for dev-topup-review smoke")
	}
	if _, err := normalizedAPIURL(baseURL, "/healthz"); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn, err := platformdb.Open(ctx, platformdb.Config{DriverName: platformdb.DefaultDriverName, DSN: dsn})
	if err != nil {
		return fmt.Errorf("open smoke DB for top-up review checks: %w", err)
	}
	defer conn.Close()

	scenario := newBillingMutationScenario()
	fixture, err := createTopupReviewFixture(ctx, conn, scenario)
	if err != nil {
		return err
	}
	baseline, err := readTopupReviewBaseline(ctx, conn, fixture.WalletID)
	if err != nil {
		return err
	}
	client := &http.Client{Timeout: timeout}

	approved, err := createAndApproveTopup(ctx, client, baseURL, scenario, fixture)
	if err != nil {
		return err
	}
	approvedLedger, err := verifyTopupLedgerCredit(ctx, conn, approved, fixture)
	if err != nil {
		return err
	}
	approvedAudit, err := verifyTopupReviewAudit(ctx, conn, approved, auditActionTopup)
	if err != nil {
		return err
	}

	rejected, err := createAndRejectTopup(ctx, client, baseURL, scenario, fixture)
	if err != nil {
		return err
	}
	if err := verifyRejectedTopupHasNoLedger(ctx, conn, rejected, fixture); err != nil {
		return err
	}
	rejectedAudit, err := verifyTopupReviewAudit(ctx, conn, rejected, "wallet.topup.rejected")
	if err != nil {
		return err
	}

	finalState, err := readTopupReviewBaseline(ctx, conn, fixture.WalletID)
	if err != nil {
		return err
	}
	if err := verifyTopupReviewFinalState(baseline, finalState); err != nil {
		return err
	}

	fmt.Printf("topup review smoke passed: wallet_display_id=%d approve_topup_display_id=%d approve_ledger_display_id=%d approve_audit_display_id=%d reject_topup_display_id=%d reject_ledger_count=0 reject_audit_display_id=%d wallet_balance_delta_minor=%d provider_side_effects=none\n",
		fixture.WalletDisplayID,
		approved.DisplayID,
		approvedLedger.DisplayID,
		approvedAudit.DisplayID,
		rejected.DisplayID,
		rejectedAudit.DisplayID,
		finalState.WalletBalanceMinor-baseline.WalletBalanceMinor,
	)
	fmt.Println("Top-up review smoke output intentionally excludes DB_DSN, raw backend IDs, provider payloads, and credentials.")
	return nil
}

func createAndApproveTopup(ctx context.Context, client *http.Client, baseURL string, scenario billingMutationScenario, fixture topupReviewFixture) (topupResponse, error) {
	headers := actorHeaders(fixture.ClientID, "client")
	headers["Idempotency-Key"] = "smoke-topup-review-approve-" + scenario.RunID
	created, err := doJSON[topupResponse](ctx, client, http.MethodPost, baseURL, "/client/topup-requests", headers, topupRequestBody{
		WalletID:         fixture.WalletID,
		AmountMinor:      topupReviewApproveAmount,
		Currency:         smokeOrderCurrency,
		PaymentMethod:    "manual",
		PaymentReference: "SMOKE-TOPUP-REVIEW-APPROVE-" + scenario.RunID,
	}, http.StatusCreated)
	if err != nil {
		return topupResponse{}, err
	}
	if created.DisplayID <= 0 || created.Status != "submitted" {
		return topupResponse{}, fmt.Errorf("expected submitted approve top-up with display id, got display=%d status=%s", created.DisplayID, created.Status)
	}

	approved, err := doJSON[topupResponse](ctx, client, http.MethodPost, baseURL, "/reseller/topup-requests/"+created.ID+"/approve", resellerHeaders(), reviewTopupBody{
		ReviewNote: "Top-up review smoke approval " + scenario.RunID,
	}, http.StatusOK)
	if err != nil {
		return topupResponse{}, err
	}
	if approved.Status != "approved" || approved.DisplayID != created.DisplayID || approved.LedgerEntryID == "" {
		return topupResponse{}, fmt.Errorf("expected approved top-up with ledger entry, got display=%d status=%s ledger_present=%t", approved.DisplayID, approved.Status, approved.LedgerEntryID != "")
	}
	return approved, nil
}

func createAndRejectTopup(ctx context.Context, client *http.Client, baseURL string, scenario billingMutationScenario, fixture topupReviewFixture) (topupResponse, error) {
	headers := actorHeaders(fixture.ClientID, "client")
	headers["Idempotency-Key"] = "smoke-topup-review-reject-" + scenario.RunID
	created, err := doJSON[topupResponse](ctx, client, http.MethodPost, baseURL, "/client/topup-requests", headers, topupRequestBody{
		WalletID:         fixture.WalletID,
		AmountMinor:      topupReviewRejectAmount,
		Currency:         smokeOrderCurrency,
		PaymentMethod:    "manual",
		PaymentReference: "SMOKE-TOPUP-REVIEW-REJECT-" + scenario.RunID,
	}, http.StatusCreated)
	if err != nil {
		return topupResponse{}, err
	}
	if created.DisplayID <= 0 || created.Status != "submitted" {
		return topupResponse{}, fmt.Errorf("expected submitted reject top-up with display id, got display=%d status=%s", created.DisplayID, created.Status)
	}

	rejected, err := doJSON[topupResponse](ctx, client, http.MethodPost, baseURL, "/reseller/topup-requests/"+created.ID+"/reject", resellerHeaders(), reviewTopupBody{
		ReviewNote: "Top-up review smoke rejection " + scenario.RunID,
	}, http.StatusOK)
	if err != nil {
		return topupResponse{}, err
	}
	if rejected.Status != "rejected" || rejected.DisplayID != created.DisplayID || rejected.LedgerEntryID != "" {
		return topupResponse{}, fmt.Errorf("expected rejected top-up without ledger entry, got display=%d status=%s ledger_present=%t", rejected.DisplayID, rejected.Status, rejected.LedgerEntryID != "")
	}
	return rejected, nil
}

func createTopupReviewFixture(ctx context.Context, conn *sql.DB, scenario billingMutationScenario) (topupReviewFixture, error) {
	fixture := topupReviewFixture{}
	email := "smoke-topup-review-" + scenario.RunID + "@local.billing"
	err := conn.QueryRowContext(ctx, `
WITH inserted_user AS (
  INSERT INTO users (tenant_id, email, email_verified_at, password_hash, full_name, user_type, status)
  VALUES ($1, $2, NOW(), 'smoke-disabled-password-hash', 'Top-up Review Smoke', 'client', 'active')
  RETURNING user_id
),
assigned_role AS (
  INSERT INTO user_roles (user_id, tenant_id, role_id)
  SELECT inserted_user.user_id, $1, role.role_id
  FROM inserted_user
  JOIN roles role ON role.role_key = 'customer_catalog_viewer' AND role.is_system = TRUE
  ON CONFLICT (user_id, tenant_id, role_id) DO NOTHING
),
inserted_wallet AS (
  INSERT INTO wallets (tenant_id, owner_type, owner_id, currency, status, available_balance_minor, locked_balance_minor, metadata)
  SELECT $1, 'user', inserted_user.user_id, 'USD', 'active', 0, 0, jsonb_build_object('source', 'topup_review_smoke')
  FROM inserted_user
  RETURNING wallet_id, display_id
)
SELECT inserted_user.user_id, inserted_wallet.wallet_id, inserted_wallet.display_id
FROM inserted_user
CROSS JOIN inserted_wallet`,
		demoTenantID,
		email,
	).Scan(&fixture.ClientID, &fixture.WalletID, &fixture.WalletDisplayID)
	if err != nil {
		return topupReviewFixture{}, fmt.Errorf("create top-up review fixture: %w", err)
	}
	return fixture, nil
}

func readTopupReviewBaseline(ctx context.Context, conn *sql.DB, walletID string) (topupReviewBaseline, error) {
	var state topupReviewBaseline
	err := conn.QueryRowContext(ctx, `
SELECT
  COALESCE((SELECT available_balance_minor FROM wallets WHERE wallet_id = $1), 0),
  (SELECT COUNT(*) FROM orders),
  (SELECT COUNT(*) FROM jobs WHERE job_type = 'provider.provision'),
  (SELECT COUNT(*) FROM service_instances)`,
		walletID,
	).Scan(&state.WalletBalanceMinor, &state.OrderCount, &state.ProviderJobCount, &state.ServiceCount)
	if err != nil {
		return topupReviewBaseline{}, fmt.Errorf("read top-up review baseline: %w", err)
	}
	return state, nil
}

func verifyTopupLedgerCredit(ctx context.Context, conn *sql.DB, topup topupResponse, fixture topupReviewFixture) (topupLedgerEvidence, error) {
	var evidence topupLedgerEvidence
	err := conn.QueryRowContext(ctx, `
SELECT COUNT(*), COALESCE(MAX(display_id), 0), COALESCE(MAX(amount_minor), 0), COALESCE(BOOL_OR(ledger_entry_id = $2::uuid), false)
FROM wallet_ledger_entries
WHERE tenant_id = $1
  AND wallet_id = $3
  AND reference_type = 'topup_request'
  AND reference_id = $4::uuid
  AND entry_type = 'topup'
  AND direction = 'credit'
  AND status = 'posted'`,
		demoTenantID,
		topup.LedgerEntryID,
		fixture.WalletID,
		topup.ID,
	).Scan(&evidence.Count, &evidence.DisplayID, &evidence.AmountMinor, &evidence.EntryIDMatch)
	if err != nil {
		return topupLedgerEvidence{}, fmt.Errorf("verify approved top-up ledger: %w", err)
	}
	if evidence.Count != 1 || evidence.DisplayID <= 0 || evidence.AmountMinor != topupReviewApproveAmount || !evidence.EntryIDMatch {
		return topupLedgerEvidence{}, fmt.Errorf("expected one approved top-up ledger credit, got count=%d display=%d amount=%d ledger_match=%t", evidence.Count, evidence.DisplayID, evidence.AmountMinor, evidence.EntryIDMatch)
	}
	return evidence, nil
}

func verifyRejectedTopupHasNoLedger(ctx context.Context, conn *sql.DB, topup topupResponse, fixture topupReviewFixture) error {
	var count int
	err := conn.QueryRowContext(ctx, `
SELECT COUNT(*)
FROM wallet_ledger_entries
WHERE tenant_id = $1
  AND wallet_id = $2
  AND reference_type = 'topup_request'
  AND reference_id = $3::uuid`,
		demoTenantID,
		fixture.WalletID,
		topup.ID,
	).Scan(&count)
	if err != nil {
		return fmt.Errorf("verify rejected top-up ledger absence: %w", err)
	}
	if count != 0 {
		return fmt.Errorf("expected rejected top-up to have no ledger credit, got %d ledger row(s)", count)
	}
	return nil
}

func verifyTopupReviewAudit(ctx context.Context, conn *sql.DB, topup topupResponse, action string) (topupAuditEvidence, error) {
	var evidence topupAuditEvidence
	err := conn.QueryRowContext(ctx, `
SELECT COUNT(*), COALESCE(MAX(display_id), 0)
FROM audit_logs
WHERE tenant_id = $1
  AND action = $2
  AND target_type = 'topup_request'
  AND target_id = $3::uuid`,
		demoTenantID,
		action,
		topup.ID,
	).Scan(&evidence.Count, &evidence.DisplayID)
	if err != nil {
		return topupAuditEvidence{}, fmt.Errorf("verify top-up review audit %s: %w", action, err)
	}
	if evidence.Count != 1 || evidence.DisplayID <= 0 {
		return topupAuditEvidence{}, fmt.Errorf("expected one %s audit log, got count=%d display=%d", action, evidence.Count, evidence.DisplayID)
	}
	return evidence, nil
}

func verifyTopupReviewFinalState(before topupReviewBaseline, after topupReviewBaseline) error {
	if after.WalletBalanceMinor-before.WalletBalanceMinor != topupReviewApproveAmount {
		return fmt.Errorf("expected wallet balance delta %d, got %d", topupReviewApproveAmount, after.WalletBalanceMinor-before.WalletBalanceMinor)
	}
	if after.OrderCount != before.OrderCount || after.ProviderJobCount != before.ProviderJobCount || after.ServiceCount != before.ServiceCount {
		return fmt.Errorf(
			"expected no order/provider/service side effects, got orders %d->%d provider_jobs %d->%d services %d->%d",
			before.OrderCount,
			after.OrderCount,
			before.ProviderJobCount,
			after.ProviderJobCount,
			before.ServiceCount,
			after.ServiceCount,
		)
	}
	return nil
}
