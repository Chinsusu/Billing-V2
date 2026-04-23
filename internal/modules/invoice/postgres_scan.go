package invoice

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Chinsusu/Billing-V2/internal/modules/identity"
	"github.com/Chinsusu/Billing-V2/internal/modules/order"
	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

type invoiceScanner interface {
	Scan(dest ...interface{}) error
}

func scanInvoice(row invoiceScanner) (Invoice, error) {
	var record Invoice
	var id, tenantID, buyerUserID, status string
	var orderID sql.NullString
	var issuedAt, dueAt, paidAt, voidedAt sql.NullTime
	var metadata []byte
	if err := row.Scan(
		&id, &record.DisplayID, &tenantID, &buyerUserID, &orderID, &status, &record.Currency,
		&record.SubtotalMinor, &record.TaxMinor, &record.DiscountMinor, &record.TotalMinor,
		&issuedAt, &dueAt, &paidAt, &voidedAt, &metadata, &record.CreatedAt, &record.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Invoice{}, ErrInvoiceNotFound
		}
		return Invoice{}, fmt.Errorf("scan invoice: %w", err)
	}
	record.ID = InvoiceID(id)
	record.TenantID = tenant.ID(tenantID)
	record.BuyerUserID = identity.UserID(buyerUserID)
	record.OrderID = order.OrderID(orderID.String)
	record.Status = Status(status)
	if issuedAt.Valid {
		record.IssuedAt = issuedAt.Time
	}
	if dueAt.Valid {
		record.DueAt = dueAt.Time
	}
	if paidAt.Valid {
		record.PaidAt = paidAt.Time
	}
	if voidedAt.Valid {
		record.VoidedAt = voidedAt.Time
	}
	record.Metadata = append(record.Metadata, metadata...)
	return record, nil
}

func scanInvoiceItem(row invoiceScanner) (Item, error) {
	var record Item
	var id, invoiceID, tenantID string
	var orderID, orderItemID, serviceID sql.NullString
	var metadata []byte
	if err := row.Scan(
		&id, &invoiceID, &tenantID, &orderID, &orderItemID, &serviceID, &record.Description,
		&record.Quantity, &record.UnitPriceMinor, &record.TaxMinor, &record.DiscountMinor,
		&record.LineTotalMinor, &metadata, &record.CreatedAt, &record.UpdatedAt,
	); err != nil {
		return Item{}, fmt.Errorf("scan invoice item: %w", err)
	}
	record.ID = InvoiceItemID(id)
	record.InvoiceID = InvoiceID(invoiceID)
	record.TenantID = tenant.ID(tenantID)
	record.OrderID = order.OrderID(orderID.String)
	record.OrderItemID = OrderItemID(orderItemID.String)
	record.ServiceID = order.ServiceID(serviceID.String)
	record.Metadata = append(record.Metadata, metadata...)
	return record, nil
}
