# T036 - Provisioning queue service

Status: TODO
Owner: -
Branch: feat/provisioning-queue-service
PR: -
Risk: order/provisioning
Created: 2026-04-23
Updated: 2026-04-23

## Summary

Add a service that decides when a paid order should be queued for provisioning and records that queue request once.

## Scope

- Add a small provisioning queue service on top of the order and outbox/job stores.
- Queue provisioning only for orders that are paid and not already queued.
- Keep the queue operation tenant-scoped and idempotent.
- Add tests for duplicate queue requests and wrong-status orders.
- Out of scope: calling real providers or creating service instances.

## Acceptance Criteria

- Paid orders can be queued for provisioning once.
- Pending, canceled, failed, or refunded orders are rejected.
- Repeated queue calls do not create duplicate jobs.
- The queued work contains enough order, tenant, account, and provider data for the worker task.
- Full validation passes: `make fmt`, `make test`, `make build`, `make migrate-validate`.

## Notes

- This task should start after T035 if it relies on order lifecycle events.

## Agent Log

- 2026-04-23: Task created for the next backend batch.
