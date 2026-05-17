package seed

const seedBillingFlowSQL = `
INSERT INTO wallets (wallet_id, display_id, tenant_id, owner_type, owner_id, currency, status, available_balance_minor, locked_balance_minor, metadata)
SELECT
    '00000000-0000-0000-0000-000000000901',
    41001,
    reseller.tenant_id,
    'user',
    customer.user_id,
    'USD',
    'active',
    3600,
    0,
    '{"seed":"billing_flow","label":"Demo customer wallet"}'::jsonb
FROM tenants reseller
JOIN users customer ON customer.tenant_id = reseller.tenant_id AND customer.email = 'customer@local.billing'
WHERE reseller.slug = 'demo-reseller'
ON CONFLICT (owner_type, owner_id, currency) DO UPDATE
SET status = EXCLUDED.status,
    available_balance_minor = EXCLUDED.available_balance_minor,
    locked_balance_minor = EXCLUDED.locked_balance_minor,
    metadata = EXCLUDED.metadata,
    updated_at = NOW();

INSERT INTO wallet_ledger_entries (ledger_entry_id, display_id, wallet_id, tenant_id, direction, amount_minor, currency, entry_type, status, balance_after_minor, reference_type, reference_id, idempotency_key, created_by, reason, correlation_id)
SELECT
    '00000000-0000-0000-0000-000000000902',
    50001,
    wallet.wallet_id,
    wallet.tenant_id,
    'credit',
    5000,
    'USD',
    'topup',
    'posted',
    5000,
    'topup_request',
    '00000000-0000-0000-0000-000000000908',
    'seed-topup-credit-1',
    customer.user_id,
    'Demo seed top-up',
    '00000000-0000-0000-0000-000000000908'
FROM wallets wallet
JOIN users customer ON customer.user_id = wallet.owner_id
WHERE wallet.wallet_id = '00000000-0000-0000-0000-000000000901'
ON CONFLICT (wallet_id, idempotency_key) DO UPDATE
SET amount_minor = EXCLUDED.amount_minor,
    balance_after_minor = EXCLUDED.balance_after_minor,
    reason = EXCLUDED.reason;

INSERT INTO topup_requests (topup_request_id, display_id, tenant_id, wallet_id, requested_by, amount_minor, currency, payment_method, payment_reference, status, reviewed_by, reviewed_at, review_note, ledger_entry_id, idempotency_key)
SELECT
    '00000000-0000-0000-0000-000000000908',
    52001,
    wallet.tenant_id,
    wallet.wallet_id,
    customer.user_id,
    5000,
    'USD',
    'manual',
    'DEV-TOPUP-0001',
    'approved',
    reviewer.user_id,
    '2026-04-23T00:10:00Z'::timestamptz,
    'Approved by local seed data.',
    '00000000-0000-0000-0000-000000000902',
    'seed-topup-request-1'
FROM wallets wallet
JOIN users customer ON customer.user_id = wallet.owner_id
JOIN users reviewer ON reviewer.email = 'reseller@local.billing'
WHERE wallet.wallet_id = '00000000-0000-0000-0000-000000000901'
ON CONFLICT (tenant_id, requested_by, idempotency_key) DO UPDATE
SET status = EXCLUDED.status,
    reviewed_by = EXCLUDED.reviewed_by,
    reviewed_at = EXCLUDED.reviewed_at,
    review_note = EXCLUDED.review_note,
    ledger_entry_id = EXCLUDED.ledger_entry_id,
    updated_at = NOW();

INSERT INTO orders (order_id, display_id, tenant_id, buyer_user_id, tenant_plan_id, quantity, currency, unit_price_minor, discount_minor, total_minor, order_status, billing_status, idempotency_key, product_snapshot, plan_snapshot, price_snapshot, created_at, updated_at)
SELECT
    '00000000-0000-0000-0000-000000000903',
    42001,
    reseller.tenant_id,
    customer.user_id,
    tenant_plan.tenant_plan_id,
    1,
    tenant_plan.currency,
    tenant_plan.selling_price_minor,
    0,
    tenant_plan.selling_price_minor,
    'paid',
    'paid',
    'seed-order-billing-flow-1',
    tenant_plan.product_snapshot,
    tenant_plan.plan_snapshot,
    tenant_plan.price_snapshot,
    '2026-04-23T00:20:00Z'::timestamptz,
    NOW()
FROM tenants reseller
JOIN users customer ON customer.tenant_id = reseller.tenant_id AND customer.email = 'customer@local.billing'
JOIN tenant_plans tenant_plan ON tenant_plan.tenant_id = reseller.tenant_id
WHERE reseller.slug = 'demo-reseller'
  AND tenant_plan.tenant_plan_id = '00000000-0000-0000-0000-000000000801'
ON CONFLICT (tenant_id, idempotency_key) DO UPDATE
SET order_status = EXCLUDED.order_status,
    billing_status = EXCLUDED.billing_status,
    unit_price_minor = EXCLUDED.unit_price_minor,
    discount_minor = EXCLUDED.discount_minor,
    total_minor = EXCLUDED.total_minor,
    product_snapshot = EXCLUDED.product_snapshot,
    plan_snapshot = EXCLUDED.plan_snapshot,
    price_snapshot = EXCLUDED.price_snapshot,
    updated_at = NOW();

INSERT INTO service_instances (service_instance_id, display_id, tenant_id, order_id, tenant_plan_id, provider_source_id, external_resource_id, status, billing_status, term_start, term_end, created_at, updated_at)
SELECT
    '00000000-0000-0000-0000-000000000909',
    43001,
    orders.tenant_id,
    orders.order_id,
    orders.tenant_plan_id,
    '00000000-0000-0000-0000-000000000301',
    'local-vps-405910',
    'active',
    'paid',
    '2026-04-23T00:30:00Z'::timestamptz,
    '2026-05-23T00:30:00Z'::timestamptz,
    '2026-04-23T00:30:00Z'::timestamptz,
    NOW()
FROM orders
WHERE orders.order_id = '00000000-0000-0000-0000-000000000903'
ON CONFLICT (order_id) DO UPDATE
SET status = EXCLUDED.status,
    billing_status = EXCLUDED.billing_status,
    term_start = EXCLUDED.term_start,
    term_end = EXCLUDED.term_end,
    updated_at = NOW();

INSERT INTO invoices (invoice_id, display_id, tenant_id, buyer_user_id, order_id, status, currency, subtotal_minor, tax_minor, discount_minor, total_minor, issued_at, paid_at, metadata, created_at, updated_at)
SELECT
    '00000000-0000-0000-0000-000000000904',
    44001,
    orders.tenant_id,
    orders.buyer_user_id,
    orders.order_id,
    'paid',
    orders.currency,
    orders.total_minor,
    0,
    0,
    orders.total_minor,
    '2026-04-23T00:35:00Z'::timestamptz,
    '2026-04-23T00:40:00Z'::timestamptz,
    '{"seed":"billing_flow","source":"order"}'::jsonb,
    '2026-04-23T00:35:00Z'::timestamptz,
    NOW()
FROM orders
WHERE orders.order_id = '00000000-0000-0000-0000-000000000903'
ON CONFLICT (tenant_id, order_id) WHERE order_id IS NOT NULL DO UPDATE
SET status = EXCLUDED.status,
    subtotal_minor = EXCLUDED.subtotal_minor,
    tax_minor = EXCLUDED.tax_minor,
    discount_minor = EXCLUDED.discount_minor,
    total_minor = EXCLUDED.total_minor,
    paid_at = EXCLUDED.paid_at,
    metadata = EXCLUDED.metadata,
    updated_at = NOW();

INSERT INTO invoice_items (invoice_item_id, invoice_id, tenant_id, order_id, service_instance_id, description, quantity, unit_price_minor, tax_minor, discount_minor, line_total_minor, metadata)
SELECT
    '00000000-0000-0000-0000-000000000905',
    invoice.invoice_id,
    invoice.tenant_id,
    invoice.order_id,
    service.service_instance_id,
    'CX23 VPS 40GB monthly service',
    1,
    invoice.total_minor,
    0,
    0,
    invoice.total_minor,
    '{"seed":"billing_flow"}'::jsonb
FROM invoices invoice
JOIN service_instances service ON service.order_id = invoice.order_id
WHERE invoice.invoice_id = '00000000-0000-0000-0000-000000000904'
ON CONFLICT (invoice_item_id) DO UPDATE
SET description = EXCLUDED.description,
    unit_price_minor = EXCLUDED.unit_price_minor,
    line_total_minor = EXCLUDED.line_total_minor,
    metadata = EXCLUDED.metadata,
    updated_at = NOW();

INSERT INTO wallet_ledger_entries (ledger_entry_id, display_id, wallet_id, tenant_id, direction, amount_minor, currency, entry_type, status, balance_after_minor, reference_type, reference_id, idempotency_key, created_by, reason, correlation_id)
SELECT
    '00000000-0000-0000-0000-000000000906',
    50002,
    wallet.wallet_id,
    invoice.tenant_id,
    'debit',
    invoice.total_minor,
    invoice.currency,
    'purchase',
    'posted',
    3600,
    'invoice',
    invoice.invoice_id,
    'invoice-payment:00000000-0000-0000-0000-000000000904:seed-payment-1',
    invoice.buyer_user_id,
    'Demo invoice wallet payment',
    invoice.invoice_id
FROM invoices invoice
JOIN wallets wallet ON wallet.owner_id = invoice.buyer_user_id AND wallet.currency = invoice.currency
WHERE invoice.invoice_id = '00000000-0000-0000-0000-000000000904'
ON CONFLICT (wallet_id, idempotency_key) DO UPDATE
SET amount_minor = EXCLUDED.amount_minor,
    balance_after_minor = EXCLUDED.balance_after_minor,
    reason = EXCLUDED.reason;

INSERT INTO payment_transactions (payment_transaction_id, display_id, tenant_id, account_user_id, order_id, invoice_id, transaction_type, status, currency, amount_minor, description, idempotency_key, metadata, created_at, updated_at)
SELECT
    '00000000-0000-0000-0000-000000000907',
    51001,
    invoice.tenant_id,
    invoice.buyer_user_id,
    invoice.order_id,
    invoice.invoice_id,
    'charge',
    'posted',
    invoice.currency,
    invoice.total_minor,
    'Demo invoice wallet payment',
    'seed-payment-1',
    jsonb_build_object(
        'source', 'wallet',
        'provider', 'wallet',
        'wallet_id', '00000000-0000-0000-0000-000000000901',
        'ledger_entry_id', '00000000-0000-0000-0000-000000000906'
    ),
    '2026-04-23T00:45:00Z'::timestamptz,
    NOW()
FROM invoices invoice
WHERE invoice.invoice_id = '00000000-0000-0000-0000-000000000904'
ON CONFLICT (tenant_id, idempotency_key) DO UPDATE
SET status = EXCLUDED.status,
    amount_minor = EXCLUDED.amount_minor,
    description = EXCLUDED.description,
    metadata = EXCLUDED.metadata,
    updated_at = NOW();

INSERT INTO jobs (job_id, display_id, tenant_id, job_type, reference_type, reference_id, source_id, payload_json, status, priority, idempotency_key, attempt_count, max_attempts, next_attempt_at, last_error_code, last_error_message_redacted, manual_review_reason, correlation_id, created_at, updated_at)
SELECT
    '00000000-0000-0000-0000-000000000910',
    53001,
    orders.tenant_id,
    'provider.provision',
    'order',
    orders.order_id,
    '00000000-0000-0000-0000-000000000301',
    '{"seed":"billing_flow","resource":"local-vps-405910"}'::jsonb,
    'manual_review',
    50,
    'seed-provider-provision-1',
    2,
    5,
    '2026-04-23T01:00:00Z'::timestamptz,
    'provider_timeout',
    'Provider timed out in seed smoke.',
    'Verify provider state before retry.',
    '00000000-0000-0000-0000-000000000910',
    '2026-04-23T00:50:00Z'::timestamptz,
    NOW()
FROM orders
WHERE orders.order_id = '00000000-0000-0000-0000-000000000903'
ON CONFLICT (tenant_id, job_type, idempotency_key) WHERE tenant_id IS NOT NULL DO UPDATE
SET status = EXCLUDED.status,
    priority = EXCLUDED.priority,
    attempt_count = EXCLUDED.attempt_count,
    max_attempts = EXCLUDED.max_attempts,
    next_attempt_at = EXCLUDED.next_attempt_at,
    last_error_code = EXCLUDED.last_error_code,
    last_error_message_redacted = EXCLUDED.last_error_message_redacted,
    manual_review_reason = EXCLUDED.manual_review_reason,
    updated_at = NOW();

INSERT INTO audit_logs (audit_id, display_id, tenant_id, actor_id, actor_type, action, target_type, target_id, before_snapshot_redacted, after_snapshot_redacted, metadata_redacted, correlation_id, created_at)
SELECT
    '00000000-0000-0000-0000-000000000911',
    70001,
    job.tenant_id,
    reviewer.user_id,
    'user',
    'job.retry',
    'job',
    job.job_id,
    NULL,
    NULL,
    '{"seed":"billing_flow","job_display_id":53001}'::jsonb,
    job.correlation_id,
    '2026-04-23T00:55:00Z'::timestamptz
FROM jobs job
JOIN users reviewer ON reviewer.email = 'reseller@local.billing'
WHERE job.job_id = '00000000-0000-0000-0000-000000000910'
ON CONFLICT (audit_id) DO UPDATE
SET actor_id = EXCLUDED.actor_id,
    actor_type = EXCLUDED.actor_type,
    action = EXCLUDED.action,
    target_type = EXCLUDED.target_type,
    target_id = EXCLUDED.target_id,
    metadata_redacted = EXCLUDED.metadata_redacted,
    correlation_id = EXCLUDED.correlation_id,
    created_at = EXCLUDED.created_at;
`
