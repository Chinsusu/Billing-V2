package wallet

import (
	"context"
	"errors"
	"testing"

	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

func TestServiceApproveTopupRequestCreditsWallet(t *testing.T) {
	store := &fakeTopupReviewStore{
		request: TopupRequest{
			ID: "00000000-0000-0000-0000-000000000001", DisplayID: 90001, TenantID: "tenant-1",
			WalletID: "wallet-1", AmountMinor: 1000, Currency: "USD", Status: TopupStatusSubmitted,
		},
		entry: LedgerEntry{ID: "entry-1"},
	}
	service := NewService(store)

	request, err := service.ApproveTopupRequest(context.Background(), ApproveTopupRequestInput{
		ID: "00000000-0000-0000-0000-000000000001", TenantID: tenant.ID("tenant-1"), ReviewedBy: identity.UserID("admin-1"),
	})
	if err != nil {
		t.Fatalf("expected approved request: %v", err)
	}
	if request.Status != TopupStatusApproved {
		t.Fatalf("expected approved status, got %+v", request)
	}
	if store.postCalls != 1 || store.approveCalls != 1 {
		t.Fatalf("expected post and approve, got post=%d approve=%d", store.postCalls, store.approveCalls)
	}
	if store.postInput.EntryType != EntryTypeTopup ||
		store.postInput.Direction != DirectionCredit ||
		store.postInput.AmountMinor != 1000 {
		t.Fatalf("unexpected ledger post input: %+v", store.postInput)
	}
}

func TestServiceApproveTopupRequestReturnsExistingApproved(t *testing.T) {
	store := &fakeTopupReviewStore{request: TopupRequest{ID: "topup-1", Status: TopupStatusApproved}}
	service := NewService(store)

	request, err := service.ApproveTopupRequest(context.Background(), ApproveTopupRequestInput{
		ID: "topup-1", TenantID: tenant.ID("tenant-1"), ReviewedBy: identity.UserID("admin-1"),
	})
	if err != nil {
		t.Fatalf("expected existing approved request: %v", err)
	}
	if request.Status != TopupStatusApproved {
		t.Fatalf("expected approved status, got %+v", request)
	}
	if store.postCalls != 0 || store.approveCalls != 0 {
		t.Fatalf("expected no duplicate post/update, got post=%d approve=%d", store.postCalls, store.approveCalls)
	}
}

func TestServiceRejectTopupRequestStoresReview(t *testing.T) {
	store := &fakeTopupReviewStore{request: TopupRequest{ID: "topup-1", Status: TopupStatusUnderReview}}
	service := NewService(store)

	request, err := service.RejectTopupRequest(context.Background(), RejectTopupRequestInput{
		ID: "topup-1", TenantID: tenant.ID("tenant-1"), ReviewedBy: identity.UserID("admin-1"), ReviewNote: "missing proof",
	})
	if err != nil {
		t.Fatalf("expected rejected request: %v", err)
	}
	if request.Status != TopupStatusRejected {
		t.Fatalf("expected rejected status, got %+v", request)
	}
	if store.postCalls != 0 || store.rejectCalls != 1 {
		t.Fatalf("expected only reject update, got post=%d reject=%d", store.postCalls, store.rejectCalls)
	}
}

func TestServiceRejectTopupRequestRejectsInvalidStatus(t *testing.T) {
	store := &fakeTopupReviewStore{request: TopupRequest{ID: "topup-1", Status: TopupStatusCancelled}}
	service := NewService(store)

	_, err := service.RejectTopupRequest(context.Background(), RejectTopupRequestInput{
		ID: "topup-1", TenantID: tenant.ID("tenant-1"), ReviewedBy: identity.UserID("admin-1"), ReviewNote: "cancelled",
	})
	if !errors.Is(err, ErrTopupStatusConflict) {
		t.Fatalf("expected status conflict, got %v", err)
	}
}

type fakeTopupReviewStore struct {
	request      TopupRequest
	entry        LedgerEntry
	postInput    PostLedgerEntryInput
	approveInput ApproveTopupRequestInput
	rejectInput  RejectTopupRequestInput
	postCalls    int
	approveCalls int
	rejectCalls  int
}

func (store *fakeTopupReviewStore) ListWallets(ctx context.Context, filter WalletFilter) ([]Wallet, error) {
	return nil, nil
}

func (store *fakeTopupReviewStore) GetWallet(ctx context.Context, lookup WalletLookup) (Wallet, error) {
	return Wallet{}, nil
}

func (store *fakeTopupReviewStore) CreateLedgerEntry(ctx context.Context, input CreateLedgerEntryInput) (LedgerEntry, error) {
	return LedgerEntry{}, nil
}

func (store *fakeTopupReviewStore) PostLedgerEntry(ctx context.Context, input PostLedgerEntryInput) (LedgerEntry, error) {
	store.postCalls++
	store.postInput = input
	return store.entry, nil
}

func (store *fakeTopupReviewStore) ListLedgerEntries(ctx context.Context, filter LedgerEntryFilter) ([]LedgerEntry, error) {
	return nil, nil
}

func (store *fakeTopupReviewStore) GetLedgerEntry(ctx context.Context, lookup LedgerEntryLookup) (LedgerEntry, error) {
	return LedgerEntry{}, nil
}

func (store *fakeTopupReviewStore) CreateTopupRequest(ctx context.Context, input CreateTopupRequestInput) (TopupRequest, error) {
	return TopupRequest{}, nil
}

func (store *fakeTopupReviewStore) ListTopupRequests(ctx context.Context, filter TopupRequestFilter) ([]TopupRequest, error) {
	return nil, nil
}

func (store *fakeTopupReviewStore) GetTopupRequest(ctx context.Context, lookup TopupRequestLookup) (TopupRequest, error) {
	return store.request, nil
}

func (store *fakeTopupReviewStore) ApproveTopupRequest(ctx context.Context, input ApproveTopupRequestInput, ledgerEntryID LedgerEntryID) (TopupRequest, error) {
	store.approveCalls++
	store.approveInput = input
	store.request.Status = TopupStatusApproved
	store.request.LedgerEntryID = ledgerEntryID
	return store.request, nil
}

func (store *fakeTopupReviewStore) RejectTopupRequest(ctx context.Context, input RejectTopupRequestInput) (TopupRequest, error) {
	store.rejectCalls++
	store.rejectInput = input
	store.request.Status = TopupStatusRejected
	return store.request, nil
}
