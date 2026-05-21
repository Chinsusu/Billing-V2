# T294 - Scope intake and preflight runbook

Status: IN_PROGRESS
Owner: Codex
Branch: codex/t294-scope-intake-preflight
PR: -
Risk: production scope, private beta scope, launch evidence, credentials, provider, notification
Created: 2026-05-21
Updated: 2026-05-21

## Summary

Add a production/private-beta scope intake and preflight runbook so future GO requests have concrete owner fields, evidence commands, guardrails, and pause criteria before any broader launch decision.

## Scope

- In scope: add a redacted intake/preflight runbook for broader private beta or production scope.
- In scope: link the runbook from the production/private-beta decision packet and docs index.
- Out of scope: approving production, running production commands, collecting new runtime evidence, or storing secrets/customer data.

## Acceptance Criteria

- Runbook defines required scope fields, owner approvals, safe preflight commands, evidence boundaries, and pause criteria.
- Runbook includes a fillable evidence packet template that excludes raw secrets, DSNs, customer data, cookies, provider payloads, and credentials.
- Existing decision docs link to the runbook as the next step before any NO-GO row can be changed to GO.
- Docs-only validation passes before PR: `go run ./cmd/taskguard`, `git diff --check`, touched-file line count, added-line secret-pattern scan, and added-line docs/task UUID scan.

## Notes

- This is a documentation/runbook task only. It does not deploy, mutate provider state, send notifications, or approve a broader launch.

## Agent Log

- 2026-05-21: Task created and claimed on `codex/t294-scope-intake-preflight`.
- 2026-05-21: Added scope intake/preflight runbook and linked it from doc 86 and docs index without approving broader launch.
- 2026-05-21: Validation passed: `go run ./cmd/taskguard`, `git diff --check`, touched-file line count, added-line secret-pattern scan, and added-line docs/task UUID scan.
