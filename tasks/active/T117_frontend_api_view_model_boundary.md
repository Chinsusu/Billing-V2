# T117 - Frontend API view model boundary

Status: TODO
Owner: -
Branch: codex/t117-frontend-api-view-model-boundary
PR: -
Risk: frontend/API
Created: 2026-04-25
Updated: 2026-04-25

## Summary

Keep frontend screens from depending directly on risky raw API fields by adding a small view model boundary for admin data.

## Scope

- Review admin frontend data usage for account, provider, order, invoice, transaction, and log identifiers.
- Add mapping helpers where UI should display numeric public IDs and redact or ignore sensitive backend IDs.
- Keep shared mapping code outside screen components when it is reused.
- Avoid large UI redesigns in this task.
- Keep each edited file under 500 lines.

## Acceptance Criteria

- Reused API-to-view mapping lives in a shared frontend module.
- Admin screens display safe public identifiers where available.
- Sensitive internal identifiers do not appear in user-facing labels unless explicitly intended.
- Frontend build and existing smoke/check commands pass.

## Notes

- This task builds on the ID display policy and T110 sensitive text guard.

## Agent Log

- 2026-04-25: Task created in the board and delivery hardening batch.
