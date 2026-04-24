# 56 - Billing API Operational Reference

## 1. Purpose

This document is the working contract for the billing API that the current Go backend exposes.

Use it when:
- frontend screens need stable request and response shapes,
- coding agents need the exact route and filter names,
- reviewers need to confirm that a change matches the live API surface.

This is not generated OpenAPI. It is a compact operational reference built from the current route, handler, filter, and response code.

## 2. Shared Rules

### 2.1 Path IDs vs display IDs

- Path parameters use UUID-style ids such as `order_id`, `wallet_id`, `invoice_id`, `transaction_id`, and `audit_log_id`.
- Human-facing tables should use `display_id` from the response body.
- `display_id` is also available as a list filter on most billing read APIs.

### 2.2 Headers

Required on tenant-scoped billing routes:

| Header | Required | Notes |
|---|---|---|
| `X-Tenant-Id` | yes | Sets the effective tenant context. |
| `X-Actor-Id` | yes on authenticated routes | Actor id used for client or admin access checks. |
| `X-Actor-Type` | yes when `X-Actor-Id` is present | Actor type such as client or reseller owner. |
| `X-Actor-Tenant-Id` | recommended | Keeps actor tenant explicit. |
| `X-Request-ID` | optional | Echoed back in the response header and body. |
| `Idempotency-Key` | required on create money-changing client routes | Used on order create, checkout create, top-up create, and invoice wallet payment. |

Notes:
- Route authorization is currently backed by the database role and permission store, not by `X-Actor-Role-Ids`.
- If `X-Actor-Id` is missing, the API returns `auth.actor_required`.

### 2.3 Response envelope

Success:

```json
{
  "data": {},
  "request_id": "req_..."
}
```

List success:

```json
{
  "data": [],
  "page": {
    "limit": 20,
    "next_cursor": null
  },
  "request_id": "req_..."
}
```

Validation error:

```json
{
  "error": {
    "code": "validation.failed",
    "message": "Request validation failed.",
    "fields": [
      {
        "field": "display_id",
        "code": "request.display_id_invalid",
        "message": "Display id must be a positive number."
      }
    ]
  },
  "request_id": "req_..."
}
```

Operation error:

```json
{
  "error": {
    "code": "wallet.insufficient_balance",
    "message": "Wallet balance is insufficient."
  },
  "request_id": "req_..."
}
```

### 2.4 Pagination and range filters

- All list routes use `limit` and optional `cursor`.
- Default `limit` is `20`.
- Max `limit` is `100`.
- Money range filters use minor units and the pair `amount_min` and `amount_max`.
- When both are present, `amount_max` must be greater than or equal to `amount_min`.
- Time range filters use RFC3339 values.

### 2.5 Permission map

| Route group | Permission |
|---|---|
| Admin catalog and provider readiness | `catalog.manage` |
| Client orders | `order.create` |
| Client checkout | `order.create` |
| Admin order read | `order.view` |
| Admin order status mutation | `order.manage` |
| Admin and reseller job read | `order.view` |
| Admin job retry | `provisioning.job.retry` |
| Admin job manual review or cancel | `provisioning.manual_review.resolve` |
| Client and admin services | `service.view` |
| Client and admin invoices | `wallet.view` |
| Client and admin wallets | `wallet.view` |
| Client and admin transactions | `wallet.view` |
| Client top-up create and read | `wallet.view` |
| Admin top-up read | `wallet.view` |
| Admin top-up approve or reject | `wallet.topup.approve` |
| Admin payment reconciliation | `wallet.view` |
| Admin audit logs | `audit.view` |

## 3. Resource Shapes

### 3.1 Order

`order` response fields:

`id`, `display_id`, `tenant_id`, `buyer_user_id`, `tenant_plan_id`, `quantity`, `currency`, `unit_price_minor`, `discount_minor`, `total_minor`, `order_status`, `billing_status`, `product_snapshot`, `plan_snapshot`, `price_snapshot`, `created_at`, `updated_at`

### 3.2 Service instance

`service` response fields:

`id`, `display_id`, `tenant_id`, `order_id`, `tenant_plan_id`, `provider_source_id`, `external_resource_id`, `status`, `billing_status`, `suspension_reason`, `term_start`, `term_end`, `created_at`, `updated_at`

### 3.3 Invoice

`invoice` list fields:

`id`, `display_id`, `tenant_id`, `buyer_user_id`, `order_id`, `status`, `currency`, `subtotal_minor`, `tax_minor`, `discount_minor`, `total_minor`, `issued_at`, `due_at`, `paid_at`, `voided_at`, `metadata`, `created_at`, `updated_at`

Invoice detail adds:

`items[]` with `id`, `invoice_id`, `tenant_id`, `order_id`, `order_item_id`, `service_id`, `description`, `quantity`, `unit_price_minor`, `tax_minor`, `discount_minor`, `line_total_minor`, `metadata`, `created_at`, `updated_at`

### 3.4 Wallet and ledger

`wallet` response fields:

`id`, `display_id`, `tenant_id`, `owner_type`, `owner_id`, `currency`, `status`, `available_balance_minor`, `locked_balance_minor`, `metadata`, `created_at`, `updated_at`

`ledger` response fields:

`id`, `display_id`, `wallet_id`, `tenant_id`, `direction`, `amount_minor`, `currency`, `entry_type`, `status`, `balance_after_minor`, `reference_type`, `reference_id`, `created_by`, `reason`, `correlation_id`, `created_at`

### 3.5 Top-up request

`topup_request` response fields:

`id`, `display_id`, `tenant_id`, `wallet_id`, `requested_by`, `amount_minor`, `currency`, `payment_method`, `payment_reference`, `status`, `reviewed_by`, `reviewed_at`, `review_note`, `ledger_entry_id`, `created_at`, `updated_at`

### 3.6 Transaction and reconciliation

`transaction` response fields:

`id`, `display_id`, `tenant_id`, `account_user_id`, `order_id`, `invoice_id`, `type`, `status`, `currency`, `amount_minor`, `description`, `metadata`, `created_at`, `updated_at`

`payment_reconciliation` response fields:

- `transaction`
- `provider`
- optional `invoice` with `id`, `display_id`, `status`, `total_minor`, `paid_at`
- optional `ledger` with `id`, `display_id`, `wallet_id`, `wallet_display_id`, `direction`, `entry_type`, `status`, `balance_after_minor`

### 3.7 Audit log

Audit list item fields:

`id`, `display_id`, `tenant_id`, `actor_id`, `actor_type`, `action`, `target_type`, `target_id`, `ip_address`, `correlation_id`, `created_at`

Audit detail adds:

`before_snapshot_redacted`, `after_snapshot_redacted`, `metadata_redacted`, `user_agent`

### 3.8 Job

`job` response fields:

`id`, `display_id`, `tenant_id`, `job_type`, `reference_type`, `reference_id`, `source_id`, `status`, `priority`, `attempt_count`, `max_attempts`, `next_attempt_at`, `locked_by`, `locked_until`, `last_error_code`, `last_error_message_redacted`, `manual_review_reason`, `correlation_id`, `created_at`, `updated_at`, `finished_at`

The job read API does not expose `payload_json` or `idempotency_key`.

`job_attempt` response fields:

`id`, `display_id`, `job_id`, `worker_id`, `attempt_number`, `started_at`, `finished_at`, `result`, `error_code`, `error_message_redacted`, `duration_ms`, `correlation_id`

## 4. Route Reference

### 4.0 Catalog Operations

- `GET /admin/catalog/provider-readiness`
  - auth: admin actor, `catalog.manage`
  - query: `product_type`, `status`, `limit`, `cursor`
  - response: list of provider readiness rows with `plan_display_id`, `plan_code`, `plan_name`, `product_type`, `plan_status`, optional `plan_source_display_id`, optional `source_display_id`, `source_name`, `source_type`, `source_status`, `inventory_mode`, `state`, and `reason`
  - states: `ready`, `inactive_source`, `missing_plan_source`, `unsupported_capability`, `fake_provider_only`
  - note: this route does not expose provider credentials, raw provider payloads, or capability JSON

### 4.1 Orders

- `GET /client/orders`
  - auth: client actor, `order.create`
  - query: `display_id`, `status`, `billing_status`, `amount_min`, `amount_max`, `limit`, `cursor`
  - response: list of `order`
  - note: buyer scope is forced to the current actor

- `POST /client/orders`
  - auth: client actor, `order.create`
  - headers: `Idempotency-Key` required
  - body: `tenant_plan_id`, `quantity`, `currency`, `unit_price_minor`, `discount_minor`, `total_minor`, `product_snapshot`, `plan_snapshot`, `price_snapshot`
  - response: one `order`

- `GET /client/orders/{order_id}`
  - auth: client actor, `order.create`
  - response: one `order`
  - note: buyer scope is forced to the current actor

- `POST /client/checkouts`
  - auth: client actor, `order.create`
  - headers: `Idempotency-Key` required
  - body: `order_id`
  - response: invoice detail with `items[]`
  - notes:
    - tenant and buyer scope are forced from request context
    - order must be `order_status=pending_payment` and `billing_status=unpaid`
    - duplicate submits for the same order return the existing invoice instead of creating another invoice

- `GET /admin/orders`
  - auth: admin actor, `order.view`
  - query: `buyer_user_id`, `display_id`, `status`, `billing_status`, `amount_min`, `amount_max`, `limit`, `cursor`
  - response: list of `order`

- `GET /reseller/orders`
  - auth: reseller actor, `order.view`
  - query: `buyer_user_id`, `display_id`, `status`, `billing_status`, `amount_min`, `amount_max`, `limit`, `cursor`
  - response: list of `order`
  - note: tenant scope is forced to the current reseller tenant

- `GET /admin/orders/{order_id}`
  - auth: admin actor, `order.view`
  - response: one `order`

- `PATCH /admin/orders/{order_id}/status`
  - auth: admin actor, `order.manage`
  - body: `from_status`, `to_status`, `billing_status`
  - response: one `order`
  - note: use this route for manual order state changes such as `pending_payment -> paid`

### 4.2 Services

- `GET /client/services`
  - auth: client actor, `service.view`
  - query: `display_id`, `order_id`, `order_display_id`, `status`, `limit`, `cursor`
  - response: list of `service`
  - note: buyer scope is forced to the current actor

- `GET /client/services/{service_id}`
  - auth: client actor, `service.view`
  - response: one `service`

- `GET /admin/services`
  - auth: admin actor, `service.view`
  - query: `buyer_user_id`, `display_id`, `order_id`, `order_display_id`, `status`, `limit`, `cursor`
  - response: list of `service`

- `GET /reseller/services`
  - auth: reseller actor, `service.view`
  - query: `buyer_user_id`, `display_id`, `order_id`, `order_display_id`, `status`, `limit`, `cursor`
  - response: list of `service`
  - note: tenant scope is forced to the current reseller tenant

- `GET /admin/services/{service_id}`
  - auth: admin actor, `service.view`
  - response: one `service`

### 4.3 Invoices

- `GET /client/invoices`
  - auth: client actor, `wallet.view`
  - query: `display_id`, `order_id`, `status`, `amount_min`, `amount_max`, `limit`, `cursor`
  - response: list of `invoice`
  - note: buyer scope is forced to the current actor

- `GET /client/invoices/{invoice_id}`
  - auth: client actor, `wallet.view`
  - response: invoice detail

- `GET /admin/invoices`
  - auth: admin actor, `wallet.view`
  - query: `buyer_user_id`, `display_id`, `order_id`, `status`, `amount_min`, `amount_max`, `limit`, `cursor`
  - response: list of `invoice`

- `GET /reseller/invoices`
  - auth: reseller actor, `wallet.view`
  - query: `buyer_user_id`, `display_id`, `order_id`, `status`, `amount_min`, `amount_max`, `limit`, `cursor`
  - response: list of `invoice`
  - note: tenant scope is forced to the current reseller tenant

- `GET /admin/invoices/{invoice_id}`
  - auth: admin actor, `wallet.view`
  - response: invoice detail

### 4.4 Wallets and ledger

- `GET /client/wallets`
  - auth: client actor, `wallet.view`
  - query: `display_id`, `status`, `limit`, `cursor`
  - response: list of `wallet`
  - note: owner scope is forced to the current actor's user wallet

- `GET /client/wallets/{wallet_id}`
  - auth: client actor, `wallet.view`
  - response: one `wallet`

- `GET /client/wallets/{wallet_id}/ledger`
  - auth: client actor, `wallet.view`
  - query: `display_id`, `direction`, `entry_type`, `status`, `amount_min`, `amount_max`, `limit`, `cursor`
  - response: list of `ledger`

- `GET /admin/wallets`
  - auth: admin actor, `wallet.view`
  - query: `display_id`, `owner_type`, `owner_id`, `status`, `limit`, `cursor`
  - response: list of `wallet`

- `GET /admin/wallets/{wallet_id}`
  - auth: admin actor, `wallet.view`
  - response: one `wallet`

- `GET /admin/wallets/{wallet_id}/ledger`
  - auth: admin actor, `wallet.view`
  - query: `display_id`, `direction`, `entry_type`, `status`, `amount_min`, `amount_max`, `limit`, `cursor`
  - response: list of `ledger`

- `GET /reseller/wallets`
  - auth: reseller actor, `wallet.view`
  - query: `display_id`, `owner_type`, `owner_id`, `status`, `limit`, `cursor`
  - response: list of `wallet`
  - note: tenant scope is forced to the current reseller tenant

- `GET /reseller/wallets/{wallet_id}`
  - auth: reseller actor, `wallet.view`
  - response: one `wallet`

- `GET /reseller/wallets/{wallet_id}/ledger`
  - auth: reseller actor, `wallet.view`
  - query: `display_id`, `direction`, `entry_type`, `status`, `amount_min`, `amount_max`, `limit`, `cursor`
  - response: list of `ledger`

### 4.5 Top-up requests

- `POST /client/topup-requests`
  - auth: client actor, `wallet.view`
  - headers: `Idempotency-Key` required
  - body: `wallet_id`, `amount_minor`, `currency`, `payment_method`, `payment_reference`
  - response: one `topup_request`
  - notes:
    - `payment_method` valid values: `bank_transfer`, `crypto`, `manual`, `other`
    - the wallet must belong to the current actor

- `GET /client/topup-requests`
  - auth: client actor, `wallet.view`
  - query: `display_id`, `wallet_id`, `payment_method`, `status`, `amount_min`, `amount_max`, `limit`, `cursor`
  - response: list of `topup_request`
  - note: requester scope is forced to the current actor

- `GET /client/topup-requests/{topup_request_id}`
  - auth: client actor, `wallet.view`
  - response: one `topup_request`

- `GET /admin/topup-requests`
  - auth: admin actor, `wallet.view`
  - query: `requested_by`, `display_id`, `wallet_id`, `payment_method`, `status`, `amount_min`, `amount_max`, `limit`, `cursor`
  - response: list of `topup_request`

- `GET /admin/topup-requests/{topup_request_id}`
  - auth: admin actor, `wallet.view`
  - response: one `topup_request`

- `GET /reseller/topup-requests`
  - auth: reseller actor, `wallet.view`
  - query: `requested_by`, `display_id`, `wallet_id`, `payment_method`, `status`, `amount_min`, `amount_max`, `limit`, `cursor`
  - response: list of `topup_request`
  - note: tenant scope is forced to the current reseller tenant

- `POST /admin/topup-requests/{topup_request_id}/approve`
  - auth: admin actor, `wallet.topup.approve`
  - body: `review_note`
  - response: one `topup_request`
  - note: approval credits the wallet and returns `ledger_entry_id`

- `POST /admin/topup-requests/{topup_request_id}/reject`
  - auth: admin actor, `wallet.topup.approve`
  - body: `review_note`
  - response: one `topup_request`
  - note: `review_note` is required on reject

### 4.6 Transactions and wallet payment

- `GET /client/transactions`
  - auth: client actor, `wallet.view`
  - query: `display_id`, `order_id`, `invoice_id`, `type`, `status`, `amount_min`, `amount_max`, `limit`, `cursor`
  - response: list of `transaction`
  - note: account scope is forced to the current actor

- `GET /client/transactions/{transaction_id}`
  - auth: client actor, `wallet.view`
  - response: one `transaction`

- `GET /admin/transactions`
  - auth: admin actor, `wallet.view`
  - query: `account_user_id`, `display_id`, `order_id`, `invoice_id`, `type`, `status`, `amount_min`, `amount_max`, `limit`, `cursor`
  - response: list of `transaction`

- `GET /reseller/transactions`
  - auth: reseller actor, `wallet.view`
  - query: `account_user_id`, `display_id`, `order_id`, `invoice_id`, `type`, `status`, `amount_min`, `amount_max`, `limit`, `cursor`
  - response: list of `transaction`
  - note: tenant scope is forced to the current reseller tenant

- `GET /admin/transactions/{transaction_id}`
  - auth: admin actor, `wallet.view`
  - response: one `transaction`

- `POST /client/invoice-wallet-payments`
  - auth: client actor, `wallet.view`
  - headers: `Idempotency-Key` required
  - body: `invoice_id`, `wallet_id`
  - response:
    - `invoice`
    - `transaction`
    - optional `ledger`
    - optional `order` with `id`, `display_id`, `order_status`, `billing_status`
  - notes:
    - wallet and invoice currency must match
    - invoice must still be payable
    - when the invoice has `order_id`, successful payment finalizes the order to `order_status=paid` and `billing_status=paid`
    - paid orders are queued once in `jobs` with `job_type=provider.provision`, `reference_type=order`, and `reference_id=<order_id>`
    - insufficient funds returns `wallet.insufficient_balance`

### 4.7 Payment reconciliation

- `GET /admin/payment-reconciliation`
  - auth: admin actor, `wallet.view`
  - query: `account_user_id`, `display_id`, `status`, `provider`, `invoice_id`, `invoice_display_id`, `wallet_id`, `wallet_display_id`, `amount_min`, `amount_max`, `created_from`, `created_to`, `limit`, `cursor`
  - response: list of `payment_reconciliation`

- `GET /admin/payment-reconciliation/{transaction_id}`
  - auth: admin actor, `wallet.view`
  - response: one `payment_reconciliation`

### 4.8 Audit logs

- `GET /admin/audit-logs`
  - auth: admin actor, `audit.view`
  - query: `actor_id`, `actor_type`, `display_id`, `action`, `target_type`, `target_id`, `created_from`, `created_to`, `limit`, `cursor`
  - response: list of audit log summaries

- `GET /admin/audit-logs/{audit_log_id}`
  - auth: admin actor, `audit.view`
  - response: one audit log detail

### 4.9 Jobs

- `GET /admin/jobs`
  - auth: admin actor, `order.view`
  - query: `display_id`, `job_type`, `status`, `reference_type`, `reference_id`, `source_id`, `limit`, `cursor`
  - response: list of `job`

- `GET /admin/jobs/summary`
  - auth: admin actor, `provisioning.job.view`
  - query: `job_type`, default `provider.provision`
  - response: one job summary with `job_type`, `total`, `attention_count`, `counts`, `oldest_queued_at`, `oldest_queued_age_seconds`, `latest_failure`, `generated_at`
  - note: `latest_failure` only includes display id, status, redacted error fields, manual-review reason, and timestamps

- `GET /admin/jobs/{job_id}`
  - auth: admin actor, `order.view`
  - response: one `job`

- `GET /admin/jobs/{job_id}/attempts`
  - auth: admin actor, `order.view`
  - query: `limit`, `cursor`
  - response: list of `job_attempt`
  - note: attempts are scoped through the parent job tenant

- `POST /admin/jobs/{job_id}/retry`
  - auth: admin actor, `provisioning.job.retry`
  - body: optional `next_attempt_at` as RFC3339; empty body retries now
  - response: one `job`
  - note: only `failed_retryable` and `manual_review` jobs can be requeued

- `POST /admin/jobs/{job_id}/manual-review`
  - auth: admin actor, `provisioning.manual_review.resolve`
  - body: `reason`
  - response: one `job`
  - note: active worker states and completed jobs return `job.status_conflict`

- `POST /admin/jobs/{job_id}/cancel`
  - auth: admin actor, `provisioning.manual_review.resolve`
  - body: optional `reason`
  - response: one `job`
  - note: only safe non-active jobs can be cancelled; succeeded jobs return `job.status_conflict`

- `GET /reseller/jobs`
  - auth: reseller actor, `order.view`
  - query: `display_id`, `job_type`, `status`, `reference_type`, `reference_id`, `source_id`, `limit`, `cursor`
  - response: list of `job`
  - note: tenant scope is forced to the current reseller tenant

- `GET /reseller/jobs/{job_id}`
  - auth: reseller actor, `order.view`
  - response: one `job`
  - note: tenant scope is forced to the current reseller tenant

- `GET /reseller/jobs/{job_id}/attempts`
  - auth: reseller actor, `order.view`
  - query: `limit`, `cursor`
  - response: list of `job_attempt`
  - note: attempts are scoped through the parent job tenant

## 5. Filter and Enum Reference

### 5.1 Order and service status values

- order `status`: `draft`, `pending_payment`, `paid`, `cancelled`, `failed`, `refunded`
- order `billing_status`: `unpaid`, `paid`, `overdue`, `refunded`, `partially_refunded`
- service `status`: `active`, `suspended`, `expired`, `cancelled`, `terminated`

### 5.2 Invoice values

- invoice `status`: `draft`, `issued`, `paid`, `partially_paid`, `overdue`, `voided`

### 5.3 Wallet values

- wallet `owner_type`: `tenant`, `user`, `reseller_settlement`, `platform`
- wallet `status`: `active`, `frozen`, `closed`
- ledger `direction`: `credit`, `debit`
- ledger `entry_type`: `topup`, `purchase`, `reseller_cost`, `refund`, `adjustment`, `reversal`, `commission`, `lock`, `unlock`
- ledger `status`: `posted`, `voided_by_reversal`
- top-up `status`: `draft`, `submitted`, `under_review`, `approved`, `rejected`, `expired`, `cancelled`

### 5.4 Payment values

- transaction `type`: `charge`, `refund`, `adjustment`
- transaction `status`: `pending`, `posted`, `failed`, `voided`

### 5.5 Job values

- job `status`: `queued`, `claimed`, `running`, `succeeded`, `failed_retryable`, `failed_terminal`, `manual_review`, `cancelled`
- job attempt `result`: `succeeded`, `failed_retryable`, `failed_terminal`, `manual_review`, `cancelled`

## 6. Common Error Codes

Shared errors:

- `validation.failed`
- `request.invalid_json`
- `request.limit_invalid`
- `request.limit_too_large`
- `request.display_id_invalid`
- `request.amount_invalid`
- `request.amount_range_invalid`
- `tenant.context_missing`
- `tenant.context_invalid`
- `auth.actor_required`
- `auth.permission_denied`

Route-specific errors that frontend and agents should expect:

- orders: `order.not_found`, `order.status_conflict`, `order.status_transition_invalid`
- services: `service.not_found`, `service.status_invalid`
- invoices: `invoice.not_found`, `invoice.status_conflict`
- wallets: `wallet.not_found`, `wallet.ledger_not_found`
- top-up: `wallet.topup_not_found`, `wallet.topup_status_conflict`, `wallet.payment_method_invalid`
- checkout: `checkout.order_not_checkoutable`
- payment: `payment.transaction_not_found`, `payment.invoice_not_payable`, `payment.idempotency_conflict`, `payment.wallet_currency_mismatch`, `wallet.insufficient_balance`, `order.status_conflict`, `order.provisioning_source_not_found`
- jobs: `job.not_found`, `job.status_invalid`, `job.status_conflict`, `job.manual_review_reason_missing`
- audit: `audit.created_time_invalid`

## 7. Practical Notes For Frontend And Agents

- Use `display_id` in tables and badges for human operators.
- Keep UUID ids for route navigation, row actions, and follow-up API calls.
- Treat all money amounts as integer minor units.
- For client list routes, do not rely on user id filters to widen scope. The backend forces actor-owned scope.
- For money-changing client routes, always send a fresh `Idempotency-Key`.
- For checkout, submit `POST /client/checkouts` after `POST /client/orders`, then pay the returned invoice through `POST /client/invoice-wallet-payments`.
- For audit and reconciliation date filters, send RFC3339 timestamps.
- For list screens, build UI filters around the exact query names in this document rather than inventing aliases.
