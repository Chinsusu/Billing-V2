package order

import (
	"encoding/json"
	"errors"
	"strings"
)

var (
	ErrOrderIDMissing             = errors.New("order id missing")
	ErrReservationIDMissing       = errors.New("reservation id missing")
	ErrProvisioningJobIDMissing   = errors.New("provisioning job id missing")
	ErrServiceIDMissing           = errors.New("service id missing")
	ErrBuyerIDMissing             = errors.New("buyer id missing")
	ErrTenantPlanIDMissing        = errors.New("tenant plan id missing")
	ErrProviderSourceIDMissing    = errors.New("provider source id missing")
	ErrIdempotencyKeyMissing      = errors.New("idempotency key missing")
	ErrCurrencyMissing            = errors.New("currency missing")
	ErrCurrencyInvalid            = errors.New("currency invalid")
	ErrAmountInvalid              = errors.New("amount invalid")
	ErrQuantityInvalid            = errors.New("quantity invalid")
	ErrAttemptInvalid             = errors.New("attempt number invalid")
	ErrOrderStatusInvalid         = errors.New("order status invalid")
	ErrBillingStatusInvalid       = errors.New("billing status invalid")
	ErrReservationStatusInvalid   = errors.New("reservation status invalid")
	ErrProvisioningStatusInvalid  = errors.New("provisioning status invalid")
	ErrServiceStatusInvalid       = errors.New("service status invalid")
	ErrSuspensionReasonInvalid    = errors.New("suspension reason invalid")
	ErrReservationExpiryMissing   = errors.New("reservation expiry missing")
	ErrTermWindowInvalid          = errors.New("service term window invalid")
	ErrStatusTransitionInvalid    = errors.New("status transition invalid")
	ErrExternalResourceIDMissing  = errors.New("external resource id missing")
	ErrProviderOperationIDMissing = errors.New("provider operation id missing")
	ErrStoreExecutorMissing       = errors.New("order store executor missing")
	ErrServiceStoreMissing        = errors.New("order service store missing")
)

type OrderID string
type ReservationID string
type ProvisioningJobID string
type ServiceID string
type IdempotencyKey string
type ProviderOperationID string

func (id OrderID) Empty() bool           { return strings.TrimSpace(string(id)) == "" }
func (id ReservationID) Empty() bool     { return strings.TrimSpace(string(id)) == "" }
func (id ProvisioningJobID) Empty() bool { return strings.TrimSpace(string(id)) == "" }
func (id ServiceID) Empty() bool         { return strings.TrimSpace(string(id)) == "" }

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
