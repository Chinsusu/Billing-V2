# T115 - API error code drift guard

Status: TODO
Owner: -
Branch: codex/t115-api-error-code-drift-guard
PR: -
Risk: API/docs/CI
Created: 2026-04-25
Updated: 2026-04-25

## Summary

Add a guard that keeps stable API error codes and response envelope documentation aligned with backend behavior.

## Scope

- Identify stable backend error codes and response envelope conventions.
- Check that operational API docs include the tracked codes and response fields.
- Add the guard to the existing validation path if it is fast and deterministic.
- Document how to update the guard when intentional API errors are added or renamed.
- Keep each edited file under 500 lines.

## Acceptance Criteria

- Guard fails when a tracked stable error code is missing from docs.
- Guard output points to the missing code or response convention.
- Existing backend tests and contract guards pass.
- No API behavior changes are introduced.

## Notes

- This task is about docs drift, not designing a full error catalog.

## Agent Log

- 2026-04-25: Task created in the board and delivery hardening batch.
