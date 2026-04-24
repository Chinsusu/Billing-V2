# Task Files

Task files reduce merge conflicts between coding agents.

`TASKS.md` is only the index. Each active task owns one file under `tasks/active/`, and the assigned agent updates only that file while working.

## Claim Flow

1. Pull latest `main`.
2. Create the task branch from latest `origin/main`, not from the current branch.
3. Open `TASKS.md` and pick a task.
4. Open the linked task file.
5. Read the task file from latest `origin/main`.
6. Claim only if that file still says `Status: TODO`.
7. Edit only that task file:
   - set `Status: IN_PROGRESS`
   - set `Owner`
   - set `Branch`
   - add a short log entry

If the branch was accidentally created from another task branch, stop and recreate it from `origin/main`; cherry-pick only the commits for this task.

If a task says `IN_PROGRESS` or `REVIEW` but the branch or PR no longer exists, treat it as a stale claim. Do not silently reuse that branch name; reset or clean up the task explicitly first.

## Review Flow

When a PR is open:

- set `Status: REVIEW`
- set `PR`
- add validation commands and results

## Done Flow

After the PR is merged into `main`:

- set `Status: DONE`
- keep the merged PR link
- add the completion date
- add any follow-up task IDs if needed

Do not mark a task `DONE` before its PR is merged.

## Conflict Rules

- Do not edit another agent's task file.
- Do not create child branches from another agent's feature/task branch.
- Do not merge or rebase another task branch into your task branch unless the task owner asks for that handoff.
- Do not edit `TASKS.md` for normal claim/review/done updates.
- If a rebase conflicts in `TASKS.md`, preserve all unrelated new task rows.
- If a rebase conflicts in your own task file, re-read the latest file and apply only your current status update.

## Template

```markdown
# T010 - Short task title

Status: TODO
Owner: -
Branch: type/scope-short-name
PR: -
Risk: risk area
Created: YYYY-MM-DD
Updated: YYYY-MM-DD

## Summary

Clear task summary.

## Scope

- In scope item.
- Out of scope item if needed.

## Acceptance Criteria

- Required behavior or deliverable.
- Required tests or validation.

## Notes

- Important coordination notes.

## Agent Log

- YYYY-MM-DD: Task created.
```
