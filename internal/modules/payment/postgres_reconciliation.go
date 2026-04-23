package payment

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/invoice"
	"github.com/Chinsusu/Billing-V2/internal/modules/order"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
	"github.com/Chinsusu/Billing-V2/internal/modules/wallet"
)

const reconciliationReadColumns = transactionReadColumns + `,
COALESCE(txn.metadata->>'provider', '') AS provider,
inv.invoice_id, inv.display_id, inv.status, inv.total_minor, inv.paid_at,
ledger.ledger_entry_id, ledger.display_id, ledger.wallet_id, ledger.direction, ledger.entry_type, ledger.status, ledger.balance_after_minor,
linked_wallet.display_id`

func (store *PostgresStore) ListPaymentReconciliations(ctx context.Context, filter ReconciliationFilter) ([]PaymentReconciliation, error) {
	if err := store.ready(); err != nil {
		return nil, err
	}
	query, args, err := buildListPaymentReconciliationsQuery(filter)
	if err != nil {
		return nil, err
	}
	rows, err := store.executor.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list payment reconciliations: %w", err)
	}
	defer rows.Close()
	records := make([]PaymentReconciliation, 0)
	for rows.Next() {
		record, err := scanPaymentReconciliation(rows)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("read payment reconciliations: %w", err)
	}
	return records, nil
}

func (store *PostgresStore) GetPaymentReconciliation(ctx context.Context, lookup ReconciliationLookup) (PaymentReconciliation, error) {
	if err := store.ready(); err != nil {
		return PaymentReconciliation{}, err
	}
	query, args, err := buildGetPaymentReconciliationQuery(lookup)
	if err != nil {
		return PaymentReconciliation{}, err
	}
	return scanPaymentReconciliation(store.executor.QueryRowContext(ctx, query, args...))
}

func buildListPaymentReconciliationsQuery(filter ReconciliationFilter) (string, []interface{}, error) {
	filter = normalizeReconciliationFilter(filter)
	if err := validateReconciliationFilter(filter); err != nil {
		return "", nil, err
	}
	query := paymentReconciliationBaseQuery() + `
WHERE txn.tenant_id = $1`
	args := []interface{}{filter.TenantID}
	if filter.Status != "" {
		args = append(args, filter.Status)
		query += fmt.Sprintf("\n  AND txn.status = $%d", len(args))
	}
	if filter.Provider != "" {
		args = append(args, filter.Provider)
		query += fmt.Sprintf("\n  AND txn.metadata ->> 'provider' = $%d", len(args))
	}
	if filter.InvoiceID != "" {
		args = append(args, filter.InvoiceID)
		query += fmt.Sprintf("\n  AND txn.invoice_id = $%d", len(args))
	}
	if !filter.WalletID.Empty() {
		args = append(args, filter.WalletID)
		query += fmt.Sprintf("\n  AND ledger.wallet_id = $%d", len(args))
	}
	if !filter.CreatedFrom.IsZero() {
		args = append(args, filter.CreatedFrom)
		query += fmt.Sprintf("\n  AND txn.created_at >= $%d", len(args))
	}
	if !filter.CreatedTo.IsZero() {
		args = append(args, filter.CreatedTo)
		query += fmt.Sprintf("\n  AND txn.created_at <= $%d", len(args))
	}
	args = append(args, filter.Limit)
	query += fmt.Sprintf("\nORDER BY txn.created_at DESC\nLIMIT $%d", len(args))
	return query, args, nil
}

func buildGetPaymentReconciliationQuery(lookup ReconciliationLookup) (string, []interface{}, error) {
	if err := validateReconciliationLookup(lookup); err != nil {
		return "", nil, err
	}
	query := paymentReconciliationBaseQuery() + `
WHERE txn.tenant_id = $1
  AND txn.payment_transaction_id = $2`
	return query, []interface{}{lookup.TenantID, lookup.TransactionID}, nil
}

func paymentReconciliationBaseQuery() string {
	return `SELECT ` + reconciliationReadColumns + `
FROM payment_transactions txn
LEFT JOIN invoices inv
  ON inv.invoice_id = txn.invoice_id
 AND inv.tenant_id = txn.tenant_id
LEFT JOIN LATERAL (
    SELECT ledger_entry_id, display_id, wallet_id, direction, entry_type, status, balance_after_minor
    FROM wallet_ledger_entries ledger
    WHERE ledger.tenant_id = txn.tenant_id
      AND ledger.reference_type = 'invoice'
      AND ledger.reference_id = txn.invoice_id
      AND ledger.entry_type = 'purchase'
    ORDER BY ledger.created_at DESC
    LIMIT 1
) ledger ON TRUE
LEFT JOIN wallets linked_wallet
  ON linked_wallet.wallet_id = ledger.wallet_id
 AND linked_wallet.tenant_id = txn.tenant_id`
}

func scanPaymentReconciliation(row transactionScanner) (PaymentReconciliation, error) {
	var record PaymentReconciliation
	var id, tenantID, accountUserID, transactionType, status, idempotencyKey string
	var orderID, transactionInvoiceID, description sql.NullString
	var metadata []byte
	var invoiceID, invoiceStatus sql.NullString
	var invoiceDisplayID, invoiceTotalMinor sql.NullInt64
	var invoicePaidAt sql.NullTime
	var ledgerID, ledgerWalletID, ledgerDirection, ledgerEntryType, ledgerStatus sql.NullString
	var ledgerDisplayID, ledgerBalanceAfterMinor, walletDisplayID sql.NullInt64
	if err := row.Scan(
		&id, &record.Transaction.DisplayID, &tenantID, &accountUserID, &orderID, &transactionInvoiceID,
		&transactionType, &status, &record.Transaction.Currency, &record.Transaction.AmountMinor,
		&description, &idempotencyKey, &metadata, &record.Transaction.CreatedAt, &record.Transaction.UpdatedAt,
		&record.Provider,
		&invoiceID, &invoiceDisplayID, &invoiceStatus, &invoiceTotalMinor, &invoicePaidAt,
		&ledgerID, &ledgerDisplayID, &ledgerWalletID, &ledgerDirection, &ledgerEntryType, &ledgerStatus,
		&ledgerBalanceAfterMinor, &walletDisplayID,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return PaymentReconciliation{}, ErrTransactionNotFound
		}
		return PaymentReconciliation{}, fmt.Errorf("scan payment reconciliation: %w", err)
	}
	record.Transaction.ID = TransactionID(id)
	record.Transaction.TenantID = tenant.ID(tenantID)
	record.Transaction.AccountUserID = identity.UserID(accountUserID)
	record.Transaction.OrderID = order.OrderID(orderID.String)
	record.Transaction.InvoiceID = invoice.InvoiceID(transactionInvoiceID.String)
	record.Transaction.Type = TransactionType(transactionType)
	record.Transaction.Status = TransactionStatus(status)
	record.Transaction.Description = description.String
	record.Transaction.IdempotencyKey = IdempotencyKey(idempotencyKey)
	record.Transaction.Metadata = append(record.Transaction.Metadata, metadata...)
	if invoiceID.Valid {
		record.Invoice = ReconciliationInvoice{
			ID:         invoice.InvoiceID(invoiceID.String),
			DisplayID:  invoiceDisplayID.Int64,
			Status:     invoice.Status(invoiceStatus.String),
			TotalMinor: invoiceTotalMinor.Int64,
		}
		if invoicePaidAt.Valid {
			record.Invoice.PaidAt = invoicePaidAt.Time
		}
	}
	if ledgerID.Valid {
		record.Ledger = ReconciliationLedger{
			ID:                wallet.LedgerEntryID(ledgerID.String),
			DisplayID:         ledgerDisplayID.Int64,
			WalletID:          wallet.WalletID(ledgerWalletID.String),
			WalletDisplayID:   walletDisplayID.Int64,
			Direction:         wallet.Direction(ledgerDirection.String),
			EntryType:         wallet.EntryType(ledgerEntryType.String),
			Status:            wallet.LedgerStatus(ledgerStatus.String),
			BalanceAfterMinor: ledgerBalanceAfterMinor.Int64,
		}
	}
	return record, nil
}
