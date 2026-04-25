# T109 - API contract drift guard

Status: TODO
Owner: -
Branch: codex/t109-api-contract-drift-guard
PR: -
Risk: docs/API/CI
Created: 2026-04-25
Updated: 2026-04-25

## Summary

Add a lightweight guard that helps detect when backend route changes drift from the operational API reference.

## Scope

- Check that key billing routes in `cmd/api` or route wiring are represented in `docs/05_development_standards/56_Billing_API_Operational_Reference.md`.
- Focus on stable route groups, permission names, response redaction notes, and query names.
- Keep the guard simple and maintainable; avoid generating a full OpenAPI spec in this task.
- Add a documented command for agents and CI to run.
- Keep each edited file under 500 lines.

## Acceptance Criteria

- Guard fails when a tracked backend route is missing from the operational reference.
- Guard output is readable enough for another agent to fix docs quickly.
- Existing build/test commands pass.
- Docs explain how to update the guard when intentional API changes are made.

## Notes

- This task should not change API behavior.

## Agent Log

- 2026-04-25: Task created in the post-readiness hardening batch.
