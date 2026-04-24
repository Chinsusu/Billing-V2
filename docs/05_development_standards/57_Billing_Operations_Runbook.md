# 57 - Billing Operations Runbook

## 1. Purpose

Use this runbook when a billing flow needs to be checked from order creation through invoice payment and provisioning follow-up.

It is written for local/dev and sandbox investigation. Do not run these commands against production unless the production runbook and approval process explicitly allow it.

Related references:

- `docs/05_development_standards/55_Local_Development_Runbook.md`
- `docs/05_development_standards/56_Billing_API_Operational_Reference.md`
- `docs/05_development_standards/58_Provisioning_Ops_Readiness_Checklist.md`
- `cmd/smoke dev-billing`

## 2. Required Context

Set these values for local/dev checks:

```bash
export API_BASE_URL="http://localhost:8080"
export DB_DSN="postgres://billing:billing@localhost:5432/billing?sslmode=disable"
export TENANT_ID="00000000-0000-0000-0000-000000000010"
export ADMIN_ID="00000000-0000-0000-0000-000000000102"
export CUSTOMER_ID="00000000-0000-0000-0000-000000000103"
```

Client headers:

```bash
CLIENT_HEADERS=(
  -H "X-Tenant-Id: $TENANT_ID"
  -H "X-Actor-Id: $CUSTOMER_ID"
  -H "X-Actor-Type: client"
  -H "X-Actor-Tenant-Id: $TENANT_ID"
)
```

Admin/reseller headers:

```bash
ADMIN_HEADERS=(
  -H "X-Tenant-Id: $TENANT_ID"
  -H "X-Actor-Id: $ADMIN_ID"
  -H "X-Actor-Type: reseller_owner"
  -H "X-Actor-Tenant-Id: $TENANT_ID"
)
```

Run the deterministic smoke before deeper checks. It creates and pays a new order, runs the fake-provider provisioning worker against the same database, and verifies that the service is visible.

```bash
go run ./cmd/smoke -dsn "$DB_DSN" -base-url "$API_BASE_URL" dev-billing
```

The smoke must use the same database as the running API. It creates a top-up, creates an order, checks out an invoice, pays it from the wallet, verifies the paid order, processes the `provider.provision` job with the fake worker, and checks the resulting service.

## 3. Normal Flow

Expected state sequence:

| Step | Owner | Expected state |
|---|---|---|
| Client creates order | `POST /client/orders` | `order_status=pending_payment`, `billing_status=unpaid` |
| Client checks out order | `POST /client/checkouts` | invoice `status=issued` |
| Client pays invoice from wallet | `POST /client/invoice-wallet-payments` | invoice `paid`, transaction `posted`, ledger debit `posted` |
| Payment finalizes order | payment service | `order_status=paid`, `billing_status=paid` |
| Provisioning is queued | order/payment service | one `jobs` row with `job_type=provider.provision` |
| Worker provisions service | worker/provider flow | service appears under `GET /client/services?order_id=<order_id>` |

Important rule: payment success means the order is paid. It does not mean the provider has already created the service.

## 4. Quick Inspection

Order:

```bash
curl -s "${CLIENT_HEADERS[@]}" "$API_BASE_URL/client/orders/$ORDER_ID"
```

Invoices for an order:

```bash
curl -s "${CLIENT_HEADERS[@]}" "$API_BASE_URL/client/invoices?order_id=$ORDER_ID&limit=20"
```

Payment transactions for an order:

```bash
curl -s "${CLIENT_HEADERS[@]}" "$API_BASE_URL/client/transactions?order_id=$ORDER_ID&limit=20"
```

Services for an order:

```bash
curl -s "${CLIENT_HEADERS[@]}" "$API_BASE_URL/client/services?order_id=$ORDER_ID&limit=20"
```

Admin payment reconciliation:

```bash
curl -s "${ADMIN_HEADERS[@]}" "$API_BASE_URL/admin/payment-reconciliation?invoice_id=$INVOICE_ID&limit=20"
```

Audit logs:

```bash
curl -s "${ADMIN_HEADERS[@]}" "$API_BASE_URL/admin/audit-logs?target_id=$INVOICE_ID&limit=20"
```

Provisioning job:

```bash
curl -s "${ADMIN_HEADERS[@]}" "$API_BASE_URL/admin/jobs?job_type=provider.provision&reference_type=order&reference_id=$ORDER_ID&limit=20"
```

Provisioning job attempts after selecting a `JOB_ID` from the job response:

```bash
curl -s "${ADMIN_HEADERS[@]}" "$API_BASE_URL/admin/jobs/$JOB_ID/attempts?limit=20"
```

Run one local provisioning worker pass with the fake provider registry:

```bash
go run ./cmd/worker provision-once -dsn "$DB_DSN"
```

Run the local/sandbox worker as a loop when you want it to keep polling claimable jobs. The loop prints one summary per pass and waits on `-interval` after an idle pass, so it does not spin when there is no work.

```bash
go run ./cmd/worker provision-loop -dsn "$DB_DSN" -interval 5s -batch-size 10
```

Stop the loop with `Ctrl+C`, or use `-timeout` for a bounded sandbox run. Keep `APP_ENV` set to `local`, `dev`, or another non-production value.

Retry a retryable or manual-review provisioning job after fixing the provider/source issue:

```bash
curl -s -X POST "${ADMIN_HEADERS[@]}" \
  -H "Content-Type: application/json" \
  "$API_BASE_URL/admin/jobs/$JOB_ID/retry" \
  -d '{"next_attempt_at":"2026-04-24T00:00:00Z"}'
```

Move a safe stuck job to manual review with the operator reason:

```bash
curl -s -X POST "${ADMIN_HEADERS[@]}" \
  -H "Content-Type: application/json" \
  "$API_BASE_URL/admin/jobs/$JOB_ID/manual-review" \
  -d '{"reason":"provider response needs manual verification"}'
```

Cancel a safe non-active job when the order should not provision:

```bash
curl -s -X POST "${ADMIN_HEADERS[@]}" \
  -H "Content-Type: application/json" \
  "$API_BASE_URL/admin/jobs/$JOB_ID/cancel" \
  -d '{"reason":"order was voided before provisioning"}'
```

If the API is unavailable, inspect the database:

```bash
psql "$DB_DSN" -c "
SELECT display_id, job_type, reference_type, reference_id, source_id, status,
       attempt_count, max_attempts, last_error_code, last_error_message_redacted,
       manual_review_reason, created_at, updated_at, finished_at
FROM jobs
WHERE job_type = 'provider.provision'
  AND reference_type = 'order'
  AND reference_id = '$ORDER_ID'
ORDER BY created_at DESC;"
```

Service record:

```bash
psql "$DB_DSN" -c "
SELECT display_id, order_id, tenant_plan_id, provider_source_id, external_resource_id,
       status, billing_status, term_start, term_end, updated_at
FROM service_instances
WHERE order_id = '$ORDER_ID'
ORDER BY created_at DESC;"
```

## 5. Common Failures

### 5.1 Insufficient Balance

API code: `wallet.insufficient_balance`

Meaning: the invoice is payable, but the customer wallet cannot cover the invoice total.

Inspect:

```bash
curl -s "${CLIENT_HEADERS[@]}" "$API_BASE_URL/client/wallets"
curl -s "${CLIENT_HEADERS[@]}" "$API_BASE_URL/client/invoices/$INVOICE_ID"
curl -s "${CLIENT_HEADERS[@]}" "$API_BASE_URL/client/transactions?invoice_id=$INVOICE_ID&limit=20"
```

Recovery:

- Do not mark the invoice paid by hand.
- Ask the customer to top up, or create and approve a local/dev top-up when testing.
- Retry `POST /client/invoice-wallet-payments` with a new `Idempotency-Key` after balance is available.

### 5.2 Duplicate Submit

Related routes: `POST /client/orders`, `POST /client/checkouts`, `POST /client/invoice-wallet-payments`

Meaning: the same action was submitted more than once. The expected result is reuse of the same created record when the idempotency key and payload match.

Inspect:

```bash
curl -s "${CLIENT_HEADERS[@]}" "$API_BASE_URL/client/invoices?order_id=$ORDER_ID&limit=20"
curl -s "${CLIENT_HEADERS[@]}" "$API_BASE_URL/client/transactions?invoice_id=$INVOICE_ID&limit=20"
```

DB check:

```bash
psql "$DB_DSN" -c "
SELECT display_id, order_id, status, idempotency_key, created_at, updated_at
FROM invoices
WHERE order_id = '$ORDER_ID'
ORDER BY created_at DESC;"
```

Recovery:

- If the first request succeeded, use the existing order, invoice, transaction, or job.
- If the API returns `payment.idempotency_conflict`, the same key was reused with different payment details. Stop and create a new request key only after confirming no duplicate charge exists.

### 5.3 Checkout Conflict

API code: `checkout.order_not_checkoutable`

Meaning: the order is not in `pending_payment/unpaid`, or it is otherwise not ready to become an invoice.

Inspect:

```bash
curl -s "${CLIENT_HEADERS[@]}" "$API_BASE_URL/client/orders/$ORDER_ID"
curl -s "${CLIENT_HEADERS[@]}" "$API_BASE_URL/client/invoices?order_id=$ORDER_ID&limit=20"
```

Recovery:

- If an invoice already exists, continue with that invoice instead of creating another one.
- If the order is already `paid/paid`, do not run checkout again. Inspect provisioning instead.
- If the order was cancelled, failed, or refunded, create a new order for the customer flow.

### 5.4 Invoice Already Paid

API code: `payment.invoice_not_payable`

Meaning: the invoice cannot be paid again. A paid invoice with a different payment key should not create another wallet debit.

Inspect:

```bash
curl -s "${CLIENT_HEADERS[@]}" "$API_BASE_URL/client/invoices/$INVOICE_ID"
curl -s "${CLIENT_HEADERS[@]}" "$API_BASE_URL/client/transactions?invoice_id=$INVOICE_ID&limit=20"
curl -s "${ADMIN_HEADERS[@]}" "$API_BASE_URL/admin/payment-reconciliation?invoice_id=$INVOICE_ID&limit=20"
```

DB check:

```bash
psql "$DB_DSN" -c "
SELECT inv.display_id AS invoice_display_id, inv.status AS invoice_status,
       tx.display_id AS transaction_display_id, tx.status AS transaction_status,
       tx.idempotency_key, tx.amount_minor
FROM invoices inv
LEFT JOIN payment_transactions tx ON tx.invoice_id = inv.invoice_id
WHERE inv.invoice_id = '$INVOICE_ID'
ORDER BY tx.created_at DESC;"
```

Recovery:

- Use the existing paid invoice and posted transaction.
- Do not create a second transaction manually.
- If the UI lost the response, read invoice, transaction, ledger, and order status from the API paths above.

### 5.5 Provisioning Stuck

Symptoms:

- order is `paid/paid`,
- invoice is `paid`,
- no service is visible for the order,
- `jobs` has `provider.provision` stuck in `queued`, `claimed`, `running`, `failed_retryable`, `failed_terminal`, or `manual_review`.

Inspect:

```bash
curl -s "${CLIENT_HEADERS[@]}" "$API_BASE_URL/client/orders/$ORDER_ID"
curl -s "${CLIENT_HEADERS[@]}" "$API_BASE_URL/client/services?order_id=$ORDER_ID&limit=20"
curl -s "${ADMIN_HEADERS[@]}" "$API_BASE_URL/admin/jobs/$JOB_ID/attempts?limit=20"
psql "$DB_DSN" -c "
SELECT display_id, source_id, status, attempt_count, max_attempts,
       last_error_code, last_error_message_redacted, manual_review_reason, updated_at
FROM jobs
WHERE job_type = 'provider.provision'
  AND reference_type = 'order'
  AND reference_id = '$ORDER_ID'
ORDER BY created_at DESC;"
```

Check that the order has an active provider source:

```bash
psql "$DB_DSN" -c "
SELECT tp.display_id AS tenant_plan_display_id, ps.source_id AS provider_source_id,
       src.display_id AS provider_source_display_id, src.status AS provider_source_status,
       ps.status AS plan_source_status
FROM orders ord
JOIN tenant_plans tp ON tp.tenant_plan_id = ord.tenant_plan_id
JOIN plan_sources ps ON ps.plan_id = tp.master_plan_id
JOIN provider_sources src ON src.source_id = ps.source_id
WHERE ord.order_id = '$ORDER_ID';"
```

Recovery:

- If the job is `queued`, wait for the worker or start the local/sandbox worker.
- In local/dev, run `go run ./cmd/worker provision-once -dsn "$DB_DSN"` to process claimable `provider.provision` jobs once, or `go run ./cmd/worker provision-loop -dsn "$DB_DSN" -interval 5s` while testing repeated fulfillment.
- If the job is `failed_retryable`, inspect `last_error_code` and provider source config before retrying with `POST /admin/jobs/$JOB_ID/retry`.
- If the job is `manual_review`, record the reason, confirm provider state, then retry or cancel through the admin job action routes.
- If the job is `failed_terminal`, keep the order paid and hand it to operations. Move it to manual review or cancel only after confirming provider state.
- If no `provider.provision` job exists for a paid order, treat it as a backend defect and create a fix task. Do not pay the invoice again.

## 6. Rollback And Recovery Rules

- Money rows are append-only in normal operation. Do not edit wallet ledger, payment transaction, or invoice rows by hand.
- Use API reads first, then DB reads only for job and cross-table investigation.
- Use `display_id` when talking to support or operators, and UUID ids for API paths.
- A paid invoice and posted transaction must match the same invoice, wallet, currency, and amount.
- A paid order must have either a provisioning job or a service record. If it has neither, open a backend task.
- Provider success must be checked through service records or provider/job logs. Do not assume the service exists because payment passed.

| Area | Recovery approach | Do not do |
|---|---|---|
| Order before payment | Create a new order if the old one is cancelled, failed, or has the wrong plan/price. | Do not force a bad order into `paid`. |
| Invoice before payment | Reuse the existing issued invoice for the same order, or create a new order if the order is no longer checkoutable. | Do not create duplicate invoices for the same live order by hand. |
| Wallet balance | Add funds through top-up flow in local/dev, then retry payment with a new key. | Do not edit `wallets.available_balance_minor` directly. |
| Payment transaction | Use existing posted transaction when the invoice is already paid. | Do not insert a second charge for the same invoice. |
| Provisioning job | Fix provider source/config first, then retry, manual-review, or cancel through the approved admin action routes. | Do not pay the invoice again to create another job. |
| Service record | If provider created the service but no local service exists, open a backend repair task with provider evidence. | Do not invent provider resource ids without proof. |

## 7. Local Validation Commands

For code or runbook changes around this flow, run:

```bash
go test ./...
go build ./cmd/api ./cmd/migrate ./cmd/seed ./cmd/smoke ./cmd/worker
cd frontend
npm audit --omit=dev
npm run lint
npm run build
```

For end-to-end local verification with a running API:

```bash
go run ./cmd/smoke -dsn "$DB_DSN" dev-db
go run ./cmd/smoke -base-url "$API_BASE_URL" dev-api
go run ./cmd/smoke -dsn "$DB_DSN" -base-url "$API_BASE_URL" dev-billing
```
