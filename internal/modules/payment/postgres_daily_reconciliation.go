package payment

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/Chinsusu/Billing-V2/internal/modules/invoice"
	"github.com/Chinsusu/Billing-V2/internal/modules/wallet"
)

const dailyWalletSummarySQL = `
WITH ledger_balance AS (
    SELECT wallet_id, tenant_id,
           COALESCE(SUM(CASE WHEN direction = 'credit' THEN amount_minor ELSE -amount_minor END), 0) AS ledger_balance_minor
    FROM wallet_ledger_entries
    WHERE tenant_id = $1
      AND status = 'posted'
    GROUP BY wallet_id, tenant_id
)
SELECT COUNT(*) AS wallets_checked
FROM wallets wallet
LEFT JOIN ledger_balance
  ON ledger_balance.wallet_id = wallet.wallet_id
 AND ledger_balance.tenant_id = wallet.tenant_id
WHERE wallet.tenant_id = $1`

const dailyWalletMismatchSQL = `
WITH ledger_balance AS (
    SELECT wallet_id, tenant_id,
           COALESCE(SUM(CASE WHEN direction = 'credit' THEN amount_minor ELSE -amount_minor END), 0) AS ledger_balance_minor
    FROM wallet_ledger_entries
    WHERE tenant_id = $1
      AND status = 'posted'
    GROUP BY wallet_id, tenant_id
)
SELECT wallet.wallet_id, wallet.display_id, wallet.currency, wallet.available_balance_minor,
       COALESCE(ledger_balance.ledger_balance_minor, 0) AS ledger_balance_minor,
       wallet.available_balance_minor - COALESCE(ledger_balance.ledger_balance_minor, 0) AS difference_minor,
       last_ledger.ledger_entry_id, last_ledger.display_id
FROM wallets wallet
LEFT JOIN ledger_balance
  ON ledger_balance.wallet_id = wallet.wallet_id
 AND ledger_balance.tenant_id = wallet.tenant_id
LEFT JOIN LATERAL (
    SELECT ledger_entry_id, display_id
    FROM wallet_ledger_entries entry
    WHERE entry.wallet_id = wallet.wallet_id
      AND entry.tenant_id = wallet.tenant_id
    ORDER BY entry.created_at DESC
    LIMIT 1
) last_ledger ON TRUE
WHERE wallet.tenant_id = $1
  AND wallet.available_balance_minor <> COALESCE(ledger_balance.ledger_balance_minor, 0)
ORDER BY wallet.display_id`

const dailyInvoiceCheckedSQL = `
SELECT COUNT(*)
FROM invoices invoice
WHERE invoice.tenant_id = $1
  AND (
    (invoice.created_at >= $2 AND invoice.created_at < $3)
    OR (invoice.paid_at >= $2 AND invoice.paid_at < $3)
    OR EXISTS (
        SELECT 1
        FROM payment_transactions txn
        WHERE txn.tenant_id = invoice.tenant_id
          AND txn.invoice_id = invoice.invoice_id
          AND txn.created_at >= $2
          AND txn.created_at < $3
    )
  )`

const dailyInvoiceMismatchSQL = `
WITH posted_charges AS (
    SELECT invoice_id,
           COUNT(*) AS transaction_count,
           COALESCE(SUM(amount_minor), 0) AS total_minor
    FROM payment_transactions
    WHERE tenant_id = $1
      AND invoice_id IS NOT NULL
      AND transaction_type = 'charge'
      AND status = 'posted'
    GROUP BY invoice_id
)
SELECT invoice.invoice_id, invoice.display_id, invoice.status, invoice.total_minor,
       COALESCE(posted_charges.total_minor, 0) AS posted_payment_total_minor,
       COALESCE(posted_charges.transaction_count, 0) AS posted_payment_transaction_count,
       CASE
         WHEN invoice.status = 'paid' AND COALESCE(posted_charges.transaction_count, 0) = 0 THEN 'paid_invoice_missing_posted_charge'
         WHEN invoice.status = 'paid' AND COALESCE(posted_charges.total_minor, 0) <> invoice.total_minor THEN 'paid_invoice_amount_mismatch'
         WHEN invoice.status <> 'paid' AND COALESCE(posted_charges.transaction_count, 0) > 0 THEN 'unpaid_invoice_has_posted_charge'
         ELSE ''
       END AS reason
FROM invoices invoice
LEFT JOIN posted_charges
  ON posted_charges.invoice_id = invoice.invoice_id
WHERE invoice.tenant_id = $1
  AND (
    (invoice.created_at >= $2 AND invoice.created_at < $3)
    OR (invoice.paid_at >= $2 AND invoice.paid_at < $3)
    OR EXISTS (
        SELECT 1
        FROM payment_transactions txn
        WHERE txn.tenant_id = invoice.tenant_id
          AND txn.invoice_id = invoice.invoice_id
          AND txn.created_at >= $2
          AND txn.created_at < $3
    )
  )
  AND (
    (invoice.status = 'paid' AND COALESCE(posted_charges.transaction_count, 0) = 0)
    OR (invoice.status = 'paid' AND COALESCE(posted_charges.total_minor, 0) <> invoice.total_minor)
    OR (invoice.status <> 'paid' AND COALESCE(posted_charges.transaction_count, 0) > 0)
  )
ORDER BY invoice.display_id`

const dailyPaymentCheckedSQL = `
SELECT COUNT(*)
FROM payment_transactions txn
WHERE txn.tenant_id = $1
  AND txn.created_at >= $2
  AND txn.created_at < $3`

const dailyDuplicatePaymentReferenceSQL = `
WITH changed_references AS (
    SELECT DISTINCT invoice_id
    FROM payment_transactions
    WHERE tenant_id = $1
      AND created_at >= $2
      AND created_at < $3
      AND invoice_id IS NOT NULL
      AND transaction_type = 'charge'
      AND status = 'posted'
)
SELECT 'invoice' AS reference_type,
       txn.invoice_id,
       invoice.display_id,
       STRING_AGG(txn.display_id::text, ',' ORDER BY txn.display_id) AS transaction_display_ids,
       COUNT(*) AS transaction_count,
       COALESCE(SUM(txn.amount_minor), 0) AS total_amount_minor
FROM payment_transactions txn
JOIN changed_references changed
  ON changed.invoice_id = txn.invoice_id
LEFT JOIN invoices invoice
  ON invoice.invoice_id = txn.invoice_id
 AND invoice.tenant_id = txn.tenant_id
WHERE txn.tenant_id = $1
  AND txn.invoice_id IS NOT NULL
  AND txn.transaction_type = 'charge'
  AND txn.status = 'posted'
GROUP BY txn.invoice_id, invoice.display_id
HAVING COUNT(*) > 1
ORDER BY invoice.display_id`

func (store *PostgresStore) GetDailyReconciliationData(ctx context.Context, input DailyReconciliationInput) (DailyReconciliationData, error) {
	if err := store.ready(); err != nil {
		return DailyReconciliationData{}, err
	}
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return DailyReconciliationData{}, err
	}
	windowTo := input.WindowTo()
	walletsChecked, err := store.countDailyReconciliationRows(ctx, dailyWalletSummarySQL, input.TenantID)
	if err != nil {
		return DailyReconciliationData{}, err
	}
	walletMismatches, err := store.listWalletBalanceMismatches(ctx, input)
	if err != nil {
		return DailyReconciliationData{}, err
	}
	invoicesChecked, err := store.countDailyReconciliationRows(ctx, dailyInvoiceCheckedSQL, input.TenantID, input.Date, windowTo)
	if err != nil {
		return DailyReconciliationData{}, err
	}
	invoiceMismatches, err := store.listInvoicePaymentMismatches(ctx, input)
	if err != nil {
		return DailyReconciliationData{}, err
	}
	paymentsChecked, err := store.countDailyReconciliationRows(ctx, dailyPaymentCheckedSQL, input.TenantID, input.Date, windowTo)
	if err != nil {
		return DailyReconciliationData{}, err
	}
	duplicateReferences, err := store.listDuplicatePaymentReferences(ctx, input)
	if err != nil {
		return DailyReconciliationData{}, err
	}
	return DailyReconciliationData{
		WalletsChecked:             walletsChecked,
		WalletMismatches:           walletMismatches,
		InvoicesChecked:            invoicesChecked,
		InvoicePaymentMismatches:   invoiceMismatches,
		PaymentsChecked:            paymentsChecked,
		DuplicatePaymentReferences: duplicateReferences,
	}, nil
}

func (store *PostgresStore) countDailyReconciliationRows(ctx context.Context, query string, args ...interface{}) (int, error) {
	var count int
	if err := store.executor.QueryRowContext(ctx, query, args...).Scan(&count); err != nil {
		return 0, fmt.Errorf("count daily reconciliation rows: %w", err)
	}
	return count, nil
}

func (store *PostgresStore) listWalletBalanceMismatches(ctx context.Context, input DailyReconciliationInput) ([]WalletBalanceMismatch, error) {
	rows, err := store.executor.QueryContext(ctx, dailyWalletMismatchSQL, input.TenantID)
	if err != nil {
		return nil, fmt.Errorf("list wallet balance mismatches: %w", err)
	}
	defer rows.Close()
	mismatches := make([]WalletBalanceMismatch, 0)
	for rows.Next() {
		record, err := scanWalletBalanceMismatch(rows)
		if err != nil {
			return nil, err
		}
		mismatches = append(mismatches, record)
	}
	return mismatches, rows.Err()
}

func (store *PostgresStore) listInvoicePaymentMismatches(ctx context.Context, input DailyReconciliationInput) ([]InvoicePaymentMismatch, error) {
	rows, err := store.executor.QueryContext(ctx, dailyInvoiceMismatchSQL, input.TenantID, input.Date, input.WindowTo())
	if err != nil {
		return nil, fmt.Errorf("list invoice payment mismatches: %w", err)
	}
	defer rows.Close()
	mismatches := make([]InvoicePaymentMismatch, 0)
	for rows.Next() {
		record, err := scanInvoicePaymentMismatch(rows)
		if err != nil {
			return nil, err
		}
		mismatches = append(mismatches, record)
	}
	return mismatches, rows.Err()
}

func (store *PostgresStore) listDuplicatePaymentReferences(ctx context.Context, input DailyReconciliationInput) ([]DuplicatePaymentReference, error) {
	rows, err := store.executor.QueryContext(ctx, dailyDuplicatePaymentReferenceSQL, input.TenantID, input.Date, input.WindowTo())
	if err != nil {
		return nil, fmt.Errorf("list duplicate payment references: %w", err)
	}
	defer rows.Close()
	duplicates := make([]DuplicatePaymentReference, 0)
	for rows.Next() {
		record, err := scanDuplicatePaymentReference(rows)
		if err != nil {
			return nil, err
		}
		duplicates = append(duplicates, record)
	}
	return duplicates, rows.Err()
}

func scanWalletBalanceMismatch(row transactionScanner) (WalletBalanceMismatch, error) {
	var record WalletBalanceMismatch
	var walletID string
	var lastLedgerID sql.NullString
	var lastLedgerDisplayID sql.NullInt64
	if err := row.Scan(
		&walletID, &record.WalletDisplayID, &record.Currency, &record.AvailableBalanceMinor,
		&record.LedgerBalanceMinor, &record.DifferenceMinor, &lastLedgerID, &lastLedgerDisplayID,
	); err != nil {
		return WalletBalanceMismatch{}, fmt.Errorf("scan wallet balance mismatch: %w", err)
	}
	record.WalletID = wallet.WalletID(walletID)
	record.LastLedgerEntryID = wallet.LedgerEntryID(lastLedgerID.String)
	if lastLedgerDisplayID.Valid {
		record.LastLedgerDisplayID = lastLedgerDisplayID.Int64
	}
	return record, nil
}

func scanInvoicePaymentMismatch(row transactionScanner) (InvoicePaymentMismatch, error) {
	var record InvoicePaymentMismatch
	var invoiceID, status string
	if err := row.Scan(
		&invoiceID, &record.InvoiceDisplayID, &status, &record.TotalMinor,
		&record.PostedPaymentTotalMinor, &record.PostedPaymentTransactionCount, &record.Reason,
	); err != nil {
		return InvoicePaymentMismatch{}, fmt.Errorf("scan invoice payment mismatch: %w", err)
	}
	record.InvoiceID = invoice.InvoiceID(invoiceID)
	record.Status = invoice.Status(status)
	return record, nil
}

func scanDuplicatePaymentReference(row transactionScanner) (DuplicatePaymentReference, error) {
	var record DuplicatePaymentReference
	var referenceID, displayIDs string
	if err := row.Scan(
		&record.ReferenceType, &referenceID, &record.ReferenceDisplayID, &displayIDs,
		&record.TransactionCount, &record.TotalAmountMinor,
	); err != nil {
		return DuplicatePaymentReference{}, fmt.Errorf("scan duplicate payment reference: %w", err)
	}
	record.ReferenceID = invoice.InvoiceID(referenceID)
	transactionDisplayIDs, err := parseInt64CSV(displayIDs)
	if err != nil {
		return DuplicatePaymentReference{}, fmt.Errorf("scan duplicate payment reference: %w", err)
	}
	record.TransactionDisplayIDs = transactionDisplayIDs
	return record, nil
}

func parseInt64CSV(value string) ([]int64, error) {
	if strings.TrimSpace(value) == "" {
		return nil, nil
	}
	parts := strings.Split(value, ",")
	values := make([]int64, 0, len(parts))
	for _, part := range parts {
		parsed, err := strconv.ParseInt(strings.TrimSpace(part), 10, 64)
		if err != nil {
			return nil, fmt.Errorf("parse display id %q: %w", strings.TrimSpace(part), err)
		}
		values = append(values, parsed)
	}
	return values, nil
}
