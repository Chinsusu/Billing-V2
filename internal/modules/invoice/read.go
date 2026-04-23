package invoice

import "github.com/Chinsusu/Billing-V2/internal/modules/tenant"

const defaultInvoiceListLimit = 100
const maxInvoiceListLimit = 500

func normalizeInvoiceFilter(filter InvoiceFilter) InvoiceFilter {
	if filter.Limit <= 0 {
		filter.Limit = defaultInvoiceListLimit
	}
	if filter.Limit > maxInvoiceListLimit {
		filter.Limit = maxInvoiceListLimit
	}
	return filter
}

func validateInvoiceFilter(filter InvoiceFilter) error {
	if filter.TenantID.Empty() {
		return tenant.ErrTenantIDMissing
	}
	if filter.Status != "" && !filter.Status.Valid() {
		return ErrStatusInvalid
	}
	if amountRangeInvalid(filter.AmountMinMinor, filter.AmountMaxMinor) {
		return ErrAmountInvalid
	}
	return nil
}

func validateInvoiceLookup(lookup InvoiceLookup) error {
	if lookup.ID.Empty() {
		return ErrInvoiceIDMissing
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
