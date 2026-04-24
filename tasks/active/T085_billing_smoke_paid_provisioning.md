# T085 - Billing smoke paid provisioning

Status: TODO
Owner: -
Branch: codex/t085-billing-smoke-paid-provisioning
PR: -
Risk: QA/smoke
Created: 2026-04-24
Updated: 2026-04-24

## Summary

Extend the billing smoke path so it verifies the post-payment order/provisioning state introduced by T081-T082.

## Scope

- Work mainly in `cmd/smoke/**/*`, smoke runbook docs, and deterministic seed notes.
- Extend the existing `dev-billing` flow rather than adding brittle UI automation.
- Keep failures specific enough for another agent to diagnose quickly.

## Acceptance Criteria

- Smoke confirms order state after wallet invoice payment.
- Smoke confirms provisioning job or service state when T082 makes it available.
- Required seed data and environment variables remain documented.
- Backend and frontend validation commands still pass.

## Notes

- This task depends on T081 and likely T082.
- Keep API smoke deterministic; avoid relying on real provider credentials.

## Agent Log

- 2026-04-24: Task created after T080 updated the smoke path to use live checkout.

