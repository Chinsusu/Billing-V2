package payment

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/invoice"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
	"github.com/Chinsusu/Billing-V2/internal/modules/wallet"
)

func TestBuildListPaymentReconciliationsQueryAddsFilters(t *testing.T) {
	createdFrom := time.Date(2026, 4, 23, 0, 0, 0, 0, time.UTC)
	createdTo := createdFrom.Add(24 * time.Hour)
	query, args, err := buildListPaymentReconciliationsQuery(ReconciliationFilter{
		TenantID:         tenant.ID("tenant-1"),
		AccountUserID:    identity.UserID("account-1"),
		DisplayID:        51001,
		Status:           TransactionStatusPosted,
		Provider:         "wallet",
		InvoiceID:        invoice.InvoiceID("invoice-1"),
		InvoiceDisplayID: 44001,
		WalletID:         wallet.WalletID("wallet-1"),
		WalletDisplayID:  41001,
		AmountMinMinor:   int64Ptr(100),
		AmountMaxMinor:   int64Ptr(900),
		CreatedFrom:      createdFrom,
		CreatedTo:        createdTo,
		Limit:            25,
	})
	if err != nil {
		t.Fatalf("expected query: %v", err)
	}
	for _, clause := range []string{
		"LEFT JOIN LATERAL",
		"txn.tenant_id = $1",
		"txn.account_user_id = $2",
		"txn.display_id = $3",
		"txn.status = $4",
		"txn.metadata ->> 'provider' = $5",
		"txn.invoice_id = $6",
		"inv.display_id = $7",
		"ledger.wallet_id = $8",
		"linked_wallet.display_id = $9",
		"txn.amount_minor >= $10",
		"txn.amount_minor <= $11",
		"txn.created_at >= $12",
		"txn.created_at <= $13",
		"LIMIT $14",
	} {
		if !strings.Contains(query, clause) {
			t.Fatalf("expected %q in query: %s", clause, query)
		}
	}
	if len(args) != 14 || args[13] != 25 {
		t.Fatalf("unexpected args: %#v", args)
	}
}

func TestBuildListPaymentReconciliationsRejectsBadWindow(t *testing.T) {
	_, _, err := buildListPaymentReconciliationsQuery(ReconciliationFilter{
		TenantID:    tenant.ID("tenant-1"),
		CreatedFrom: time.Date(2026, 4, 24, 0, 0, 0, 0, time.UTC),
		CreatedTo:   time.Date(2026, 4, 23, 0, 0, 0, 0, time.UTC),
	})
	if !errors.Is(err, ErrCreatedTimeWindowInvalid) {
		t.Fatalf("expected created time window error, got %v", err)
	}
}

func TestBuildGetPaymentReconciliationQueryUsesTenantAndTransaction(t *testing.T) {
	query, args, err := buildGetPaymentReconciliationQuery(ReconciliationLookup{
		TenantID:      tenant.ID("tenant-1"),
		TransactionID: TransactionID("txn-1"),
	})
	if err != nil {
		t.Fatalf("expected query: %v", err)
	}
	for _, clause := range []string{"txn.tenant_id = $1", "txn.payment_transaction_id = $2"} {
		if !strings.Contains(query, clause) {
			t.Fatalf("expected %q in query: %s", clause, query)
		}
	}
	if len(args) != 2 {
		t.Fatalf("unexpected args: %#v", args)
	}
}
