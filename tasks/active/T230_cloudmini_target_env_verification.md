# T230 - Cloudmini target-env hardening verification

Status: REVIEW
Owner: Codex
Branch: codex/t230-target-env-t229-verification
PR: https://github.com/Chinsusu/Billing-V2/pull/493
Risk: provider/provisioning/lifecycle/credential/ops
Created: 2026-05-17
Updated: 2026-05-17

## Summary

Deploy or verify the T229 Cloudmini hardening on the approved Billing test server and record redacted target-environment evidence that the backend, worker, and frontend build/runtime are healthy without running a new mutating provider pilot.

## Scope

- Verify the test server has the T229 fail-closed Cloudmini provisioning behavior and provider-backed lifecycle-worker cleanup code.
- Verify target services, ports, build/tests, and redacted runtime configuration needed for the Cloudmini provider registry path.
- Record only redacted target-environment evidence in launch/provider docs.
- Do not call Cloudmini create, delete, action, or Billing checkout/provisioning mutating flows in this task.
- Do not broaden pilot readiness or mark launch GO.

## Acceptance Criteria

- Target environment evidence confirms the deployed code includes T229 Cloudmini fail-closed status handling and lifecycle-worker provider cleanup.
- Target validation runs pass for the relevant Go provider/order/worker packages, taskguard, and any changed docs checks.
- Service health and listening ports are checked on the test server without printing secrets, DSNs, raw provider IDs, raw provider payloads, or credentials.
- Documentation records the target-env verification result and explicitly states that no mutating provider route was called.

## Notes

- This task is a follow-up from T229 and does not replace the still-missing live duplicate-create and timeout-after-send provider evidence.
- Use the approved test server information from the local ops credential file without committing or printing it.

## Agent Log

- 2026-05-17: Task created and claimed by Codex from latest `origin/main` on branch `codex/t230-target-env-t229-verification`.
- 2026-05-17: Synced T229 code to the approved test server at `/opt/Billing`, preserving local env files, credential files, `frontend/node_modules`, frontend build output, and `bin`.
- 2026-05-17: Target server confirmed T229 source markers, passed focused provider/order/worker tests, taskguard, Go API/worker builds, frontend build, service restart, `/healthz`, `/readyz`, frontend HTTP checks, and port checks. No Cloudmini mutating route or Billing provisioning mutation was run.
- 2026-05-17: Local validation passed: `go test ./internal/modules/provider ./internal/modules/order ./cmd/worker`; `go run ./cmd/taskguard`; `git diff --check`; changed-file secret scan found only existing documentation text about avoiding `?token=` query credentials.
- 2026-05-17: Opened PR https://github.com/Chinsusu/Billing-V2/pull/493.
