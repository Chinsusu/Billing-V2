# T089 - Provisioning worker command

Status: TODO
Owner: -
Branch: codex/t089-provisioning-worker-command
PR: -
Risk: backend/worker
Created: 2026-04-24
Updated: 2026-04-24

## Summary

Add a runnable worker command for `provider.provision` jobs so local/dev can process queued paid-order provisioning without hand-driving internals.

## Scope

- Work mainly in `cmd/worker/**/*`, `internal/app` only if needed, and runbook docs.
- Wire the existing jobs runner and order provider provisioning handler.
- Provide a deterministic `run-once` mode suitable for smoke tests.
- Keep fake/sandbox provider defaults; do not require real provider credentials.
- Keep each file under 500 lines.

## Acceptance Criteria

- `go run ./cmd/worker provision-once` or equivalent processes claimable `provider.provision` jobs once.
- Command supports DB DSN and local safety guard equivalent to smoke commands.
- Successful fake provisioning creates or updates the expected service state.
- Runbook documents how to run the worker locally.
- Backend and frontend validation commands pass.

## Notes

- This should build on existing `jobs.Runner` and `order.ProviderProvisioningHandler`.
- Avoid long-running daemon behavior unless the code already has a simple pattern.

## Agent Log

- 2026-04-24: Task created in the provisioning operations batch after T086.
