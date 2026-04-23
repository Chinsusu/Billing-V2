package invoice

import (
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

func (store *PostgresStore) ready() error {
	if store == nil || store.executor == nil {
		return ErrStoreExecutorMissing
	}
	return nil
}
