# T008 - Local Development Runbook

Status: TODO
Owner: -
Branch: docs/local-dev-runbook
PR: -
Risk: docs
Created: 2026-04-22
Updated: 2026-04-22

## Summary

Add local development runbook after DB and migration commands exist.

## Scope

- Document local setup commands.
- Document test/build commands.
- Document migration runner usage after it is ready.
- Avoid fake secrets or production credentials.

## Acceptance Criteria

- A new developer can run the backend locally from the documented steps.
- Commands match the repository scripts.
- Secret/config handling follows `.env.example` and config docs.
- `make test` passes.
- `make build` passes.

## Notes

- This task should wait until the DB and migration commands are stable enough to document.

## Agent Log

- 2026-04-22: Task file created from `TASKS.md`.
