# T096 - Provisioning ops readiness docs

Status: TODO
Owner: -
Branch: codex/t096-provisioning-ops-readiness-docs
PR: -
Risk: docs/ops
Created: 2026-04-24
Updated: 2026-04-24

## Summary

Consolidate the provisioning operations flow into a concise operator checklist for local/sandbox readiness.

## Scope

- Work mainly in `docs/05_development_standards/**/*` and task references.
- Document paid-order fulfillment, worker run modes, job recovery actions, smoke verification, and common failure decisions.
- Prefer short operator steps and exact commands over long theory.
- Keep new docs under 500 lines.

## Acceptance Criteria

- Operators have one concise checklist for local/sandbox provisioning readiness.
- The checklist links to API reference, worker command, smoke command, and recovery action docs.
- It states what not to do: do not pay invoices twice, do not edit money rows by hand, and do not retry without provider-state checks.
- Documentation-only validation is still backed by backend/frontend gates if CI requires them.

## Notes

- Should follow T086, T091, and T092.

## Agent Log

- 2026-04-24: Task created after T092 completed and the active board was empty.
