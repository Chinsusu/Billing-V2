CREATE UNIQUE INDEX idx_invoices_tenant_order_unique
    ON invoices(tenant_id, order_id)
    WHERE order_id IS NOT NULL;
