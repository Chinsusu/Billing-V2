# Validation Command Matrix

**Scope:** One place for local and CI validation commands by change type.

Use this before opening a PR. Add the commands you actually ran to the task log and PR body.

## Ordering Rules

- Run commands from the repo root unless a row says `cd frontend`.
- Run focused tests first when they exist, then broader repo checks.
- Do not run `npm run build`, `npm run smoke:admin`, or `npm run smoke:admin:ci` in parallel. They read or write `.next`.
- For frontend CI smoke, run `npm run build` before `npm run smoke:admin:ci`.
- If `make` is not available on Windows, run the Go equivalent in this doc.
- Do not use production DSNs, provider accounts, or real customer data for local smoke.

## Command Reference

| Purpose | Preferred command | Windows/no-make equivalent | When to run |
| --- | --- | --- | --- |
| Go format | `make fmt` | `go fmt ./...` | Go code changed. |
| Go tests | `make test` | `go test ./...` | Backend, DB, provider, shared code, or full-stack changed. |
| Go build | `make build` | `go build ./cmd/api ./cmd/migrate ./cmd/seed ./cmd/smoke ./cmd/worker` | Backend entrypoints or shared Go packages changed. |
| Migration syntax | `make migrate-validate` | `go run ./cmd/migrate validate` | Migration files or migrator changed. |
| API contract guard | `make contract-guard` | `go run ./cmd/contractguard` | Backend route, permission, query, response shape, or API docs changed. |
| Error code guard | `make error-code-guard` | `go run ./cmd/errorcodeguard` | Error envelope, public error code, handler error behavior, or error docs changed. |
| Task board guard | `make task-guard` | `go run ./cmd/taskguard` | `TASKS.md`, task files, task batches, or board cleanup changed. |
| Frontend install | `cd frontend && npm ci` | same | Fresh checkout, dependency change, or CI parity check. |
| Frontend dependency audit | `cd frontend && npm audit --omit=dev` | same | Dependency or lockfile changed; useful before frontend PRs. |
| Frontend sensitive text guard | `cd frontend && npm run check:sensitive-text` | same | Frontend source, mocks, API mapping, or smoke scripts changed. |
| Frontend lint | `cd frontend && npm run lint` | same | Frontend source changed. |
| Frontend build | `cd frontend && npm run build` | same | Frontend source, config, dependency, or app shell changed. |
| Frontend admin browser smoke | `cd frontend && npm run smoke:admin` | same | Admin navigation, screens, API adapter, mock data, or response mapping changed. |
| Frontend CI browser smoke | `cd frontend && npm run smoke:admin:ci` | same, after `npm run build` | CI parity or standalone output changed. |
| Install Playwright Chromium | `cd frontend && npx playwright install chromium` | same | Browser smoke fails because Chromium is missing locally. |
| DB smoke | `make smoke-dev-db` | `go run ./cmd/smoke dev-db` | Migration, seed, DB lifecycle, or billing seed data changed. |
| API smoke | `make smoke-dev-api` | `go run ./cmd/smoke dev-api` | API read paths or seeded billing API behavior changed and local API is running. |
| Billing mutation smoke | `make smoke-dev-billing` | `go run ./cmd/smoke dev-billing` | Checkout, wallet, payment, order finalization, job creation, provisioning, or service activation changed. |
| Provider sandbox contract | n/a | `go test ./internal/modules/provider -run SandboxContract` | Provider adapter behavior or provider sandbox readiness changed. |
| Whitespace check | n/a | `git diff --check` | Every PR before commit or review. |

## Matrix By Change Type

| Change type | Required local validation | Add when touched |
| --- | --- | --- |
| Docs only | `git diff --check` | `go run ./cmd/taskguard` for task docs; `go run ./cmd/contractguard` or `go run ./cmd/errorcodeguard` if API/error docs changed. |
| Task board or task batch | `go run ./cmd/taskguard`, `git diff --check` | None unless code changed too. |
| Backend API or service | `go fmt ./...`, `go test ./...`, `go build ./cmd/api ./cmd/migrate ./cmd/seed ./cmd/smoke ./cmd/worker`, `git diff --check` | `go run ./cmd/contractguard`, `go run ./cmd/errorcodeguard`, smoke command matching the flow. |
| Shared Go platform | `go fmt ./...`, `go test ./...`, Go build command above, `git diff --check` | API and error guards if HTTP response, middleware, or routing changed. |
| DB migration or seed | `go run ./cmd/migrate validate`, `go test ./cmd/migrate ./internal/platform/db ./internal/seed`, `git diff --check` | `go run ./cmd/smoke dev-db` on local DB; `go run ./cmd/smoke dev-billing` for billing seed flows. |
| Provider adapter | `go test ./internal/modules/provider`, `go test ./internal/modules/provider -run SandboxContract`, `go test ./...`, `git diff --check` | Billing smoke if provisioning outcome changes; update provider sandbox docs. |
| Frontend UI or app shell | `cd frontend && npm run check:sensitive-text`, `npm run lint`, `npm run build`, `git diff --check` | `npm run smoke:admin` for admin screens/nav/API mapping; `npm audit --omit=dev` for dependency changes. |
| Frontend CI or standalone smoke | `cd frontend && npm run check:sensitive-text`, `npm run lint`, `npm run build`, `npm run smoke:admin:ci`, `git diff --check` | `npx playwright install chromium` when local browser runtime is missing. |
| Full-stack billing flow | Backend API/service commands, frontend commands, `go run ./cmd/smoke dev-billing`, `git diff --check` | `go run ./cmd/smoke dev-db` when DB setup changed. |
| CI or workflow | Local equivalent of the changed job plus `git diff --check` | Run the exact frontend/backend commands the workflow invokes. |

## Smoke Prerequisites

`dev-db` requires a local or approved sandbox PostgreSQL database and `DB_DSN`.

`dev-api` requires the API to be running against the same seeded database. Use `-base-url` when it is not on `http://localhost:8080`.

`dev-billing` requires the same database for the API, smoke command, and local fake-provider worker. It must never point at production.

Frontend browser smoke uses mock/intercepted data and does not need a backend or provider account.

## PR Validation Note Template

Use a short list, not a paragraph:

```text
Validation:
- go test ./...
- go run ./cmd/contractguard
- cd frontend && npm run build
- cd frontend && npm run smoke:admin:ci
- git diff --check
```

If a command was not run, write the reason and the risk.

## Related Docs

- Git workflow: `docs/05_development_standards/47_Git_Workflow_Build_Test_PR_Merge_Guide.md`
- Testing strategy: `docs/05_development_standards/49_Testing_Strategy_And_Quality_Gates.md`
- Database workflow: `docs/05_development_standards/52_Database_Migration_Seed_Data_Workflow.md`
- Frontend standard: `docs/05_development_standards/53_Frontend_App_Shell_And_UI_Implementation_Standard.md`
- Local runbook: `docs/05_development_standards/55_Local_Development_Runbook.md`
- API contract guard: `docs/05_development_standards/59_API_Contract_Drift_Guard.md`
- Provider sandbox: `docs/05_development_standards/60_Provider_Sandbox_Contract_Checklist.md`
- Task board guard: `docs/05_development_standards/61_Task_Board_Consistency_Guard.md`
- API error guard: `docs/05_development_standards/62_API_Error_Code_Drift_Guard.md`
