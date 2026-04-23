package wallet

import (
	"net/http"
	"strings"

	"github.com/Chinsusu/Billing-V2/internal/platform/httpserver"
)

func walletFilterFromRequest(w http.ResponseWriter, r *http.Request) (WalletFilter, httpserver.CursorPageRequest, bool) {
	page, ok := pageFromRequest(w, r)
	if !ok {
		return WalletFilter{}, httpserver.CursorPageRequest{}, false
	}
	filter := WalletFilter{Limit: page.Limit}
	query := r.URL.Query()
	if displayID, present, ok := walletPositiveInt64Query(w, r, "display_id"); !ok {
		return WalletFilter{}, httpserver.CursorPageRequest{}, false
	} else if present {
		filter.DisplayID = displayID
	}
	if ownerType := OwnerType(strings.TrimSpace(query.Get("owner_type"))); ownerType != "" {
		if !ownerType.Valid() {
			writeWalletError(w, r, ErrOwnerTypeInvalid)
			return WalletFilter{}, httpserver.CursorPageRequest{}, false
		}
		filter.OwnerType = ownerType
	}
	if ownerID := OwnerID(strings.TrimSpace(query.Get("owner_id"))); ownerID != "" {
		filter.OwnerID = ownerID
	}
	if status := Status(strings.TrimSpace(query.Get("status"))); status != "" {
		if !status.Valid() {
			writeWalletError(w, r, ErrStatusInvalid)
			return WalletFilter{}, httpserver.CursorPageRequest{}, false
		}
		filter.Status = status
	}
	return filter, page, true
}

func ledgerFilterFromRequest(w http.ResponseWriter, r *http.Request) (LedgerEntryFilter, httpserver.CursorPageRequest, bool) {
	page, ok := pageFromRequest(w, r)
	if !ok {
		return LedgerEntryFilter{}, httpserver.CursorPageRequest{}, false
	}
	filter := LedgerEntryFilter{Limit: page.Limit}
	query := r.URL.Query()
	if displayID, present, ok := walletPositiveInt64Query(w, r, "display_id"); !ok {
		return LedgerEntryFilter{}, httpserver.CursorPageRequest{}, false
	} else if present {
		filter.DisplayID = displayID
	}
	if direction := Direction(strings.TrimSpace(query.Get("direction"))); direction != "" {
		if !direction.Valid() {
			writeWalletError(w, r, ErrDirectionInvalid)
			return LedgerEntryFilter{}, httpserver.CursorPageRequest{}, false
		}
		filter.Direction = direction
	}
	if entryType := EntryType(strings.TrimSpace(query.Get("entry_type"))); entryType != "" {
		if !entryType.Valid() {
			writeWalletError(w, r, ErrEntryTypeInvalid)
			return LedgerEntryFilter{}, httpserver.CursorPageRequest{}, false
		}
		filter.EntryType = entryType
	}
	if status := LedgerStatus(strings.TrimSpace(query.Get("status"))); status != "" {
		if !status.Valid() {
			writeWalletError(w, r, ErrLedgerStatusInvalid)
			return LedgerEntryFilter{}, httpserver.CursorPageRequest{}, false
		}
		filter.Status = status
	}
	amountMin, amountMax, ok := walletAmountRangeQuery(w, r)
	if !ok {
		return LedgerEntryFilter{}, httpserver.CursorPageRequest{}, false
	}
	filter.AmountMinMinor = amountMin
	filter.AmountMaxMinor = amountMax
	return filter, page, true
}
