package wallet

import (
	"encoding/json"
	"errors"
	"strings"
)

var (
	ErrWalletIDMissing      = errors.New("wallet id missing")
	ErrOwnerTypeInvalid     = errors.New("wallet owner type invalid")
	ErrOwnerIDMissing       = errors.New("wallet owner id missing")
	ErrStatusInvalid        = errors.New("wallet status invalid")
	ErrCurrencyMissing      = errors.New("wallet currency missing")
	ErrCurrencyInvalid      = errors.New("wallet currency invalid")
	ErrBalanceInvalid       = errors.New("wallet balance invalid")
	ErrStoreExecutorMissing = errors.New("wallet store executor missing")
	ErrServiceStoreMissing  = errors.New("wallet service store missing")
	ErrWalletNotFound       = errors.New("wallet not found")
)

type WalletID string
type OwnerID string

func (id WalletID) Empty() bool { return strings.TrimSpace(string(id)) == "" }
func (id OwnerID) Empty() bool  { return strings.TrimSpace(string(id)) == "" }

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
