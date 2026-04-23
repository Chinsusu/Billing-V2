package invoice

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	platformdb "github.com/Chinsusu/Billing-V2/internal/platform/db"
)

type PostgresStore struct {
	executor platformdb.Executor
}

var _ Store = (*PostgresStore)(nil)

func NewPostgresStore(executor platformdb.Executor) *PostgresStore {
	return &PostgresStore{executor: executor}
}

const invoiceReadColumns = `inv.invoice_id, inv.display_id, inv.tenant_id, inv.buyer_user_id, inv.order_id, inv.status, inv.currency, inv.subtotal_minor, inv.tax_minor, inv.discount_minor, inv.total_minor, inv.issued_at, inv.due_at, inv.paid_at, inv.voided_at, inv.metadata, inv.created_at, inv.updated_at`
const invoiceItemReadColumns = `item.invoice_item_id, item.invoice_id, item.tenant_id, item.order_id, item.order_item_id, item.service_instance_id, item.description, item.quantity, item.unit_price_minor, item.tax_minor, item.discount_minor, item.line_total_minor, item.metadata, item.created_at, item.updated_at`
const invoiceColumns = `invoice_id, display_id, tenant_id, buyer_user_id, order_id, status, currency, subtotal_minor, tax_minor, discount_minor, total_minor, issued_at, due_at, paid_at, voided_at, metadata, created_at, updated_at`

const createInvoiceFromOrderSQL = `
WITH inserted AS (
INSERT INTO invoices (tenant_id, buyer_user_id, order_id, status, currency, subtotal_minor, tax_minor, discount_minor, total_minor, issued_at, metadata)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11::jsonb)
ON CONFLICT (tenant_id, order_id) WHERE order_id IS NOT NULL
DO NOTHING
RETURNING ` + invoiceColumns + `
), selected AS (
SELECT ` + invoiceColumns + ` FROM inserted
UNION ALL
SELECT ` + invoiceColumns + `
FROM invoices
WHERE tenant_id = $1 AND order_id = $3
LIMIT 1
), inserted_item AS (
INSERT INTO invoice_items (invoice_id, tenant_id, order_id, order_item_id, service_instance_id, description, quantity, unit_price_minor, tax_minor, discount_minor, line_total_minor, metadata)
SELECT invoice_id, tenant_id, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21::jsonb
FROM selected
WHERE EXISTS (SELECT 1 FROM inserted)
RETURNING invoice_item_id
), event AS (
INSERT INTO outbox_events (tenant_id, aggregate_type, aggregate_id, event_type, payload_json, dedupe_key, correlation_id)
SELECT
    tenant_id,
    '` + AggregateTypeInvoice + `',
    invoice_id,
    '` + EventInvoiceGenerated + `',
    jsonb_build_object(
        'invoice_id', invoice_id,
        'display_id', display_id,
        'tenant_id', tenant_id,
        'order_id', order_id,
        'order_display_id', $22,
        'buyer_user_id', buyer_user_id,
        'total_minor', total_minor,
        'currency', currency,
        'idempotency_key', $23
    ),
    '` + EventInvoiceGenerated + `:' || order_id::text,
    invoice_id
FROM selected
WHERE EXISTS (SELECT 1 FROM inserted)
ON CONFLICT (dedupe_key) DO NOTHING
)
SELECT ` + invoiceColumns + ` FROM selected`

func (store *PostgresStore) CreateInvoiceFromOrder(ctx context.Context, input CreateInvoiceFromOrderInput) (InvoiceDetail, error) {
	if err := store.ready(); err != nil {
		return InvoiceDetail{}, err
	}
	args, err := createInvoiceFromOrderArgs(input)
	if err != nil {
		return InvoiceDetail{}, err
	}
	record, err := scanInvoice(store.executor.QueryRowContext(ctx, createInvoiceFromOrderSQL, args...))
	if err != nil {
		return InvoiceDetail{}, err
	}
	items, err := store.listInvoiceItems(ctx, record.TenantID, record.ID)
	if err != nil {
		return InvoiceDetail{}, err
	}
	return InvoiceDetail{Invoice: record, Items: items}, nil
}

const markInvoicePaidSQL = `
WITH updated AS (
UPDATE invoices
SET status = 'paid',
    paid_at = COALESCE(paid_at, $3, NOW()),
    updated_at = NOW()
WHERE invoice_id = $1
  AND tenant_id = $2
  AND status IN ('issued', 'overdue')
RETURNING ` + invoiceColumns + `
), event AS (
INSERT INTO outbox_events (tenant_id, aggregate_type, aggregate_id, event_type, payload_json, dedupe_key, correlation_id)
SELECT
    tenant_id,
    '` + AggregateTypeInvoice + `',
    invoice_id,
    '` + EventInvoicePaid + `',
    jsonb_build_object(
        'invoice_id', invoice_id,
        'display_id', display_id,
        'tenant_id', tenant_id,
        'order_id', order_id,
        'buyer_user_id', buyer_user_id,
        'total_minor', total_minor,
        'currency', currency,
        'payment_transaction_id', $4,
        'wallet_id', $5,
        'ledger_entry_id', $6,
        'idempotency_key', $7
    ),
    '` + EventInvoicePaid + `:' || invoice_id::text,
    invoice_id
FROM updated
ON CONFLICT (dedupe_key) DO NOTHING
)
SELECT ` + invoiceColumns + ` FROM updated`

func (store *PostgresStore) MarkInvoicePaid(ctx context.Context, input MarkInvoicePaidInput) (InvoiceDetail, error) {
	if err := store.ready(); err != nil {
		return InvoiceDetail{}, err
	}
	args, err := markInvoicePaidArgs(input)
	if err != nil {
		return InvoiceDetail{}, err
	}
	record, err := scanInvoice(store.executor.QueryRowContext(ctx, markInvoicePaidSQL, args...))
	if errors.Is(err, ErrInvoiceNotFound) {
		if _, lookupErr := store.GetInvoice(ctx, InvoiceLookup{ID: input.ID, TenantID: input.TenantID}); lookupErr == nil {
			return InvoiceDetail{}, ErrInvoiceStatusConflict
		}
	}
	if err != nil {
		return InvoiceDetail{}, err
	}
	items, err := store.listInvoiceItems(ctx, record.TenantID, record.ID)
	if err != nil {
		return InvoiceDetail{}, err
	}
	return InvoiceDetail{Invoice: record, Items: items}, nil
}

func markInvoicePaidArgs(input MarkInvoicePaidInput) ([]interface{}, error) {
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return nil, err
	}
	return []interface{}{
		input.ID,
		input.TenantID,
		nullableTime(input.PaidAt),
		nullableString(input.PaymentTransactionID),
		nullableString(input.WalletID),
		nullableString(input.LedgerEntryID),
		input.IdempotencyKey,
	}, nil
}

func createInvoiceFromOrderArgs(input CreateInvoiceFromOrderInput) ([]interface{}, error) {
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return nil, err
	}
	return []interface{}{
		input.Invoice.TenantID,
		input.Invoice.BuyerUserID,
		input.Invoice.OrderID,
		input.Invoice.Status,
		input.Invoice.Currency,
		input.Invoice.SubtotalMinor,
		input.Invoice.TaxMinor,
		input.Invoice.DiscountMinor,
		input.Invoice.TotalMinor,
		input.Invoice.IssuedAt,
		json.RawMessage(input.Invoice.Metadata),
		input.Item.OrderID,
		nullableString(string(input.Item.OrderItemID)),
		nullableString(string(input.Item.ServiceID)),
		input.Item.Description,
		input.Item.Quantity,
		input.Item.UnitPriceMinor,
		input.Item.TaxMinor,
		input.Item.DiscountMinor,
		input.Item.LineTotalMinor,
		json.RawMessage(input.Item.Metadata),
		input.OrderDisplayID,
		input.IdempotencyKey,
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

func nullableTime(value time.Time) sql.NullTime {
	if value.IsZero() {
		return sql.NullTime{}
	}
	return sql.NullTime{Time: value, Valid: true}
}
