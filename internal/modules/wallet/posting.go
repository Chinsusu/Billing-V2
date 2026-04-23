package wallet

import (
	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

type PostLedgerEntryInput struct {
	WalletID       WalletID
	TenantID       tenant.ID
	Direction      Direction
	AmountMinor    int64
	Currency       string
	EntryType      EntryType
	ReferenceType  ReferenceType
	ReferenceID    ReferenceID
	IdempotencyKey IdempotencyKey
	CreatedBy      identity.UserID
	Reason         string
	CorrelationID  CorrelationID
}

func (input PostLedgerEntryInput) Normalize() PostLedgerEntryInput {
	output := input
	output.Currency = upperTrim(output.Currency)
	output.ReferenceType = ReferenceType(trim(string(output.ReferenceType)))
	output.ReferenceID = ReferenceID(trim(string(output.ReferenceID)))
	output.IdempotencyKey = IdempotencyKey(trim(string(output.IdempotencyKey)))
	output.Reason = trim(output.Reason)
	output.CorrelationID = CorrelationID(trim(string(output.CorrelationID)))
	return output
}

func (input PostLedgerEntryInput) Validate() error {
	if input.WalletID.Empty() {
		return ErrWalletIDMissing
	}
	if input.TenantID.Empty() {
		return tenant.ErrTenantIDMissing
	}
	if !input.Direction.Valid() {
		return ErrDirectionInvalid
	}
	if err := validatePositiveAmount(input.AmountMinor); err != nil {
		return err
	}
	if err := validateCurrency(input.Currency); err != nil {
		return err
	}
	if !input.EntryType.Valid() {
		return ErrEntryTypeInvalid
	}
	if input.ReferenceType == "" {
		return ErrReferenceTypeMissing
	}
	if input.ReferenceID == "" {
		return ErrReferenceIDMissing
	}
	if input.IdempotencyKey == "" {
		return ErrIdempotencyKeyMissing
	}
	if input.EntryType == EntryTypeAdjustment && input.Reason == "" {
		return ErrReasonMissing
	}
	if input.CorrelationID == "" {
		return ErrCorrelationIDMissing
	}
	return nil
}
