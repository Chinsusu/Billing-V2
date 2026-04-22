# Shared Task Board

This file is the shared coordination board for coding agents.

Agents must update this file when they claim a task, open a PR, or finish work. A task is only done after its PR is merged into `main`.

## Status Rules

Use one of these statuses:

```text
TODO         ready to claim
IN_PROGRESS claimed and being worked on
REVIEW       PR is open and waiting for review/merge
BLOCKED      cannot continue without a decision
DONE         PR merged into main
```

Rules:

- Do not mark a task `DONE` before the PR is merged.
- Do not take a task already marked `IN_PROGRESS` unless the owner says it is abandoned.
- Keep edits small to avoid conflicts when multiple agents update this file.
- If a task touches money, tenant, RBAC, credential, provider, provisioning, audit, or migration, write that in `Risk`.
- If a task creates follow-up work, add a new task instead of hiding it in PR comments.

## How To Claim A Task

1. Pull latest `main`.
2. Pick a task with `TODO`.
3. Change `Status` to `IN_PROGRESS`.
4. Fill `Owner` and `Branch`.
5. Commit the task board update with the code PR, or make a small docs commit first if coordination is needed.

Example:

```text
| T001 | IN_PROGRESS | agent-name | chore/ci-workflow | - | CI | Add GitHub Actions |
```

## Active Tasks

| ID | Status | Owner | Branch | PR | Risk | Task |
| --- | --- | --- | --- | --- | --- | --- |
| T003 | TODO | - | feat/http-middleware-base | - | API/logging | Add recover middleware, request logging middleware, method guard helper, and tests. |
| T004 | TODO | - | feat/identity-tenant-rbac-skeleton | - | tenant/RBAC | Add skeleton interfaces/types for identity, tenant context, and RBAC checks without persistence. |
| T005 | TODO | - | chore/initial-db-migrations | - | migration/tenant | Add initial migration files for tenants, users, roles, permissions, and audit shell after DB skeleton exists. |
| T006 | TODO | - | feat/outbox-job-skeleton | - | worker/migration | Add outbox/jobs table model and worker claim interface after DB skeleton exists. |
| T007 | TODO | - | feat/provider-adapter-interface | - | provider/credential | Add provider adapter interface, normalized provider error types, and fake adapter tests. |
| T008 | TODO | - | docs/local-dev-runbook | - | docs | Add local development runbook after DB and migration commands exist. |
| T009 | TODO | - | feat/frontend-app-shell | - | frontend | Build a runnable Next.js/React/TypeScript frontend app shell with package scripts, working navigation, screen registry, mock data, and build validation. Node is toolchain only; static HTML alone is not accepted. |

## Completed Tasks

| ID | Status | Owner | Branch | PR | Risk | Task |
| --- | --- | --- | --- | --- | --- | --- |
| T001 | DONE | Codex | chore/ci-workflow | [#5](https://github.com/Chinsusu/Billing-V2/pull/5) | CI | Add GitHub Actions workflow for `make fmt`, `make test`, `make build`, and basic secret scan. |
| T002 | DONE | Codex | chore/db-migration-skeleton | [#7](https://github.com/Chinsusu/Billing-V2/pull/7) | migration | Add `internal/platform/db` skeleton and migration runner entrypoint without domain tables yet. |
| T000 | DONE | Codex | chore/bootstrap-go-app | [#3](https://github.com/Chinsusu/Billing-V2/pull/3) | API/config/logging | Bootstrap Go API foundation with config, logger, HTTP helpers, health endpoints, Makefile, and tests. |

## Blocked Tasks

Move blocked tasks here only when they cannot continue without a decision.

| ID | Owner | Branch | Blocker | Needed Decision |
| --- | --- | --- | --- | --- |

## Task Template

Use this when adding a new task:

```text
| T010 | TODO | - | type/scope-short-name | - | risk area | Clear task summary. |
```

## Done Checklist

Before moving a task to `DONE`:

- PR is merged into `main`.
- Required tests passed.
- File line limit is respected.
- Docs or `.env.example` are updated if behavior/config changed.
- No secret, `.env`, key, database dump, or customer data was committed.
- Follow-up tasks are added to this file if needed.
