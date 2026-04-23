package wallet

import (
	"errors"
	"strings"
	"testing"

	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

func TestCreateLedgerEntryArgsNormalizeValidate(t *testing.T) {
	args, err := createLedgerEntryArgs(CreateLedgerEntryInput{
		WalletID:          WalletID("wallet-1"),
		TenantID:          tenant.ID("tenant-1"),
		Direction:         DirectionCredit,
		AmountMinor:       1000,
		Currency:          " usd ",
		EntryType:         EntryTypeTopup,
		BalanceAfterMinor: 1000,
		ReferenceType:     ReferenceType(" topup_request "),
		ReferenceID:       ReferenceID(" request-1 "),
		IdempotencyKey:    IdempotencyKey(" key-1 "),
		CorrelationID:     CorrelationID("00000000-0000-0000-0000-000000000001"),
	})
	if err != nil {
		t.Fatalf("expected args: %v", err)
	}
	if args[4] != "USD" || args[6] != LedgerStatusPosted || args[8] != ReferenceType("topup_request") {
		t.Fatalf("unexpected normalized args: %#v", args)
	}
}

func TestCreateLedgerEntryArgsRejectsMissingCorrelation(t *testing.T) {
	_, err := createLedgerEntryArgs(CreateLedgerEntryInput{
		WalletID:          WalletID("wallet-1"),
		TenantID:          tenant.ID("tenant-1"),
		Direction:         DirectionCredit,
		AmountMinor:       1000,
		Currency:          "USD",
		EntryType:         EntryTypeTopup,
		BalanceAfterMinor: 1000,
		ReferenceType:     ReferenceType("topup_request"),
		ReferenceID:       ReferenceID("request-1"),
		IdempotencyKey:    IdempotencyKey("key-1"),
	})
	if !errors.Is(err, ErrCorrelationIDMissing) {
		t.Fatalf("expected correlation error, got %v", err)
	}
}

func TestPostLedgerEntryArgsNormalizeValidate(t *testing.T) {
	args, err := postLedgerEntryArgs(PostLedgerEntryInput{
		WalletID:       WalletID("wallet-1"),
		TenantID:       tenant.ID("tenant-1"),
		Direction:      DirectionDebit,
		AmountMinor:    500,
		Currency:       " usd ",
		EntryType:      EntryTypePurchase,
		ReferenceType:  ReferenceType(" order "),
		ReferenceID:    ReferenceID("00000000-0000-0000-0000-000000000001"),
		IdempotencyKey: IdempotencyKey(" debit-key-1 "),
		CreatedBy:      identity.UserID("buyer-1"),
		Reason:         " invoice payment ",
		CorrelationID:  CorrelationID("00000000-0000-0000-0000-000000000002"),
	})
	if err != nil {
		t.Fatalf("expected args: %v", err)
	}
	if args[4] != "USD" || args[6] != ReferenceType("order") || args[8] != IdempotencyKey("debit-key-1") {
		t.Fatalf("unexpected normalized args: %#v", args)
	}
}

func TestPostLedgerEntryArgsRejectsAdjustmentWithoutReason(t *testing.T) {
	_, err := postLedgerEntryArgs(PostLedgerEntryInput{
		WalletID:       WalletID("wallet-1"),
		TenantID:       tenant.ID("tenant-1"),
		Direction:      DirectionCredit,
		AmountMinor:    500,
		Currency:       "USD",
		EntryType:      EntryTypeAdjustment,
		ReferenceType:  ReferenceType("manual_adjustment"),
		ReferenceID:    ReferenceID("00000000-0000-0000-0000-000000000001"),
		IdempotencyKey: IdempotencyKey("adjust-key-1"),
		CorrelationID:  CorrelationID("00000000-0000-0000-0000-000000000002"),
	})
	if !errors.Is(err, ErrReasonMissing) {
		t.Fatalf("expected reason error, got %v", err)
	}
}

func TestPostLedgerEntrySQLUpdatesBalanceOnce(t *testing.T) {
	for _, clause := range []string{
		"UPDATE wallets wallet",
		"wallet.available_balance_minor + $4",
		"wallet.available_balance_minor - $4",
		"wallet.available_balance_minor >= $4",
		"NOT EXISTS (SELECT 1 FROM existing)",
		"ON CONFLICT (wallet_id, idempotency_key)",
	} {
		if !strings.Contains(postLedgerEntrySQL, clause) {
			t.Fatalf("expected %q in post ledger SQL", clause)
		}
	}
}

func TestTopupReviewSQLRestrictsReviewableStatuses(t *testing.T) {
	for _, query := range []string{approveTopupRequestSQL, rejectTopupRequestSQL} {
		if !strings.Contains(query, "status IN ('submitted', 'under_review')") {
			t.Fatalf("expected reviewable status guard in query: %s", query)
		}
	}
	if !strings.Contains(approveTopupRequestSQL, "ledger_entry_id = $5") {
		t.Fatalf("expected approve query to link ledger entry: %s", approveTopupRequestSQL)
	}
}
