# T196 - Reservation TTL and concurrency proof

Status: REVIEW
Owner: Codex
Branch: codex/t196-reservation-ttl-concurrency
PR: https://github.com/Chinsusu/Billing-V2/pull/423
Risk: order, inventory locking, checkout, provider provisioning, and money safety
Created: 2026-05-13
Updated: 2026-05-13

## Summary

Verify and complete reservation TTL and concurrency behavior for checkout and inventory safety.

## Scope

- Audit current reservation and checkout code against the launch checklist.
- Add or complete reservation TTL expiry behavior if missing.
- Add concurrency tests proving limited stock cannot be over-reserved.
- Ensure failed/expired reservations do not incorrectly debit or provision.

## Acceptance Criteria

- Concurrent checkout for limited stock allows only the expected number of reservations.
- Expired reservations release stock according to documented policy.
- Tests cover success, out-of-stock, expiration, and concurrency.
- Relevant backend validation and CI pass.

## Notes

- Stop and ask if TTL duration or release policy is not already defined by docs/code.

## Agent Log

- 2026-05-13: Task created by Codex backlog planning.
- 2026-05-13: Claimed by Codex on branch `codex/t196-reservation-ttl-concurrency`.
- 2026-05-13: Implemented provider inventory counters, reservation quantity, atomic reserve SQL, idempotent expiry release SQL, and reservation concurrency tests. Local validation passed: `make fmt`, `make test`, `make build`, `make migrate-validate`, `make contract-guard`, `make error-code-guard`, `make task-guard`, `git diff --check`.
- 2026-05-13: Opened PR https://github.com/Chinsusu/Billing-V2/pull/423 and moved task to REVIEW.
