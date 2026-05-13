package wallet

import (
	"context"

	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

type CreateWalletRefundInput struct {
	TenantID       tenant.ID
	WalletID       WalletID
	AmountMinor    int64
	Currency       string
	ReferenceType  ReferenceType
	ReferenceID    ReferenceID
	IdempotencyKey IdempotencyKey
	CreatedBy      identity.UserID
	Reason         string
	CorrelationID  CorrelationID
}

type CreateWalletAdjustmentInput struct {
	TenantID       tenant.ID
	WalletID       WalletID
	Direction      Direction
	AmountMinor    int64
	Currency       string
	ReferenceType  ReferenceType
	ReferenceID    ReferenceID
	IdempotencyKey IdempotencyKey
	CreatedBy      identity.UserID
	Reason         string
	CorrelationID  CorrelationID
}

func (input CreateWalletRefundInput) Normalize() CreateWalletRefundInput {
	output := input
	output.Currency = upperTrim(output.Currency)
	output.ReferenceType = ReferenceType(trim(string(output.ReferenceType)))
	output.ReferenceID = ReferenceID(trim(string(output.ReferenceID)))
	output.IdempotencyKey = IdempotencyKey(trim(string(output.IdempotencyKey)))
	output.Reason = trim(output.Reason)
	output.CorrelationID = CorrelationID(trim(string(output.CorrelationID)))
	return output
}

func (input CreateWalletRefundInput) Validate() error {
	if input.TenantID.Empty() {
		return tenant.ErrTenantIDMissing
	}
	if input.WalletID.Empty() {
		return ErrWalletIDMissing
	}
	if input.CreatedBy == "" {
		return identity.ErrActorIDMissing
	}
	if err := validatePositiveAmount(input.AmountMinor); err != nil {
		return err
	}
	if err := validateCurrency(input.Currency); err != nil {
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
	if input.Reason == "" {
		return ErrReasonMissing
	}
	if input.CorrelationID == "" {
		return ErrCorrelationIDMissing
	}
	return nil
}

func (input CreateWalletAdjustmentInput) Normalize() CreateWalletAdjustmentInput {
	output := input
	output.Currency = upperTrim(output.Currency)
	output.ReferenceType = ReferenceType(trim(string(output.ReferenceType)))
	output.ReferenceID = ReferenceID(trim(string(output.ReferenceID)))
	output.IdempotencyKey = IdempotencyKey(trim(string(output.IdempotencyKey)))
	output.Reason = trim(output.Reason)
	output.CorrelationID = CorrelationID(trim(string(output.CorrelationID)))
	return output
}

func (input CreateWalletAdjustmentInput) Validate() error {
	if input.TenantID.Empty() {
		return tenant.ErrTenantIDMissing
	}
	if input.WalletID.Empty() {
		return ErrWalletIDMissing
	}
	if input.CreatedBy == "" {
		return identity.ErrActorIDMissing
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
	if input.ReferenceType == "" {
		return ErrReferenceTypeMissing
	}
	if input.ReferenceID == "" {
		return ErrReferenceIDMissing
	}
	if input.IdempotencyKey == "" {
		return ErrIdempotencyKeyMissing
	}
	if input.Reason == "" {
		return ErrReasonMissing
	}
	if input.CorrelationID == "" {
		return ErrCorrelationIDMissing
	}
	return nil
}

func (service *Service) CreateWalletRefund(ctx context.Context, input CreateWalletRefundInput) (LedgerEntry, error) {
	if err := service.ready(); err != nil {
		return LedgerEntry{}, err
	}
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return LedgerEntry{}, err
	}
	if err := service.validateManualLedgerWallet(ctx, input.WalletID, input.TenantID, input.Currency); err != nil {
		return LedgerEntry{}, err
	}
	postInput := refundPostInput(input)
	result, err := service.store.PostLedgerEntryResult(ctx, postInput)
	if err != nil {
		return LedgerEntry{}, err
	}
	if err := ensureManualLedgerEntryMatches(result.Entry, postInput); err != nil {
		return LedgerEntry{}, err
	}
	if result.Created {
		if err := service.appendManualLedgerAudit(ctx, walletAuditActionRefundCreated, result.Entry); err != nil {
			return LedgerEntry{}, err
		}
	}
	return result.Entry, nil
}

func (service *Service) CreateWalletAdjustment(ctx context.Context, input CreateWalletAdjustmentInput) (LedgerEntry, error) {
	if err := service.ready(); err != nil {
		return LedgerEntry{}, err
	}
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return LedgerEntry{}, err
	}
	if err := service.validateManualLedgerWallet(ctx, input.WalletID, input.TenantID, input.Currency); err != nil {
		return LedgerEntry{}, err
	}
	postInput := adjustmentPostInput(input)
	result, err := service.store.PostLedgerEntryResult(ctx, postInput)
	if err != nil {
		return LedgerEntry{}, err
	}
	if err := ensureManualLedgerEntryMatches(result.Entry, postInput); err != nil {
		return LedgerEntry{}, err
	}
	if result.Created {
		if err := service.appendManualLedgerAudit(ctx, walletAuditActionAdjustmentCreated, result.Entry); err != nil {
			return LedgerEntry{}, err
		}
	}
	return result.Entry, nil
}

func (service *Service) validateManualLedgerWallet(ctx context.Context, walletID WalletID, tenantID tenant.ID, currency string) error {
	record, err := service.store.GetWallet(ctx, WalletLookup{ID: walletID, TenantID: tenantID})
	if err != nil {
		return err
	}
	if record.Status != StatusActive {
		return ErrWalletStatusConflict
	}
	if record.Currency != currency {
		return ErrWalletCurrencyMismatch
	}
	return nil
}

func refundPostInput(input CreateWalletRefundInput) PostLedgerEntryInput {
	return PostLedgerEntryInput{
		WalletID:       input.WalletID,
		TenantID:       input.TenantID,
		Direction:      DirectionCredit,
		AmountMinor:    input.AmountMinor,
		Currency:       input.Currency,
		EntryType:      EntryTypeRefund,
		ReferenceType:  input.ReferenceType,
		ReferenceID:    input.ReferenceID,
		IdempotencyKey: input.IdempotencyKey,
		CreatedBy:      input.CreatedBy,
		Reason:         input.Reason,
		CorrelationID:  input.CorrelationID,
	}
}

func adjustmentPostInput(input CreateWalletAdjustmentInput) PostLedgerEntryInput {
	return PostLedgerEntryInput{
		WalletID:       input.WalletID,
		TenantID:       input.TenantID,
		Direction:      input.Direction,
		AmountMinor:    input.AmountMinor,
		Currency:       input.Currency,
		EntryType:      EntryTypeAdjustment,
		ReferenceType:  input.ReferenceType,
		ReferenceID:    input.ReferenceID,
		IdempotencyKey: input.IdempotencyKey,
		CreatedBy:      input.CreatedBy,
		Reason:         input.Reason,
		CorrelationID:  input.CorrelationID,
	}
}

func ensureManualLedgerEntryMatches(entry LedgerEntry, input PostLedgerEntryInput) error {
	if entry.TenantID != input.TenantID ||
		entry.WalletID != input.WalletID ||
		entry.Direction != input.Direction ||
		entry.AmountMinor != input.AmountMinor ||
		entry.Currency != input.Currency ||
		entry.EntryType != input.EntryType ||
		entry.Status != LedgerStatusPosted ||
		entry.ReferenceType != input.ReferenceType ||
		entry.ReferenceID != input.ReferenceID ||
		entry.CreatedBy != input.CreatedBy ||
		entry.Reason != input.Reason ||
		entry.CorrelationID != input.CorrelationID {
		return ErrIdempotencyConflict
	}
	return nil
}
