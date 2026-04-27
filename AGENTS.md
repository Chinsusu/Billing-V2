# Agent Rules and Workflow

This file is the required briefing for any coding agent working in this repository.

The project is a Go modular monolith for a VPS/Proxy billing platform. Money, tenant isolation, provider provisioning, credentials, RBAC, and audit behavior are high-risk areas. Do not guess in those areas.

## Behavioral Guidelines

Behavioral guidelines to reduce common LLM coding mistakes.
Merge with project-specific instructions as needed.

Core principle: optimize for **correct**, then **verifiable**, then **minimal** changes.
When these conflict, correctness wins, then honest verification, then minimalism. A small change that is wrong or unverified is worse than a slightly larger one that is right and proven.

Prefer caution over speed for non-trivial work, but do not block on ambiguity that can be resolved safely from the codebase.

### 1. Understand Before Coding

Do not assume silently. Do not hide uncertainty. Do not ask questions the codebase can answer.

Before implementing:

- Inspect the smallest relevant context: caller, callee, tests, types, config, and existing patterns.
- State assumptions when they materially affect the solution.
- If multiple interpretations exist, present the tradeoff.
- If ambiguity affects product behavior, data integrity, security, public APIs, or user-visible behavior, ask before changing.
- If a reasonable low-risk assumption can unblock the task, state it briefly and proceed.
- When both rules above could apply, ask if the decision is hard to reverse; assume if reverting is cheap.

Push back when the requested approach is more complex, risky, or unnecessary than a likely simpler alternative. Disagreement, stated respectfully and with reasoning, is part of the job — not a deviation from it. Defaulting to compliance on a bad approach is a failure mode, not politeness.

Avoid:

- Guessing about business rules.
- Inventing requirements.
- Asking the user for details that can be discovered in the repo.

The codebase answers codebase questions. The user answers product decisions.

### 2. Define Success Before Editing

Turn every task into a verifiable goal before touching code.

Examples:

- "Fix the bug" → reproduce the bug, fix it, verify the fix.
- "Add validation" → test invalid inputs, implement validation, verify expected errors.
- "Refactor X" → confirm behavior before and after remains equivalent.
- "Improve performance" → identify baseline, change bottleneck, compare result if possible.
- "Add a feature" → identify expected behavior, implement, add or update tests, verify integration points.

For multi-step tasks, use a brief plan:

1. Inspect relevant code → verify: identify current behavior and affected files.
2. Make minimal change → verify: targeted tests or checks.
3. Clean up only own changes → verify: no unused imports/types and no unrelated diff.

For trivial changes (under ~10 lines, self-contained, no behavior risk), skip the plan and implement directly.

### 3. Honest Verification and Reporting

This is the most consequential rule in this document.

**Do not claim verification unless it was performed.** Writing a test is not the same as running it. Reading code is not the same as executing it. "This should work" is not verification.

Report only what you actually executed. If you wrote tests but did not run them, say so. If you ran a subset, say which subset. If you relied on type-checking instead of runtime tests, say that.

When a test fails:

- Investigate whether the bug is in the implementation or the test before changing either.
- Do not modify a passing assertion to make a failing test pass unless the assertion itself is provably wrong.
- Do not weaken, skip, or delete tests to clear a red build. If a test is wrong, fix it deliberately and explain why.

When verification is blocked (no test runner, missing fixtures, environment unavailable):

- Explain exactly what blocked it.
- Run the next best available check (type-check, lint, dry-run, manual trace).
- State the remaining risk explicitly.

If you skipped something the user asked for, lead the report with that, not with what you did.

### 4. Simplicity First

Write the minimum code that solves the requested problem.

- No features beyond what was asked.
- No abstractions for single-use code.
- No speculative configurability.
- No broad refactors unless explicitly requested.
- No defensive handling for scenarios made impossible by proven invariants.
- No try/catch wrapping unless there is a concrete failure mode and a real recovery action. Catching to log-and-rethrow is noise.
- Do handle external inputs, I/O, network calls, permissions, API boundaries, and untrusted data.
- If the solution becomes large, pause and reassess whether a simpler path exists.

Self-test:

> "Would a senior engineer consider this overcomplicated for the request?"

If yes, simplify before submitting.

### 5. Surgical Changes

Touch only what is necessary.

When editing existing code:

- Do not improve adjacent code, comments, or formatting unless required.
- Do not refactor unrelated code.
- Do not rename existing variables, functions, or files "for clarity" unless asked. Renames break git blame, search history, and reviewer muscle memory.
- Match existing style, naming, structure, and patterns.
- Avoid broad search-and-replace unless every affected location is reviewed.
- Do not change public APIs, schemas, migrations, auth behavior, or error contracts unless required.

**Never mix behavioral and cosmetic changes in the same diff.** If reformatting is needed, do it in a separate commit — before or after the behavioral change, never alongside. A reviewer should be able to read the behavioral diff without filtering noise.

When your changes create orphans:

- Remove imports, variables, functions, or tests made unused *by your changes*.
- Do not remove pre-existing dead code unless it sits in the immediate vicinity of your change and is clearly safe to remove. When in doubt, mention it separately.
- If you notice unrelated issues, mention them as separate notes — do not fix them in this diff.

The test:

> Every changed line should trace directly to the user's request.

### 6. Security and Data Safety

Never weaken safety to make a task easier. Security regressions are not acceptable shortcuts, even temporarily.

- Do not log secrets, tokens, passwords, private keys, PII, or sensitive payloads.
- Do not commit credentials, even temporarily "for testing." Use environment variables, secret stores, or test fixtures.
- Do not bypass authentication, authorization, validation, rate limits, CSRF/CORS, encryption, or certificate checks unless explicitly requested and justified.
- Do not loosen permissions, scopes, or access controls to make a failing call succeed. Investigate why it fails first.
- Treat user input, API responses, files, database records, and environment variables as untrusted.
- Prefer failing safely over silently accepting unsafe states.
- For destructive or irreversible actions (deletes, migrations, force-pushes, production writes), confirm intent before executing.

### 7. Communication

Be concise but transparent.

Before non-trivial changes:

- State the goal.
- State important assumptions.
- State the verification approach.

After changes, in this order:

1. Anything you skipped, could not do, or left incomplete.
2. What changed.
3. Files touched.
4. Tests/checks actually run, and their results.
5. Anything not verified, with the reason.
6. Unrelated issues noticed, as separate notes.

Do not bury caveats. Do not pad with summary of the user's own request. Do not claim completeness when parts were skipped.

### How To Know This Is Working

These guidelines are succeeding if:

- Diffs are smaller and trace cleanly to the request.
- Fewer unrelated changes appear.
- Fewer rewrites happen due to overcomplication.
- Clarifying questions happen before risky implementation, not after.
- Verification claims are explicit, accurate, and bounded.
- Skipped or incomplete work is reported honestly and up front.

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
docs/05_development_standards/64_Public_Display_ID_And_Backend_Reference_Policy.md
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
docs/05_development_standards/64_Public_Display_ID_And_Backend_Reference_Policy.md
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
- Follow `docs/05_development_standards/64_Public_Display_ID_And_Backend_Reference_Policy.md` for visible IDs.
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
- Use public numeric IDs for visible resource labels; do not render backend UUID references as UI labels.
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
