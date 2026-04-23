package invoice

import (
	"encoding/json"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/order"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

type Item struct {
	ID             InvoiceItemID
	InvoiceID      InvoiceID
	TenantID       tenant.ID
	OrderID        order.OrderID
	OrderItemID    OrderItemID
	ServiceID      order.ServiceID
	Description    string
	Quantity       int
	UnitPriceMinor int64
	TaxMinor       int64
	DiscountMinor  int64
	LineTotalMinor int64
	Metadata       json.RawMessage
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type CreateItemInput struct {
	InvoiceID      InvoiceID
	TenantID       tenant.ID
	OrderID        order.OrderID
	OrderItemID    OrderItemID
	ServiceID      order.ServiceID
	Description    string
	Quantity       int
	UnitPriceMinor int64
	TaxMinor       int64
	DiscountMinor  int64
	LineTotalMinor int64
	Metadata       json.RawMessage
}

func (input CreateItemInput) Normalize() CreateItemInput {
	output := input
	output.Description = trim(output.Description)
	output.Metadata = defaultJSON(output.Metadata)
	if output.Quantity == 0 {
		output.Quantity = 1
	}
	return output
}

func (input CreateItemInput) Validate() error {
	if input.InvoiceID.Empty() {
		return ErrInvoiceIDMissing
	}
	if input.TenantID.Empty() {
		return tenant.ErrTenantIDMissing
	}
	if input.Description == "" {
		return ErrDescriptionMissing
	}
	if input.Quantity <= 0 {
		return ErrQuantityInvalid
	}
	if err := validateMinorAmount(input.UnitPriceMinor); err != nil {
		return err
	}
	if err := validateMinorAmount(input.TaxMinor); err != nil {
		return err
	}
	if err := validateMinorAmount(input.DiscountMinor); err != nil {
		return err
	}
	if err := validateMinorAmount(input.LineTotalMinor); err != nil {
		return err
	}
	expectedTotal := int64(input.Quantity)*input.UnitPriceMinor + input.TaxMinor - input.DiscountMinor
	if expectedTotal != input.LineTotalMinor {
		return ErrTotalInvalid
	}
	return nil
}
