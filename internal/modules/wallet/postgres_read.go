package wallet

import (
	"context"
	"fmt"
)

const ledgerEntryReadColumns = `entry.ledger_entry_id, entry.display_id, entry.wallet_id, entry.tenant_id, entry.direction, entry.amount_minor, entry.currency, entry.entry_type, entry.status, entry.balance_after_minor, entry.reference_type, entry.reference_id, entry.idempotency_key, entry.created_by, entry.reason, entry.correlation_id, entry.created_at`

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
