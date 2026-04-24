# T079 - Checkout payment orchestration

Status: TODO
Owner: -
Branch: codex/t079-checkout-payment-orchestration
PR: -
Risk: backend/billing
Created: 2026-04-24
Updated: 2026-04-24

## Summary

Define and implement the next backend step that turns a client order into a payable billing flow without manual UUID work.

## Scope

- Read the current order, invoice, wallet, and payment modules before choosing the smallest safe implementation.
- Prefer an explicit backend orchestration endpoint or service that can create the expected invoice/payment state from a client order.
- Preserve idempotency for any mutation.
- Document the chosen flow in the API contract.

## Acceptance Criteria

- A client checkout can move from order creation to a payable invoice or equivalent next state through backend APIs.
- The flow does not trust tenant IDs from request bodies.
- Existing invoice and wallet payment behavior remains compatible.
- Tests cover idempotency, tenant scoping, and duplicate-submit behavior.
- `go test ./...` passes and backend binaries build.

## Notes

- Keep provisioning and real provider allocation out of scope unless already supported by existing services.
- If the current domain model blocks a clean implementation, write the blocker clearly and split the smallest follow-up task.

## Agent Log

- 2026-04-24: Task created after T074 exposed the UI action surface and the backend checkout flow became the next gap.
