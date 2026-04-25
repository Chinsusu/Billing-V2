package invoice

import (
	"errors"
	"strings"
	"testing"

	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/order"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

func TestBuildListInvoicesQueryAddsTenantBuyerAndStatusFilters(t *testing.T) {
	query, args, err := buildListInvoicesQuery(InvoiceFilter{
		TenantID:       tenant.ID("tenant-1"),
		BuyerUserID:    identity.UserID("buyer-1"),
		BuyerDisplayID: 10002,
		DisplayID:      44001,
		OrderID:        order.OrderID("order-1"),
		OrderDisplayID: 30001,
		Status:         StatusIssued,
		AmountMinMinor: int64Ptr(100),
		AmountMaxMinor: int64Ptr(900),
		Limit:          25,
	})
	if err != nil {
		t.Fatalf("expected query: %v", err)
	}
	for _, clause := range []string{
		"inv.tenant_id = $1",
		"buyer_display_id",
		"order_display_id",
		"inv.buyer_user_id = $2",
		"buyer.user_id = inv.buyer_user_id",
		"buyer.display_id = $3",
		"inv.display_id = $4",
		"inv.order_id = $5",
		"ord.order_id = inv.order_id",
		"ord.display_id = $6",
		"inv.status = $7",
		"inv.total_minor >= $8",
		"inv.total_minor <= $9",
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

func TestBuildListInvoicesQueryDefaultsLimit(t *testing.T) {
	query, args, err := buildListInvoicesQuery(InvoiceFilter{TenantID: tenant.ID("tenant-1")})
	if err != nil {
		t.Fatalf("expected query: %v", err)
	}
	if !strings.Contains(query, "LIMIT $2") {
		t.Fatalf("expected default limit placeholder: %s", query)
	}
	if len(args) != 2 || args[1] != defaultInvoiceListLimit {
		t.Fatalf("unexpected args: %#v", args)
	}
}

func TestBuildListInvoicesQueryRejectsBadStatus(t *testing.T) {
	_, _, err := buildListInvoicesQuery(InvoiceFilter{
		TenantID: tenant.ID("tenant-1"),
		Status:   Status("bad"),
	})
	if !errors.Is(err, ErrStatusInvalid) {
		t.Fatalf("expected status error, got %v", err)
	}
}

func TestBuildGetInvoiceQueryAddsBuyerScope(t *testing.T) {
	query, args, err := buildGetInvoiceQuery(InvoiceLookup{
		ID:          InvoiceID("invoice-1"),
		TenantID:    tenant.ID("tenant-1"),
		BuyerUserID: identity.UserID("buyer-1"),
	})
	if err != nil {
		t.Fatalf("expected query: %v", err)
	}
	for _, clause := range []string{"inv.invoice_id = $1", "inv.tenant_id = $2", "inv.buyer_user_id = $3"} {
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
