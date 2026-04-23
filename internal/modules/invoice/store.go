package invoice

import (
	"context"

	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/order"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

type InvoiceFilter struct {
	TenantID    tenant.ID
	BuyerUserID identity.UserID
	OrderID     order.OrderID
	Status      Status
	Limit       int
}

type InvoiceLookup struct {
	ID          InvoiceID
	TenantID    tenant.ID
	BuyerUserID identity.UserID
}

type InvoiceDetail struct {
	Invoice Invoice
	Items   []Item
}

type Store interface {
	ListInvoices(ctx context.Context, filter InvoiceFilter) ([]Invoice, error)
	GetInvoice(ctx context.Context, lookup InvoiceLookup) (InvoiceDetail, error)
	CreateInvoiceFromOrder(ctx context.Context, input CreateInvoiceFromOrderInput) (InvoiceDetail, error)
	MarkInvoicePaid(ctx context.Context, input MarkInvoicePaidInput) (InvoiceDetail, error)
}
