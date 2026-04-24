# T084 - Reseller fulfillment visibility

Status: IN_PROGRESS
Owner: Codex
Branch: codex/t084-reseller-fulfillment-visibility
PR: -
Risk: frontend/reseller
Created: 2026-04-24
Updated: 2026-04-24

## Summary

Make reseller views show the fulfillment state of paid client orders and provisioning jobs clearly enough for support follow-up.

## Scope

- Work mainly in reseller frontend screens and existing order/service/provisioning read APIs.
- Prefer live read APIs already available; create a backend follow-up only if a necessary read model is missing.
- Keep rows numeric-display-ID first.
- Keep file size under 500 lines.

## Acceptance Criteria

- Reseller can distinguish pending payment, paid pending provisioning, provisioning, active service, failed/manual review states.
- Reseller order/service screens link the relevant order display ID, service display ID, and job/status where available.
- API errors are visible with explicit fallback behavior.
- Frontend validation commands pass.

## Notes

- This task should follow T082 if it depends on newly queued provisioning jobs.
- Do not build new mutation controls in this task.

## Agent Log

- 2026-04-24: Task created to close the reseller support visibility gap after live billing flow work.
- 2026-04-24: Codex claimed the task after T082/T083 merged and started inspecting reseller order/service visibility against current read APIs.
- 2026-04-24: Added reseller fulfillment visibility from existing order/service reads: billing and service screens now show order display IDs, fulfillment status, service display IDs or inferred provider.provision jobs, with partial-data fallback text. Local frontend/backend gates passed.
