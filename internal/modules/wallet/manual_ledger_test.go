package wallet

import (
	"context"
	"errors"
	"testing"

	"github.com/Chinsusu/Billing-V2/internal/modules/audit"
	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

func TestServiceCreateWalletRefundPostsCreditAndAudits(t *testing.T) {
	store := &fakeManualLedgerStore{wallet: Wallet{ID: "wallet-1", TenantID: "tenant-1", Currency: "USD", Status: StatusActive}}
	auditLog := &fakeWalletAuditAppender{}
	service := NewServiceWithAudit(store, auditLog)

	entry, err := service.CreateWalletRefund(context.Background(), walletRefundInput())
	if err != nil {
		t.Fatalf("expected refund ledger entry: %v", err)
	}
	if store.getWalletCalls != 1 || store.postCalls != 1 {
		t.Fatalf("expected wallet lookup and ledger post, got lookup=%d post=%d", store.getWalletCalls, store.postCalls)
	}
	if store.postInput.EntryType != EntryTypeRefund ||
		store.postInput.Direction != DirectionCredit ||
		store.postInput.CreatedBy != identity.UserID("admin-1") {
		t.Fatalf("unexpected refund post input: %+v", store.postInput)
	}
	if entry.ID != LedgerEntryID("ledger-1") {
		t.Fatalf("expected ledger entry result, got %+v", entry)
	}
	if auditLog.calls != 1 ||
		auditLog.input.Action != walletAuditActionRefundCreated ||
		auditLog.input.TargetID != audit.TargetID("ledger-1") ||
		auditLog.input.ActorID != audit.ActorID("admin-1") {
		t.Fatalf("unexpected audit input: %+v", auditLog.input)
	}
}

func TestServiceCreateWalletAdjustmentRejectsInsufficientContext(t *testing.T) {
	store := &fakeManualLedgerStore{wallet: Wallet{ID: "wallet-1", TenantID: "tenant-1", Currency: "USD", Status: StatusActive}}
	service := NewService(store)

	missingActor := walletAdjustmentInput()
	missingActor.CreatedBy = ""
	_, err := service.CreateWalletAdjustment(context.Background(), missingActor)
	if !errors.Is(err, identity.ErrActorIDMissing) {
		t.Fatalf("expected actor error, got %v", err)
	}

	missingReason := walletAdjustmentInput()
	missingReason.Reason = " "
	_, err = service.CreateWalletAdjustment(context.Background(), missingReason)
	if !errors.Is(err, ErrReasonMissing) {
		t.Fatalf("expected reason error, got %v", err)
	}
	if store.postCalls != 0 {
		t.Fatalf("expected no ledger post, got %d", store.postCalls)
	}
}

func TestServiceCreateWalletAdjustmentReturnsDuplicateWithoutAudit(t *testing.T) {
	input := walletAdjustmentInput()
	postInput := adjustmentPostInput(input.Normalize())
	store := &fakeManualLedgerStore{
		wallet: Wallet{ID: "wallet-1", TenantID: "tenant-1", Currency: "USD", Status: StatusActive},
		result: PostLedgerEntryResult{
			Entry:   manualLedgerEntryFromPost(postInput),
			Created: false,
		},
	}
	auditLog := &fakeWalletAuditAppender{}
	service := NewServiceWithAudit(store, auditLog)

	entry, err := service.CreateWalletAdjustment(context.Background(), input)
	if err != nil {
		t.Fatalf("expected duplicate adjustment result: %v", err)
	}
	if entry.ID != LedgerEntryID("ledger-1") || store.postCalls != 1 {
		t.Fatalf("expected existing ledger result, entry=%+v post=%d", entry, store.postCalls)
	}
	if auditLog.calls != 0 {
		t.Fatalf("expected no duplicate audit event, got %d", auditLog.calls)
	}
}

func TestServiceCreateWalletAdjustmentRejectsDuplicateConflict(t *testing.T) {
	input := walletAdjustmentInput()
	postInput := adjustmentPostInput(input.Normalize())
	conflictingEntry := manualLedgerEntryFromPost(postInput)
	conflictingEntry.AmountMinor = 9900
	store := &fakeManualLedgerStore{
		wallet: Wallet{ID: "wallet-1", TenantID: "tenant-1", Currency: "USD", Status: StatusActive},
		result: PostLedgerEntryResult{
			Entry:   conflictingEntry,
			Created: false,
		},
	}
	auditLog := &fakeWalletAuditAppender{}
	service := NewServiceWithAudit(store, auditLog)

	_, err := service.CreateWalletAdjustment(context.Background(), input)
	if !errors.Is(err, ErrIdempotencyConflict) {
		t.Fatalf("expected idempotency conflict, got %v", err)
	}
	if auditLog.calls != 0 {
		t.Fatalf("expected no audit on conflict, got %d", auditLog.calls)
	}
}

func walletRefundInput() CreateWalletRefundInput {
	return CreateWalletRefundInput{
		TenantID:       tenant.ID("tenant-1"),
		WalletID:       WalletID("wallet-1"),
		AmountMinor:    1200,
		Currency:       " usd ",
		ReferenceType:  ReferenceType(" invoice "),
		ReferenceID:    ReferenceID("00000000-0000-0000-0000-000000000101"),
		IdempotencyKey: IdempotencyKey(" refund-key-1 "),
		CreatedBy:      identity.UserID("admin-1"),
		Reason:         "support ticket RF-101",
		CorrelationID:  CorrelationID("00000000-0000-0000-0000-000000000201"),
	}
}

func walletAdjustmentInput() CreateWalletAdjustmentInput {
	return CreateWalletAdjustmentInput{
		TenantID:       tenant.ID("tenant-1"),
		WalletID:       WalletID("wallet-1"),
		Direction:      DirectionDebit,
		AmountMinor:    700,
		Currency:       " usd ",
		ReferenceType:  ReferenceType(" manual_adjustment "),
		ReferenceID:    ReferenceID("00000000-0000-0000-0000-000000000102"),
		IdempotencyKey: IdempotencyKey(" adjustment-key-1 "),
		CreatedBy:      identity.UserID("admin-1"),
		Reason:         "finance correction ADJ-102",
		CorrelationID:  CorrelationID("00000000-0000-0000-0000-000000000202"),
	}
}

type fakeManualLedgerStore struct {
	wallet         Wallet
	result         PostLedgerEntryResult
	postInput      PostLedgerEntryInput
	getWalletCalls int
	postCalls      int
}

func (store *fakeManualLedgerStore) ListWallets(ctx context.Context, filter WalletFilter) ([]Wallet, error) {
	return nil, nil
}

func (store *fakeManualLedgerStore) GetWallet(ctx context.Context, lookup WalletLookup) (Wallet, error) {
	store.getWalletCalls++
	return store.wallet, nil
}

func (store *fakeManualLedgerStore) CreateLedgerEntry(ctx context.Context, input CreateLedgerEntryInput) (LedgerEntry, error) {
	return LedgerEntry{}, nil
}

func (store *fakeManualLedgerStore) PostLedgerEntry(ctx context.Context, input PostLedgerEntryInput) (LedgerEntry, error) {
	result, err := store.PostLedgerEntryResult(ctx, input)
	return result.Entry, err
}

func (store *fakeManualLedgerStore) PostLedgerEntryResult(ctx context.Context, input PostLedgerEntryInput) (PostLedgerEntryResult, error) {
	store.postCalls++
	store.postInput = input
	if store.result.Entry.ID.Empty() {
		return PostLedgerEntryResult{Entry: manualLedgerEntryFromPost(input), Created: true}, nil
	}
	return store.result, nil
}

func (store *fakeManualLedgerStore) ListLedgerEntries(ctx context.Context, filter LedgerEntryFilter) ([]LedgerEntry, error) {
	return nil, nil
}

func (store *fakeManualLedgerStore) GetLedgerEntry(ctx context.Context, lookup LedgerEntryLookup) (LedgerEntry, error) {
	return LedgerEntry{}, nil
}

func (store *fakeManualLedgerStore) CreateTopupRequest(ctx context.Context, input CreateTopupRequestInput) (TopupRequest, error) {
	return TopupRequest{}, nil
}

func (store *fakeManualLedgerStore) ListTopupRequests(ctx context.Context, filter TopupRequestFilter) ([]TopupRequest, error) {
	return nil, nil
}

func (store *fakeManualLedgerStore) GetTopupRequest(ctx context.Context, lookup TopupRequestLookup) (TopupRequest, error) {
	return TopupRequest{}, nil
}

func (store *fakeManualLedgerStore) ApproveTopupRequest(ctx context.Context, input ApproveTopupRequestInput, ledgerEntryID LedgerEntryID) (TopupRequest, error) {
	return TopupRequest{}, nil
}

func (store *fakeManualLedgerStore) RejectTopupRequest(ctx context.Context, input RejectTopupRequestInput) (TopupRequest, error) {
	return TopupRequest{}, nil
}

func manualLedgerEntryFromPost(input PostLedgerEntryInput) LedgerEntry {
	return LedgerEntry{
		ID:                "ledger-1",
		WalletID:          input.WalletID,
		TenantID:          input.TenantID,
		Direction:         input.Direction,
		AmountMinor:       input.AmountMinor,
		Currency:          input.Currency,
		EntryType:         input.EntryType,
		Status:            LedgerStatusPosted,
		BalanceAfterMinor: 4300,
		ReferenceType:     input.ReferenceType,
		ReferenceID:       input.ReferenceID,
		IdempotencyKey:    input.IdempotencyKey,
		CreatedBy:         input.CreatedBy,
		Reason:            input.Reason,
		CorrelationID:     input.CorrelationID,
	}
}
