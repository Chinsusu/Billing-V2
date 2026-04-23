# T037 - Provider provisioning worker

Status: TODO
Owner: -
Branch: feat/provider-provisioning-worker
PR: -
Risk: worker/provider
Created: 2026-04-23
Updated: 2026-04-23

## Summary

Add a worker that claims provisioning jobs, calls the provider adapter registry, and records the provisioning result.

## Scope

- Add a worker entry point for provisioning jobs.
- Use the existing provider registry and normalized provider errors.
- Record success, retryable failure, and permanent failure in provisioning records.
- Keep fake provider tests deterministic.
- Out of scope: real provider credentials, background scheduler deployment, or frontend changes.

## Acceptance Criteria

- The worker claims one job at a time through the existing job store pattern.
- Successful provider calls create or update provisioning state for the order item.
- Retryable provider errors are kept retryable without losing the job payload.
- Permanent provider errors are recorded clearly for admin review.
- Full validation passes: `make fmt`, `make test`, `make build`, `make migrate-validate`.

## Notes

- This task should start after T036.

## Agent Log

- 2026-04-23: Task created for the next backend batch.
