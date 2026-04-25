# Task Board Consistency Guard

**Scope:** Lightweight guard for task board metadata, task-file schema, and claimable/in-flight task rows.

## Command

Run:

```bash
make task-guard
```

This runs:

```bash
go run ./cmd/taskguard
```

On Windows hosts without `make`, run the Go command directly.

## What It Checks

The guard checks:

- task files under `tasks/active/` and `tasks/removed/` use valid `Txxx_*.md` names;
- task file headings match the task ID in the file name;
- required fields exist: `Status`, `Owner`, `Branch`, `PR`, `Risk`, `Created`, and `Updated`;
- required sections exist: `Summary`, `Scope`, `Acceptance Criteria`, `Notes`, and `Agent Log`;
- active task files do not use `REMOVED`;
- removed task files use `REMOVED`;
- `TASKS.md` board snapshot counts match task-file statuses;
- claimable rows point only to active `TODO` task files;
- in-flight rows point only to active `IN_PROGRESS` or `REVIEW` task files.

## When To Run

Run the guard when a PR intentionally changes the task registry, creates a new task batch, archives tasks, or cleans the board snapshot.

Normal task claim, review, and done updates should still edit only the owned task file unless the task explicitly includes a board cleanup. This keeps multi-agent conflicts low.

## Failure Handling

If `make task-guard` fails:

1. Read the failure line; it names the mismatched task, row, or snapshot count.
2. Fix the task file if the task metadata is wrong.
3. Fix `TASKS.md` if the high-level snapshot or claimable/in-flight rows are stale.
4. Re-run the guard before opening the PR.

Do not resolve a guard failure by deleting another agent's task row or changing another agent's task file without a handoff.
