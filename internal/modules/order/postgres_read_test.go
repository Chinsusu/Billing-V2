package order

import (
	"errors"
	"strings"
	"testing"

	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

func TestBuildListOrdersQueryAddsFilters(t *testing.T) {
	query, args, err := buildListOrdersQuery(OrderFilter{
		TenantID:       tenant.ID("tenant-1"),
		BuyerUserID:    identity.UserID("buyer-1"),
		DisplayID:      30001,
		OrderStatus:    OrderStatusPendingPayment,
		BillingStatus:  BillingStatusUnpaid,
		AmountMinMinor: int64Ptr(100),
		AmountMaxMinor: int64Ptr(900),
		Limit:          25,
	})
	if err != nil {
		t.Fatalf("expected query: %v", err)
	}
	for _, clause := range []string{
		"tenant_id = $1",
		"buyer_user_id = $2",
		"display_id = $3",
		"order_status = $4",
		"billing_status = $5",
		"total_minor >= $6",
		"total_minor <= $7",
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

func TestBuildListOrdersQueryDefaultsLimit(t *testing.T) {
	query, args, err := buildListOrdersQuery(OrderFilter{TenantID: tenant.ID("tenant-1")})
	if err != nil {
		t.Fatalf("expected query: %v", err)
	}
	if !strings.Contains(query, "LIMIT $2") {
		t.Fatalf("expected limit placeholder in query: %s", query)
	}
	if len(args) != 2 || args[1] != defaultOrderListLimit {
		t.Fatalf("unexpected args: %#v", args)
	}
}

func TestBuildListOrdersQueryRejectsBadStatus(t *testing.T) {
	_, _, err := buildListOrdersQuery(OrderFilter{
		TenantID:    tenant.ID("tenant-1"),
		OrderStatus: OrderStatus("bad"),
	})
	if !errors.Is(err, ErrOrderStatusInvalid) {
		t.Fatalf("expected order status error, got %v", err)
	}
}

func TestBuildGetOrderQueryAddsBuyerScope(t *testing.T) {
	query, args, err := buildGetOrderQuery(OrderLookup{
		ID:          OrderID("order-1"),
		TenantID:    tenant.ID("tenant-1"),
		BuyerUserID: identity.UserID("buyer-1"),
	})
	if err != nil {
		t.Fatalf("expected query: %v", err)
	}
	if !strings.Contains(query, "order_id = $1") || !strings.Contains(query, "tenant_id = $2") ||
		!strings.Contains(query, "buyer_user_id = $3") {
		t.Fatalf("expected scoped lookup query, got %s", query)
	}
	if len(args) != 3 {
		t.Fatalf("unexpected args: %#v", args)
	}
}

func int64Ptr(value int64) *int64 {
	return &value
}
