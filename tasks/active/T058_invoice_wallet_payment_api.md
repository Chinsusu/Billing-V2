# T058 - Invoice wallet payment API

Status: DONE
Owner: Codex
Branch: feat/invoice-wallet-payment-api
PR: https://github.com/Chinsusu/Billing-V2/pull/132
Risk: API/payment/money
Created: 2026-04-23
Updated: 2026-04-23

## Summary

Expose a client API endpoint that lets a customer pay an issued or overdue invoice from a tenant-scoped wallet.

## Scope

- Add a client HTTP mutation route for invoice wallet payment.
- Require tenant context, actor context, RBAC, and idempotency key.
- Validate wallet ownership through the payment service flow.
- Return invoice, payment transaction, and ledger entry display IDs where practical.
- Keep audit event behavior from T057 active for the mutation.

## Acceptance Criteria

- Payment endpoint succeeds for a valid client wallet and payable invoice.
- Duplicate idempotency requests do not double charge.
- Cross-tenant or wrong-owner wallet attempts fail.
- Handler and service tests cover success and validation errors.
- Backend quality gates pass.

## Notes

- UUID remains the API path identifier; display ID is returned for support/UI.

## Agent Log

- 2026-04-23: Task created for the next backend payment action.
- 2026-04-23: Codex started the client invoice wallet payment API endpoint.
- 2026-04-23: PR #132 merged. Added the client wallet payment endpoint and passed backend gates.
