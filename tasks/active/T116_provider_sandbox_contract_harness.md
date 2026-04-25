# T116 - Provider sandbox contract harness

Status: TODO
Owner: -
Branch: codex/t116-provider-sandbox-contract-harness
PR: -
Risk: provider/testing
Created: 2026-04-25
Updated: 2026-04-25

## Summary

Create the first provider sandbox contract test harness shape so future real provider integrations have a consistent validation path.

## Scope

- Add a small harness for provider adapter sandbox contract cases without requiring real credentials.
- Cover at least one fake provider scenario for quote, order, status, cancel, and idempotency behavior where the current adapter interface supports it.
- Document how real provider sandbox credentials should be plugged in later without committing secrets.
- Keep each edited file under 500 lines.

## Acceptance Criteria

- The harness can run in local/CI mode without external network credentials.
- It defines clear pass/fail expectations for provider adapter behavior.
- Docs link the harness to the provider sandbox checklist.
- Existing tests pass.

## Notes

- Do not add real provider credentials.
- This task should prefer test helpers over production abstractions unless a production type is already needed.

## Agent Log

- 2026-04-25: Task created in the board and delivery hardening batch.
