package wallet

import (
	"time"

	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

type LedgerEntry struct {
	ID                 LedgerEntryID
	DisplayID          int64
	WalletID           WalletID
	TenantID           tenant.ID
	Direction          Direction
	AmountMinor        int64
	Currency           string
	EntryType          EntryType
	Status             LedgerStatus
	BalanceAfterMinor  int64
	ReferenceType      ReferenceType
	ReferenceID        ReferenceID
	ReferenceDisplayID int64
	IdempotencyKey     IdempotencyKey
	CreatedBy          identity.UserID
	Reason             string
	CorrelationID      CorrelationID
	CreatedAt          time.Time
}

type CreateLedgerEntryInput struct {
	WalletID          WalletID
	TenantID          tenant.ID
	Direction         Direction
	AmountMinor       int64
	Currency          string
	EntryType         EntryType
	Status            LedgerStatus
	BalanceAfterMinor int64
	ReferenceType     ReferenceType
	ReferenceID       ReferenceID
	IdempotencyKey    IdempotencyKey
	CreatedBy         identity.UserID
	Reason            string
	CorrelationID     CorrelationID
}

func (input CreateLedgerEntryInput) Normalize() CreateLedgerEntryInput {
	output := input
	output.Currency = upperTrim(output.Currency)
	output.ReferenceType = ReferenceType(trim(string(output.ReferenceType)))
	output.ReferenceID = ReferenceID(trim(string(output.ReferenceID)))
	output.IdempotencyKey = IdempotencyKey(trim(string(output.IdempotencyKey)))
	output.Reason = trim(output.Reason)
	output.CorrelationID = CorrelationID(trim(string(output.CorrelationID)))
	if output.Status == "" {
		output.Status = LedgerStatusPosted
	}
	return output
}

func (input CreateLedgerEntryInput) Validate() error {
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
	if !input.Status.Valid() {
		return ErrLedgerStatusInvalid
	}
	if err := validateBalance(input.BalanceAfterMinor); err != nil {
		return err
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
	if input.EntryType == EntryTypeAdjustment && input.CreatedBy == "" {
		return identity.ErrActorIDMissing
	}
	if input.CorrelationID == "" {
		return ErrCorrelationIDMissing
	}
	return nil
}
