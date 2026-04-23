package wallet

import (
	"errors"
	"testing"

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
