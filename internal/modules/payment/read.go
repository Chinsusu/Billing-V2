package payment

import "github.com/Chinsusu/Billing-V2/internal/modules/tenant"

const defaultTransactionListLimit = 100
const maxTransactionListLimit = 500

func normalizeTransactionFilter(filter TransactionFilter) TransactionFilter {
	if filter.Limit <= 0 {
		filter.Limit = defaultTransactionListLimit
	}
	if filter.Limit > maxTransactionListLimit {
		filter.Limit = maxTransactionListLimit
	}
	return filter
}

func validateTransactionFilter(filter TransactionFilter) error {
	if filter.TenantID.Empty() {
		return tenant.ErrTenantIDMissing
	}
	if filter.Type != "" && !filter.Type.Valid() {
		return ErrTypeInvalid
	}
	if filter.Status != "" && !filter.Status.Valid() {
		return ErrStatusInvalid
	}
	return nil
}

func validateTransactionLookup(lookup TransactionLookup) error {
	if lookup.ID.Empty() {
		return ErrTransactionIDMissing
	}
	if lookup.TenantID.Empty() {
		return tenant.ErrTenantIDMissing
	}
	return nil
}
