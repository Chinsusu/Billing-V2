package invoice

import (
	"encoding/json"
	"errors"
	"strings"
)

var (
	ErrInvoiceIDMissing      = errors.New("invoice id missing")
	ErrInvoiceItemIDMissing  = errors.New("invoice item id missing")
	ErrBuyerIDMissing        = errors.New("invoice buyer id missing")
	ErrCurrencyMissing       = errors.New("invoice currency missing")
	ErrCurrencyInvalid       = errors.New("invoice currency invalid")
	ErrAmountInvalid         = errors.New("invoice amount invalid")
	ErrTotalInvalid          = errors.New("invoice total invalid")
	ErrQuantityInvalid       = errors.New("invoice item quantity invalid")
	ErrDescriptionMissing    = errors.New("invoice item description missing")
	ErrStatusInvalid         = errors.New("invoice status invalid")
	ErrIdempotencyKeyMissing = errors.New("invoice idempotency key missing")
	ErrOrderReaderMissing    = errors.New("invoice order reader missing")
	ErrOrderNotPaid          = errors.New("invoice generation requires paid order")
	ErrOrderNotCheckoutable  = errors.New("invoice checkout requires pending unpaid order")
	ErrInvoiceStatusConflict = errors.New("invoice status conflict")
	ErrStoreExecutorMissing  = errors.New("invoice store executor missing")
	ErrServiceStoreMissing   = errors.New("invoice service store missing")
	ErrInvoiceNotFound       = errors.New("invoice not found")
)

type InvoiceID string
type InvoiceItemID string
type OrderItemID string
type IdempotencyKey string

func (id InvoiceID) Empty() bool     { return strings.TrimSpace(string(id)) == "" }
func (id InvoiceItemID) Empty() bool { return strings.TrimSpace(string(id)) == "" }
func (id OrderItemID) Empty() bool   { return strings.TrimSpace(string(id)) == "" }

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

func validateMinorAmount(value int64) error {
	if value < 0 {
		return ErrAmountInvalid
	}
	return nil
}
