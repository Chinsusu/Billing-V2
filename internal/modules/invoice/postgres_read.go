package invoice

import (
	"context"
	"fmt"

	"github.com/Chinsusu/Billing-V2/internal/modules/tenant"
)

func (store *PostgresStore) ListInvoices(ctx context.Context, filter InvoiceFilter) ([]Invoice, error) {
	if err := store.ready(); err != nil {
		return nil, err
	}
	query, args, err := buildListInvoicesQuery(filter)
	if err != nil {
		return nil, err
	}
	rows, err := store.executor.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list invoices: %w", err)
	}
	defer rows.Close()
	invoices := make([]Invoice, 0)
	for rows.Next() {
		invoice, err := scanInvoiceRead(rows)
		if err != nil {
			return nil, err
		}
		invoices = append(invoices, invoice)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("read invoices: %w", err)
	}
	return invoices, nil
}

func (store *PostgresStore) GetInvoice(ctx context.Context, lookup InvoiceLookup) (InvoiceDetail, error) {
	if err := store.ready(); err != nil {
		return InvoiceDetail{}, err
	}
	query, args, err := buildGetInvoiceQuery(lookup)
	if err != nil {
		return InvoiceDetail{}, err
	}
	record, err := scanInvoiceRead(store.executor.QueryRowContext(ctx, query, args...))
	if err != nil {
		return InvoiceDetail{}, err
	}
	items, err := store.listInvoiceItems(ctx, record.TenantID, record.ID)
	if err != nil {
		return InvoiceDetail{}, err
	}
	return InvoiceDetail{Invoice: record, Items: items}, nil
}

func (store *PostgresStore) listInvoiceItems(ctx context.Context, tenantID tenant.ID, invoiceID InvoiceID) ([]Item, error) {
	rows, err := store.executor.QueryContext(ctx, `SELECT `+invoiceItemReadColumns+`
FROM invoice_items item
WHERE item.tenant_id = $1
  AND item.invoice_id = $2
ORDER BY item.created_at ASC`, tenantID, invoiceID)
	if err != nil {
		return nil, fmt.Errorf("list invoice items: %w", err)
	}
	defer rows.Close()
	items := make([]Item, 0)
	for rows.Next() {
		item, err := scanInvoiceItem(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("read invoice items: %w", err)
	}
	return items, nil
}

func buildListInvoicesQuery(filter InvoiceFilter) (string, []interface{}, error) {
	filter = normalizeInvoiceFilter(filter)
	if err := validateInvoiceFilter(filter); err != nil {
		return "", nil, err
	}
	query := `SELECT ` + invoiceReadColumns + `
FROM invoices inv
WHERE inv.tenant_id = $1`
	args := []interface{}{filter.TenantID}
	if filter.BuyerUserID != "" {
		args = append(args, filter.BuyerUserID)
		query += fmt.Sprintf("\n  AND inv.buyer_user_id = $%d", len(args))
	}
	if filter.BuyerDisplayID > 0 {
		args = append(args, filter.BuyerDisplayID)
		query += fmt.Sprintf(`
  AND EXISTS (
    SELECT 1
    FROM users buyer
    WHERE buyer.user_id = inv.buyer_user_id
      AND buyer.tenant_id = inv.tenant_id
      AND buyer.display_id = $%d
  )`, len(args))
	}
	if filter.DisplayID > 0 {
		args = append(args, filter.DisplayID)
		query += fmt.Sprintf("\n  AND inv.display_id = $%d", len(args))
	}
	if filter.OrderID != "" {
		args = append(args, filter.OrderID)
		query += fmt.Sprintf("\n  AND inv.order_id = $%d", len(args))
	}
	if filter.OrderDisplayID > 0 {
		args = append(args, filter.OrderDisplayID)
		query += fmt.Sprintf(`
  AND EXISTS (
    SELECT 1
    FROM orders ord
    WHERE ord.order_id = inv.order_id
      AND ord.tenant_id = inv.tenant_id
      AND ord.display_id = $%d
  )`, len(args))
	}
	if filter.Status != "" {
		args = append(args, filter.Status)
		query += fmt.Sprintf("\n  AND inv.status = $%d", len(args))
	}
	if filter.AmountMinMinor != nil {
		args = append(args, *filter.AmountMinMinor)
		query += fmt.Sprintf("\n  AND inv.total_minor >= $%d", len(args))
	}
	if filter.AmountMaxMinor != nil {
		args = append(args, *filter.AmountMaxMinor)
		query += fmt.Sprintf("\n  AND inv.total_minor <= $%d", len(args))
	}
	args = append(args, filter.Limit)
	query += fmt.Sprintf("\nORDER BY inv.created_at DESC\nLIMIT $%d", len(args))
	return query, args, nil
}

func buildGetInvoiceQuery(lookup InvoiceLookup) (string, []interface{}, error) {
	if err := validateInvoiceLookup(lookup); err != nil {
		return "", nil, err
	}
	query := `SELECT ` + invoiceReadColumns + `
FROM invoices inv
WHERE inv.invoice_id = $1
  AND inv.tenant_id = $2`
	args := []interface{}{lookup.ID, lookup.TenantID}
	if lookup.BuyerUserID != "" {
		args = append(args, lookup.BuyerUserID)
		query += fmt.Sprintf("\n  AND inv.buyer_user_id = $%d", len(args))
	}
	return query, args, nil
}
