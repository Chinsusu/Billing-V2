package invoice

import (
	"database/sql"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/order"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

func TestCreateInvoiceFromOrderArgsNormalizeValidate(t *testing.T) {
	args, err := createInvoiceFromOrderArgs(CreateInvoiceFromOrderInput{
		Invoice: CreateInvoiceInput{
			TenantID:      tenant.ID("tenant-1"),
			BuyerUserID:   identity.UserID("buyer-1"),
			OrderID:       order.OrderID("order-1"),
			Status:        StatusIssued,
			Currency:      " usd ",
			SubtotalMinor: 2000,
			DiscountMinor: 200,
			TotalMinor:    1800,
			Metadata:      json.RawMessage(`{"source":"order"}`),
		},
		Item: GeneratedInvoiceItemInput{
			OrderID:        order.OrderID("order-1"),
			Description:    " Service ",
			Quantity:       2,
			UnitPriceMinor: 1000,
			DiscountMinor:  200,
			LineTotalMinor: 1800,
			Metadata:       json.RawMessage(`{"source":"order"}`),
		},
		IdempotencyKey: IdempotencyKey(" key-1 "),
		OrderDisplayID: 60001,
	})
	if err != nil {
		t.Fatalf("expected args: %v", err)
	}
	if args[4] != "USD" || args[14] != "Service" || args[22] != IdempotencyKey("key-1") {
		t.Fatalf("unexpected normalized args: %#v", args)
	}
	if _, ok := args[12].(sql.NullString); !ok {
		t.Fatalf("expected nullable order item id, got %#v", args[12])
	}
}

func TestCreateInvoiceFromOrderArgsRejectsMissingIdempotency(t *testing.T) {
	_, err := createInvoiceFromOrderArgs(CreateInvoiceFromOrderInput{
		Invoice: CreateInvoiceInput{
			TenantID:      tenant.ID("tenant-1"),
			BuyerUserID:   identity.UserID("buyer-1"),
			OrderID:       order.OrderID("order-1"),
			Status:        StatusIssued,
			Currency:      "USD",
			SubtotalMinor: 1000,
			TotalMinor:    1000,
		},
		Item: GeneratedInvoiceItemInput{
			OrderID:        order.OrderID("order-1"),
			Description:    "Service",
			Quantity:       1,
			UnitPriceMinor: 1000,
			LineTotalMinor: 1000,
		},
	})
	if !errors.Is(err, ErrIdempotencyKeyMissing) {
		t.Fatalf("expected idempotency error, got %v", err)
	}
}

func TestCreateInvoiceFromOrderSQLEmitsOutboxAndUsesOrderConflict(t *testing.T) {
	for _, clause := range []string{
		"ON CONFLICT (tenant_id, order_id) WHERE order_id IS NOT NULL",
		"INSERT INTO invoice_items",
		"INSERT INTO outbox_events",
		EventInvoiceGenerated,
	} {
		if !strings.Contains(createInvoiceFromOrderSQL, clause) {
			t.Fatalf("expected %q in create SQL", clause)
		}
	}
}
