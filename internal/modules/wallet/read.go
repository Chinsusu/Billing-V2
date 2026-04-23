package wallet

import "github.com/Chinsusu/Billing-V2/internal/modules/tenant"

const defaultLedgerEntryListLimit = 100
const maxLedgerEntryListLimit = 500

func normalizeLedgerEntryFilter(filter LedgerEntryFilter) LedgerEntryFilter {
	if filter.Limit <= 0 {
		filter.Limit = defaultLedgerEntryListLimit
	}
	if filter.Limit > maxLedgerEntryListLimit {
		filter.Limit = maxLedgerEntryListLimit
	}
	return filter
}

func validateLedgerEntryFilter(filter LedgerEntryFilter) error {
	if filter.TenantID.Empty() {
		return tenant.ErrTenantIDMissing
	}
	if filter.WalletID.Empty() {
		return ErrWalletIDMissing
	}
	if filter.Direction != "" && !filter.Direction.Valid() {
		return ErrDirectionInvalid
	}
	if filter.EntryType != "" && !filter.EntryType.Valid() {
		return ErrEntryTypeInvalid
	}
	if filter.Status != "" && !filter.Status.Valid() {
		return ErrLedgerStatusInvalid
	}
	return nil
}

func validateLedgerEntryLookup(lookup LedgerEntryLookup) error {
	if lookup.ID.Empty() {
		return ErrLedgerEntryIDMissing
	}
	if lookup.TenantID.Empty() {
		return tenant.ErrTenantIDMissing
	}
	if lookup.WalletID.Empty() {
		return ErrWalletIDMissing
	}
	return nil
}
