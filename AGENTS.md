# Agent Rules and Workflow

This file is the required briefing for any coding agent working in this repository.

The project is a Go modular monolith for a VPS/Proxy billing platform. Money, tenant isolation, provider provisioning, credentials, RBAC, and audit behavior are high-risk areas. Do not guess in those areas.

## Required Reading Before Coding

Read these files before starting any code change:

```text
README.md
docs/00_README.md
docs/05_development_standards/46_Coding_Standards_Module_Component_Guide.md
docs/05_development_standards/47_Git_Workflow_Build_Test_PR_Merge_Guide.md
docs/05_development_standards/48_Definition_Of_Ready_Done_And_Task_Workflow.md
docs/05_development_standards/49_Testing_Strategy_And_Quality_Gates.md
docs/05_development_standards/50_API_Response_Error_Logging_Standard.md
docs/05_development_standards/56_Billing_API_Operational_Reference.md
docs/05_development_standards/54_Multi_Agent_Task_Board_Conflict_Workflow.md
docs/04_architecture_deep_dive/37_Go_Backend_Architecture_And_Module_Boundaries.md
docs/04_architecture_deep_dive/45_Architecture_Decision_Records_ADR.md
```

If the task touches config, env, secrets, or credentials, also read:

```text
docs/05_development_standards/51_Environment_Config_Secrets_Guide.md
docs/04_architecture_deep_dive/42_Secrets_Credential_Encryption_And_Audit_Architecture.md
```

If the task touches database, migrations, stores, transactions, wallet, ledger, order, or tenant-scoped data, also read:

```text
docs/05_development_standards/52_Database_Migration_Seed_Data_Workflow.md
docs/04_architecture_deep_dive/38_PostgreSQL_Data_Consistency_Transaction_Design.md
docs/02_technical_handoff/15_Database_Schema_And_ERD.md
```

If the task touches workers, scheduler, jobs, retry, or outbox, also read:

```text
docs/04_architecture_deep_dive/39_Async_Worker_Outbox_And_Job_Architecture.md
docs/02_technical_handoff/19_Worker_Queue_And_Cron_Jobs_Spec.md
```

If the task touches provider, provisioning, inventory, capability snapshot, or provider errors, also read:

```text
docs/04_architecture_deep_dive/40_Provider_Adapter_Runtime_And_Error_Taxonomy.md
docs/01_product_foundation/05_Provisioning_Provider_Adapter_Spec.md
docs/01_product_foundation/11_Provisioning_Idempotency_And_Inventory_Locking.md
docs/02_technical_handoff/18_Provider_Adapter_Technical_Spec.md
```

If the task touches tenant, user roles, RBAC, storefront access, or permissions, also read:

```text
docs/04_architecture_deep_dive/41_Tenant_Isolation_RBAC_And_Security_Architecture.md
docs/01_product_foundation/02_Tenant_Model_Role_Architecture.md
docs/01_product_foundation/10_Tenant_Security_Access_Control_Spec.md
docs/02_technical_handoff/17_RBAC_Permission_Matrix.md
```

If the task touches frontend or API contract, also read:

```text
docs/02_technical_handoff/16_API_Contract_And_Permission_Spec.md
docs/02_technical_handoff/20_UI_Wireflow_And_Screen_Spec.md
docs/05_development_standards/50_API_Response_Error_Logging_Standard.md
docs/05_development_standards/56_Billing_API_Operational_Reference.md
docs/05_development_standards/53_Frontend_App_Shell_And_UI_Implementation_Standard.md
```

## Start-of-Task Checklist

Before editing files, summarize for yourself:

- What is the task goal?
- Which module owns the behavior?
- Which files are expected to change?
- What acceptance criteria apply?
- What tests are required?
- Does the task touch money, tenant, RBAC, credentials, provider, provisioning, audit, or migration?

If any high-risk behavior is unclear, stop and ask for clarification before coding.

## Shared Task Board

Use `TASKS.md` as the shared task index, and use one file per task under `tasks/active/` as the mutable status record.

Rules:

- Read `TASKS.md` before starting work.
- Open the linked file under `tasks/active/` before claiming.
- Fetch `origin/main` and create the task branch from `origin/main`; never branch from another agent's task branch.
- If you are already on a feature/task branch, create a separate worktree or switch to `main` before creating a new branch.
- If a branch was created from the wrong base, stop and recreate it from `origin/main`, then cherry-pick only the intended commits.
- Claim a `TODO` task by editing only that task file: set `Status: IN_PROGRESS`, `Owner`, `Branch`, and an `Agent Log` entry.
- Do not edit `TASKS.md` just to claim, review, block, or finish an existing task.
- Do not claim a task already marked `IN_PROGRESS` or `REVIEW` unless the owner says it is abandoned.
- Change status to `REVIEW` in the task file when the PR is open.
- Change status to `DONE` in the task file only after the PR is merged into `main`.
- Add follow-up tasks by creating a new `tasks/active/Txxx_*.md` file and adding one row to `TASKS.md`.
- If `TASKS.md` conflicts, keep all unrelated task rows and preserve task-file status.

## Repository Structure Rules

Use the existing layout:

```text
cmd/          process entrypoints only
internal/app application wiring only
internal/platform shared infrastructure only
internal/modules business modules
migrations/  PostgreSQL migrations
scripts/     local/dev/ops scripts
docs/        project documentation
```

Rules:

- `cmd/*` starts processes and calls `internal/app`.
- `internal/app` wires dependencies and does not contain business rules.
- `internal/platform/*` must not import `internal/modules/*`.
- Business rules live in `internal/modules/<module>`.
- HTTP handlers parse input and map output; services own business rules; stores own SQL.
- Do not put business rules in scripts, handlers, worker loops, scheduler loops, or entrypoints.

## Code Rules

- Keep every file under 500 lines.
- At 350 lines, consider splitting before merge.
- Do not create packages or files named `common`, `utils`, `helpers`, `misc`, `base`, or `core`.
- Do not copy shared logic across modules.
- Extract shared logic when used in 3 places, or in 2 places if it touches money, tenant, RBAC, credentials, provider, provisioning, audit, idempotency, or rate limit.
- Use clear names. Avoid vague names like `manager`, `processor`, `data`, `item`, `obj`, or `handle`.
- Prefer simple code over clever abstractions.
- Add comments only when they explain non-obvious decisions or high-risk behavior.

## API, Error, and Logging Rules

- Follow `docs/05_development_standards/50_API_Response_Error_Logging_Standard.md`.
- Use one response envelope across modules.
- Error responses must include stable `code`, readable `message`, and `request_id`.
- Validation errors must include field-level details.
- Logs must be structured when possible.
- Logs must include `request_id` and `tenant_id` when relevant.
- Never log secrets, tokens, passwords, provider credentials, private keys, cookies, authorization headers, or raw provider responses that may contain secrets.

## Database Rules

- All schema changes go through migrations.
- Do not edit migrations that have already run in a shared environment.
- Tenant-scoped tables need `tenant_id` unless explicitly documented as global.
- Ledger-style money records are append-only.
- Do not update historical ledger amounts.
- Use reversal or adjustment entries for money corrections.
- Migration PRs must explain data impact and rollback.

## Testing Rules

- Add or update tests required by `docs/05_development_standards/49_Testing_Strategy_And_Quality_Gates.md`.
- Use `docs/05_development_standards/63_Validation_Command_Matrix.md` to choose validation commands for each change type.
- Run relevant tests before PR.
- Run full tests when changing shared platform code.
- Test money flows for idempotency and double-debit prevention.
- Test tenant/RBAC flows for allowed and denied access.
- Test provider/provisioning flows for success, fail, timeout, partial success, retry, and manual review where relevant.

## Frontend Rules

- Do not submit only a static HTML file for frontend app-shell tasks.
- Frontend work must follow `docs/05_development_standards/53_Frontend_App_Shell_And_UI_Implementation_Standard.md`.
- Default frontend stack is Next.js App Router, React, TypeScript, and Tailwind CSS.
- Node.js is only the frontend runtime/toolchain for install, dev, build, and preview/start.
- Do not create a Node.js backend, Express/Nest/Fastify service, Next API route, or Next Server Action for Billing business logic.
- A frontend app shell must have `frontend/package.json` with `dev`, `build`, and `preview` scripts.
- It must have working navigation, screen registry, layout shell, mock data layer, and build validation.
- Do not wire production backend routes during app-shell phase unless the task explicitly asks for it.
- Run the frontend build command and include the result in the PR.

## Git Workflow

Use short-lived branches from latest `origin/main`. Do not create a new task branch while sitting on another task branch.

```bash
git fetch origin --prune
git switch -c <type>/<scope>-<short-name> origin/main
```

For parallel agents, prefer isolated worktrees:

```bash
git fetch origin --prune
git worktree add -b <type>/<scope>-<short-name> /tmp/Billing-<task-id> origin/main
```

Before opening PR, update from `main` without pulling another task branch into your branch:

```bash
git fetch origin --prune
git rebase origin/main
```

If the PR diff contains commits or files from another task, close/recreate the branch from `origin/main` instead of trying to merge it as-is.

Use commit messages like:

```text
feat(wallet): add topup service
fix(provider): classify timeout as retryable
docs(workflow): add agent rules
test(ledger): cover reversal cases
```

Open a PR for every change. Do not push directly to `main`.

## Pull Request Rules

Before opening PR:

- Check `git diff`.
- Ensure the diff is scoped.
- Ensure no file exceeds 500 lines.
- Ensure no secret or `.env` file is committed.
- Run required build/test commands.
- Update docs if behavior, API, config, migration, or workflow changed.

PR body must include:

- What changed.
- Why it changed.
- Tests run.
- Risks.
- Rollback notes for risky changes.

## Stop and Ask

Stop and ask before coding if:

- Acceptance criteria are missing.
- The task touches wallet, ledger, payment, settlement, tenant isolation, RBAC, credential handling, provider provisioning, audit, or migration and the expected behavior is unclear.
- A docs rule conflicts with the requested implementation.
- The change requires breaking the 500-line rule.
- The change requires introducing a new shared abstraction without a clear owner.
- The change may expose secrets or customer data.

## Definition of Done

A task is not done until:

- Code is committed on a branch.
- PR is open or merged according to the workflow.
- Required tests pass or blockers are documented.
- Docs are updated if behavior changed.
- No file exceeds 500 lines.
- No secrets or customer data are committed.
- Review comments are addressed.

## Bootstrap Recommendation

For the first backend implementation PR, keep scope small:

```text
go.mod
Makefile
internal/platform/config
internal/platform/logger
internal/platform/httpserver
internal/app
cmd/api
health endpoint
basic tests
```

Do not start wallet, ledger, order, provider, or provisioning implementation until the base app structure and standards are in place.
