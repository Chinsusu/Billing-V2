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
	if amountRangeInvalid(filter.AmountMinMinor, filter.AmountMaxMinor) {
		return ErrAmountInvalid
	}
	return nil
}

func validateTransactionLookup(lookup TransactionLookup) error {
	lookup = normalizeTransactionLookup(lookup)
	if lookup.ID.Empty() && lookup.IdempotencyKey == "" {
		return ErrTransactionIDMissing
	}
	if lookup.TenantID.Empty() {
		return tenant.ErrTenantIDMissing
	}
	return nil
}

func normalizeTransactionLookup(lookup TransactionLookup) TransactionLookup {
	output := lookup
	output.IdempotencyKey = IdempotencyKey(trim(string(output.IdempotencyKey)))
	return output
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
