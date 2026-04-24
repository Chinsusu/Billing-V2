# T082 - Paid order provisioning enqueue

Status: TODO
Owner: -
Branch: codex/t082-paid-order-provisioning-enqueue
PR: -
Risk: backend/provisioning
Created: 2026-04-24
Updated: 2026-04-24

## Summary

After an order becomes paid, enqueue the existing provisioning path so paid client orders can move toward service creation without manual admin action.

## Scope

- Work mainly in `internal/modules/order/**/*`, `internal/modules/jobs/**/*`, provider worker wiring, and smoke/docs where needed.
- Reuse existing provisioning queue and provider adapter interfaces.
- Keep real provider allocation out of scope unless already supported by the current fake/demo provider path.
- Preserve idempotency so a paid order is queued once.

## Acceptance Criteria

- A newly paid order can create or reuse a provisioning job through backend code.
- Duplicate payment/order-finalization does not create duplicate provisioning jobs.
- Tests cover queue creation, duplicate behavior, and invalid order state.
- Operational docs explain what happens after payment and where to inspect stuck provisioning jobs.
- `go test ./...` passes and backend binaries build.

## Notes

- This task depends on T081.
- If the current queue model blocks automatic enqueue safely, document the blocker and split the smallest follow-up.

## Agent Log

- 2026-04-24: Task created as the next backend step after checkout/payment/order state consistency.

