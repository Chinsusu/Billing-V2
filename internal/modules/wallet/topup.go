package wallet

import (
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

type TopupRequest struct {
	ID                   TopupRequestID
	DisplayID            int64
	TenantID             tenant.ID
	WalletID             WalletID
	WalletDisplayID      int64
	RequestedBy          identity.UserID
	RequestedByDisplayID int64
	AmountMinor          int64
	Currency             string
	PaymentMethod        PaymentMethod
	PaymentReference     string
	Status               TopupStatus
	ReviewedBy           identity.UserID
	ReviewedByDisplayID  int64
	ReviewedAt           *time.Time
	ReviewNote           string
	LedgerEntryID        LedgerEntryID
	IdempotencyKey       IdempotencyKey
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

type CreateTopupRequestInput struct {
	TenantID         tenant.ID
	WalletID         WalletID
	RequestedBy      identity.UserID
	AmountMinor      int64
	Currency         string
	PaymentMethod    PaymentMethod
	PaymentReference string
	Status           TopupStatus
	IdempotencyKey   IdempotencyKey
}

func (input CreateTopupRequestInput) Normalize() CreateTopupRequestInput {
	output := input
	output.Currency = upperTrim(output.Currency)
	output.PaymentReference = trim(output.PaymentReference)
	output.IdempotencyKey = IdempotencyKey(trim(string(output.IdempotencyKey)))
	if output.Status == "" {
		output.Status = TopupStatusSubmitted
	}
	return output
}

func (input CreateTopupRequestInput) Validate() error {
	if input.TenantID.Empty() {
		return tenant.ErrTenantIDMissing
	}
	if input.WalletID.Empty() {
		return ErrWalletIDMissing
	}
	if input.RequestedBy == "" {
		return identity.ErrActorIDMissing
	}
	if err := validatePositiveAmount(input.AmountMinor); err != nil {
		return err
	}
	if err := validateCurrency(input.Currency); err != nil {
		return err
	}
	if !input.PaymentMethod.Valid() {
		return ErrPaymentMethodInvalid
	}
	if !input.Status.Valid() {
		return ErrTopupStatusInvalid
	}
	if input.IdempotencyKey == "" {
		return ErrIdempotencyKeyMissing
	}
	return nil
}
