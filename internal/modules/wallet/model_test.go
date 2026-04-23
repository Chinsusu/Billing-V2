package wallet

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

func TestCreateWalletInputNormalizeValidate(t *testing.T) {
	input := CreateWalletInput{
		TenantID:              tenant.ID("tenant-1"),
		OwnerType:             OwnerTypeUser,
		OwnerID:               OwnerID(" user-1 "),
		Currency:              " usd ",
		AvailableBalanceMinor: 1000,
		LockedBalanceMinor:    200,
		Metadata:              json.RawMessage(`{"source":"seed"}`),
	}.Normalize()

	if input.OwnerID != OwnerID("user-1") {
		t.Fatalf("expected trimmed owner id, got %q", input.OwnerID)
	}
	if input.Currency != "USD" {
		t.Fatalf("expected normalized currency, got %q", input.Currency)
	}
	if input.Status != StatusActive {
		t.Fatalf("expected active status, got %q", input.Status)
	}
	if err := input.Validate(); err != nil {
		t.Fatalf("expected valid wallet input: %v", err)
	}

	input.AvailableBalanceMinor = -1
	if err := input.Validate(); !errors.Is(err, ErrBalanceInvalid) {
		t.Fatalf("expected balance error, got %v", err)
	}
}

func TestCreateWalletInputRequiresOwner(t *testing.T) {
	err := CreateWalletInput{
		TenantID: tenant.ID("tenant-1"),
		Currency: "USD",
	}.Normalize().Validate()
	if !errors.Is(err, ErrOwnerTypeInvalid) {
		t.Fatalf("expected owner type error, got %v", err)
	}

	err = CreateWalletInput{
		TenantID:  tenant.ID("tenant-1"),
		OwnerType: OwnerTypeTenant,
		Currency:  "USD",
	}.Normalize().Validate()
	if !errors.Is(err, ErrOwnerIDMissing) {
		t.Fatalf("expected owner id error, got %v", err)
	}
}

func TestWalletEnums(t *testing.T) {
	for _, ownerType := range []OwnerType{OwnerTypeTenant, OwnerTypeUser, OwnerTypeResellerSettlement, OwnerTypePlatform} {
		if !ownerType.Valid() {
			t.Fatalf("expected valid owner type %q", ownerType)
		}
	}
	if OwnerType("bad").Valid() {
		t.Fatal("unexpected valid owner type")
	}
	if !StatusActive.Valid() || !StatusFrozen.Valid() || !StatusClosed.Valid() {
		t.Fatal("expected core statuses to be valid")
	}
	if Status("bad").Valid() {
		t.Fatal("unexpected valid status")
	}
}

func TestCreateLedgerEntryInputNormalizeValidate(t *testing.T) {
	input := CreateLedgerEntryInput{
		WalletID:          WalletID("wallet-1"),
		TenantID:          tenant.ID("tenant-1"),
		Direction:         DirectionDebit,
		AmountMinor:       500,
		Currency:          " usd ",
		EntryType:         EntryTypePurchase,
		BalanceAfterMinor: 1500,
		ReferenceType:     ReferenceType(" order "),
		ReferenceID:       ReferenceID(" order-1 "),
		IdempotencyKey:    IdempotencyKey(" ledger-key-1 "),
		CorrelationID:     CorrelationID("00000000-0000-0000-0000-000000000001"),
	}.Normalize()

	if input.Currency != "USD" {
		t.Fatalf("expected normalized currency, got %q", input.Currency)
	}
	if input.Status != LedgerStatusPosted {
		t.Fatalf("expected posted status, got %q", input.Status)
	}
	if input.ReferenceType != ReferenceType("order") || input.ReferenceID != ReferenceID("order-1") {
		t.Fatalf("expected trimmed reference, got %+v", input)
	}
	if err := input.Validate(); err != nil {
		t.Fatalf("expected valid ledger input: %v", err)
	}

	input.AmountMinor = 0
	if err := input.Validate(); !errors.Is(err, ErrAmountInvalid) {
		t.Fatalf("expected amount error, got %v", err)
	}
}

func TestCreateLedgerEntryInputRequiresAdjustmentReason(t *testing.T) {
	err := CreateLedgerEntryInput{
		WalletID:          WalletID("wallet-1"),
		TenantID:          tenant.ID("tenant-1"),
		Direction:         DirectionCredit,
		AmountMinor:       100,
		Currency:          "USD",
		EntryType:         EntryTypeAdjustment,
		BalanceAfterMinor: 100,
		ReferenceType:     ReferenceType("manual_adjustment"),
		ReferenceID:       ReferenceID("adjustment-1"),
		IdempotencyKey:    IdempotencyKey("ledger-key-1"),
		CorrelationID:     CorrelationID("00000000-0000-0000-0000-000000000001"),
	}.Normalize().Validate()
	if !errors.Is(err, ErrReasonMissing) {
		t.Fatalf("expected reason error, got %v", err)
	}
}

func TestLedgerEnums(t *testing.T) {
	if !DirectionCredit.Valid() || !DirectionDebit.Valid() {
		t.Fatal("expected core directions to be valid")
	}
	if Direction("bad").Valid() {
		t.Fatal("unexpected valid direction")
	}
	for _, entryType := range []EntryType{EntryTypeTopup, EntryTypePurchase, EntryTypeResellerCost, EntryTypeRefund, EntryTypeAdjustment, EntryTypeReversal, EntryTypeCommission, EntryTypeLock, EntryTypeUnlock} {
		if !entryType.Valid() {
			t.Fatalf("expected valid entry type %q", entryType)
		}
	}
	if !LedgerStatusPosted.Valid() || !LedgerStatusVoidedByReversal.Valid() {
		t.Fatal("expected core ledger statuses to be valid")
	}
}
