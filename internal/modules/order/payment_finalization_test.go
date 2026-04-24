package order

import (
	"errors"
	"strings"
	"testing"

	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

func TestFinalizePaymentInputNormalizeValidate(t *testing.T) {
	input := FinalizePaymentInput{
		ID:          " order-1 ",
		TenantID:    tenant.ID(" tenant-1 "),
		BuyerUserID: identity.UserID(" buyer-1 "),
	}.Normalize()

	if input.ID != OrderID("order-1") ||
		input.TenantID != tenant.ID("tenant-1") ||
		input.BuyerUserID != identity.UserID("buyer-1") {
		t.Fatalf("unexpected normalized input: %+v", input)
	}
	if err := input.Validate(); err != nil {
		t.Fatalf("expected valid finalization input: %v", err)
	}
}

func TestFinalizePaymentInputRejectsMissingBuyer(t *testing.T) {
	err := FinalizePaymentInput{
		ID:       "order-1",
		TenantID: tenant.ID("tenant-1"),
	}.Validate()
	if !errors.Is(err, ErrBuyerIDMissing) {
		t.Fatalf("expected buyer error, got %v", err)
	}
}

func TestFinalizePaymentArgsNormalizeAndValidate(t *testing.T) {
	args, err := finalizePaymentArgs(FinalizePaymentInput{
		ID:          " order-1 ",
		TenantID:    tenant.ID(" tenant-1 "),
		BuyerUserID: identity.UserID(" buyer-1 "),
	})
	if err != nil {
		t.Fatalf("expected finalize args: %v", err)
	}
	if len(args) != 3 ||
		args[0] != OrderID("order-1") ||
		args[1] != tenant.ID("tenant-1") ||
		args[2] != identity.UserID("buyer-1") {
		t.Fatalf("unexpected finalize args: %#v", args)
	}
}

func TestFinalizePaymentSQLScopesAndGuardsPaymentState(t *testing.T) {
	for _, clause := range []string{
		"UPDATE orders",
		"order_id = $1",
		"tenant_id = $2",
		"buyer_user_id = $3",
		"order_status = 'pending_payment'",
		"billing_status = 'unpaid'",
		"order_status = 'paid'",
		"billing_status = 'paid'",
		"already_paid AS",
	} {
		if !strings.Contains(finalizePaymentSQL, clause) {
			t.Fatalf("expected %q in finalization SQL: %s", clause, finalizePaymentSQL)
		}
	}
}

func TestFinalizePaymentSQLEmitsStatusChangedOutboxEvent(t *testing.T) {
	for _, clause := range []string{
		"WITH updated AS",
		"INSERT INTO outbox_events",
		OrderEventStatusChanged,
		"'from_status', 'pending_payment'",
		"'to_status', order_status",
		"'billing_status', billing_status",
		":payment_paid",
		"FROM updated",
	} {
		if !strings.Contains(finalizePaymentSQL, clause) {
			t.Fatalf("expected %q in finalization SQL: %s", clause, finalizePaymentSQL)
		}
	}
}

func TestGetOrderPaymentStateSQLScopesBuyer(t *testing.T) {
	for _, clause := range []string{
		"SELECT order_status, billing_status",
		"order_id = $1",
		"tenant_id = $2",
		"buyer_user_id = $3",
	} {
		if !strings.Contains(getOrderPaymentStateSQL, clause) {
			t.Fatalf("expected %q in payment state SQL: %s", clause, getOrderPaymentStateSQL)
		}
	}
}
