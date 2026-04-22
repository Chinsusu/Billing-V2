# T013 - Provider Registry

Status: DONE
Owner: Codex
Branch: feat/provider-registry
PR: https://github.com/Chinsusu/Billing-V2/pull/34
Risk: provider
Created: 2026-04-23
Updated: 2026-04-23

## Summary

Add a provider adapter registry and fake registry wiring so provisioning code can resolve adapters by provider type without hardcoding fake adapters in services.

## Scope

- Add registry registration and lookup APIs for provider adapters.
- Add adapter lookup from `OperationContext.ProviderSourceSnapshot.ProviderType`.
- Add default capability profiles for manual, VPS, and proxy provider categories.
- Wire `NewFakeAdapter` and `NewFakeRegistry` for dev/test flows.
- Keep real provider clients, credentials, source persistence, and provisioning service wiring out of scope.

## Acceptance Criteria

- Registry rejects nil, missing provider type, and duplicate provider adapters.
- Registry returns adapters by provider type and operation source snapshot.
- Fake registry creates a stable default provider set.
- Tests cover registry lookup, duplicate handling, fake defaults, and capability category differences.
- `make fmt`, `make test`, `make build`, and `make migrate-validate` pass.

## Notes

- T013 builds on T007 provider adapter interfaces.

## Agent Log

- 2026-04-23: Codex claimed task from `origin/main` using isolated worktree `/tmp/Billing-T013`.
- 2026-04-23: Opened PR #34 after `make fmt`, `make test`, `make build`, and `make migrate-validate` passed.
- 2026-04-23: PR #34 merged; T013 marked DONE.
