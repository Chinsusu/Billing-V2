package payment

import (
	"context"
	"fmt"
)

const transactionReadColumns = `txn.payment_transaction_id, txn.display_id, txn.tenant_id, txn.account_user_id, txn.order_id, txn.invoice_id, txn.transaction_type, txn.status, txn.currency, txn.amount_minor, txn.description, txn.idempotency_key, txn.metadata, txn.created_at, txn.updated_at`

func (store *PostgresStore) ListTransactions(ctx context.Context, filter TransactionFilter) ([]Transaction, error) {
	if err := store.ready(); err != nil {
		return nil, err
	}
	query, args, err := buildListTransactionsQuery(filter)
	if err != nil {
		return nil, err
	}
	rows, err := store.executor.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list payment transactions: %w", err)
	}
	defer rows.Close()
	transactions := make([]Transaction, 0)
	for rows.Next() {
		transaction, err := scanTransaction(rows)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, transaction)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("read payment transactions: %w", err)
	}
	return transactions, nil
}

func (store *PostgresStore) GetTransaction(ctx context.Context, lookup TransactionLookup) (Transaction, error) {
	if err := store.ready(); err != nil {
		return Transaction{}, err
	}
	query, args, err := buildGetTransactionQuery(lookup)
	if err != nil {
		return Transaction{}, err
	}
	return scanTransaction(store.executor.QueryRowContext(ctx, query, args...))
}

func buildListTransactionsQuery(filter TransactionFilter) (string, []interface{}, error) {
	filter = normalizeTransactionFilter(filter)
	if err := validateTransactionFilter(filter); err != nil {
		return "", nil, err
	}
	query := `SELECT ` + transactionReadColumns + `
FROM payment_transactions txn
WHERE txn.tenant_id = $1`
	args := []interface{}{filter.TenantID}
	if filter.AccountUserID != "" {
		args = append(args, filter.AccountUserID)
		query += fmt.Sprintf("\n  AND txn.account_user_id = $%d", len(args))
	}
	if filter.AccountDisplayID > 0 {
		args = append(args, filter.AccountDisplayID)
		query += fmt.Sprintf(`
  AND EXISTS (
    SELECT 1
    FROM users account
    WHERE account.user_id = txn.account_user_id
      AND account.tenant_id = txn.tenant_id
      AND account.display_id = $%d
  )`, len(args))
	}
	if filter.DisplayID > 0 {
		args = append(args, filter.DisplayID)
		query += fmt.Sprintf("\n  AND txn.display_id = $%d", len(args))
	}
	if filter.OrderID != "" {
		args = append(args, filter.OrderID)
		query += fmt.Sprintf("\n  AND txn.order_id = $%d", len(args))
	}
	if filter.OrderDisplayID > 0 {
		args = append(args, filter.OrderDisplayID)
		query += fmt.Sprintf(`
  AND EXISTS (
    SELECT 1
    FROM orders ord
    WHERE ord.order_id = txn.order_id
      AND ord.tenant_id = txn.tenant_id
      AND ord.display_id = $%d
  )`, len(args))
	}
	if filter.InvoiceID != "" {
		args = append(args, filter.InvoiceID)
		query += fmt.Sprintf("\n  AND txn.invoice_id = $%d", len(args))
	}
	if filter.InvoiceDisplayID > 0 {
		args = append(args, filter.InvoiceDisplayID)
		query += fmt.Sprintf(`
  AND EXISTS (
    SELECT 1
    FROM invoices inv
    WHERE inv.invoice_id = txn.invoice_id
      AND inv.tenant_id = txn.tenant_id
      AND inv.display_id = $%d
  )`, len(args))
	}
	if filter.Type != "" {
		args = append(args, filter.Type)
		query += fmt.Sprintf("\n  AND txn.transaction_type = $%d", len(args))
	}
	if filter.Status != "" {
		args = append(args, filter.Status)
		query += fmt.Sprintf("\n  AND txn.status = $%d", len(args))
	}
	if filter.AmountMinMinor != nil {
		args = append(args, *filter.AmountMinMinor)
		query += fmt.Sprintf("\n  AND txn.amount_minor >= $%d", len(args))
	}
	if filter.AmountMaxMinor != nil {
		args = append(args, *filter.AmountMaxMinor)
		query += fmt.Sprintf("\n  AND txn.amount_minor <= $%d", len(args))
	}
	args = append(args, filter.Limit)
	query += fmt.Sprintf("\nORDER BY txn.created_at DESC\nLIMIT $%d", len(args))
	return query, args, nil
}

func buildGetTransactionQuery(lookup TransactionLookup) (string, []interface{}, error) {
	lookup = normalizeTransactionLookup(lookup)
	if err := validateTransactionLookup(lookup); err != nil {
		return "", nil, err
	}
	if lookup.ID.Empty() {
		query := `SELECT ` + transactionReadColumns + `
FROM payment_transactions txn
WHERE txn.tenant_id = $1
  AND txn.idempotency_key = $2`
		args := []interface{}{lookup.TenantID, lookup.IdempotencyKey}
		if lookup.AccountUserID != "" {
			args = append(args, lookup.AccountUserID)
			query += fmt.Sprintf("\n  AND txn.account_user_id = $%d", len(args))
		}
		return query, args, nil
	}
	query := `SELECT ` + transactionReadColumns + `
FROM payment_transactions txn
WHERE txn.payment_transaction_id = $1
  AND txn.tenant_id = $2`
	args := []interface{}{lookup.ID, lookup.TenantID}
	if lookup.AccountUserID != "" {
		args = append(args, lookup.AccountUserID)
		query += fmt.Sprintf("\n  AND txn.account_user_id = $%d", len(args))
	}
	return query, args, nil
}
