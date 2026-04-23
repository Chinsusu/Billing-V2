package wallet

import (
	"context"

	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

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

type Store interface {
	CreateLedgerEntry(ctx context.Context, input CreateLedgerEntryInput) (LedgerEntry, error)
	ListLedgerEntries(ctx context.Context, filter LedgerEntryFilter) ([]LedgerEntry, error)
	GetLedgerEntry(ctx context.Context, lookup LedgerEntryLookup) (LedgerEntry, error)
}
