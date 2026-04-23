package payment

import (
	"encoding/json"
	"errors"
	"strings"
)

var (
	ErrTransactionIDMissing     = errors.New("payment transaction id missing")
	ErrAccountIDMissing         = errors.New("payment account id missing")
	ErrCurrencyMissing          = errors.New("payment currency missing")
	ErrCurrencyInvalid          = errors.New("payment currency invalid")
	ErrAmountInvalid            = errors.New("payment amount invalid")
	ErrTypeInvalid              = errors.New("payment transaction type invalid")
	ErrStatusInvalid            = errors.New("payment transaction status invalid")
	ErrIdempotencyKeyMissing    = errors.New("payment idempotency key missing")
	ErrIdempotencyConflict      = errors.New("payment idempotency conflict")
	ErrInvoiceNotPayable        = errors.New("payment invoice not payable")
	ErrWalletCurrencyMismatch   = errors.New("payment wallet currency mismatch")
	ErrBillingDependencyMissing = errors.New("payment billing dependency missing")
	ErrCreatedTimeInvalid       = errors.New("payment created time invalid")
	ErrCreatedTimeWindowInvalid = errors.New("payment created time window invalid")
	ErrStoreExecutorMissing     = errors.New("payment store executor missing")
	ErrServiceStoreMissing      = errors.New("payment service store missing")
	ErrTransactionNotFound      = errors.New("payment transaction not found")
)

type TransactionID string
type IdempotencyKey string

func (id TransactionID) Empty() bool { return strings.TrimSpace(string(id)) == "" }

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

func validatePositiveMinorAmount(value int64) error {
	if value <= 0 {
		return ErrAmountInvalid
	}
	return nil
}
