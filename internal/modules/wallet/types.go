package wallet

import (
	"encoding/json"
	"errors"
	"strings"
)

var (
	ErrWalletIDMissing       = errors.New("wallet id missing")
	ErrLedgerEntryIDMissing  = errors.New("wallet ledger entry id missing")
	ErrTopupRequestIDMissing = errors.New("wallet top-up request id missing")
	ErrOwnerTypeInvalid      = errors.New("wallet owner type invalid")
	ErrOwnerIDMissing        = errors.New("wallet owner id missing")
	ErrStatusInvalid         = errors.New("wallet status invalid")
	ErrDirectionInvalid      = errors.New("wallet ledger direction invalid")
	ErrEntryTypeInvalid      = errors.New("wallet ledger entry type invalid")
	ErrLedgerStatusInvalid   = errors.New("wallet ledger status invalid")
	ErrTopupStatusInvalid    = errors.New("wallet top-up status invalid")
	ErrPaymentMethodInvalid  = errors.New("wallet top-up payment method invalid")
	ErrCurrencyMissing       = errors.New("wallet currency missing")
	ErrCurrencyInvalid       = errors.New("wallet currency invalid")
	ErrBalanceInvalid        = errors.New("wallet balance invalid")
	ErrInsufficientBalance   = errors.New("wallet insufficient balance")
	ErrAmountInvalid         = errors.New("wallet amount invalid")
	ErrReferenceTypeMissing  = errors.New("wallet reference type missing")
	ErrReferenceIDMissing    = errors.New("wallet reference id missing")
	ErrIdempotencyKeyMissing = errors.New("wallet idempotency key missing")
	ErrCorrelationIDMissing  = errors.New("wallet correlation id missing")
	ErrReasonMissing         = errors.New("wallet ledger reason missing")
	ErrStoreExecutorMissing  = errors.New("wallet store executor missing")
	ErrServiceStoreMissing   = errors.New("wallet service store missing")
	ErrWalletNotFound        = errors.New("wallet not found")
	ErrLedgerEntryNotFound   = errors.New("wallet ledger entry not found")
	ErrTopupRequestNotFound  = errors.New("wallet top-up request not found")
)

type WalletID string
type LedgerEntryID string
type TopupRequestID string
type OwnerID string
type ReferenceType string
type ReferenceID string
type IdempotencyKey string
type CorrelationID string

func (id WalletID) Empty() bool       { return strings.TrimSpace(string(id)) == "" }
func (id LedgerEntryID) Empty() bool  { return strings.TrimSpace(string(id)) == "" }
func (id TopupRequestID) Empty() bool { return strings.TrimSpace(string(id)) == "" }
func (id OwnerID) Empty() bool        { return strings.TrimSpace(string(id)) == "" }

func trim(value string) string {
	return strings.TrimSpace(value)
}

func upperTrim(value string) string {
	return strings.ToUpper(strings.TrimSpace(value))
}

func defaultJSON(value json.RawMessage) json.RawMessage {
	if len(value) == 0 {
		return json.RawMessage(`{}`)
	}
	return append(json.RawMessage(nil), value...)
}

func validateCurrency(value string) error {
	if value == "" {
		return ErrCurrencyMissing
	}
	if len(value) != 3 {
		return ErrCurrencyInvalid
	}
	for _, letter := range value {
		if letter < 'A' || letter > 'Z' {
			return ErrCurrencyInvalid
		}
	}
	return nil
}

func validateBalance(value int64) error {
	if value < 0 {
		return ErrBalanceInvalid
	}
	return nil
}

func validatePositiveAmount(value int64) error {
	if value <= 0 {
		return ErrAmountInvalid
	}
	return nil
}
