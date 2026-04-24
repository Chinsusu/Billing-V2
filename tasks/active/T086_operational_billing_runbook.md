# T086 - Operational billing runbook

Status: REVIEW
Owner: Codex
Branch: codex/t086-operational-billing-runbook
PR: https://github.com/Chinsusu/Billing-V2/pull/196
Risk: docs/ops
Created: 2026-04-24
Updated: 2026-04-24

## Summary

Document the end-to-end billing operations flow, common failure modes, and rollback/recovery steps for order, invoice, wallet, payment, and provisioning states.

## Scope

- Work mainly in `docs/05_development_standards/**/*` and any existing operational reference docs.
- Use simple wording; avoid overly specialized terms where a direct explanation is clearer.
- Include commands and API paths that agents can run locally.

## Acceptance Criteria

- Runbook explains the live flow from catalog/order to checkout invoice to wallet payment to provisioning follow-up.
- Runbook lists common failure cases: insufficient balance, duplicate submit, checkout conflict, invoice already paid, provisioning stuck.
- Runbook gives clear inspect commands or API routes for each failure case.
- Documentation stays consistent with current route names and smoke commands.

## Notes

- This can be done after T081-T085 if those tasks change the final flow.
- Keep it as an operational guide, not a broad architecture rewrite.

## Agent Log

- 2026-04-24: Task created after the checkout/payment/smoke batch to keep operational docs aligned with the live flow.
- 2026-04-24: Codex claimed the task after T085 merged and started drafting the operational billing runbook from current API routes and smoke commands.
- 2026-04-24: Added the billing operations runbook with normal flow, inspection commands, common failure handling, and rollback/recovery rules. Linked it from the local runbook and passed Go/frontend gates locally.
- 2026-04-24: Opened PR https://github.com/Chinsusu/Billing-V2/pull/196 for review.
