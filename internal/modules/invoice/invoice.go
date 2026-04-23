package invoice

import (
	"encoding/json"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/order"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

type Invoice struct {
	ID            InvoiceID
	DisplayID     int64
	TenantID      tenant.ID
	BuyerUserID   identity.UserID
	OrderID       order.OrderID
	Status        Status
	Currency      string
	SubtotalMinor int64
	TaxMinor      int64
	DiscountMinor int64
	TotalMinor    int64
	IssuedAt      time.Time
	DueAt         time.Time
	PaidAt        time.Time
	VoidedAt      time.Time
	Metadata      json.RawMessage
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type CreateInvoiceInput struct {
	TenantID      tenant.ID
	BuyerUserID   identity.UserID
	OrderID       order.OrderID
	Status        Status
	Currency      string
	SubtotalMinor int64
	TaxMinor      int64
	DiscountMinor int64
	TotalMinor    int64
	IssuedAt      time.Time
	DueAt         time.Time
	Metadata      json.RawMessage
}

func (input CreateInvoiceInput) Normalize() CreateInvoiceInput {
	output := input
	output.Currency = upperTrim(output.Currency)
	output.Metadata = defaultJSON(output.Metadata)
	if output.Status == "" {
		output.Status = StatusDraft
	}
	return output
}

func (input CreateInvoiceInput) Validate() error {
	if input.TenantID.Empty() {
		return tenant.ErrTenantIDMissing
	}
	if input.BuyerUserID == "" {
		return ErrBuyerIDMissing
	}
	if !input.Status.Valid() {
		return ErrStatusInvalid
	}
	if err := validateCurrency(input.Currency); err != nil {
		return err
	}
	if err := validateMinorAmount(input.SubtotalMinor); err != nil {
		return err
	}
	if err := validateMinorAmount(input.TaxMinor); err != nil {
		return err
	}
	if err := validateMinorAmount(input.DiscountMinor); err != nil {
		return err
	}
	if err := validateMinorAmount(input.TotalMinor); err != nil {
		return err
	}
	if input.SubtotalMinor+input.TaxMinor-input.DiscountMinor != input.TotalMinor {
		return ErrTotalInvalid
	}
	return nil
}
