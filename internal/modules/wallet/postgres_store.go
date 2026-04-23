package wallet

import (
	"context"
	"database/sql"
	"errors"

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
const topupRequestColumns = `topup_request_id, display_id, tenant_id, wallet_id, requested_by, amount_minor, currency, payment_method, payment_reference, status, reviewed_by, reviewed_at, review_note, ledger_entry_id, idempotency_key, created_at, updated_at`

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

const postLedgerEntrySQL = `
WITH existing AS (
SELECT ` + ledgerEntryColumns + `
FROM wallet_ledger_entries
WHERE wallet_id = $1 AND idempotency_key = $9
), updated_wallet AS (
UPDATE wallets wallet
SET available_balance_minor = CASE
        WHEN $3::wallet_ledger_direction = 'credit' THEN wallet.available_balance_minor + $4
        ELSE wallet.available_balance_minor - $4
    END,
    updated_at = NOW()
WHERE wallet.wallet_id = $1
  AND wallet.tenant_id = $2
  AND wallet.currency = $5
  AND wallet.status = 'active'
  AND NOT EXISTS (SELECT 1 FROM existing)
  AND ($3::wallet_ledger_direction = 'credit' OR wallet.available_balance_minor >= $4)
RETURNING wallet.available_balance_minor
), inserted AS (
INSERT INTO wallet_ledger_entries (wallet_id, tenant_id, direction, amount_minor, currency, entry_type, status, balance_after_minor, reference_type, reference_id, idempotency_key, created_by, reason, correlation_id)
SELECT $1, $2, $3, $4, $5, $6, 'posted', updated_wallet.available_balance_minor, $7, $8, $9, $10, $11, $12
FROM updated_wallet
ON CONFLICT (wallet_id, idempotency_key)
DO NOTHING
RETURNING ` + ledgerEntryColumns + `
)
SELECT ` + ledgerEntryColumns + ` FROM inserted
UNION ALL
SELECT ` + ledgerEntryColumns + ` FROM existing
LIMIT 1`

func (store *PostgresStore) PostLedgerEntry(ctx context.Context, input PostLedgerEntryInput) (LedgerEntry, error) {
	if err := store.ready(); err != nil {
		return LedgerEntry{}, err
	}
	args, err := postLedgerEntryArgs(input)
	if err != nil {
		return LedgerEntry{}, err
	}
	entry, err := scanLedgerEntry(store.executor.QueryRowContext(ctx, postLedgerEntrySQL, args...))
	if errors.Is(err, ErrLedgerEntryNotFound) && input.Normalize().Direction == DirectionDebit {
		return LedgerEntry{}, ErrInsufficientBalance
	}
	return entry, err
}

func postLedgerEntryArgs(input PostLedgerEntryInput) ([]interface{}, error) {
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return nil, err
	}
	return []interface{}{
		input.WalletID, input.TenantID, input.Direction, input.AmountMinor, input.Currency,
		input.EntryType, input.ReferenceType, input.ReferenceID, input.IdempotencyKey,
		nullableString(string(input.CreatedBy)), nullableString(input.Reason), input.CorrelationID,
	}, nil
}

const createTopupRequestSQL = `
WITH inserted AS (
INSERT INTO topup_requests (tenant_id, wallet_id, requested_by, amount_minor, currency, payment_method, payment_reference, status, idempotency_key)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
ON CONFLICT (tenant_id, requested_by, idempotency_key)
DO NOTHING
RETURNING ` + topupRequestColumns + `
)
SELECT ` + topupRequestColumns + ` FROM inserted
UNION ALL
SELECT ` + topupRequestColumns + `
FROM topup_requests
WHERE tenant_id = $1 AND requested_by = $3 AND idempotency_key = $9
LIMIT 1`

func (store *PostgresStore) CreateTopupRequest(ctx context.Context, input CreateTopupRequestInput) (TopupRequest, error) {
	if err := store.ready(); err != nil {
		return TopupRequest{}, err
	}
	args, err := createTopupRequestArgs(input)
	if err != nil {
		return TopupRequest{}, err
	}
	return scanTopupRequest(store.executor.QueryRowContext(ctx, createTopupRequestSQL, args...))
}

func createTopupRequestArgs(input CreateTopupRequestInput) ([]interface{}, error) {
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return nil, err
	}
	return []interface{}{
		input.TenantID, input.WalletID, input.RequestedBy, input.AmountMinor, input.Currency,
		input.PaymentMethod, nullableString(input.PaymentReference), input.Status, input.IdempotencyKey,
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
