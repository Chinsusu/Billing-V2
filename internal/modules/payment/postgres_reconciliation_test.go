package payment

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/invoice"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
	"github.com/Chinsusu/Billing-V2/internal/modules/wallet"
)

func TestBuildListPaymentReconciliationsQueryAddsFilters(t *testing.T) {
	createdFrom := time.Date(2026, 4, 23, 0, 0, 0, 0, time.UTC)
	createdTo := createdFrom.Add(24 * time.Hour)
	query, args, err := buildListPaymentReconciliationsQuery(ReconciliationFilter{
		TenantID:    tenant.ID("tenant-1"),
		Status:      TransactionStatusPosted,
		Provider:    "wallet",
		InvoiceID:   invoice.InvoiceID("invoice-1"),
		WalletID:    wallet.WalletID("wallet-1"),
		CreatedFrom: createdFrom,
		CreatedTo:   createdTo,
		Limit:       25,
	})
	if err != nil {
		t.Fatalf("expected query: %v", err)
	}
	for _, clause := range []string{
		"LEFT JOIN LATERAL",
		"txn.tenant_id = $1",
		"txn.status = $2",
		"txn.metadata ->> 'provider' = $3",
		"txn.invoice_id = $4",
		"ledger.wallet_id = $5",
		"txn.created_at >= $6",
		"txn.created_at <= $7",
		"LIMIT $8",
	} {
		if !strings.Contains(query, clause) {
			t.Fatalf("expected %q in query: %s", clause, query)
		}
	}
	if len(args) != 8 || args[7] != 25 {
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
