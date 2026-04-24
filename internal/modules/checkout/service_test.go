package checkout

import (
	"context"
	"errors"
	"testing"

	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/invoice"
	"github.com/Chinsusu/Billing-V2/internal/modules/order"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

func TestCheckoutOrderUsesInvoiceIssuerWithNormalizedInput(t *testing.T) {
	issuer := &fakeInvoiceIssuer{detail: testInvoiceDetail()}
	service := NewService(issuer)

	detail, err := service.CheckoutOrder(context.Background(), CheckoutOrderInput{
		TenantID:       tenant.ID("tenant_1"),
		BuyerUserID:    identity.UserID("buyer_1"),
		OrderID:        order.OrderID("order_1"),
		IdempotencyKey: invoice.IdempotencyKey(" checkout-key "),
	})
	if err != nil {
		t.Fatalf("expected checkout invoice: %v", err)
	}
	if issuer.calls != 1 {
		t.Fatalf("expected issuer once, got %d", issuer.calls)
	}
	if issuer.input.TenantID != tenant.ID("tenant_1") ||
		issuer.input.BuyerUserID != identity.UserID("buyer_1") ||
		issuer.input.OrderID != order.OrderID("order_1") ||
		issuer.input.IdempotencyKey != invoice.IdempotencyKey("checkout-key") {
		t.Fatalf("unexpected issuer input: %+v", issuer.input)
	}
	if detail.Invoice.DisplayID != 70001 {
		t.Fatalf("unexpected invoice detail: %+v", detail)
	}
}

func TestCheckoutOrderRequiresBuyerAndIdempotency(t *testing.T) {
	service := NewService(&fakeInvoiceIssuer{})

	_, err := service.CheckoutOrder(context.Background(), CheckoutOrderInput{
		TenantID: tenant.ID("tenant_1"),
		OrderID:  order.OrderID("order_1"),
	})
	if !errors.Is(err, invoice.ErrBuyerIDMissing) {
		t.Fatalf("expected buyer error, got %v", err)
	}

	_, err = service.CheckoutOrder(context.Background(), CheckoutOrderInput{
		TenantID:    tenant.ID("tenant_1"),
		BuyerUserID: identity.UserID("buyer_1"),
		OrderID:     order.OrderID("order_1"),
	})
	if !errors.Is(err, invoice.ErrIdempotencyKeyMissing) {
		t.Fatalf("expected idempotency error, got %v", err)
	}
}

type fakeInvoiceIssuer struct {
	calls  int
	input  invoice.IssueInvoiceForOrderInput
	detail invoice.InvoiceDetail
	err    error
}

func (issuer *fakeInvoiceIssuer) IssueInvoiceForOrder(ctx context.Context, input invoice.IssueInvoiceForOrderInput) (invoice.InvoiceDetail, error) {
	issuer.calls++
	issuer.input = input
	if issuer.err != nil {
		return invoice.InvoiceDetail{}, issuer.err
	}
	return issuer.detail, nil
}
