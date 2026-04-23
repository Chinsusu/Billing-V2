package wallet

import (
	"context"
	"database/sql"

	platformdb "github.com/Chinsusu/Billing-V2/internal/platform/db"
)

type PostgresStore struct {
	executor platformdb.Executor
}

var _ Store = (*PostgresStore)(nil)

func NewPostgresStore(executor platformdb.Executor) *PostgresStore {
	return &PostgresStore{executor: executor}
}

const ledgerEntryColumns = `ledger_entry_id, display_id, wallet_id, tenant_id, direction, amount_minor, currency, entry_type, status, balance_after_minor, reference_type, reference_id, idempotency_key, created_by, reason, correlation_id, created_at`

const createLedgerEntrySQL = `
WITH inserted AS (
INSERT INTO wallet_ledger_entries (wallet_id, tenant_id, direction, amount_minor, currency, entry_type, status, balance_after_minor, reference_type, reference_id, idempotency_key, created_by, reason, correlation_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
ON CONFLICT (wallet_id, idempotency_key)
DO NOTHING
RETURNING ` + ledgerEntryColumns + `
)
SELECT ` + ledgerEntryColumns + ` FROM inserted
UNION ALL
SELECT ` + ledgerEntryColumns + `
FROM wallet_ledger_entries
WHERE wallet_id = $1 AND idempotency_key = $11
LIMIT 1`

func (store *PostgresStore) CreateLedgerEntry(ctx context.Context, input CreateLedgerEntryInput) (LedgerEntry, error) {
	if err := store.ready(); err != nil {
		return LedgerEntry{}, err
	}
	args, err := createLedgerEntryArgs(input)
	if err != nil {
		return LedgerEntry{}, err
	}
	return scanLedgerEntry(store.executor.QueryRowContext(ctx, createLedgerEntrySQL, args...))
}

func createLedgerEntryArgs(input CreateLedgerEntryInput) ([]interface{}, error) {
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return nil, err
	}
	return []interface{}{
		input.WalletID, input.TenantID, input.Direction, input.AmountMinor, input.Currency,
		input.EntryType, input.Status, input.BalanceAfterMinor, input.ReferenceType, input.ReferenceID,
		input.IdempotencyKey, nullableString(string(input.CreatedBy)), nullableString(input.Reason), input.CorrelationID,
	}, nil
}

func (store *PostgresStore) ready() error {
	if store == nil || store.executor == nil {
		return ErrStoreExecutorMissing
	}
	return nil
}

func nullableString(value string) sql.NullString {
	if value == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: value, Valid: true}
}
