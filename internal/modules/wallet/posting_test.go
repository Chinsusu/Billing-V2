package wallet

import (
	"context"
	"errors"
	"testing"

	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

func TestPostLedgerEntryInputNormalizeValidate(t *testing.T) {
	input := PostLedgerEntryInput{
		WalletID:       WalletID("wallet-1"),
		TenantID:       tenant.ID("tenant-1"),
		Direction:      DirectionCredit,
		AmountMinor:    1000,
		Currency:       " usd ",
		EntryType:      EntryTypeTopup,
		ReferenceType:  ReferenceType(" topup_request "),
		ReferenceID:    ReferenceID("00000000-0000-0000-0000-000000000001"),
		IdempotencyKey: IdempotencyKey(" key-1 "),
		CreatedBy:      identity.UserID("admin-1"),
		CorrelationID:  CorrelationID("00000000-0000-0000-0000-000000000002"),
	}.Normalize()

	if input.Currency != "USD" ||
		input.ReferenceType != ReferenceType("topup_request") ||
		input.IdempotencyKey != IdempotencyKey("key-1") {
		t.Fatalf("unexpected normalized input: %+v", input)
	}
	if err := input.Validate(); err != nil {
		t.Fatalf("expected valid posting input: %v", err)
	}

	input.AmountMinor = 0
	if err := input.Validate(); !errors.Is(err, ErrAmountInvalid) {
		t.Fatalf("expected amount error, got %v", err)
	}
}

func TestPostLedgerEntryInputRequiresAdjustmentActor(t *testing.T) {
	err := PostLedgerEntryInput{
		WalletID:       WalletID("wallet-1"),
		TenantID:       tenant.ID("tenant-1"),
		Direction:      DirectionCredit,
		AmountMinor:    100,
		Currency:       "USD",
		EntryType:      EntryTypeAdjustment,
		ReferenceType:  ReferenceType("manual_adjustment"),
		ReferenceID:    ReferenceID("adjustment-1"),
		IdempotencyKey: IdempotencyKey("ledger-key-1"),
		Reason:         "finance correction",
		CorrelationID:  CorrelationID("00000000-0000-0000-0000-000000000001"),
	}.Normalize().Validate()
	if !errors.Is(err, identity.ErrActorIDMissing) {
		t.Fatalf("expected actor error, got %v", err)
	}
}

func TestServicePostLedgerEntryCallsStore(t *testing.T) {
	store := &fakePostingStore{entry: LedgerEntry{ID: "entry-1"}}
	service := NewService(store)

	entry, err := service.PostLedgerEntry(context.Background(), PostLedgerEntryInput{
		WalletID:       WalletID("wallet-1"),
		TenantID:       tenant.ID("tenant-1"),
		Direction:      DirectionCredit,
		AmountMinor:    1000,
		Currency:       " usd ",
		EntryType:      EntryTypeTopup,
		ReferenceType:  ReferenceType("topup_request"),
		ReferenceID:    ReferenceID("00000000-0000-0000-0000-000000000001"),
		IdempotencyKey: IdempotencyKey(" key-1 "),
		CorrelationID:  CorrelationID("00000000-0000-0000-0000-000000000002"),
	})
	if err != nil {
		t.Fatalf("expected posted entry: %v", err)
	}
	if entry.ID != LedgerEntryID("entry-1") || store.postCalls != 1 {
		t.Fatalf("expected store result, got entry=%+v calls=%d", entry, store.postCalls)
	}
	if store.input.Currency != "USD" || store.input.IdempotencyKey != IdempotencyKey("key-1") {
		t.Fatalf("expected normalized input, got %+v", store.input)
	}
}

func TestServicePostLedgerEntryRejectsInvalidInput(t *testing.T) {
	store := &fakePostingStore{}
	service := NewService(store)

	_, err := service.PostLedgerEntry(context.Background(), PostLedgerEntryInput{})
	if !errors.Is(err, ErrWalletIDMissing) {
		t.Fatalf("expected wallet id error, got %v", err)
	}
	if store.postCalls != 0 {
		t.Fatalf("expected no store call, got %d", store.postCalls)
	}
}

type fakePostingStore struct {
	entry     LedgerEntry
	input     PostLedgerEntryInput
	postCalls int
}

func (store *fakePostingStore) ListWallets(ctx context.Context, filter WalletFilter) ([]Wallet, error) {
	return nil, nil
}

func (store *fakePostingStore) GetWallet(ctx context.Context, lookup WalletLookup) (Wallet, error) {
	return Wallet{}, nil
}

func (store *fakePostingStore) CreateLedgerEntry(ctx context.Context, input CreateLedgerEntryInput) (LedgerEntry, error) {
	return LedgerEntry{}, nil
}

func (store *fakePostingStore) PostLedgerEntry(ctx context.Context, input PostLedgerEntryInput) (LedgerEntry, error) {
	store.postCalls++
	store.input = input
	return store.entry, nil
}

func (store *fakePostingStore) PostLedgerEntryResult(ctx context.Context, input PostLedgerEntryInput) (PostLedgerEntryResult, error) {
	entry, err := store.PostLedgerEntry(ctx, input)
	return PostLedgerEntryResult{Entry: entry, Created: true}, err
}

func (store *fakePostingStore) ListLedgerEntries(ctx context.Context, filter LedgerEntryFilter) ([]LedgerEntry, error) {
	return nil, nil
}

func (store *fakePostingStore) GetLedgerEntry(ctx context.Context, lookup LedgerEntryLookup) (LedgerEntry, error) {
	return LedgerEntry{}, nil
}

func (store *fakePostingStore) CreateTopupRequest(ctx context.Context, input CreateTopupRequestInput) (TopupRequest, error) {
	return TopupRequest{}, nil
}

func (store *fakePostingStore) ListTopupRequests(ctx context.Context, filter TopupRequestFilter) ([]TopupRequest, error) {
	return nil, nil
}

func (store *fakePostingStore) GetTopupRequest(ctx context.Context, lookup TopupRequestLookup) (TopupRequest, error) {
	return TopupRequest{}, nil
}

func (store *fakePostingStore) ApproveTopupRequest(ctx context.Context, input ApproveTopupRequestInput, ledgerEntryID LedgerEntryID) (TopupRequest, error) {
	return TopupRequest{}, nil
}

func (store *fakePostingStore) RejectTopupRequest(ctx context.Context, input RejectTopupRequestInput) (TopupRequest, error) {
	return TopupRequest{}, nil
}
