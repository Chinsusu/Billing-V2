# T187 - Launch gap task backlog

Status: DONE
Owner: Codex
Branch: codex/t187-launch-gap-task-backlog
PR: https://github.com/Chinsusu/Billing-V2/pull/405
Risk: task planning for launch-readiness P0/P1 work and frontend dependency audit unblock
Created: 2026-05-13
Updated: 2026-05-13

## Summary

Create the remaining launch-readiness task backlog so implementation can proceed one task at a time through the normal branch, PR, CI, review, and merge workflow.

## Scope

- Add TODO task files for the remaining MVP/pilot gaps identified from the roadmap, MVP scope, launch checklist, and current code structure.
- Add one row per new TODO task to `TASKS.md`.
- Update the frontend Next.js patch dependency only if required to clear the existing CI audit blocker.
- Do not implement backend, frontend feature, schema, or runtime behavior in this backlog PR.

## Acceptance Criteria

- The task board lists the remaining work as claimable TODO tasks.
- Each new task has scope, acceptance criteria, risk, and validation notes.
- `go run ./cmd/taskguard`, frontend dependency audit, and `git diff --check` pass locally.

## Notes

- Implementation starts from T188 after this backlog task is merged and marked DONE.
- High-risk tasks must stop and ask when product/security behavior is unclear.
- The Next.js patch update is included only because GitHub Actions blocks all PRs on the existing audit finding.

## Agent Log

- 2026-05-13: Task created and claimed by Codex.
- 2026-05-13: Added launch-readiness TODO backlog T188-T205; taskguard and diff check pass locally.
- 2026-05-13: Opened PR #405 for review.
- 2026-05-13: Updated Next.js patch dependency to clear the GitHub Actions frontend audit blocker.
- 2026-05-13: PR #405 passed GitHub Actions and was merged into main.
