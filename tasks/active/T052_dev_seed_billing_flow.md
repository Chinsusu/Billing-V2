# T052 - Dev seed billing flow

Status: DONE
Owner: Codex
Branch: feat/dev-seed-billing-flow
PR: https://github.com/Chinsusu/Billing-V2/pull/118
Risk: seed/dev
Created: 2026-04-23
Updated: 2026-04-23

## Summary

Extend seed data so local development has a realistic billing flow across catalog, order, wallet, invoice, and payment records.

## Scope

- Seed a customer wallet, ledger history, invoice, invoice items, and payment transaction examples.
- Keep all seeded entities tenant-scoped and linked by UUID plus numeric display IDs.
- Document test actors and sample API calls in the seed/runbook notes.
- Avoid secrets or real provider credentials.

## Acceptance Criteria

- `go run ./cmd/seed` creates a coherent billing demo flow.
- Seed can be rerun without duplicate logical records.
- Sample data supports frontend/admin screens for wallet, invoice, payment, and order views.
- Focused seed tests and full validation pass.

## Notes

- This task should follow core schema/API tasks so seed uses real stores where possible.

## Agent Log

- 2026-04-23: Task created to make end-to-end local testing practical.
- 2026-04-23: Added demo customer and linked wallet/order/service/invoice/payment seed flow plus runbook notes.
- 2026-04-23: PR #118 merged after GitHub checks passed.
