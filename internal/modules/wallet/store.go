package wallet

import (
	"context"

	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

type WalletFilter struct {
	TenantID  tenant.ID
	OwnerType OwnerType
	OwnerID   OwnerID
	Status    Status
	Limit     int
}

type WalletLookup struct {
	ID        WalletID
	TenantID  tenant.ID
	OwnerType OwnerType
	OwnerID   OwnerID
}

type LedgerEntryFilter struct {
	TenantID  tenant.ID
	WalletID  WalletID
	Direction Direction
	EntryType EntryType
	Status    LedgerStatus
	Limit     int
}

type LedgerEntryLookup struct {
	ID       LedgerEntryID
	TenantID tenant.ID
	WalletID WalletID
}

type TopupRequestFilter struct {
	TenantID      tenant.ID
	WalletID      WalletID
	RequestedBy   identity.UserID
	PaymentMethod PaymentMethod
	Status        TopupStatus
	Limit         int
}

type TopupRequestLookup struct {
	ID          TopupRequestID
	TenantID    tenant.ID
	RequestedBy identity.UserID
}

type Store interface {
	ListWallets(ctx context.Context, filter WalletFilter) ([]Wallet, error)
	GetWallet(ctx context.Context, lookup WalletLookup) (Wallet, error)
	CreateLedgerEntry(ctx context.Context, input CreateLedgerEntryInput) (LedgerEntry, error)
	PostLedgerEntry(ctx context.Context, input PostLedgerEntryInput) (LedgerEntry, error)
	ListLedgerEntries(ctx context.Context, filter LedgerEntryFilter) ([]LedgerEntry, error)
	GetLedgerEntry(ctx context.Context, lookup LedgerEntryLookup) (LedgerEntry, error)
	CreateTopupRequest(ctx context.Context, input CreateTopupRequestInput) (TopupRequest, error)
	ListTopupRequests(ctx context.Context, filter TopupRequestFilter) ([]TopupRequest, error)
	GetTopupRequest(ctx context.Context, lookup TopupRequestLookup) (TopupRequest, error)
	ApproveTopupRequest(ctx context.Context, input ApproveTopupRequestInput, ledgerEntryID LedgerEntryID) (TopupRequest, error)
	RejectTopupRequest(ctx context.Context, input RejectTopupRequestInput) (TopupRequest, error)
}

func UserOwnerID(userID identity.UserID) OwnerID {
	return OwnerID(userID)
}
