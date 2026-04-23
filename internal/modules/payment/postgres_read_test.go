package payment

import (
	"errors"
	"strings"
	"testing"

	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/invoice"
	"github.com/Chinsusu/Billing-V2/internal/modules/order"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

func TestBuildListTransactionsQueryAddsAccountScopeAndFilters(t *testing.T) {
	query, args, err := buildListTransactionsQuery(TransactionFilter{
		TenantID:       tenant.ID("tenant-1"),
		AccountUserID:  identity.UserID("account-1"),
		DisplayID:      51001,
		OrderID:        order.OrderID("order-1"),
		InvoiceID:      invoice.InvoiceID("invoice-1"),
		Type:           TransactionTypeCharge,
		Status:         TransactionStatusPosted,
		AmountMinMinor: int64Ptr(100),
		AmountMaxMinor: int64Ptr(900),
		Limit:          25,
	})
	if err != nil {
		t.Fatalf("expected query: %v", err)
	}
	for _, clause := range []string{
		"txn.tenant_id = $1",
		"txn.account_user_id = $2",
		"txn.display_id = $3",
		"txn.order_id = $4",
		"txn.invoice_id = $5",
		"txn.transaction_type = $6",
		"txn.status = $7",
		"txn.amount_minor >= $8",
		"txn.amount_minor <= $9",
		"LIMIT $10",
	} {
		if !strings.Contains(query, clause) {
			t.Fatalf("expected %q in query: %s", clause, query)
		}
	}
	if len(args) != 10 || args[9] != 25 {
		t.Fatalf("unexpected args: %#v", args)
	}
}

func TestBuildListTransactionsQueryDefaultsLimit(t *testing.T) {
	query, args, err := buildListTransactionsQuery(TransactionFilter{TenantID: tenant.ID("tenant-1")})
	if err != nil {
		t.Fatalf("expected query: %v", err)
	}
	if !strings.Contains(query, "LIMIT $2") {
		t.Fatalf("expected default limit placeholder: %s", query)
	}
	if len(args) != 2 || args[1] != defaultTransactionListLimit {
		t.Fatalf("unexpected args: %#v", args)
	}
}

func TestBuildListTransactionsQueryRejectsBadType(t *testing.T) {
	_, _, err := buildListTransactionsQuery(TransactionFilter{
		TenantID: tenant.ID("tenant-1"),
		Type:     TransactionType("bad"),
	})
	if !errors.Is(err, ErrTypeInvalid) {
		t.Fatalf("expected type error, got %v", err)
	}
}

func TestBuildGetTransactionQueryAddsAccountScope(t *testing.T) {
	query, args, err := buildGetTransactionQuery(TransactionLookup{
		ID:            TransactionID("txn-1"),
		TenantID:      tenant.ID("tenant-1"),
		AccountUserID: identity.UserID("account-1"),
	})
	if err != nil {
		t.Fatalf("expected query: %v", err)
	}
	for _, clause := range []string{"txn.payment_transaction_id = $1", "txn.tenant_id = $2", "txn.account_user_id = $3"} {
		if !strings.Contains(query, clause) {
			t.Fatalf("expected %q in query: %s", clause, query)
		}
	}
	if len(args) != 3 {
		t.Fatalf("unexpected args: %#v", args)
	}
}

func TestBuildGetTransactionQuerySupportsIdempotencyLookup(t *testing.T) {
	query, args, err := buildGetTransactionQuery(TransactionLookup{
		TenantID:       tenant.ID("tenant-1"),
		IdempotencyKey: IdempotencyKey(" key-1 "),
	})
	if err != nil {
		t.Fatalf("expected query: %v", err)
	}
	for _, clause := range []string{"txn.tenant_id = $1", "txn.idempotency_key = $2"} {
		if !strings.Contains(query, clause) {
			t.Fatalf("expected %q in query: %s", clause, query)
		}
	}
	if len(args) != 2 || args[1] != IdempotencyKey("key-1") {
		t.Fatalf("unexpected args: %#v", args)
	}
}

func int64Ptr(value int64) *int64 {
	return &value
}
