package checkout

import (
	"context"

	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/invoice"
	"github.com/Chinsusu/Billing-V2/internal/modules/order"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

type InvoiceIssuer interface {
	IssueInvoiceForOrder(ctx context.Context, input invoice.IssueInvoiceForOrderInput) (invoice.InvoiceDetail, error)
}

type Service struct {
	invoiceIssuer InvoiceIssuer
}

type CheckoutOrderInput struct {
	TenantID       tenant.ID
	BuyerUserID    identity.UserID
	OrderID        order.OrderID
	IdempotencyKey invoice.IdempotencyKey
}

func NewService(invoiceIssuer InvoiceIssuer) *Service {
	return &Service{invoiceIssuer: invoiceIssuer}
}

func (service *Service) CheckoutOrder(ctx context.Context, input CheckoutOrderInput) (invoice.InvoiceDetail, error) {
	if service == nil || service.invoiceIssuer == nil {
		return invoice.InvoiceDetail{}, invoice.ErrServiceStoreMissing
	}
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return invoice.InvoiceDetail{}, err
	}
	return service.invoiceIssuer.IssueInvoiceForOrder(ctx, invoice.IssueInvoiceForOrderInput{
		TenantID:       input.TenantID,
		BuyerUserID:    input.BuyerUserID,
		OrderID:        input.OrderID,
		IdempotencyKey: input.IdempotencyKey,
	})
}

func (input CheckoutOrderInput) Normalize() CheckoutOrderInput {
	output := input
	output.IdempotencyKey = invoice.IdempotencyKey(trim(string(output.IdempotencyKey)))
	return output
}

func (input CheckoutOrderInput) Validate() error {
	if input.TenantID.Empty() {
		return tenant.ErrTenantIDMissing
	}
	if input.BuyerUserID == "" {
		return invoice.ErrBuyerIDMissing
	}
	if input.OrderID.Empty() {
		return order.ErrOrderIDMissing
	}
	if input.IdempotencyKey == "" {
		return invoice.ErrIdempotencyKeyMissing
	}
	return nil
}
