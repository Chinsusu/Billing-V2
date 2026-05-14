# T211 - Proxy Cloudmini API V3 provider adapter

Status: REVIEW
Owner: Codex
Branch: codex/t211-cloudmini-v3-provider
PR: https://github.com/Chinsusu/Billing-V2/pull/453
Risk: provider/provisioning/credential
Created: 2026-05-14
Updated: 2026-05-14

## Summary

Add a Billing provider adapter for the Proxy Cloudmini API V3 contract found in `/opt/proxy-cloudmini`, using local contract tests only.

## Scope

- Add a provider type for Proxy Cloudmini API V3.
- Implement adapter behavior for V3 health/capabilities, inventory, create/poll provision, status read, terminate, suspend, unsuspend, and change IP where supported.
- Encrypt returned proxy credentials through an injected cipher before returning `OperationResult`.
- Add unit tests with `httptest` for success, idempotency headers, out-of-stock, auth denial, timeout/manual-review mapping, and unsupported actions.
- Update provider sandbox evidence with the Cloudmini V3 code-read contract and remaining real-sandbox blockers.
- Do not copy provider credentials, production URLs, raw provider responses, or customer data from `/opt/proxy-cloudmini`.
- Do not call real Cloudmini production or sandbox endpoints in this task.

## Acceptance Criteria

- Adapter sends V3 auth and idempotency headers without logging secrets.
- Provision maps async V3 operations to Billing `OperationResult` safely.
- Confirmed capacity errors map to `PROVIDER_OUT_OF_STOCK`.
- Unknown create/action outcomes map to manual review instead of blind retry.
- Credential plaintext is only used transiently in memory and tests verify the returned envelope is encrypted.
- Required provider tests pass locally.

## Notes

- Source reference is `/opt/proxy-cloudmini` API V3 code. That repo is read-only for this task.
- Real sandbox evidence remains blocked until a non-production base URL, scoped credential, owner, quota, SKU/group mapping, timeout policy, and cleanup owner are approved outside git.

## Agent Log

- 2026-05-14: Task created and claimed on `codex/t211-cloudmini-v3-provider`.
- 2026-05-14: Added Cloudmini V3 adapter, provider type registration, local contract tests, and code-read evidence notes.
- 2026-05-14: Validation passed: `make fmt`; `go test ./internal/modules/provider`; `go test ./internal/modules/provider -run SandboxContract`; `go test ./internal/modules/catalog`; `go test ./internal/modules/order -run ProviderProvisioningHandler`; `make test`; `make build`; `go run ./cmd/taskguard`; `git diff --check`.
- 2026-05-14: Opened PR https://github.com/Chinsusu/Billing-V2/pull/453 and moved task to `REVIEW`.
