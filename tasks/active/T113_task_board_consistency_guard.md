# T113 - Task board consistency guard

Status: IN_PROGRESS
Owner: Codex
Branch: codex/t113-task-board-consistency-guard
PR: -
Risk: workflow/CI
Created: 2026-04-25
Updated: 2026-04-25

## Summary

Add a lightweight guard that checks task board metadata for obvious drift after parallel agent updates.

## Scope

- Check active and removed task files for valid status values, IDs, and required metadata fields.
- Check that `TASKS.md` board snapshot counts match task file statuses.
- Check that `TASKS.md` claimable rows only list `TODO` task files.
- Add a documented command agents can run before PR creation.
- Keep each edited file under 500 lines.

## Acceptance Criteria

- Guard fails when a `TODO`, `IN_PROGRESS`, `REVIEW`, `BLOCKED`, `DONE`, or `REMOVED` count is inconsistent with task files.
- Guard fails when a claimable row points to a non-`TODO` task.
- Existing build/test/contract checks pass.
- Docs explain when to run the guard and how to fix failures.

## Notes

- This task should not change product behavior.
- Keep parsing simple and predictable; task files are small Markdown records.

## Agent Log

- 2026-04-25: Task created in the board and delivery hardening batch.
- 2026-04-25: Codex claimed the task; adding a task board consistency guard and docs.
