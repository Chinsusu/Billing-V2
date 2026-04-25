package wallet

import (
	"errors"
	"strings"
	"testing"

	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

func TestBuildListWalletsQueryAddsOwnerScopeAndFilters(t *testing.T) {
	query, args, err := buildListWalletsQuery(WalletFilter{
		TenantID:  tenant.ID("tenant-1"),
		DisplayID: 70001,
		OwnerType: OwnerTypeUser,
		OwnerID:   OwnerID("user-1"),
		Status:    StatusActive,
		Limit:     25,
	})
	if err != nil {
		t.Fatalf("expected query: %v", err)
	}
	for _, clause := range []string{
		"wallet.tenant_id = $1",
		"owner_display_id",
		"wallet.display_id = $2",
		"wallet.owner_type = $3",
		"wallet.owner_id = $4",
		"wallet.status = $5",
		"LIMIT $6",
	} {
		if !strings.Contains(query, clause) {
			t.Fatalf("expected %q in query: %s", clause, query)
		}
	}
	if len(args) != 6 || args[5] != 25 {
		t.Fatalf("unexpected args: %#v", args)
	}
}

func TestBuildGetWalletQueryAddsOwnerScope(t *testing.T) {
	query, args, err := buildGetWalletQuery(WalletLookup{
		ID:        WalletID("wallet-1"),
		TenantID:  tenant.ID("tenant-1"),
		OwnerType: OwnerTypeUser,
		OwnerID:   OwnerID("user-1"),
	})
	if err != nil {
		t.Fatalf("expected query: %v", err)
	}
	for _, clause := range []string{"owner_display_id", "wallet.wallet_id = $1", "wallet.tenant_id = $2", "wallet.owner_type = $3", "wallet.owner_id = $4"} {
		if !strings.Contains(query, clause) {
			t.Fatalf("expected %q in query: %s", clause, query)
		}
	}
	if len(args) != 4 {
		t.Fatalf("unexpected args: %#v", args)
	}
}

func TestBuildListLedgerEntriesQueryAddsWalletScopeAndFilters(t *testing.T) {
	query, args, err := buildListLedgerEntriesQuery(LedgerEntryFilter{
		TenantID:       tenant.ID("tenant-1"),
		WalletID:       WalletID("wallet-1"),
		DisplayID:      71001,
		Direction:      DirectionCredit,
		EntryType:      EntryTypeTopup,
		Status:         LedgerStatusPosted,
		AmountMinMinor: int64Ptr(100),
		AmountMaxMinor: int64Ptr(900),
		Limit:          25,
	})
	if err != nil {
		t.Fatalf("expected query: %v", err)
	}
	for _, clause := range []string{
		"entry.tenant_id = $1",
		"entry.wallet_id = $2",
		"entry.display_id = $3",
		"entry.direction = $4",
		"entry.entry_type = $5",
		"entry.status = $6",
		"entry.amount_minor >= $7",
		"entry.amount_minor <= $8",
		"LIMIT $9",
	} {
		if !strings.Contains(query, clause) {
			t.Fatalf("expected %q in query: %s", clause, query)
		}
	}
	if len(args) != 9 || args[8] != 25 {
		t.Fatalf("unexpected args: %#v", args)
	}
}

func TestBuildListLedgerEntriesQueryDefaultsLimit(t *testing.T) {
	query, args, err := buildListLedgerEntriesQuery(LedgerEntryFilter{
		TenantID: tenant.ID("tenant-1"),
		WalletID: WalletID("wallet-1"),
	})
	if err != nil {
		t.Fatalf("expected query: %v", err)
	}
	if !strings.Contains(query, "LIMIT $3") {
		t.Fatalf("expected default limit placeholder: %s", query)
	}
	if len(args) != 3 || args[2] != defaultLedgerEntryListLimit {
		t.Fatalf("unexpected args: %#v", args)
	}
}

func TestBuildListLedgerEntriesQueryRejectsBadStatus(t *testing.T) {
	_, _, err := buildListLedgerEntriesQuery(LedgerEntryFilter{
		TenantID: tenant.ID("tenant-1"),
		WalletID: WalletID("wallet-1"),
		Status:   LedgerStatus("bad"),
	})
	if !errors.Is(err, ErrLedgerStatusInvalid) {
		t.Fatalf("expected status error, got %v", err)
	}
}

func TestBuildGetLedgerEntryQueryRequiresWalletScope(t *testing.T) {
	query, args, err := buildGetLedgerEntryQuery(LedgerEntryLookup{
		ID:       LedgerEntryID("entry-1"),
		TenantID: tenant.ID("tenant-1"),
		WalletID: WalletID("wallet-1"),
	})
	if err != nil {
		t.Fatalf("expected query: %v", err)
	}
	for _, clause := range []string{"entry.ledger_entry_id = $1", "entry.tenant_id = $2", "entry.wallet_id = $3"} {
		if !strings.Contains(query, clause) {
			t.Fatalf("expected %q in query: %s", clause, query)
		}
	}
	if len(args) != 3 {
		t.Fatalf("unexpected args: %#v", args)
	}
}

func int64Ptr(value int64) *int64 {
	return &value
}
