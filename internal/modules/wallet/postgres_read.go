package wallet

import (
	"context"
	"fmt"
)

const ledgerEntryReadColumns = `entry.ledger_entry_id, entry.display_id, entry.wallet_id, entry.tenant_id, entry.direction, entry.amount_minor, entry.currency, entry.entry_type, entry.status, entry.balance_after_minor, entry.reference_type, entry.reference_id, entry.idempotency_key, entry.created_by, entry.reason, entry.correlation_id, entry.created_at`
const topupRequestReadColumns = `topup.topup_request_id, topup.display_id, topup.tenant_id, topup.wallet_id, topup.requested_by, topup.amount_minor, topup.currency, topup.payment_method, topup.payment_reference, topup.status, topup.reviewed_by, topup.reviewed_at, topup.review_note, topup.ledger_entry_id, topup.idempotency_key, topup.created_at, topup.updated_at`
const walletReadColumns = `wallet.wallet_id, wallet.display_id, wallet.tenant_id, wallet.owner_type, wallet.owner_id, wallet.currency, wallet.status, wallet.available_balance_minor, wallet.locked_balance_minor, wallet.metadata, wallet.created_at, wallet.updated_at`

func (store *PostgresStore) ListWallets(ctx context.Context, filter WalletFilter) ([]Wallet, error) {
	if err := store.ready(); err != nil {
		return nil, err
	}
	query, args, err := buildListWalletsQuery(filter)
	if err != nil {
		return nil, err
	}
	rows, err := store.executor.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list wallets: %w", err)
	}
	defer rows.Close()
	wallets := make([]Wallet, 0)
	for rows.Next() {
		wallet, err := scanWallet(rows)
		if err != nil {
			return nil, err
		}
		wallets = append(wallets, wallet)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("read wallets: %w", err)
	}
	return wallets, nil
}

func (store *PostgresStore) GetWallet(ctx context.Context, lookup WalletLookup) (Wallet, error) {
	if err := store.ready(); err != nil {
		return Wallet{}, err
	}
	query, args, err := buildGetWalletQuery(lookup)
	if err != nil {
		return Wallet{}, err
	}
	return scanWallet(store.executor.QueryRowContext(ctx, query, args...))
}

func buildListWalletsQuery(filter WalletFilter) (string, []interface{}, error) {
	filter = normalizeWalletFilter(filter)
	if err := validateWalletFilter(filter); err != nil {
		return "", nil, err
	}
	query := `SELECT ` + walletReadColumns + `
FROM wallets wallet
WHERE wallet.tenant_id = $1`
	args := []interface{}{filter.TenantID}
	if filter.DisplayID > 0 {
		args = append(args, filter.DisplayID)
		query += fmt.Sprintf("\n  AND wallet.display_id = $%d", len(args))
	}
	if filter.OwnerType != "" {
		args = append(args, filter.OwnerType)
		query += fmt.Sprintf("\n  AND wallet.owner_type = $%d", len(args))
	}
	if filter.OwnerID != "" {
		args = append(args, filter.OwnerID)
		query += fmt.Sprintf("\n  AND wallet.owner_id = $%d", len(args))
	}
	if filter.Status != "" {
		args = append(args, filter.Status)
		query += fmt.Sprintf("\n  AND wallet.status = $%d", len(args))
	}
	args = append(args, filter.Limit)
	query += fmt.Sprintf("\nORDER BY wallet.created_at DESC\nLIMIT $%d", len(args))
	return query, args, nil
}

func buildGetWalletQuery(lookup WalletLookup) (string, []interface{}, error) {
	if err := validateWalletLookup(lookup); err != nil {
		return "", nil, err
	}
	query := `SELECT ` + walletReadColumns + `
FROM wallets wallet
WHERE wallet.wallet_id = $1
  AND wallet.tenant_id = $2`
	args := []interface{}{lookup.ID, lookup.TenantID}
	if lookup.OwnerType != "" {
		args = append(args, lookup.OwnerType)
		query += fmt.Sprintf("\n  AND wallet.owner_type = $%d", len(args))
	}
	if lookup.OwnerID != "" {
		args = append(args, lookup.OwnerID)
		query += fmt.Sprintf("\n  AND wallet.owner_id = $%d", len(args))
	}
	return query, args, nil
}

func (store *PostgresStore) ListLedgerEntries(ctx context.Context, filter LedgerEntryFilter) ([]LedgerEntry, error) {
	if err := store.ready(); err != nil {
		return nil, err
	}
	query, args, err := buildListLedgerEntriesQuery(filter)
	if err != nil {
		return nil, err
	}
	rows, err := store.executor.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list wallet ledger entries: %w", err)
	}
	defer rows.Close()
	entries := make([]LedgerEntry, 0)
	for rows.Next() {
		entry, err := scanLedgerEntry(rows)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("read wallet ledger entries: %w", err)
	}
	return entries, nil
}

func (store *PostgresStore) GetLedgerEntry(ctx context.Context, lookup LedgerEntryLookup) (LedgerEntry, error) {
	if err := store.ready(); err != nil {
		return LedgerEntry{}, err
	}
	query, args, err := buildGetLedgerEntryQuery(lookup)
	if err != nil {
		return LedgerEntry{}, err
	}
	return scanLedgerEntry(store.executor.QueryRowContext(ctx, query, args...))
}

func buildListLedgerEntriesQuery(filter LedgerEntryFilter) (string, []interface{}, error) {
	filter = normalizeLedgerEntryFilter(filter)
	if err := validateLedgerEntryFilter(filter); err != nil {
		return "", nil, err
	}
	query := `SELECT ` + ledgerEntryReadColumns + `
FROM wallet_ledger_entries entry
WHERE entry.tenant_id = $1
  AND entry.wallet_id = $2`
	args := []interface{}{filter.TenantID, filter.WalletID}
	if filter.DisplayID > 0 {
		args = append(args, filter.DisplayID)
		query += fmt.Sprintf("\n  AND entry.display_id = $%d", len(args))
	}
	if filter.Direction != "" {
		args = append(args, filter.Direction)
		query += fmt.Sprintf("\n  AND entry.direction = $%d", len(args))
	}
	if filter.EntryType != "" {
		args = append(args, filter.EntryType)
		query += fmt.Sprintf("\n  AND entry.entry_type = $%d", len(args))
	}
	if filter.Status != "" {
		args = append(args, filter.Status)
		query += fmt.Sprintf("\n  AND entry.status = $%d", len(args))
	}
	if filter.AmountMinMinor != nil {
		args = append(args, *filter.AmountMinMinor)
		query += fmt.Sprintf("\n  AND entry.amount_minor >= $%d", len(args))
	}
	if filter.AmountMaxMinor != nil {
		args = append(args, *filter.AmountMaxMinor)
		query += fmt.Sprintf("\n  AND entry.amount_minor <= $%d", len(args))
	}
	args = append(args, filter.Limit)
	query += fmt.Sprintf("\nORDER BY entry.created_at DESC\nLIMIT $%d", len(args))
	return query, args, nil
}

func buildGetLedgerEntryQuery(lookup LedgerEntryLookup) (string, []interface{}, error) {
	if err := validateLedgerEntryLookup(lookup); err != nil {
		return "", nil, err
	}
	query := `SELECT ` + ledgerEntryReadColumns + `
FROM wallet_ledger_entries entry
WHERE entry.ledger_entry_id = $1
  AND entry.tenant_id = $2
  AND entry.wallet_id = $3`
	return query, []interface{}{lookup.ID, lookup.TenantID, lookup.WalletID}, nil
}

func (store *PostgresStore) ListTopupRequests(ctx context.Context, filter TopupRequestFilter) ([]TopupRequest, error) {
	if err := store.ready(); err != nil {
		return nil, err
	}
	query, args, err := buildListTopupRequestsQuery(filter)
	if err != nil {
		return nil, err
	}
	rows, err := store.executor.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list wallet top-up requests: %w", err)
	}
	defer rows.Close()
	requests := make([]TopupRequest, 0)
	for rows.Next() {
		request, err := scanTopupRequest(rows)
		if err != nil {
			return nil, err
		}
		requests = append(requests, request)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("read wallet top-up requests: %w", err)
	}
	return requests, nil
}

func (store *PostgresStore) GetTopupRequest(ctx context.Context, lookup TopupRequestLookup) (TopupRequest, error) {
	if err := store.ready(); err != nil {
		return TopupRequest{}, err
	}
	query, args, err := buildGetTopupRequestQuery(lookup)
	if err != nil {
		return TopupRequest{}, err
	}
	return scanTopupRequest(store.executor.QueryRowContext(ctx, query, args...))
}

func buildListTopupRequestsQuery(filter TopupRequestFilter) (string, []interface{}, error) {
	filter = normalizeTopupRequestFilter(filter)
	if err := validateTopupRequestFilter(filter); err != nil {
		return "", nil, err
	}
	query := `SELECT ` + topupRequestReadColumns + `
FROM topup_requests topup
WHERE topup.tenant_id = $1`
	args := []interface{}{filter.TenantID}
	if filter.DisplayID > 0 {
		args = append(args, filter.DisplayID)
		query += fmt.Sprintf("\n  AND topup.display_id = $%d", len(args))
	}
	if filter.WalletID != "" {
		args = append(args, filter.WalletID)
		query += fmt.Sprintf("\n  AND topup.wallet_id = $%d", len(args))
	}
	if filter.WalletDisplayID > 0 {
		args = append(args, filter.WalletDisplayID)
		query += fmt.Sprintf(`
  AND EXISTS (
    SELECT 1
    FROM wallets wallet
    WHERE wallet.wallet_id = topup.wallet_id
      AND wallet.tenant_id = topup.tenant_id
      AND wallet.display_id = $%d
  )`, len(args))
	}
	if filter.RequestedBy != "" {
		args = append(args, filter.RequestedBy)
		query += fmt.Sprintf("\n  AND topup.requested_by = $%d", len(args))
	}
	if filter.RequestedByDisplayID > 0 {
		args = append(args, filter.RequestedByDisplayID)
		query += fmt.Sprintf(`
  AND EXISTS (
    SELECT 1
    FROM users requester
    WHERE requester.user_id = topup.requested_by
      AND requester.tenant_id = topup.tenant_id
      AND requester.display_id = $%d
  )`, len(args))
	}
	if filter.PaymentMethod != "" {
		args = append(args, filter.PaymentMethod)
		query += fmt.Sprintf("\n  AND topup.payment_method = $%d", len(args))
	}
	if filter.Status != "" {
		args = append(args, filter.Status)
		query += fmt.Sprintf("\n  AND topup.status = $%d", len(args))
	}
	if filter.AmountMinMinor != nil {
		args = append(args, *filter.AmountMinMinor)
		query += fmt.Sprintf("\n  AND topup.amount_minor >= $%d", len(args))
	}
	if filter.AmountMaxMinor != nil {
		args = append(args, *filter.AmountMaxMinor)
		query += fmt.Sprintf("\n  AND topup.amount_minor <= $%d", len(args))
	}
	args = append(args, filter.Limit)
	query += fmt.Sprintf("\nORDER BY topup.created_at DESC\nLIMIT $%d", len(args))
	return query, args, nil
}

func buildGetTopupRequestQuery(lookup TopupRequestLookup) (string, []interface{}, error) {
	if err := validateTopupRequestLookup(lookup); err != nil {
		return "", nil, err
	}
	query := `SELECT ` + topupRequestReadColumns + `
FROM topup_requests topup
WHERE topup.topup_request_id = $1
  AND topup.tenant_id = $2`
	args := []interface{}{lookup.ID, lookup.TenantID}
	if lookup.RequestedBy != "" {
		args = append(args, lookup.RequestedBy)
		query += fmt.Sprintf("\n  AND topup.requested_by = $%d", len(args))
	}
	return query, args, nil
}
