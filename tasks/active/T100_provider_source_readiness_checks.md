# T100 - Provider source readiness checks

Status: TODO
Owner: -
Branch: codex/t100-provider-source-readiness-checks
PR: -
Risk: backend/provider
Created: 2026-04-24
Updated: 2026-04-24

## Summary

Add local/sandbox checks that make provider source readiness visible before paid orders depend on provisioning.

## Scope

- Work mainly in provider/catalog backend modules, admin provider docs, and focused tests.
- Add a read-only readiness check for active provider sources used by plans.
- Report simple states such as ready, inactive source, missing plan source, unsupported capability, or fake-provider only.
- Keep checks local/sandbox friendly and avoid credential disclosure.
- Keep each file under 500 lines.

## Acceptance Criteria

- Admin/operator can inspect provider source readiness without reading raw database rows.
- Readiness output identifies sources/plans by display ID first.
- Checks do not leak provider credentials or raw provider payloads.
- Local fake provider remains supported for smoke and worker runs.
- Backend and frontend validation commands pass.

## Notes

- Should follow T068, T082, T092, and T096.
- Do not add production provider polling or deployment automation in this task.

## Agent Log

- 2026-04-24: Task created after T096 completed and the active board was fully DONE.
