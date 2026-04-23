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
