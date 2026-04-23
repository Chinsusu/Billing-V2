package invoice

import (
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

type MarkInvoicePaidInput struct {
	ID                   InvoiceID
	TenantID             tenant.ID
	PaidAt               time.Time
	PaymentTransactionID string
	WalletID             string
	LedgerEntryID        string
	IdempotencyKey       IdempotencyKey
}

func (input MarkInvoicePaidInput) Normalize() MarkInvoicePaidInput {
	output := input
	output.PaymentTransactionID = trim(output.PaymentTransactionID)
	output.WalletID = trim(output.WalletID)
	output.LedgerEntryID = trim(output.LedgerEntryID)
	output.IdempotencyKey = IdempotencyKey(trim(string(output.IdempotencyKey)))
	if !output.PaidAt.IsZero() {
		output.PaidAt = output.PaidAt.UTC()
	}
	return output
}

func (input MarkInvoicePaidInput) Validate() error {
	if input.ID.Empty() {
		return ErrInvoiceIDMissing
	}
	if input.TenantID.Empty() {
		return tenant.ErrTenantIDMissing
	}
	if input.IdempotencyKey == "" {
		return ErrIdempotencyKeyMissing
	}
	return nil
}
