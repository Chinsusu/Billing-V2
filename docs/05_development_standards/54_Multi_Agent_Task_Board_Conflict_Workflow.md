# Multi-Agent Task Board Conflict Workflow

**Version:** v1.9
**Date:** 2026-04-22
**Scope:** Multi-agent task coordination, task claiming, PR status updates, conflict reduction, and task completion.

## Problem

When many agents edit the same `TASKS.md` table, normal work creates merge conflicts:

- Agent A claims or finishes task T004.
- Agent B adds task T010.
- Agent C marks T009 in review.
- All three changed the same central table.

The result is unnecessary conflict even when the actual code changes are unrelated.

## Decision

Use one file per active task.

`TASKS.md` remains the shared index, but it is no longer the mutable source of truth for task status. The mutable source of truth is the task file under:

```text
tasks/active/
```

Each agent edits only the file for the task they own.

## File Ownership

Use this ownership model:

```text
TASKS.md
  Registry/index only. Edit when adding a new task or changing task summary.

tasks/README.md
  Task-file workflow and template.

tasks/active/Txxx_short_name.md
  Source of truth for one active task.
```

Do not use one shared table as the primary status store for all active agents.

## Task File Fields

Each active task file must include:

```text
Status
Owner
Branch
PR
Risk
Created
Updated
Summary
Scope
Acceptance Criteria
Notes
Agent Log
```

Keep the field names stable so agents can scan files quickly.

## Claim Flow

Before claiming:

1. Pull latest `main`.
2. Read `TASKS.md`.
3. Open the linked task file.
4. Confirm `Status: TODO`.
5. Check GitHub PRs or remote branches if the branch name already exists.

To claim:

1. Edit only that task file.
2. Set `Status: IN_PROGRESS`.
3. Set `Owner`.
4. Set `Branch`.
5. Update `Updated`.
6. Add one `Agent Log` entry.

Do not edit `TASKS.md` just to claim a task.

## Review Flow

When opening a PR:

1. Edit only that task file.
2. Set `Status: REVIEW`.
3. Set `PR` to the PR link.
4. Add validation commands and results to `Notes` or `Agent Log`.

Do not update unrelated task files.

## Done Flow

After the PR is merged:

1. Pull latest `main`.
2. Edit only that task file.
3. Set `Status: DONE`.
4. Keep the merged PR link.
5. Update `Updated`.
6. Add the merge date and follow-up task IDs if any.

Do not mark a task `DONE` before the PR is merged.

## Adding A New Task

Adding a new task is the main time when `TASKS.md` should change.

Steps:

1. Pick the next task ID.
2. Create `tasks/active/Txxx_short_name.md` from the template in `tasks/README.md`.
3. Add one row in `TASKS.md` with a short summary and link to the task file.
4. Keep details in the task file, not the index row.

If multiple agents add tasks and `TASKS.md` conflicts, preserve all new rows and sort by task ID.

## Blocking A Task

If blocked:

1. Edit only the task file.
2. Set `Status: BLOCKED`.
3. Describe the blocker, needed decision, and safe follow-up work.

Do not hide blockers in PR comments only.

## Rebase Conflict Rules

If `TASKS.md` conflicts:

- Keep all task rows from both sides.
- Do not overwrite a task file status based on the old table.
- Sort active rows by task ID.
- Keep completed historical rows unless the PR intentionally changes them.

If a task file conflicts:

- Only the owning agent should resolve it.
- Re-read the latest task file.
- Apply the smallest status/log update needed.
- Preserve other agent log entries unless they are clearly accidental duplicates.

## Pull Request Expectations

Each PR should mention:

- task ID
- task file path
- status before and after the PR
- validation commands
- follow-up task IDs if any

Example:

```text
Task: T009
Task file: tasks/active/T009_frontend_app_shell.md
Status: IN_PROGRESS -> REVIEW
Validation: npm run build, make test, make build
```

## Anti-Patterns

Do not:

- edit `TASKS.md` for every status update
- put all active task state in one markdown table
- update another agent's task file without handoff
- mark a task done before merge
- resolve a conflict by deleting unrelated task rows
- hide follow-up tasks in PR comments only

## Definition Of Done

This workflow is followed when:

- every active task has a task file
- agents update only their task file for normal status changes
- `TASKS.md` stays mostly stable
- new tasks are added with both an index row and a task file
- rebase conflicts preserve unrelated task rows and task-file status
