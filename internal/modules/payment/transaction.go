package payment

import (
	"encoding/json"
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/invoice"
	"github.com/Chinsusu/Billing-V2/internal/modules/order"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

type Transaction struct {
	ID             TransactionID
	DisplayID      int64
	TenantID       tenant.ID
	AccountUserID  identity.UserID
	OrderID        order.OrderID
	InvoiceID      invoice.InvoiceID
	Type           TransactionType
	Status         TransactionStatus
	Currency       string
	AmountMinor    int64
	Description    string
	IdempotencyKey IdempotencyKey
	Metadata       json.RawMessage
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type CreateTransactionInput struct {
	TenantID       tenant.ID
	AccountUserID  identity.UserID
	OrderID        order.OrderID
	InvoiceID      invoice.InvoiceID
	Type           TransactionType
	Status         TransactionStatus
	Currency       string
	AmountMinor    int64
	Description    string
	IdempotencyKey IdempotencyKey
	Metadata       json.RawMessage
}

func (input CreateTransactionInput) Normalize() CreateTransactionInput {
	output := input
	output.Currency = upperTrim(output.Currency)
	output.Description = trim(output.Description)
	output.IdempotencyKey = IdempotencyKey(trim(string(output.IdempotencyKey)))
	output.Metadata = defaultJSON(output.Metadata)
	if output.Status == "" {
		output.Status = TransactionStatusPosted
	}
	return output
}

func (input CreateTransactionInput) Validate() error {
	if input.TenantID.Empty() {
		return tenant.ErrTenantIDMissing
	}
	if input.AccountUserID == "" {
		return ErrAccountIDMissing
	}
	if !input.Type.Valid() {
		return ErrTypeInvalid
	}
	if !input.Status.Valid() {
		return ErrStatusInvalid
	}
	if err := validateCurrency(input.Currency); err != nil {
		return err
	}
	if err := validatePositiveMinorAmount(input.AmountMinor); err != nil {
		return err
	}
	if input.IdempotencyKey == "" {
		return ErrIdempotencyKeyMissing
	}
	return nil
}
