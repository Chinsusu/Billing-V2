# API Error Code Drift Guard

**Scope:** Lightweight guard for keeping stable API error response fields and public error codes aligned with backend handlers.

## Command

Run:

```bash
make error-code-guard
```

This runs:

```bash
go run ./cmd/errorcodeguard
```

CI runs the same command on pull requests and pushes to `main`.

## What It Checks

The guard checks:

- `internal/platform/httpserver/response.go` still exposes the standard error envelope fields: `error`, `code`, `message`, `details`, `fields`, and `request_id`;
- `docs/05_development_standards/50_API_Response_Error_Logging_Standard.md` still documents the envelope shape;
- `docs/05_development_standards/56_Billing_API_Operational_Reference.md` still lists tracked stable public error codes;
- selected backend handler files still contain the stable error codes that frontend and agents are expected to handle.

It is not a full static analyzer for every validation field code. Track codes that are public, stable, and useful for frontend/support behavior.

## Tracked Stable API Error Codes

### Shared codes

- `validation.failed`
- `request.invalid_json`
- `request.method_not_allowed`
- `request.limit_invalid`
- `request.limit_too_large`
- `request.display_id_invalid`
- `request.amount_invalid`
- `request.amount_range_invalid`
- `tenant.context_missing`
- `tenant.context_invalid`
- `tenant.context_mismatch`
- `auth.actor_required`
- `auth.permission_denied`
- `auth.reason_required`

### Route-specific codes

- catalog: `catalog.not_found`
- orders: `order.not_found`, `order.status_conflict`, `order.status_transition_invalid`, `order.idempotency_conflict`, `order.provisioning_source_not_found`
- services: `service.not_found`, `service.status_invalid`, `service.status_conflict`, `service.status_transition_invalid`, `service.lifecycle_action_invalid`, `service.reason_missing`, `service.suspension_reason_invalid`, `service.billing_cycle_invalid`, `service.billing_cycle_value_invalid`, `service.renewal_unavailable`, `credential.not_found`, `credential.reveal_rate_limited`, `credential.reveal_denied`
- invoices: `invoice.not_found`, `invoice.status_conflict`
- wallets: `wallet.not_found`, `wallet.ledger_not_found`, `wallet.topup_not_found`, `wallet.topup_status_conflict`, `wallet.payment_method_invalid`, `wallet.status_conflict`, `wallet.currency_mismatch`, `wallet.idempotency_conflict`, `wallet.insufficient_balance`
- checkout: `checkout.order_not_checkoutable`
- payment: `payment.transaction_not_found`, `payment.invoice_not_payable`, `payment.idempotency_conflict`, `payment.wallet_currency_mismatch`
- jobs: `job.not_found`, `job.status_invalid`, `job.status_conflict`, `job.manual_review_reason_missing`
- audit: `audit.created_time_invalid`, `audit.log_not_found`

## When To Update

Update `cmd/errorcodeguard/main.go` in the same PR when you intentionally:

- add a new stable public error code;
- rename or remove a tracked code;
- move a tracked code to another handler file;
- change the error response envelope fields;
- update the operational API reference error section.

## Failure Handling

If `make error-code-guard` fails:

1. Check whether backend behavior changed intentionally.
2. If yes, update the operational API reference and guard manifest together.
3. If no, restore the missing error code or response envelope field.
