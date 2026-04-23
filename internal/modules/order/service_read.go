package order

import "github.com/Chinsusu/Billing-V2/internal/modules/tenant"

const defaultServiceInstanceListLimit = 100
const maxServiceInstanceListLimit = 500

func normalizeServiceInstanceFilter(filter ServiceInstanceFilter) ServiceInstanceFilter {
	if filter.Limit <= 0 {
		filter.Limit = defaultServiceInstanceListLimit
	}
	if filter.Limit > maxServiceInstanceListLimit {
		filter.Limit = maxServiceInstanceListLimit
	}
	return filter
}

func validateServiceInstanceFilter(filter ServiceInstanceFilter) error {
	if filter.TenantID.Empty() {
		return tenant.ErrTenantIDMissing
	}
	if filter.Status != "" && !filter.Status.Valid() {
		return ErrServiceStatusInvalid
	}
	return nil
}

func validateServiceInstanceLookup(lookup ServiceInstanceLookup) error {
	if lookup.ID.Empty() {
		return ErrServiceIDMissing
	}
	if lookup.TenantID.Empty() {
		return tenant.ErrTenantIDMissing
	}
	return nil
}
