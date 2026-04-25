package payment

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/invoice"
	"github.com/Chinsusu/Billing-V2/internal/modules/order"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

type transactionScanner interface {
	Scan(dest ...interface{}) error
}

func scanTransaction(row transactionScanner) (Transaction, error) {
	return scanTransactionFields(row, false)
}

func scanTransactionRead(row transactionScanner) (Transaction, error) {
	return scanTransactionFields(row, true)
}

func scanTransactionFields(row transactionScanner, includeRelatedDisplayIDs bool) (Transaction, error) {
	var record Transaction
	var id, tenantID, accountUserID, transactionType, status, idempotencyKey string
	var orderID, invoiceID, description sql.NullString
	var accountDisplayID, orderDisplayID, invoiceDisplayID sql.NullInt64
	var metadata []byte
	destinations := []interface{}{
		&id, &record.DisplayID, &tenantID, &accountUserID, &orderID, &invoiceID, &transactionType, &status,
		&record.Currency, &record.AmountMinor, &description, &idempotencyKey, &metadata, &record.CreatedAt, &record.UpdatedAt,
	}
	if includeRelatedDisplayIDs {
		destinations = append(destinations, &accountDisplayID, &orderDisplayID, &invoiceDisplayID)
	}
	if err := row.Scan(destinations...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Transaction{}, ErrTransactionNotFound
		}
		return Transaction{}, fmt.Errorf("scan payment transaction: %w", err)
	}
	record.ID = TransactionID(id)
	record.TenantID = tenant.ID(tenantID)
	record.AccountUserID = identity.UserID(accountUserID)
	record.OrderID = order.OrderID(orderID.String)
	record.InvoiceID = invoice.InvoiceID(invoiceID.String)
	if accountDisplayID.Valid {
		record.AccountDisplayID = accountDisplayID.Int64
	}
	if orderDisplayID.Valid {
		record.OrderDisplayID = orderDisplayID.Int64
	}
	if invoiceDisplayID.Valid {
		record.InvoiceDisplayID = invoiceDisplayID.Int64
	}
	record.Type = TransactionType(transactionType)
	record.Status = TransactionStatus(status)
	record.Description = description.String
	record.IdempotencyKey = IdempotencyKey(idempotencyKey)
	record.Metadata = append(json.RawMessage(nil), metadata...)
	return record, nil
}
