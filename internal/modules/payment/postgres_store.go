package payment

import (
	"context"
	"database/sql"

	"github.com/Chinsusu/Billing-V2/internal/modules/invoice"
	"github.com/Chinsusu/Billing-V2/internal/modules/wallet"
	platformdb "github.com/Chinsusu/Billing-V2/internal/platform/db"
)

type PostgresStore struct {
	executor platformdb.Executor
}

var _ Store = (*PostgresStore)(nil)

func NewPostgresStore(executor platformdb.Executor) *PostgresStore {
	return &PostgresStore{executor: executor}
}

const transactionColumns = `payment_transaction_id, display_id, tenant_id, account_user_id, order_id, invoice_id, transaction_type, status, currency, amount_minor, description, idempotency_key, metadata, created_at, updated_at`

const createTransactionSQL = `
INSERT INTO payment_transactions (tenant_id, account_user_id, order_id, invoice_id, transaction_type, status, currency, amount_minor, description, idempotency_key, metadata)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11::jsonb)
ON CONFLICT (tenant_id, idempotency_key)
DO UPDATE SET idempotency_key = EXCLUDED.idempotency_key
RETURNING ` + transactionColumns

func (store *PostgresStore) CreateTransaction(ctx context.Context, input CreateTransactionInput) (Transaction, error) {
	if err := store.ready(); err != nil {
		return Transaction{}, err
	}
	args, err := createTransactionArgs(input)
	if err != nil {
		return Transaction{}, err
	}
	return scanTransaction(store.executor.QueryRowContext(ctx, createTransactionSQL, args...))
}

func (store *PostgresStore) PayInvoiceFromWallet(ctx context.Context, input PayInvoiceFromWalletInput) (WalletInvoicePayment, error) {
	if err := store.ready(); err != nil {
		return WalletInvoicePayment{}, err
	}
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return WalletInvoicePayment{}, err
	}
	if conn, ok := store.executor.(*sql.DB); ok {
		var result WalletInvoicePayment
		err := platformdb.WithTx(ctx, conn, func(ctx context.Context, tx *sql.Tx) error {
			var runErr error
			result, runErr = payInvoiceFromWallet(
				ctx,
				NewPostgresStore(tx),
				invoice.NewPostgresStore(tx),
				wallet.NewService(wallet.NewPostgresStore(tx)),
				input,
			)
			return runErr
		})
		if err != nil {
			return WalletInvoicePayment{}, err
		}
		return result, nil
	}
	return payInvoiceFromWallet(
		ctx,
		store,
		invoice.NewPostgresStore(store.executor),
		wallet.NewService(wallet.NewPostgresStore(store.executor)),
		input,
	)
}

func createTransactionArgs(input CreateTransactionInput) ([]interface{}, error) {
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return nil, err
	}
	return []interface{}{
		input.TenantID, input.AccountUserID, nullableString(string(input.OrderID)), nullableString(string(input.InvoiceID)),
		input.Type, input.Status, input.Currency, input.AmountMinor, nullableString(input.Description),
		input.IdempotencyKey, string(input.Metadata),
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
