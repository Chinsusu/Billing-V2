# Provisioning Ops Readiness Checklist

**Version:** v1.11
**Date:** 2026-04-24
**Scope:** Local and sandbox readiness checks for paid-order fulfillment, provisioning worker runs, job recovery, and smoke verification.

## 1. Read This First

Use this checklist before calling local/sandbox provisioning ready.

References:

- API routes and response fields: `docs/05_development_standards/56_Billing_API_Operational_Reference.md`
- Full billing operations runbook: `docs/05_development_standards/57_Billing_Operations_Runbook.md`
- Local setup and smoke commands: `docs/05_development_standards/55_Local_Development_Runbook.md`
- Worker command: `cmd/worker`
- Smoke command: `cmd/smoke dev-billing`

## 2. Required Context

Set these values for local/sandbox checks:

```bash
export API_BASE_URL="http://localhost:8080"
export DB_DSN="postgres://billing:billing@localhost:5432/billing?sslmode=disable"
export TENANT_ID="00000000-0000-0000-0000-000000000010"
export ADMIN_ID="00000000-0000-0000-0000-000000000102"
export CUSTOMER_ID="00000000-0000-0000-0000-000000000103"
```

Rules:

- API and worker must use the same `DB_DSN`.
- `APP_ENV` must not be `prod` or `production`.
- Use local fake provider unless a sandbox provider is explicitly approved.
- Do not paste provider credentials, wallet secrets, or production DSNs into task files, PRs, logs, or docs.

## 3. Provider Source Readiness

Before paid-order provisioning tests, inspect active plan/source readiness through the API instead of reading catalog tables by hand:

```bash
curl -s "$API_BASE_URL/admin/catalog/provider-readiness?status=active&limit=100" \
  -H "X-Tenant-Id: $TENANT_ID" \
  -H "X-Actor-Id: $ADMIN_ID" \
  -H "X-Actor-Type: reseller_owner" \
  -H "X-Actor-Tenant-Id: $TENANT_ID"
```

States:

- `ready`: active plan source, active provider source, and automatic provisioning capability match the product type.
- `inactive_source`: the plan source or provider source is not active.
- `missing_plan_source`: the active plan has no provider source link.
- `unsupported_capability`: the source cannot auto-provision the plan product type.
- `fake_provider_only`: local fake/manual path can run smoke, but this is not production provider readiness.

The response uses plan/source display IDs first and does not expose provider credentials, raw provider payloads, or capability JSON.

Fresh local seed examples:

| State | Plan display ID | Source display ID | Plan code | Use |
|---|---:|---:|---|---|
| `ready` | `10000` | `10001` | `vps-cx23-40gb-monthly` | Green-path paid-order smoke. |
| `fake_provider_only` | `10001` | `10000` | `vps-cx33-80gb-monthly` | Manual local fallback example. |
| `unsupported_capability` | `10002` | `10001` | `proxy-static-10gb-monthly` | Active source that cannot provision this product type. |
| `inactive_source` | `10003` | `10002` | `vps-maintenance-example-monthly` | Non-cloned maintenance example for ops checks. |

If an existing local database has been reseeded many times, display IDs may differ. Use the plan code to find the same scenario.

## 4. Green Path Checklist

Run the deterministic smoke first:

```bash
go run ./cmd/smoke -dsn "$DB_DSN" -base-url "$API_BASE_URL" dev-billing
```

Expected result:

- client top-up is approved;
- client order is created with an `Idempotency-Key`;
- checkout creates an issued invoice;
- wallet payment marks invoice paid and creates a posted transaction;
- order becomes `order_status=paid` and `billing_status=paid`;
- one `provider.provision` job exists for the paid order;
- fake-provider worker processes the job;
- `GET /client/services?order_id=<order_id>` returns the active paid service.

If this smoke fails, stop and inspect the failed step before retrying payment or provisioning.

## 5. Worker Run Modes

Run one pass:

```bash
go run ./cmd/worker provision-once -dsn "$DB_DSN"
```

Run a local/sandbox loop:

```bash
go run ./cmd/worker provision-loop -dsn "$DB_DSN" -interval 5s -batch-size 10
```

Use `provision-once` when checking a known stuck job. Use `provision-loop` while testing repeated local fulfillment. Stop loop mode with `Ctrl+C`, or add `-timeout 5m` for a bounded run.

## 6. Inspect A Paid Order

Use display IDs in human discussion, but UUIDs in API paths and filters.

Client checks:

```bash
curl -s "$API_BASE_URL/client/orders/$ORDER_ID" \
  -H "X-Tenant-Id: $TENANT_ID" \
  -H "X-Actor-Id: $CUSTOMER_ID" \
  -H "X-Actor-Type: client" \
  -H "X-Actor-Tenant-Id: $TENANT_ID"

curl -s "$API_BASE_URL/client/services?order_id=$ORDER_ID&limit=20" \
  -H "X-Tenant-Id: $TENANT_ID" \
  -H "X-Actor-Id: $CUSTOMER_ID" \
  -H "X-Actor-Type: client" \
  -H "X-Actor-Tenant-Id: $TENANT_ID"
```

Admin job checks:

```bash
curl -s "$API_BASE_URL/admin/jobs/summary?job_type=provider.provision" \
  -H "X-Tenant-Id: $TENANT_ID" \
  -H "X-Actor-Id: $ADMIN_ID" \
  -H "X-Actor-Type: reseller_owner" \
  -H "X-Actor-Tenant-Id: $TENANT_ID"

curl -s "$API_BASE_URL/admin/jobs?job_type=provider.provision&reference_type=order&reference_id=$ORDER_ID&limit=20" \
  -H "X-Tenant-Id: $TENANT_ID" \
  -H "X-Actor-Id: $ADMIN_ID" \
  -H "X-Actor-Type: reseller_owner" \
  -H "X-Actor-Tenant-Id: $TENANT_ID"

curl -s "$API_BASE_URL/admin/jobs/$JOB_ID/attempts?limit=20" \
  -H "X-Tenant-Id: $TENANT_ID" \
  -H "X-Actor-Id: $ADMIN_ID" \
  -H "X-Actor-Type: reseller_owner" \
  -H "X-Actor-Tenant-Id: $TENANT_ID"
```

The attempts view must show redacted errors only. If an error includes a secret or raw provider credential, stop and open a security fix task.

## 7. Recovery Decision Table

| Situation | Do | Do not |
|---|---|---|
| Order is paid and job is `queued` | Start `provision-once` or `provision-loop`. | Do not pay the invoice again. |
| Job is `failed_retryable` | Check provider source/config and attempt error, then retry only if provider state is known. | Do not retry blindly when external resource state is unclear. |
| Job is `manual_review` | Record the reason, verify provider state, then retry or cancel through admin job actions. | Do not clear manual review without evidence. |
| Job is `failed_terminal` | Keep the order paid, escalate, and only manual-review/cancel after provider-state check. | Do not force a service record by hand. |
| No job exists for a paid order | Open a backend defect task with order/invoice/job evidence. | Do not create a second invoice or second payment. |
| Provider created a resource but service is missing | Open a repair task with provider evidence and job attempts. | Do not invent `external_resource_id`. |

## 8. Recovery Actions

Retry after provider-state check:

```bash
curl -s -X POST "$API_BASE_URL/admin/jobs/$JOB_ID/retry" \
  -H "Content-Type: application/json" \
  -H "X-Tenant-Id: $TENANT_ID" \
  -H "X-Actor-Id: $ADMIN_ID" \
  -H "X-Actor-Type: reseller_owner" \
  -H "X-Actor-Tenant-Id: $TENANT_ID" \
  -d '{"next_attempt_at":"2026-04-24T00:00:00Z"}'
```

Move to manual review:

```bash
curl -s -X POST "$API_BASE_URL/admin/jobs/$JOB_ID/manual-review" \
  -H "Content-Type: application/json" \
  -H "X-Tenant-Id: $TENANT_ID" \
  -H "X-Actor-Id: $ADMIN_ID" \
  -H "X-Actor-Type: reseller_owner" \
  -H "X-Actor-Tenant-Id: $TENANT_ID" \
  -d '{"reason":"provider response needs manual verification"}'
```

Cancel only when the job is safe to stop:

```bash
curl -s -X POST "$API_BASE_URL/admin/jobs/$JOB_ID/cancel" \
  -H "Content-Type: application/json" \
  -H "X-Tenant-Id: $TENANT_ID" \
  -H "X-Actor-Id: $ADMIN_ID" \
  -H "X-Actor-Type: reseller_owner" \
  -H "X-Actor-Tenant-Id: $TENANT_ID" \
  -d '{"reason":"order should not provision after provider-state check"}'
```

## 9. Hard No-Go Rules

- Do not pay an invoice twice to create another job.
- Do not edit wallet ledger, payment transaction, invoice, or wallet balance rows by hand.
- Do not retry when provider state is unknown.
- Do not mark provider success from payment success alone.
- Do not run local worker commands against production.
- Do not expose provider credentials or raw provider responses in UI, API errors, logs, tasks, or PRs.

## 10. Ready For Handoff

Before handoff, record:

- order display ID and UUID;
- invoice display ID and UUID;
- payment transaction display ID and UUID;
- job display ID and UUID;
- latest attempt display ID, result, error code, and redacted message;
- worker command used;
- recovery action taken, if any;
- smoke command result.

For CI-backed docs changes, run:

```bash
go test ./...
go build ./cmd/api ./cmd/migrate ./cmd/seed ./cmd/smoke ./cmd/worker
cd frontend
npm ci
npm audit --omit=dev
npm run lint
npm run build
```
