package order

import "github.com/Chinsusu/Billing-V2/internal/modules/tenant"

const defaultOrderListLimit = 100
const maxOrderListLimit = 500

func normalizeOrderFilter(filter OrderFilter) OrderFilter {
	if filter.Limit <= 0 {
		filter.Limit = defaultOrderListLimit
	}
	if filter.Limit > maxOrderListLimit {
		filter.Limit = maxOrderListLimit
	}
	return filter
}

func validateOrderFilter(filter OrderFilter) error {
	if filter.TenantID.Empty() {
		return tenant.ErrTenantIDMissing
	}
	if filter.OrderStatus != "" && !filter.OrderStatus.Valid() {
		return ErrOrderStatusInvalid
	}
	if filter.BillingStatus != "" && !filter.BillingStatus.Valid() {
		return ErrBillingStatusInvalid
	}
	if amountRangeInvalid(filter.AmountMinMinor, filter.AmountMaxMinor) {
		return ErrAmountInvalid
	}
	return nil
}

func validateOrderLookup(lookup OrderLookup) error {
	if lookup.ID.Empty() {
		return ErrOrderIDMissing
	}
	if lookup.TenantID.Empty() {
		return tenant.ErrTenantIDMissing
	}
	return nil
}

func amountRangeInvalid(minorMin *int64, minorMax *int64) bool {
	if minorMin != nil && *minorMin < 0 {
		return true
	}
	if minorMax != nil && *minorMax < 0 {
		return true
	}
	return minorMin != nil && minorMax != nil && *minorMax < *minorMin
}
