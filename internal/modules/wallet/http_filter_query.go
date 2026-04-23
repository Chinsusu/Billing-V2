package wallet

import (
	"net/http"

	"github.com/Chinsusu/Billing-V2/internal/platform/httpserver"
)

func walletPositiveInt64Query(w http.ResponseWriter, r *http.Request, field string) (int64, bool, bool) {
	value, present, err := httpserver.ParseOptionalPositiveInt64Query(r, field)
	if err != nil {
		httpserver.WriteValidationError(w, r, []httpserver.ValidationField{
			validationField(field, "request.display_id_invalid", "Display id must be a positive number."),
		})
		return 0, present, false
	}
	return value, present, true
}

func walletAmountRangeQuery(w http.ResponseWriter, r *http.Request) (*int64, *int64, bool) {
	minValue, minPresent, ok := walletNonNegativeInt64Query(w, r, "amount_min")
	if !ok {
		return nil, nil, false
	}
	maxValue, maxPresent, ok := walletNonNegativeInt64Query(w, r, "amount_max")
	if !ok {
		return nil, nil, false
	}
	if minPresent && maxPresent && maxValue < minValue {
		httpserver.WriteValidationError(w, r, []httpserver.ValidationField{
			validationField("amount_max", "request.amount_range_invalid", "Amount max must be greater than or equal to amount min."),
		})
		return nil, nil, false
	}
	return optionalInt64Pointer(minValue, minPresent), optionalInt64Pointer(maxValue, maxPresent), true
}

func walletNonNegativeInt64Query(w http.ResponseWriter, r *http.Request, field string) (int64, bool, bool) {
	value, present, err := httpserver.ParseOptionalNonNegativeInt64Query(r, field)
	if err != nil {
		httpserver.WriteValidationError(w, r, []httpserver.ValidationField{
			validationField(field, "request.amount_invalid", "Amount must be a non-negative number."),
		})
		return 0, present, false
	}
	return value, present, true
}

func optionalInt64Pointer(value int64, present bool) *int64 {
	if !present {
		return nil
	}
	return &value
}
