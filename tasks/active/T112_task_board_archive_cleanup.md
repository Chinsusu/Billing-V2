# T112 - Task board archive cleanup

Status: TODO
Owner: -
Branch: codex/t112-task-board-archive-cleanup
PR: -
Risk: workflow/docs
Created: 2026-04-25
Updated: 2026-04-25

## Summary

Clean up task board documentation so active, done, and removed task records are easier for parallel agents to navigate.

## Scope

- Update `TASKS.md` snapshot and active task table after the T107-T112 batch exists.
- Preserve every task file and historical link.
- Keep conflict-safe one-file-per-task rules intact.
- Do not move or delete task files unless the workflow docs explicitly support the chosen archive layout.
- Keep each edited file under 500 lines.

## Acceptance Criteria

- `TASKS.md` clearly shows current TODO tasks and does not imply older DONE tasks are still claimable.
- Removed task history remains visible.
- Workflow rules still tell agents not to edit unrelated task files.
- Existing validation commands pass.

## Notes

- This is documentation/workflow-only.

## Agent Log

- 2026-04-25: Task created in the post-readiness hardening batch.
