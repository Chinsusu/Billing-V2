# 68 - Full E2E Quality Gate Runbook

**Date:** 2026-05-14  
**Scope:** Repeatable local/dev launch-critical quality gate for backend, database smoke, API smoke, billing mutation smoke, worker provisioning, and frontend smoke.

## Purpose

This runbook gives operators and reviewers one command for the current launch-critical local/dev proof. It reuses the existing smoke commands instead of replacing them.

The gate proves the current fake-provider local flow:

- database migrations and deterministic seed are valid;
- seeded API read paths and RBAC negative checks pass;
- client top-up request is created and approved;
- client order is created, checked out, invoiced, and paid from wallet;
- payment queues exactly one `provider.provision` job;
- fake-provider worker provisions the service;
- service activation is visible through API;
- frontend install, audit, sensitive-text guard, lint, build, and admin browser smoke pass.

It does not prove real provider sandbox readiness, production auth hardening, production delivery channels, or launch owner sign-off. Those remain separate T205 inputs.

## Safety Boundary

Allowed:

```text
APP_ENV=local
APP_ENV=dev
local/dev PostgreSQL DB_DSN
fake provider registry
frontend smoke with mocked/intercepted data
```

Forbidden:

```text
APP_ENV=staging for this dev-header smoke
APP_ENV=production or APP_ENV=prod
production DB_DSN
provider production credentials
unmasked customer data
real payment or provider endpoints
```

The script refuses production-like markers in `DB_DSN` and `API_BASE_URL`, and refuses `APP_ENV=staging`, `APP_ENV=prod`, or `APP_ENV=production` because the smoke uses local dev actor headers.
The script logs smoke commands with `DB_DSN` redacted.

## Prerequisites

```text
Go
make
Git
curl
Node.js 20-compatible npm
frontend dependencies installable from package-lock.json
Playwright Chromium installed or installable by npm ci / CI setup
local/dev PostgreSQL database approved for mutation smoke
```

The database must be safe to mutate. The gate runs seed idempotently and creates additional local smoke records for top-up, order, invoice, payment, job, and service.

## Command

Set a local/dev DSN and run:

```bash
export APP_ENV=local
export DB_DSN='postgres://billing:billing@localhost:5432/billing_e2e?sslmode=disable'
make full-e2e-quality-gate
```

The script starts a local API on `127.0.0.1:18080` by default. Override if needed:

```bash
export BILLING_E2E_API_ADDR='127.0.0.1:18081'
export API_BASE_URL='http://127.0.0.1:18081'
```

If `frontend/node_modules` is already trusted and the operator wants a faster local rerun:

```bash
export BILLING_E2E_SKIP_NPM_CI=1
make full-e2e-quality-gate
```

Do not use the skip flag for launch evidence unless the exact dependency install was already captured in the same evidence packet.

## Gate Steps

```text
make task-guard
env -u DB_DSN make test
env -u DB_DSN make contract-guard
env -u DB_DSN make error-code-guard
env -u DB_DSN make build
git diff --check
go run ./cmd/smoke -dsn "$DB_DSN" dev-db
start local API with the same DB_DSN
go run ./cmd/smoke -dsn "$DB_DSN" -base-url "$API_BASE_URL" dev-api
go run ./cmd/smoke -dsn "$DB_DSN" -base-url "$API_BASE_URL" dev-billing
npm --prefix frontend ci
npm --prefix frontend audit --omit=dev
npm --prefix frontend run check:sensitive-text
npm --prefix frontend run lint
npm --prefix frontend run build
npm --prefix frontend run smoke:admin:ci
```

## Evidence Template

```text
Gate ID:
Date/time UTC:
Operator:
Environment: local/dev
DB classification: local/dev, no production data
API base URL: local URL only
Command: make full-e2e-quality-gate
Backend result:
DB smoke result:
API smoke result:
Billing mutation result:
Frontend result:
CI result:
Result: pass/fail
Issues:
Follow-up:
```

## T204 Local Evidence

```text
Gate ID: T204-local-20260514T032446Z
Operator: Codex
Environment: local
DB classification: local temporary database billing_t204_e2e_20260514032311, no production data
API base URL: http://127.0.0.1:18080
Command: make full-e2e-quality-gate
Backend result: taskguard, test, contract guard, error code guard, build, and diff check passed
DB smoke result: dev-db passed, 23 migrations present, 19 checks
API smoke result: dev-api passed, 35 checks including RBAC negative checks
Billing mutation result: dev-billing passed; top-up, checkout, wallet payment, provisioning job, worker fulfillment, audit checks, and active service verified
Frontend result: npm ci, audit, sensitive-text guard, lint, build, and admin browser smoke passed
Result: pass
Issues found and fixed: catalog admin auth missing tenant header context; smoke job JSON UUID cast; jobs claim RETURNING ambiguity; DB_DSN redaction in gate logs
Follow-up: repeat against approved shared staging inputs before final T205 launch sign-off
```
