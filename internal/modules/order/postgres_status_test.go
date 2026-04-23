package order

import (
	"errors"
	"strings"
	"testing"

	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

func TestTransitionOrderStatusArgsNormalizeAndValidate(t *testing.T) {
	args, err := transitionOrderStatusArgs(TransitionOrderStatusInput{
		ID:            " order-1 ",
		TenantID:      tenant.ID(" tenant-1 "),
		ActorID:       " admin-1 ",
		FromStatus:    OrderStatusPendingPayment,
		ToStatus:      OrderStatusPaid,
		BillingStatus: BillingStatusPaid,
	})
	if err != nil {
		t.Fatalf("expected transition args: %v", err)
	}
	if len(args) != 5 {
		t.Fatalf("expected 5 args, got %d", len(args))
	}
	if args[0] != OrderID("order-1") || args[1] != tenant.ID("tenant-1") || args[4] != BillingStatusPaid {
		t.Fatalf("unexpected transition args: %#v", args)
	}
}

func TestTransitionOrderStatusArgsRejectsBadTransition(t *testing.T) {
	_, err := transitionOrderStatusArgs(TransitionOrderStatusInput{
		ID:            "order-1",
		TenantID:      tenant.ID("tenant-1"),
		ActorID:       "admin-1",
		FromStatus:    OrderStatusPendingPayment,
		ToStatus:      OrderStatusRefunded,
		BillingStatus: BillingStatusRefunded,
	})
	if !errors.Is(err, ErrStatusTransitionInvalid) {
		t.Fatalf("expected transition error, got %v", err)
	}
}

func TestTransitionOrderStatusSQLScopesTenantAndExpectedStatus(t *testing.T) {
	for _, clause := range []string{"UPDATE orders", "order_id = $1", "tenant_id = $2", "order_status = $3", "RETURNING"} {
		if !strings.Contains(transitionOrderStatusSQL, clause) {
			t.Fatalf("expected %q in transition SQL: %s", clause, transitionOrderStatusSQL)
		}
	}
}

func TestTransitionOrderStatusSQLEmitsStatusChangedOutboxEvent(t *testing.T) {
	for _, clause := range []string{
		"WITH updated AS",
		"INSERT INTO outbox_events",
		OrderEventStatusChanged,
		"'from_status', $3::text",
		"'to_status', order_status",
		"'display_id', display_id",
		"FROM updated",
	} {
		if !strings.Contains(transitionOrderStatusSQL, clause) {
			t.Fatalf("expected %q in transition SQL: %s", clause, transitionOrderStatusSQL)
		}
	}
}
